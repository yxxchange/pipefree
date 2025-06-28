package pipe_exec

import (
	"encoding/json"
	"fmt"
	"github.com/yxxchange/pipefree/infra/dal/model"
)

const KeyTemplate = "/namespace/%s/name/%s/version/%d/exec_id/%d"

func KeyGen(node *model.NodeExec) string {
	return fmt.Sprintf(KeyTemplate, node.Namespace, node.Name, node.PipeVersion, node.Id)
}

func ValueGen(node *model.NodeExec) (string, error) {
	b, err := json.Marshal(node)
	return string(b), err
}
