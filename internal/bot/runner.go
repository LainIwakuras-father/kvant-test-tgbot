package bot

import (
	"log/slog"

	"github.com/LainIwakuras-father/kvant-test-tgbot/internal/core/interfaces"
	bot_handlers"github.com/LainIwakuras-father/kvant-test-tgbot/internal/bot/handlers"
)
// runTelegramBot запускает Telegram бота в отдельной горутине
func Run(bot interfaces.IBot, db interfaces.IStorage) {

	handler_bot := bot_handlers.NewHandler(bot,db)
	// Запускаем прослушивание обновлений
	updates := bot.ListenUpdates()
	slog.Info(" Бот запущен и ожидает сообщения...")

	for update := range updates {
		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID
		userID := update.Message.From.ID
		username := update.Message.From.UserName

		// Сохраняем пользователя
	
		

		// Обработка команд
		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				handler_bot.HandleStart(username,chatID)
			default:
				// Можно добавить обработку неизвестных команд
				if err := bot.SendMessage(chatID, "Неизвестная команда. Используй /start"); err != nil {
					slog.Error("Ошибка отправки: %v", err)
				}
			}
			continue
		}

		// Обработка обычных сообщений
		if update.Message.Text != "" {
            handler_bot.HandleTextMessage(userID, chatID, update.Message.Text)
        }
	}
}