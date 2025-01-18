package storage

import (
	"context"
	"messenger/internal/domain/models"
	"messenger/internal/infrastructure/errors"
	"time"

	"github.com/jmoiron/sqlx"
)

type MessageStorage struct {
	*AbstractStorage[models.MessageModel]
}

func NewMessageStorage(pg *sqlx.DB) *MessageStorage {
	return &MessageStorage{
		&AbstractStorage[models.MessageModel]{pg},
	}
}

func (r *MessageStorage) GetChatMessages(
	ctx context.Context,
	firstUserId, secondUserId, limit, offset int,
) (
	[]*models.MessageModel,
	*errors.Error,
) {
	sql := `
	select * from messages
	where (sender_id = $1 and receiver_id = $2) 
		or (sender_id = $2 and receiver_id = $1)
	limit $3
	offset $4;
	`

	messages, err := r.findSlice(
		ctx, 
		sql, 
		firstUserId, 
		secondUserId, 
		limit, 
		offset,
	)
	if err != nil {
		return nil, err.WithField("firstUserId", firstUserId).
			WithField("secondUserId", secondUserId).
			WithLogMessage("failed to MessageStorage.GetChatMessages")
	}

	return messages, nil
}

func (r *MessageStorage) SaveMessage(
	ctx context.Context,
	createMessageModel *models.CreateMessageModel,
	createdAt time.Time,
) (int, *errors.Error) {
	sql := `
		insert into messages (sender_id, receiver_id, text, unread, created_at)
		values ($1, $2, $3, $4, $5)
		returning id;
	`

	var messageId int
	if err := r.queryRow(
		ctx,
		sql,
		createMessageModel.SenderId,
		createMessageModel.ReceiverId,
		createMessageModel.Text,
		true,
		createdAt,
	).Scan(&messageId); err != nil {
		return 0, errors.InternalError().
			WithLogMessage(err.Error()).
			WithLogMessage("failed to save message").
			WithField("CreateMessageRequest", createMessageModel).
			WithField("createdAt", createdAt)
	}

	return messageId, nil
}

func (r *MessageStorage) MakeMessageRead(
	ctx context.Context,
	messageId int,
) *errors.Error {
	sql := `
		update messages
		set unread = false
		where id = $1;
	`

	if err := r.exec(ctx, sql, messageId); err != nil {
		return errors.InternalError().
			WithLogMessage(err.Error()).
			WithLogMessage("failed to make message read").
			WithField("messageId", messageId)
	}

	return nil
}
