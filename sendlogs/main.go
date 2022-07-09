package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DataDog/datadog-api-client-go/api/v2/datadog"
)

func main() {
	var apiKey, appKey string
	if apiKey = os.Getenv("DD_API_KEY"); apiKey == "" {
		log.Fatalln("DD_API_KEY environment variable is not set")
	}
	if appKey = os.Getenv("DD_APP_KEY"); appKey == "" {
		log.Fatalln("DD_APP_KEY environment variable is not set")
	}
	host, err := os.Hostname()
	if err != nil {
		log.Fatalf("get hostname: %s", err)
	}
	ctx := context.WithValue(
		context.Background(),
		datadog.ContextAPIKeys,
		map[string]datadog.APIKey{
			"apiKeyAuth": {Key: apiKey},
			"appKeyAuth": {Key: appKey},
		},
	)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)

	go func() {
		<-ch
		cancel()
	}()

	cli := datadog.NewAPIClient(datadog.NewConfiguration())
	rand.Seed(time.Now().UnixNano())

	min, max := 0, 4
	dogBreeds := map[int]string{
		0: "pomeranian",
		1: "poodle",
		2: "bulldog",
		3: "doberman",
		4: "samoyed",
	}

	ticker := time.NewTicker(500 * time.Millisecond)
	var submitCount int64
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				n := rand.Intn(max-min+1) + min
				breed := dogBreeds[n]
				_, _, err := cli.LogsApi.SubmitLog(ctx, []datadog.HTTPLogItem{
					{
						Ddsource: datadog.PtrString("go"),
						Ddtags:   datadog.PtrString(fmt.Sprintf("env:local,level:info,breed:%s", breed)),
						Hostname: datadog.PtrString(host),
						Message:  fmt.Sprintf("%s bow wow", breed),
						Service:  datadog.PtrString("dog_dialog"),
					},
				})
				if err != nil {
					log.Printf("submit log: %s", err)
				}
				submitCount++
			}
		}
	}()

	<-ctx.Done()
	log.Printf("submitted %d logs", submitCount)
}
