# gate

### 协议
#### 客户端到Gate
`消息长度4+消息id4+序列号4+时间戳8+protobuf消息体`  
#### Gate到游戏服
`消息长度4+玩家ID8+消息id4+序列号4+时间戳8+protobuf消息体`

### 实现
* 客户端基于Mirror-kcp2k修改
* 只实现可靠传输

#### Mirror kcp
* 定义了`Handshake、Ping、Data、Disconnect` 等不同类型的消息封包模式
* 封包是`可靠消息标识1+cookie4+消息内容`，cookie判断每个消息是否合法；消息内容还有子封包`时间戳+消息ID`
* 心跳消息自定义实现在KcpPeer中
* 可靠消息经过Kcp处理，不可靠消息直接使用udp
* 封包实现功能较多，很多回调函数，较复杂