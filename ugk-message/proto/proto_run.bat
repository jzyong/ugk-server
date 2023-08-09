:common message
tool\protoc --go_out="..\message" --go_opt=paths=source_relative --go-grpc_out="..\message" --go-grpc_opt=paths=source_relative constant.proto
tool\protoc --go_out="..\message" --go_opt=paths=source_relative --go-grpc_out="..\message" --go-grpc_opt=paths=source_relative common.proto
tool\protoc --go_out="..\message" --go_opt=paths=source_relative --go-grpc_out="..\message" --go-grpc_opt=paths=source_relative login.proto

pause
