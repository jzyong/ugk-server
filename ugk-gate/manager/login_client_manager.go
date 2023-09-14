package manager

import (
	"encoding/json"
	"fmt"
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	config2 "github.com/jzyong/ugk/common/config"
	"github.com/jzyong/ugk/common/manager"
	"github.com/jzyong/ugk/gate/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"strconv"
	"sync"
	"time"
)

// LoginClientManager 登录客户端

type LoginClientManager struct {
	util.DefaultModule
	Clients           []*LoginClient //登录服列表
	watchLoginService bool           //是否在监听
}

var loginClientManager *LoginClientManager
var loginClientSingletonOnce sync.Once

func GetLoginClientManager() *LoginClientManager {
	loginClientSingletonOnce.Do(func() {
		loginClientManager = &LoginClientManager{
			Clients: make([]*LoginClient, 0, 2),
		}
	})
	return loginClientManager
}

func (m *LoginClientManager) Init() error {
	log.Info("LoginClientManager 初始化......")
	m.watchService()
	go m.run()

	return nil
}

func (m *LoginClientManager) Run() {
}

func (m *LoginClientManager) Stop() {
}

func (m *LoginClientManager) run() {
	ticker := time.Tick(time.Second * 5)
	for {
		select {
		case <-ticker:
			//定时检测是否在监听,如果监听
			if m.watchLoginService == false {
				m.watchService()
				log.Warn("无可用登录服")
			}
		}
	}
}

// watchService 监听大厅grpc服务，报异常会终止监听
func (m *LoginClientManager) watchService() {
	baseConfig := config.BaseConfig
	path := config2.GetZKServicePath(baseConfig.Profile, "login", 0)
	conn := manager.GetZookeeperManager().GetConn()
	children, errors := util.ZKWatchChildrenW(conn, path, true)
	m.watchLoginService = true
	go func() {
		for m.watchLoginService {
			select {
			case serverIds := <-children:
				log.Info("登录服变更为：%v", serverIds)
				m.updateLoginClient(serverIds)
			case err := <-errors:
				//如果启动服务器监听节点或大厅全部关闭出现：zk: node does not exist，则后面无法再进行监听？多层父级监听？
				log.Warn("登录服监听报错：%v", err)
				m.watchLoginService = false
				return
			}
		}
	}()
}

func (m *LoginClientManager) updateLoginClient(serverIds []string) {
	//遍历添加新连接
	for _, serverIdStr := range serverIds {
		m.addLoginClient(serverIdStr)
	}
	//删除已关闭的网关  note遍历中删除了对象
	if len(serverIds) < len(m.Clients) {
		var deleteClient *LoginClient = nil
		for _, c := range m.Clients {
			if util.SliceContains(serverIds, c.ServiceConfig.Id) < 0 {
				deleteClient = m.GetClient(c.Id)
			}
		}
		if deleteClient != nil {
			m.removeLoginClient(deleteClient)
		}
	}
}

func (m *LoginClientManager) GetClient(id uint32) *LoginClient {
	for _, client := range m.Clients {
		if client.Id == id {
			return client
		}
	}
	return nil
}

func (m *LoginClientManager) RandomClient() *LoginClient {
	len := len(m.Clients)
	if len < 1 {
		log.Warn("无可用登录服")
		return nil
	}
	return m.Clients[int(util.RandomInt32(0, int32(len-1)))]
}

// 添加大厅 客户端
func (m *LoginClientManager) addLoginClient(serverIdStr string) {
	serverId, err := strconv.ParseInt(serverIdStr, 10, 32)
	if err != nil {
		log.Warn("登录服ID错误： %v =>%v", serverIdStr, err)
		return
	}
	client := m.GetClient(uint32(serverId))
	if client != nil {
		if client.ClientConn.GetState() == connectivity.Shutdown { //移除老连接
			log.Info("连接已关闭")
			m.removeLoginClient(client)
		} else {
			return
		}
	}
	//连接服务器
	conn := manager.GetZookeeperManager().GetConn()
	path := config2.GetZKServicePath(config.BaseConfig.Profile, "login", uint32(serverId))
	serviceConfigStr := util.ZKGet(conn, path)
	if serviceConfigStr == "" {
		log.Warn("%v 登录服配置未找到", serverIdStr)
		return
	}

	var serviceConfig = &util.ServiceConfig{}
	json.Unmarshal([]byte(serviceConfigStr), serviceConfig)
	var serverUrl = fmt.Sprintf("%s:%d", serviceConfig.Address, serviceConfig.Port)
	dialOption := grpc.WithTransportCredentials(insecure.NewCredentials())
	clientConn, err := grpc.Dial(serverUrl, dialOption)
	if err != nil {
		log.Error("连接登录服%v错误:%v", serverIdStr, err)
		return
	}
	client = &LoginClient{
		Id:            uint32(serverId),
		ServiceConfig: serviceConfig,
		ClientConn:    clientConn,
	}
	m.Clients = append(m.Clients, client)
	log.Info("添加登录服：%v url：%v ", serverIdStr, serverUrl)
}

// 移除客户端
func (m *LoginClientManager) removeLoginClient(client *LoginClient) {
	for i, c := range m.Clients {
		if c.Id == client.Id {
			m.Clients = append(m.Clients[:i], m.Clients[i+1:]...)
			log.Info("登录服 %d-%s 连接移除", client.Id, client.ServiceConfig.Address)
			break
		}
	}
}

// LoginClient 登录客户端
type LoginClient struct {
	Id            uint32              //服务id
	ServiceConfig *util.ServiceConfig //登录服配置
	ClientConn    *grpc.ClientConn    //连接
}
