package config

import (
	"fmt"
	"github.com/spf13/viper"
	"testing"
)

func TestConfig(t *testing.T) {
	app := viper.GetString("app.name")
	fmt.Println(app)
}
