rem build image
cd ../../../../../ugk-game/game-galactic-kittens/galactic-kittens-game
docker image build -t game-galactic-kittens:develop .
set UnityParam="grpcUrl=192.168.110.2:4000"

rem run docker
docker stop game-galactic-kittens-develop
docker rm  game-galactic-kittens-develop
rem docker run -dit --name game-galactic-kittens-develop game-galactic-kittens:develop
docker run -dit --name game-galactic-kittens-develop -e UnityParam=%UnityParam% game-galactic-kittens:develop

pause

