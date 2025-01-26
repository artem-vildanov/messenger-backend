package tests

import (
	"fmt"
	"log"
	"messenger/internal/domain/models"
	appErrors "messenger/internal/infrastructure/errors"
	"messenger/internal/infrastructure/utils/ws_utils"
	"messenger/internal/presentation/dto"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

func Test_GetChatMessages(t *testing.T) {
	firstUserData := createUser(t, "test-getchatmessages-1", "qwe123")
	secondUserData := createUser(t, "test-getchatmessages-2", "qwe123")
	thirdUserData := createUser(t, "test-getchatmessages-3", "qwe123")
	fourthUserData := createUser(t, "test-getchatmessages-4", "qwe123")
	fifthUserData := createUser(t, "test-getchatmessages-5", "qwe123")

	fromFirstToSecText1 := "msg from first to second 1"
	fromFirsttToSecId1 := createMessage(
		t,
		firstUserData.UserModel.Id,
		secondUserData.UserModel.Id,
		fromFirstToSecText1,
		time.Now().In(time.UTC).Add(time.Minute),
	)

	fromFirstToSecText2 := "msg from first to second 2"
	fromFirstToSecId2 := createMessage(
		t,
		firstUserData.UserModel.Id,
		secondUserData.UserModel.Id,
		fromFirstToSecText2,
		time.Now().In(time.UTC).Add(time.Minute*2),
	)

	fromFirstToThirdText1 := "msg from first to third 1"
	fromFirstToThirdId1 := createMessage(
		t,
		firstUserData.UserModel.Id,
		thirdUserData.UserModel.Id,
		fromFirstToThirdText1,
		time.Now().In(time.UTC).Add(time.Minute*3),
	)

	fromFourthToFifthText1 := "msg from fourth to fifth 1"
	fromFourthToFifthId1 := createMessage(
		t,
		fourthUserData.UserModel.Id,
		fifthUserData.UserModel.Id,
		fromFourthToFifthText1,
		time.Now().In(time.UTC).Add(time.Minute*4),
	)

	fromFourthToFifthText2 := "msg from fourth to fifth 2"
	fromFourthToFifthId2 := createMessage(
		t,
		fourthUserData.UserModel.Id,
		fifthUserData.UserModel.Id,
		fromFourthToFifthText2,
		time.Now().In(time.UTC).Add(time.Minute*5),
	)

	fromFifthToFourthText1 := "msg from fifth to fourth 1"
	fromFifthToFourthId1 := createMessage(
		t,
		fifthUserData.UserModel.Id,
		fourthUserData.UserModel.Id,
		fromFifthToFourthText1,
		time.Now().In(time.UTC).Add(time.Minute*6),
	)

	testCases := []struct {
		name                  string
		userId                int
		userSessionId         string
		partnerId             int
		expectedMessagesIds   []int
		expectedMessagesTexts []string
	}{
		{
			name:                  "messages between first and second",
			userId:                firstUserData.UserModel.Id,
			userSessionId:         firstUserData.SessionId,
			partnerId:             secondUserData.UserModel.Id,
			expectedMessagesIds:   []int{fromFirsttToSecId1, fromFirstToSecId2},
			expectedMessagesTexts: []string{fromFirstToSecText1, fromFirstToSecText2},
		},
		{
			name:                  "messages between first and third",
			userId:                firstUserData.UserModel.Id,
			userSessionId:         firstUserData.SessionId,
			partnerId:             thirdUserData.UserModel.Id,
			expectedMessagesIds:   []int{fromFirstToThirdId1},
			expectedMessagesTexts: []string{fromFirstToThirdText1},
		},
		{
			name:          "messages between fourth and fifth",
			userId:        fourthUserData.UserModel.Id,
			userSessionId: fourthUserData.SessionId,
			partnerId:     fifthUserData.UserModel.Id,
			expectedMessagesIds: []int{
				fromFourthToFifthId1,
				fromFourthToFifthId2,
				fromFifthToFourthId1,
			},
			expectedMessagesTexts: []string{
				fromFourthToFifthText1,
				fromFourthToFifthText2,
				fromFifthToFourthText1,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			res := Get[[]dto.MessageResponse](
				t,
				fmt.Sprintf("/chats/%d", testCase.partnerId),
				WithQueryParams([]QueryParam{
					{
						Key:   "limit",
						Value: "10",
					},
					{
						Key:   "offset",
						Value: "0",
					},
				}...),
				WithSessionId(testCase.userSessionId),
			)

			assert.Equal(t, http.StatusOK, res.response.StatusCode)

			for actualMessageIndex := range res.body {
				assert.Equal(
					t,
					testCase.expectedMessagesIds[actualMessageIndex],
					res.body[actualMessageIndex].Id,
				)
				assert.Equal(
					t,
					testCase.expectedMessagesTexts[actualMessageIndex],
					res.body[actualMessageIndex].Text,
				)
			}
		})
	}
}

func Test_GetChats(t *testing.T) {
	firstUserData := createUser(t, "test-getchats-1", "qwe123")
	secondUserData := createUser(t, "test-getchats-2", "qwe123")
	thirdUserData := createUser(t, "test-getchats-3", "qwe123")
	fourthUserData := createUser(t, "test-getchats-4", "qwe123")
	fifthUserData := createUser(t, "test-getchats-5", "qwe123")

	fromFirstToSecText1 := "msg from first to second 1"
	fromFirstToSecTime1 := time.Now().In(time.UTC).Add(-24 * time.Hour)
	_ = createMessage(
		t,
		firstUserData.UserModel.Id,
		secondUserData.UserModel.Id,
		fromFirstToSecText1,
		fromFirstToSecTime1,
	)

	fromFirstToSecText2 := "msg from first to second 2"
	fromFirstToSecTime2 := time.Now().In(time.UTC).Add(-12 * time.Hour)
	_ = createMessage(
		t,
		firstUserData.UserModel.Id,
		secondUserData.UserModel.Id,
		fromFirstToSecText2,
		fromFirstToSecTime2,
	)

	fromFirstToThirdText1 := "msg from first to third 1"
	fromFirstToThirdTime1 := time.Now().In(time.UTC).Add(-24 * time.Hour)
	_ = createMessage(
		t,
		firstUserData.UserModel.Id,
		thirdUserData.UserModel.Id,
		fromFirstToThirdText1,
		fromFirstToThirdTime1,
	)

	fromFourthToFifthText1 := "msg from fourth to fifth 1"
	fromFourthToFifthTime1 := time.Now().In(time.UTC).Add(-24 * time.Hour)
	_ = createMessage(
		t,
		fourthUserData.UserModel.Id,
		fifthUserData.UserModel.Id,
		fromFourthToFifthText1,
		fromFourthToFifthTime1,
	)

	fromFourthToFifthText2 := "msg from fourth to fifth 2"
	fromFourthToFifthTime2 := time.Now().In(time.UTC).Add(-12 * time.Hour)
	_ = createMessage(
		t,
		fourthUserData.UserModel.Id,
		fifthUserData.UserModel.Id,
		fromFourthToFifthText2,
		fromFourthToFifthTime2,
	)

	fromFifthToFourthText1 := "msg from fifth to fourth 1"
	fromFifthToFourthTime1 := time.Now().In(time.UTC).Add(-6 * time.Hour)
	_ = createMessage(
		t,
		fifthUserData.UserModel.Id,
		fourthUserData.UserModel.Id,
		fromFifthToFourthText1,
		fromFifthToFourthTime1,
	)

	testCases := []struct {
		name          string
		user          *UserData
		expectedChats []*models.ChatModel
	}{
		{
			name: "first user chats",
			user: firstUserData,
			expectedChats: []*models.ChatModel{
				{
					UserID:              secondUserData.UserModel.Id,
					Username:            secondUserData.UserModel.Username,
					LastMessageDate:     fromFirstToSecTime2,
					LastMessageText:     fromFirstToSecText2,
					UnreadMessagesCount: 0,
				},
				{
					UserID:              thirdUserData.UserModel.Id,
					Username:            thirdUserData.UserModel.Username,
					LastMessageDate:     fromFirstToThirdTime1,
					LastMessageText:     fromFirstToThirdText1,
					UnreadMessagesCount: 0,
				},
			},
		},
		{
			name: "second user chats",
			user: secondUserData,
			expectedChats: []*models.ChatModel{
				{
					UserID:              firstUserData.UserModel.Id,
					Username:            firstUserData.UserModel.Username,
					LastMessageDate:     fromFirstToSecTime2,
					LastMessageText:     fromFirstToSecText2,
					UnreadMessagesCount: 2,
				},
			},
		},
		{
			name: "third user chats",
			user: thirdUserData,
			expectedChats: []*models.ChatModel{
				{
					UserID:              firstUserData.UserModel.Id,
					Username:            firstUserData.UserModel.Username,
					LastMessageDate:     fromFirstToThirdTime1,
					LastMessageText:     fromFirstToThirdText1,
					UnreadMessagesCount: 1,
				},
			},
		},
		{
			name: "fourth user chats",
			user: fourthUserData,
			expectedChats: []*models.ChatModel{
				{
					UserID:              fifthUserData.UserModel.Id,
					Username:            fifthUserData.UserModel.Username,
					LastMessageDate:     fromFifthToFourthTime1,
					LastMessageText:     fromFifthToFourthText1,
					UnreadMessagesCount: 1,
				},
			},
		},
		{
			name: "fifth user chats",
			user: fifthUserData,
			expectedChats: []*models.ChatModel{
				{
					UserID:              fourthUserData.UserModel.Id,
					Username:            fourthUserData.UserModel.Username,
					LastMessageDate:     fromFifthToFourthTime1,
					LastMessageText:     fromFifthToFourthText1,
					UnreadMessagesCount: 2,
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			res := Get[[]*dto.ChatResponse](
				t,
				"/chats",
				WithSessionId(testCase.user.SessionId),
				WithQueryParams([]QueryParam{
					{
						Key:   "limit",
						Value: "10",
					},
					{
						Key:   "offset",
						Value: "0",
					},
				}...),
			)

			assert.Equal(t, http.StatusOK, res.response.StatusCode)

			for actualChatIndex := range res.body {
				assert.Equal(
					t,
					testCase.expectedChats[actualChatIndex].UserID,
					res.body[actualChatIndex].UserID,
				)
				assert.Equal(
					t,
					testCase.expectedChats[actualChatIndex].Username,
					res.body[actualChatIndex].Username,
				)
				assert.Equal(
					t,
					testCase.expectedChats[actualChatIndex].LastMessageText,
					res.body[actualChatIndex].LastMessageText,
				)
				assert.Equal(
					t,
					testCase.expectedChats[actualChatIndex].UnreadMessagesCount,
					res.body[actualChatIndex].UnreadMessagesCount,
				)
			}
		})
	}
}

func Test_WS(t *testing.T) {
	var group errgroup.Group

	firstUserData := createUser(t, "test-ws-1", "qwe123")
	secondUserData := createUser(t, "test-ws-2", "qwe123")
	thirdUserData := createUser(t, "test-ws-3", "qwe123")

	msgFromFirstToSec1 := "msg from first to sec 1"
	msgFromFirstToSec2 := "msg from first to sec 2"
	msgFromSecToFirst1 := "msg from sec to first 1"
	msgFromSecToFirst2 := "msg from sec to first 2"
	/*
		# первая секунда

		запускаем три горутины в которых -
		первый заходит в чат к второму
		второй к первому
		третий к второму

		# вторая секунда

		первый отправляет сообщения второму

		# третья секунда

		проверяем, получил ли второй сообщения от первого
		второй шлет сообщения первому

		# четвертая секунда

		проверяем, получил ли первый сообщения от второго
		проверям, что у третьего нет чужих сообщений

		завершаем тест
	*/

	// first
	group.Go(func() error {
		conn := Ws(
			t,
			fmt.Sprintf(
				"/chat-messages/%d",
				secondUserData.UserModel.Id,
			),
			firstUserData.SessionId,
		)
		defer func() {
			closeWsConn(conn)
			log.Println("first exited")
		}()

		time.Sleep(time.Second)

		err := ws_utils.Write(conn, &dto.CreateMessageRequest{
			SenderId:   firstUserData.UserModel.Id,
			ReceiverId: secondUserData.UserModel.Id,
			Text:       msgFromFirstToSec1,
		})
		if err != nil {
			log.Println("Write error: ", err.Error())
			return err
		}

		err = ws_utils.Write(conn, &dto.CreateMessageRequest{
			SenderId:   firstUserData.UserModel.Id,
			ReceiverId: secondUserData.UserModel.Id,
			Text:       msgFromFirstToSec2,
		})
		if err != nil {
			log.Println("Write error: ", err.Error())
			return err
		}

		log.Println("first sent messages to second")

		time.Sleep(time.Second)
		time.Sleep(time.Second)

		log.Println("first started reading messages from second")
		msg1, err := ws_utils.Read[*dto.CreateMessageRequest](conn)
		if err != nil {
			log.Println("Read error: ", err.Error())
			return err
		}

		msg2, err := ws_utils.Read[*dto.CreateMessageRequest](conn)
		if err != nil {
			log.Println("Read error: ", err.Error())
			return err
		}
		log.Println("first read messages from second")

		assert.Equal(t, msgFromFirstToSec1, msg1.Text)
		assert.Equal(t, msgFromFirstToSec2, msg2.Text)

		return nil
	})

	// second
	group.Go(func() error {
		conn := Ws(
			t,
			fmt.Sprintf(
				"/chat-messages/%d",
				firstUserData.UserModel.Id,
			),
			secondUserData.SessionId,
		)
		defer func() {
			log.Println("second exited")
			closeWsConn(conn)
		}()

		time.Sleep(time.Second)
		time.Sleep(time.Second)

		log.Println("second started reading messages from first")
		msg1, err := ws_utils.Read[*dto.CreateMessageRequest](conn)
		if err != nil {
			return err
		}

		msg2, err := ws_utils.Read[*dto.CreateMessageRequest](conn)
		if err != nil {
			return err
		}

		assert.Equal(t, msgFromFirstToSec1, msg1.Text)
		assert.Equal(t, msgFromFirstToSec2, msg2.Text)

		ws_utils.Write(conn, &dto.CreateMessageRequest{
			SenderId:   secondUserData.UserModel.Id,
			ReceiverId: firstUserData.UserModel.Id,
			Text:       msgFromSecToFirst1,
		})

		ws_utils.Write(conn, &dto.CreateMessageRequest{
			SenderId:   secondUserData.UserModel.Id,
			ReceiverId: firstUserData.UserModel.Id,
			Text:       msgFromSecToFirst2,
		})

		log.Println("second wrote messages to first")
		time.Sleep(time.Second)

		return nil
	})

	// third
	group.Go(func() error {
		conn := Ws(
			t,
			fmt.Sprintf(
				"/chat-messages/%d",
				secondUserData.UserModel.Id,
			),
			thirdUserData.SessionId,
		)
		defer func() {
			closeWsConn(conn)
			log.Println("third exited")
		}()

		time.Sleep(time.Second)
		time.Sleep(time.Second)
		time.Sleep(time.Second)

		conn.SetReadDeadline(time.Now().Add(1 * time.Second))

		msg, err := ws_utils.Read[*dto.CreateMessageRequest](conn)

		if err != nil {
			unwrappedErr, _ := appErrors.Unwrap(err)

			if unwrappedErr.Error.ResponseMessage == 
				appErrors.ErrTimeout.ResponseMessage {
				return nil
			} 

			log.Printf("third got unexpected error: %v", err)
			return err
		}

		return fmt.Errorf(
			`the third user should not have received the message;
			got message - %s`,
			msg.Text,
		)
	})

	if err := group.Wait(); err != nil {
		t.Errorf("got error: %s", err.Error())
	}
}
