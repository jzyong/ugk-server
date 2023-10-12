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
docker stop ugk-lobby-develop
docker rm ugk-lobby-develop
docker run -dit -p 3021:3021 --name ugk-lobby-develop -m 100M -e GO_OPTS=%goRunParams% ugk-lobby:develop
