package fetcher

import (
	"log"
	"path/filepath"

	"github.com/zelenin/go-tdlib/client"
)

// константы для авторизации в Telegram
const (
	apiId   = 23451646
	apiHash = "53404d162eebc76acd958a2a55185872"
)

// Задает параметры клиента TDLib и возвращает его
func NewClient() *client.Client {
	tdlibParameters := &client.SetTdlibParametersRequest{
		UseTestDc:           false,
		DatabaseDirectory:   filepath.Join(".", "tdlib-db"),
		FilesDirectory:      filepath.Join(".", "tdlib-files"),
		UseFileDatabase:     true,
		UseChatInfoDatabase: true,
		UseMessageDatabase:  true,
		UseSecretChats:      false,
		ApiId:               apiId,
		ApiHash:             apiHash,
		SystemLanguageCode:  "en",
		DeviceModel:         "Desktop",
		SystemVersion:       "Unknown",
		ApplicationVersion:  "1.0",
	}

	authorizer := client.ClientAuthorizer(tdlibParameters)
	go client.CliInteractor(authorizer)

	_, err := client.SetLogVerbosityLevel(&client.SetLogVerbosityLevelRequest{
		NewVerbosityLevel: 1,
	})
	if err != nil {
		log.Fatalf("SetLogVerbosityLevel error: %s", err)
	}

	tdlibClient, err := client.NewClient(authorizer)
	if err != nil {
		log.Fatalf("NewClient error: %s", err)
	}

	return tdlibClient
}

// печатает используемую версию TDLib
func PrintVersionInfo() {
	versionOption, err := client.GetOption(&client.GetOptionRequest{Name: "version"})
	if err != nil {
		log.Printf("Failed to get TDLib version: %s", err)
		return
	}

	commitOption, err := client.GetOption(&client.GetOptionRequest{Name: "commit_hash"})
	if err != nil {
		log.Printf("Failed to get commit hash: %s", err)
		return
	}

	log.Printf("TDLib version: %s (commit: %s)",
		versionOption.(*client.OptionValueString).Value,
		commitOption.(*client.OptionValueString).Value)
}
