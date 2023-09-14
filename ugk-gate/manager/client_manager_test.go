package manager

import (
	"bytes"
	"encoding/binary"
	"github.com/jzyong/golib/util"
	"github.com/jzyong/ugk/message/message"
	"github.com/xtaci/kcp-go/v5"
	"google.golang.org/protobuf/proto"
	"log"
	"testing"
	"time"
)

// 测试创建连接
func TestConnect(t *testing.T) {
	// dial to the echo server

	if sess, err := kcp.DialWithOptions("127.0.0.1:5000", nil, 0, 0); err == nil {
		sess.SetMtu(config.MTU)
		sess.SetStreamMode(true) //true 流模式：使每个段数据填充满,避免浪费
		sess.SetNoDelay(1, 10, 2, 1)
		var seq uint32 = 1
		for {
			data := time.Now().String()
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

			// 循环发送几次
			sendBytes := sendDatas.Bytes()
			for i := 0; i < 3; i++ {
				if _, err := sess.Write(sendBytes); err == nil {
					//buf := make([]byte, len(data))
					////read back the data
					//if _, err := io.ReadFull(sess, buf); err == nil {
					//	log.Println("recv:", string(buf))
					//} else {
					//	log.Fatal(err)
					//}
				} else {
					log.Fatal(err)
				}
				time.Sleep(time.Second * 3)
				seq++
				sendMsg(sess, &message.HeartRequest{}, message.MID_HeartReq, seq)
			}

			time.Sleep(time.Minute * 60)
		}
	} else {
		log.Fatal(err)
	}
}

// 发送消息
func sendMsg(session *kcp.UDPSession, message proto.Message, mid message.MID, seq uint32) {
	//`消息长度4+消息id4+序列号4+时间戳8+protobuf消息体`
	var data, err = proto.Marshal(message)
	if err != nil {
		log.Fatal(err)
		return
	}
	length := 20 + len(data)
	sendDatas := bytes.NewBuffer([]byte{}) //每个玩家可以缓存？
	var messageLength = uint32(length - 4)
	binary.Write(sendDatas, binary.LittleEndian, messageLength)
	binary.Write(sendDatas, binary.LittleEndian, int32(mid))
	binary.Write(sendDatas, binary.LittleEndian, seq)
	binary.Write(sendDatas, binary.LittleEndian, util.CurrentTimeMillisecond())
	binary.Write(sendDatas, binary.LittleEndian, []byte(data))
	seq++
	if _, err := session.Write(sendDatas.Bytes()); err != nil {
		log.Fatal(err)
	}
}
