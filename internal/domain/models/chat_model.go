package models

import (
	"time"
)

type ChatModel struct {
	UserID              int       `json:"userId" db:"second_user_id"`
	Username            string    `json:"username" db:"second_user_name"`
	LastMessageDate     time.Time `json:"lastMessageDate" db:"last_message_date"`
	LastMessageText     string    `json:"lastMessageText" db:"last_message_text"`
	UnreadMessagesCount int       `json:"unreadMessagesCount" db:"unread_messages_count"`
}
