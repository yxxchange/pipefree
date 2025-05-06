package infra

import (
	"github.com/yxxchange/pipefree/pkg/infra/etcd"
	"github.com/yxxchange/pipefree/pkg/infra/mongoDB"
	"github.com/yxxchange/pipefree/pkg/infra/nebula"
)

func Init() {
	// Initialize MongoDB
	mongoDB.Init()

	// Initialize Nebula
	nebula.Init()

	// Initialize Etcd
	etcd.Init()
}

func Close() {
	// Close MongoDB
	mongoDB.Close()

	// Close Nebula
	nebula.Close()

	// Close Etcd
	etcd.Close()
}
