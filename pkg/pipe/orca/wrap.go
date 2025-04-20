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

func WrapNodeToEvent(eventType model.EventType, node model.Node) model.Event {
	return model.Event{
		EventType: eventType,
		Node:      node,
	}
}
