rem build app
cd ../../../../ugk-gate
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -o ugk-gate

rem build image
docker image build -t ugk-gate:develop .

rem run docker
set goRunParams="-config /go/src/ugk-gate/config/app_config_develop2.json"
docker stop ugk-gate-develop2
docker rm ugk-gate-develop2
docker run -dit -p 3002:3002 -p 40020:40020/udp -p 5003:5003/udp --name ugk-gate-develop2 -m 100M -e GO_OPTS=%goRunParams% ugk-gate:develop
