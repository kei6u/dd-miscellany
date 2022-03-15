package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(time.Second)
			fmt.Fprintln(os.Stderr, "hello in stderr")
			fmt.Fprintln(os.Stdout, "hello in stdout")
		}
	}
}
