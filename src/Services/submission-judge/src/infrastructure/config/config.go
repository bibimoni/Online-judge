package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Database struct {
		Uri  string
		Name string
	}
	Redis struct {
		Uri                string
		Password           string
		SubmissionQueueKey string
	}
	LogLevel        string
	SandboxLogLevel string
	Enviroment      string
	Server          struct {
		Port         string
		Host         string
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
	}
	ProblemServerAddr string
	Judge             struct {
		IDOffset int
		Amount   int
		IsMain   bool
	}
	ProblemsDir       string
	CheckerBinName    string
	InteractorBinName string
	CrossRunJarName   string
	InternalSecret string
}

func Load() (*Config, error) {
	godotenv.Load()
	cfg := &Config{}

	cfg.Database.Uri = getEnv("SUBMISSION_MONGODB_URI", "mongodb://mongosubmissiondb:37017/submissionjudgedb")
	cfg.Database.Name = getEnv("SUBMISSION_MONGODB_DATABASE_NAME", "submissionjudgedb")
	cfg.Redis.Uri = getEnv("SUBMISSION_REDIS_URI", "redissubmissionjudge:6379")
	cfg.Redis.Password = getEnv("SUBMISSION_REDIS_PASSWORD", "")
	cfg.Redis.SubmissionQueueKey = "SubmissionQueue"

	cfg.Enviroment = getEnv("SUBMISSION_ENV", "Development")

	cfg.LogLevel = getEnv("SUBMISSION_LOG_LEVEL", "debug")
	cfg.SandboxLogLevel = getEnv("SUBMISSION_SANDBOX_LOG_LEVEL", "debug")

	cfg.ProblemsDir = getEnv("SUBMISSION_JUDGE_PROBLEM_DIR", "problems_dir")

	cfg.Server.Port = getEnv("SUBMISSION_PORT", "8000")
	cfg.Server.Host = getEnv("SUBMISSION_HOST", "0.0.0.0")
	cfg.Server.ReadTimeout = time.Second * 15
	cfg.Server.WriteTimeout = time.Second * 15

	// cfg.ProblemServerAddr = "http://problem" + ":" + getEnv("PROBLEM_PORT", "3000") + "/problem/"
	cfg.ProblemServerAddr = getEnv("PROBLEM_ENDPOINT", "http://problem"+":"+getEnv("PROBLEM_PORT", "3000")) + "/problem/"

	numberOfJudges, err := strconv.Atoi(getEnv("SUBMISSION_NUMBER_OF_JUDGE", "10"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse env variable [SUBMISSION_NUMBER_OF_JUDGE] into int: %v", err)
	}
	cfg.Judge.Amount = numberOfJudges

	judgeIdOffset, err := strconv.Atoi(getEnv("SUBMISSION_JUDGE_ID_OFFSET", "0"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse env variable [SUBMISSION_JUDGE_ID_OFFSET] into int: %v", err)
	}
	cfg.Judge.IDOffset = judgeIdOffset

	cfg.CheckerBinName = "checker"
	cfg.InteractorBinName = "interactor"
	cfg.CrossRunJarName = "CrossRun.jar"
	cfg.Judge.IsMain, err = strconv.ParseBool(getEnv("SUBMISSION_IS_MAIN_JUDGE", "false"))
	if err != nil {
		GetLogger().Error().Err(err).Msgf("failed to get SUBMISSION_IS_MAIN_JUDGE field from .env file")
		return nil, err
	}

	cfg.InternalSecret = getEnv("INTERNAL_SECRET", "internal-secret")

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
