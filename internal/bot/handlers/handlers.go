package handlers_bot

import (
	"log"
	"strconv"

	"github.com/LainIwakuras-father/kvant-test-tgbot/internal/core/interfaces"
)

// поменять на интерфейсы
type Handler struct {
	adapter interfaces.IBot
	db      interfaces.IStorage
}

func NewHandler(adapter interfaces.IBot, db interfaces.IStorage) *Handler {
	return &Handler{
		adapter: adapter,
		db:      db,
	}
}

// HandleStart обрабатывает /start
func (h *Handler) HandleStart(username string, chatID int64) {
	
	if err:=h.db.AddUser(chatID,username); err !=nil{
		log.Printf("Ошибка добавления пользователя в базу")
	}

	log.Printf("👤 User added: %d (@%s)", chatID, username)


	text := `Привествую Друг!
Я бот сохраняющий и отправляющий твои сообщения кому либо!
Напиши Юзернейм или ID пользователя и текст который хочешь ему отправить
Формат:
9038487587 Привет, Друг!
	`
	if err := h.adapter.SendMessage(chatID, text); err != nil {
		log.Printf("Ошибка отправки /start: %v", err)
	}
}

// HandleTextMessage обрабатывает текстовые сообщения любые кроме команд
func (h *Handler) HandleTextMessage(userID int64, chatID int64, text string) {

	// Сохранить сообщение в памяти
	h.SaveMessage(userID, chatID, text)
	//Отправить ответ подтверждение
	totalmsg := strconv.Itoa(h.db.Count())
	confirmation := "Cообщение сохранено! Всего Сообщений:" + totalmsg + "\n\n"

	if err := h.adapter.SendMessage(chatID, confirmation); err != nil {
		log.Printf("Ошибка сохранения сообщения %v", err)
		return
	}

}

func (h *Handler) SaveMessage(userID int64, chatID int64, text string) {
	h.db.Save(userID, chatID, text)
}
