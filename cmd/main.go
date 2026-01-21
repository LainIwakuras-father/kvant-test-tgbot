package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	api_handlers "github.com/LainIwakuras-father/kvant-test-tgbot/internal/api/handlers"
	bot_handlers "github.com/LainIwakuras-father/kvant-test-tgbot/internal/bot/handlers"
	"github.com/LainIwakuras-father/kvant-test-tgbot/pkg/logger"

	"github.com/LainIwakuras-father/kvant-test-tgbot/internal/api/middleware"
	"github.com/LainIwakuras-father/kvant-test-tgbot/internal/api/service"
	"github.com/LainIwakuras-father/kvant-test-tgbot/internal/bot/adapter"

	"github.com/LainIwakuras-father/kvant-test-tgbot/internal/core/interfaces"
	"github.com/LainIwakuras-father/kvant-test-tgbot/internal/core/storage"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/LainIwakuras-father/kvant-test-tgbot/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/sirupsen/logrus"
)

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
	// Инициализация логгера
	logger.Init()

	// Создаем логгеры для модулей
	apiLogger := logger.NewModuleLogger("api")
	botLogger := logger.NewModuleLogger("bot")
	storageLogger := logger.NewModuleLogger("storage")

	// Создаем контекст
	ctx := context.Background()
	
	// Загружаем .env файл
	if err := godotenv.Load(); err != nil {
		apiLogger.Warn("Warning: .env file not found")
	}

	// Получаем токен бота
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		apiLogger.WithField("variable", "BOT_TOKEN").Fatal("Required environment variable is missing")
	}

	// Инициализация Telegram адаптера с логгером
	apiLogger.Info("Initializing Telegram adapter")
	telegramAdapter, err := adapter.NewTelegramAdapter(botToken, botLogger)
	if err != nil {
		apiLogger.WithError(err).Fatal("Failed to create Telegram adapter")
	}

	botUsername := telegramAdapter.GetBotUsername()
	apiLogger.WithFields(logrus.Fields{
		"username": botUsername,
		
	}).Info("Bot authorized successfully")

	// Инициализация хранилища с логгером
	apiLogger.Info("Initializing storage")
	db := storage.NewMemoryStorage(storageLogger)

	// Запускаем бота в горутине с контекстом
	botErrChan := make(chan error, 1)
	botCtx := logger.NewContextWithLogger(ctx, botLogger)
	go runTelegramBot(botCtx, telegramAdapter, db, botErrChan)

	// Инициализация сервиса сообщений для API
	apiLogger.Info("Initializing message service")
	messageSvc := service.NewMessageService(telegramAdapter, db, apiLogger)

	// Получаем секретный ключ
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		apiLogger.WithField("variable", "SECRET_KEY").Fatal("Required environment variable is missing")
	}

	// Инициализация middleware для проверки secret_key
	apiLogger.Info("Initializing auth middleware")
	secretKeyMiddleware := middleware.NewAuthMiddleware(secretKey, apiLogger)

	// Инициализация обработчиков API
	apiLogger.Info("Initializing API handlers")
	messageHandler := api_handlers.NewMessageHandler(messageSvc, apiLogger)

	// Настройка роутера
	router := gin.Default()
	
	// Middleware для логирования запросов
	router.Use(func(c *gin.Context) {
		start := time.Now()
		
		// Создаем request_id
		requestID := fmt.Sprintf("%d", time.Now().UnixNano())
		
		// Создаем логгер для этого запроса
		log := apiLogger.WithFields(logrus.Fields{
			"request_id": requestID,
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"ip":         c.ClientIP(),
		})
		
		// Обновляем контекст
		ctx := logger.NewContextWithLogger(c.Request.Context(), log)
		c.Request = c.Request.WithContext(ctx)
		
		// Логируем начало запроса
		log.Info("Request started")
		
		// Обрабатываем запрос
		c.Next()
		
		// Логируем завершение
		duration := time.Since(start)
		log.WithFields(logrus.Fields{
			"status":   c.Writer.Status(),
			"duration": duration.String(),
			"size":     c.Writer.Size(),
		}).Info("Request completed")
	})

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
			log := logger.FromContext(c.Request.Context())
			log.WithError(err).Error("Failed to get users count")
			c.JSON(500, gin.H{"status": "error", "error": err.Error()})
			return
		}
		
		c.JSON(200, gin.H{
			"status":      "ok",
			"bot":         "running",
			"users_count": count,
			"timestamp":   time.Now().UTC(),
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

	// Настройка и запуск HTTP сервера с graceful shutdown
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	addr := fmt.Sprintf(":%s", port)
	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Запуск сервера в отдельной горутине
	go func() {
		apiLogger.WithFields(logrus.Fields{
			"port":    port,
			"address": "http://localhost:" + port,
			"bot":     botUsername,
		}).Info("Starting API server")
		
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			apiLogger.WithError(err).Fatal("Server failed to start")
		}
	}()

	// Ожидание сигнала завершения
	<-quit
	apiLogger.Info("Shutting down server...")
	
	// Создаем контекст с таймаутом для graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	if err := server.Shutdown(shutdownCtx); err != nil {
		apiLogger.WithError(err).Error("Server forced to shutdown")
	}
	
	apiLogger.Info("Server exited properly")
}

// runTelegramBot запускает Telegram бота в отдельной горутине
func runTelegramBot(ctx context.Context, bot interfaces.IBot, db interfaces.IStorage, errChan chan<- error) {
	log := logger.FromContext(ctx)
	
	defer func() {
		if r := recover(); r != nil {
			log.WithField("recover", r).Error("Bot panicked")
			errChan <- fmt.Errorf("bot panic: %v", r)
		}
	}()

	// Создаем handler с логгером
	handlerBot := bot_handlers.NewHandler(bot, db, log)
	
	// Запускаем прослушивание обновлений
	updates := bot.ListenUpdates()
	log.Info("Bot started and waiting for messages...")

	for update := range updates {
		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID
		userID := update.Message.From.ID
		username := update.Message.From.UserName

		// Создаем контекст для этого сообщения
		msgCtx := logger.WithFields(ctx, logrus.Fields{
			"user_id":   userID,
			"chat_id":   chatID,
			"username":  username,
		})

		// Обработка команд
		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				handlerBot.HandleStart(msgCtx, username, chatID)
			default:
				// Обработка неизвестных команд
				log.WithFields(logrus.Fields{
					"command": update.Message.Command(),
					"user_id": userID,
					"chat_id": chatID,
				}).Warn("Unknown command received")
				
				if err := bot.SendMessage(msgCtx, chatID, "Неизвестная команда. Используй /start"); err != nil {
					log.WithError(err).Error("Failed to send unknown command response")
				}
			}
			continue
		}

		// Обработка обычных сообщений
		if update.Message.Text != "" {
			handlerBot.HandleTextMessage(msgCtx, userID, chatID, update.Message.Text)
		}
	}
	
	// Если цикл завершился (канал закрыт), отправляем ошибку
	errChan <- fmt.Errorf("updates channel closed")
}