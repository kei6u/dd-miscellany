package main

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	counter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "golang",
			Name:      "my_counter",
			Help:      "This is my counter",
		})

	gauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "golang",
			Name:      "my_gauge",
			Help:      "This is my gauge",
		})

	histogram = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: "golang",
			Name:      "my_histogram",
			Help:      "This is my histogram",
		})

	summary = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Namespace: "golang",
			Name:      "my_summary",
			Help:      "This is my summary",
		})
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	rand.Seed(time.Now().Unix())

	http.Handle("/metrics", promhttp.Handler())
	// RemoteDisconnected exception will occur in the Datadog Agent side.
	http.Handle("/remote_disconnected", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if wr, ok := w.(http.Hijacker); ok {
			conn, _, err := wr.Hijack()
			if err != nil {
				log.Printf("Error hijacking connection: %v", err)
				return
			}
			log.Println("Hijacked connection and close it")
			conn.Close()
		}
	}))
	prometheus.MustRegister(counter)
	prometheus.MustRegister(gauge)
	prometheus.MustRegister(histogram)
	prometheus.MustRegister(summary)
	s := &http.Server{Addr: ":8080", Handler: http.DefaultServeMux}
	go func() {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			_ = s.Shutdown(context.Background())
			log.Println("Received signal, exiting")
			return
		default:
			counter.Add(rand.Float64() * 5)
			gauge.Add(rand.Float64()*15 - 5)
			histogram.Observe(rand.Float64() * 10)
			summary.Observe(rand.Float64() * 10)
			time.Sleep(time.Second)
		}
	}
}
