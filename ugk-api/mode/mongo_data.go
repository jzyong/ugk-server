package mode

// ServerInfo 服务器全局信息
type ServerInfo struct {
	Id    int32 `_id` //唯一ID 1
	dirty bool  //数据是否修改
}
