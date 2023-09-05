module github.com/jzyong/ugk/common

go 1.19

require (
	github.com/jzyong/golib v0.0.21
	github.com/xtaci/kcp-go/v5 v5.6.2
	github.com/jzyong/ugk/message v0.0.0-00010101000000-000000000000
)

require (
	github.com/go-zookeeper/zk v1.0.2 // indirect
	github.com/klauspost/cpuid/v2 v2.2.5 // indirect
	github.com/klauspost/reedsolomon v1.11.8 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/templexxx/cpu v0.1.0 // indirect
	github.com/templexxx/xorsimd v0.4.2 // indirect
	github.com/tjfoc/gmsm v1.4.1 // indirect
	golang.org/x/crypto v0.12.0 // indirect
	golang.org/x/net v0.14.0 // indirect
	golang.org/x/sys v0.11.0 // indirect
)

replace github.com/jzyong/ugk/message => ../ugk-message
