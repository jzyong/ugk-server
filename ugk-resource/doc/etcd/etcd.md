# etcd
[官方文档](https://etcd.io/docs/v3.5/)

## 安装
[Docker安装参考文档](https://hub.docker.com/r/bitnami/etcd)
```shell
docker run -dit --name Etcd --env ALLOW_NONE_AUTHENTICATION=yes --publish 2379:2379 --publish 2380:2380 bitnami/etcd
```