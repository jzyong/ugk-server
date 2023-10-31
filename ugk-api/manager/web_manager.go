package manager

import (
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	"github.com/jzyong/ugk/api/config"
	"sync"
	"time"
)

// WebManager web
type WebManager struct {
	util.DefaultModule
}

var webManager *WebManager
var webSingletonOnce sync.Once

func GetWebManager() *WebManager {
	webSingletonOnce.Do(func() {
		webManager = &WebManager{}
	})
	return webManager
}

func (m *WebManager) Init() error {
	log.Info("WebManager 初始化......")
	return nil
}

func (m *WebManager) Run() {
	go m.start()
}

func (m *WebManager) Stop() {
}

func (m *WebManager) start() {
	//设置配置
	web.BConfig.CopyRequestBody = true
	web.BConfig.WebConfig.Session.SessionOn = true
	web.BConfig.AppName = config.BaseConfig.Name
	web.BConfig.ServerName = config.BaseConfig.Name
	web.BConfig.Log.FileLineNum = true
	if config.BaseConfig.Profile == "develop" {
		web.BConfig.RunMode = web.DEV //开发模式消耗更多，如每次render都构建模板等
	} else {
		web.BConfig.RunMode = web.PROD
	}
	if config.BaseConfig.Profile == "online" {
		logs.SetLogger(logs.AdapterFile, `{"filename":"../log/api.log","maxsize":102400000,"maxbackup":7}`)
		loc, err := time.LoadLocation("America/Atikokan")
		if err != nil {
			log.Warn("修改时区错误：%v")
		}
		time.Local = loc
	}

	//web.ErrorController(&controller.ErrorController{})

	//http://localhost:5041
	//外测服不能绑定[ip:port]形式，只能[:port] beego的bug？还是系统配置问题？
	web.Run(config.BaseConfig.HttpUrl)
}
