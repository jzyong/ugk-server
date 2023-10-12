rem build app
cd ../../../../ugk-login
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -o ugk-login

rem build image
docker image build -t ugk-login:develop .

rem run docker
set goRunParams="-config /go/src/ugk-login/config/app_config_develop.json"
docker stop ugk-login-develop
docker rm ugk-login-develop
docker run -dit -p 3011:3011 --name ugk-login-develop -m 100M -e GO_OPTS=%goRunParams% ugk-login:develop
