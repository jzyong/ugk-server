# UGK-Server

&emsp;&emsp;快节奏多人联网游戏Demo，UGK-Server：unity、go、kcp server 。
服务器使用微服务架构，服务器游戏逻辑需要物理碰撞、寻路的使用Unity、C#开发，其他使用Go开发。
对应客户端[ugk-client](https://github.com/jzyong/ugk-client)。 开发中......
![ugk-architecture](ugk-resource/img/ugk_architecture.png)

## 特性
* 自动化编排Docker容器(Unity dedicated server)
* 前后端分离
* 快照同步，状态同步预测插值

## 服务
### 通用
| 服务	                | 描述                        |
|--------------------|---------------------------|
| ugk-agent          | 执行unity服务器docker进程的创建销毁   |
| ugk-agent-manager  | 管理ugk-agent服务，为玩家房间分配游戏进程 |
| ugk-api            | HTTP Restful请求接口          |
| ugk-charge         | 充值                        |
| ugk-chat           | 聊天                        |
| ugk-common         | 公共逻辑封装                    |
| ugk-game           | 游戏微服务                     |
| ugk-gate           | 网关，消息转换                   |
| ugk-lobby          | 大厅，一般逻辑                   |
| ugk-login          | 登录、账号                     |
| ugk-message        | 协议消息                      |
| ugk-resource       | 文档、脚本                     |
| ugk-stress-testing | 压力测试客户端集群                 |

### 游戏
| 游戏	                                                               | 描述         |
|-------------------------------------------------------------------|------------|
| [game-galactic-kittens](ugk-game/game-galactic-kittens/README.md) | 2D多人飞船射击游戏 |



## 技术选择
* Unity、C# 客户端和服务器
* Go 服务器
* Kcp 网络通信 客户端和网关，网关和后端服务通信
* Protobuf+Grpc 内部网络通信
* Zookeeper 服务发现注册
* Mongodb,Redis 数据存储
* Docker、Jenkins 进行CI/CD




### 参考资料
#### 网络
* [Unity Multiplayer Networking](https://github.com/Unity-Technologies/com.unity.netcode.gameobjects)
* [FishNet](https://github.com/FirstGearGames/FishNet/)
* [Mirror](https://github.com/MirrorNetworking/Mirror)
* [Telepathy](https://github.com/vis2k/Telepathy) TCP前端
* [kcp2k](https://github.com/vis2k/kcp2k) unity前端
* [kcp-go](https://github.com/xtaci/kcp-go) go服务器
* [grpc](https://grpc.io/) 服务器之间通信
* [可靠UDP，KCP协议快在哪？](https://wetest.qq.com/lab/view/391.html)
#### 同步
* [Prediction,Reconciliation,Lag Compensation](https://www.gabrielgambetta.com/client-server-game-architecture.html)
* [Latency Compensating Methods in Client/Server In-game Protocol Design and Optimization](https://developer.valvesoftware.com/wiki/Latency_Compensating_Methods_in_Client/Server_In-game_Protocol_Design_and_Optimization)
* [无畏契约网络代码](https://technology.riotgames.com/news/peeking-valorants-netcode)
#### Unity
* [com.unity.multiplayer.samples.coop](https://github.com/Unity-Technologies/com.unity.multiplayer.samples.coop)3D rpg示例demo
* [com.unity.multiplayer.samples.bitesize](https://github.com/Unity-Technologies/com.unity.multiplayer.samples.bitesize)小游戏示例demo
* [GalacticKittens](https://github.com/UnityTechnologies/GalacticKittens) 2D示例demo
* [ECS-Network-Racing-Sample](https://github.com/Unity-Technologies/ECS-Network-Racing-Sample) ECS 赛车demo


## TODO
* 物理同步，动画同步
* aoi管理
* 客户端同步封装，参考NetworkTransform、NetworkAnimator、NetworkRigidbody？（参考Mirror）
* 服务器NetworkTransform只起配置作用，通过对象manager进行批量同步，子弹碰撞服务器检测，需要回滚
* 客户端和服务器配置都通过go服务器拉取
* 完整的GalacticKittensMatch流程
* 玩家进入游戏后端已经可以通过agent创建unity docker容器；下一步需要前端进入游戏场景，后端unity服务器开发，
* 场景消息同步定义，后端刷出小怪及客户端显示角色 ，给速度，客户端模拟，服务器出发了碰撞等在同步
* GalacticKittens客户端控制的使用快照插值同步，服务器控制的使用预测插值；Push coin使用快照插值同步
* 显示网络同步统计信息
* GalacticKittens 让服务器进行子弹命中判断，需要将客户端发射子弹的延迟加上



### 计划
* Websocket网络通信
* 压力测试客户端使用ugk-web开发界面（vue3）
* 添加聊天、排行、匹配、房间（Mirror）服务
* 广告、充值接取
* 断线重连(ugk-client,unity后端服务器与网关等)
* 使用c#开发导表等图形化工具
* 服务器unity提取公共包，unity的package
* ugk-client 弹窗增加tween动画
* 后台管理系统查看unity docker服务器


感谢
---------
<img src="https://resources.jetbrains.com/storage/products/company/brand/logos/jb_beam.png" width = "150" height = "150" div align=left />
