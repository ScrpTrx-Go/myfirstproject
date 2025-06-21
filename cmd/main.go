package main

import (
	"context"
	"log"
	"time"

	"github.com/ScrpTrx-Go/myfirstproject/internal/infra/analyzer"
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
	PostAnalyzer := analyzer.NewPostAnalyzer(zaplogger)
	tdlibFetcher, err := fetcher.NewTDLibFetcher(tdlibclient, zaplogger)
	if err != nil {
		zaplogger.Error("tdlibfetcher error", err)
	}

	from := time.Date(2024, time.January, 01, 0, 0, 0, 0, time.Local)
	to := time.Date(2025, time.January, 01, 0, 0, 0, 0, time.Local)

	outputFromFetch := tdlibFetcher.Fetch(ctx, from, to)
	PostAnalyzer.Analyze(ctx, outputFromFetch)
}
