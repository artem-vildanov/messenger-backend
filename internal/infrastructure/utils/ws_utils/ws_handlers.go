package ws_utils

import (
	"context"
	"errors"
	appError "messenger/internal/infrastructure/errors"
	"messenger/internal/infrastructure/utils/handler_utils"
	"messenger/internal/infrastructure/utils/router_utils"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

func WsHandlers(handlers ...WsHandler) router_utils.Handler {
	return func(handlerContext *handler_utils.HandlerContext) error {
		if len(handlers) == 0 {
			return nil
		}

		conn, err := NewWsConn(handlerContext)
		if err != nil {
			return err
		}

		ctx, cancel := context.WithCancel(context.Background())

		wsContext := NewWsContext(ctx, handlerContext, conn)
		errorChannel := make(chan error, 1)
		wg := new(sync.WaitGroup)

		for _, handler := range handlers {
			wg.Add(1)
			go RunHandler(
				ctx,
				cancel,
				wg,
				errorChannel,
				wsContext,
				handler,
			)
		}

		wg.Wait()
		cancel()

		select {
		case err := <-errorChannel:
			return err
		default:
			return nil
		}
	}
}

func NewWsConn(handlerContext *handler_utils.HandlerContext) (*websocket.Conn, error) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // todo add origin check
		},
	}

	conn, err := upgrader.Upgrade(
		handlerContext.ResponseWriter,
		handlerContext.Request,
		nil,
	)

	if err != nil {
		return nil, appError.Wrap(
			appError.ErrInternal,
			err,
			errors.New("NewWsConn"),
		)
	}

	return conn, nil
}

func RunHandler(
	ctx context.Context,
	cancel context.CancelFunc,
	wg *sync.WaitGroup,
	errorChannel chan error,
	wsContext *WsContext,
	wsHandler WsHandler,
) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			wsContext.Conn.Close()
			return
		default:
			if err := wsHandler(wsContext); err != nil {
				errorChannel <- err
				cancel()
				wsContext.Conn.Close()
				return
			}
		}
	}
}
