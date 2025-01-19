package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"messenger/internal/bootstrap"
	"messenger/internal/domain/models"
	"net/http"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

type UserData struct {
	UserModel *models.UserModel
	SessionId string
}

type Response[T any] struct {
	body     T
	response *http.Response
}

func EmptyResponse[T any](response *http.Response) *Response[T] {
	return &Response[T]{
		body:     *new(T),
		response: response,
	}
}

type RequestOptionFunc func(*http.Request)

func WithSessionId(sessionId string) RequestOptionFunc {
	return func(r *http.Request) {
		r.Header.Add("Cookie", "sessionId="+sessionId)
	}
}

type QueryParam struct {
	Key   string
	Value string
}

func WithQueryParams(queries ...QueryParam) RequestOptionFunc {
	return func(r *http.Request) {
		q := r.URL.Query()
		for _, query := range queries {
			q.Add(query.Key, query.Value)
		}
		r.URL.RawQuery = q.Encode()
	}
}

func Ws(
	t *testing.T,
	url string,
	sessionId string,
) *websocket.Conn {
	addr := fmt.Sprintf(
		"ws://%s:%s/ws%s",
		app.Env.AppHost,
		app.Env.AppPort,
		url,
	)
	log.Println(addr)
	header := http.Header{}
	header.Add("Cookie", "sessionId="+sessionId)
	conn, _, err := websocket.DefaultDialer.Dial(addr, header)
	require.NoError(t, err)
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	return conn
}

func Post[T any](
	t *testing.T,
	url string,
	body map[string]interface{},
	options ...RequestOptionFunc,
) *Response[T] {
	req, err := createPostReq(url, body)
	if err != nil {
		require.NoError(t, err)
	}

	for _, option := range options {
		option(req)
	}

	return execRequest[T](t, req)
}

func Get[T any](
	t *testing.T,
	url string,
	options ...RequestOptionFunc,
) *Response[T] {
	req, err := createGetReq(url)
	require.NoError(t, err)

	for _, option := range options {
		option(req)
	}

	return execRequest[T](t, req)
}

func execRequest[T any](t *testing.T, req *http.Request) *Response[T] {
	req.Header.Set("Content-Type", "application/json")

	response, err := client.Do(req)
	require.NoError(t, err)
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return EmptyResponse[T](response)
	}

	if isVoidType[T]() {
		return EmptyResponse[T](response)
	}

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	var parsedBody struct {
		Message T `json:"message"`
	}
	err = json.Unmarshal(body, &parsedBody)
	require.NoError(t, err)

	return &Response[T]{
		body:     parsedBody.Message,
		response: response,
	}
}

func isVoidType[T any]() bool {
	var v T
	t := reflect.TypeOf(v)
	return t.Kind() == reflect.Struct && t.NumField() == 0
}

func createPostReq(url string, body map[string]any) (*http.Request, error) {
	addr := fmt.Sprintf(
		"http://%s:%s/api%s",
		app.Env.AppHost,
		app.Env.AppPort,
		url,
	)

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		"POST",
		addr,
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func createGetReq(url string) (*http.Request, error) {
	addr := fmt.Sprintf(
		"http://%s:%s/api%s",
		app.Env.AppHost,
		app.Env.AppPort,
		url,
	)

	req, err := http.NewRequest(
		"GET",
		addr,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func createUser(
	t *testing.T,
	username,
	password string,
) *UserData {
	response := Post[struct {
		UserId int `json:"userId"`
	}](
		t,
		"/auth/register",
		map[string]any{
			"username": username,
			"password": password,
		},
	)

	sessionId := getSessionIdFromResponse(t, response.response.Cookies())

	userModel, err := app.StorageRegistry.
		UserStorage.
		GetById(
			context.Background(),
			response.body.UserId,
		)
	require.NoError(t, err)

	return &UserData{
		UserModel: userModel,
		SessionId: sessionId,
	}
}

func createMessage(
	t *testing.T,
	senderId,
	receiverId int,
	text string,
	createdAt time.Time,
) (messageId int) {
	model := &models.CreateMessageModel{
		SenderId:   senderId,
		ReceiverId: receiverId,
		Text:       text,
	}

	messageId, err := app.StorageRegistry.MessageStorage.SaveMessage(
		context.Background(),
		model,
		createdAt,
	)
	require.NoError(t, err)

	return
}

func getSessionIdFromResponse(t *testing.T, cookies []*http.Cookie) string {
	for _, cookie := range cookies {
		if cookie.Name == "sessionId" {
			return cookie.Value
		}
	}
	t.Errorf("no session cookie provided in response: %v", cookies)
	t.FailNow()
	return ""
}

func resetTestDb() error {
	defer func() {
		if r := recover(); r != nil {
			log.Println("db reset paniced")
		}
	}()
	if _, err := app.StorageRegistry.PgConn.Exec(`
	DO
	$$
	DECLARE
		r RECORD;
	BEGIN
		-- Отключаем внешние ключи
		EXECUTE 'SET session_replication_role = ''replica''';

		-- Очищаем все таблицы
		FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = 'public') LOOP
			EXECUTE 'TRUNCATE TABLE ' || quote_ident(r.tablename) || ' CASCADE';
		END LOOP;

		-- Включаем внешние ключи обратно
		EXECUTE 'SET session_replication_role = ''origin''';
	END
	$$;
	`); err != nil {
		log.Printf("failed to reset test db: %s", err.Error())
		return err
	}

	log.Println("test db successfully reseted")
	return nil
}

func closeWsConn(conn *websocket.Conn) {
	err := conn.WriteMessage(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, "bye"),
	)
	if err != nil {
		log.Printf("failed to close connection: %s", err)
	}
	conn.Close()
}

var (
	client *http.Client
	app    *bootstrap.App
	mu     sync.Mutex
)
