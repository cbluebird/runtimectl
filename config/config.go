package config

import (
	"github.com/spf13/viper"
	"log"
)

var Config = viper.New()

func init() {
	Config.SetConfigName("db")
	Config.SetConfigType("yaml")
	Config.AddConfigPath(".")
	Config.WatchConfig()
	if err := Config.ReadInConfig(); err != nil {
		log.Fatalf("Config not found: %v", err)
	} else {
		log.Println("Config loaded successfully")
	}
}
