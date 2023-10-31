package controller

import "github.com/jzyong/golib/log"

// GateController 处理网关逻辑
type GateController struct {
	BaseController
}

// Url 获取网关地址
// http://localhost:3046/gate/url
func (m *GateController) Url() {
	log.Debug("%v请求获取gate", m.Ctx.Input.IP())
	m.replayJson("ugk-请求获取gate")
}
