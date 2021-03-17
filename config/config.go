package config

import (
	"log"

	"github.com/caarlos0/env"
)

//Config ...
type Config struct {
	TokenPassword string `env:"TOKEN_PASSWORD" envDefault:"only_for_testing"`
	Host          string `env:"HOST" envDefault:":8000"`
	DBHost        string `env:"DB_HOST"`
	DBName        string `env:"DB_NAME"`
	DBPassword    string `env:"DB_PASSWORD"`
}

//Cfg - parsed instance of Config
var Cfg Config

func init() {
	if err := env.Parse(&Cfg); err != nil {
		log.Fatalf("when parsing env: %v", err)
	}
}
