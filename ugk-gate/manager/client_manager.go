package manager

import (
	"errors"
	"fmt"
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	config2 "github.com/jzyong/ugk/common/config"
	"github.com/jzyong/ugk/common/mode"
	"github.com/jzyong/ugk/gate/config"
	"github.com/xtaci/kcp-go/v5"
	"sync"
)

// ClientManager 客户端-网络
type ClientManager struct {
	util.DefaultModule
}

var clientManager *ClientManager
var clientSingletonOnce sync.Once

func GetClientManager() *ClientManager {
	clientSingletonOnce.Do(func() {
		clientManager = &ClientManager{}
	})
	return clientManager
}

// 消息执行函数
type handFunc func(user *User, msg *mode.UgkMessage)

// ClientHandlers 客户端消息处理器
var ClientHandlers = make(map[uint32]handFunc)

func (m *ClientManager) Init() error {
	log.Info("ClientManager 初始化......")
	go m.runKcpServer()
	return nil
}

// 启动kcp服务器
func (m *ClientManager) runKcpServer() {
	url := fmt.Sprintf("%v:%v", "0.0.0.0", config.BaseConfig.ClientPort)
	log.Info("玩家udp监听地址：%s", url)
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
			user := channelActive(s)
			go channelRead(user)

		}
	} else {
		log.Error("kcp启动失败：%v", err)
	}

}

// 连接激活
func channelActive(session *kcp.UDPSession) *User {
	log.Info("%s 连接创建", session.RemoteAddr().String())
	return NewUser(session)
}

// 连接关闭
// 客户端强制杀进程，服务器不知道连接断开。kcp-go源码没有示例,因此使用自定义心跳（每2s请求一次心跳，超过10s断开连接）
func channelInactive(user *User, err error) {
	//不断的创建异常的socket，出现user.CloseChane已关闭问题
	defer func() {
		if r := recover(); r != nil {
			log.Error("关闭 chan异常：%v", r)
		}
	}()
	log.Info("%d - %s 连接关闭:%s", user.Id, user.ClientSession.RemoteAddr(), err)
	close(user.CloseChan)
}

// 处理接收消息
// UDPSession Read和Write都可能阻塞，共用routine是否会被阻塞自定义逻辑？
func channelRead(user *User) {
	session := user.ClientSession
	for user.State != Offline {
		//会阻塞
		//最多读取mtu - kcp消息头 字节
		//UDPSession mtu最大限制为1500，发送消息大于1500字节kcp底层默认分为几段进行消息发送（标识每段frg=0），
		//但是接收端每次只能读取1段（因为每段frg=0）， 需要自己截取几段字节流封装
		n, err := session.Read(user.ReceiveReadCache)
		if err != nil {
			log.Error("kcp读取报错：%v", err)
			channelInactive(user, err)
			return
		}
		user.ReceiveBuffer.Write(user.ReceiveReadCache[0:n])

		// 转发消息到User routine
		// 通过比较n和 length循环获取批量消息包
		//`消息长度4+消息id4+序列号4+时间戳8+protobuf消息体`
		receiveBytes := user.ReceiveBuffer.Bytes()
		user.ReceiveBuffer.Reset()
		remainBytes := len(receiveBytes)
		index := 0
		//解析批量消息包
		for remainBytes > 0 {
			//小端
			length := int(uint32(receiveBytes[index]) | uint32(receiveBytes[index+1])<<8 | uint32(receiveBytes[index+2])<<16 | uint32(receiveBytes[index+3])<<24)
			length += 4 //客户端请求长度不包含自身
			if length > config2.MessageLimit {
				channelInactive(user, errors.New(fmt.Sprintf("消息太长")))
				return
			}
			//消息不够,缓存下次使用
			if length > remainBytes {
				user.ReceiveBuffer.Write(receiveBytes[index:])
				break
			}

			//packetData := make([]byte, length)
			packetData := mode.GetBytes()[:length] //用缓存池，减少垃圾回收，能有性能提升？
			copy(packetData, receiveBytes[index:index+length])
			user.ReceiveBytes <- packetData
			remainBytes = remainBytes - length
			index += length
			//log.Info("收到消息：读取长度=%v 消息长度=%v 剩余长度=%v", n, length, remainBytes)
		}

		//n, err = session.Write(user.ReceiveReadCache[:n])
		//if err != nil {
		//	log.Error("kcp启动失败：%v", err)
		//	channelInactive(session, err)
		//	return
		//}
	}
}

func (m *ClientManager) Run() {
}

func (m *ClientManager) Stop() {
}
