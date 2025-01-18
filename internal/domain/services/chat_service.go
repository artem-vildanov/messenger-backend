package services

import (
	"context"
	"fmt"
	"messenger/internal/domain/models"
	"messenger/internal/infrastructure/errors"
	"messenger/internal/infrastructure/pubsub"
	"messenger/internal/infrastructure/pubsub/dto"
	"messenger/internal/infrastructure/utils/ws_utils"
	"time"

	"github.com/gorilla/websocket"
)

const (
	failedToPublishMessage = "failed to publish message"
)

type MessageSaver interface {
	SaveMessage(
		context.Context,
		*models.CreateMessageModel,
		time.Time,
	) (int, *errors.Error)
}

type PubsubClient interface {
	SubscribeMessageNotifications(
		ctx context.Context,
		userId int,
	) <-chan *dto.PubsubDto[*dto.MessageDto]
	SubscribeMessages(
		ctx context.Context,
		chatId string,
	) <-chan *dto.PubsubDto[*dto.MessageDto]
	PublishMessage(
		ctx context.Context,
		messageModel *models.MessageModel,
		chatId string,
	) *errors.Error
}

type ChatService struct {
	pubsubClient PubsubClient
	messageSaver MessageSaver
}

func NewChatService(
	pubsubClient *pubsub.PubsubClient,
	messageSaver MessageSaver,
) *ChatService {
	return &ChatService{
		pubsubClient,
		messageSaver,
	}
}

func (s *ChatService) SubscribeMessages(
	ctx context.Context,
	conn *websocket.Conn,
	firstUserId,
	secondUserId int,
) *errors.Error {
	pubsubDtos := s.pubsubClient.SubscribeMessages(ctx, s.createChatId(firstUserId, secondUserId))
	if err := sendPubsubDtos(conn, pubsubDtos); err != nil {
		return err.WithField("firstUserId", firstUserId).
			WithField("secondUserId", secondUserId).
			WithLogMessage("failed while subscribed to messages")
	}
	return nil
}

func (s *ChatService) SubscribeMessageNotifications(
	ctx context.Context,
	conn *websocket.Conn,
	userId int,
) *errors.Error {
	pubsubDtos := s.pubsubClient.SubscribeMessageNotifications(ctx, userId)
	if err := sendPubsubDtos(conn, pubsubDtos); err != nil {
		return err.WithField("userId", userId).
			WithLogMessage("failed while subscribed to message notifications")
	}
	return nil
}

func (s *ChatService) PublishMessage(
	ctx context.Context,
	createdMessage *models.CreateMessageModel,
) *errors.Error {
	createdAt := time.Now().In(time.UTC)

	messageId, err := s.messageSaver.SaveMessage(ctx, createdMessage, createdAt)
	if err != nil {
		return err.WithLogMessage(failedToPublishMessage)
	}

	messageModel := models.NewMessageModel(
		messageId,
		createdAt,
		createdMessage,
	)

	chatId := s.createChatId(
		createdMessage.SenderId,
		createdMessage.ReceiverId,
	)

	if err := s.pubsubClient.PublishMessage(
		ctx, 
		messageModel, 
		chatId,
	); err != nil {
		return err.WithLogMessage(failedToPublishMessage)
	}

	return nil
}

func (s *ChatService) createChatId(firstUserId, secondUserId int) string {
	return fmt.Sprintf(
		"%d-%d",
		max(firstUserId, secondUserId),
		min(firstUserId, secondUserId),
	)
}

func sendPubsubDtos[T any](
	conn *websocket.Conn,
	pubsubDtosChannel <-chan *dto.PubsubDto[T],
) *errors.Error {
	for pubsubDto := range pubsubDtosChannel {
		if pubsubDto.Error != nil {
			return pubsubDto.Error
		}
		if err := ws_utils.Write(conn, pubsubDto.Payload); err != nil {
			return err
		}
	}
	return nil
}
