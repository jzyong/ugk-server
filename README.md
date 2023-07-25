# UGK-Server

多人联网游戏demo，UGK-Server：unity、go、kcp server 。开发中......

## 服务

| 服务	         | 描述               |
|-------------|------------------|
| ugk-api     | HTTP Restful请求接口 |
| ugk-common  | 公共逻辑封装           |
| ugk-game    | 游戏微服务            |
| ugk-gate    | 网关，消息转换          |
| ugk-lobby   | 大厅，一般逻辑          |
| ugk-login   | 登录、账号            |
| ugk-message | 协议消息             |



## 技术选择
* Unity 客户端
* go 服务器
* kcp 网络通信

### 参考资料
#### 网络
* [Unity Multiplayer Networking](https://github.com/Unity-Technologies/com.unity.netcode.gameobjects)
* [FishNet](https://github.com/FirstGearGames/FishNet/)
* [Mirror](https://github.com/MirrorNetworking/Mirror)
* [Telepathy](https://github.com/vis2k/Telepathy) TCP前端
* [kcp2k](https://github.com/vis2k/kcp2k) unity前端
* [kcp-go](https://github.com/xtaci/kcp-go) go服务器
* [grpc](https://grpc.io/) 服务器之间通信
#### Unity
* [com.unity.multiplayer.samples.coop](https://github.com/Unity-Technologies/com.unity.multiplayer.samples.coop)3D rpg示例demo
* [com.unity.multiplayer.samples.bitesize](https://github.com/Unity-Technologies/com.unity.multiplayer.samples.bitesize)小游戏示例demo
* [GalacticKittens](https://github.com/UnityTechnologies/GalacticKittens) 2D示例demo
* [ECS-Network-Racing-Sample](https://github.com/Unity-Technologies/ECS-Network-Racing-Sample) ECS 赛车demo


## TODO
* 添加kcp-go
* Unity客户端开发，使用kcp，添加插值，网络同步，使用unity的官方demo改造？
* 架构类似slots，登录大厅后，可选择多个小游戏进行玩耍
* ugk-client kcp 参考Mirror，版本管理使用git-lfs
* 网关消息转发，消息ID右移20位switch1 2 4 8判断，奇数客户端请求，偶数服务器返回，可能模块100个消息，游戏模块1000个消息
* 网络同步、延迟、插值、回退（延迟补偿），物理同步，动画同步，位置方向同步

