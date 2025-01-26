package dto

import (
	"messenger/internal/domain/models"
	"time"
)


type MessageResponse struct {
	Id         int       `json:"id"`
	Unread     bool      `json:"unread"`
	SenderId   int       `json:"senderId"`
	ReceiverId int       `json:"receiverId"`
	Text       string    `json:"text"`
	CreatedAt  time.Time `json:"createdAt"`
}

func NewMessageResponse(messageModel *models.MessageModel) *MessageResponse {
	return &MessageResponse{
		Id: messageModel.Id,
		Unread: messageModel.Unread,
		SenderId: messageModel.SenderId,
		ReceiverId: messageModel.ReceiverId,
		Text: messageModel.Text,
		CreatedAt: messageModel.CreatedAt,
	}
}

func NewMultipleMessagesResponse(messagesModels []*models.MessageModel) []*MessageResponse {
	responsesModels := make([]*MessageResponse, 0, len(messagesModels))
	for _, messageModel := range messagesModels {
		responsesModels = append(responsesModels, NewMessageResponse(messageModel))
	}
	return responsesModels
}

type CreateMessageRequest struct {
	SenderId   int    `json:"senderId" validate:"required,gt=0"`
	ReceiverId int    `json:"receiverId" validate:"required,gt=0"`
	Text       string `json:"text" validate:"required,min=0"`
}

func (r *CreateMessageRequest) ToDomain() *models.CreateMessageModel {
	return &models.CreateMessageModel{
		SenderId: r.SenderId,
		ReceiverId: r.ReceiverId,
		Text: r.Text,
	}
}
