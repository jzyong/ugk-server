# gate

### 协议
#### 客户端到Gate
`消息长度4+消息id4+序列号4+时间戳8+protobuf消息体`
#### Gate到游戏服
`消息长度4+玩家ID8+消息id4+序列号4+时间戳8+protobuf消息体`