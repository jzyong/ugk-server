# UGK-Server

多人联网游戏demo，UGK-Server：unity、go、kcp server 

## 服务

| 服务	        | 描述      |
|------------|---------|
| ugk-common | 公共逻辑封装  |
| ugk-game   | 游戏      |
| ugk-gate   | 网关，消息转换 |
| ugk-lobby  | 大厅，一般逻辑 |



## 技术选择
* Unity 客户端
* go 服务器
* kcp 网络通信


## TODO
* 添加kcp-go
* Unity客户端开发，使用kcp，添加插值，网络同步，使用unity的官方demo改造？
* 架构类似slots，登录大厅后，可选择多个小游戏进行玩耍
* ugk-client kcp 参考Mirror

