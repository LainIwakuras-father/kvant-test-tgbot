package service

import (
	"context"
	"fmt"
	"log"

	"github.com/LainIwakuras-father/kvant-test-tgbot/internal/core/interfaces"
	"github.com/sirupsen/logrus"
)

type MessageService struct{
	bot 	interfaces.IBot
	storage interfaces.IStorage
    logger *logrus.Entry
}

func NewMessageService(bot interfaces.IBot, storage interfaces.IStorage, logger *logrus.Entry) *MessageService {
    return &MessageService{
        bot:     bot,
        storage: storage,
        logger:  logger,
    }
}
func (s *MessageService) SendMessage(ctx context.Context, chatID int64, message string) error {
    s.logger.WithFields(logrus.Fields{
		"user_id": chatID,
		"message": message,
	}).Info("Sending message to user")
    // Проверяем, существует ли пользователь
    exists, err := s.storage.UserExists(ctx, chatID)
    if err != nil {
        return fmt.Errorf("failed to check user: %w", err)
    }
    
    if !exists {
        return fmt.Errorf("user with chat_id %d not found", chatID)
    }
    
    // Отправляем сообщение через бота
    if err := s.bot.SendMessage(ctx, chatID, message); err != nil {
        return fmt.Errorf("failed to send message: %w", err)
    }
    
    log.Printf("Message sent to chat_id %d: %s", chatID, message)
    return nil
}

func (s *MessageService) BroadcastMessage(ctx context.Context, message string) (int, int, error) {
    // Получаем всех пользователей
    users, err := s.storage.GetAllUsers(ctx)
    if err != nil {
        return 0, 0, fmt.Errorf("failed to get users: %w", err)
    }
    
    total := len(users)
    sent := 0
    
    // Отправляем сообщение каждому пользователю
    for id,_ := range users {
        if err := s.bot.SendMessage(ctx,id, message); err != nil {
            log.Printf("Failed to send message to chat_id %d: %v", id, err)
            continue
        }
        sent++
    }
    
    log.Printf("Broadcast sent: %d/%d users", sent, total)
    return sent, total, nil
}
