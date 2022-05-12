package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	var message1 string
	var message2 string
	var message3 string
	flag.StringVar(&message1, "message1", "message1", "")
	flag.StringVar(&message2, "message2", "message2", "")
	flag.StringVar(&message3, "message3", "message3", "")
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	f, err := os.Create("test.log")
	if err != nil {
		log.Fatalln(err)
	}

	nLines := 0
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("wrote %d lines\n", nLines)
			return
		default:
			now := time.Now().UnixMilli()
			for _, m := range []string{message1, message2, message3} {
				_, _ = f.WriteString(fmt.Sprintf(`{"message": %q, "timestamp": "%d"}`, m, now))
				_, _ = f.WriteString("\n")
				nLines++
			}
		}
		time.Sleep(time.Second)
	}
}
