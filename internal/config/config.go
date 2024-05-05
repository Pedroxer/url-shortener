package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	PgStorage  Postgres `mapstructure:"postgres"`
	ServerPort string   `mapstructure:"server_port"`
	Env        string   `mapstructure:"env"`
}

type Postgres struct {
	Addr     string `mapstructure:"address"`
	Port     int    `mapstructure:"port"`
	Database string `mapstructure:"database"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	SSLMode  string `mapstructure:"ssl_mode"`
}

func LoadConfig(path string) (*Config, error) {
	cfg := Config{}
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("json")
	v.AddConfigPath(path)

	err := v.ReadInConfig()
	if err != nil {
		return &Config{}, err
	}
	err = v.Unmarshal(&cfg)

	if err != nil {
		return &Config{}, err
	}
	return &cfg, err

}
