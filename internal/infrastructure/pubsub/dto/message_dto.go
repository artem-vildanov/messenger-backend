package dto

import (
	"messenger/internal/domain/models"
	"time"
)

type MessageDto struct {
	Id         int       `json:"id"`
	Unread     bool      `json:"unread"`
	SenderId   int       `json:"senderId"`
	ReceiverId int       `json:"receiverId"`
	Text       string    `json:"text"`
	CreatedAt  time.Time `json:"createdAt"`
}

func NewMessageDto(messageModel *models.MessageModel) *MessageDto {
	return &MessageDto{
		Id:         messageModel.Id,
		Unread:     messageModel.Unread,
		SenderId:   messageModel.SenderId,
		ReceiverId: messageModel.ReceiverId,
		Text:       messageModel.Text,
		CreatedAt:  messageModel.CreatedAt,
	}
}
