package metrics

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

func RecordGRPCMetrics(ctx context.Context, meter metric.Meter, method string, err error, duration time.Duration) {
	// Create a counter for counting the number of gRPC requests
	counter, _ := meter.Int64Counter(
		"grpc_requests_total",
		metric.WithUnit("1"),
		metric.WithDescription("Counts the total number of gRPC requests"),
	)

	// Create a histogram for measuring the duration of gRPC requests
	histogram, _ := meter.Int64Histogram(
		"grpc_request_duration",
		metric.WithUnit("ms"),
		metric.WithDescription("The duration of gRPC requests in milliseconds"),
	)

	// Increment the gRPC requests counter
	counter.Add(ctx, 1)

	// Record the duration in the histogram
	histogram.Record(ctx, int64(duration.Milliseconds()), metric.WithAttributes(attribute.String("method", method)))

	// If there was an error, you might want to record an error metric
	if err != nil {
		errorCounter, _ := meter.Int64Counter(
			"grpc_errors_total",
			metric.WithUnit("1"),
			metric.WithDescription("Counts the total number of gRPC errors"),
		)
		errorCounter.Add(ctx, 1)
	}
}
