package utils

import (
	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The value are read by viper from a config file or environment variables.
type Config struct {
	Environment                 string `mapstructure:"ENVIRONMENT"`
	DBSource                    string `mapstructure:"DB_SOURCE"`
	GRPCServerAddress           string `mapstructure:"GRPCServerAddress"`
	HTTPServerAddress           string `mapstructure:"HTTPServerAddress"`
	EmployerServiceAddress      string `mapstructure:"EMPLOYER_SERVICE_ADDRESS"`
	CognitoRegionCandidates     string `mapstructure:"COGNITO_REGION_CANDIDATES"`
	CognitoUserPoolIDCandidates string `mapstructure:"COGNITO_USER_POOL_ID_CANDIDATES"`
	CognitoRegionEmployers      string `mapstructure:"COGNITO_REGION_EMPLOYERS"`
	CognitoUserPoolIDEmployers  string `mapstructure:"COGNITO_USER_POOL_ID_EMPLOYERS"`
	RedisAddress                string `mapstructure:"redisAddress"`
	RedisUsername               string `mapstructure:"redisUsername"`
	RedisPassword               string `mapstructure:"redisPassword"`
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
