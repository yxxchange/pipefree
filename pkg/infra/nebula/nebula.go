package nebula

import (
	"fmt"
	"github.com/spf13/viper"
	nebula_go "github.com/vesoft-inc/nebula-go/v3"
	"github.com/yxxchange/pipefree/helper/log"
	"strconv"
	"strings"
	"sync"
	"time"
)

var pool *SessionManager
var once sync.Once

func Init() {
	once.Do(func() {
		initNebula()
	})
}

func initNebula() {
	user := viper.GetString("nebula.username")
	passwd := viper.GetString("nebula.password")
	spaces := viper.GetStringSlice("nebula.spaces")
	endpoints := parseAddress(viper.GetStringSlice("nebula.address"))
	initPool()
	if len(spaces) == 0 {
		panic("spaces is empty")
	}
	for _, space := range spaces {
		_, err := pool.open(space, user, passwd, endpoints)
		if err != nil {
			panic(fmt.Errorf("open space %s error: %v", space, err))
		}
	}
}

func Use(space string) *nebula_go.SessionPool {
	session, ok := pool.Get(space)
	if !ok {
		panic("must open the space before use")
	}
	return session.SessionPool
}

type Session struct {
	*nebula_go.SessionPool
}

func (s *Session) Close() {
	if s.SessionPool != nil {
		s.SessionPool.Close()
	}
}

type SessionManager struct {
	mu    sync.Mutex
	store map[string]*Session
}

func (s *SessionManager) Close() {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	for _, session := range pool.store {
		session.Close()
	}
	pool.store = make(map[string]*Session)
}

func (s *SessionManager) Get(space string) (*Session, bool) {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	pool, ok := pool.store[space]
	return pool, ok
}

func (s *SessionManager) open(space, user, passwd string, endpoints []nebula_go.HostAddress) (*Session, error) {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	opts := []nebula_go.SessionPoolConfOption{
		nebula_go.WithTimeOut(viper.GetDuration("nebula.timeout") * time.Second),
		nebula_go.WithIdleTime(viper.GetDuration("nebula.idleTime") * time.Second),
		nebula_go.WithMaxSize(viper.GetInt("nebula.maxConnSize")),
		nebula_go.WithMinSize(viper.GetInt("nebula.minConnSize")),
	}
	conf, err := nebula_go.NewSessionPoolConf(user, passwd, endpoints, space, opts...)
	sessionPool, err := nebula_go.NewSessionPool(*conf, log.AsNebularLoggerPlugin())
	if err != nil {
		return nil, err
	}
	session := &Session{
		SessionPool: sessionPool,
	}
	pool.store[space] = session
	return session, nil
}

func initPool() {
	pool = &SessionManager{
		store: make(map[string]*Session),
	}
}

func Close() {
	if pool != nil {
		pool.Close()
	}
	log.Info("close nebula session")
}

func parseAddress(endpoints []string) []nebula_go.HostAddress {
	var addresses []nebula_go.HostAddress
	for _, endpoint := range endpoints {
		s := strings.Split(endpoint, ":")
		if len(s) != 2 {
			log.Errorf("invalid nebula endpoint %s", endpoint)
			continue
		}
		host := s[0]
		port, err := strconv.Atoi(s[1])
		if err != nil {
			log.Errorf("invalid nebula endpoint %s", endpoint)
			continue
		}
		addresses = append(addresses, nebula_go.HostAddress{
			Host: host,
			Port: port,
		})
	}
	return addresses
}

type Result struct {
	Res *nebula_go.ResultSet
	Err error
}

func HandleSQL(space string, sql string, parameters map[string]interface{}) (result Result) {
	var err error
	result.Res, err = Use(space).ExecuteWithParameter(sql, parameters)
	if err != nil {
		result.Err = err
		return
	}
	if result.Res.GetErrorCode() != nebula_go.ErrorCode_SUCCEEDED {
		result.Err = fmt.Errorf("code: %v, errMsg: %s", result.Res.GetErrorMsg(), result.Res.GetErrorMsg())
		return
	}
	return
}
