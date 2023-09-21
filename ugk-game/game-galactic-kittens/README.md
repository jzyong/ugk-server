# GalacticKittens
&emsp;&emsp;改造unity官方网络示例[GalacticKittens](https://github.com/UnityTechnologies/GalacticKittens)。
需要用到unity物理相关的在unity中开发，其他都是用go开发。
![galactic-kittens_architecture](../../ugk-resource/img/game/galactic-kittens_architecture.png)

## 服务
### galactic-kittens-game

* Unity 游戏逻辑开发，需要物理碰撞监测，不需要渲染，摄像机等
* 每一个房间创建一个Unity进程，和gate、lobby、galactic-kittens-match服务连接
* Unity单线程执行

### galactic-kittens-match
* 房间匹配、创建、管理、结算，go开发
* 进入房间调用ugk-agent创建 galactic-kittens-game Docker进程，分配端口


## TODO 
* 需要处理连接多网关问题
* 每个玩家连接维护

