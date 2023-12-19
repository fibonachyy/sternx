package service

import (
	"context"
	"time"

	"github.com/fibonachyy/sternx/internal/logger"
	"github.com/fibonachyy/sternx/internal/metrics"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnaryInterceptor(myLogger logger.Logger, meter metric.Meter) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Set the logger in the context
		ctx = logger.WithLogger(ctx, myLogger)
		ctx = metrics.WithMeter(ctx, meter)

		startTime := time.Now()

		// Perform pre-handler operations or logging if needed
		myLogger.Info(ctx, "Received gRPC request")

		tracer := otel.Tracer("grpc-server")
		ctx, span := tracer.Start(ctx, info.FullMethod)
		defer span.End()

		// Add span to context
		ctx = trace.ContextWithSpan(ctx, span)

		// Add attributes specific to the service method level
		span.SetAttributes(attribute.String("interceptor.method", info.FullMethod))

		// Invoke the next middleware or handler
		resp, err := handler(ctx, req)

		duration := time.Since(startTime)

		statusCode := codes.Unknown
		if st, ok := status.FromError(err); ok {
			statusCode = st.Code()
		}

		// Record gRPC metrics
		metrics.RecordGRPCMetrics(ctx, meter, info.FullMethod, err, duration)

		myLogger.Info(ctx, "Completed gRPC request", "method", info.FullMethod, "statusCode", statusCode.String(), "duration", duration)

		return resp, err
	}
}
