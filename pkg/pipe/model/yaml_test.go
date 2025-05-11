package model

import (
	"fmt"
	"github.com/yxxchange/pipefree/helper/serialize"
	"gopkg.in/yaml.v3"
	"os"
	"testing"
)

func Test_Deserialize(t *testing.T) {
	// path := "../examples/1.yaml"
	path := "../examples/2.yaml"
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}
	defer file.Close()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	var pipe PipeFlow
	err = serialize.YamlDeserialize(raw, &pipe)
	if err != nil {
		t.Fatalf("failed to unmarshal json: %v", err)
	}
	fmt.Printf("pipe: %+v\n", pipe)
	b, err := yaml.Marshal(pipe)
	if err != nil {
		t.Fatalf("failed to marshal json: %v", err)
	}
	fmt.Printf("yaml: %s\n", string(b))
}
