package internal

import (
	"context"
	"github.com/yxxchange/pipefree/config"
	"github.com/yxxchange/pipefree/helper/log"
	"github.com/yxxchange/pipefree/helper/serialize"
	"github.com/yxxchange/pipefree/pkg/infra"
	"github.com/yxxchange/pipefree/pkg/pipe/model"
	"github.com/yxxchange/pipefree/pkg/repo"
	"os"
	"testing"
)

var flow model.PipeFlow

func TestMain(m *testing.M) {
	config.Init("../../../config.yaml")
	infra.Init()
	path := "../../pipe/examples/2.yaml"
	b, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	err = serialize.YamlDeserialize(b, &flow)
	if err != nil {
		panic(err)
	}
	m.Run()
}

func TestCreatePipe(t *testing.T) {
	err := CreatePipe(context.TODO(), flow.ToPipeCfg())
	if err != nil {
		panic(err)
	}
	log.Info("create pipe success")
}

func TestRunPipe(t *testing.T) {
	testID := "6829bc78480667e83c1cc940"
	err := RunPipe(context.TODO(), testID)
	if err != nil {
		panic(err)
	}
	log.Info("run pipe success")
}

func TestFindNextNode(t *testing.T) {
	curID := "0b66182c-199a-4ccb-ac94-43fc65a2ab16"
	ver, err := repo.PipeRepo.FindVertex(curID, 1, model.NebulaPipeExecBasicEdge)
	if err != nil {
		panic(err)
	}
	log.Infof("%+v", ver)
}
