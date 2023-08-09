tool\protoc -I ./ --csharp_out=c# StressTestingOuterRpc.proto --grpc_out=c# --plugin=protoc-gen-grpc=tool\grpc_csharp_plugin.exe

pause

