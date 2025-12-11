package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Database struct {
		Uri  string
		Name string
	}
	Redis struct {
		Host     string
		Port     string
		Password string
	}
	LogLevel   string
	Enviroment string
	Server     struct {
		Port         string
		Host         string
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
	}
}

func Load() (*Config, error) {
	godotenv.Load()
	cfg := &Config{}

	cfg.Database.Uri = getEnv("CONTEST_MONGO_URI", "mongodb://mongocontest:27017/contestdb")
	cfg.Database.Name = getEnv("CONTEST_MONGO_DATABASE_NAME", "contestdb")

	cfg.Redis.Host = getEnv("REDIS_HOST", "rediscontest")
	cfg.Redis.Port = getEnv("REDIS_PORT", "6379")
	cfg.Redis.Password = getEnv("REDIS_PASSWORD", "")

	cfg.Enviroment = getEnv("CONTEST_ENV", "Development")
	cfg.LogLevel = getEnv("CONTEST_LOG_LEVEL", "debug")

	cfg.Server.Port = getEnv("CONTEST_PORT", "8001")
	cfg.Server.Host = getEnv("CONTEST_HOST", "0.0.0.0")
	cfg.Server.ReadTimeout = time.Second * 15
	cfg.Server.WriteTimeout = time.Second * 15

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
