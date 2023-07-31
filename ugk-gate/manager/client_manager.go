package manager

import (
	"fmt"
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
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

func (m *ClientManager) Init() error {
	log.Info("ClientManager 初始化......")
	go m.runKcpServer()
	return nil
}

// 启动kcp服务器
func (m *ClientManager) runKcpServer() {
	url := fmt.Sprintf("%v:%v", config.AppConfigManager.PublicIp, config.AppConfigManager.ClientPort)
	log.Info("玩家udp监听地址：%s", url)
	if listener, err := kcp.ListenWithOptions(url, nil, 0, 0); err == nil {
		for {
			s, err := listener.AcceptKCP()
			if err != nil {
				log.Error("kcp启动失败：%v", err)
			}

			//设置参数 https://github.com/skywind3000/kcp/blob/master/README.en.md#protocol-configuration
			//s.SetACKNoDelay(true)
			s.SetMtu(4096)
			s.SetWindowSize(4096, 4096)
			s.SetReadBuffer(8 * 1024 * 1024)
			s.SetWriteBuffer(16 * 1024 * 1024)
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
func channelInactive(session *kcp.UDPSession, err error) {
	log.Info("%s 连接关闭:%s", session.RemoteAddr(), err)
}

// 处理接收消息
// UDPSession Read和Write都可能阻塞，共用routine是否会被阻塞自定义逻辑？
func channelRead(user *User) {
	session := user.ClientSession
	for {
		//会阻塞
		//TODO 最多读取1352 字节，超过的消息包怎么读取？
		n, err := session.Read(user.ReceiveReadCache)
		if err != nil {
			log.Error("kcp启动失败：%v", err)
			channelInactive(session, err)
			return
		}

		// 转发消息到User routine
		// 通过比较n和 length循环获取批量消息包
		//`消息长度4+消息id4+序列号4+时间戳8+protobuf消息体`
		remainBytes := n
		index := 0
		//解析批量消息包
		for remainBytes > 0 {
			//小端
			length := int(uint32(user.ReceiveReadCache[index]) | uint32(user.ReceiveReadCache[index+1])<<8 | uint32(user.ReceiveReadCache[index+2])<<16 | uint32(user.ReceiveReadCache[index+3])<<24)
			length += 4 //客户端请求长度不包含自身
			packetData := make([]byte, length)
			copy(packetData, user.ReceiveReadCache[index:index+length])
			user.ReceiveBytes <- packetData
			remainBytes = remainBytes - length
			index += length
			log.Info("收到消息：读取长度=%v 消息长度=%v 剩余长度=%v", n, length, remainBytes)
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
