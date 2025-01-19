package ws_utils

import (
	"context"
	"errors"
	"log"
	appErrors "messenger/internal/infrastructure/errors"
	"messenger/internal/infrastructure/utils/handler_utils"
	"messenger/internal/infrastructure/utils/mapping_utils"

	"github.com/gorilla/websocket"
)

type WsHandler func(*WsContext) error

type WsContext struct {
	handler_utils.HandlerContext
	Conn  *websocket.Conn
	WsCtx context.Context // separate context for websockets
}

func NewWsContext(
	ctx context.Context,
	handlerContext *handler_utils.HandlerContext,
	conn *websocket.Conn,
) *WsContext {
	return &WsContext{
		HandlerContext: *handlerContext,
		Conn:           conn,
		WsCtx:          ctx,
	}
}

// blocking operation
func Read[T any](conn *websocket.Conn) (T, error) {
	incoming := new(T)

	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Println("failed to read: ", err.Error())
		return *incoming, appErrors.ErrBadRequest
	}

	log.Println("success", string(msg))
	return *incoming, appErrors.ErrBadRequest

	if err := conn.ReadJSON(&incoming); err != nil {
		log.Printf("failed to read from JSON: %s", err.Error())
		return *incoming, appErrors.Wrap(
			appErrors.ErrInternal,
			err,
			errors.New("Read"),
		)
	}

	if err := mapping_utils.ValidateRequestModel(incoming); err != nil {
		return *incoming, err
	}

	return *incoming, nil
}

// blocking operation
func Write[T any](conn *websocket.Conn, outgoing T) error {
	if err := conn.WriteJSON(outgoing); err != nil {
		log.Println("Write error json parse: ", err.Error())
		return appErrors.Wrap(
			appErrors.ErrInternal,
			err,
			errors.New("Write"),
		)
	}
	return nil
}
