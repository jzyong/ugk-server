package config

import (
	"time"
)

//网络相关常量

const (
	MTU                 = 1200             //mtu 长度 ，(reduced to 1200 to fit all cases: https://en.wikipedia.org/wiki/Maximum_transmission_unit ; steam uses 1200 too!)
	MessageLimit        = 4000             // 消息长度限制
	WindowSize          = 4096             //窗口大小
	ClientHeartInterval = time.Minute * 5  //客户端心跳时间  //Unity 默认超时为10s，最大重连次数为10次 https://docs-multiplayer.unity3d.com/netcode/current/learn/bossroom/optimizing-bossroom/#disconnect-timeout
	ServerHeartInterval = time.Minute * 10 //客户端心跳时间
)
