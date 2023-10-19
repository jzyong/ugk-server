package manager

import (
	"context"
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	"github.com/jzyong/ugk/common/manager"
	"github.com/jzyong/ugk/lobby/config"
	"github.com/jzyong/ugk/lobby/mode"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
	"time"
)

// DataManager mongodb客户端  暂时没有配置读取
type DataManager struct {
	util.DefaultModule
	DataProcessChan chan func() //数据处理
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

func (dataManager *DataManager) FindPlayer(id int64) *mode.Player {

	collection := manager.GetMongoManager().GetProductionDB().Database(config.BaseConfig.MongoDbName).Collection("player")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result := collection.FindOne(ctx, bson.M{"_id": id})
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return nil
		}
		log.Error("查询玩家错误：%v", result.Err())
		return nil
	}
	var player = &mode.Player{}
	result.Decode(&player)
	return player
}

func (dataManager *DataManager) InsertPlayer(player *mode.Player) {
	collection := manager.GetMongoManager().GetProductionDB().Database(config.BaseConfig.MongoDbName).Collection("player")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := collection.InsertOne(ctx, player)
	if err != nil {
		log.Error("创建玩家错误：%v", err)
	}
	log.Debug("创建玩家 --> %v ", res.InsertedID)
}

func (dataManager *DataManager) SavePlayer(id int64, update *bson.M) {
	collection := manager.GetMongoManager().GetProductionDB().Database(config.BaseConfig.MongoDbName).Collection("player")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := collection.UpdateByID(ctx, id, bson.D{{Key: "$set", Value: update}})
	//_, err := collection.UpdateOne(ctx, bson.M{"_id": receiver.serverInfo.Id}, bson.D{{Key: "$set", Value: update}})
	if err != nil {
		log.Error("更新玩家错误：%v", err)
		return
	}
	log.Debug("更新玩家%d --> %v ", id, res.ModifiedCount)
}
