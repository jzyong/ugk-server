#FROM alpine:latest 镜像不行，缺少组件
FROM oraclelinux:8-slim

# 后面单独复制，减少更新文件大小,添加启动参数传入
COPY Build/Linux /usr/src/game

WORKDIR /usr/src/game

CMD ./game.x86_64
