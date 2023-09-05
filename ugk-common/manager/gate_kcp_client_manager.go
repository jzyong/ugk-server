package manager

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	"github.com/jzyong/ugk/common/constant"
	"github.com/jzyong/ugk/message/message"
	"github.com/xtaci/kcp-go/v5"
	"google.golang.org/protobuf/proto"
	"sync"
	"time"
)

// GateKcpClientManager 网络，kcp连接网关
type GateKcpClientManager struct {
	util.DefaultModule
	IpClients          map[string]*GateKcpClient   //网络连接客户端 key=ip
	MessageHandFunc    messageHandFunc             //消息分发处理函数
	ServerHeartRequest *message.ServerHeartRequest //服务器心跳消息
}

var gateKcpClientManager *GateKcpClientManager
var gateKcpClientSingletonOnce sync.Once

func GetGateKcpClientManager() *GateKcpClientManager {
	gateKcpClientSingletonOnce.Do(func() {
		gateKcpClientManager = &GateKcpClientManager{IpClients: make(map[string]*GateKcpClient, 2)}
	})
	return gateKcpClientManager
}

func (m *GateKcpClientManager) Init() error {
	log.Info("GateKcpClientManager 初始化......")
	return nil
}

// TODO 连接网关，需要服务发现
// TODO 定时发送服务器信息，充当心跳

// ConnectKcpServer 连接kcp服务器
func (m *GateKcpClientManager) ConnectKcpServer(url string) {

	log.Info("连接网关：%s", url)
	if udpSession, err := kcp.DialWithOptions(url, nil, 0, 0); err == nil {
		//kcp文档 https://wetest-qq-com.translate.goog/labs/391?_x_tr_sl=auto&_x_tr_tl=en&_x_tr_hl=zh-CN
		//设置参数 https://github.com/skywind3000/kcp/blob/master/README.en.md#protocol-configuration
		//UDPSession mtu最大限制为1500，发送消息大于1500字节kcp底层默认分为几段进行消息发送（标识每段frg=0），
		//但是接收端每次只能读取1段（因为每段frg=0）， 需要自己截取几段字节流封装
		udpSession.SetMtu(constant.MTU)
		udpSession.SetWindowSize(constant.WindowSize, constant.WindowSize)
		udpSession.SetReadBuffer(8 * 1024 * 1024)
		udpSession.SetWriteBuffer(16 * 1024 * 1024)
		udpSession.SetStreamMode(true) //true 流模式：使每个段数据填充满,避免浪费; false 消息模式 每个消息一个数据段
		//nodelay : Whether nodelay mode is enabled, 0 is not enabled; 1 enabled.
		//interval ：Protocol internal work interval, in milliseconds, such as 10 ms or 20 ms.
		//resend ：Fast retransmission mode, 0 represents off by default, 2 can be set (2 ACK spans will result in direct retransmission)
		//nc ：Whether to turn off flow control, 0 represents “Do not turn off” by default, 1 represents “Turn off”.
		//Normal Mode: ikcp_nodelay(kcp, 0, 40, 0, 0);
		//Turbo Mode： ikcp_nodelay(kcp, 1, 10, 2, 1);
		//s.SetNoDelay(0, 40, 0, 0)
		udpSession.SetNoDelay(1, 10, 2, 1)
		client := channelActive(udpSession, 1, url) //TODO 或者id
		go channelRead(client)

		//for {
		//	data := time.Now().String()
		//	buf := make([]byte, len(data))
		//	log.Println("sent:", data)
		//	if _, err := udpSession.Write([]byte(data)); err == nil {
		//		// read back the data
		//		if _, err := io.ReadFull(udpSession, buf); err == nil {
		//			log.Println("recv:", string(buf))
		//		} else {
		//			log.Fatal(err)
		//		}
		//	} else {
		//		log.Fatal(err)
		//	}
		//	time.Sleep(time.Second)
		//
		//}
	} else {
		log.Error("连接kcp服务器失败：%v", err)
	}
}

// 连接激活
func channelActive(session *kcp.UDPSession, serverId uint32, url string) *GateKcpClient {
	log.Info("%s 连接创建", session.RemoteAddr().String())
	return NewGateKcpClient(session, serverId, url)
}

// 连接关闭
func channelInactive(client *GateKcpClient, err error) {
	log.Info("%d - %s 连接关闭:%s", client.Id, client.UdpSession.RemoteAddr(), err)
	close(client.CloseChan)
}

