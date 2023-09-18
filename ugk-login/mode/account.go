package mode

// Account 账号
type Account struct {
	Id       string `_id`      //唯一id，暂时用账号
	Password string `password` //密码
	PlayerId int64  `playerId` //玩家id
}
