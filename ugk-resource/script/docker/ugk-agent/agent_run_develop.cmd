rem build app
cd ../../../../ugk-agent
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -o ugk-agent

rem build image
docker image build -t ugk-agent:develop .

rem run docker
set goRunParams="-config /go/src/ugk-agent/config/app_config_develop.json"
docker stop ugk-agent-develop
docker rm ugk-agent-develop
docker run -dit -p 3031:3031 --name ugk-agent-develop -m 100M -e GO_OPTS=%goRunParams% ugk-agent:develop
