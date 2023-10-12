# UGK-Server

&emsp;&emsp;快节奏多人联网游戏Demo，UGK-Server：unity、go、kcp server 。
服务器使用微服务架构，服务器游戏逻辑需要物理碰撞、寻路的使用Unity、C#开发，其他使用Go开发。
对应客户端[ugk-client](https://github.com/jzyong/ugk-client)。 开发中......
![ugk-architecture](ugk-resource/img/ugk_architecture.png)


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
* Kcp 网络通信 忘记和网关，网关和后端服务通信
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
* 网络同步、延迟、插值、回退（延迟补偿），物理同步，动画同步，位置方向同步
* 场景同步消息，transform同步消息，aoi管理
* 场景对象同步分为消除和隐藏？
* 请求返回消息需要添加时间戳（延迟补偿需要？），客户端需要维护同步服务器时间戳？（Mirror怎样实现的？）[延迟补偿](https://www.gabrielgambetta.com/lag-compensation.html)
* 客户端封装NetworkBehavior？（参考Mirror）
* ugk-client 网络时间封装 NetworkTime
* 添加断开连接消息包处理逻辑
* 优先使用GalacticKittens 进行改造，服务器帧率30
* unity 专用服务器中实现和网关kcp相互连接（心跳），在Linux测试环境部署
* 网关服务器连接分为三个大类：lobby、功能微服务、游戏微服务
* 网关和大厅通信联调(完成登录流程，显示大厅游戏列表)
* 时间同步使用Mirror的（NetworkConnectionToClient（服务器）、Mathd、SnapshotInterpolation、NetworkTime(客户端)，NetworkClient_TimeInterpolation，ClientCube）时间同步计算只是在服务器？客户端通过心跳同步服务器时间NetworkPingMessage？
* unity double时间网络传输转换为毫秒使用int64传输
* GalacticKittensGame 与 Match、lobby、Gate的grpc连接创建管理,使用状态机管理房间状态
* agent获取系统CPU、内存、磁盘定时上报(3s), 用于agent-manager分配服务器启动游戏容器
* 使用agent创建GalacticKittensGame容器
* protobuf和grpc能否通过包管理引入
* unity 日志封装（写文件使用单独线程从队列中获取，可设置输出级别，debug，info，warn，error）
* 测试打包成Linux平台，docker镜像制作
* Unity 后台脚本构建项目，容器传参数到unity进程？
* go程序能否在docker容器中执行系统命令作用与宿主系统？
* GalacticKittensGame 断线重连网关，网关重启后，连接断开了
* docker 中程序获取宿主CPU，内存，磁盘等信息
* ugk-agent使用docker复制可执行程序，挂载到宿主系统，然后调用脚本让程序在宿主系统中运行？



### 计划
* Websocket网络通信
* ugk-agent 游戏服务进程管理
* 压力测试客户端使用ugk-web开发界面（vue3）
* 添加聊天、排行、匹配、房间（Mirror）服务
* 广告、充值接取
* 断线重连(ugk-client,unity后端服务器与网关等)
* 使用c#开发导表等图形化工具
* 服务器unity提取公共包，unity的package
* ugk-client 弹窗增加tween动画
