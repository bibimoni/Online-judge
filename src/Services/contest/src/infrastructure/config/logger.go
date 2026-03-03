package config

import (
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/rs/zerolog"
)

var log *zerolog.Logger = nil

func GetLogger() *zerolog.Logger {
	if log == nil {
		cfg, err := Load()
		if err != nil {
			panic("Can't load config")
		}
		log = NewLogger(cfg.LogLevel)
	}
	return log
}

func NewLogger(level string) *zerolog.Logger {
	color.NoColor = false
	logLevel := getLogLevel(level)
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		NoColor:    false,
		TimeFormat: time.RFC3339,
		FormatLevel: func(i any) string {
			lvl := strings.Trim(i.(string), "[]")
			switch lvl {
			case "INF":
				return color.New(color.FgGreen, color.Bold).Sprintf("[%s]", lvl)
			case "WRN":
				return color.New(color.FgYellow, color.Bold).Sprintf("[%s]", lvl)
			case "ERR":
				return color.New(color.FgRed, color.Bold).Sprintf("[%s]", lvl)
			default:
				return color.New(color.Bold).Sprintf("[%s]", lvl)
			}
		},
		FormatTimestamp: func(i any) string {
			return color.New(color.FgCyan).Sprint(i)
		},
		FormatCaller: func(i any) string {
			return color.New(color.FgMagenta).Sprintf("[%s]", i)
		},
		FormatMessage: func(i any) string {
			return color.New(color.FgWhite).Sprintf("{ %-20s }", i)
		},
	}
	logger := zerolog.New(consoleWriter).
		Level(logLevel).
		With().
		Caller().
		Timestamp().
		Logger()

	log = &logger
	return &logger
}

func getLogLevel(level string) zerolog.Level {
	var logLevel zerolog.Level
	switch level {
	case "info":
		logLevel = zerolog.InfoLevel
	case "debug":
		logLevel = zerolog.DebugLevel
	case "warn":
		logLevel = zerolog.WarnLevel
	case "fatal":
		logLevel = zerolog.FatalLevel
	case "panic":
		logLevel = zerolog.PanicLevel
	case "no_level":
		logLevel = zerolog.NoLevel
	case "trace":
		logLevel = zerolog.TraceLevel
	default:
		logLevel = zerolog.Disabled
	}
	return logLevel
}
