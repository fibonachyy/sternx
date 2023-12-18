package logger

import (
	"context"
	"os"

	"github.com/rs/zerolog"
	"google.golang.org/grpc/status"
)

type DatabaseLogger struct {
	logger zerolog.Logger
}

func NewDatabaseLogger() *DatabaseLogger {
	return &DatabaseLogger{
		logger: zerolog.New(os.Stdout).With().Timestamp().Str("protocol", "database").Logger(),
	}
}

func (db *DatabaseLogger) LogError(ctx context.Context, err error, format string, args ...interface{}) {
	// Check if it's a gRPC status error
	st, ok := status.FromError(err)
	if ok {
		db.logger.
			Err(err).
			Int("status_code", int(st.Code())).
			Str("status_text", st.Code().String()).
			Str("user_message", "An error occurred. Please try again later."). // Provide a user-friendly message
			Msgf(format, args...)
		return
	}

	// Log the error without gRPC status information
	db.logger.
		Err(err).
		Str("user_message", "An error occurred. Please try again later."). // Provide a user-friendly message
		Str("err_message", err.Error()).                                   // Provide a user-friendly message
		Msgf(format, args...)
}

func (db *DatabaseLogger) LogInfo(ctx context.Context, format string, args ...interface{}) {
	db.logger.Info().Msgf(format, args...)
}
