package adapter

import (
	"log"
	
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	
)

type TelegramAdapter struct {
	Bot *tgbotapi.BotAPI
}

func NewTelegramAdapter(botToken string) (*TelegramAdapter, error) {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return nil, err
	}

	bot.Debug = true
	log.Printf("Бот авторизован: %s", bot.Self.UserName)

	return &TelegramAdapter{
		Bot: bot,
		}, nil
}

func (ta *TelegramAdapter) SendMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := ta.Bot.Send(msg)
	return err
}



func (ta *TelegramAdapter) ListenUpdates() tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	return ta.Bot.GetUpdatesChan(u)
}
func (ta *TelegramAdapter) GetBotUsername() string {
	return ta.Bot.Self.UserName
}