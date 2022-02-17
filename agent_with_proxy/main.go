package main

import (
	"context"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		logger.Error("listen to 8080", zap.Error(err))
	}
	s := http.Server{Handler: &Handler{Logger: logger}}
	logger.Info("proxy starts serving on 8080")
	if err := s.Serve(lis); err != nil {
		logger.Error("proxy serves", zap.Error(err))
	}

	<-ctx.Done()
	_ = lis.Close()
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
