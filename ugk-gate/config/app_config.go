package config

import (
	"encoding/json"
	"flag"
	"github.com/jzyong/golib/log"
	"os"
)

// AppConfigManager 配置
var AppConfigManager *AppConfig

// FilePath 配置文件路径
var FilePath string

// AppConfig 配置
type AppConfig struct {
	Id            uint16   `json:"id"`            //服务器ID
	RpcUrl        string   `json:"rpcUrl"`        //rpc 地址
	ZookeeperUrls []string `json:"zookeeperUrls"` //zookeeper 地址
	Profile       string   `json:"profile"`       //个性化配置
	LogLevel      string   `json:"logLevel"`      //日志级别
	PublicIp      string   `json:"publicIp"`      //公网IP
	PrivateIp     string   `json:"privateIp"`     //内网IP
	ClientPort    uint16   `json:"clientPort"`    //客户端端口 KCP
	GamePort      uint16   `json:"gamePort"`      //内网游戏连接端口 TCP
}

func init() {
	AppConfigManager = &AppConfig{
		Id:       1,
		LogLevel: "DEBUG",
		Profile:  "develop",
	}
}

// 初始化项目配置和日志
func InitConfigAndLog() {
	//1.配置文件路径
	configPath := flag.String("config", "D:\\Go\\ugk-server\\ugk-gate\\config\\app_config_develop.json", "配置文件加载路径")
	flag.Parse()
	FilePath = *configPath
	AppConfigManager.Reload()

	//2.关闭debug
	if "DEBUG" != AppConfigManager.LogLevel {
		log.CloseDebug()
	}
	log.SetLogFile("log", "gate")
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
