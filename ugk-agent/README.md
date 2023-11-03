# ugk-agent
&emsp;&emsp;创建Unity专用服务器docker进程。
在docker中运行可以通过镜像安装docker程序，挂载docker和宿主系统通信调用docker命令，
但是在容器中运行暂时没有好的办法获取宿主系统的CPU，磁盘，内存等监控信息。

## TODO
* ugk-agent-manager 需要ugk-api后台监控，网页查看（有多少ugk-agent，有多少游戏进程）？