// 处理接收消息
// UDPSession Read和Write都可能阻塞，共用routine是否会被阻塞自定义逻辑？
func channelRead(client *GateKcpClient) {
	session := client.UdpSession
	for client.State != Offline {
		//会阻塞
		//最多读取mtu - kcp消息头 字节
		//UDPSession mtu最大限制为1500，发送消息大于1500字节kcp底层默认分为几段进行消息发送（标识每段frg=0），
		//但是接收端每次只能读取1段（因为每段frg=0）， 需要自己截取几段字节流封装
		n, err := session.Read(client.ReceiveReadCache)
		if err != nil {
			log.Error("kcp启动失败：%v", err)
			channelInactive(client, err)
			return
		}
		client.ReceiveBuffer.Write(client.ReceiveReadCache[0:n])

		// 转发消息到Client routine
		// 通过比较n和 length循环获取批量消息包
		//`消息长度4+玩家ID8+消息id4+序列号4+时间戳8+protobuf消息体`
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
				channelInactive(client, errors.New(fmt.Sprintf("消息太长:%d", length)))
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

func (m *GateKcpClientManager) Run() {
}

func (m *GateKcpClientManager) Stop() {
}

// 消息处理函数
type messageHandFunc func(playerId int64, messageId uint32, seq uint32, timeStamp int64, data []byte, client *GateKcpClient)

// NetState 用户状态
type NetState int

const (
	NetWorkActive NetState = 0 //网络激活
	Connected     NetState = 1 //已登录
	Offline       NetState = 2 //已离线
)

// GateKcpClient 网关kcp客户端
type GateKcpClient struct {
	Id               uint32          //唯一id
	Url              string          //地址
	UdpSession       *kcp.UDPSession //客户端连接会话
	SendBuffer       *bytes.Buffer   //发送缓冲区，单线程调用
	ReceiveBuffer    *bytes.Buffer   //接收缓冲区
	ReceiveBytes     chan []byte     //接收到的消息
	ReceiveReadCache []byte          // 接收端读取Byte缓存
	CloseChan        chan struct{}   //离线等关闭Chan
	State            NetState        //用户状态
	HeartTime        time.Time       //心跳时间
}

// NewGateKcpClient 构建
func NewGateKcpClient(udpSession *kcp.UDPSession, serverId uint32, url string) *GateKcpClient {
	client := &GateKcpClient{UdpSession: udpSession,
		Id:               serverId,
		Url:              url,
		SendBuffer:       bytes.NewBuffer([]byte{}),
		ReceiveBuffer:    bytes.NewBuffer([]byte{}),
		ReceiveBytes:     make(chan []byte, 1024),
		ReceiveReadCache: make([]byte, 1500), //每次最多读取1500-消息头字节
		CloseChan:        make(chan struct{}),
		State:            NetWorkActive,
		HeartTime:        time.Now(),
	}
	//只在此处添加
	GetGateKcpClientManager().IpClients[udpSession.RemoteAddr().String()] = client
	go client.run()
	return client
}

// 每个客户端一个routine运行
func (client *GateKcpClient) run() {
	secondTicker := time.Tick(time.Second)
	for {
		select {
		case receiveByte := <-client.ReceiveBytes:
			client.messageDistribute(receiveByte)
		case <-client.CloseChan:
			log.Info("%v %v chan关闭", client.Id, client.UdpSession.RemoteAddr().String())
			client.State = Offline
			return
		case <-secondTicker:
			client.secondUpdate()
		}
	}
}

// 玩家更新逻辑
func (client *GateKcpClient) secondUpdate() {
	// 心跳监测
	if time.Now().Sub(client.HeartTime) > constant.ClientHeartInterval {
		channelInactive(client, errors.New(fmt.Sprintf("心跳超时%f", time.Now().Sub(client.HeartTime).Seconds())))
	}
	if GetGateKcpClientManager().ServerHeartRequest != nil {
		client.SendToGate(0, message.MID_ServerHeartReq, GetGateKcpClientManager().ServerHeartRequest, 0)
	}

}

func (client *GateKcpClient) messageDistribute(data []byte) {
	client.HeartTime = time.Now()
	//`消息长度4+玩家ID8+消息id4+序列号4+时间戳8+protobuf消息体`
	//截取消息
	dataReader := bytes.NewReader(data)
	var messageLength int32
	if err := binary.Read(dataReader, binary.LittleEndian, &messageLength); err != nil {
		channelInactive(client, errors.New("读取消息长度错误"))
		return
	}
	var playerId int64
	if err := binary.Read(dataReader, binary.LittleEndian, &playerId); err != nil {
		channelInactive(client, errors.New("读取玩家ID错误"))
		return
	}
	var messageId uint32
	if err := binary.Read(dataReader, binary.LittleEndian, &messageId); err != nil {
		channelInactive(client, errors.New("读取消息ID错误"))
		return
	}
	var seq uint32
	if err := binary.Read(dataReader, binary.LittleEndian, &seq); err != nil {
		channelInactive(client, errors.New("读取消息seq错误"))
		return
	}
	var timeStamp int64
	if err := binary.Read(dataReader, binary.LittleEndian, &timeStamp); err != nil {
		channelInactive(client, errors.New("读取消息timeStamp错误"))
		return
	}
	protoData := make([]byte, messageLength-24)
	if err := binary.Read(dataReader, binary.LittleEndian, &protoData); err != nil {
		channelInactive(client, errors.New("读取消息proto数据错误"))
		return
	}
	if messageId == uint32(message.MID_ServerHeartRes) {
		client.heartRes()
	} else {
		GetGateKcpClientManager().MessageHandFunc(playerId, messageId, seq, timeStamp, protoData, client)
	}
}

func (client *GateKcpClient) heartRes() {
	client.HeartTime = time.Now()
	log.Info("收到心跳返回消息包")
}

// SendToGate 发送消息到网关
func (client *GateKcpClient) SendToGate(playerId int64, mid message.MID, msg proto.Message, seq uint32) error {

	protoData, err := proto.Marshal(msg)
	if err != nil {
		log.Error("%d - %s 发送消息 %d 失败：%v", client.Id, client.UdpSession.RemoteAddr().String(), mid, err)
		return err
	}
	protoLength := len(protoData)
	if protoLength > constant.MessageLimit {
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
