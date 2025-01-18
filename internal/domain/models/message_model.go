package models

import (
	"time"
)

type CreateMessageModel struct {
	SenderId   int
	ReceiverId int
	Text       string
}

type MessageModel struct {
	Id         int       `json:"id" db:"id"`
	Unread     bool      `json:"unread" db:"unread"`
	SenderId   int       `json:"senderId" db:"sender_id"`
	ReceiverId int       `json:"receiverId" db:"receiver_id"`
	Text       string    `json:"text" db:"text"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
}

func NewMessageModel(
	id int, 
	createdAt time.Time, 
	createMessageModel *CreateMessageModel,
) *MessageModel {
	return &MessageModel{
		Id: id,
		Unread: true,
		SenderId: createMessageModel.SenderId,
		ReceiverId: createMessageModel.ReceiverId,
		Text: createMessageModel.Text,
		CreatedAt: createdAt,
	}
}