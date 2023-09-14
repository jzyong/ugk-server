module github.com/jzyong/ugk/login

go 1.19

require (
	github.com/jzyong/golib v0.0.27
	github.com/jzyong/ugk/common v0.0.0-00010101000000-000000000000
	github.com/jzyong/ugk/message v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.57.0
)

require (
	github.com/go-zookeeper/zk v1.0.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/klauspost/cpuid/v2 v2.2.5 // indirect
	github.com/klauspost/reedsolomon v1.11.8 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/templexxx/cpu v0.1.0 // indirect
	github.com/templexxx/xorsimd v0.4.2 // indirect
	github.com/tjfoc/gmsm v1.4.1 // indirect
	github.com/xtaci/kcp-go/v5 v5.6.2 // indirect
	golang.org/x/crypto v0.12.0 // indirect
	golang.org/x/net v0.14.0 // indirect
	golang.org/x/sys v0.11.0 // indirect
	golang.org/x/text v0.12.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230525234030-28d5490b6b19 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
)

replace github.com/jzyong/ugk/common => ../ugk-common

replace github.com/jzyong/ugk/message => ../ugk-message
