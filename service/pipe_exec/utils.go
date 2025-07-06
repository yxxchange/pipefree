package pipe_exec

import (
	"encoding/json"
	"fmt"
	"github.com/yxxchange/pipefree/infra/dal/model"
)

const KeyTemplate = "/namespace/%s/kind/%s/version/%s/node_exec/%d"
const KeyPrefixTemplate = "/namespace/%s/kind/%s/version/%s/node_exec/"

func KeyGen(node *model.NodeExec) string {
	return fmt.Sprintf(KeyTemplate, node.Namespace, node.Kind, node.Version, node.Id)
}

func ValueGen(node *model.NodeExec) (string, error) {
	b, err := json.Marshal(node)
	return string(b), err
}
