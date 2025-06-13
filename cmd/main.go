package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ScrpTrx-Go/myfirstproject/internal/infra/fetcher"
	"github.com/ScrpTrx-Go/myfirstproject/internal/infra/logger"
)

func main() {
	zaplogger, err := logger.NewZapLogger(false, "app.log")
	if err != nil {
		log.Fatalf("Error initialize logger: %v", err)
	}
	defer zaplogger.Sync()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tdlibclient := fetcher.NewClient()
	tdlibFetcher, err := fetcher.NewTDLibFetcher(tdlibclient, zaplogger)
	if err != nil {
		zaplogger.Error("tdlibfetcher error", err)
	}

	from := time.Date(2025, time.June, 13, 0, 0, 0, 0, time.Local)
	to := time.Date(2025, time.June, 14, 0, 0, 0, 0, time.Local)

	outputFromFetch := tdlibFetcher.Fetch(ctx, from, to)

	for {
		select {
		case msg, ok := <-outputFromFetch:
			if !ok {
				return
			}
			fmt.Println(msg.ID)
		}
	}
}
