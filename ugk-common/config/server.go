package config

import "fmt"

// zookeeper服务注册地址

const (
	ZKServicePath          = "/ugk/%s/service/%s"  //服务地址
	ZKServiceConfigPath    = "/ugk/%s/config/%s"   //服务配置地址
	ZKMongoExcelConfigPath = "/ugk/%s/mongo/excel" //MongoDB Excel 配置数据地址
	ZKInfluxPath           = "/ugk/%s/influx"      //InfluxDB地址
	ZKRedisPath            = "/ugk/%s/redis"       //redis地址
)

// GetZKServicePath 获取服务地址 serverId >0拼接服务器id路径
func GetZKServicePath(profile string, name string, serverId uint32) string {
	path := fmt.Sprintf(ZKServicePath, profile, name)
	if serverId > 0 {
		path = fmt.Sprintf("%s/%d", path, serverId)
	}
	return path
}

// GetZKServiceConfigPath 获取服务地址
func GetZKServiceConfigPath(profile string, name string, serverId uint32) string {
	path := fmt.Sprintf(ZKServiceConfigPath, profile, name)
	if serverId > 0 {
		path = fmt.Sprintf("%s/%d", path, serverId)
	}
	return path
}

// ServiceConfig 服务配置，所有服务实现次接口
type ServiceConfig interface {
	GetId() uint32              //服务器id
	GetName() string            //服务名称
	GetProfile() string         //配置文件选项|平台
	GetRpcUrl() string          //rpc 地址
	GetLogLevel() string        //日志级别
	GetZookeeperUrls() []string //zookeeper地址
}

// ServiceConfigImpl 服务配置默认实现
type ServiceConfigImpl struct {
	Id            uint32   `json:"id"`            //服务器ID
	Name          string   `json:"name"`          //服务名字
	Profile       string   `json:"profile"`       //个性化配置
	RpcUrl        string   `json:"rpcUrl"`        //rpc 地址
	LogLevel      string   `json:"logLevel"`      //日志级别
	ZookeeperUrls []string `json:"zookeeperUrls"` //zookeeper 地址

}

func (s *ServiceConfigImpl) GetId() uint32 {
	return s.Id
}

func (s *ServiceConfigImpl) GetName() string {
	return s.Name
}

func (s *ServiceConfigImpl) GetProfile() string {
	return s.Profile
}

func (s *ServiceConfigImpl) GetRpcUrl() string {
	return s.RpcUrl
}

func (s *ServiceConfigImpl) GetZookeeperUrls() []string {
	return s.ZookeeperUrls
}

func (s *ServiceConfigImpl) GetLogLevel() string {
	return s.LogLevel
}
