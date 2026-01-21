package models

type SendMessageRequest struct {
    ChatID  int64  `json:"chat_id" binding:"required"`
    Message string `json:"message" binding:"required"`
}

type BroadcastMessageRequest struct {
    Message string `json:"message" binding:"required"`
}

type MessageResponse struct {
    Success bool   `json:"success"`
    Message string `json:"message,omitempty"`
    Error   string `json:"error,omitempty"`
}
