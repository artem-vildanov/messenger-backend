package storage

import (
	"context"
	"messenger/internal/domain/models"
	"messenger/internal/infrastructure/errors"

	"github.com/jmoiron/sqlx"
)

type ChatStorage struct {
	*AbstractStorage[models.ChatModel]
}

func NewChatStorage(pg *sqlx.DB) *ChatStorage {
	return &ChatStorage{
		&AbstractStorage[models.ChatModel]{pg},
	}
}

func (r *ChatStorage) GetChatsByUserId(
	ctx context.Context,
	userId, limit, offset int,
) (
	[]*models.ChatModel,
	*errors.Error,
) {
	// $1 - userId
	// $2 - limit
	// $3 - offset
	sql := `
WITH messages_with_second_user AS (
    SELECT 
        m.*,
        CASE 
            WHEN m.sender_id = $1 THEN m.receiver_id
            ELSE m.sender_id
        END AS second_user_id
    FROM messages m
    WHERE $1 IN (m.sender_id, m.receiver_id)
),

grouped_messages AS (
    SELECT 
        second_user_id,
        MAX(created_at) AS last_message_date,
        -- Замена функции get_last_message_text на подзапрос
        (SELECT text
         FROM messages
         WHERE (sender_id = $1 AND receiver_id = second_user_id)
            OR (sender_id = second_user_id AND receiver_id = $1)
         ORDER BY created_at DESC
         LIMIT 1) AS last_message_text,
        -- Замена функции get_unread_messages_count на подзапрос
        (SELECT COUNT(*)
         FROM messages
         WHERE receiver_id = $1 
           AND sender_id = second_user_id
           AND unread = true) AS unread_messages_count
    FROM messages_with_second_user
    GROUP BY second_user_id
)

SELECT 
    g.second_user_id,
    u.username AS second_user_name,
    g.last_message_date,
    g.last_message_text,
    g.unread_messages_count
FROM grouped_messages g
JOIN users u ON u.id = g.second_user_id
ORDER BY g.last_message_date DESC
LIMIT $2
OFFSET $3;
	`

	chats, err := r.findSlice(ctx, sql, userId, limit, offset)
	if err != nil {
		return nil, err.WithField("userId", userId).
			WithLogMessage("failed to find chats by userId")
	}

	return chats, nil
}
