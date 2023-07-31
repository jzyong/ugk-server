package manager

import (
	"bytes"
	"encoding/binary"
	"github.com/jzyong/golib/util"
	"github.com/xtaci/kcp-go/v5"
	"io"
	"log"
	"testing"
	"time"
)

// 测试创建连接
func TestConnect(t *testing.T) {
	// dial to the echo server
	if sess, err := kcp.DialWithOptions("127.0.0.1:5000", nil, 0, 0); err == nil {
		sess.SetMtu(4096)
		var seq uint32 = 1
		for {
			data := time.Now().String()
			buf := make([]byte, len(data))
			log.Println("sent:", data)

			//`消息长度4+消息id4+序列号4+时间戳8+protobuf消息体`
			length := 20 + len(data)
			sendDatas := bytes.NewBuffer([]byte{}) //每个玩家可以缓存？
			var messageLength uint32 = uint32(length - 4)
			var messageId uint32 = 1
			binary.Write(sendDatas, binary.LittleEndian, messageLength)
			binary.Write(sendDatas, binary.LittleEndian, messageId)
			binary.Write(sendDatas, binary.LittleEndian, seq)
			binary.Write(sendDatas, binary.LittleEndian, util.CurrentTimeMillisecond())
			binary.Write(sendDatas, binary.LittleEndian, []byte(data))
			seq++

			//其他合并的 消息
			for i := 0; i < 50; i++ {
				length = 20 + len(data)
				messageLength = uint32(length - 4)
				messageId = uint32(1048577 + i) //2^20+1
				binary.Write(sendDatas, binary.LittleEndian, messageLength)
				binary.Write(sendDatas, binary.LittleEndian, messageId)
				binary.Write(sendDatas, binary.LittleEndian, seq)
				binary.Write(sendDatas, binary.LittleEndian, util.CurrentTimeMillisecond())
				binary.Write(sendDatas, binary.LittleEndian, []byte(data))
				seq++
			}

			if _, err := sess.Write(sendDatas.Bytes()); err == nil {
				// read back the data
				if _, err := io.ReadFull(sess, buf); err == nil {
					log.Println("recv:", string(buf))
				} else {
					log.Fatal(err)
				}
			} else {
				log.Fatal(err)
			}
			time.Sleep(time.Second)
		}
	} else {
		log.Fatal(err)
	}
}
