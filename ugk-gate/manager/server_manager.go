package manager

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	config2 "github.com/jzyong/ugk/common/config"
	"github.com/jzyong/ugk/common/mode"
	"github.com/jzyong/ugk/gate/config"
	"github.com/jzyong/ugk/message/message"
	"github.com/xtaci/kcp-go/v5"
	"google.golang.org/protobuf/proto"
	"sync"
	"time"
)

// ServerManager 服务器-网络
type ServerManager struct {
	util.DefaultModule
	GameClients map[string]map[uint32]*GameKcpClient //游戏客户端 key=游戏名称 ==》游戏id
	mutex       sync.RWMutex
}

var serverManager *ServerManager
var serverSingletonOnce sync.Once

func GetServerManager() *ServerManager {
	serverSingletonOnce.Do(func() {
		serverManager = &ServerManager{
			GameClients: make(map[string]map[uint32]*GameKcpClient, 10),
		}
	})
	return serverManager
}

// 消息执行函数
type serverHandFunc func(user *User, client *GameKcpClient, msg *mode.UgkMessage)

// ServerHandlers 客户端消息处理器
var ServerHandlers = make(map[uint32]serverHandFunc)

func (m *ServerManager) Init() error {
	log.Info("ServerManager 初始化......")
	go m.runKcpServer()
	return nil
}

// 启动kcp服务器
func (m *ServerManager) runKcpServer() {
	url := fmt.Sprintf("%v:%v", "0.0.0.0", config.BaseConfig.GamePort)
	log.Info("游戏udp监听地址：%s", url)
	if listener, err := kcp.ListenWithOptions(url, nil, 0, 0); err == nil {
		for {
			s, err := listener.AcceptKCP()
			if err != nil {
				log.Error("kcp启动失败：%v", err)
			}

			//kcp文档 https://wetest-qq-com.translate.goog/labs/391?_x_tr_sl=auto&_x_tr_tl=en&_x_tr_hl=zh-CN
			//设置参数 https://github.com/skywind3000/kcp/blob/master/README.en.md#protocol-configuration
			//UDPSession mtu最大限制为1500，发送消息大于1500字节kcp底层默认分为几段进行消息发送（标识每段frg=0），
			//但是接收端每次只能读取1段（因为每段frg=0）， 需要自己截取几段字节流封装
			s.SetMtu(config2.MTU)
			s.SetWindowSize(config2.WindowSize, config2.WindowSize)
			s.SetReadBuffer(8 * 1024 * 1024)
			s.SetWriteBuffer(16 * 1024 * 1024)
			s.SetStreamMode(true) //true 流模式：使每个段数据填充满,避免浪费; false 消息模式 每个消息一个数据段
			//nodelay : Whether nodelay mode is enabled, 0 is not enabled; 1 enabled.
			//interval ：Protocol internal work interval, in milliseconds, such as 10 ms or 20 ms.
			//resend ：Fast retransmission mode, 0 represents off by default, 2 can be set (2 ACK spans will result in direct retransmission)
			//nc ：Whether to turn off flow control, 0 represents “Do not turn off” by default, 1 represents “Turn off”.
			//Normal Mode: ikcp_nodelay(kcp, 0, 40, 0, 0);
			//Turbo Mode： ikcp_nodelay(kcp, 1, 10, 2, 1);
			//s.SetNoDelay(0, 40, 0, 0)
			s.SetNoDelay(1, 10, 2, 1)
			client := gameChannelActive(s)
			go gameChannelRead(client)

		}
	} else {
		log.Error("kcp启动失败：%v", err)
	}
}

// UpdateGameServer 更新游戏服务器信息
func (m *ServerManager) UpdateGameServer(serverInfo *message.ServerInfo, client *GameKcpClient) {
	defer m.mutex.Unlock()
	m.mutex.Lock()
	if serverList, ok := m.GameClients[serverInfo.GetName()]; ok {
		if c, ok2 := serverList[serverInfo.GetId()]; ok2 {
			c.Id = serverInfo.GetId()
			c.Name = serverInfo.GetName()
		} else {
			serverList[serverInfo.GetId()] = client
			m.GameClients[serverInfo.GetName()] = serverList
			log.Info("后端服务：%s-%d 注册", client.Name, client.Id)
		}
	} else {
		serverList = make(map[uint32]*GameKcpClient, 2)
		serverList[serverInfo.GetId()] = client
		m.GameClients[serverInfo.GetName()] = serverList
		log.Info("后端服务：%s-%d 注册", client.Name, client.Id)
	}
}

