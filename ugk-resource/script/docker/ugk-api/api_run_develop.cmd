rem build app
cd ../../../../ugk-api
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -o ugk-api

rem build image
docker image build -t ugk-api:develop .

rem run docker
set goRunParams="-config /go/src/ugk-api/config/app_config_develop.json"
docker stop ugk-api-develop
docker rm ugk-api-develop
docker run -dit -p 3041:3041 -p 3046:3046 --name ugk-api-develop -m 100M -e GO_OPTS=%goRunParams% ugk-api:develop
