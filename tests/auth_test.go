package tests

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func authUrl(url string) string {
	return fmt.Sprintf("/auth%s", url)
}

func Test_Register(t *testing.T) {
	testCases := []struct {
		name       string
		username   string
		password   string
		expectCode int
	}{
		{
			name:       "successfull registrarion",
			username:   "test-register",
			password:   "123qwe",
			expectCode: http.StatusOK,
		},
		{
			name:       "invalid username",
			username:   "iv",
			password:   "qwe123",
			expectCode: http.StatusBadRequest,
		},
		{
			name:       "invalid password",
			username:   "ivan",
			password:   "qw",
			expectCode: http.StatusBadRequest,
		},
		{
			name:       "user already exists",
			username:   "test-register",
			password:   "qwe123",
			expectCode: http.StatusBadRequest,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := Post[struct{}](
				t,
				authUrl("/register"),
				map[string]any{
					"username": testCase.username,
					"password": testCase.password,
				},
			)

			ctx := context.Background()

			assert.Equal(t, testCase.expectCode, result.response.StatusCode)

			if result.response.StatusCode == http.StatusOK {
				sessionId := getSessionIdFromResponse(t, result.response.Cookies())

				user, appErr := app.StorageRegistry.
					UserStorage.
					GetByUsername(ctx, testCase.username)
				requireNoError(t, appErr)

				session, appErr := app.StorageRegistry.
					SessionStorage.
					GetSessionByUserId(ctx, user.Id)
				requireNoError(t, appErr)

				assert.Equal(t, sessionId, session.Id)
			}
		})
	}
}

func Test_Login(t *testing.T) {
	username := "login-test"
	password := "12345"
	_ = createUser(t, username, password)

	testCases := []struct {
		name       string
		username   string
		password   string
		expectCode int
	}{
		{
			name:       "successfull login",
			username:   username,
			password:   password,
			expectCode: http.StatusOK,
		},
		{
			name:       "invalid username",
			username:   "iv",
			password:   "qwe123",
			expectCode: http.StatusBadRequest,
		},
		{
			name:       "invalid password",
			username:   "ivan",
			password:   "qw",
			expectCode: http.StatusBadRequest,
		},
		{
			name:       "wrong password",
			username:   username,
			password:   "qwe123",
			expectCode: http.StatusUnauthorized,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			res := Post[struct{}](
				t,
				authUrl("/login"),
				map[string]any{
					"username": testCase.username,
					"password": testCase.password,
				},
			)

			assert.Equal(t, testCase.expectCode, res.response.StatusCode)

			ctx := context.Background()

			if res.response.StatusCode == http.StatusOK {
				sessionId := getSessionIdFromResponse(t, res.response.Cookies())

				session, appErr := app.StorageRegistry.
					SessionStorage.
					GetSessionById(ctx, sessionId)
				requireNoError(t, appErr)

				assert.Equal(t, sessionId, session.Id)
			}
		})
	}
}

func Test_Logout(t *testing.T) {
	username := "logout-test"
	password := "12345"
	userData := createUser(t, username, password)

	testCases := []struct {
		name       string
		sessionId  string
		expectCode int
	}{
		{
			name:       "successfull logout",
			sessionId:  userData.SessionId,
			expectCode: http.StatusNoContent,
		},
		{
			name:       "unauthorized",
			sessionId:  "123123",
			expectCode: http.StatusUnauthorized,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			res := Post[struct{}](
				t,
				authUrl("/logout"),
				map[string]any{},
				WithSessionId(testCase.sessionId),
			)

			assert.Equal(t, testCase.expectCode, res.response.StatusCode)

			ctx := context.Background()

			if res.response.StatusCode == http.StatusNoContent {
				_, appErr := app.StorageRegistry.
					SessionStorage.
					GetSessionById(ctx, userData.SessionId)
				if appErr == nil {
					t.Error("session must be removed after logout")
				}
			}
		})
	}
}
