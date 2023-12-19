package metrics

import (
	"context"

	"go.opentelemetry.io/otel/metric"
)

type key int

const meterKey key = iota

func WithMeter(ctx context.Context, meter metric.Meter) context.Context {
	return context.WithValue(ctx, meterKey, meter)
}

func FromContext(ctx context.Context) metric.Meter {
	meter, _ := ctx.Value(meterKey).(metric.Meter)
	return meter
}
