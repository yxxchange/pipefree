package connection

import "github.com/spf13/viper"

type Connection struct {
	Endpoint string `json:"endpoint"` // 连接的服务端点
	Scheme   Schema `json:"scheme"`
}

type Schema struct {
	Namespace string `json:"namespace"` // 节点命名空间
	Kind      string `json:"kind"`      // 节点资源类型
}

func NewConnectionConfig(schema Schema) Connection {
	return Connection{
		Endpoint: viper.GetString("http.endpoint"),
		Scheme:   schema,
	}
}

func (c *Connection) ListAndWatch() string {
	//
}
