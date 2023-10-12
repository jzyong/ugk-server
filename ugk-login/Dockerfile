FROM alpine:latest

MAINTAINER  jzy
WORKDIR /go/src/ugk-login
# 必须配置，windows平台制作镜像可运行，centos不添加如下代码报 /bin/sh: /go/src/activity-service/activity-service: not found
RUN mkdir /lib64 \
    && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

COPY ugk-login /go/src/ugk-login/ugk-login
COPY config/*.json /go/src/ugk-login/config/

CMD /go/src/ugk-login/ugk-login ${GO_OPTS}
