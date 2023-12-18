package logger

import (
	"context"
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func init() {
	// Use JSON formatter for structured logging
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})

	// Set log level based on environment variable or default to InfoLevel
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel != "" {
		parsedLevel, err := logrus.ParseLevel(logLevel)
		if err == nil {
			log.SetLevel(parsedLevel)
		}
	} else {
		log.SetLevel(logrus.InfoLevel)
	}

	// Log to standard output (console)
	log.SetOutput(os.Stdout)

	// Rotate logs with file-based rotation using lfshook
	log.AddHook(lfshook.NewHook(
		lfshook.WriterMap{
			logrus.InfoLevel:  rotatingLog("logs/info.log"),
			logrus.WarnLevel:  rotatingLog("logs/warn.log"),
			logrus.ErrorLevel: rotatingLog("logs/error.log"),
			logrus.FatalLevel: rotatingLog("logs/fatalog"),
		},
		&logrus.JSONFormatter{},
	))
}

// Helper function to create a rotating log file
func rotatingLog(filename string) *rotatelogs.RotateLogs {
	// Rotate log files daily
	rotationTime := 24 * time.Hour

	// Keep log files for 7 days
	maxAge := 7 * 24 * time.Hour

	// Create a new rotating log instance
	rl, err := rotatelogs.New(
		filename+".%Y%m%d",
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithRotationTime(rotationTime),
		rotatelogs.WithMaxAge(maxAge),
	)
	if err != nil {
		panic(err)
	}

	return rl
}

func NewLogrus() Logger {
	return &LogrusLogger{
		logger: log,
	}
}

type LogrusLogger struct {
	logger *logrus.Logger
}

func (l *LogrusLogger) Info(ctx context.Context, args ...interface{}) {
	log.Info(args...)
}

func (l *LogrusLogger) Infof(ctx context.Context, format string, args ...interface{}) {
	log.Infof(format, args...)
}

func (l *LogrusLogger) Debug(ctx context.Context, args ...interface{}) {
	log.Debug(args...)
}

func (l *LogrusLogger) Debugf(ctx context.Context, format string, args ...interface{}) {
	log.Debugf(format, args...)
}

func (l *LogrusLogger) Warn(ctx context.Context, args ...interface{}) {
	log.Warn(args...)
}

func (l *LogrusLogger) Warnf(ctx context.Context, format string, args ...interface{}) {
	log.Warnf(format, args...)
}

func (l *LogrusLogger) Error(ctx context.Context, args ...interface{}) {
	log.Error(args...)
}

func (l *LogrusLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	log.Errorf(format, args...)
}

func (l *LogrusLogger) Fatal(ctx context.Context, args ...interface{}) {
	log.Fatal(args...)
}

func (l *LogrusLogger) Fatalf(ctx context.Context, format string, args ...interface{}) {
	log.Fatalf(format, args...)
}
