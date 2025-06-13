package app

import (
	"fmt"
	"strings"

	"github.com/ScrpTrx-Go/myfirstproject/internal/model"
)

type App struct { // приложение состоит из:
	fetcher  PostFetcher     // парсера постов (интерфейс)
	analyzer Analyzer        // анализатора постов (интерфейс)
	reporter ReportGenerator // генератора отчетов (интерфейс)
	log      Logger          // логгера
}

func NewApp(f PostFetcher, l Logger) *App {
	return &App{
		fetcher:  f,
		analyzer: &MockAnalyzer{},
		reporter: &MockReporter{},
		log:      l,
	}
}

/* func (a *App) ProcessPeriod(from, to time.Time) error { // метод приложения App. тк app включает в себя вышеуказ. структуры, то имеет доступ к методам этих структур
	if to.Before(from) {
		return errors.New("конечная дата раньше начальной") // валидация даты
	}

	posts, err := a.fetcher.Fetch(from, to) // запуск парсера
	if err != nil {
		return fmt.Errorf("ошибка при получении постов: %w", err)
	}

	stats := a.analyzer.Analyze(posts) // запуск анализатора

	if err := a.reporter.Generate(posts, stats); err != nil {
		return fmt.Errorf("ошибка генерации отчета: %w", err)
	}

	return nil
} */

type MockAnalyzer struct{}

func (m *MockAnalyzer) Analyze(posts []model.Post) Stats {
	wordCount := 0
	for _, p := range posts {
		wordCount += len(splitWords(p.Text))
	}
	return Stats{
		TotalPosts: len(posts),
		WordCount:  wordCount,
	}
}

type MockReporter struct{}

func (m *MockReporter) Generate(posts []model.Post, stats Stats) error {
	fmt.Println("Сгенерирован отчет:")
	fmt.Println("Постов:", stats.TotalPosts)
	fmt.Println("Слов:", stats.WordCount)
	return nil
}

func splitWords(s string) []string {
	return strings.Fields(s)
}

type Stats struct {
	TotalPosts int
	WordCount  int
}
