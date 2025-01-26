package ws_utils

import (
	"context"
	"errors"
	appErrors "messenger/internal/infrastructure/errors"
	"messenger/internal/infrastructure/utils/handler_utils"
	"messenger/internal/infrastructure/utils/mapping_utils"
	"os"

	"github.com/gorilla/websocket"
)

type MessageType string

const (
	message MessageType = "message"
	ping    MessageType = "ping"
)

type WsHandler func(*WsContext) error

type WsContext struct {
	handler_utils.HandlerContext
	Conn  *websocket.Conn
	WsCtx context.Context // separate context for websockets
}

type WsMessage[T any] struct {
	Type MessageType `json:"type"`
	Dto  T           `json:"dto"`
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
	var incoming WsMessage[T]

	if err := conn.ReadJSON(&incoming); err != nil {
		if websocket.IsCloseError(
			err,
			websocket.CloseNormalClosure,
			websocket.CloseGoingAway,
		) {
			return incoming.Dto, appErrors.Wrap(
				appErrors.WsConnClosed,
				err,
				errors.New("Read"),
			)
		}
		if os.IsTimeout(err) {
			return incoming.Dto, appErrors.Wrap(
				appErrors.ErrTimeout,
				err,
				errors.New("Read"),
			)
		}
		return incoming.Dto, appErrors.Wrap(
			appErrors.ErrInternal,
			err,
			errors.New("Read"),
		)
	}

	if incoming.Type == ping {
		if err := Write(conn, "Pong"); err != nil {
			return incoming.Dto, appErrors.Wrap(err, errors.New("Read"))
		}
		return Read[T](conn)
	}

	if err := mapping_utils.ValidateRequestModel(incoming); err != nil {
		return incoming.Dto, appErrors.Wrap(err, errors.New("Read"))
	}

	return incoming.Dto, nil
}

// blocking operation
func Write(conn *websocket.Conn, outgoing any) error {
	if err := conn.WriteJSON(WsMessage[any]{
		Type: message,
		Dto:  outgoing,
	}); err != nil {
		return appErrors.Wrap(
			appErrors.ErrInternal,
			err,
			errors.New("Write"),
		)
	}
	return nil
}
