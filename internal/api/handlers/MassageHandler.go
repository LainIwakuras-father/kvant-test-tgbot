// internal/adapter/api/handlers/message_handler.go
package handlers_api

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/LainIwakuras-father/kvant-test-tgbot/internal/core/interfaces"
    "github.com/LainIwakuras-father/kvant-test-tgbot/internal/core/models"
)

type MessageHandler struct {
    messageService interfaces.IMessageService
}

func NewMessageHandler(messageService interfaces.IMessageService) *MessageHandler {
    return &MessageHandler{
        messageService: messageService,
    }
}

// SendMessage отправляет сообщение конкретному пользователю
func (h *MessageHandler) SendMessage(c *gin.Context) {
    var req models.SendMessageRequest
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, models.MessageResponse{
            Success: false,
            Error:   "Invalid request format: " + err.Error(),
        })
        return
    }
    
    if req.Message == "" {
        c.JSON(http.StatusBadRequest, models.MessageResponse{
            Success: false,
            Error:   "Message cannot be empty",
        })
        return
    }
    
    if err := h.messageService.SendMessage(c.Request.Context(), req.ChatID, req.Message); err != nil {
        c.JSON(http.StatusInternalServerError, models.MessageResponse{
            Success: false,
            Error:   "Failed to send message: " + err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, models.MessageResponse{
        Success: true,
        Message: "Message sent successfully to chat_id " + strconv.FormatInt(req.ChatID, 10),
    })
}

// BroadcastMessage отправляет сообщение всем пользователям
func (h *MessageHandler) BroadcastMessage(c *gin.Context, ) {
    var req models.BroadcastMessageRequest
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, models.MessageResponse{
            Success: false,
            Error:   "Invalid request format: " + err.Error(),
        })
        return
    }
    
    if req.Message == "" {
        c.JSON(http.StatusBadRequest, models.MessageResponse{
            Success: false,
            Error:   "Message cannot be empty",
        })
        return
    }
    
    sent, total, err := h.messageService.BroadcastMessage(c.Request.Context(), req.Message)
    if err != nil {
        c.JSON(http.StatusInternalServerError, models.MessageResponse{
            Success: false,
            Error:   "Failed to broadcast message: " + err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "message": "Message broadcasted successfully",
        "stats": gin.H{
            "sent":  sent,
            "total": total,
            "failed": total - sent,
        },
    })
}

