package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	LogLevel   string
	Enviroment string
	Server     struct {
		Port         string
		Host         string
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
	}
	Endpoints struct {
		Submission string
		Auth       string
		Problem    string
		Contest    string
	}
}

var cfg *Config = nil

func Load() *Config {
	godotenv.Load()
	if cfg != nil {
		return cfg
	}
	cfg = &Config{}

	cfg.Enviroment = getEnv("GATEWAY_ENV", "Development")

	cfg.LogLevel = getEnv("GATEWAY_LOG_LEVEL", "debug")

	cfg.Server.Port = getEnv("GATEWAY_PORT", "80")
	cfg.Server.Host = getEnv("GATEWAY_HOST", "0.0.0.0")
	cfg.Server.ReadTimeout = time.Second * 15
	cfg.Server.WriteTimeout = time.Second * 15

	cfg.Endpoints.Submission = getEnv("SUBMISSION_ENDPOINT", "")
	cfg.Endpoints.Auth = getEnv("AUTH_ENDPOINT", "")
	cfg.Endpoints.Problem = getEnv("PROBLEM_ENDPOINT", "")
	cfg.Endpoints.Contest = getEnv("CONTEST_ENDPOINT", "")
	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
