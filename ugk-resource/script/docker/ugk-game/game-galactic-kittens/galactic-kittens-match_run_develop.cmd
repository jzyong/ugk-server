rem build app
cd ../../../../../ugk-game/game-galactic-kittens/galactic-kittens-match
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -o galactic-kittens-match

rem build image
docker image build -t game-galactic-kittens-match:develop .

rem run docker
set goRunParams="-config /go/src/galactic-kittens-match/config/app_config_develop.json"
docker stop game-galactic-kittens-match-develop
docker rm game-galactic-kittens-match-develop
docker run -dit -p 4000:4000 --name game-galactic-kittens-match-develop -m 100M -e GO_OPTS=%goRunParams% game-galactic-kittens-match:develop
