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

	"github.com/DataDog/datadog-api-client-go/api/v1/datadog"
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

	idxMin, idxMax := 0, 4
	dogBreeds := map[int]string{
		0: "pomeranian",
		1: "poodle",
		2: "bulldog",
		3: "doberman",
		4: "samoyed",
	}

	s := datadog.NewSeriesWithDefaults()
	s.SetHost(host)
	s.SetMetric("dog.heartrate")
	min, max := 40, 60

	ticker := time.NewTicker(300 * time.Millisecond)
	var submitCount int64
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				n := rand.Intn(idxMax-idxMin+1) + idxMin
				hr := float64(rand.Intn(max-min+1) + min)
				now := float64(time.Now().Unix())
				breed := dogBreeds[n]
				s.SetTags([]string{fmt.Sprintf("breed:%s", breed)})
				s.SetPoints([][]*float64{{&now, &hr}})
				_, _, err := cli.MetricsApi.SubmitMetrics(ctx, datadog.MetricsPayload{Series: []datadog.Series{*s}})
				if err != nil {
					log.Printf("submit metrics: %s", err)
				}
				submitCount++
			}
		}
	}()

	<-ctx.Done()
	log.Printf("submitted %d metrics", submitCount)
}
