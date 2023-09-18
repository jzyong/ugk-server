package manager

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	"github.com/jzyong/ugk/common/config"
	"strings"
	"sync"
	"time"
)

// RedisManager redis
type RedisManager struct {
	util.DefaultModule
	redisClient  *redis.Client        //redis 单点 客户端
	redisCluster *redis.ClusterClient //redis cluster 客户端
	CmdAble      redis.Cmdable        //命令操作接口，部分特殊命令不支持 数据操作使用CmdAble，避免单点和集群的切换
}

var redisManager *RedisManager
var redisManagerOnce sync.Once

func GetRedisManager() *RedisManager {
	redisManagerOnce.Do(func() {
		redisManager = &RedisManager{}
	})
	return redisManager
}

// Init 开始启动
func (m *RedisManager) Init() error {
	log.Info("redis init started ......")
	return nil
}

// Start 启动
func (m *RedisManager) Start() {
	m.StartByZookeeper()
}

// StartByZookeeper 启动，通过 profile从zookeeper获取地址
func (m *RedisManager) StartByZookeeper() {

	redisConfigStr := util.ZKGet(GetZookeeperManager().GetConn(), fmt.Sprintf(config.ZKRedisPath, GetZookeeperManager().ServiceCfg.GetProfile()))

	var redisConfig = &config.RedisConfig{}
	json.Unmarshal([]byte(redisConfigStr), redisConfig)
	log.Info("redis database :%v", redisConfig)

	//集群连接
	if strings.Contains(redisConfig.Host, ":") {
		addrs := strings.Split(redisConfig.Host, ",")
		m.redisCluster = redis.NewClusterClient(&redis.ClusterOptions{Addrs: addrs})
		m.CmdAble = m.redisCluster
		//单独连接
	} else {
		m.redisClient = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%v:%v", redisConfig.Host, redisConfig.Port),
			Password: redisConfig.Password,
			DB:       0,
		})
		m.CmdAble = m.redisClient
	}

}

// HGet 获取map 对象
func (m *RedisManager) HGet(key, filed string) string {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	valueStr, err := m.CmdAble.HGet(ctx, key, filed).Result()
	if err == redis.Nil {
		valueStr = ""
		log.Info("key:%v filed:%v not exist", key, filed)
	} else if err != nil {
		log.Error("key:%v filed:%v error:%v", key, filed, err)
		valueStr = ""
	}
	return valueStr
}

// HGetAll 获取map 所有对象
func (m *RedisManager) HGetAll(key string) *redis.StringStringMapCmd {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	result := m.CmdAble.HGetAll(ctx, key)
	return result
}

// HGetObject 获取map对象 并且json反序列化为对象
func (m *RedisManager) HGetObject(key, field string, out interface{}) {
	result := m.HGet(key, field)
	if result == "" {
		return
	}
	json.Unmarshal([]byte(result), out)
}

// Stop 关闭连接
func (m *RedisManager) Stop() {
	if m.redisClient != nil {
		m.redisClient.Close()
	} else if m.redisCluster != nil {
		m.redisCluster.Close()
	}
}
