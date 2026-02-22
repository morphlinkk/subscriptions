package logger

import (
	"log/slog"
	"os"

	"github.com/morphlinkk/subscriptions/internal/config"
)

func SetLogLevel(c *config.Config) {
	switch c.LogLevel {
	case "debug":
		slog.SetLogLoggerLevel(slog.LevelDebug)
	case "info":
		slog.SetLogLoggerLevel(slog.LevelInfo)
	case "warn":
		slog.SetLogLoggerLevel(slog.LevelWarn)
	case "error":
		slog.SetLogLoggerLevel(slog.LevelError)
	}
	slog.Debug("Set logger level", "level", c.LogLevel)
}

func Fatal(msg string, args ...any) {
	slog.Error(msg, args...)
	os.Exit(1)
}
