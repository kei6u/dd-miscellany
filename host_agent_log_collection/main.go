package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var duration int
	flag.IntVar(&duration, "duration-ms", 1000, "duration to sleep before append log in milliseconds")
	flag.Parse()

	sig := make(chan os.Signal, 1)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sig
		cancel()
	}()

	f, err := os.Create("test.log")
	if err != nil {
		log.Fatalln(err)
	}

	b, err := os.ReadFile("log.txt")
	if err != nil {
		log.Fatalln(err)
	}

	ticker := time.NewTicker(time.Duration(duration) * time.Millisecond)
	go func() {
		for {
			select {
			case <-ctx.Done():
				break
			case <-ticker.C:
				_, err := f.Write(b)
				if err != nil {
					log.Fatalln(err)
				}
			}
		}
	}()

	<-ctx.Done()
	_ = f.Close()
	ticker.Stop()
}
