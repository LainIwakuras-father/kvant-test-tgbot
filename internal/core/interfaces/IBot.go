package interfaces

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type IBot interface {
	SendMessage(ctx context.Context, chatID int64, message string) error
	ListenUpdates() tgbotapi.UpdatesChannel
	GetBotUsername() string
}
