package pipe_cfg

import (
	"context"
	"github.com/yxxchange/pipefree/config"
	"github.com/yxxchange/pipefree/helper/serialize"
	"github.com/yxxchange/pipefree/infra/dal"
	"github.com/yxxchange/pipefree/infra/dal/model"
	"io"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Initialize database connection or any other setup required before tests
	// This is a placeholder; actual implementation may vary
	config.Init("../../config.yaml")
	dal.InitDB()
	os.Exit(m.Run())
}

func TestService_Create(t *testing.T) {
	path := "../../example_001.yaml"
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("Failed to open file %s: %v", path, err)
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("Failed to read file %s: %v", path, err)
	}
	tmp := struct {
		PipeCfg     *model.PipeCfg   `yaml:"pipe"`
		NodeCfgList []*model.NodeCfg `yaml:"nodes"`
	}{}
	err = serialize.YamlDeserialize(b, &tmp)
	if err != nil {
		t.Fatalf("Failed to deserialize YAML: %v", err)
	}
	svc := NewService(context.Background())
	err = svc.Create(tmp.PipeCfg, tmp.NodeCfgList)
	if err != nil {
		t.Fatalf("Failed to create pipe configuration: %v", err)
	}
}
