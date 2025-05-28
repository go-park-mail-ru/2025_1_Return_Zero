package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/config"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/ctxExtractor"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/customErrors"
	loggerPkg "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
	mock_jam "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/jam/mocks"
	deliveryModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/delivery"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func setupTestJamHandler(t *testing.T) (*mock_jam.MockUsecase, *JamHandler, *config.Config) {
	ctrl := gomock.NewController(t)
	mockUsecase := mock_jam.NewMockUsecase(ctrl)
	cfg := &config.Config{}
	handler := NewJamHandler(mockUsecase, cfg)
	return mockUsecase, handler, cfg
}

func TestCreateRoom(t *testing.T) {
	mockUsecase, handler, _ := setupTestJamHandler(t)

	tests := []struct {
		name           string
		userID         int64
		requestBody    interface{}
		mockBehavior   func()
		expectedStatus int
		expectedBody   interface{}
		hasUserInCtx   bool
	}{
		{
			name:   "Success",
			userID: 123,
			requestBody: deliveryModel.CreateJamRequest{
				TrackID:  "track123",
				Position: 0,
			},
			mockBehavior: func() {
				mockUsecase.EXPECT().CreateJam(
					gomock.Any(),
					&usecaseModel.CreateJamRequest{
						UserID:   "123",
						TrackID:  "track123",
						Position: 0,
					},
				).Return(&usecaseModel.CreateJamResponse{
					RoomID: "room123",
					HostID: "123",
				}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody: &deliveryModel.CreateJamResponse{
				RoomID: "room123",
				HostID: "123",
			},
			hasUserInCtx: true,
		},
		{
			name:           "Unauthorized - No User in Context",
			userID:         0,
			requestBody:    deliveryModel.CreateJamRequest{TrackID: "track123"},
			mockBehavior:   func() {},
			expectedStatus: http.StatusForbidden,
			expectedBody:   nil,
			hasUserInCtx:   false,
		},
		{
			name:   "Usecase Error",
			userID: 123,
			requestBody: deliveryModel.CreateJamRequest{
				TrackID:  "track123",
				Position: 0,
			},
			mockBehavior: func() {
				mockUsecase.EXPECT().CreateJam(
					gomock.Any(),
					gomock.Any(),
				).Return(nil, errors.New("usecase error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   nil,
			hasUserInCtx:   true,
		},
		{
			name:   "Custom Error - Not All Data Provided",
			userID: 123,
			requestBody: deliveryModel.CreateJamRequest{
				TrackID:  "track123",
				Position: 0,
			},
			mockBehavior: func() {
				mockUsecase.EXPECT().CreateJam(
					gomock.Any(),
					gomock.Any(),
				).Return(nil, customErrors.ErrCreateRoomNotAllDataProvided)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
			hasUserInCtx:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/jam/create", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			ctx := context.WithValue(req.Context(), loggerPkg.LoggerKey{}, zap.NewNop().Sugar())
			if tt.hasUserInCtx {
				ctx = context.WithValue(ctx, ctxExtractor.UserContextKey{}, tt.userID)
			}
			req = req.WithContext(ctx)

			handler.CreateRoom(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedStatus == http.StatusCreated && tt.expectedBody != nil {
				var response deliveryModel.APIResponse
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)

				var jamResponse deliveryModel.CreateJamResponse
				respBodyBytes, err := json.Marshal(response.Body)
				assert.NoError(t, err)

				err = json.Unmarshal(respBodyBytes, &jamResponse)
				assert.NoError(t, err)

				expectedResponse := tt.expectedBody.(*deliveryModel.CreateJamResponse)
				assert.Equal(t, expectedResponse.RoomID, jamResponse.RoomID)
				assert.Equal(t, expectedResponse.HostID, jamResponse.HostID)
			}
		})
	}
}

func TestWSHandler_Unauthorized(t *testing.T) {
	_, handler, _ := setupTestJamHandler(t)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/jam/room123/ws", nil)
	req = req.WithContext(context.WithValue(req.Context(), loggerPkg.LoggerKey{}, zap.NewNop().Sugar()))

	handler.WSHandler(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), customErrors.ErrUnauthorized.Error())
}

func TestWSHandler_MissingRoomID(t *testing.T) {
	_, handler, _ := setupTestJamHandler(t)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/jam/ws", nil)

	ctx := context.WithValue(req.Context(), loggerPkg.LoggerKey{}, zap.NewNop().Sugar())
	ctx = context.WithValue(ctx, ctxExtractor.UserContextKey{}, int64(123))
	req = req.WithContext(ctx)

	handler.WSHandler(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), customErrors.ErrRoomIDRequired.Error())
}

func TestWSHandler_JoinJamError(t *testing.T) {
	mockUsecase, handler, _ := setupTestJamHandler(t)

	// Create a test server to handle WebSocket upgrade
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), loggerPkg.LoggerKey{}, zap.NewNop().Sugar())
		ctx = context.WithValue(ctx, ctxExtractor.UserContextKey{}, int64(123))
		r = r.WithContext(ctx)

		// Set up mux vars
		vars := map[string]string{"id": "room123"}
		r = mux.SetURLVars(r, vars)

		mockUsecase.EXPECT().JoinJam(
			gomock.Any(),
			&usecaseModel.JoinJamRequest{
				RoomID: "room123",
				UserID: "123",
			},
		).Return(nil, errors.New("join error"))

		handler.WSHandler(w, r)
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer conn.Close()

	var msg deliveryModel.JamMessage
	err = conn.ReadJSON(&msg)
	assert.NoError(t, err)
	assert.Equal(t, "error", msg.Type)
	assert.Equal(t, "join error", msg.Error)
}

func TestWSHandler_Success(t *testing.T) {
	mockUsecase, handler, _ := setupTestJamHandler(t)

	messageChan := make(chan *usecaseModel.JamMessage, 1)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), loggerPkg.LoggerKey{}, zap.NewNop().Sugar())
		ctx = context.WithValue(ctx, ctxExtractor.UserContextKey{}, int64(123))
		r = r.WithContext(ctx)

		vars := map[string]string{"id": "room123"}
		r = mux.SetURLVars(r, vars)

		mockUsecase.EXPECT().JoinJam(
			gomock.Any(),
			&usecaseModel.JoinJamRequest{
				RoomID: "room123",
				UserID: "123",
			},
		).Return(&usecaseModel.JamMessage{
			Type:   "join",
			UserID: "123",
		}, nil)

		mockUsecase.EXPECT().SubscribeToJamMessages(
			gomock.Any(),
			"room123",
		).Return((<-chan *usecaseModel.JamMessage)(messageChan), nil)

		mockUsecase.EXPECT().HandleClientMessage(
			gomock.Any(),
			"room123",
			"123",
			gomock.Any(),
		).Return(nil).AnyTimes()

		mockUsecase.EXPECT().LeaveJam(
			gomock.Any(),
			"room123",
			"123",
		).Return(nil).AnyTimes()

		handler.WSHandler(w, r)
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer conn.Close()

	var joinMsg deliveryModel.JamMessage
	err = conn.ReadJSON(&joinMsg)
	assert.NoError(t, err)
	assert.Equal(t, "join", joinMsg.Type)
	assert.Equal(t, "123", joinMsg.UserID)

	clientMsg := deliveryModel.JamMessage{
		Type:     "play",
		TrackID:  "track123",
		Position: 100,
	}
	err = conn.WriteJSON(clientMsg)
	assert.NoError(t, err)

	go func() {
		messageChan <- &usecaseModel.JamMessage{
			Type:     "pause",
			TrackID:  "track123",
			Position: 150,
			Paused:   true,
		}
		close(messageChan)
	}()

	var broadcastMsg deliveryModel.JamMessage
	err = conn.ReadJSON(&broadcastMsg)
	assert.NoError(t, err)
	assert.Equal(t, "pause", broadcastMsg.Type)
	assert.Equal(t, "track123", broadcastMsg.TrackID)
	assert.Equal(t, int64(150), broadcastMsg.Position)
	assert.True(t, broadcastMsg.Paused)
}

func TestWSHandler_SubscriptionError(t *testing.T) {
	mockUsecase, handler, _ := setupTestJamHandler(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), loggerPkg.LoggerKey{}, zap.NewNop().Sugar())
		ctx = context.WithValue(ctx, ctxExtractor.UserContextKey{}, int64(123))
		r = r.WithContext(ctx)

		vars := map[string]string{"id": "room123"}
		r = mux.SetURLVars(r, vars)

		mockUsecase.EXPECT().JoinJam(
			gomock.Any(),
			&usecaseModel.JoinJamRequest{
				RoomID: "room123",
				UserID: "123",
			},
		).Return(&usecaseModel.JamMessage{
			Type:   "join",
			UserID: "123",
		}, nil)

		mockUsecase.EXPECT().SubscribeToJamMessages(
			gomock.Any(),
			"room123",
		).Return(nil, errors.New("subscription error"))

		handler.WSHandler(w, r)
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer conn.Close()

	var joinMsg deliveryModel.JamMessage
	err = conn.ReadJSON(&joinMsg)
	assert.NoError(t, err)
	assert.Equal(t, "join", joinMsg.Type)

	_, _, err = conn.ReadMessage()
	assert.Error(t, err)
}
