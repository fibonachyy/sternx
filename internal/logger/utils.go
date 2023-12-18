package logger

import "context"

type key int

const loggerKey key = iota

// WithLogger adds a logger to the context.
func WithLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// FromContext extracts the logger from the context.
func FromContext(ctx context.Context) Logger {
	logger, _ := ctx.Value(loggerKey).(Logger)
	return logger
}
