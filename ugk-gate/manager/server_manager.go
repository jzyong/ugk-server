package manager

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	"github.com/jzyong/ugk/common/constant"
	"github.com/jzyong/ugk/gate/config"
	"github.com/xtaci/kcp-go/v5"
	"sync"
	"time"
)

// ServerManager 服务器-网络
type ServerManager struct {
	util.DefaultModule
	GameClients map[string]map[uint16]*GameKcpClient //游戏客户端 key=游戏名称 ==》游戏id
}

var serverManager *ServerManager
var serverSingletonOnce sync.Once

func GetServerManager() *ServerManager {
	serverSingletonOnce.Do(func() {
		serverManager = &ServerManager{
			GameClients: make(map[string]map[uint16]*GameKcpClient, 10),
		}
	})
	return serverManager
}

// 消息执行函数
type serverHandFunc func(user *User, data []byte, seq uint32, client *GameKcpClient)

// ServerHandlers 客户端消息处理器
var ServerHandlers = make(map[uint32]serverHandFunc)

func (m *ServerManager) Init() error {
	log.Info("ServerManager 初始化......")
	go m.runKcpServer()
	return nil
}

// 启动kcp服务器
func (m *ServerManager) runKcpServer() {
	url := fmt.Sprintf("%v:%v", config.AppConfigManager.PrivateIp, config.AppConfigManager.GamePort)
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
			s.SetMtu(constant.MTU)
			s.SetWindowSize(constant.WindowSize, constant.WindowSize)
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

// 连接激活
func gameChannelActive(session *kcp.UDPSession) *GameKcpClient {
	log.Info("%s 连接创建", session.RemoteAddr().String())
	return NewGameKcpClient(session)
}

// 连接关闭
// 客户端强制杀进程，服务器不知道连接断开。kcp-go源码没有示例,因此使用自定义心跳（每2s请求一次心跳，超过10s断开连接）
func gameChannelInactive(client *GameKcpClient, err error) {
	log.Info("%d - %s 连接关闭:%s", client.Id, client.UdpSession.RemoteAddr(), err)
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
			if length > constant.MessageLimit {
				gameChannelInactive(client, errors.New(fmt.Sprintf("消息太长")))
				return
			}
			//消息不够,缓存下次使用
			if length > remainBytes {
				client.ReceiveBuffer.Write(receiveBytes[index:])
				break
			}

			packetData := make([]byte, length)
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
	if time.Now().Sub(client.HeartTime) > constant.ServerHeartInterval {
		gameChannelInactive(client, errors.New(fmt.Sprintf("心跳超时%f", time.Now().Sub(client.HeartTime).Seconds())))
	}

}

func (client *GameKcpClient) messageDistribute(data []byte) {
	client.HeartTime = time.Now()
	//`消息长度4+玩家ID8+消息id4+序列号4+时间戳8+protobuf消息体`
	//截取消息
	//小端
	messageId := uint32(data[12]) | uint32(data[13])<<8 | uint32(data[14])<<16 | uint32(data[15])<<24
	//log.Info("收到消息：%d", messageId)
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
		protoData := make([]byte, messageLength-24)
		if err := binary.Read(dataReader, binary.LittleEndian, &protoData); err != nil {
			gameChannelInactive(client, errors.New("读取消息proto数据错误"))
			return
		}
		//TODO 用户消息转发到用户routine
		handFunc(nil, protoData, seq, client)
	} else { //转发给用户
		//TODO
	}
}

//TODO 发送返回消息
