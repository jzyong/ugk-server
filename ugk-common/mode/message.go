package mode

import (
	"github.com/jzyong/ugk/common/constant"
	"sync"
)

//使用对象池，缓存byte和消息

// ugkMessagePool UgkMessage 对象池
var ugkMessagePool sync.Pool
var bytePool sync.Pool

func init() {
	ugkMessagePool.New = func() any {
		return &UgkMessage{}
	}
	bytePool.New = func() any {
		return make([]byte, constant.MessageLimit)
	}
}

// UgkMessage 内部自定义消息
type UgkMessage struct {
	MessageId uint32 //消息id
	Seq       uint32 //序列号
	Bytes     []byte //Byte数据
	TimeStamp int64  //时间戳
	Client    any    //客户端
}

// Reset 重置
func (msg *UgkMessage) Reset() {
	ReturnBytes(msg.Bytes)
	msg.MessageId = 0
	msg.Seq = 0
	msg.Bytes = nil
	msg.TimeStamp = 0
	msg.Client = nil
}

// ReturnUgkMessage 归还消息
func ReturnUgkMessage(msg *UgkMessage) {
	msg.Reset()
	ugkMessagePool.Put(msg)
}

// GetUgkMessage 获取Ugk消息
func GetUgkMessage() *UgkMessage {
	return ugkMessagePool.Get().(*UgkMessage)
}

// GetBytes 获取byte
func GetBytes() []byte {
	return bytePool.Get().([]byte)
}

// ReturnBytes 归还byte
func ReturnBytes(data []byte) {
	bytePool.Put(data)
}
