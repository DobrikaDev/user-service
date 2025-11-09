package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/spf13/viper"
)

type Config struct {
	Port string `mapstructure:"port" env:"PORT"`

	SQL DB `mapstructure:"sql" env-prefix:"POSTGRES_"`
}

type DB struct {
	Host     string `mapstructure:"host" env:"HOST"`
	Port     int    `mapstructure:"port" env:"PORT"`
	Username string `mapstructure:"username" env:"USERNAME"`
	Password string `mapstructure:"password" env:"PASSWORD"`
	Name     string `mapstructure:"name" env:"NAME"`
}

func LoadConfigFromFile(path string) (*Config, error) {
	config := new(Config)
	viper.SetConfigFile(path)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func MustLoadConfigFromFile(path string) *Config {
	config, err := LoadConfigFromFile(path)
	if err != nil {
		panic(err)
	}

	return config
}

func LoadConfigFromEnv() (*Config, error) {
	config := new(Config)
	err := cleanenv.ReadEnv(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func MustLoadConfigFromEnv() *Config {
	config, err := LoadConfigFromEnv()
	if err != nil {
		panic(err)
	}

	return config
}
