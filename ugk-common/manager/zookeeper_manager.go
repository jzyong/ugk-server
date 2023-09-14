package manager

import (
	"encoding/json"
	"github.com/go-zookeeper/zk"
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	"github.com/jzyong/ugk/common/config"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ZookeeperConfig zookeeper 配置
type ZookeeperConfig struct {
}

// ZookeeperManager  zookeeper服务发现
type ZookeeperManager struct {
	util.DefaultModule
	conn       *zk.Conn             //zookeeper连接
	ServiceCfg config.ServiceConfig //服务配置
}

var zookeeperManager *ZookeeperManager
var zookeeperManagerOnce sync.Once

func GetZookeeperManager() *ZookeeperManager {
	zookeeperManagerOnce.Do(func() {
		zookeeperManager = &ZookeeperManager{}
	})
	return zookeeperManager
}

func (m *ZookeeperManager) Init() error {
	log.Info("ZookeeperManager init start......")

	log.Info("ZookeeperManager init end......")
	return nil
}

// GetConn 获取zookeeper连接
func (m *ZookeeperManager) GetConn() *zk.Conn {
	if m.conn == nil {
		log.Warn("Zookeeper %v 连接未创建", m.ServiceCfg.GetZookeeperUrls())
		return nil
	}
	return m.conn
}

// Start 启动Zookeeper
func (m *ZookeeperManager) Start(serviceCfg config.ServiceConfig) {
	//1.设置服务器配置信息
	zkConnect := util.ZKCreateConnect(serviceCfg.GetZookeeperUrls())
	m.conn = zkConnect
	m.ServiceCfg = serviceCfg
	m.writeServerConfigAndRpcService(serviceCfg)
}

// WriteServerConfigAndRpcService 写服务器配置和服务grpc监听端口服务
func (m *ZookeeperManager) writeServerConfigAndRpcService(serviceCfg config.ServiceConfig) {
	configBytes, _ := json.Marshal(serviceCfg)
	serviceConfigPath := config.GetZKServiceConfigPath(serviceCfg.GetProfile(), serviceCfg.GetName(), serviceCfg.GetId())
	util.ZKUpdate(m.GetConn(), serviceConfigPath, string(configBytes))
	//2.设置grpc服务
	addressPort := strings.Split(serviceCfg.GetRpcUrl(), ":")
	port, _ := strconv.ParseInt(addressPort[1], 10, 32)
	serviceConfig := &util.ServiceConfig{
		Name:                serviceCfg.GetName(),
		Id:                  strconv.Itoa(int(serviceCfg.GetId())),
		Address:             addressPort[0],
		Port:                int32(port),
		ServiceType:         "DYNAMIC",
		RegistrationTimeUTC: time.Now().Unix() * 1000,
	}
	serviceBytes, _ := json.Marshal(serviceConfig)
	servicePath := config.GetZKServicePath(serviceCfg.GetProfile(), serviceCfg.GetName(), serviceCfg.GetId())
	util.ZKDelete(m.GetConn(), servicePath) //先删除一下，重启启动太快，还没删除，添加失败
	util.ZKAdd(m.GetConn(), servicePath, string(serviceBytes), zk.FlagEphemeral)
}

func (m *ZookeeperManager) Stop() {
	if m.conn != nil {
		m.conn.Close()
	}
}
