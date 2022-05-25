package main

import (
	"fmt"

	"github.com/lzzzzl/basketball-go/cmd"
	"github.com/lzzzzl/basketball-go/modules/log"

	"github.com/spf13/viper"
)

func init() {
	// get config by viper
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			fmt.Println("no such config file")
		} else {
			// Config file was found but another error was produced
			fmt.Println("read config error")
		}
		panic(err)
	}
}

func main() {
	log.LogInit(viper.GetString("log.level"))

	cmd.Execute()
}
