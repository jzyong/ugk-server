package constant

import "time"

//网络相关常量

const (
	MTU                 = 1500             //mtu 长度
	MessageLimit        = 4000             // 消息长度限制
	WindowSize          = 4096             //窗口大小
	ClientHeartInterval = time.Second * 10 //客户端心跳时间
)
