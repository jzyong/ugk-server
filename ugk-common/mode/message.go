package mode

//TODO 添加对象池

// UgkMessage 内部自定义消息
type UgkMessage struct {
	MessageId uint32 //消息id
	Seq       uint32 //序列号
	Bytes     []byte //Byte数据
	TimeStamp int64  //时间戳
}
