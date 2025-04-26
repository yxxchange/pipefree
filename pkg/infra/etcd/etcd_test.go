package etcd

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"github.com/yxxchange/pipefree/config"
	"github.com/yxxchange/pipefree/helper/log"
	"testing"
	"time"
)

func TestEtcd(t *testing.T) {
	config.InitConfig("../../../config.yaml")
	defer func() {
		_ = Close()
	}()
	err := NewLauncher().SetEndpoint(viper.GetStringSlice("etcd.endpoints")).
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
		panic(err)
	}
	v, err := Get(context.Background(), "k1")
	if err != nil {
		log.Errorf("err: %v", err)
		return
	}
	fmt.Printf("%v\n", v)
}
