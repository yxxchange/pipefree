package method

import (
	"github.com/yxxchange/pipefree/pkg/pipe/data"
)

func WrapNodeToView(node data.Node) data.NodeView {
	return data.NodeView{
		ApiVersion: node.ApiVersion,
		Kind:       node.Kind,
		MetaData:   node.MetaData,
		Spec:       node.Spec,
		Status:     node.Status,
	}
}

func WrapNodeToEvent(eventType data.EventType, node data.Node) data.Event {
	return data.Event{
		EventType: eventType,
		Kind:      node.Kind,
		Node:      node,
	}
}
