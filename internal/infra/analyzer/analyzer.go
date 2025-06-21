package analyzer

import (
	"context"
	"fmt"
	"strings"

	"github.com/ScrpTrx-Go/myfirstproject/internal/app"
	"github.com/ScrpTrx-Go/myfirstproject/internal/model"
)

type PostAnalyzer struct {
	log app.Logger
}

func NewPostAnalyzer(log app.Logger) *PostAnalyzer {
	return &PostAnalyzer{log: log}
}

func (p *PostAnalyzer) Analyze(ctx context.Context, postFromFetch <-chan model.Post) <-chan model.Post {
	counter := 0
	for msg := range postFromFetch {
		lines := strings.SplitN(msg.Text, "\n", 2)
		title := strings.TrimSpace(lines[0])
		fmt.Println("Заголовок:", title)
		counter++
		_ = msg
	}
	p.log.Info("posts", "count", counter)

	return nil
}

// Метод IsErrand - возвращает bool является ли сообщение поручением
// Метод ErrandRegion - определяет регион поручения
// Метод ErrandNature - определяет суть поручения (ВУД, доклад или иное)
// после этого итоговый model.Post отправляется в канал и сохраняется другим методом в БД со всеми нужными флагами
// Reporter сохраняет все в word и excel, дополнительно указать какие регионы выделил код и к какой категории поручений отнес (для самопроверки)
