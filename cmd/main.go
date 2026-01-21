package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	
	bot_handlers "github.com/LainIwakuras-father/kvant-test-tgbot/internal/bot/handlers"
	api_handlers "github.com/LainIwakuras-father/kvant-test-tgbot/internal/api/handlers"


	"github.com/LainIwakuras-father/kvant-test-tgbot/internal/bot/adapter"
	"github.com/LainIwakuras-father/kvant-test-tgbot/internal/api/middleware"
	"github.com/LainIwakuras-father/kvant-test-tgbot/internal/api/service"
	
	"github.com/LainIwakuras-father/kvant-test-tgbot/internal/core/interfaces"
	"github.com/LainIwakuras-father/kvant-test-tgbot/internal/core/storage"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/LainIwakuras-father/kvant-test-tgbot/docs"
	swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
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
	// Загружаем .env файл
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Получаем токен бота
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatal("BOT_TOKEN environment variable is required")
	}

	// Инициализация Telegram адаптера
	telegramAdapter, err := adapter.NewTelegramAdapter(botToken)
	if err != nil {
		log.Fatalf("Failed to create Telegram adapter: %v", err)
	}

	log.Printf("Bot authorized as @%s", telegramAdapter.GetBotUsername())

	// Инициализация хранилища в памяти (ОБЩЕЕ для бота и API)
	db := storage.NewMemoryStorage()

	// ЗАПУСКАЕМ БОТА В ГОРУТИНЕ
	go runTelegramBot(telegramAdapter,db)

	// Инициализация сервиса сообщений для API
	messageSvc := service.NewMessageService(telegramAdapter, db)

	
	secret_key := os.Getenv("SECRET_KEY")
	if secret_key == "" {
		log.Fatal("SECRET_KEY environment variable is required")
	}
	// Инициализация middleware для проверки secret_key
	secretKeyMiddleware:= middleware.NewAuthMiddleware(secret_key)
	

	// Инициализация обработчиков API
	messageHandler := api_handlers.NewMessageHandler(messageSvc)

	// Настройка роутера
	router := gin.Default()
	// Swagger документация
	docs.SwaggerInfo.BasePath = "/"
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Публичные маршруты (без аутентификации)
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": "ValentinkaBot API",
			"version": "1.0.0",
			"bot":     telegramAdapter.GetBotUsername(),
			"status":  "running",
		})
	})

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	router.GET("/health", func(c *gin.Context) {
		// Проверяем состояние бота
		count, err := db.GetAllUsers(c.Request.Context())
		if err != nil {
			c.JSON(500, gin.H{"status": "error", "error": err.Error()})
			return
		}
		
		c.JSON(200, gin.H{
			"status":      "ok",
			"bot":         "running",
			"users_count": count,
		})
	})

	// Защищенные маршруты (требуют secret_key)
	api := router.Group("/api")
	api.Use(secretKeyMiddleware.Validate)
	{
		// Отправка сообщения конкретному пользователю
		api.POST("/send/:id", messageHandler.SendMessage)
		
		// Рассылка сообщения всем пользователям
		api.POST("/send/broadcast", messageHandler.BroadcastMessage)
		
	}

	// Запуск HTTP сервера
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	addr := fmt.Sprintf(":%s", port)
	log.Printf(" Starting API server on http://localhost:%s",port)
	log.Printf(" Bot running in background: @%s", telegramAdapter.GetBotUsername())
	log.Printf(" API protected with secret_key")
	
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal("Server failed:", err)
	}
}

// runTelegramBot запускает Telegram бота в отдельной горутине
func runTelegramBot(bot interfaces.IBot, db interfaces.IStorage) {

	handler_bot := bot_handlers.NewHandler(bot,db)
	// Запускаем прослушивание обновлений
	updates := bot.ListenUpdates()
	log.Println(" Бот запущен и ожидает сообщения...")

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
					log.Printf("Ошибка отправки: %v", err)
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