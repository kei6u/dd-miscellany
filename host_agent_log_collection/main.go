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
			for i := 0; i < 2; i++ {
				id := xid.New()
				_, _ = f.WriteString(fmt.Sprintf(`{"message": %q, "timestamp": "%d"}`, id.String(), now))
				_, _ = f.WriteString("\n")
			}
			nLines++
		}
		time.Sleep(time.Second)
	}
}
