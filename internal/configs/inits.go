package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

func InitializeViper() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()

	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading config file, ", err)
	}
}
