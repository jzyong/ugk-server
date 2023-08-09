module github.com/jzyong/ugk/gate

go 1.19

require (
	github.com/jzyong/golib v0.0.21
	github.com/jzyong/ugk/common v0.0.0-00010101000000-000000000000
	github.com/jzyong/ugk/message v0.0.0-00010101000000-000000000000
	github.com/xtaci/kcp-go/v5 v5.6.2
	google.golang.org/protobuf v1.28.1
)

require (
	github.com/go-zookeeper/zk v1.0.2 // indirect
	github.com/klauspost/cpuid/v2 v2.0.14 // indirect
	github.com/klauspost/reedsolomon v1.10.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/templexxx/cpu v0.0.9 // indirect
	github.com/templexxx/xorsimd v0.4.1 // indirect
	github.com/tjfoc/gmsm v1.4.1 // indirect
	golang.org/x/crypto v0.0.0-20220622213112-05595931fe9d // indirect
	golang.org/x/net v0.0.0-20220624214902-1bab6f366d9e // indirect
	golang.org/x/sys v0.0.0-20220624220833-87e55d714810 // indirect
)

replace github.com/jzyong/ugk/common => ../ugk-common

replace github.com/jzyong/ugk/message => ../ugk-message
