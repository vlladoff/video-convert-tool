package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"time"
)

type Config struct {
	Env string `env:"VCT_ENV" env-default:"prod"`
	HTTPServer
}

type HTTPServer struct {
	Address     string        `env:"VCT_HTTP_SERVER_ADDRESS" env-default:"localhost:8080"`
	Timeout     time.Duration `env:"VCT_HTTP_SERVER_TIMEOUT" env-default:"30s"`
	IdleTimeout time.Duration `env:"VCT_HTTP_SERVER_IDLE_TIMEOUT" env-default:"60s"`
}

func MustLoad() *Config {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
