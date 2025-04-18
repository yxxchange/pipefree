package pool

import (
	"github.com/yxxchange/pipefree/pkg/view"
	"sync"
)

type Connection interface {
	ConnectionType() string
	Write(data interface{}) error
}

var pool *ConnectionPool

type ConnectionPool struct {
	mutex sync.Mutex
	pool  map[string][]Connection
}

func GetPool() *ConnectionPool {
	if pool == nil {
		pool = &ConnectionPool{
			pool: make(map[string][]Connection),
		}
	}
	return pool
}

func (c *ConnectionPool) Register(conn Connection) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	connType := conn.ConnectionType()
	c.pool[connType] = append(c.pool[connType], conn)
}

func (c *ConnectionPool) Transport(event view.Event) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	connArr := c.pool[event.Kind]
	wg := sync.WaitGroup{}
	wg.Add(len(connArr))
	for _, conn := range connArr {
		go func(conn Connection) {
			defer wg.Done()
			_ = conn.Write(event)
		}(conn)
	}
	wg.Wait()
}

func Register(connection Connection) {
	GetPool().Register(connection)
	return
}
