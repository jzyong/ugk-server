package manager

import (
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	"sync"
)

// LoginManager  入口
type LoginManager struct {
	util.DefaultModule
}

var loginManager *LoginManager
var loginSingletonOnce sync.Once

func GetLoginManager() *LoginManager {
	loginSingletonOnce.Do(func() {
		loginManager = &LoginManager{}
	})
	return loginManager
}

func (m *LoginManager) Init() error {
	log.Info("LoginManager 初始化......")
	return nil
}

func (m *LoginManager) Run() {
}

func (m *LoginManager) Stop() {
}
