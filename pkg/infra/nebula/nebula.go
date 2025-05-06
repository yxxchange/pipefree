package nebula

import (
	"github.com/haysons/nebulaorm"
	"github.com/spf13/viper"
	"github.com/yxxchange/pipefree/helper/log"
	"time"

	"sync"
)

var Pool *SessionManager
var once sync.Once

func Open(user, passwd, space string) error {
	if Pool == nil {
		initPool()
	}

	_, ok := Pool.Get(space)
	if !ok {
		_, err := Pool.Open(space, user, passwd)
		return err
	}
	return nil
}

func Use(space string) *nebulaorm.DB {
	session, ok := Pool.Get(space)
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
	Pool.mu.Lock()
	defer Pool.mu.Unlock()
	for _, session := range Pool.store {
		session.Close()
	}
	Pool.store = make(map[string]*Session)
}

func (s *SessionManager) Get(space string) (*Session, bool) {
	Pool.mu.Lock()
	defer Pool.mu.Unlock()
	pool, ok := Pool.store[space]
	return pool, ok
}

func (s *SessionManager) Open(space, user, passwd string) (*Session, error) {
	Pool.mu.Lock()
	defer Pool.mu.Unlock()
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
	Pool.store[space] = session
	return session, nil
}

func initPool() {
	once.Do(func() {
		Pool = &SessionManager{
			store: make(map[string]*Session),
		}
	})
}
