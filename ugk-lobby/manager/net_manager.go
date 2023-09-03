package manager

import (
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	"github.com/jzyong/ugk/lobby/mode"
	"sync"
)

// NetManager 网络，kcp连接网关 TODO
type NetManager struct {
	util.DefaultModule
}

var netManager *NetManager
var netSingletonOnce sync.Once

func GetNetManager() *NetManager {
	netSingletonOnce.Do(func() {
		netManager = &NetManager{}
	})
	return netManager
}

// 消息执行函数
type handFunc func(player *mode.Player, data []byte, seq uint32)

// NetHandlers 客户端消息处理器
var NetHandlers = make(map[uint32]handFunc)

func (m *NetManager) Init() error {
	log.Info("NetManager 初始化......")
	return nil
}

// TODO 连接网关，需要服务发现

func (m *NetManager) Run() {
}

func (m *NetManager) Stop() {
}
