package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	HttpPort             string        `mapstructure:"HTTP_PORT"`
	UsersServiceAddress  string        `mapstructure:"USERS_SERVICE_ADDRESS"`
	ShopsServiceAddress  string        `mapstructure:"SHOPS_SERVICE_ADDRESS"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	LogLevel             string        `mapstructure:"LOG_LEVEL"`
}

func LoadConfig(path string) (config *Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
