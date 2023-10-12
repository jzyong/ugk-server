rem build app
cd ../../../../ugk-agent-manager
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -o ugk-agent-manager

rem build image
docker image build -t ugk-agent-manager:develop .

rem run docker
set goRunParams="-config /go/src/ugk-agent-manager/config/app_config_develop.json"
docker stop ugk-agent-manager-develop
docker rm ugk-agent-manager-develop
docker run -dit -p 3030:3030 --name ugk-agent-manager-develop -m 100M -e GO_OPTS=%goRunParams% ugk-agent-manager:develop
