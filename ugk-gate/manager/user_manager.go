package manager

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	"github.com/jzyong/ugk/common/config"
	"github.com/jzyong/ugk/common/mode"
	config2 "github.com/jzyong/ugk/gate/config"
	"github.com/jzyong/ugk/message/message"
	"github.com/xtaci/kcp-go/v5"
	"google.golang.org/protobuf/proto"
	"sync"
	"time"
)

// UserManager 客户端-用户
type UserManager struct {
	util.DefaultModule
	//ipUsers map[string]*User //IP用户
	idUsers map[int64]*User //登录后的玩家ID用户
	mutex   sync.RWMutex
}

var userManager *UserManager
var userSingletonOnce sync.Once

func GetUserManager() *UserManager {
	userSingletonOnce.Do(func() {
		userManager = &UserManager{
			//ipUsers: make(map[string]*User, 1024),
			idUsers: make(map[int64]*User, 1024),
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

func (m *UserManager) AddUser(user *User) {
	defer m.mutex.Unlock()
	m.mutex.Lock()
	m.idUsers[user.Id] = user
}

func (m *UserManager) GetUser(id int64) *User {
	defer m.mutex.RUnlock()
	m.mutex.RLock()
	return m.idUsers[id]
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
	Id               int64                 //唯一id
	ClientSession    *kcp.UDPSession       //客户端连接会话
	LobbyClient      *GameKcpClient        //大厅连接会话
	GameClient       *GameKcpClient        //游戏连接会话
	SendBuffer       *bytes.Buffer         //发送缓冲区，消息批量合并，单线程调用
	ReceiveBuffer    *bytes.Buffer         //接收缓冲区,
	ReceiveBytes     chan []byte           //接收到的消息
	GameMessages     chan *mode.UgkMessage //游戏服返回的消息 proto字节流、消息id、序列号、客户端
	ReceiveReadCache []byte                // 接收端读取Byte缓存
	CloseChan        chan struct{}         //离线等关闭Chan
	State            UserState             //用户状态
	HeartTime        time.Time             //心跳时间
}

// NewUser 构建用户
func NewUser(clientSession *kcp.UDPSession) *User {
	user := &User{ClientSession: clientSession,
		SendBuffer:       bytes.NewBuffer([]byte{}),
		ReceiveBuffer:    bytes.NewBuffer([]byte{}),
		ReceiveBytes:     make(chan []byte, 1024),
		GameMessages:     make(chan *mode.UgkMessage, 1024), //proto字节流、消息id、序列号、客户端
		ReceiveReadCache: make([]byte, 1500),                //每次最多读取1500-消息头字节
		CloseChan:        make(chan struct{}),
		State:            NetWorkActive,
		HeartTime:        time.Now(),
	}
	//只在此处添加
	//GetUserManager().ipUsers[clientSession.RemoteAddr().String()] = user
	go user.run()
	return user
}

// 玩家routine运行
func (user *User) run() {
	messageMergeTicker := time.Tick(time.Millisecond * 10) //最小10ms进行一次心跳
	secondTicker := time.Tick(time.Second)
	for {
		select {
		case receiveByte := <-user.ReceiveBytes: //用户消息分发
			user.messageDistribute(receiveByte)
		case gameMessage := <-user.GameMessages: //处理服务器返回消息
			handFunc := ServerHandlers[gameMessage.MessageId]
			handFunc(user, gameMessage.Client.(*GameKcpClient), gameMessage)
			mode.ReturnUgkMessage(gameMessage)
		case <-user.CloseChan: //关闭用户chan
			log.Info("%v %v chan关闭", user.Id, user.ClientSession.RemoteAddr().String())
			user.State = Offline
			return
		case <-messageMergeTicker:
			user.batchTransmitToClient()
		case <-secondTicker:
			user.secondUpdate()
		}

	}
}

// 玩家更新逻辑
func (user *User) secondUpdate() {
	// 心跳监测
	if time.Now().Sub(user.HeartTime) > config.ClientHeartInterval {
		channelInactive(user, errors.New(fmt.Sprintf("心跳超时%f", time.Now().Sub(user.HeartTime).Seconds())))
	}

}

func (user *User) messageDistribute(data []byte) {
	defer mode.ReturnBytes(data)
	user.HeartTime = time.Now()
	//`消息长度4+消息id4+序列号4+时间戳8+protobuf消息体`
	//小端
	messageId := uint32(data[4]) | uint32(data[5])<<8 | uint32(data[6])<<16 | uint32(data[7])<<24
	handlerServer := messageId >> 20 //0截取本地、1lobby、2功能微服务，3游戏微服务
	//log.Info("%v 消息Id=%v 处理服务=%v", user.Id, messageId, handlerServer)
	//除了大厅（1<<20+编号）和unity服务器（4<<20+编号），其他服务的消息需要zookeeper进行注册转发[(2<<20)<消息id<(4<<20)]
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
		//protoData := make([]byte, messageLength-16) //使用对象池,减少内存分配回收
		protoData := mode.GetBytes()[:messageLength-16]
		if err := binary.Read(dataReader, binary.LittleEndian, &protoData); err != nil {
			channelInactive(user, errors.New("读取消息proto数据错误"))
			return
		}
		ugkMessage := mode.GetUgkMessage()
		ugkMessage.MessageId = messageId
		ugkMessage.Seq = seq
		ugkMessage.Bytes = protoData
		ugkMessage.TimeStamp = timeStamp
		defer mode.ReturnUgkMessage(ugkMessage)
		handFunc(user, ugkMessage)
		break
	case 1: // 大厅
		user.TransmitToLobby(data, messageId)
	case 2: //公共服务
		user.TransmitToBackend(data, messageId)
		log.Trace("收到消息：%v", messageId)
	case 3: //子游戏go服务
		user.TransmitToBackend(data, messageId)
		log.Trace("收到消息：%v", messageId)
	default:
		// 转发到玩家绑定的游戏服
		user.TransmitToGame(data, messageId)
		log.Trace("%d - %s 收到未知消息mid=%d", user.Id, user.ClientSession.RemoteAddr().String(), messageId)
	}
}

// TransmitToLobby 转发到大厅
func (user *User) TransmitToLobby(clientData []byte, messageId uint32) error {
	if user.LobbyClient == nil {
		log.Warn("玩家：%d未分配大厅", user.Id)
		return errors.New(fmt.Sprintf("玩家：%d未分配大厅", user.Id))
	}
	bytes, err := user.toGameBytes(clientData, messageId)
	if err != nil {
		return err
	}
	_, err = user.LobbyClient.UdpSession.Write(bytes)
	if err != nil {
		log.Error("%d - %s 发送消息 %d 失败：%v", user.Id, user.ClientSession.RemoteAddr().String(), messageId, err)
		return err
	}
	return nil
}

// TransmitToBackend 转发后端服务
func (user *User) TransmitToBackend(clientData []byte, messageId uint32) error {
	serviceName := GetServerManager().MessageIdService[messageId]
	if serviceName == "" {
		msg := fmt.Sprintf("玩家%d 消息：%v注册服务未找到", user.Id, messageId)
		log.Warn(msg)
		return errors.New(msg)
	}
	// 后期需要跟进id获取，否则user缓存客户端，减少加锁获取
	client := GetServerManager().GetGameClient(serviceName, 0)
	if client == nil {
		msg := fmt.Sprintf("玩家：%d 消息%v 获取服务%v失败", user.Id, messageId, serviceName)
		log.Warn(msg)
		return errors.New(msg)
	}
	bytes, err := user.toGameBytes(clientData, messageId)
	if err != nil {
		return err
	}
	_, err = client.UdpSession.Write(bytes)
	if err != nil {
		log.Error("%d - %s 发送消息 %d 失败：%v", user.Id, user.ClientSession.RemoteAddr().String(), messageId, err)
		return err
	}
	return nil
}

// TransmitToGame 转发unity游戏服务
func (user *User) TransmitToGame(clientData []byte, messageId uint32) error {
	if user.GameClient == nil {
		msg := fmt.Sprintf("玩家%d 消息：%v游戏服务未注册", user.Id, messageId)
		log.Warn(msg)
		return errors.New(msg)
	}
	bytes, err := user.toGameBytes(clientData, messageId)
	if err != nil {
		return err
	}
	_, err = user.GameClient.UdpSession.Write(bytes)
	if err != nil {
		log.Error("%d - %s 发送消息 %d 失败：%v", user.Id, user.ClientSession.RemoteAddr().String(), messageId, err)
		return err
	}
	return nil
}

// 客户端byte转换为服务器byte流
func (user *User) toGameBytes(clientData []byte, messageId uint32) ([]byte, error) {
	clientLength := len(clientData)
	if clientLength > config.MessageLimit {
		log.Error("%d - %s 发送消息 %d  失败：消息太长", user.Id, user.ClientSession.RemoteAddr().String(), messageId)
		return nil, errors.New("消息超长")
	}

	//消息长度4+玩家ID+消息id4+序列号4+时间戳8+protobuf消息体
	buffer := bytes.NewBuffer([]byte{})
	//写dataLen 不包含自身长度,比客户端长度多8字节玩家id clientLength+8-4
	if err := binary.Write(buffer, binary.LittleEndian, uint32(clientLength+4)); err != nil {
		log.Error("%d - %s 发送消息 %d 失败：%v", user.Id, user.ClientSession.RemoteAddr().String(), messageId, err)
		return nil, err
	}
	//写玩家ID
	if err := binary.Write(buffer, binary.LittleEndian, user.Id); err != nil {
		log.Error("%d - %s 发送消息 %d 失败：%v", user.Id, user.ClientSession.RemoteAddr().String(), messageId, err)
		return nil, err
	}
	data := clientData[4:]
	//写data数据
	if err := binary.Write(buffer, binary.LittleEndian, data); err != nil {
		log.Error("%d - %s 发送消息 %d 失败：%v", user.Id, user.ClientSession.RemoteAddr().String(), messageId, err)
		return nil, err
	}
	return buffer.Bytes(), nil
}

// TransmitToClient 转发到客户端
func (user *User) TransmitToClient(gameData []byte, messageId uint32) error {
	if user.ClientSession == nil {
		log.Warn("玩家：%d无客户端连接", user.Id)
		return errors.New(fmt.Sprintf("玩家：%d无客户端连接", user.Id))
	}
	gameLength := len(gameData)
	if gameLength > config.MessageLimit {
		log.Error("%d - %s 发送消息 %d  失败：消息太长", user.Id, user.ClientSession.RemoteAddr().String(), messageId)
		return errors.New("消息超长")
	}

	//直接发送
	if config2.BaseConfig.BatchMessage == false {
		return user.immediateTransmitToClient(gameLength, messageId, gameData)
		// 消息合并批量发送，将小于MTU的消息合并，大于的直接发送，减少丢包重传率，减少IO
	} else {
		//当前消息大于mtu，先将老的发送，然后立即发送当前消息
		if gameLength > config.MTU {
			user.batchTransmitToClient()
			user.immediateTransmitToClient(gameLength, messageId, gameData)
			//将消息写入缓存
		} else {
			//写dataLen 不包含自身长度,比客户端长度多8字节玩家id gameLength-8-4
			if err := binary.Write(user.SendBuffer, binary.LittleEndian, uint32(gameLength-12)); err != nil {
				log.Error("%d - %s 发送消息 %d 失败：%v", user.Id, user.ClientSession.RemoteAddr().String(), messageId, err)
				return err
			}
			data := gameData[12:]
			//写data数据
			if err := binary.Write(user.SendBuffer, binary.LittleEndian, data); err != nil {
				log.Error("%d - %s 发送消息 %d 失败：%v", user.Id, user.ClientSession.RemoteAddr().String(), messageId, err)
				return err
			}

			//已经超过mtu ，直接推送
			if user.SendBuffer.Len() > config.MTU {
				user.batchTransmitToClient()
			}
		}
	}

	return nil
}

// 立即发送客户端
func (user *User) immediateTransmitToClient(gameLength int, messageId uint32, gameData []byte) error {
	//消息长度4+玩家ID+消息id4+序列号4+时间戳8+protobuf消息体 ==》消息长度4+消息id4+序列号4+时间戳8+protobuf消息体
	buffer := user.SendBuffer
	//写dataLen 不包含自身长度,比客户端长度多8字节玩家id gameLength-8-4
	if err := binary.Write(buffer, binary.LittleEndian, uint32(gameLength-12)); err != nil {
		log.Error("%d - %s 发送消息 %d 失败：%v", user.Id, user.ClientSession.RemoteAddr().String(), messageId, err)
		return err
	}

	data := gameData[12:]
	//写data数据
	if err := binary.Write(buffer, binary.LittleEndian, data); err != nil {
		log.Error("%d - %s 发送消息 %d 失败：%v", user.Id, user.ClientSession.RemoteAddr().String(), messageId, err)
		return err
	}

	_, err := user.ClientSession.Write(buffer.Bytes())
	buffer.Reset()
	if err != nil {
		log.Error("%d - %s 发送消息 %d 失败：%v", user.Id, user.ClientSession.RemoteAddr().String(), messageId, err)
		return err
	}

	return nil
}

// 批量发送缓存的消息到客户端
func (user *User) batchTransmitToClient() error {
	if user.SendBuffer.Len() < 1 {
		return nil
	}
	_, err := user.ClientSession.Write(user.SendBuffer.Bytes())
	user.SendBuffer.Reset()
	if err != nil {
		log.Error("%d - %s 批量发送消息失败：%v", user.Id, user.ClientSession.RemoteAddr().String(), err)
		return err
	}

	return nil
}

// SendToClient 发送消息到客户端
func (user *User) SendToClient(mid message.MID, msg proto.Message, seq uint32) error {
	protoData, err := proto.Marshal(msg)
	if err != nil {
		log.Error("%d - %s 发送消息 %d 失败：%v", user.Id, user.ClientSession.RemoteAddr().String(), mid, err)
		return err
	}
	protoLength := len(protoData)
	if protoLength > config.MessageLimit {
		log.Error("%d - %s 发送消息 %d 失败：%v", user.Id, user.ClientSession.RemoteAddr().String(), mid, err)
		return errors.New("消息超长")
	}

	//消息长度4+消息id4+序列号4+时间戳8+protobuf消息体
	buffer := bytes.NewBuffer([]byte{})
	//写dataLen 不包含自身长度
	if err := binary.Write(buffer, binary.LittleEndian, uint32(16+protoLength)); err != nil {
		log.Error("%d - %s 发送消息 %d 失败：%v", user.Id, user.ClientSession.RemoteAddr().String(), mid, err)
		return err
	}
	//写msgID
	if err := binary.Write(buffer, binary.LittleEndian, uint32(mid)); err != nil {
		log.Error("%d - %s 发送消息 %d 失败：%v", user.Id, user.ClientSession.RemoteAddr().String(), mid, err)
		return err
	}
	//写 序列号
	if err := binary.Write(buffer, binary.LittleEndian, seq); err != nil {
		log.Error("%d - %s 发送消息 %d 失败：%v", user.Id, user.ClientSession.RemoteAddr().String(), mid, err)
		return err
	}
	//写时间戳
	if err := binary.Write(buffer, binary.LittleEndian, util.CurrentTimeMillisecond()); err != nil {
		log.Error("%d - %s 发送消息 %d 失败：%v", user.Id, user.ClientSession.RemoteAddr().String(), mid, err)
		return err
	}
	//写data数据
	if err := binary.Write(buffer, binary.LittleEndian, protoData); err != nil {
		log.Error("%d - %s 发送消息 %d 失败：%v", user.Id, user.ClientSession.RemoteAddr().String(), mid, err)
		return err
	}
	_, err = user.ClientSession.Write(buffer.Bytes())
	if err != nil {
		log.Error("%d - %s 发送消息 %d 失败：%v", user.Id, user.ClientSession.RemoteAddr().String(), mid, err)
		return err
	}
	return nil
}
