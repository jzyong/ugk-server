package constant

import (
	"fmt"
	"time"
)

//网络相关常量

const (
	MTU                 = 1200             //mtu 长度 ，(reduced to 1200 to fit all cases: https://en.wikipedia.org/wiki/Maximum_transmission_unit ; steam uses 1200 too!)
	MessageLimit        = 4000             // 消息长度限制
	WindowSize          = 4096             //窗口大小
	ClientHeartInterval = time.Second * 10 //客户端心跳时间
	ServerHeartInterval = time.Minute * 10 //客户端心跳时间
)

// zookeeper服务注册地址

const (
	ZKServicePath          = "/ugk/%s/service/%s"  //服务地址
	ZKMongoExcelConfigPath = "/ugk/%s/mongo/excel" //MongoDB Excel 配置数据地址
	ZKInfluxPath           = "/ugk/%s/influx"      //InfluxDB地址
	ZKRedisPath            = "/ugk/%s/redis"       //redis地址
)

// GetZKServicePath 获取服务地址
func GetZKServicePath(profile string, name ServiceName, serverId uint32) string {
	path := fmt.Sprintf(ZKServicePath, profile, name)
	if serverId > 0 {
		path = fmt.Sprintf("%s/%d", path, serverId)
	}
	return path
}
