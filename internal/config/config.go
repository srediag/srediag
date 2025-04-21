// internal/config/config.go
package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Debug bool   `mapstructure:"debug"`
	Port  int    `mapstructure:"port"`
	DSN   string `mapstructure:"datasource"`
}

func Load() (*Config, error) {
	viper.SetConfigFile(viper.GetString("config"))
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
