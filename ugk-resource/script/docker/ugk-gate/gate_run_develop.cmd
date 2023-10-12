rem build app
cd ../../../../ugk-gate
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -o ugk-gate

rem build image
docker image build -t ugk-gate:develop .

rem run docker
set goRunParams="-config /go/src/ugk-gate/config/app_config_develop.json"
docker stop ugk-gate-develop
docker rm ugk-gate-develop
docker run -dit -p 3001:3001 -p 5000:5000/udp -p 5001:5001/udp --name ugk-gate-develop -m 100M -e GO_OPTS=%goRunParams% ugk-gate:develop
