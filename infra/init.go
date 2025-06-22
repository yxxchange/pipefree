package infra

import (
	"github.com/yxxchange/pipefree/infra/dal"
	"github.com/yxxchange/pipefree/infra/etcd"
)

func Init() {
	dal.InitDB()
	etcd.InitEtcd()
}

func Close() {
	etcd.Close()
}
