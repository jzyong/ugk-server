package controller

// IndexController 首页
type IndexController struct {
	BaseController
}

// Index 首页
// http://localhost:3046
func (m *IndexController) Index() {
	m.replayJson("ugk-后台系统")
}
