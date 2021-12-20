package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	var duration int
	var messages string
	var silent bool
	flag.IntVar(&duration, "duration-ms", 1000, "duration to sleep before append log in milliseconds")
	flag.StringVar(&messages, "messages", "", "messages for logging, split by ,")
	flag.BoolVar(&silent, "silent", false, "suppress logging when appending message to body")
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

	msgs := strings.Split(messages, ",")

	ticker := time.NewTicker(time.Duration(duration) * time.Millisecond)
	go func() {
		for {
			select {
			case <-ctx.Done():
				break
			case <-ticker.C:
				for _, msg := range msgs {
					_, err := f.WriteString(fmt.Sprintf(msgFormat, msg))
					if silent {
						continue
					}
					if err != nil {
						log.Printf("failed to append %s to the log file: %s\n", msg, err)
						continue
					}
					log.Printf("write %s in the log file\n", msg)
				}
			}
		}
	}()

	<-ctx.Done()
	_ = f.Close()
	ticker.Stop()
}

var msgFormat = `%s
`
