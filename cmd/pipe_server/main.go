package main

import (
	"github.com/yxxchange/pipefree/config"
	"github.com/yxxchange/pipefree/infra/dal"
	"github.com/yxxchange/pipefree/infra/etcd"
)

func main() {
	config.Init("./config.yaml")
	dal.InitDB()
	etcd.InitEtcd()
}
