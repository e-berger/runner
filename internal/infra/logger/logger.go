package logger

import (
	"log/slog"
	"os"
)

func SetupLog() {
	lvl := new(slog.LevelVar)
	logLevel := os.Getenv("LOGLEVEL")
	lvl.Set(slog.LevelInfo)
	if logLevel != "" {
		slog.Info("Logger", "loglevel", logLevel)
		lvl.UnmarshalText([]byte(logLevel))
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: lvl,
	}))
	slog.SetDefault(logger)
}
