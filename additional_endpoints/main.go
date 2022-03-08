package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	s := http.Server{
		Addr:    ":50000",
		Handler: &Handler{Logger: logger},
	}
	go func() {
		logger.Info("logger server starts")
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("logger serves", zap.Error(err))
		}
	}()

	<-sig
	_ = s.Shutdown(context.Background())
}

type Handler struct {
	*zap.Logger
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Info(
		"receive Datadog Agent request",
		zap.String("request_url", r.URL.String()),
		zap.Any("request_headers", r.Header),
	)
	defer r.Body.Close()
	w.WriteHeader(http.StatusOK)
}
