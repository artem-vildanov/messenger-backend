package dto

import (
	"messenger/internal/domain/models"
	"time"
)

type ChatResponse struct {
	UserID              int       `json:"userId"`
	Username            string    `json:"username"`
	LastMessageDate     time.Time `json:"lastMessageDate"`
	LastMessageText     string    `json:"lastMessageText"`
	UnreadMessagesCount int       `json:"unreadMessagesCount"`
}

func NewChatResponse(chatModel *models.ChatModel) *ChatResponse {
	return &ChatResponse{
		UserID: chatModel.UserID,
		Username: chatModel.Username,
		LastMessageDate: chatModel.LastMessageDate,
		LastMessageText: chatModel.LastMessageText,
		UnreadMessagesCount: chatModel.UnreadMessagesCount,
	}
}

func NewMultipleChatsResponse(chatsModels []*models.ChatModel) []*ChatResponse {
	responseModels := make([]*ChatResponse, 0, len(chatsModels))
	for _, chatModel := range chatsModels {
		responseModels = append(responseModels, NewChatResponse(chatModel))
	}
	return responseModels
}
