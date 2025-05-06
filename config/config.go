package config

import (
	"github.com/spf13/viper"
	"sync"
)

var once sync.Once

func Init(path string) {
	once.Do(func() {
		viper.SetConfigFile(path)
		err := viper.ReadInConfig()
		if err != nil {
			panic(err)
		}
	})
}
