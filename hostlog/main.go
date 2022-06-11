package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/rs/xid"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	testLog1, err := os.Create("test.log.1")
	if err != nil {
		log.Fatalln(err)
	}
	testLog2, err := os.Create("test.log.2")
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
			_, _ = testLog1.WriteString(fmt.Sprintf(`{"message": "%s", "timestamp": "%d"}`, xid.New().String(), time.Now().UnixMilli()))
			_, _ = testLog1.WriteString("\n")
			_, _ = testLog2.WriteString(fmt.Sprintf(`{"message": "%s", "timestamp": "%d"}`, xid.New().String(), time.Now().UnixMilli()))
			_, _ = testLog2.WriteString("\n")
			nLines++
		}
		time.Sleep(time.Second)
	}
}
