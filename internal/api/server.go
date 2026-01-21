package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/LainIwakuras-father/kvant-test-tgbot/docs"
	"github.com/LainIwakuras-father/kvant-test-tgbot/internal/api/middleware"
	"github.com/LainIwakuras-father/kvant-test-tgbot/internal/api/service"
	"github.com/LainIwakuras-father/kvant-test-tgbot/internal/core/interfaces"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	api_handlers "github.com/LainIwakuras-father/kvant-test-tgbot/internal/api/handlers"
	"github.com/gin-gonic/gin"
)



func StartServer(
	bot interfaces.IBot, 
	db interfaces.IStorage, 
	secretKey string,
	) error {
// Инициализация сервиса сообщений для API
	messageSvc := service.NewMessageService(bot, db)

	
	
	// Инициализация middleware для проверки secret_key
	secretKeyMiddleware:= middleware.NewAuthMiddleware(secretKey)
	

	// Инициализация обработчиков API
	messageHandler := api_handlers.NewMessageHandler(messageSvc)

	// Настройка роутера
	router := gin.Default()
	// Swagger документация
	docs.SwaggerInfo.BasePath = "/"
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))


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
	slog.Info("Запуск API сервера", "port", 8080, "docs", "http://localhost:8080/docs/index.html")
	if err := http.ListenAndServe(addr, router); err != nil {
		slog.Error("Server failed:", err)
	}
	return nil
	}