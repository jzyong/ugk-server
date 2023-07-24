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

// 连接关闭 TODO
func channelInactive(session *kcp.UDPSession) {

}

// handleEcho send back everything it received  TODO 编写自定义逻辑,连接建立和断开事件监测
func channelRead(conn *kcp.UDPSession) {
	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Error("kcp启动失败：%v", err)
			return
		}

		n, err = conn.Write(buf[:n])
		if err != nil {
			log.Error("kcp启动失败：%v", err)
			return
		}
	}
}

func (m *ClientManager) Run() {
}

func (m *ClientManager) Stop() {
}
