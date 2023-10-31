package manager

import (
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	"github.com/jzyong/ugk/api/config"
	config2 "github.com/jzyong/ugk/common/config"
	"github.com/jzyong/ugk/common/manager"
	"sync"
)

// GateClientManager gateClient
type GateClientManager struct {
	util.DefaultModule
	gateClients []*manager.ServiceClient
}

var gateClientManager *GateClientManager
var gateClientSingletonOnce sync.Once

func GetGateClientManager() *GateClientManager {
	gateClientSingletonOnce.Do(func() {
		gateClientManager = &GateClientManager{
			gateClients: make([]*manager.ServiceClient, 0, 2),
		}
	})
	return gateClientManager
}

func (m *GateClientManager) Init() error {
	log.Info("GateClientManager 初始化......")
	return nil
}

func (m *GateClientManager) Run() {
	path := config2.GetZKServicePath(config.BaseConfig.Profile, config2.GateName, 0)
	manager.GetServiceClientManager().WatchGrpcService(path)
	manager.GetServiceClientManager().ClientAddAction = GateAdd
	manager.GetServiceClientManager().ClientRemoveAction = GateRemove
}

func (m *GateClientManager) Stop() {
}

// HashModGate 通过hash值取余获得网关
func (m *GateClientManager) HashModGate(ip string) string {
	if len(m.gateClients) < 1 {
		log.Info("无可用网关服务器")
		return "192.168.110.2:5000"
	}
	hash := util.GetJavaIntHash(ip)
	index := int(hash) % len(m.gateClients)
	return m.gateClients[index].Url
}

//// 初始化
//func (m *GateClientManager) initGate() {
//	path := config2.GetZKServicePath(config.BaseConfig.Profile, config2.GateName, 0)
//	manager.GetServiceClientManager().ReloadDefaultService(path)
//	manager.GetServiceClientManager().g
//}

// GateAdd 网关添加
func GateAdd(client *manager.ServiceClient) {
	path := config2.GetZKServicePath(config.BaseConfig.Profile, config2.GateName, 0)
	if client.Path == path {
		GetGateClientManager().gateClients = append(GetGateClientManager().gateClients, client)
		log.Info("网关%v-%v 添加", client.Id, client.Url)
	}
}

// GateRemove 网关移除
func GateRemove(client *manager.ServiceClient) {
	path := config2.GetZKServicePath(config.BaseConfig.Profile, config2.GateName, 0)
	if client.Path == path {
		for index, gateClient := range GetGateClientManager().gateClients {
			if gateClient.Id == client.Id {
				GetGateClientManager().gateClients = append(GetGateClientManager().gateClients[:index], GetGateClientManager().gateClients[index+1:]...)
				log.Info("网关%v-%v 移除", client.Id, client.Url)
			}
		}
	}
}
