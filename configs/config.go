package configs

import "github.com/spf13/viper"

var cfg *conf

type conf struct {
	WeatherAPIKey string `mapstructure:"WEATHER_API_KEY"`
}

func LoadConfig(path string) (*conf, error) {
	if cfg != nil {
		return cfg, nil
	}

	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
