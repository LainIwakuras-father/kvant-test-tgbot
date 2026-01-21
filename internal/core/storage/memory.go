package storage

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"sync"

	"github.com/LainIwakuras-father/kvant-test-tgbot/internal/core/models"
	"github.com/sirupsen/logrus"
)

// MemoryStorage - временное хранилище в памяти
type MemoryStorage struct {
	messages map[string]*models.Message
	users map[int64]string
	mu sync.RWMutex
	logger *logrus.Entry
}

func NewMemoryStorage(logger *logrus.Entry) *MemoryStorage {
	return &MemoryStorage{
		messages: make(map[string]*models.Message),
		users: make(map[int64]string),
		logger: logger,
	}
}

func (m *MemoryStorage) Save(ctx context.Context, userID int64, chatID int64, text string) error {

	select {
	case <-ctx.Done():
		m.logger.WithField("user_id", userID).Warn("Context canceled while saving user")
		return ctx.Err()
	default:
	}
	// Генерируем уникальный ID
	m.mu.Lock()
	defer m.mu.Unlock()
	generatedID := generateID()

	msg := &models.Message{
		ChatID: chatID,
		UserID: userID,
		Text:   text,
	}

	m.messages[generatedID] = msg
	m.logger.WithFields(logrus.Fields{
		"user_id":     userID,
		"total_users": len(m.users),
	}).Debug("User saved")
	return  nil

}

func (m *MemoryStorage) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.messages)
}



func (m *MemoryStorage) AddUser(chatID int64,username string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.users[chatID] = username
	return nil
}

func (m *MemoryStorage) UserExists(ctx context.Context, chatID int64) (bool,error){
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, exists := m.users[chatID]
	return exists, nil
}

func (m *MemoryStorage) GetAllUsers(ctx context.Context) (map[int64]string,error) {
	select {
	case <-ctx.Done():
		m.logger.Warn("Context canceled while getting users")
		return nil, ctx.Err()
	default:
	}

	m.mu.RLock()
	defer m.mu.RUnlock()
	
	// Возвращаем копию
	result := make(map[int64]string, len(m.users))
	for k, v := range m.users {
		result[k] = v
	}
	return result, nil
}

// generateID - генерирует уникальный ID
func generateID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}