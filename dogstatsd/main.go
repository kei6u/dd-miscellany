package main

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DataDog/datadog-go/v5/statsd"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	viper.AutomaticEnv()

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	c := make(chan os.Signal, 1)

	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-c
		cancel()
	}()

	statsdcli, err := statsd.New(viper.GetString("AGENT_ADDR"))
	if err != nil {
		logger.Fatal("failed to create DogStatsD client", zap.Error(err))
	}
	logger.Info("create DogStatsD client")

	go func() {
		e := statsd.NewEvent("dogstatsd_client.daily_event", "Daily event by dogstatsd client")
		e.Timestamp = time.Now()
		e.Tags = []string{"source:go"}
		e.Priority = statsd.Normal
		e.AlertType = statsd.Info
		next := time.Now()
		for {
			if time.Now().After(next) {
				err := statsdcli.Event(e)
				if err != nil {
					logger.Error("failed to send an event", zap.Any("event", e), zap.Error(err))
					continue
				}
				logger.Info("send an event", zap.Any("event", e))
				now := time.Now()
				next = randClock(now.Year(), now.Month(), now.Add(24*time.Hour).Day())
				d := next.Sub(now)
				logger.Info(fmt.Sprintf("sleep until %s", next), zap.Float64("sleep_hour", d.Hours()))
				time.Sleep(d)
			}
		}
	}()

	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				v := randFloat64Between(1, 10)
				err := statsdcli.Gauge(
					"dogstatsd_client.random_float64",
					v,
					[]string{"source:go"},
					1,
				)
				if err == nil {
					logger.Info("send a metric successfully", zap.Float64("randomfloat64", v))
				} else {
					logger.Error("failed to send a metric", zap.Error(err))
				}
			}
		}
	}()

	<-ctx.Done()
	ticker.Stop()
	_ = statsdcli.Close()
}

func randFloat64Between(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func randClock(year int, month time.Month, day int) time.Time {
	h := int(math.Trunc(randFloat64Between(1, 24)))
	m := int(math.Trunc(randFloat64Between(1, 60)))
	s := int(math.Trunc(randFloat64Between(1, 60)))
	return time.Date(
		year,
		month,
		day,
		h,
		m,
		s,
		0,
		time.Local,
	)
}
