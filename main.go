package pipefree

import (
	"flag"
	"github.com/yxxchange/pipefree/config"
	"github.com/yxxchange/pipefree/pkg/infra"
)

var configPath string

func main() {
	initFlag()
	config.Init(configPath)
	infra.Init()
	defer infra.Close()
	// TODO： Add your main logic here
}

func initFlag() {
	flag.StringVar(&configPath, "config", "./config.yaml", "config file path")
}
