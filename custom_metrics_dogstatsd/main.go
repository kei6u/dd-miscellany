package main

import (
	"context"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DataDog/datadog-go/v5/statsd"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-c
		cancel()
	}()

	statsd, err := statsd.New("datadog-agent:8125")
	if err != nil {
		logger.Fatal("create DogStatsD client", zap.Error(err))
	}

	logger.Info("create DogStatsD client")
	ticker := time.NewTicker(1 * time.Second)
	rand.Seed(time.Now().UnixNano())
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				v := 1+rand.Float64()*(10-1)
				err := statsd.Gauge(
					"randomfloat64",
					v,
					[]string{"source:go"},
					1,
				)
				if err == nil {
					logger.Info("send a metric successfully", zap.Float64("randomfloat64", v))
				} else {
					logger.Error("send a metric", zap.Error(err))
				}
			}
		}
	}()

	<-ctx.Done()
	ticker.Stop()
	_ = statsd.Close()
}
