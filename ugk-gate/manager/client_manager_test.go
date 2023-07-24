package manager

import (
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
		for {
			data := time.Now().String()
			buf := make([]byte, len(data))
			log.Println("sent:", data)
			if _, err := sess.Write([]byte(data)); err == nil {
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
