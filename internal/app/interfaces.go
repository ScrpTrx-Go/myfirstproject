package app

import (
	"context"
	"time"

	"github.com/ScrpTrx-Go/myfirstproject/internal/model"
)

// PostFetcher отвечает за получение постов по диапазону дат
type PostFetcher interface {
	Fetch(ctx context.Context, from, to time.Time, outputChan chan model.Post) error
}

// Analyzer обрабатывает список постов и возвращает статистику
type Analyzer interface {
	Analyze(posts []model.Post) Stats
}

// ReportGenerator создает документ с постами и статистикой
type ReportGenerator interface {
	Generate(posts []model.Post, stats Stats) error
}
