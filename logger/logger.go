package logger

import (
	"log/slog"
	"os"
	"strings"
)

var log *slog.Logger

func init() {
	log = slog.Default()
}

// Init configures the global logger. level: debug|info|warn|error. format: text|json.
func Init(level string, format string) {
	slogLevel := parseLevel(level)
	opts := &slog.HandlerOptions{
		Level:     slogLevel,
		AddSource: slogLevel == slog.LevelDebug,
	}

	var handler slog.Handler
	switch strings.ToLower(strings.TrimSpace(format)) {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, opts)
	default:
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	log = slog.New(handler)
	slog.SetDefault(log)

	Info("logger initialized", "level", level, "format", format)
}

func parseLevel(level string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	case "info", "":
		return slog.LevelInfo
	default:
		return slog.LevelInfo
	}
}

func Debug(msg string, args ...any) {
	log.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	log.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	log.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	log.Error(msg, args...)
}

func Fatal(msg string, args ...any) {
	log.Error(msg, args...)
	os.Exit(1)
}
