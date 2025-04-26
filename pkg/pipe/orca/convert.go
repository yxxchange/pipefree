package orca

import (
	"github.com/yxxchange/pipefree/pkg/pipe/model"
)

func WrapNodeToView(node model.Node) model.NodeView {
	return model.NodeView{
		ApiVersion: node.ApiVersion,
		Kind:       node.Kind,
		MetaData:   node.MetaData,
		Spec:       node.Spec,
		Status:     node.Status,
	}
}

func ConvertNodeToIdentifier(node model.Node) model.NodeIdentifier {
	return model.NodeIdentifier{
		ApiVersion: node.ApiVersion,
		Namespace:  node.Namespace,
		Kind:       node.Kind,
		Operation:  node.Operation,
	}
}
