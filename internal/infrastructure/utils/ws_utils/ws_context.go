package ws_utils

import (
	"context"
	"log"
	"messenger/internal/infrastructure/errors"
	"messenger/internal/infrastructure/utils/handler_utils"
	"messenger/internal/infrastructure/utils/mapping_utils"

	"github.com/gorilla/websocket"
)

type WsHandler func(*WsContext) *errors.Error

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
func Read[T any](conn *websocket.Conn) (T, *errors.Error) {
	incoming := new(T)

	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Println("failed to read: ", err.Error())
		return *incoming, errors.BadRequestError()
	}

	log.Println("success", string(msg))
	return *incoming, errors.BadRequestError()

	if err := conn.ReadJSON(&incoming); err != nil {
		log.Printf("failed to read from JSON: %s", err.Error())
		return *incoming, errors.BadRequestError().
			WithResponseMessage("failed to parse json").
			WithLogMessage(err.Error()).
			WithOriginalError(err)
	}

	if err := mapping_utils.ValidateRequestModel(incoming); err != nil {
		return *incoming, err
	}

	return *incoming, nil
}

// blocking operation
func Write[T any](conn *websocket.Conn, outgoing T) *errors.Error {
	if err := conn.WriteJSON(outgoing); err != nil {
		log.Println("Write error json parse: ", err.Error())
		return errors.InternalError().
			WithLogMessage(err.Error(), "failed to write json")
	}
	return nil
}
