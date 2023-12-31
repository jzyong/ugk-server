package manager

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-zookeeper/zk"
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	config2 "github.com/jzyong/ugk/common/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"strconv"
	"sync"
	"time"
)

// ServiceClient 微服务客户端
type ServiceClient struct {
	Id         uint32           //服务器id
	Url        string           //连接地址
	Path       string           //路径或类型
	ClientConn *grpc.ClientConn //客户端连接
}

// ServiceClientManager 处理微服务客户端逻辑，管理连接
type ServiceClientManager struct {
	util.DefaultModule
	clients            map[string]map[uint32]*ServiceClient //grpc客户端  类型（路径）：id：客户端
	watchService       map[string]bool                      //服务是否在监听 key：路径
	clientsLock        sync.RWMutex                         //读写锁
	ClientAddAction    func(client *ServiceClient)          //客户端增加回调
	ClientRemoveAction func(client *ServiceClient)          //客户端移除回调
}

var serviceClientManager *ServiceClientManager
var serviceClientManagerOnce sync.Once

func GetServiceClientManager() *ServiceClientManager {
	serviceClientManagerOnce.Do(func() {
		serviceClientManager = &ServiceClientManager{
			clients:      make(map[string]map[uint32]*ServiceClient),
			watchService: make(map[string]bool),
		}
	})

	return serviceClientManager
}

func (m *ServiceClientManager) Init() error {
	log.Info("ServiceClientManager init start......")

	log.Info("ServiceClientManager init end......")
	return nil
}

func (m *ServiceClientManager) Run() {
	go m.run()
}

func (m *ServiceClientManager) run() {
	ticker := time.Tick(time.Second * 5)
	for {
		select {
		case <-ticker:
			//定时检测是否在监听,如果监听
			for path, watch := range m.watchService {
				if watch == false {
					go m.WatchGrpcService(path)
					log.Warn("%v 无可用服务", path)
				}
			}
		}
	}
}

// UpdateClient 更新微服务客户端
func (m *ServiceClientManager) UpdateClient(serverIds []string, zkConnect *zk.Conn, path string) {
	if m.clients == nil {
		log.Warn("clients map 还未初始化就收到更新大厅客户端，注意初始化顺序")
		return
	}

	//遍历添加新连接
	for _, serverIdStr := range serverIds {
		m.addClient(serverIdStr, zkConnect, path)
	}
	clients := m.clients[path]
	if clients == nil {
		return
	}
	//删除已关闭的服务  note遍历中删除了对象
	if len(serverIds) < len(m.clients) {
		for serverId, _ := range clients {
			idStr := strconv.Itoa(int(serverId))
			if util.SliceContains(serverIds, idStr) < 0 {
				client := clients[serverId]
				log.Info("server Id:%d already close", serverId)
				m.removeClient(client)
			}
		}
	}
}

// 添加 客户端
func (m *ServiceClientManager) addClient(serverIdStr string, zkConnect *zk.Conn, path string) {
	serverId, err := strconv.ParseInt(serverIdStr, 10, 32)
	if err != nil {
		log.Warn("hall error Id： %v =>%v", serverIdStr, err)
		return
	}
	defer m.clientsLock.Unlock()
	m.clientsLock.Lock()
	clients := m.clients[path]
	if clients == nil {
		clients = make(map[uint32]*ServiceClient)
		m.clients[path] = clients
	}

	if client, ok := clients[uint32(serverId)]; ok {
		if client.ClientConn.GetState() == connectivity.Shutdown { //移除老连接
			log.Info("server MessageId:%d already close,open new connection")
			m.removeClient(client)
		} else {
			log.Info("%v service %v already connected", client.Path, serverIdStr)
			return
		}
	}
	//连接服务器
	serverRpcConfigStr := util.ZKGet(zkConnect, fmt.Sprintf("%v/%v", path, serverIdStr))
	type RpcConfig struct {
		Address string `json:"address"`
		Port    int    `json:"port"`
	}
	//log.Info("%v 服务器地址：%v", fmt.Sprintf("%v/%v", Path, serverIdStr), serverRpcConfigStr)
	if serverRpcConfigStr == "" {
		log.Warn("%v service config is nil", serverIdStr)
		return
	}

	var serverRpcConfig = &RpcConfig{}
	json.Unmarshal([]byte(serverRpcConfigStr), serverRpcConfig)
	var serverUrl = fmt.Sprintf("%s:%d", serverRpcConfig.Address, serverRpcConfig.Port)
	dialOption := grpc.WithTransportCredentials(insecure.NewCredentials())
	clientConn, err := grpc.Dial(serverUrl, dialOption)
	if err != nil {
		log.Error("add  service %v-%v fail:%v", path, serverIdStr, err)
		return
	}
	client := &ServiceClient{
		Id:         uint32(serverId),
		Url:        serverUrl,
		Path:       path,
		ClientConn: clientConn,
	}
	clients[client.Id] = client
	if m.ClientAddAction != nil {
		m.ClientAddAction(client)
	}
	log.Info("add Client：%v-%v Url=%v ", path, serverIdStr, serverUrl)
}

// GetClients 获取客户端，参数为路径
func (m *ServiceClientManager) GetClients(path string) map[uint32]*ServiceClient {
	defer m.clientsLock.RUnlock()
	m.clientsLock.RLock()
	return m.clients[path]
}

