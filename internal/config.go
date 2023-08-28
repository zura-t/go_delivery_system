package internal

import "github.com/spf13/viper"

type Config struct {
	HttpServerAddress   string `mapstructure:"HTTP_SERVER_ADDRESS"`
	UsersServiceAddress string `mapstructure:"USERS_SERVICE_ADDRESS"`
}

func LoadConfig(path string) (config Config, err error) {
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