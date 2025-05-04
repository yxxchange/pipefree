package nebula

import (
	"github.com/yxxchange/pipefree/config"
	"github.com/yxxchange/pipefree/helper/log"
	"testing"
)

func TestNebula(t *testing.T) {
	config.InitConfig("../../../config.yaml")
	err := UseSpace("test").Err
	if err != nil {
		log.Errorf("err: %v", err)
		return
	}
	log.Info("nebula test ok")
}
