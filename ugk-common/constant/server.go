package constant

// ServiceName 服务器名称
type ServiceName string

const (
	Api           ServiceName = "api"
	Charge        ServiceName = "charge"
	Chat          ServiceName = "chat"
	Gate          ServiceName = "gate"
	Lobby         ServiceName = "lobby"
	Login         ServiceName = "login"
	StressTesting ServiceName = "stress-testing"
)
