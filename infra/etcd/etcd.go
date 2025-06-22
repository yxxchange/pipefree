package etcd

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"github.com/yxxchange/pipefree/helper/log"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

var cli *clientv3.Client

type Response struct {
	Key      []byte `json:"key"`
	Value    []byte `json:"value"`
	Revision int64  `json:"revision"` // 全局版本号
}

func InitEtcd() {
	if cli != nil {
		return // 如果已经初始化，直接返回
	}
	endpoints := viper.GetStringSlice("etcd.endpoints")
	timeout := viper.GetDuration("etcd.timeout")
	var err error
	cli, err = clientv3.New(clientv3.Config{
		Endpoints:            endpoints,
		DialTimeout:          timeout * time.Second,
		TLS:                  nil, // 如果需要 TLS 支持，可以在这里配置
		Username:             viper.GetString("etcd.auth.username"),
		Password:             viper.GetString("etcd.auth.password"),
		DialKeepAliveTimeout: time.Second * 10, // 设置心跳保持时间
		DialKeepAliveTime:    time.Second * 30, // 设置心跳间隔时间
		MaxCallSendMsgSize:   2 * 1024 * 1024,  // 设置最大发送消息大小为2MB
		MaxCallRecvMsgSize:   2 * 1024 * 1024,  // 设置最大接收消息大小为2MB
		AutoSyncInterval:     time.Second * 5,  // 自动同步间隔时间
	})
	if err != nil {
		panic("failed to connect to etcd: " + err.Error())
	}
	log.Infof("ETCD client initialized with endpoints: %v", endpoints)
	return
}

// Put 设置键值对
func Put(ctx context.Context, key, value string) error {
	_, err := cli.Put(ctx, key, value)
	return err
}

// Get 获取指定键的值
func Get(ctx context.Context, key string) (Response, error) {
	resp, err := cli.Get(ctx, key)
	if err != nil {
		return Response{}, err
	}
	if len(resp.Kvs) == 0 {
		return Response{}, nil
	}
	if len(resp.Kvs) > 1 {
		return Response{}, fmt.Errorf("unexpected multiple values for key %s", key)
	}
	return Response{
		Key:      resp.Kvs[0].Key,
		Value:    resp.Kvs[0].Value,
		Revision: resp.Header.Revision,
	}, nil
}

// Delete 删除指定键
func Delete(ctx context.Context, key string) error {
	_, err := cli.Delete(ctx, key)
	return err
}

// PutWithLease 带租约的设置键值对
func PutWithLease(ctx context.Context, key, value string, leaseID clientv3.LeaseID) error {
	_, err := cli.Put(ctx, key, value, clientv3.WithLease(leaseID))
	return err
}

// GetWithPrefix 获取指定前缀的所有键值对
func GetWithPrefix(ctx context.Context, prefix string) (map[string]Response, error) {
	resp, err := cli.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	result := make(map[string]Response, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		result[string(kv.Key)] = Response{
			Key:      kv.Key,
			Value:    kv.Value,
			Revision: resp.Header.Revision,
		}
	}
	return result, nil
}

// Watch 监听指定前缀的变化
func Watch(ctx context.Context, prefix string, callback func(key, value string, isDelete bool)) {
	rch := cli.Watch(ctx, prefix, clientv3.WithPrefix())
	for {
		select {
		case <-ctx.Done():
			return
		case wresp := <-rch:
			for _, ev := range wresp.Events {
				var isDelete bool
				val := string(ev.Kv.Value)
				if ev.Type == mvccpb.DELETE {
					isDelete = true
					val = ""
				}
				callback(string(ev.Kv.Key), val, isDelete)
			}
		}
	}
}

// Close 关闭客户端连接
func Close() {
	_ = cli.Close()
}
