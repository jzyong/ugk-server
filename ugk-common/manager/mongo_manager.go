package manager

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	"github.com/jzyong/ugk/common/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"sync"
	"time"
)

// MongoManager Mongo
type MongoManager struct {
	util.DefaultModule
	configClient     *mongo.Client //配置数据 mongodb客户端
	ConfigDatabase   string        //配置表数据库名称
	productionClient *mongo.Client //生产数据库
}

var mongoManager *MongoManager
var mongoManagerOnce sync.Once

func GetMongoManager() *MongoManager {
	mongoManagerOnce.Do(func() {
		mongoManager = &MongoManager{}
	})
	return mongoManager
}

// Init 开始启动
func (m *MongoManager) Init() error {
	log.Info("mongo init started ......")
	return nil
}

// StartConfigDB 启动配置文件数据库
func (m *MongoManager) StartConfigDB() error {
	//1.连接数据库
	mongoConfigStr := util.ZKGet(GetZookeeperManager().GetConn(), fmt.Sprintf(config.ZKMongoExcelConfigPath, GetZookeeperManager().ServiceCfg.GetProfile()))
	var mongoConfig = &config.MongoConfig{}
	err := json.Unmarshal([]byte(mongoConfigStr), mongoConfig)
	if err != nil {
		log.Error("解析配置：%v", err)
		return err
	}
	log.Info("Excel mongo database :%v", mongoConfig)
	m.ConfigDatabase = mongoConfig.Database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoConfig.Url))
	if err != nil {
		log.Warn("mongodb客户端初始化失败")
		return err
	}
	m.configClient = client

	//err = m.configClient.Connect(ctx)
	//if err != nil {
	//	return err
	//}
	//检测连接状态
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Warn("mongodb客户端 断开连接")
		return err
	}
	return nil
}

// StartProductionDB 启动生产数据库
func (m *MongoManager) StartProductionDB(dbUrl string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbUrl))
	if err != nil {
		log.Error("启动mongodb：%v", err)
		return err
	}
	m.productionClient = client
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	//err = m.productionClient.Connect(ctx)
	//if err != nil {
	//	GetSentryManager().Exception(err)
	//	return err
	//}
	//检测连接状态
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	}
	return nil
}

func (m *MongoManager) GetConfigDB() *mongo.Client {
	if m.configClient == nil {
		log.Warn("配置数据库 mongoDB未创建")
		return nil
	}
	return m.configClient
}

func (m *MongoManager) GetProductionDB() *mongo.Client {
	if m.productionClient == nil {
		log.Warn("生产数据库 mongoDB未创建")
		return nil
	}
	return m.productionClient
}

// StructToM 将结构体转换为更新的M
func (m *MongoManager) StructToM(o interface{}) *bson.M {
	bytes, err := bson.Marshal(o)
	if err != nil {
		log.Error("序列化存储对象：%v", o)
		return nil
	}
	var update bson.M
	err = bson.Unmarshal(bytes, &update)
	if err != nil {
		return nil
	}
	delete(update, "_id") //删除主键
	return &update
}

// Stop 关闭连接
func (m *MongoManager) Stop() {
	if m.configClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if err := m.configClient.Disconnect(ctx); err != nil {
			log.Error("关闭配置表mongodb:%v", err)
		}
	}
	if m.productionClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if err := m.productionClient.Disconnect(ctx); err != nil {
			log.Error("关闭mongodb:%v", err)
		}
	}
}
