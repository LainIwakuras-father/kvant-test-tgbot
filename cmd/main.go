package main

import (
	"log"
	"log/slog"

	"os"

	"github.com/LainIwakuras-father/kvant-test-tgbot/internal/api"


	"github.com/LainIwakuras-father/kvant-test-tgbot/internal/bot"

	"github.com/LainIwakuras-father/kvant-test-tgbot/internal/bot/adapter"

	"github.com/LainIwakuras-father/kvant-test-tgbot/internal/core/config"
	"github.com/LainIwakuras-father/kvant-test-tgbot/internal/core/storage"
	
)

// main.go или создай файл docs.go в корне проекта

// @title           ValentinkaBot API
// @version         1.0
// @description     API для Valentinka Telegram-бота. Позволяет отправлять сообщения пользователям и делать рассылки.
// @termsOfService  https://example.com/terms/

// @contact.name    API Support
// @contact.url     https://example.com/support
// @contact.email   support@example.com

// @license.name    MIT
// @license.url     https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /
// @schemes   http https

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-Secret-Key

func main() {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	slog.SetDefault(slog.New(handler))

	cfg, err := config.Load()
	if err != nil {
		slog.Error("Ошибка загрузки конфигурации", "err", err)
		os.Exit(1)
	}
	
	// Инициализация Telegram адаптера
	telegramAdapter, err := adapter.NewTelegramAdapter(cfg.BotToken)
	if err != nil {
		log.Fatalf("Failed to create Telegram adapter: %v", err)
	}

	slog.Info("Бот успешно авторизован", "username", telegramAdapter.GetBotUsername())

	// Инициализация хранилища в памяти (ОБЩЕЕ для бота и API)
	db := storage.NewMemoryStorage()

	// ЗАПУСКАЕМ БОТА В ГОРУТИНЕ на фоне
	go bot.Run(telegramAdapter, db)
	
	

	//Запускаем http сервер 
	if err := api.StartServer(telegramAdapter, db, cfg.SecretKey); err != nil {
		slog.Error("Ошибка запуска API сервера", "err", err)
		os.Exit(1)
	}
}

