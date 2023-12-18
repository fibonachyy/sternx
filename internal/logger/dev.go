package logger

import (
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog"
)

// DevLogger is a logger implementation using zerolog for development.
type DevLogger struct {
	logger zerolog.Logger
}

// NewDevLogger creates a new DevLogger instance using zerolog.
func NewDevLogger() Logger {
	return &DevLogger{
		logger: zerolog.New(os.Stdout).With().Timestamp().Logger(),
	}
}

// Helper function to format log entries consistently
func formatLogEntry(level string, msg string, args ...interface{}) string {
	if msg == "" {
		return fmt.Sprint(args...)
	}
	return fmt.Sprintf("[%s] %s", level, fmt.Sprintf(msg, args...))
}

// Helper function to log with level and context
func logWithLevel(logger zerolog.Logger, level string, msg string, args ...interface{}) {
	switch level {
	case "info":
		logger.Info().Msg(formatLogEntry(level, msg, args...))
	case "debug":
		logger.Debug().Msg(formatLogEntry(level, msg, args...))
	case "warn":
		logger.Warn().Msg(formatLogEntry(level, msg, args...))
	case "error":
		logger.Error().Msg(formatLogEntry(level, msg, args...))
	case "fatal":
		logger.Fatal().Msg(formatLogEntry(level, msg, args...))
	}
}

// Implementations of the Logger interface
func (l *DevLogger) Info(ctx context.Context, args ...interface{}) {
	logWithLevel(l.logger, "info", "", args...)
}

func (l *DevLogger) Infof(ctx context.Context, format string, args ...interface{}) {
	logWithLevel(l.logger, "info", format, args...)
}

func (l *DevLogger) Debug(ctx context.Context, args ...interface{}) {
	logWithLevel(l.logger, "debug", "", args...)
}

func (l *DevLogger) Debugf(ctx context.Context, format string, args ...interface{}) {
	logWithLevel(l.logger, "debug", format, args...)
}

func (l *DevLogger) Warn(ctx context.Context, args ...interface{}) {
	logWithLevel(l.logger, "warn", "", args...)
}

func (l *DevLogger) Warnf(ctx context.Context, format string, args ...interface{}) {
	logWithLevel(l.logger, "warn", format, args...)
}

func (l *DevLogger) Error(ctx context.Context, args ...interface{}) {
	logWithLevel(l.logger, "error", "", args...)
}

func (l *DevLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	logWithLevel(l.logger, "error", format, args...)
}

func (l *DevLogger) Fatal(ctx context.Context, args ...interface{}) {
	logWithLevel(l.logger, "fatal", "", args...)
}

func (l *DevLogger) Fatalf(ctx context.Context, format string, args ...interface{}) {
	logWithLevel(l.logger, "fatal", format, args...)
}
