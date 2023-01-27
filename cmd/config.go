package cmd

import "github.com/spf13/viper"

var vc *viper.Viper

func loadConfig() {
	vc = viper.New()
	vc.SetConfigFile(".env")
	vc.ReadInConfig()
}

func GetConfig(key string) interface{} {
	if vc == nil {
		loadConfig()
	}
	return vc.Get(key)
}
