package method

import (
	"github.com/yxxchange/pipefree/pkg/pipe/data"
	"github.com/yxxchange/pipefree/pkg/view"
)

func WrapNodeToView(node data.Node) view.NodeView {
	return view.NodeView{
		ApiVersion: node.ApiVersion,
		Kind:       node.Kind,
		MetaData:   node.MetaData,
		Spec:       node.Spec,
		Status:     node.Status,
	}
}

func WrapNodeToEvent(eventType view.EventType, node data.Node) view.Event {
	return view.Event{
		EventType: eventType,
		Kind:      node.Kind,
		Node:      node,
	}
}
