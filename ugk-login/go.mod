module github.com/jzyong/ugk/login

go 1.19

require github.com/jzyong/golib v0.0.21

require github.com/go-zookeeper/zk v1.0.2 // indirect

replace github.com/jzyong/ugk/common => ../ugk-common

replace github.com/jzyong/ugk/message => ../ugk-message
