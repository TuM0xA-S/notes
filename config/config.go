package config

import (
	"log"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

//Config ...
type Config struct {
	TokenPassword string `env:"TOKEN_PASSWORD,required"`
	DBURI         string `env:"DB_URI,required"`
	Host          string `env:"HOST" envDefault:":8000"`
}

//Cfg - parsed instance of Config
var Cfg Config

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("when loading .env: %v", err)
	}
	if err := env.Parse(&Cfg); err != nil {
		log.Fatalf("when parsing .env: %v", err)
	}
}