// 移除game服务器连接
func (m *ServerManager) removeGameServer(client *GameKcpClient) {
	defer m.mutex.Unlock()
	m.mutex.Lock()
	if serverList, ok := m.GameClients[client.Name]; ok {
		if _, ok2 := serverList[client.Id]; ok2 {
			delete(serverList, client.Id)
			m.GameClients[client.Name] = serverList
			log.Info("后端服务：%s-%d 移除", client.Name, client.Id)
		}
	}
}

// AssignLobby 分配大厅 TODO 一致性hash？
func (m *ServerManager) AssignLobby(playerId int64) *GameKcpClient {
	//TODO 暂时只有一个，随机
	if serverList, ok := m.GameClients["lobby"]; ok {
		if len(serverList) > 0 {
			return serverList[1]
		}
		return nil
	} else {
		return nil
	}
}

// 连接激活
func gameChannelActive(session *kcp.UDPSession) *GameKcpClient {
	log.Info("%s 连接创建", session.RemoteAddr().String())
	return NewGameKcpClient(session)
}

// 连接关闭
// 客户端强制杀进程，服务器不知道连接断开。kcp-go源码没有示例,因此使用自定义心跳（每2s请求一次心跳，超过10s断开连接）
func gameChannelInactive(client *GameKcpClient, err error) {
	log.Info("%d - %s 连接关闭:%s", client.Id, client.UdpSession.RemoteAddr(), err)
	// 移除客户端对象
	GetServerManager().removeGameServer(client)
	close(client.CloseChan)
}

// 处理接收消息
// UDPSession Read和Write都可能阻塞，共用routine是否会被阻塞自定义逻辑？
func gameChannelRead(client *GameKcpClient) {
	session := client.UdpSession
	for client.State != Closed {
		//会阻塞
		//最多读取mtu - kcp消息头 字节
		//UDPSession mtu最大限制为1500，发送消息大于1500字节kcp底层默认分为几段进行消息发送（标识每段frg=0），
		//但是接收端每次只能读取1段（因为每段frg=0）， 需要自己截取几段字节流封装
		n, err := session.Read(client.ReceiveReadCache)
		if err != nil {
			log.Error("kcp启动失败：%v", err)
			gameChannelInactive(client, err)
			return
		}
		client.ReceiveBuffer.Write(client.ReceiveReadCache[0:n])

		// 转发消息到User routine
		// 通过比较n和 length循环获取批量消息包
		//`消息长度4+消息id4+序列号4+时间戳8+protobuf消息体`
		receiveBytes := client.ReceiveBuffer.Bytes()
		client.ReceiveBuffer.Reset()
		remainBytes := len(receiveBytes)
		index := 0
		//解析批量消息包
		for remainBytes > 0 {
			//小端
			length := int(uint32(receiveBytes[index]) | uint32(receiveBytes[index+1])<<8 | uint32(receiveBytes[index+2])<<16 | uint32(receiveBytes[index+3])<<24)
			length += 4 //客户端请求长度不包含自身
			if length > config2.MessageLimit {
				gameChannelInactive(client, errors.New(fmt.Sprintf("消息太长")))
				return
			}
			//消息不够,缓存下次使用
			if length > remainBytes {
				client.ReceiveBuffer.Write(receiveBytes[index:])
				break
			}

			//packetData := make([]byte, length)
			packetData := mode.GetBytes()[:length]
			copy(packetData, receiveBytes[index:index+length])
			client.ReceiveBytes <- packetData
			remainBytes = remainBytes - length
			index += length
			//log.Info("收到消息：读取长度=%v 消息长度=%v 剩余长度=%v", n, length, remainBytes)
		}

	}
}

func (m *ServerManager) Run() {
}

func (m *ServerManager) Stop() {
}

