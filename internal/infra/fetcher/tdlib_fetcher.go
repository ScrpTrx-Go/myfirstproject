package fetcher

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/ScrpTrx-Go/myfirstproject/internal/app"
	"github.com/ScrpTrx-Go/myfirstproject/internal/model"
	"github.com/zelenin/go-tdlib/client"
)

type TDLibFetcher struct {
	client *client.Client
	me     *client.User
	log    app.Logger
}

func NewTDLibFetcher(tdlibClient *client.Client, log app.Logger) (*TDLibFetcher, error) {
	log.Info("Creating new TDLibFetcher...")
	me, err := tdlibClient.GetMe()
	if err != nil {
		log.Error("Failed to get authorized user", "error", err)
		return nil, fmt.Errorf("GetMe error: %w", err)
	}
	log.Info("Authorized succesfully",
		"user_id", me.Id,
		"first_name", me.FirstName,
	)
	log.Info("New TDLibFetcher was created")
	return &TDLibFetcher{
		client: tdlibClient,
		me:     me,
		log:    log}, nil
}

func (f *TDLibFetcher) Fetch(ctx context.Context, from, to time.Time) <-chan model.Post {
	out := make(chan model.Post)
	go func() {
		defer close(out)

		chatID, err := f.FindChat("sledcom_press")
		if err != nil {
			f.log.Error("failed to find chat", "err", err)
			return
		}

		resultCh, errCh := f.RunPipeline(ctx, chatID, from, to)

		for {
			select {
			case <-ctx.Done():
				f.log.Info("context canceled in Fetch")
				return
			case post, ok := <-resultCh:
				if !ok {
					f.log.Info("resultCh closed")
					return
				}
				out <- post
			case err, ok := <-errCh:
				if ok {
					f.log.Error("pipeline error", "err", err)
				}
			}
		}
	}()
	return out
}

func (f *TDLibFetcher) RunPipeline(ctx context.Context, chatID int64, from, to time.Time) (<-chan model.Post, <-chan error) {
	const numWorkers = 5

	rawOut := make(chan *client.Message)
	postOut := make(chan model.Post)
	errCh := make(chan error, 1)

	// Запускаем продюсера
	go func() {
		defer close(rawOut)
		f.log.Info("Producer: GetHistoryByPeriod started")
		err := f.GetHistoryByPeriod(ctx, chatID, from, to, rawOut)
		if err != nil {
			errCh <- fmt.Errorf("GetHistory failed %w", err)
		}
	}()

	// Запускаем воркеров
	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for i := 0; i < numWorkers; i++ {
		go func(workerID int) {
			defer wg.Done()
			f.log.Info("Worker started", "worker", workerID)
			for raw := range rawOut {
				post, ok := f.ValidateMessage(raw)
				if !ok {
					continue
				}
				select {
				case <-ctx.Done():
					return
				case postOut <- post:
				}
			}
		}(i)
	}

	// Закрываем канал, когда все воркеры закончат
	go func() {
		wg.Wait()
		f.log.Info("All workers finished")
		close(postOut)
	}()

	return postOut, errCh
}

func (f *TDLibFetcher) FindChat(username string) (int64, error) {
	f.log.Info("FindChat started", "username", username)
	chat, err := f.client.SearchPublicChat(&client.SearchPublicChatRequest{
		Username: username,
	})
	if err != nil {
		f.log.Error("Failed to SearchChat", "error", err)
		return 0, fmt.Errorf("SearchPublicChat error: %w", err)
	}
	if chat == nil {
		f.log.Error("ChatID is nil", "error", err)
		return 0, fmt.Errorf("chat is nil after SearchPublicChat")
	}
	f.log.Info("chat founded", "chatID", chat.Id)
	return chat.Id, nil

}

func (f *TDLibFetcher) GetHistoryByPeriod(ctx context.Context, chatID int64, from, to time.Time, out chan<- *client.Message) error {
	var fromMessageID int64
	stop := false

	for {
		select {
		case <-ctx.Done():
			f.log.Info("Producer context cancelled")
			return nil
		default:
		}

		history, err := f.client.GetChatHistory(&client.GetChatHistoryRequest{
			ChatId:        chatID,
			FromMessageId: fromMessageID,
			Offset:        0,
			Limit:         50,
			OnlyLocal:     false,
		})
		if err != nil {
			f.log.Error("GetChatHistory error", "error", err)
			return err
		}
		if len(history.Messages) == 0 || stop {
			return nil
		}

		for _, msg := range history.Messages {
			t := time.Unix(int64(msg.Date), 0)
			if t.After(to) {
				continue
			}
			if t.Before(from) {
				stop = true
				break
			}
			select {
			case <-ctx.Done():
				f.log.Info("Producer context cancelled")
				return nil
			case out <- msg:
			}
		}

		fromMessageID = history.Messages[len(history.Messages)-1].Id
	}
}

func (f *TDLibFetcher) ValidateMessage(raw *client.Message) (model.Post, bool) {
	msg := raw
	var text string

	switch content := msg.Content.(type) {
	case *client.MessageText:
		text = content.Text.Text
	case *client.MessagePhoto:
		text = content.Caption.Text
	case *client.MessageVideo:
		text = content.Caption.Text
	default:
		f.log.Warn("Unhandled message content", "type", fmt.Sprintf("%T", msg.Content))
		return model.Post{}, false
	}

	text = strings.TrimSpace(text)
	if text == "" {
		return model.Post{}, false
	}

	return model.Post{
		Text:      text,
		Timestamp: time.Unix(int64(msg.Date), 0),
	}, true
}
