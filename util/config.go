package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver          string        `mapstructure:"DB_DRIVER"`
	DBSource          string        `mapstructure:"DB_SOURCE"`
	ServerAddress     string        `mapstructure:"SERVER_ADDRESS"`
	SymmetricKey      string        `mapstructure:"SYMMETRIC_KEY"`
	AccessTokenExpiry time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
}

func LoadConfig(path string) (config Config, errr error) {

	viper.AddConfigPath(path)  // path to that config file
	viper.SetConfigName("app") // name of config file
	viper.SetConfigType("env") // json, xml

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
