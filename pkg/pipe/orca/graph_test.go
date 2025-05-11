package orca

import (
	"github.com/yxxchange/pipefree/helper/log"
	"github.com/yxxchange/pipefree/helper/serialize"
	"github.com/yxxchange/pipefree/pkg/pipe/model"
	"os"
	"testing"
)

func TestGraphParser(t *testing.T) {
	// path := "../examples/1.yaml"
	path := "../examples/2.yaml"
	// path := "../examples/3.yaml"
	// path := "../examples/4.yaml"
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	var pipe model.PipeConfig
	err = serialize.YamlDeserialize(raw, &pipe)
	if err != nil {
		panic(err)
	}

	parser := NewGraphParser()
	err = parser.Parse(pipe.PipeFlow).IsValid()
	if err != nil {
		log.Errorf("parse graph error: %v", err)
		return
	}
	origin, err := parser.FindTheOrigin()
	if err != nil {
		log.Errorf("find the origin error: %v", err)
		return
	}
	log.Infof("origin node: %s", origin.ToString())
	log.Info("parse graph success")
}
