package manager

import (
	"context"
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	"github.com/jzyong/ugk/client-drive-match/config"
	"github.com/jzyong/ugk/client-drive-match/mode"
	"github.com/jzyong/ugk/common/manager"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
	"time"
)

// DataManager mongodb客户端  暂时没有配置读取
type DataManager struct {
	util.DefaultModule
	DataProcessChan chan func()  //数据处理
	server          *mode.Server //服务器全局数据
}

var dataManager = &DataManager{}
var dataSingletonOnce sync.Once

func GetDataManager() *DataManager {
	dataSingletonOnce.Do(func() {
		dataManager = &DataManager{
			DataProcessChan: make(chan func(), 1024),
		}
	})
	return dataManager
}

// Init 开始启动
func (dataManager *DataManager) Init() error {
	//1.连接数据库
	manager.GetMongoManager().StartConfigDB()
	manager.GetMongoManager().StartProductionDB(config.BaseConfig.MongoUrl)

	//2. 加载配置
	dataManager.ReloadConfig(context.Background())

	//3. 启动redis
	manager.GetRedisManager().Start()

	log.Info("data init started ......")
	return nil
}

func (dataManager *DataManager) Run() {

	//单独routine处理数据
	go dataManager.run()

}

// 单独routine处理数据，避免并发问题
func (dataManager *DataManager) run() {
	saveTicker := time.Tick(time.Minute * 10)
	for {
		select {
		case saveFun := <-dataManager.DataProcessChan:
			saveFun()

		case <-saveTicker:
			dataManager.batchSave()
		}
	}
}

// 批量存储数据
func (dataManager *DataManager) batchSave() {
	if dataManager.GetServer().GetDirty() {
		dataManager.SaveServer(manager.GetMongoManager().StructToM(dataManager.server))
	}
}

// Stop 关闭连接
func (dataManager *DataManager) Stop() {
}

// ReloadConfig 加载配置 暂时不用，没有配置
func (dataManager *DataManager) ReloadConfig(ctx context.Context) string {

	configDB := manager.GetMongoManager().GetConfigDB()
	if configDB == nil {
		return "加载配置失败，配置数据库连接未创建"
	}

	// 暂时没有配置加载
	return "加载配置成功"
}

// GetServer 获取服务器信息
func (dataManager *DataManager) GetServer() *mode.Server {
	if dataManager.server != nil {
		return dataManager.server
	}

	collection := manager.GetMongoManager().GetProductionDB().Database(config.BaseConfig.MongoDbName).Collection("server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result := collection.FindOne(ctx, bson.M{"_id": 1})

	if result.Err() != nil {
		//插入新的实例对象
		if result.Err() == mongo.ErrNoDocuments {
			server := &mode.Server{
				Id: 1,
			}
			_, err := collection.InsertOne(ctx, server)
			if err != nil {
				log.Error("创建账号错误：%v", err)
			}
			return server
		}
		log.Error("查询错误：%v", result.Err())
		return nil
	}
	var server = &mode.Server{}
	result.Decode(&server)
	return server
}

func (dataManager *DataManager) SaveServer(update *bson.M) {
	collection := manager.GetMongoManager().GetProductionDB().Database(config.BaseConfig.MongoDbName).Collection("server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := collection.UpdateByID(ctx, 1, bson.D{{Key: "$set", Value: update}})
	if err != nil {
		log.Error("更新错误：%v", err)
		return
	}
	log.Debug("更新server%d --> %v ", 1, res.ModifiedCount)
}
