# GalacticKittens
&emsp;&emsp;改造unity官方网络示例[GalacticKittens](https://github.com/UnityTechnologies/GalacticKittens)。
需要用到unity物理相关的在unity中开发，其他都是用go开发。
![galactic-kittens_architecture](../../ugk-resource/img/game/galactic-kittens_architecture.png)

1. 玩家在大厅选择GalacticKittens进入匹配服`galactic-kittens-match`
2. 准备完成请求`agent-manager`分配服务器，调用对应服务器的`agent`创建unity游戏docker进程`galactic-kittens-game`
3. `galactic-kittens-game`主动创建与`lobby`,`galactic-kittens-match`,`gate`的网络连接
4. 玩家进入`galactic-kittens-game`unity游戏场景，向`lobby`请求玩家基础数据，然后进行游戏
5. 游戏结束请求`galactic-kittens-match`进行结算
6. `galactic-kittens-match`请求`lobby`进行数据存储更新
7. `galactic-kittens-match`请求`agent-manager`执行游戏进程`galactic-kittens-game`的结束销毁
8. 玩家返回游戏大厅



## 服务
### galactic-kittens-game

* Unity 游戏逻辑开发，需要物理碰撞监测，不需要渲染，摄像机等
* 每一个房间创建一个Unity进程，和gate、lobby、galactic-kittens-match服务连接
* Unity单线程执行

### galactic-kittens-match
* 房间匹配、创建、管理、结算，go开发
* 进入房间调用ugk-agent创建 galactic-kittens-game Docker进程，分配端口


## TODO 
* 窗口大小分辨率不同，影响相对位置
* 

