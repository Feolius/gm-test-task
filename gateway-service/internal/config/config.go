package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Port           string `json:"port" env:"PORT" env-default:"3000"`
	AuthScheme     string `json:"auth_scheme" env:"AUTH_SCHEME" env-default:"https"`
	AuthHost       string `json:"auth_host" env:"AUTH_HOST" env-default:"localhost"`
	AuthPort       string `json:"auth_port" env:"AUTH_PORT" env-default:"8080"`
	CurrencyScheme string `json:"currency_scheme" env:"CURRENCY_SCHEME" env-default:"https"`
	CurrencyHost   string `json:"currency_host" env:"CURRENCY_HOST"`
	CurrencyPort   string `json:"currency_port" env:"CURRENCY_PORT"`
}

func GetConfig() *Config {
	cfg := &Config{}
	err := cleanenv.ReadConfig("configs/config.json", cfg)
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}
	return cfg
}
