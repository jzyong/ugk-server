# login

&emsp;&emsp;登录，账号管理

## 实现
* 网关和登录服交互协议较少，且需要在网关处理，因此全部用grpc实现通信

## TODO
* 暂时不使用数据库，写死四个账号 test1、test2、test3、test4，密码123，玩家id 1,2,3,4