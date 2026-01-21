package handlers_bot

import (
	"context"
	
	"strconv"

	"github.com/LainIwakuras-father/kvant-test-tgbot/internal/core/interfaces"
	"github.com/sirupsen/logrus"
)

// поменять на интерфейсы
type Handler struct {
	adapter interfaces.IBot
	db      interfaces.IStorage
	logger *logrus.Entry
}

func NewHandler(adapter interfaces.IBot, db interfaces.IStorage, logger *logrus.Entry) *Handler {
	return &Handler{
		adapter: adapter,
		db:      db,
		logger:  logger,
	}
}

// HandleStart обрабатывает /start
func (h *Handler) HandleStart(ctx context.Context, username string, chatID int64) {
	
	log := h.logger.WithFields(logrus.Fields{
		"handler":  "HandleStart",
		"username": username,
		"chat_id":  chatID,
	})
	if err:=h.db.AddUser(chatID,username); err !=nil{
		log.Printf("Ошибка добавления пользователя в базу")
	}

	log.Printf("👤 User added: %d (@%s)", chatID, username)


	text := `Приветствую Друг!
Я бот сохраняющий и отправляющий твои сообщения кому либо!
Напиши Юзернейм или ID пользователя и текст который хочешь ему отправить
Формат:
9038487587 Привет, Друг!
	`
	if err := h.adapter.SendMessage(ctx, chatID, text); err != nil {
		log.WithError(err).Error("Failed to send welcome message")
		return
	}
	log.Info("Welcome message sent successfully")
}

// HandleTextMessage обрабатывает текстовые сообщения любые кроме команд
func (h *Handler) HandleTextMessage(ctx context.Context, userID int64, chatID int64, text string) {
	log := h.logger.WithFields(logrus.Fields{
		"handler": "HandleTextMessage",
		"user_id": userID,
		"chat_id": chatID,
		"text":    text,
	})
	// Сохранить сообщение в памяти
	h.SaveMessage(ctx, userID, chatID, text)
	//Отправить ответ подтверждение
	totalmsg := strconv.Itoa(h.db.Count())
	confirmation := "Cообщение сохранено! Всего Сообщений:" + totalmsg + "\n\n"

	if err := h.adapter.SendMessage(ctx, chatID, confirmation); err != nil {
		log.WithError(err).Error("Failed to send confirmation message")
		return
	}
	log.Info("Confirmation message sent successfully")

}

func (h *Handler) SaveMessage(ctx context.Context, userID int64, chatID int64, text string) {
	h.db.Save(ctx, userID, chatID, text)
}
