package manager

import (
	"context"
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	"github.com/jzyong/ugk/common/manager"
	"github.com/jzyong/ugk/login/config"
	"github.com/jzyong/ugk/login/mode"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
)

// DataManager mongodb客户端  暂时没有配置读取
type DataManager struct {
	util.DefaultModule
	DataProcessChan chan func() //数据处理
	snowflake       *util.Snowflake
}

var dataManager = &DataManager{}
var dataSingletonOnce sync.Once

func GetDataManager() *DataManager {
	dataSingletonOnce.Do(func() {
		dataManager = &DataManager{
			DataProcessChan: make(chan func(), 1024),
			snowflake:       util.NewSnowflake(int16(config.BaseConfig.Id)),
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
	manager.GetMongoManager().Stop()
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

// 获取聊天室
func (dataManager *DataManager) FindAccount(id string) *mode.Account {

	collection := manager.GetMongoManager().GetProductionDB().Database(config.BaseConfig.MongoDbName).Collection("account")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result := collection.FindOne(ctx, bson.M{"_id": id})

	if result.Err() != nil {
		//插入新的实例对象
		if result.Err() == mongo.ErrNoDocuments {
			playerId, err := dataManager.GetNextSequence("playerId")
			if err != nil {
				log.Error("获取玩家ID失败：%v", err)
			}
			account := &mode.Account{
				Id:       id,
				PlayerId: playerId,
			}
			_, err = collection.InsertOne(ctx, account)
			if err != nil {
				log.Error("创建账号错误：%v", err)
			}
			return account
		}
		log.Error("查询账号错误：%v", result.Err())
		return nil
	}
	var account = &mode.Account{}
	result.Decode(&account)
	return account
}

// SaveAccount 保存账号
func (dataManager *DataManager) SaveAccount(id int64, update *bson.M) {
	collection := manager.GetMongoManager().GetProductionDB().Database(config.BaseConfig.MongoDbName).Collection("account")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := collection.UpdateByID(ctx, id, bson.D{{Key: "$set", Value: update}})
	//_, err := collection.UpdateOne(ctx, bson.M{"_id": receiver.serverInfo.Id}, bson.D{{Key: "$set", Value: update}})
	if err != nil {
		log.Error("更新账号错误：%v", err)
		return
	}
	log.Debug("更新账号%d --> %v ", id, res.ModifiedCount)
}

// DeleteAccount 删除聊天室
func (dataManager *DataManager) DeleteAccount(id int64) {
	collection := manager.GetMongoManager().GetProductionDB().Database(config.BaseConfig.MongoDbName).Collection("account")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	result, err := collection.DeleteOne(ctx, bson.D{{Key: "_id", Value: id}})
	defer cancel()

	if err != nil {
		log.Error("删除账号错误：%v", err)
		return
	}
	log.Debug("删除账号%v %v", id, result.DeletedCount)
}

func (dataManager *DataManager) GetNextSequence(sequenceName string) (int64, error) {
	collection := manager.GetMongoManager().GetProductionDB().Database(config.BaseConfig.MongoDbName).Collection("counter")

	filter := bson.M{"_id": sequenceName}
	update := bson.M{"$inc": bson.M{"value": 1}}

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	result := collection.FindOneAndUpdate(context.Background(), filter, update, opts)

	var counter mode.Counter
	if err := result.Decode(&counter); err != nil {
		return 0, err
	}

	return counter.Value, nil
}
