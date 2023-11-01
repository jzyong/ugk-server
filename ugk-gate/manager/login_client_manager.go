package manager

import (
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	config2 "github.com/jzyong/ugk/common/config"
	"github.com/jzyong/ugk/common/manager"
	"github.com/jzyong/ugk/gate/config"
	"sync"
)

// LoginClientManager 登录客户端

type LoginClientManager struct {
	util.DefaultModule
	Clients []*manager.ServiceClient //登录服列表
}

var loginClientManager *LoginClientManager
var loginClientSingletonOnce sync.Once

func GetLoginClientManager() *LoginClientManager {
	loginClientSingletonOnce.Do(func() {
		loginClientManager = &LoginClientManager{
			Clients: make([]*manager.ServiceClient, 0, 2),
		}
	})
	return loginClientManager
}

func (m *LoginClientManager) Init() error {
	log.Info("LoginClientManager 初始化......")

	return nil
}

func (m *LoginClientManager) Run() {

	path := config2.GetZKServicePath(config.BaseConfig.Profile, config2.LoginName, 0)
	manager.GetServiceClientManager().WatchGrpcService(path)
	manager.GetServiceClientManager().ClientAddAction = LoginAdd
	manager.GetServiceClientManager().ClientRemoveAction = LoginRemove
}

func (m *LoginClientManager) Stop() {
}

func (m *LoginClientManager) GetClient(id uint32) *manager.ServiceClient {
	for _, client := range m.Clients {
		if client.Id == id {
			return client
		}
	}
	return nil
}

// RandomClient 随机获取
func (m *LoginClientManager) RandomClient() *manager.ServiceClient {
	len := len(m.Clients)
	if len < 1 {
		log.Warn("无可用登录服")
		return nil
	}
	return m.Clients[int(util.RandomInt32(0, int32(len-1)))]
}

// HashModClient hash取余随机
func (m *LoginClientManager) HashModClient(token string) *manager.ServiceClient {
	len := len(m.Clients)
	if len < 1 {
		log.Warn("无可用登录服")
		return nil
	}
	hash := util.GetJavaIntHash(token)
	index := int(hash) % len
	return m.Clients[index]
}

func LoginAdd(client *manager.ServiceClient) {
	path := config2.GetZKServicePath(config.BaseConfig.Profile, config2.LoginName, 0)
	if client.Path == path {
		GetLoginClientManager().Clients = append(GetLoginClientManager().Clients, client)
		log.Info("登录客户端%v-%v 添加", client.Id, client.Url)
	}
}

func LoginRemove(client *manager.ServiceClient) {
	path := config2.GetZKServicePath(config.BaseConfig.Profile, config2.GateName, 0)
	if client.Path == path {
		for index, loginClient := range GetLoginClientManager().Clients {
			if loginClient.Id == client.Id {
				GetLoginClientManager().Clients = append(GetLoginClientManager().Clients[:index], GetLoginClientManager().Clients[index+1:]...)
				log.Info("登录客户端%v-%v 移除", client.Id, client.Url)
			}
		}
	}
}
