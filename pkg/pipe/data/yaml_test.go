package data

import (
	"fmt"
	yaml "gopkg.in/yaml.v3"
	"os"
	"testing"
)

func Test_Deserialize(t *testing.T) {
	path := "../examples/1.yaml"
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}
	defer file.Close()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	var node Node
	err = yaml.Unmarshal(raw, &node)
	if err != nil {
		t.Fatalf("failed to unmarshal json: %v", err)
	}
	fmt.Printf("node: %+v\n", node)
	b, err := yaml.Marshal(node)
	if err != nil {
		t.Fatalf("failed to marshal json: %v", err)
	}
	fmt.Printf("yaml: %s\n", string(b))
}
