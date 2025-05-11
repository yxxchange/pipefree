package nebula

import (
	"fmt"
	"github.com/haysons/nebulaorm"
	"github.com/spf13/viper"
	"github.com/yxxchange/pipefree/helper/log"
	"time"

	"sync"
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
	initPool()
	if len(spaces) == 0 {
		panic("spaces is empty")
	}
	for _, space := range spaces {
		_, err := pool.open(space, user, passwd)
		if err != nil {
			panic(fmt.Errorf("open space %s error: %v", space, err))
		}
	}
}

func Use(space string) *nebulaorm.DB {
	session, ok := pool.Get(space)
	if !ok {
		panic("must open the space before use")
	}
	return session.DB
}

type Session struct {
	*nebulaorm.DB
	*nebulaorm.Config
}

func (s *Session) Close() {
	if s.DB != nil {
		err := s.DB.Close()
		if err != nil {
			log.Errorf("close session error: %v", err)
		}
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

func (s *SessionManager) open(space, user, passwd string) (*Session, error) {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	conf := &nebulaorm.Config{
		Username:        user,
		Password:        passwd,
		SpaceName:       space,
		Addresses:       viper.GetStringSlice("nebula.address"),
		ConnTimeout:     viper.GetDuration("nebula.timeout") * time.Second,
		ConnMaxIdleTime: viper.GetDuration("nebula.idleTime") * time.Second,
		MaxOpenConns:    viper.GetInt("nebula.maxConnSize"),
		MinOpenConns:    viper.GetInt("nebula.minConnSize"),
	}
	db, err := nebulaorm.Open(conf)
	if err != nil {
		return nil, err
	}
	session := &Session{
		DB:     db,
		Config: conf,
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
