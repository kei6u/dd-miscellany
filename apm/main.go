package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	tracer.Start(
		tracer.WithAgentAddr("datadog-agent:8126"),
		tracer.WithAnalyticsRate(1.0),
		tracer.WithLogStartup(true),
	)
	defer tracer.Stop()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go generateSpan(ctx, logger, "gopher1", "custom")
	go generateSpan(ctx, logger, "gopher2", "custom")
	go generateSpan(ctx, logger, "gopher3", "custom")

	<-ctx.Done()
}

func generateSpan(ctx context.Context, logger *zap.Logger, name, spanType string) {
	for {
		if ctx.Err() != nil {
			return
		}
		span := tracer.StartSpan(name, tracer.SpanType(spanType))
		time.Sleep(time.Duration(rand.Intn(10) * int(time.Second)))
		logger.Info(fmt.Sprintf("%s generate span", name), zap.Uint64("dd.span_id", span.Context().SpanID()), zap.Uint64("dd.trace_id", span.Context().TraceID()))
		span.Finish()
	}
}
