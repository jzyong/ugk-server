# UGK-Server

快节奏多人联网游戏demo，UGK-Server：unity、go、kcp server 。开发中......

## 服务

| 服务	                | 描述               |
|--------------------|------------------|
| ugk-api            | HTTP Restful请求接口 |
| ugk-common         | 公共逻辑封装           |
| ugk-game           | 游戏微服务            |
| ugk-gate           | 网关，消息转换          |
| ugk-lobby          | 大厅，一般逻辑          |
| ugk-login          | 登录、账号            |
| ugk-message        | 协议消息             |
| ugk-resource       | 文档、脚本            |
| ugk-stress-testing | 压力测试客户端集群        |



## 技术选择
* Unity 客户端
* go 服务器
* kcp 网络通信
* protobuf+grpc

### 参考资料
#### 网络
* [Unity Multiplayer Networking](https://github.com/Unity-Technologies/com.unity.netcode.gameobjects)
* [FishNet](https://github.com/FirstGearGames/FishNet/)
* [Mirror](https://github.com/MirrorNetworking/Mirror)
* [Telepathy](https://github.com/vis2k/Telepathy) TCP前端
* [kcp2k](https://github.com/vis2k/kcp2k) unity前端
* [kcp-go](https://github.com/xtaci/kcp-go) go服务器
* [grpc](https://grpc.io/) 服务器之间通信
#### 同步
* [Prediction,Reconciliation,Lag Compensation](https://www.gabrielgambetta.com/client-server-game-architecture.html)
* [Latency Compensating Methods in Client/Server In-game Protocol Design and Optimization](https://developer.valvesoftware.com/wiki/Latency_Compensating_Methods_in_Client/Server_In-game_Protocol_Design_and_Optimization)
#### Unity
* [com.unity.multiplayer.samples.coop](https://github.com/Unity-Technologies/com.unity.multiplayer.samples.coop)3D rpg示例demo
* [com.unity.multiplayer.samples.bitesize](https://github.com/Unity-Technologies/com.unity.multiplayer.samples.bitesize)小游戏示例demo
* [GalacticKittens](https://github.com/UnityTechnologies/GalacticKittens) 2D示例demo
* [ECS-Network-Racing-Sample](https://github.com/Unity-Technologies/ECS-Network-Racing-Sample) ECS 赛车demo


## TODO
* Unity客户端开发，使用kcp，添加插值，网络同步，使用unity的官方demo改造？
* 架构类似slots，登录大厅后，可选择多个小游戏进行玩耍，一个玩家网关可以同时和多个后端服务通信
* ugk-client kcp 参考Mirror，版本管理使用git-lfs,unity使用2023版本
* 网络同步、延迟、插值、回退（延迟补偿），物理同步，动画同步，位置方向同步
* 压测客户端嵌入到ugk-client中？
* 压力测试客户端使用ugk-web开发界面（vue3）
* 服务器帧率30
* ugk-client 消息处理器参考slots-tool 同服务器一致（优化：注册方法而不是类，继承MessageHandler,注解加在方法上，方法不止传消息，还传上下文，参考grpc和springBoot实现）
    Mirror是手动注册的，自己实现使用注解扫描包自动注册
* 心跳2s每次，10s超时
* 场景同步消息，transform同步消息，aoi管理
* 场景对象同步分为消除和隐藏？
* 请求返回消息需要添加时间戳（延迟补偿需要？），客户端需要维护同步服务器时间戳？（Mirror怎样实现的？）[延迟补偿](https://www.gabrielgambetta.com/lag-compensation.html)
* 客户端封装NetworkBehavior？（参考Mirror）

