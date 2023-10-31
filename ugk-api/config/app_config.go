package config

import (
	"encoding/json"
	"flag"
	"github.com/jzyong/golib/log"
	config2 "github.com/jzyong/ugk/common/config"
	"os"
)

// BaseConfig 配置
var BaseConfig *AppConfig

// FilePath 配置文件路径
var FilePath string

// AppConfig 配置
type AppConfig struct {
	config2.ServiceConfigImpl
	MongoUrl    string `json:"mongoUrl"`
	MongoDbName string `json:"mongoDbName"`
	HttpUrl     string `json:"httpUrl"`
}

func init() {
	BaseConfig = &AppConfig{
		MongoDbName: "ugk-api",
	}
}

// InitConfigAndLog 初始化项目配置和日志
func InitConfigAndLog() {
	//1.配置文件路径
	configPath := flag.String("config", "D:\\Go\\ugk-server\\ugk-api\\config\\app_config_develop.json", "配置文件加载路径")
	flag.Parse()
	FilePath = *configPath
	BaseConfig.Reload()

	//2.关闭debug
	if "DEBUG" != BaseConfig.LogLevel {
		log.CloseDebug()
	}
	log.SetLogFile("log", "api")
}

// PathExists 判断一个文件是否存在
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// Reload 读取用户的配置文件
func (appConfig *AppConfig) Reload() {
	if confFileExists, _ := pathExists(FilePath); confFileExists != true {
		log.Warn("config file ", FilePath, "not find, use default config")
		return
	}
	//log.Info("加载配置文件：%v", FilePath)
	data, err := os.ReadFile(FilePath)
	if err != nil {
		panic(err)
	}
	//将json数据解析到struct中
	err = json.Unmarshal(data, appConfig)
	if err != nil {
		log.Error("%v", err)
		panic(err)
	}
}
