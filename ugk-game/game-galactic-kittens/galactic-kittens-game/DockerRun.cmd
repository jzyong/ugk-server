docker image build -t game-galactic-kittens:develop .

docker stop game-galactic-kittens-develop
docker rm  game-galactic-kittens-develop
rem docker run -dit --name game-galactic-kittens-develop game-galactic-kittens:develop
docker run --name game-galactic-kittens-develop game-galactic-kittens:develop

pause

