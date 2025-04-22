package etcd

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/yxxchange/pipefree/helper/log"
	"go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"time"
)

var (
	ErrClientClosed = errors.New("etcd client is closed")
)

// Etcd 封装客户端
type Etcd struct {
	client     *clientv3.Client
	config     clientv3.Config
	ctx        context.Context
	cancelFunc context.CancelFunc
	logger     *zap.Logger
}

// NewClient 创建客户端实例
func NewClient(endpoints []string, options ...Option) (*Etcd, error) {
	cfg := clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
		Logger:      log.AsZapLoggerPlugin(),
	}

	// 应用可选参数
	for _, opt := range options {
		opt(&cfg)
	}

	// 创建官方客户端
	cli, err := clientv3.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create etcd client: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	ec := &Etcd{
		client:     cli,
		config:     cfg,
		ctx:        ctx,
		logger:     log.AsZapLoggerPlugin(),
		cancelFunc: cancel,
	}

	// 启动连接健康检查
	go ec.healthCheck()

	return ec, nil
}

// Option 客户端配置选项
type Option func(*clientv3.Config)

func WithTLS(tlsConfig *tls.Config) Option {
	return func(c *clientv3.Config) {
		c.TLS = tlsConfig
	}
}

func WithAuth(username, password string) Option {
	return func(c *clientv3.Config) {
		c.Username = username
		c.Password = password
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(c *clientv3.Config) {
		c.DialTimeout = timeout
	}
}

func WithLogger(logger *zap.Logger) Option {
	return func(c *clientv3.Config) {
		// 需要自定义日志适配器
	}
}

// 基础操作封装
// Put 写入键值
func (ec *Etcd) Put(ctx context.Context, key, val string, opts ...clientv3.OpOption) error {
	if ec.client == nil {
		return ErrClientClosed
	}

	_, err := ec.client.Put(ctx, key, val, opts...)
	return err
}

// Get 读取键值
func (ec *Etcd) Get(ctx context.Context, key string, opts ...clientv3.OpOption) (map[string]string, error) {
	if ec.client == nil {
		return nil, ErrClientClosed
	}

	resp, err := ec.client.Get(ctx, key, opts...)
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, kv := range resp.Kvs {
		result[string(kv.Key)] = string(kv.Value)
	}
	return result, nil
}

// Delete 删除键
func (ec *Etcd) Delete(ctx context.Context, key string, opts ...clientv3.OpOption) error {
	if ec.client == nil {
		return ErrClientClosed
	}

	_, err := ec.client.Delete(ctx, key, opts...)
	return err
}

// 租约管理
type Lease struct {
	ID  clientv3.LeaseID
	ttl int64
	cli *clientv3.Client
}

// Grant 创建租约
func (ec *Etcd) Grant(ctx context.Context, ttl int64) (*Lease, error) {
	resp, err := ec.client.Grant(ctx, ttl)
	if err != nil {
		return nil, err
	}

	return &Lease{
		ID:  resp.ID,
		ttl: ttl,
		cli: ec.client,
	}, nil
}

// KeepAlive 自动续约
func (l *Lease) KeepAlive(ctx context.Context) (<-chan *clientv3.LeaseKeepAliveResponse, error) {
	return l.cli.KeepAlive(ctx, l.ID)
}

// Watch 封装
func (ec *Etcd) Watch(ctx context.Context, key string, opts ...clientv3.OpOption) <-chan clientv3.WatchResponse {
	if ec.client == nil {
		ch := make(chan clientv3.WatchResponse)
		close(ch)
		return ch
	}
	return ec.client.Watch(ctx, key, opts...)
}

// Close 关闭连接
func (ec *Etcd) Close() error {
	if ec.client != nil {
		ec.cancelFunc()
		return ec.client.Close()
	}
	return nil
}

// 健康检查
func (ec *Etcd) healthCheck() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ec.ctx.Done():
			return
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(ec.ctx, 3*time.Second)
			defer cancel()

			if _, err := ec.client.Status(ctx, ec.config.Endpoints[0]); err != nil {
				ec.logger.Error("etcd health check failed", zap.Error(err))
			}
		}
	}
}
