package utils

import (
	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The value are read by viper from a config file or environment variables.
type Config struct {
	Environment            string `mapstructure:"ENVIRONMENT"`
	GRAPHQL_SERVER_ADDRESS string `mapstructure:"GRAPHQL_SERVER_ADDRESS"`
	DB_URL                 string `mapstructure:"DB_URL"`
	DB_NAME                string `mapstructure:"DB_NAME"`
}

// LoadConfig reads configuration from file or environment variable.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("json")

	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
