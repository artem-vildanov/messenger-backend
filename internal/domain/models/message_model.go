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
	Id               int       `db:"id"`
	Unread           bool      `db:"unread"`
	SenderId         int       `db:"sender_id"`
	ReceiverId       int       `db:"receiver_id"`
	Text             string    `db:"text"`
	CreatedAt        time.Time `db:"created_at"`
}

func NewMessageModel(
	id int,
	createdAt time.Time,
	createMessageModel *CreateMessageModel,
) *MessageModel {
	return &MessageModel{
		Id:         id,
		Unread:     true,
		SenderId:   createMessageModel.SenderId,
		ReceiverId: createMessageModel.ReceiverId,
		Text:       createMessageModel.Text,
		CreatedAt:  createdAt,
	}
}

