package interfaces 
// Требования к http-серверу
//     1. Авторизация путём проверки токена (он хранится в конфигурационном файле)
//     2. Функционал отправки сообщения конкретному пользователю
//     3. Функционал отправки сообщения всем пользователям

import (
    "context"
)
type IMessageService interface {
    SendMessage(ctx context.Context, chatID int64, message string) error
    BroadcastMessage(ctx context.Context, message string) (int, int, error) // отправлено, всего, ошибка
    
}