// GetClient 获取客户端，参数为路径
func (m *ServiceClientManager) GetClient(path string, severId uint32) *ServiceClient {
	defer m.clientsLock.RUnlock()
	m.clientsLock.RLock()
	if clients, ok := m.clients[path]; ok {
		return clients[severId]
	}
	return nil
}

// 移除客户端
func (m *ServiceClientManager) removeClient(client *ServiceClient) {
	m.clientsLock.Lock()
	defer m.clientsLock.Unlock()
	if clients, ok := m.clients[client.Path]; ok {
		delete(clients, client.Id)
		if m.ClientRemoveAction != nil {
			m.ClientRemoveAction(client)
		}
		log.Info("服务 %d-%s 连接移除", client.Id, client.Url)
	}
}

// ReloadDefaultService 使用默认的zookeeper加载一下服务
func (m *ServiceClientManager) ReloadDefaultService(path string) {
	m.ReloadService(GetZookeeperManager().GetConn(), GetZookeeperManager().ServiceCfg.GetProfile(), path)
}

// ReloadService 加载一下服务
func (m *ServiceClientManager) ReloadService(conn *zk.Conn, profile, path string) {
	path = fmt.Sprintf(path, profile)
	c, _, err := conn.Children(path)
	if err != nil {
		log.Warn("加载服务异常：%v", err)
		return
	}
	//现在拥有连接大于存在连接
	if len(m.clients) >= len(c) {
		return
	}

	for _, serverIdStr := range c {
		m.addClient(serverIdStr, conn, path)
		log.Info(" service %v", serverIdStr)
	}
}

// GetGrpc  id 小于0，直接返回第一个
func (m *ServiceClientManager) GetGrpc(path string, id uint32) (*grpc.ClientConn, error) {
	m.clientsLock.RLock()
	defer m.clientsLock.RUnlock()
	if clients, ok := m.clients[path]; ok {
		if len(clients) < 1 {
			return nil, errors.New(fmt.Sprintf("路径：%s 无可用服务", path))
		}
		//id 小于1，直接返回第一个
		if id < 1 {
			for _, client := range clients {
				return client.ClientConn, nil
			}
		} else {
			if client, ok := clients[id]; ok {
				return client.ClientConn, nil
			} else {
				return nil, errors.New(fmt.Sprintf("路径：%s Id:%d 服务不可用", path, id))
			}
		}
	}
	return nil, errors.New(fmt.Sprintf("路径：%s 服务不存在", path))
}

// GetLobbyGrpcByServerId 通过服务器Id获取大厅
func (m *ServiceClientManager) GetLobbyGrpcByServerId(serverId uint32) (*grpc.ClientConn, error) {
	path := config2.GetZKServicePath(GetZookeeperManager().ServiceCfg.GetProfile(), config2.LobbyName, 0)
	return m.GetGrpc(path, serverId)
}

// GetLobbyGrpcByPlayerId 玩家获取大厅 grpc客户端 大厅Id
func (m *ServiceClientManager) GetLobbyGrpcByPlayerId(playerId int64) (*grpc.ClientConn, error, uint32) {
	path := config2.GetZKServicePath(GetZookeeperManager().ServiceCfg.GetProfile(), config2.LobbyName, 0)
	serverIdStr := GetRedisManager().HGet(config2.RedisPlayerLocation, fmt.Sprintf("%v", playerId))
	m.clientsLock.RLock()
	defer m.clientsLock.RUnlock()

	if clients, ok := m.clients[path]; ok {
		if len(clients) < 1 {
			return nil, errors.New(fmt.Sprintf("路径：%s 无可用服务", path)), 0
		}
		//根据玩家上次id获取
		if serverIdStr != "" {
			serverId := uint32(util.ParseInt32(serverIdStr))
			for _, client := range clients {
				if client.Id == serverId {
					return client.ClientConn, nil, client.Id
				}
			}
		} else {
			log.Warn("玩家%d 未能获得位置，玩家ID是否非法", playerId)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()
			for _, client := range clients {
				GetRedisManager().CmdAble.HSet(ctx, config2.RedisPlayerLocation, fmt.Sprintf("%v", playerId))
				return client.ClientConn, nil, client.Id
			}
		}
	}
	return nil, errors.New(fmt.Sprintf("路径：%s 服务不存在", path)), 0
}

func (m *ServiceClientManager) Stop() {
	// 关闭grpc服务器
	for _, clients := range m.clients {
		for _, client := range clients {
			client.ClientConn.Close()
		}
	}
}

// WatchGrpcService  监听grpc服务
func (m *ServiceClientManager) WatchGrpcService(path string) {
	children, errors := util.ZKWatchChildrenW(GetZookeeperManager().GetConn(), path, false)
	m.clientsLock.Lock()
	m.watchService[path] = true
	m.clientsLock.Unlock()
	go func(p string) {
		for {
			select {
			case serverIds := <-children:
				log.Info("%v service change to：%v", p, serverIds)
				m.UpdateClient(serverIds, GetZookeeperManager().GetConn(), p)
			case err := <-errors:
				log.Warn("%v service listen error：%v", p, err)
				m.clientsLock.Lock()
				m.watchService[p] = false
				m.clientsLock.Unlock()
				return
			}
		}
	}(path)
}