// 消息处理函数
type messageHandFunc func(playerId int64, messageId uint32, seq uint32, timeStamp int64, data []byte, client *GameKcpClient)

// NetState 用户状态
type NetState int

const (
	Active    NetState = 0 //网络激活
	Connected NetState = 1 //已登录
	Closed    NetState = 2 //已关闭
)

// GameKcpClient 后端游戏kcp客户端
type GameKcpClient struct {
	Id               uint32          //唯一id
	Name             string          //名称
	UdpSession       *kcp.UDPSession //客户端连接会话
	SendBuffer       *bytes.Buffer   //发送缓冲区，单线程调用
	ReceiveBuffer    *bytes.Buffer   //接收缓冲区
	ReceiveBytes     chan []byte     //接收到的消息
	ReceiveReadCache []byte          // 接收端读取Byte缓存
	CloseChan        chan struct{}   //离线等关闭Chan
	State            NetState        //用户状态
	HeartTime        time.Time       //心跳时间
}

// NewGameKcpClient 构建
func NewGameKcpClient(udpSession *kcp.UDPSession) *GameKcpClient {
	client := &GameKcpClient{UdpSession: udpSession,
		SendBuffer:       bytes.NewBuffer([]byte{}),
		ReceiveBuffer:    bytes.NewBuffer([]byte{}),
		ReceiveBytes:     make(chan []byte, 1024),
		ReceiveReadCache: make([]byte, 1500), //每次最多读取1500-消息头字节
		CloseChan:        make(chan struct{}),
		State:            Active,
		HeartTime:        time.Now(),
	}
	////只在此处添加
	//GetServerManager().IpClients[udpSession.RemoteAddr().String()] = client
	go client.run()
	return client
}

// 每个客户端一个routine运行
func (client *GameKcpClient) run() {
	secondTicker := time.Tick(time.Second)
	for {
		select {
		case receiveByte := <-client.ReceiveBytes:
			client.messageDistribute(receiveByte)
		case <-client.CloseChan:
			log.Info("%v %v chan关闭", client.Id, client.UdpSession.RemoteAddr().String())
			client.State = Closed
			return
		case <-secondTicker:
			client.secondUpdate()
		}
	}
}

// 玩家更新逻辑
func (client *GameKcpClient) secondUpdate() {
	// 心跳监测
	if time.Now().Sub(client.HeartTime) > config2.ClientHeartInterval {
		gameChannelInactive(client, errors.New(fmt.Sprintf("心跳超时%f", time.Now().Sub(client.HeartTime).Seconds())))
	}

}

func (client *GameKcpClient) messageDistribute(data []byte) {
	defer mode.ReturnBytes(data)
	client.HeartTime = time.Now()
	//`消息长度4+玩家ID8+消息id4+序列号4+时间戳8+protobuf消息体`
	//截取消息
	//小端
	messageId := uint32(data[12]) | uint32(data[13])<<8 | uint32(data[14])<<16 | uint32(data[15])<<24
	//log.Debug("收到消息：%d", messageId)
	handFunc := ServerHandlers[messageId]
	if handFunc != nil { //本地处理
		dataReader := bytes.NewReader(data)
		var messageLength int32
		if err := binary.Read(dataReader, binary.LittleEndian, &messageLength); err != nil {
			gameChannelInactive(client, errors.New("读取消息长度错误"))
			return
		}
		var playerId int64
		if err := binary.Read(dataReader, binary.LittleEndian, &playerId); err != nil {
			gameChannelInactive(client, errors.New("读取玩家ID错误"))
			return
		}
		if err := binary.Read(dataReader, binary.LittleEndian, &messageId); err != nil {
			gameChannelInactive(client, errors.New("读取消息ID错误"))
			return
		}
		var seq uint32
		if err := binary.Read(dataReader, binary.LittleEndian, &seq); err != nil {
			gameChannelInactive(client, errors.New("读取消息seq错误"))
			return
		}
		var timeStamp int64
		if err := binary.Read(dataReader, binary.LittleEndian, &timeStamp); err != nil {
			gameChannelInactive(client, errors.New("读取消息timeStamp错误"))
			return
		}
		//protoData := make([]byte, messageLength-24)
		protoData := mode.GetBytes()[:messageLength-24]
		if err := binary.Read(dataReader, binary.LittleEndian, &protoData); err != nil {
			gameChannelInactive(client, errors.New("读取消息proto数据错误"))
			return
		}

		ugkMessage := mode.GetUgkMessage()
		ugkMessage.MessageId = messageId
		ugkMessage.Seq = seq
		ugkMessage.TimeStamp = timeStamp
		ugkMessage.Bytes = protoData
		ugkMessage.Client = client

		// 用户消息转发到用户routine
		if playerId > 0 {
			user := GetUserManager().GetUser(playerId)
			if user == nil {
				log.Warn("玩家：%d 已离线，消息%d转发失败", playerId, messageId)
			} else {
				user.GameMessages <- ugkMessage
			}
		} else {
			handFunc(nil, client, ugkMessage)
			mode.ReturnUgkMessage(ugkMessage)
		}

	} else { //转发给用户
		playerId := int64(data[4]) | int64(data[5])<<8 | int64(data[6])<<16 | int64(data[7])<<24 | int64(data[8])<<32 | int64(data[9])<<40 | int64(data[10])<<48 | int64(data[11])<<56
		user := GetUserManager().GetUser(playerId)
		if user == nil {
			log.Warn("玩家：%d 已离线，消息%d转发失败", playerId, messageId)
		} else {
			user.TransmitToClient(data, messageId)
		}
	}
}

