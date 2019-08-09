package config

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/ravenhurst/golang-playground/consts"
)

// Config type
type Config struct {
	Port int
}

// GetConfig Gets application configuration
func GetConfig() Config {
	pflag.IntP(consts.PORT, "p", 8080, "Server port")
	pflag.Parse()
	err := viper.BindPFlag(consts.PORT, pflag.CommandLine.Lookup(consts.PORT))
	if err != nil {
		fmt.Println("Couldn't handle CLI flag => ", err)
		panic(err)
	}
	var config Config
	config.Port = viper.GetInt(consts.PORT)
	return config
}