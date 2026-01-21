package interfaces

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type IBot interface {
	SendMessage(chatID int64, message string) error
	ListenUpdates() tgbotapi.UpdatesChannel
	GetBotUsername() string
}
