package manager

import (
	"context"
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	"github.com/jzyong/ugk/api/config"
	"github.com/jzyong/ugk/common/manager"
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
