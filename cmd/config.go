package cmd

import (
	"fmt"

	"github.com/spf13/viper"
)

var vc *viper.Viper

func loadConfig() {
	vc = viper.New()
	vc.SetConfigFile(".env")
	vc.ReadInConfig()
}

func GetConfig(key string) string {
	if vc == nil {
		loadConfig()
	}
	return fmt.Sprint(vc.Get(key))
}
