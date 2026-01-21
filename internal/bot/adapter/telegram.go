package adapter

import (
	"context"
	"fmt"
	

	"github.com/LainIwakuras-father/kvant-test-tgbot/pkg/errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type TelegramAdapter struct {
	Bot *tgbotapi.BotAPI
	logger *logrus.Entry
}

func NewTelegramAdapter(botToken string, logger *logrus.Entry) (*TelegramAdapter, error) {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return nil,  errors.Wrap(err, "failed to create bot API client", "BOT_API_INIT")
	}

	bot.Debug = false
	logger.WithField("bot_username", bot.Self.UserName).Info("Telegram bot initialized")

	return &TelegramAdapter{
		Bot: bot,
		}, nil
}

func (ta *TelegramAdapter) SendMessage(ctx context.Context,chatID int64, text string) error {
	select {
	case <-ctx.Done():
		return errors.Wrap(ctx.Err(), "context canceled while sending message", "CONTEXT_CANCELED")
	default:
	}
	
	msg := tgbotapi.NewMessage(chatID, text)

	ta.logger.WithFields(logrus.Fields{
		"chat_id":      chatID,
		"text_length":  len(text),
		"text_preview": fmt.Sprintf("%.50s...", text),
	}).Debug("Sending message")

	_, err := ta.Bot.Send(msg)
	if err != nil {
		return errors.Wrap(err, "failed to send telegram message", "TELEGRAM_SEND")
	}

	ta.logger.WithField("chat_id", chatID).Info("Message sent successfully")
	return nil
}



func (ta *TelegramAdapter) ListenUpdates() tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	ta.logger.Debug("Started listening for updates")
	return ta.Bot.GetUpdatesChan(u)
	
}
func (ta *TelegramAdapter) GetBotUsername() string {
	return ta.Bot.Self.UserName
}