package main

import (
	"context"
	"io"
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
		// To test the connection,
		// curl -k https://logger:50000/
		logger.Info("logger server starts")
		if err := s.ListenAndServeTLS("go-server.crt", "go-server.key"); err != nil && err != http.ErrServerClosed {
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
	h.Info("receive Datadog Agent request", zap.String("request_url", r.URL.String()))
	defer r.Body.Close()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		h.Error("read request body", zap.Error(err))
		return
	}
	h.Info("read request body", zap.String("body", string(b)))
	w.WriteHeader(http.StatusBadGateway)
}
