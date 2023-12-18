package service

import (
	"context"
	"time"

	"github.com/fibonachyy/sternx/internal/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnaryInterceptor(myLogger logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Set the logger in the context
		ctx = logger.WithLogger(ctx, myLogger)

		startTime := time.Now()

		// Perform pre-handler operations or logging if needed
		myLogger.Info(ctx, "Received gRPC request")

		// Invoke the next middleware or handler
		resp, err := handler(ctx, req)

		duration := time.Since(startTime)

		statusCode := codes.Unknown
		if st, ok := status.FromError(err); ok {
			statusCode = st.Code()
		}

		// Log information about the completed gRPC request
		myLogger.Info(ctx, "Completed gRPC request", "method", info.FullMethod, "statusCode", statusCode.String(), "duration", duration)

		return resp, err
	}
}
