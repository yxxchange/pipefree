package nebula

import (
	"fmt"
	"github.com/spf13/viper"
	nebula "github.com/vesoft-inc/nebula-go/v3"
	"github.com/yxxchange/pipefree/helper/log"

	"sync"
	"time"
)

var Pool *SessionManager
var once sync.Once

type SessionExecutor struct {
	Space string
	Err   error
	Res   *nebula.ResultSet
}

type Session struct {
	*nebula.SessionPool
	*nebula.SessionPoolConf
}

type SessionManager struct {
	mu    sync.RWMutex
	store map[string]*Session

	hostAddressList []nebula.HostAddress
	options         []nebula.SessionPoolConfOption
}

func (s *SessionManager) Close() {
	Pool.mu.Lock()
	defer Pool.mu.Unlock()
	for _, session := range Pool.store {
		session.Close()
	}
	Pool.store = make(map[string]*Session)
}

func (s *SessionManager) get(space string) (*Session, bool) {
	Pool.mu.RLock()
	defer Pool.mu.RUnlock()
	pool, ok := Pool.store[space]
	return pool, ok
}

func (s *SessionManager) create(space string) error {
	Pool.mu.Lock()
	defer Pool.mu.Unlock()
	conf, err := nebula.NewSessionPoolConf(
		viper.GetString("nebula.userName"),
		viper.GetString("nebula.password"),
		Pool.hostAddressList,
		space,
		Pool.options...,
	)
	if err != nil {
		return err
	}
	pool, err := nebula.NewSessionPool(*conf, log.AsNebularLoggerPlugin())
	if err != nil {
		return err
	}
	session := &Session{
		SessionPool:     pool,
		SessionPoolConf: conf,
	}
	Pool.store[space] = session
	return nil
}

func initPool() {
	once.Do(func() {
		Pool = &SessionManager{
			store: make(map[string]*Session),
			hostAddressList: []nebula.HostAddress{
				{
					Host: viper.GetString("nebula.host"),
					Port: viper.GetInt("nebula.port"),
				},
			},
			options: []nebula.SessionPoolConfOption{
				nebula.WithTimeOut(time.Duration(viper.GetInt("nebula.sessionConfig.timeout")) * time.Second),
				nebula.WithIdleTime(time.Duration(viper.GetInt("nebula.sessionConfig.idleTime")) * time.Second),
				nebula.WithMaxSize(viper.GetInt("nebula.sessionConfig.poolMaxSize")),
				nebula.WithMinSize(viper.GetInt("nebula.sessionConfig.poolMinSize")),
				nebula.WithHTTP2(viper.GetBool("nebula.sessionConfig.useHttp2")),
			},
		}
	})
}

func UseSpace(space string) *SessionExecutor {
	if Pool == nil {
		initPool()
	}
	executor := &SessionExecutor{
		Space: space,
	}
	if _, ok := Pool.get(space); !ok {
		executor.Err = Pool.create(space)
	}
	return executor
}

func (s *Session) Close() {
	if s.SessionPool != nil {
		s.SessionPool.Close()
	}
}

func (s *SessionExecutor) Execute(sql string) *SessionExecutor {
	if s.Err != nil {
		return s
	}
	if session, ok := Pool.get(s.Space); !ok {
		s.Err = fmt.Errorf("must specify the nebula space")
	} else {
		s.Res, s.Err = session.Execute(sql)
	}
	return s.handleResultSet(s.Res)
}

func (s *SessionExecutor) handleResultSet(res *nebula.ResultSet) *SessionExecutor {
	if res == nil {
		return s
	}
	if res.IsSucceed() {
		return s
	}
	s.Err = fmt.Errorf("nebula execute, error code: %d, error: %s", res.GetErrorCode(), res.GetErrorMsg())
	return s
}
