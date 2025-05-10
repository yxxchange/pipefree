package internal

import (
	"github.com/yxxchange/pipefree/config"
	"github.com/yxxchange/pipefree/helper/serialize"
	"github.com/yxxchange/pipefree/pkg/infra"
	"golang.org/x/text/message/pipeline"
	"os"
	"testing"
)

var testPipe pipeline.Config

func TestMain(m *testing.M) {
	config.Init("../../../config.yaml")
	infra.Init()
	f, err := os.Open("../pipe/examples/1.yaml")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	b, err := os.ReadFile("../pipe/examples/1.yaml")
	if err != nil {
		panic(err)
	}
	err = serialize.JsonDeserialize(b, &testPipe)
	if err != nil {
		panic(err)
	}
	m.Run()
}

func TestCreatePipe(t *testing.T) {

}
