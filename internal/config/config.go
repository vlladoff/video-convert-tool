package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	Env string `env:"VCT_ENV" env-default:"prod"`
	KafkaServer
	WorkerPool
}

type KafkaServer struct {
	KafkaAddress string `env:"VCT_KAFKA_SERVER_ADDRESS" env-default:"kafka:9092"`
	GroupId      string `env:"VCT_KAFKA_GROUP_ID" env-default:"vct"`
}

type WorkerPool struct {
	WorkersCount int `env:"VCT_WORKERS_COUNT" env-default:"8"`
}

func MustLoad() *Config {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
