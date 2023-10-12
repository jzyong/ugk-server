FROM alpine:latest

MAINTAINER  jzy
# 必须配置，windows平台制作镜像可运行，centos不添加如下代码报 /bin/sh: /go/src/activity-service/activity-service: not found
RUN mkdir /lib64 \
    && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
# 安装docker
RUN apk update
RUN apk add docker
# RUN service docker start
