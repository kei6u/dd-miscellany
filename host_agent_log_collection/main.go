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

	f, err := os.OpenFile("test.log", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	msgs := strings.Split(messages, ",")

	for {
		if ctx.Err() != nil {
			break
		}
		for _, msg := range msgs {
			n, err := appendLog(f, msg)
			if silent {
				continue
			}
			if err != nil {
				log.Printf("failed to append %s to log file: %s\n", msg, err)
				continue
			}
			log.Printf("write %d bytes in log file\n", n)
		}
		time.Sleep(time.Duration(duration) * time.Millisecond)
	}
}

var msgFormat = `%s
`

func appendLog(f *os.File, msg string) (n int, err error) {
	return f.WriteString(fmt.Sprintf(msgFormat, msg))
}
