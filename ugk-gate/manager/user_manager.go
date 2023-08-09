package manager

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	"github.com/jzyong/ugk/common/constant"
	"github.com/xtaci/kcp-go/v5"
	"sync"
	"time"
)

// UserManager 客户端-用户
type UserManager struct {
	util.DefaultModule
	IpUsers map[string]*User //IP用户
	IdUsers map[int64]*User  //登录后的玩家ID用户
}

var userManager *UserManager
var userSingletonOnce sync.Once

func GetUserManager() *UserManager {
	userSingletonOnce.Do(func() {
		userManager = &UserManager{
			IpUsers: make(map[string]*User, 1024),
			IdUsers: make(map[int64]*User, 1024),
		}
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

// UserState 用户状态
type UserState int

const (
	NetWorkActive UserState = 0 //网络激活
	Login         UserState = 1 //已登录
	Offline       UserState = 2 // 已离线
)

// User 网关用户，每个用户一个routine接收消息，一个routine处理逻辑和发送消息，会不会创建太多routine？
type User struct {
	Id               int64           //唯一id
	ClientSession    *kcp.UDPSession //客户端连接会话
	LobbySession     *kcp.UDPSession //大厅连接会话 //TODO 使用服务模式，不写死，可以连接多个后端服务
	GameSession      *kcp.UDPSession //游戏连接会话
	SendBuffer       *bytes.Buffer   //发送缓冲区，单线程调用
	ReceiveBuffer    *bytes.Buffer   //接收缓冲区
	ReceiveBytes     chan []byte     //接收到的消息
	ReceiveReadCache []byte          // 接收端读取Byte缓存
	CloseChan        chan struct{}   //离线等关闭Chan
	State            UserState       //用户状态
	HeartTime        time.Time       //心跳时间
}

// NewUser 构建用户
func NewUser(clientSession *kcp.UDPSession) *User {
	user := &User{ClientSession: clientSession,
		SendBuffer:       bytes.NewBuffer([]byte{}),
		ReceiveBuffer:    bytes.NewBuffer([]byte{}),
		ReceiveBytes:     make(chan []byte, 1024),
		ReceiveReadCache: make([]byte, 1500), //每次最多读取1500-消息头字节
		CloseChan:        make(chan struct{}),
		State:            NetWorkActive,
	}
	//只在此处添加
	GetUserManager().IpUsers[clientSession.RemoteAddr().String()] = user
	go user.run()
	return user
}

// 玩家routine运行
func (user *User) run() {
	messageMergeTicker := time.Tick(time.Millisecond * 10) //最小10ms进行一次心跳
	secondTicker := time.Tick(time.Second)
	for {
		select {
		case receiveByte := <-user.ReceiveBytes:
			user.messageDistribute(receiveByte)
		case <-user.CloseChan:
			log.Info("%v %v chan关闭", user.Id, user.ClientSession.RemoteAddr().String())
			user.State = Offline
			return
		case <-messageMergeTicker:
			user.sendMergeMessage()
		case <-secondTicker:
			user.secondUpdate()
		}

	}
}

// 玩家更新逻辑
func (user *User) secondUpdate() {
	// 心跳监测
	if time.Now().Sub(user.HeartTime) > constant.ClientHeartInterval {
		channelInactive(user, errors.New(fmt.Sprintf("心跳超时%f", time.Now().Sub(user.HeartTime).Seconds())))
	}

}

// 将缓存的消息合并批量发送
func (user *User) sendMergeMessage() {
	//TODO 合并消息包等逻辑
}

func (user *User) messageDistribute(data []byte) {
	user.HeartTime = time.Now()
	//`消息长度4+消息id4+序列号4+时间戳8+protobuf消息体`
	//小端
	messageId := uint32(data[4]) | uint32(data[5])<<8 | uint32(data[6])<<16 | uint32(data[7])<<24
	handlerServer := messageId >> 20 //0本地处理，1lobby，2游戏
	//log.Info("%v 消息Id=%v 处理服务=%v", user.Id, messageId, handlerServer)
	switch handlerServer {
	case 0:
		// 本地逻辑处理
		handFunc := ClientHandlers[messageId]
		if handFunc == nil {
			log.Warn("%d mid=%d执行失败，未找到执行函数", user.Id, messageId)
			return
		}

		//截取消息
		dataReader := bytes.NewReader(data)
		var messageLength int32
		if err := binary.Read(dataReader, binary.LittleEndian, &messageLength); err != nil {
			channelInactive(user, errors.New("读取消息长度错误"))
			return
		}
		if err := binary.Read(dataReader, binary.LittleEndian, &messageId); err != nil {
			channelInactive(user, errors.New("读取消息ID错误"))
			return
		}
		var seq uint32
		if err := binary.Read(dataReader, binary.LittleEndian, &seq); err != nil {
			channelInactive(user, errors.New("读取消息seq错误"))
			return
		}
		var timeStamp int64
		if err := binary.Read(dataReader, binary.LittleEndian, &timeStamp); err != nil {
			channelInactive(user, errors.New("读取消息timeStamp错误"))
			return
		}
		protoData := make([]byte, messageLength-16)
		if err := binary.Read(dataReader, binary.LittleEndian, &protoData); err != nil {
			channelInactive(user, errors.New("读取消息proto数据错误"))
			return
		}
		handFunc(user, protoData, seq, timeStamp)

	case 1:

	case 2:
	default:
		log.Warn("%d - %s 收到未知消息mid=%d", user.Id, user.ClientSession.RemoteAddr().String(), messageId)
	}

}
