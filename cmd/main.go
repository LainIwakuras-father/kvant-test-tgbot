package main

import (
	
	"log"
	"log/slog"
	
	"os"

	"github.com/LainIwakuras-father/kvant-test-tgbot/internal/api"
	
	"github.com/LainIwakuras-father/kvant-test-tgbot/internal/bot"

	
	"github.com/LainIwakuras-father/kvant-test-tgbot/internal/bot/adapter"

	"github.com/LainIwakuras-father/kvant-test-tgbot/internal/core/storage"
	"github.com/joho/godotenv"

	
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


	// Загружаем .env файл
	if err := godotenv.Load(); err != nil {
		slog.Warn("Не найден или не загрузился .env файл", "error", err)
	}

	// переменные окружения
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatal("BOT_TOKEN environment variable is required")
	}
	secret_key := os.Getenv("SECRET_KEY")
	if secret_key == "" {
		log.Fatal("SECRET_KEY environment variable is required")
	}
	// Инициализация Telegram адаптера
	telegramAdapter, err := adapter.NewTelegramAdapter(botToken)
	if err != nil {
		log.Fatalf("Failed to create Telegram adapter: %v", err)
	}

	slog.Info("Бот успешно авторизован", "username", telegramAdapter.GetBotUsername())

	// Инициализация хранилища в памяти (ОБЩЕЕ для бота и API)
	db := storage.NewMemoryStorage()

	// ЗАПУСКАЕМ БОТА В ГОРУТИНЕ на фоне
	go bot.Run(telegramAdapter, db)
	
	

	//Запускаем http сервер 
	if err := api.StartServer(telegramAdapter, db, secret_key); err != nil {
		slog.Error("Ошибка запуска API сервера", "err", err)
		os.Exit(1)
	}
}

