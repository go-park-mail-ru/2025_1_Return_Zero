package artist

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/config"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/ctxExtractor"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/customErrors"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/errorStatus"
	jsonPkg "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/json"
	loggerPkg "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/jam"
	model "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model"
	delivery "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/delivery"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type JamHandler struct {
	usecase jam.Usecase
	cfg     *config.Config
}

func NewJamHandler(usecase jam.Usecase, cfg *config.Config) *JamHandler {
	return &JamHandler{usecase: usecase, cfg: cfg}
}

func (h *JamHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := loggerPkg.LoggerFromContext(ctx)
	userID, ok := ctxExtractor.UserFromContext(ctx)
	if !ok {
		status := errorStatus.ErrorStatus(customErrors.ErrUnauthorized)
		logger.Error("user not found in context", zap.Error(customErrors.ErrUnauthorized))
		jsonPkg.WriteErrorResponse(w, status, customErrors.ErrUnauthorized.Error(), nil)
		return
	}

	var request delivery.CreateJamRequest
	err := jsonPkg.ReadJSON(w, r, &request)
	if err != nil {
		status := errorStatus.ErrorStatus(err)
		logger.Error("failed to read json", zap.Error(err))
		jsonPkg.WriteErrorResponse(w, status, err.Error(), nil)
		return
	}

	valid, err := govalidator.ValidateStruct(request)
	if err != nil {
		status := errorStatus.ErrorStatus(err)
		logger.Error("failed to validate struct", zap.Error(err))
		jsonPkg.WriteErrorResponse(w, status, err.Error(), nil)
		return
	}
	if !valid {
		status := errorStatus.ErrorStatus(customErrors.ErrCreateRoomNotAllDataProvided)
		logger.Error("failed to validate struct", zap.Error(customErrors.ErrCreateRoomNotAllDataProvided))
		jsonPkg.WriteErrorResponse(w, status, customErrors.ErrCreateRoomNotAllDataProvided.Error(), nil)
		return
	}

	userIDStr := strconv.FormatInt(userID, 10)

	usecaseRequest := model.CreateJamRequestFromDeliveryToUsecase(&request, userIDStr)
	usecaseResponse, err := h.usecase.CreateJam(ctx, usecaseRequest)
	if err != nil {
		status := errorStatus.ErrorStatus(err)
		logger.Error("failed to create jam", zap.Error(err))
		jsonPkg.WriteErrorResponse(w, status, err.Error(), nil)
		return
	}

	response := model.CreateJamResponseFromUsecaseToDelivery(usecaseResponse)
	jsonPkg.WriteSuccessResponse(w, http.StatusCreated, response, nil)
}

func (h *JamHandler) WSHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := loggerPkg.LoggerFromContext(ctx)

	userID, ok := ctxExtractor.UserFromContext(ctx)
	if !ok {
		logger.Error("user not found in context", zap.Error(customErrors.ErrUnauthorized))
		http.Error(w, customErrors.ErrUnauthorized.Error(), http.StatusUnauthorized)
		return
	}

	roomID := mux.Vars(r)["id"]
	if roomID == "" {
		logger.Error("room id is required", zap.Error(customErrors.ErrRoomIDRequired))
		http.Error(w, customErrors.ErrRoomIDRequired.Error(), http.StatusBadRequest)
		return
	}

	upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		return true
	}}

	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("failed to upgrade to websocket", zap.Error(err))
		return
	}

	userIDStr := strconv.FormatInt(userID, 10)

	usecaseRequest := usecase.JoinJamRequest{
		RoomID: roomID,
		UserID: userIDStr,
	}
	usecaseResponse, err := h.usecase.JoinJam(ctx, &usecaseRequest)
	if err != nil {
		logger.Error("failed to join jam", zap.Error(err))
		err = wsConn.WriteJSON(delivery.JamMessage{
			Type:  "error",
			Error: err.Error(),
		})
		if err != nil {
			logger.Error("failed to write error message to websocket", zap.Error(err))
		}
		err = wsConn.Close()
		if err != nil {
			logger.Error("failed to close websocket", zap.Error(err))
		}
		return
	}

	response := model.JamMessageFromUsecaseToDelivery(usecaseResponse)
	err = wsConn.WriteJSON(response)
	if err != nil {
		logger.Error("failed to write response to websocket", zap.Error(err))
	}

	messageChan, err := h.usecase.SubscribeToJamMessages(ctx, roomID)
	if err != nil {
		logger.Error("failed to subscribe to jam messages", zap.Error(err))
		err = wsConn.Close()
		if err != nil {
			logger.Error("failed to close websocket", zap.Error(err))
		}
		return
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case usecaseMessage, ok := <-messageChan:
				if !ok {
					return
				}
				deliveryMessage := model.JamMessageFromUsecaseToDelivery(usecaseMessage)
				err := wsConn.WriteJSON(deliveryMessage)
				if err != nil {
					logger.Error("failed to write message to websocket", zap.Error(err))
					return
				}
			}
		}
	}()

	for {
		_, data, err := wsConn.ReadMessage()
		if err != nil {
			err := h.usecase.LeaveJam(ctx, roomID, userIDStr)
			if err != nil {
				logger.Error("failed to leave jam", zap.Error(err))
			}
			return
		}
		var m delivery.JamMessage
		err = json.Unmarshal(data, &m)
		if err != nil {
			logger.Error("failed to unmarshal message", zap.Error(err))
			return
		}

		usecaseMessage := model.JamMessageFromDeliveryToUsecase(&m)

		err = h.usecase.HandleClientMessage(ctx, roomID, userIDStr, usecaseMessage)
		if err != nil {
			logger.Error("failed to handle client message", zap.Error(err))
			return
		}
	}
}
