package main

import (
	"github.com/yxxchange/pipefree/config"
	"github.com/yxxchange/pipefree/infra/dal"
)

func main() {
	config.Init("./config.yaml")
	dal.InitDB()
}
