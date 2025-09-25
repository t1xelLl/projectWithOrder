package configs

import (
	"github.com/spf13/viper"
	"os"
)

type (
	Config struct {
		Postgres `yaml:"postgres"`
		Http     `yaml:"http"`
		Redis    `yaml:"redis"`
		Kafka    `yaml:"kafka"`
	}

	Http struct {
		Port string `yaml:"port"`
	}

	Postgres struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password" env:"DB_PASSWORD"`
		Database string `yaml:"database"`
		SSLMode  string `yaml:"sslmode"`
	}

	Kafka struct {
		Host    string `yaml:"host"`
		Port    string `yaml:"port"`
		Topic   string `yaml:"topic"`
		GroupID string `yaml:"group_id"`
	}

	Redis struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Password string `yaml:"password" env:"REDIS_PASSWORD"`
		DB       int    `yaml:"db"`
	}
)

func LoadConfig(path string) (Config, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return Config{}, err
	}
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		config.Postgres.Password = password
	}
	if password := os.Getenv("REDIS_PASSWORD"); password != "" {
		config.Redis.Password = password
	}

	return config, nil
}
