package manager

import (
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	"github.com/xtaci/kcp-go/v5"
	"sync"
)

// UserManager 客户端-用户
type UserManager struct {
	util.DefaultModule
}

var userManager *UserManager
var userSingletonOnce sync.Once

func GetUserManager() *UserManager {
	userSingletonOnce.Do(func() {
		userManager = &UserManager{}
	})
	return userManager
}

func (m *UserManager) Init() error {
	log.Info("UserManager 初始化......")
	return nil
}

func (m *UserManager) Run() {
}

func (m *UserManager) Stop() {
}

// User 网关用户，每个用户一个routine接收消息，一个routine处理逻辑和发送消息，会不会创建太多routine？
type User struct {
	Id            int64           //唯一id
	ClientSession *kcp.UDPSession //客户端连接会话
	LobbySession  *kcp.UDPSession //大厅连接会话
	GameSession   *kcp.UDPSession //游戏连接会话
}
