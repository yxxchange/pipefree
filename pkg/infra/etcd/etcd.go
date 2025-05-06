package etcd

import (
	"context"
	"crypto/tls"
	"errors"
	"github.com/spf13/viper"
	"github.com/yxxchange/pipefree/helper/log"
	"go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"sync"
	"time"
)

var (
	ErrClientClosed = errors.New("etcd client is closed")
)

var etcd *clientv3.Client
var once sync.Once

func Init() {
	err := NewLauncher().
		SetEndpoint(viper.GetStringSlice("etcd.endpoints")).
		RegisterLogger(log.AsZapLoggerPlugin()).
		KeepAlivePeriod(viper.GetDuration("etcd.time.keepAlivePeriod")*time.Second).
		DialTimeout(viper.GetDuration("etcd.time.dialTimeout")*time.Second).
		SetAutoSyncInterval(viper.GetDuration("etcd.time.autoSyncInterval")*time.Second).
		SetRetryCfg(
			viper.GetUint("etcd.retryCfg.maxSize"),
			viper.GetDuration("etcd.retryCfg.interval")*time.Second,
			viper.GetFloat64("etcd.retryCfg.jitter"),
		).Launch()
	if err != nil {
		log.Errorf("launch etcd client err: %v", err)
		panic(err)
	}
}

type Launcher struct {
	cfg clientv3.Config
}

func NewLauncher() *Launcher {
	return &Launcher{}
}

func (l *Launcher) Launch() (err error) {
	once.Do(func() {
		etcd, err = clientv3.New(l.cfg)
		if err != nil {
			log.Errorf("init etcd client err")
		}
	})
	return err
}

func (l *Launcher) SetEndpoint(endpoints []string) *Launcher {
	l.cfg.Endpoints = endpoints
	return l
}

func (l *Launcher) ConnectedByTLS(tlsConfig *tls.Config) *Launcher {
	l.cfg.TLS = tlsConfig
	return l
}

func (l *Launcher) ConnectedByAuth(user, passwd string) *Launcher {
	l.cfg.Username = user
	l.cfg.Password = passwd
	return l
}

func (l *Launcher) DialTimeout(dur time.Duration) *Launcher {
	l.cfg.DialTimeout = dur
	return l
}

func (l *Launcher) KeepAlivePeriod(period time.Duration) *Launcher {
	l.cfg.DialKeepAliveTime = period
	return l
}

func (l *Launcher) SetGrpcCfg(opts []grpc.DialOption) *Launcher {
	l.cfg.DialOptions = opts
	return l
}

func (l *Launcher) SetRetryCfg(max uint, interval time.Duration, jitter float64) *Launcher {
	l.cfg.MaxUnaryRetries = max
	l.cfg.BackoffWaitBetween = interval
	l.cfg.BackoffJitterFraction = jitter
	return l
}

func (l *Launcher) SetAutoSyncInterval(interval time.Duration) *Launcher {
	l.cfg.AutoSyncInterval = interval
	return l
}

func (l *Launcher) RegisterLogger(logger *zap.Logger) *Launcher {
	l.cfg.Logger = logger
	return l
}

// Put 写入键值
func Put(ctx context.Context, key, val string, opts ...clientv3.OpOption) error {
	if etcd == nil {
		return ErrClientClosed
	}

	_, err := etcd.Put(ctx, key, val, opts...)
	return err
}

// Get 读取键值
func Get(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error) {
	if etcd == nil {
		return nil, ErrClientClosed
	}
	return etcd.Get(ctx, key, opts...)
}

// Delete 删除键
func Delete(ctx context.Context, key string, opts ...clientv3.OpOption) error {
	if etcd == nil {
		return ErrClientClosed
	}

	_, err := etcd.Delete(ctx, key, opts...)
	return err
}

// 租约管理
type Lease struct {
	ID  clientv3.LeaseID
	ttl int64
	cli *clientv3.Client
}

// Grant 创建租约
func Grant(ctx context.Context, ttl int64) (*Lease, error) {
	resp, err := etcd.Grant(ctx, ttl)
	if err != nil {
		return nil, err
	}

	return &Lease{
		ID:  resp.ID,
		ttl: ttl,
		cli: etcd,
	}, nil
}

// Watch 封装
func Watch(ctx context.Context, key string, opts ...clientv3.OpOption) <-chan clientv3.WatchResponse {
	if etcd == nil {
		ch := make(chan clientv3.WatchResponse)
		close(ch)
		return ch
	}
	return etcd.Watch(ctx, key, opts...)
}

// Close 关闭连接
func Close() {
	if etcd != nil {
		err := etcd.Close()
		if err != nil {
			log.Errorf("close etcd client err: %v", err)
		}
	}
	log.Info("etcd client closed")
}
