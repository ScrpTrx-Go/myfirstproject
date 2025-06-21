package app

import (
	"context"
	"time"

	"github.com/ScrpTrx-Go/myfirstproject/internal/model"
)

// PostFetcher выдает сообщения по диапазону дат
type PostFetcher interface {
	Fetch(ctx context.Context, from, to time.Time) <-chan model.Post
}

// Analyzer принимает контекст, канал с постами, выполняет логику с постом, возвращает обновленный model.Post, который передается в SaveToDB
type Analyzer interface {
	Analyze(ctx context.Context, postFromFetch <-chan model.Post) <-chan model.Post
}

// ReportGenerator создает документ с постами и статистикой
type ReportGenerator interface {
	Generate(posts []model.Post, stats Stats) error
}

/* var ch1 chan<- int  // только запись (send-only)
var ch2 <-chan int  // только чтение (receive-only) */
