package controller

import (
	"github.com/beego/beego/v2/server/web"
	"html"
	"strconv"
	"strings"
)

// BaseController 基础控制器
type BaseController struct {
	web.Controller
	actionName string //调用方法名称
}

// Prepare 初始化参数
func (s *BaseController) Prepare() {
	_, actionName := s.GetControllerAndAction()
	//log.Info("Prepare Controller=%v Action=%v", controllerName, actionName)
	s.actionName = strings.ToLower(actionName)

}

// GetAjaxPageParams ajax table分页参数
func (s *BaseController) GetAjaxPageParams() (start, limit int) {
	return s.GetIntNoErr("offset"), s.GetIntNoErr("limit")
}

// GetIntNoErr 去掉没有err返回值的int
func (s *BaseController) GetIntNoErr(key string, def ...int) int {
	strv := s.Ctx.Input.Query(key)
	if len(strv) == 0 && len(def) > 0 {
		return def[0]
	}
	val, _ := strconv.Atoi(strv)
	return val
}

func (s *BaseController) getEscapeString(key string) string {
	return html.EscapeString(s.GetString(key))
}

// AjaxTable ajax table返回
func (s *BaseController) AjaxTable(list interface{}, cnt int, recordsTotal int, kwargs map[string]interface{}) {
	json := make(map[string]interface{})
	json["rows"] = list
	json["total"] = recordsTotal
	if kwargs != nil {
		for k, v := range kwargs {
			if v != nil {
				json[k] = v
			}
		}
	}
	s.Data["json"] = json
	s.ServeJSON()
	s.StopRun()
}

// replayJson 返回json
func (s *BaseController) replayJson(data any) {
	s.Data["json"] = data
	s.ServeJSON()
}
