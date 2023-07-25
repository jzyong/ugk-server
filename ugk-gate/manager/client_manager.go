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
			channelActive(s)
			go channelRead(s)

		}
	} else {
		log.Error("kcp启动失败：%v", err)
	}

}

// 连接激活
func channelActive(session *kcp.UDPSession) {
	//TODO
	log.Info("%s 连接创建", session.RemoteAddr().String())
}

// 连接关闭
// 客户端强制杀进程，服务器不知道连接断开。kcp-go源码没有示例,因此使用自定义心跳（每2s请求一次心跳，超过10s断开连接）
func channelInactive(session *kcp.UDPSession, err error) {
	log.Info("%s 连接关闭:%s", session.RemoteAddr(), err)
}

//	处理接收消息
//	UDPSession Read和Write都可能阻塞，公用routine是否会被阻塞自定义逻辑？
//
// TODO 编写自定义逻辑,需要关闭此routine
func channelRead(conn *kcp.UDPSession) {
	buf := make([]byte, 4096)
	for {
		//会阻塞
		n, err := conn.Read(buf)
		if err != nil {
			log.Error("kcp启动失败：%v", err)
			channelInactive(conn, err)
			return
		}

		n, err = conn.Write(buf[:n])
		if err != nil {
			log.Error("kcp启动失败：%v", err)
			channelInactive(conn, err)
			return
		}
	}
}

func (m *ClientManager) Run() {
}

func (m *ClientManager) Stop() {
}
