package controller

import "github.com/beego/beego/v2/server/web"

// RegisterController 注册处理器
func RegisterController() {
	//路由注册
	indexController := &IndexController{}
	web.AutoRouter(indexController)
	web.Router("/", indexController, "*:Index")
	web.AutoRouter(&GateController{})
}
