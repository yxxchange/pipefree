package client

import (
	"net"
	"net/http"
	"sync"
	"time"
)

var GlobalCfg *ConnectOption
var transport *http.Transport
var once sync.Once

type ConnectOption struct {
	Version  string `json:"version"`
	Endpoint string `json:"endpoint"`
	Timeout  string `json:"timeout"`

	MaxIdleConnSize int `json:"maxIdleConnSize"`
	IdleConnTimeout int `json:"idleConnTimeout"`
	// TODO: auth
}

func ConstructTransport(opt ConnectOption) {
	once.Do(func() {
		dialer := net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}
		transport = &http.Transport{
			Proxy:                 http.ProxyFromEnvironment,
			DialContext:           dialer.DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		}
		if opt.MaxIdleConnSize > 0 {
			transport.MaxIdleConns = opt.MaxIdleConnSize
		}
		if opt.IdleConnTimeout > 0 {
			transport.IdleConnTimeout = time.Duration(opt.IdleConnTimeout)
		}
	})
}
