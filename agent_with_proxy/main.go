package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100
	httpClient := &http.Client{
		Timeout:   10 * time.Minute,
		Transport: t,
	}

	s := http.Server{
		Addr: ":50000",
		Handler: &Proxy{
			logger: logger,
			client: httpClient,
		},
	}

	go func() {
		logger.Info("proxy is serving")
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("proxy failed to listen and serve", zap.Error(err))
		}
	}()

	<-sig
	_ = s.Shutdown(context.Background())
}

type Proxy struct {
	logger *zap.Logger
	client *http.Client
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqURL := fmt.Sprintf("https:%s", r.URL.String())
	fields := []zap.Field{
		zap.String("request_url", reqURL),
		zap.Any("request_headers", r.Header),
	}
	u, err := url.Parse(reqURL)
	if err != nil {
		p.logger.Error("failed to parse request url", append(fields, zap.Error(err))...)
		http.Error(w, "failed to parse request url", http.StatusInternalServerError)
		return
	}
	r, err = http.NewRequest(r.Method, u.String(), r.Body)
	if err != nil {
		p.logger.Error("failed to create new request", append(fields, zap.Error(err))...)
		http.Error(w, "failed to create new request", http.StatusInternalServerError)
		return
	}
	start := time.Now()
	res, err := p.client.Do(r)
	if err != nil {
		p.logger.Error("proxy failed to do request", append(fields, zap.Error(err))...)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	elapsed := time.Since(start)
	b, err := io.ReadAll(res.Body)
	if err != nil {
		p.logger.Warn("proxy failed to read response body", zap.Error(err))
	}
	p.logger.Info(
		"receive Datadog Agent request",
		append(
			fields,
			zap.Duration("elapsed", elapsed),
			zap.Int("status_code", res.StatusCode),
			zap.String("response_body", string(b)),
		)...,
	)
	w.Write(b)
	w.WriteHeader(http.StatusOK)
}
