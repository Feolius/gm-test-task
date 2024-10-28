package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	Port     string `json:"port" env:"PORT" env-default:"3000"`
	AuthHost string `json:"auth_host" env:"AUTH_HOST" env-default:"localhost"`
	AuthPort string `json:"auth_port" env:"AUTH_PORT" env-default:"8080"`
}

func GetConfig() *Config {
	cfg := &Config{}
	err := cleanenv.ReadConfig("configs/config.json", cfg)
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}
	return cfg
}