// SendToGame 发送消息到游戏
func (client *GameKcpClient) SendToGame(playerId int64, mid message.MID, msg proto.Message, seq uint32) error {
	protoData, err := proto.Marshal(msg)
	if err != nil {
		log.Error("%d - %s 发送消息 %d 失败：%v", client.Id, client.UdpSession.RemoteAddr().String(), mid, err)
		return err
	}
	protoLength := len(protoData)
	if protoLength > config2.MessageLimit {
		log.Error("%d - %s 发送消息 %d 失败：%v", client.Id, client.UdpSession.RemoteAddr().String(), mid, err)
		return errors.New("消息超长")
	}

	//`消息长度4+玩家ID8+消息id4+序列号4+时间戳8+protobuf消息体`
	buffer := bytes.NewBuffer([]byte{})
	//写dataLen 不包含自身长度
	if err := binary.Write(buffer, binary.LittleEndian, uint32(24+protoLength)); err != nil {
		log.Error("%d - %s 发送消息 %d 失败：%v", client.UdpSession, client.UdpSession.RemoteAddr().String(), mid, err)
		return err
	}

	//写玩家ID
	if err := binary.Write(buffer, binary.LittleEndian, playerId); err != nil {
		log.Error("%d - %s 发送消息 %d 失败：%v", playerId, client.UdpSession.RemoteAddr().String(), mid, err)
		return err
	}

	//写msgID
	if err := binary.Write(buffer, binary.LittleEndian, uint32(mid)); err != nil {
		log.Error("%d - %s 发送消息 %d 失败：%v", client.Id, client.UdpSession.RemoteAddr().String(), mid, err)
		return err
	}
	//写 序列号
	if err := binary.Write(buffer, binary.LittleEndian, seq); err != nil {
		log.Error("%d - %s 发送消息 %d 失败：%v", client.Id, client.UdpSession.RemoteAddr().String(), mid, err)
		return err
	}
	//写时间戳
	if err := binary.Write(buffer, binary.LittleEndian, util.CurrentTimeMillisecond()); err != nil {
		log.Error("%d - %s 发送消息 %d 失败：%v", client.Id, client.UdpSession.RemoteAddr().String(), mid, err)
		return err
	}
	//写data数据
	if err := binary.Write(buffer, binary.LittleEndian, protoData); err != nil {
		log.Error("%d - %s 发送消息 %d 失败：%v", client.Id, client.UdpSession.RemoteAddr().String(), mid, err)
		return err
	}
	_, err = client.UdpSession.Write(buffer.Bytes())
	if err != nil {
		log.Error("%d - %s 发送消息 %d 失败：%v", client.Id, client.UdpSession.RemoteAddr().String(), mid, err)
		return err
	}
	return nil
}
