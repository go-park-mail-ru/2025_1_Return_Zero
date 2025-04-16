package user

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
	mock_user "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/user/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func setupHandler(t *testing.T) (*mock_user.MockUsecase, *UserHandler, context.Context) {
	ctrl := gomock.NewController(t)
	mockUsecase := mock_user.NewMockUsecase(ctrl)
	handler := NewUserHandler(mockUsecase)

	logger, _ := zap.NewDevelopment()
	sugarLogger := logger.Sugar()
	ctx := context.WithValue(context.Background(), helpers.LoggerKey{}, sugarLogger)

	return mockUsecase, handler, ctx
}

func TestUserHandler_Signup(t *testing.T) {
	mockUsecase, handler, ctx := setupHandler(t)
	defer gomock.NewController(t).Finish()

	tests := []struct {
		name               string
		inputJSON          string
		mockBehavior       func()
		expectedStatusCode int
		expectedResponse   string
		expectCookie       bool // Новое поле для указания, ожидаем ли cookie
	}{
		{
			name: "Success",
			inputJSON: `{
				"username": "testuser",
				"email": "test@example.com",
				"password": "password123"
			}`,
			mockBehavior: func() {
				expectedUser := &usecaseModel.User{
					Username: "testuser",
					Email:    "test@example.com",
					Password: "password123",
				}
				returnedUser := &usecaseModel.User{
					Username:  "testuser",
					Email:     "test@example.com",
					AvatarUrl: "default_avatar.png",
				}
				mockUsecase.EXPECT().
					CreateUser(gomock.Any(), gomock.Eq(expectedUser)).
					Return(returnedUser, "test_session_id", nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"status":200,"body":{"username":"testuser","email":"test@example.com","avatar_url":"default_avatar.png"}}`,
			expectCookie:       true, // Только в успешном сценарии ожидаем cookie
		},
		{
			name:               "Invalid JSON",
			inputJSON:          `{"username": "testuser", "email": "test@example.com"`, // Неполный JSON
			mockBehavior:       func() {},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"status":400,"error":"unexpected EOF"}`,
			expectCookie:       false, // Не ожидаем cookie
		},
		{
			name: "Username is default_avatar",
			inputJSON: `{
				"username": "default_avatar",
				"email": "test@example.com",
				"password": "password123"
			}`,
			mockBehavior:       func() {},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"status":400,"error":"Wrong username"}`,
			expectCookie:       false, // Не ожидаем cookie
		},
		{
			name: "Validation Failed",
			inputJSON: `{
				"username": "",
				"email": "invalid-email",
				"password": "short"
			}`,
			mockBehavior:       func() {},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"status":400,"error":"validation failed"}`,
			expectCookie:       false, // Не ожидаем cookie
		},
		{
			name: "CreateUser Error",
			inputJSON: `{
				"username": "testuser",
				"email": "test@example.com",
				"password": "password123"
			}`,
			mockBehavior: func() {
				expectedUser := &usecaseModel.User{
					Username: "testuser",
					Email:    "test@example.com",
					Password: "password123",
				}
				mockUsecase.EXPECT().
					CreateUser(gomock.Any(), gomock.Eq(expectedUser)).
					Return(nil, "", errors.New("user already exists"))
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"status":500,"error":"user already exists"}`,
			expectCookie:       false, // Не ожидаем cookie
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Настраиваем поведение мока
			tt.mockBehavior()

			// Создаем запрос
			req := httptest.NewRequest(http.MethodPost, "/auth/signup", bytes.NewBufferString(tt.inputJSON))
			req.Header.Set("Content-Type", "application/json")

			// Добавляем логгер в контекст запроса
			req = req.WithContext(ctx)

			// Создаем ResponseRecorder для записи ответа
			w := httptest.NewRecorder()

			// Вызываем тестируемую функцию
			handler.Signup(w, req)

			// Проверяем код ответа
			assert.Equal(t, tt.expectedStatusCode, w.Code)

			// Проверяем тело ответа
			if tt.expectedResponse != "" {
				assert.JSONEq(t, tt.expectedResponse, w.Body.String())
			}

			// Проверяем cookie только если это ожидается в данном тесте
			if tt.expectCookie {
				cookies := w.Result().Cookies()
				assert.GreaterOrEqual(t, len(cookies), 1, "Должна быть хотя бы одна cookie")

				var sessionCookie *http.Cookie
				for _, cookie := range cookies {
					if cookie.Name == "session_id" {
						sessionCookie = cookie
						break
					}
				}

				assert.NotNil(t, sessionCookie, "Должна быть cookie с именем session_id")
				assert.Equal(t, "session_id", sessionCookie.Name)
				assert.Equal(t, "test_session_id", sessionCookie.Value)
				assert.Equal(t, "/", sessionCookie.Path)
				assert.True(t, sessionCookie.HttpOnly)
			}
		})
	}
}

func TestUserHandler_Login(t *testing.T) {
	mockUsecase, handler, ctx := setupHandler(t)
	defer gomock.NewController(t).Finish()

	tests := []struct {
		name               string
		inputJSON          string
		mockBehavior       func()
		expectedStatusCode int
		expectedResponse   string
		expectCookie       bool
	}{
		{
			name: "Success",
			inputJSON: `{
				"email": "test@example.com",
				"password": "password123"
			}`,
			mockBehavior: func() {
				expectedLogin := &usecaseModel.User{
					Email:    "test@example.com",
					Password: "password123",
				}
				returnedUser := &usecaseModel.User{
					Username:  "testuser",
					Email:     "test@example.com",
					AvatarUrl: "default_avatar.png",
				}
				mockUsecase.EXPECT().
					LoginUser(gomock.Any(), gomock.Eq(expectedLogin)).
					Return(returnedUser, "test_session_id", nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"status":200,"body":{"username":"testuser","email":"test@example.com","avatar_url":"default_avatar.png"}}`,
			expectCookie:       true,
		},
		{
			name:               "Invalid JSON",
			inputJSON:          `{"email": "test@example.com"`, // Неполный JSON
			mockBehavior:       func() {},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"status":400,"error":"unexpected EOF"}`,
			expectCookie:       false,
		},
		{
			name: "Validation Failed",
			inputJSON: `{
				"email": "",
				"password": ""
			}`,
			mockBehavior:       func() {},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"status":400,"error":"validation failed"}`,
			expectCookie:       false,
		},
		{
			name: "Invalid Credentials",
			inputJSON: `{
				"email": "test@example.com",
				"password": "wrong_password"
			}`,
			mockBehavior: func() {
				expectedLogin := &usecaseModel.User{
					Email:    "test@example.com",
					Password: "wrong_password",
				}
				mockUsecase.EXPECT().
					LoginUser(gomock.Any(), gomock.Eq(expectedLogin)).
					Return(nil, "", errors.New("invalid credentials"))
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"status":500,"error":"invalid credentials"}`,
			expectCookie:       false,
		},
		{
			name: "User Not Found",
			inputJSON: `{
				"email": "nonexistent@example.com",
				"password": "password123"
			}`,
			mockBehavior: func() {
				expectedLogin := &usecaseModel.User{
					Email:    "nonexistent@example.com",
					Password: "password123",
				}
				mockUsecase.EXPECT().
					LoginUser(gomock.Any(), gomock.Eq(expectedLogin)).
					Return(nil, "", errors.New("user not found"))
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"status":500,"error":"user not found"}`,
			expectCookie:       false,
		},
		{
			name: "Internal Server Error",
			inputJSON: `{
				"email": "test@example.com", 
				"password": "password123"
			}`,
			mockBehavior: func() {
				expectedLogin := &usecaseModel.User{
					Email:    "test@example.com",
					Password: "password123",
				}
				mockUsecase.EXPECT().
					LoginUser(gomock.Any(), gomock.Eq(expectedLogin)).
					Return(nil, "", errors.New("database error"))
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"status":500,"error":"database error"}`,
			expectCookie:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Настраиваем поведение мока
			tt.mockBehavior()

			// Создаем запрос
			req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(tt.inputJSON))
			req.Header.Set("Content-Type", "application/json")

			// Добавляем логгер в контекст запроса
			req = req.WithContext(ctx)

			// Создаем ResponseRecorder для записи ответа
			w := httptest.NewRecorder()

			// Вызываем тестируемую функцию
			handler.Login(w, req)

			// Проверяем код ответа
			assert.Equal(t, tt.expectedStatusCode, w.Code)

			// Проверяем тело ответа
			if tt.expectedResponse != "" {
				assert.JSONEq(t, tt.expectedResponse, w.Body.String())
			}

			// Проверяем cookie только если это ожидается в данном тесте
			if tt.expectCookie {
				cookies := w.Result().Cookies()
				assert.GreaterOrEqual(t, len(cookies), 1, "Должна быть хотя бы одна cookie")

				var sessionCookie *http.Cookie
				for _, cookie := range cookies {
					if cookie.Name == "session_id" {
						sessionCookie = cookie
						break
					}
				}

				assert.NotNil(t, sessionCookie, "Должна быть cookie с именем session_id")
				assert.Equal(t, "session_id", sessionCookie.Name)
				assert.Equal(t, "test_session_id", sessionCookie.Value)
				assert.Equal(t, "/", sessionCookie.Path)
				assert.True(t, sessionCookie.HttpOnly)
			}
		})
	}
}

func TestUserHandler_Logout(t *testing.T) {
	mockUsecase, handler, ctx := setupHandler(t)
	defer gomock.NewController(t).Finish()

	tests := []struct {
		name               string
		setupRequest       func() *http.Request
		mockBehavior       func(sessionID string)
		expectedStatusCode int
		expectedResponse   string
		checkCookie        bool
	}{
		{
			name: "Success",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
				cookie := &http.Cookie{
					Name:  "session_id",
					Value: "test_session_id",
				}
				req.AddCookie(cookie)
				return req.WithContext(ctx)
			},
			mockBehavior: func(sessionID string) {
				mockUsecase.EXPECT().
					Logout(gomock.Any(), gomock.Eq(sessionID)).
					Return(nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"status":200,"body":{"msg":"Successfully logged out"}}`,
			checkCookie:        true,
		},
		{
			name: "No Session Cookie",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
				return req.WithContext(ctx)
			},
			mockBehavior: func(sessionID string) {
				// Не ожидаем вызова Logout, так как cookie отсутствует
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"status":400,"error":"http: named cookie not present"}`,
			checkCookie:        false,
		},
		{
			name: "Logout Error",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
				cookie := &http.Cookie{
					Name:  "session_id",
					Value: "test_session_id",
				}
				req.AddCookie(cookie)
				return req.WithContext(ctx)
			},
			mockBehavior: func(sessionID string) {
				mockUsecase.EXPECT().
					Logout(gomock.Any(), gomock.Eq(sessionID)).
					Return(errors.New("session not found"))
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"status":500,"error":"session not found"}`,
			checkCookie:        false, // Изменено с true на false, так как в случае ошибки cookie не устанавливается
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Устанавливаем мок поведение
			req := tt.setupRequest()

			// Устанавливаем ожидания мока
			if req.Cookies() != nil && len(req.Cookies()) > 0 {
				tt.mockBehavior(req.Cookies()[0].Value)
			} else {
				tt.mockBehavior("")
			}

			// Создаем ResponseRecorder для записи ответа
			w := httptest.NewRecorder()

			// Вызываем тестируемую функцию
			handler.Logout(w, req)

			// Проверяем код ответа
			assert.Equal(t, tt.expectedStatusCode, w.Code)

			// Проверяем тело ответа
			if tt.expectedResponse != "" {
				assert.JSONEq(t, tt.expectedResponse, w.Body.String())
			}

			// Проверяем cookie только если это ожидается
			if tt.checkCookie {
				cookies := w.Result().Cookies()
				assert.GreaterOrEqual(t, len(cookies), 1, "Должна быть хотя бы одна cookie")

				var sessionCookie *http.Cookie
				for _, cookie := range cookies {
					if cookie.Name == "session_id" {
						sessionCookie = cookie
						break
					}
				}

				assert.NotNil(t, sessionCookie, "Должна быть cookie с именем session_id")
				assert.True(t, sessionCookie.Expires.Before(time.Now()), "Cookie должна быть просрочена")
			}
		})
	}
}

