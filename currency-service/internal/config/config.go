package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Port   string `json:"port" env:"PORT" env-default:"3030"`
	DBHost string `json:"db_host" env:"DB_HOST" env-default:"db"`
	DBPort string `json:"db_port" env:"DB_PORT" env-default:"3306"`
	DBUser string `json:"db_user" env:"DB_USER" env-default:"root"`
	DBPass string `json:"db_pass" env:"DB_PASS" env-default:"1234"`
	DBName string `json:"db_name" env:"DB_NAME" env-default:"gm_test_db"`
}

func GetConfig() *Config {
	cfg := &Config{}
	err := cleanenv.ReadConfig("configs/config.json", cfg)
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}
	return cfg
}
