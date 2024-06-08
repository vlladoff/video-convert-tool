package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	Env string `env:"VCT_ENV" env-default:"dev"`
	KafkaServer
	WorkerPool
	MinioS3
}

type KafkaServer struct {
	KafkaAddress string `env:"VCT_KAFKA_SERVER_ADDRESS" env-default:"kafka:9092"`
	GroupId      string `env:"VCT_KAFKA_GROUP_ID" env-default:"vct"`
}

type WorkerPool struct {
	WorkersCount int `env:"VCT_WORKERS_COUNT" env-default:"8"`
}

type MinioS3 struct {
	Endpoint  string `env:"VCT_MINIO_ENDPOINT" env-default:"127.0.0.1:9000"`
	AccessKey string `env:"VCT_MINIO_ACCESS_KEY" env-default:"minioadmin"`
	SecretKey string `env:"VCT_MINIO_SECRET_KEY" env-default:"minioadmin"`
	Bucket    string `env:"VCT_MINIO_BUCKET" env-default:"vct"`
	UseSSL    bool   `env:"VCT_MINIO_USE_SSL" env-default:"false"`
}

func MustLoad() *Config {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
