package client

import (
	"github.com/spf13/viper"
	"github.com/yxxchange/pipefree/pkg/pipe/model"
)

type EventLoop struct {
	handlerSet map[model.Schema]*EventHandler
	opt        ConnectOption
}

func (e *EventLoop) Subscribe(schema model.Schema) {
	if e.handlerSet != nil {
		e.handlerSet = make(map[model.Schema]*EventHandler)
	}
	handler := &EventHandler{}
	e.handlerSet[schema] = handler
	return
}

func NewEventLoop() *EventLoop {
	opt := ConnectOption{
		Version:         getVersion(),
		Endpoint:        viper.GetString("operator.endpoint"),
		Timeout:         viper.GetString("operator.timeout"),
		MaxIdleConnSize: viper.GetInt("operator.maxIdleConnSize"),
		IdleConnTimeout: viper.GetInt("operator.idleConnTimeout"),
	}
	ConstructTransport(opt)
	return &EventLoop{
		opt:        opt,
		handlerSet: make(map[model.Schema]*EventHandler),
	}
}

func getVersion() string {
	version := viper.GetString("operator.version")
	if version == "" {
		return "api/v1"
	}
	return version
}
