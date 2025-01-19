package handlers

import (
	"context"
	"errors"
	"messenger/internal/domain/models"
	"messenger/internal/domain/services"
	appErrors "messenger/internal/infrastructure/errors"
	"messenger/internal/infrastructure/utils/handler_utils"
	"messenger/internal/infrastructure/utils/ws_utils"
	"messenger/internal/presentation/dto"
)

const (
	failedWriteMessage = "failed to write chat message to client by ws"
	failedWriteNotification = "failed to write message notification to client by ws"
	failedGetMyChats = "failed to get auth user chats"
	failedGetChatMessages = "failed to get chat messages"
)

type ChatGetter interface {
	GetChatsByUserId(
		ctx context.Context,
		userId,
		limit,
		offset int,
	) ([]*models.ChatModel, error)
}

type MessageStorage interface {
	GetChatMessages(
		ctx context.Context,
		firstUserId,
		secondUserId,
		limit,
		offset int,
	) ([]*models.MessageModel, error)
}

type ChatHandler struct {
	chatGetter ChatGetter
	messageRepository MessageStorage
	chatService    *services.ChatService
}

func NewChatHandler(
	chatGetter ChatGetter,
	messageRepository MessageStorage,
	chatService *services.ChatService,
) *ChatHandler {
	return &ChatHandler{
		chatGetter,
		messageRepository,
		chatService,
	}
}

func (h *ChatHandler) GetMyChats(
	handlerContext *handler_utils.HandlerContext,
) error {
	limit, offset, err := handlerContext.GetLimitOffset()
	if err != nil {
		return appErrors.Wrap(err, errors.New("GetMyChats"))
	}

	chats, err := h.chatGetter.GetChatsByUserId(
		handlerContext.Request.Context(),
		handlerContext.AuthUserId,
		limit,
		offset,
	)
	if err != nil {
		return appErrors.Wrap(err, errors.New("GetMyChats"))
	}

	handlerContext.Response().
		WithContent(dto.NewMultipleChatsResponse(chats)).
		Json()

	return nil
}

func (h *ChatHandler) GetChatMessages(
	handlerContext *handler_utils.HandlerContext,
) error {
	limit, offset, err := handlerContext.GetLimitOffset()
	if err != nil {
		return appErrors.Wrap(err, errors.New("GetChatMessages"))
	}

	secondUserId, err := handlerContext.PathParams.GetInteger("userId")
	if err != nil {
		return appErrors.Wrap(err, errors.New("GetChatMessages"))
	}

	firstUserId := handlerContext.AuthUserId

	messages, err := h.messageRepository.GetChatMessages(
		handlerContext.Request.Context(),
		firstUserId,
		secondUserId,
		limit,
		offset,
	)
	if err != nil {
		return appErrors.Wrap(err, errors.New("GetChatMessages"))
	}

	handlerContext.Response().
		WithContent(dto.NewMultipleMessagesResponse(messages)).
		Json()

	return nil
}

// reading from client his outgoing message
func (h *ChatHandler) ReadChatMessage(
	wsContext *ws_utils.WsContext,
) error {
	// blocking operation
	createdMessage, err := ws_utils.Read[*dto.CreateMessageRequest](wsContext.Conn)
	if err != nil {
		return err
	}

	if err := h.chatService.PublishMessage(
		wsContext.WsCtx,
		createdMessage.ToDomain(),
	); err != nil {
		return appErrors.Wrap(err, errors.New("ReadChatMessage"))
	}

	return nil
}

// writing to client his incoming messages
func (h *ChatHandler) WriteChatMessages(
	wsContext *ws_utils.WsContext,
) error {
	firstUserId := wsContext.HandlerContext.AuthUserId
	secondUserId, err := wsContext.HandlerContext.PathParams.GetInteger("userId")
	if err != nil {
		return appErrors.Wrap(err, errors.New("WriteChatMessage"))
	}

	// blocking operation
	if err := h.chatService.SubscribeMessages(
		wsContext.WsCtx,
		wsContext.Conn,
		firstUserId,
		secondUserId,
	); err != nil {
		return appErrors.Wrap(err, errors.New("WriteChatMessage"))
	}

	return nil
}

func (h *ChatHandler) WriteMessageNotifications(
	wsContext *ws_utils.WsContext,
) error {
	if err := h.chatService.SubscribeMessageNotifications(
		wsContext.WsCtx,
		wsContext.Conn,
		wsContext.HandlerContext.AuthUserId,
	); err != nil {
		return appErrors.Wrap(err, errors.New("WriteMessageNotifications"))
	}
	
	return nil
}
