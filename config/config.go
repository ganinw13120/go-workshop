package config

import (
	"github.com/caarlos0/env/v7"
	"github.com/joho/godotenv"
)

type Config struct {
	Port             int    `env:"PORT" envDefault:"8080"`
	MongoDBURI       string `env:"MONGODB_URI"`
	JWTSecret        string `env:"JWT_SECRET"`
	JWTSigningMethod string `env:"JWT_SIGNING_METHOD"`
}

func NewConfig() Config {
	godotenv.Load()
	config := Config{}
	if err := env.Parse(&config); err != nil {
		panic(err)
	}
	return config
}
