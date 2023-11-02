rem build app
cd ../../../../ugk-lobby
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -o ugk-lobby

rem build image
docker image build -t ugk-lobby:develop .

rem run docker
set goRunParams="-config /go/src/ugk-lobby/config/app_config_develop.json"
docker stop ugk-lobby-develop2
docker rm ugk-lobby-develop2
docker run -dit -p 3022:3022 --name ugk-lobby-develop2 -m 100M -e GO_OPTS=%goRunParams% ugk-lobby:develop
