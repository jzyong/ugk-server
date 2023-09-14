package config

import (
	"encoding/json"
	"flag"
	"github.com/jzyong/golib/log"
	"os"
)

// BaseConfig 配置
var BaseConfig *AppConfig

// FilePath 配置文件路径
var FilePath string

// AppConfig 配置
type AppConfig struct {
	Id       uint32 `json:"id"`       //服务器ID
	RpcUrl   string `json:"rpcUrl"`   //rpc 地址
	GateUrl  string `json:"gateUrl"`  //登录服地址，TODO通过zookeeper进行
	Profile  string `json:"profile"`  //个性化配置
	LogLevel string `json:"logLevel"` //日志级别
}

func init() {
	BaseConfig = &AppConfig{
		Id:       1,
		LogLevel: "DEBUG",
		Profile:  "develop",
	}
}

// 初始化项目配置和日志
func InitConfigAndLog() {
	//1.配置文件路径
	configPath := flag.String("config", "D:\\Go\\ugk-server\\ugk-lobby\\config\\app_config_develop.json", "配置文件加载路径")
	flag.Parse()
	FilePath = *configPath
	BaseConfig.Reload()

	//2.关闭debug
	if "DEBUG" != BaseConfig.LogLevel {
		log.CloseDebug()
	}
	log.SetLogFile("log", "lobby")
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
