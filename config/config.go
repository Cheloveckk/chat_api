package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DbConfig DBConfig
}

type DBConfig struct {
	Key string
}

func GetConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		panic(err.Error())
	}
	return &Config{
		DBConfig{
			Key: os.Getenv("DB_CONN"),
		},
	}
}
