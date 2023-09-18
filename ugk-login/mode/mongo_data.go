package mode

// Account 账号
type Account struct {
	Id       string `_id`      //唯一id，暂时用账号
	Password string `password` //密码
	PlayerId int64  `playerId` //玩家id
}

// ServerInfo 服务器全局信息
type ServerInfo struct {
	Id       int32 `_id`      //唯一ID 1
	PlayerId int64 `playerId` //注册账号玩家自增ID
	dirty    bool  //数据是否修改
}