func TestUserHandler_CheckUser(t *testing.T) {
	mockUsecase, handler, ctx := setupHandler(t)
	defer gomock.NewController(t).Finish()

	tests := []struct {
		name               string
		setupRequest       func() *http.Request
		mockBehavior       func(sessionID string)
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "Success",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/user", nil)
				cookie := &http.Cookie{
					Name:  "session_id",
					Value: "test_session_id",
				}
				req.AddCookie(cookie)
				return req.WithContext(ctx)
			},
			mockBehavior: func(sessionID string) {
				returnedUser := &usecaseModel.User{
					Username:  "testuser",
					Email:     "test@example.com",
					AvatarUrl: "default_avatar.png",
				}
				mockUsecase.EXPECT().
					GetUserBySID(gomock.Any(), gomock.Eq(sessionID)).
					Return(returnedUser, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"status":200,"body":{"username":"testuser","email":"test@example.com","avatar_url":"default_avatar.png"}}`,
		},
		{
			name: "No Session Cookie",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/user", nil)
				return req.WithContext(ctx)
			},
			mockBehavior: func(sessionID string) {
				// Не ожидаем вызовов методов, так как функция должна завершиться с ошибкой раньше
			},
			expectedStatusCode: http.StatusOK, // Статус HTTP всегда 200, но в JSON будет код ошибки
			expectedResponse:   `{"status":400,"error":"http: named cookie not present"}`,
		},
		{
			name: "Invalid Session",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/user", nil)
				cookie := &http.Cookie{
					Name:  "session_id",
					Value: "invalid_session_id",
				}
				req.AddCookie(cookie)
				return req.WithContext(ctx)
			},
			mockBehavior: func(sessionID string) {
				mockUsecase.EXPECT().
					GetUserBySID(gomock.Any(), gomock.Eq(sessionID)).
					Return(nil, errors.New("session not found"))
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"status":500,"error":"session not found"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Подготавливаем запрос
			req := tt.setupRequest()

			// Устанавливаем ожидания мока
			if req.Cookies() != nil && len(req.Cookies()) > 0 {
				tt.mockBehavior(req.Cookies()[0].Value)
			} else {
				tt.mockBehavior("")
			}

			// Создаем ResponseRecorder для записи ответа
			w := httptest.NewRecorder()

			// Вызываем тестируемую функцию
			handler.CheckUser(w, req)

			// Проверяем код ответа
			assert.Equal(t, tt.expectedStatusCode, w.Code)

			// Проверяем тело ответа
			if tt.expectedResponse != "" {
				assert.JSONEq(t, tt.expectedResponse, w.Body.String())
			}
		})
	}
}
