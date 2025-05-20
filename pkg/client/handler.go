package client

import (
	"context"
	"github.com/yxxchange/pipefree/pkg/pipe/model"
)

const (
	UrlList  = "/pipe/namespace/:namespace/name/:name/list"
	UrlWatch = "/pipe/namespace/:namespace/name/:name/watch"
)

type EventHandler struct {
	ctx     context.Context
	opt     *ConnectOption
	schema  model.Schema
	eventCh chan Event
}

// Setup -> List && Watch
func (e *EventHandler) Setup() error {
	return nil
}

func (e *EventHandler) List() error {
	path := e.opt.Endpoint + "/" + e.opt.Version + UrlList

	return nil
}

func (e *EventHandler) Watch() error {
	return nil
}
