package controller

import (
	"github.com/jzyong/golib/log"
	"github.com/jzyong/ugk/api/manager"
)

// GateController 处理网关逻辑
type GateController struct {
	BaseController
}

// Url 获取网关地址
// http://localhost:3046/gate/url
func (m *GateController) Url() {
	url := manager.GetGateClientManager().HashModGate(m.Ctx.Input.IP())
	log.Debug("%v请求获取gate:%v", m.Ctx.Input.IP(), url)

	m.replayJson(url)
}
