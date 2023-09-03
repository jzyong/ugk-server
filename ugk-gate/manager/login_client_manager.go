package manager

import (
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	"github.com/jzyong/ugk/gate/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"sync"
)

// LoginClientManager 登录客户端

type LoginClientManager struct {
	util.DefaultModule
	ClientConn *grpc.ClientConn
}

var loginClientManager *LoginClientManager
var loginClientSingletonOnce sync.Once

func GetLoginClientManager() *LoginClientManager {
	loginClientSingletonOnce.Do(func() {
		loginClientManager = &LoginClientManager{}
	})
	return loginClientManager
}

func (m *LoginClientManager) Init() error {
	log.Info("LoginClientManager 初始化......")

	m.Start()
	return nil
}

func (m *LoginClientManager) Run() {
}

func (m *LoginClientManager) Stop() {
}

// Start TODO 暂时写死，后面从zookeeper中获取
func (m *LoginClientManager) Start() {
	go func() {
		loginUrl := config.AppConfigManager.LoginUrl
		dialOption := grpc.WithTransportCredentials(insecure.NewCredentials())
		clientConn, err := grpc.Dial(loginUrl, dialOption)
		if err != nil {
			log.Error("%v", err)
			return
		}
		m.ClientConn = clientConn
	}()
}
