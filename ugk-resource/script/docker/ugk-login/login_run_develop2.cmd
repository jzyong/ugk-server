rem build app
cd ../../../../ugk-login
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -o ugk-login

rem build image
docker image build -t ugk-login:develop .

rem run docker
set goRunParams="-config /go/src/ugk-login/config/app_config_develop2.json"
docker stop ugk-login-develop2
docker rm ugk-login-develop2
docker run -dit -p 3012:3012 --name ugk-login-develop2 -m 100M -e GO_OPTS=%goRunParams% ugk-login:develop
