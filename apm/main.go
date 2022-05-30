package main

import (
	"context"
	"crypto/sha512"
	"errors"
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

	go generateSpanForever(ctx, logger, "span-generator", "custom")

	<-ctx.Done()
	logger.Info("span-generator stops, bye~")
}

func generateSpanForever(ctx context.Context, logger *zap.Logger, name, spanType string) {
	spanGenTicker := NewRandomTicker(3*time.Second, time.Second*10)
	errGenTicker := NewRandomTicker(10*time.Second, time.Minute)
	for {
		select {
		case <-ctx.Done():
			return
		case t := <-spanGenTicker.C:
			span := tracer.StartSpan(name, tracer.SpanType(spanType))
			h := sha512.Sum512([]byte(t.String()))
			logger.Info(fmt.Sprintf("generate a span with a sha512 hashed timestamp(%x)", h), zap.Uint64("dd.span_id", span.Context().SpanID()), zap.Uint64("dd.trace_id", span.Context().TraceID()))
			span.Finish()
		case <-errGenTicker.C:
			span := tracer.StartSpan(name, tracer.SpanType(spanType))
			err := errors.New("generate an error span")
			logger.Error(err.Error(), zap.Uint64("dd.span_id", span.Context().SpanID()), zap.Uint64("dd.trace_id", span.Context().TraceID()))
			span.Finish(tracer.WithError(err))
		}
	}
}

// RandomTicker is similar to time.Ticker but ticks at random intervals between
// the min and max duration values (stored internally as int64 nanosecond
// counts).
type RandomTicker struct {
	C     chan time.Time
	stopc chan chan struct{}
	min   int64
	max   int64
}

// NewRandomTicker returns a pointer to an initialized instance of the
// RandomTicker. Min and max are durations of the shortest and longest allowed
// ticks. Ticker will run in a goroutine until explicitly stopped.
func NewRandomTicker(min, max time.Duration) *RandomTicker {
	rt := &RandomTicker{
		C:     make(chan time.Time),
		stopc: make(chan chan struct{}),
		min:   min.Nanoseconds(),
		max:   max.Nanoseconds(),
	}
	go rt.loop()
	return rt
}

// Stop terminates the ticker goroutine and closes the C channel.
func (rt *RandomTicker) Stop() {
	c := make(chan struct{})
	rt.stopc <- c
	<-c
}

func (rt *RandomTicker) loop() {
	defer close(rt.C)
	t := time.NewTimer(rt.nextInterval())
	for {
		// either a stop signal or a timeout
		select {
		case c := <-rt.stopc:
			t.Stop()
			close(c)
			return
		case <-t.C:
			select {
			case rt.C <- time.Now():
				t.Stop()
				t = time.NewTimer(rt.nextInterval())
			default:
				// there could be noone receiving...
			}
		}
	}
}

func (rt *RandomTicker) nextInterval() time.Duration {
	interval := rand.Int63n(rt.max-rt.min) + rt.min
	return time.Duration(interval) * time.Nanosecond
}
