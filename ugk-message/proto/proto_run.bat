:common message
tool\protoc --go_out="..\message" --go_opt=paths=source_relative --go-grpc_out="..\message" --go-grpc_opt=paths=source_relative agent.proto
tool\protoc --go_out="..\message" --go_opt=paths=source_relative --go-grpc_out="..\message" --go-grpc_opt=paths=source_relative common.proto
tool\protoc --go_out="..\message" --go_opt=paths=source_relative --go-grpc_out="..\message" --go-grpc_opt=paths=source_relative constant.proto
tool\protoc --go_out="..\message" --go_opt=paths=source_relative --go-grpc_out="..\message" --go-grpc_opt=paths=source_relative game_galactic_kittens.proto
tool\protoc --go_out="..\message" --go_opt=paths=source_relative --go-grpc_out="..\message" --go-grpc_opt=paths=source_relative gate.proto
tool\protoc --go_out="..\message" --go_opt=paths=source_relative --go-grpc_out="..\message" --go-grpc_opt=paths=source_relative login.proto
tool\protoc --go_out="..\message" --go_opt=paths=source_relative --go-grpc_out="..\message" --go-grpc_opt=paths=source_relative player.proto
tool\protoc --go_out="..\message" --go_opt=paths=source_relative --go-grpc_out="..\message" --go-grpc_opt=paths=source_relative server.proto

pause
