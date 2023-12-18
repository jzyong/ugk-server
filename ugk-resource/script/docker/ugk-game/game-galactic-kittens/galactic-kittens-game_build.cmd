rem build image
cd ../../../../../ugk-game/game-galactic-kittens/galactic-kittens-game
docker image build -t game-galactic-kittens:develop .
set UnityParam="grpcUrl=192.168.110.2:4000 serverId=1"


pause

