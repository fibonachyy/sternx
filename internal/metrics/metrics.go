package metrics

import (
	"context"

	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/fibonachyy/sternx/internal/logger"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	api "go.opentelemetry.io/otel/metric"
	metricsdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"google.golang.org/grpc/credentials"
)

var (
	insecure = os.Getenv("INSECURE_MODE")
)

// RunMetricGoroutine runs a metric-generating goroutine and waits for completion.
func RunMetricGoroutine(wg *sync.WaitGroup, ctx context.Context, meter api.Meter, log logger.Logger, metricFunc func(context.Context, api.Meter, logger.Logger)) {
	wg.Add(1)
	defer wg.Done()
	go func(ctx context.Context) {
		metricFunc(ctx, meter, log)
	}(ctx)
}

// CreateUserCounter generates createUser metrics.
func CreateUserCounter(ctx context.Context, meter api.Meter, log logger.Logger) {
	defer log.Infof(ctx, "createUserCounter goroutine stopped")
	counter, err := meter.Int64Counter("createUser", api.WithUnit("1"),
		api.WithDescription("Counts createUser since start"),
	)
	if err != nil {
		log.Infof(ctx, "Error in creating createUser counter: %v", err)
		return
	}

	for {
		select {
		case <-ctx.Done():
			return // Stop goroutine on context cancellation
		default:
			// Increment the counter by 1.
			counter.Add(ctx, 1)
			time.Sleep(time.Duration(rand.Int63n(5)) * time.Millisecond)
		}
	}
}

// LoginUserCounter generates login metrics.
func LoginUserCounter(ctx context.Context, meter api.Meter, log logger.Logger) {
	defer log.Info(ctx, "loginUserCounter goroutine stopped")
	counter, err := meter.Int64Counter("login", api.WithUnit("1"),
		api.WithDescription("Counts logins since start"),
	)
	if err != nil {
		log.Info(ctx, "Error in creating login counter: %v", err)
		return
	}

	for {
		select {
		case <-ctx.Done():
			return // Stop goroutine on context cancellation
		default:
			// Increment the counter by 1.
			counter.Add(ctx, 1)
			time.Sleep(time.Duration(rand.Int63n(5)) * time.Millisecond)
		}
	}
}

func InitMeter(logger logger.Logger, exporterHost string, serverName string) *metricsdk.MeterProvider {

	secureOption := otlpmetricgrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	if len(insecure) > 0 {
		secureOption = otlpmetricgrpc.WithInsecure()
	}

	exporter, err := otlpmetricgrpc.New(
		context.Background(),
		secureOption,
		otlpmetricgrpc.WithEndpoint(exporterHost),
	)

	if err != nil {
		logger.Fatalf(context.Background(), "Failed to create metrics exporter: %v", err)
	}

	// Set up resources
	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", serverName),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		logger.Fatalf(context.Background(), "Could not set resources: %v", err)
	}

	// Register the exporter with an SDK via a periodic reader.
	provider := metricsdk.NewMeterProvider(
		metricsdk.WithResource(res),
		metricsdk.WithReader(metricsdk.NewPeriodicReader(exporter)),
	)
	return provider
}

// GenerateMetrics runs goroutines for generating metrics.
func GenerateMetrics(ctx context.Context, meter api.Meter, logger logger.Logger) {
	var wg sync.WaitGroup

	// Run goroutines for generating metrics
	RunMetricGoroutine(&wg, ctx, meter, logger, CreateUserCounter)
	RunMetricGoroutine(&wg, ctx, meter, logger, LoginUserCounter)

	wg.Wait()
}
