package user

import (
	"errors"
	"net/http"
	"time"

	"github.com/asaskevich/govalidator"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/middleware"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers"
	deliveryModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/delivery"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/user"
	"go.uber.org/zap"
)

var (
	ErrValidationFailed = errors.New("validation failed")
)

type UserHandler struct {
	usecase user.Usecase
}

func NewUserHandler(usecase user.Usecase) *UserHandler {
	return &UserHandler{
		usecase: usecase,
	}
}

func registerToUsecaseModel(user *deliveryModel.RegisterData) *usecaseModel.User {
	return &usecaseModel.User{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	}
}

func loginToUsecaseModel(user *deliveryModel.LoginData) *usecaseModel.User {
	return &usecaseModel.User{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	}
}

func toUserToFront(user *usecaseModel.User) *deliveryModel.UserToFront {
	return &deliveryModel.UserToFront{
		Username: user.Username,
		Email:    user.Email,
	}
}

func validateData(data interface{}) (bool, error) {
	result, err := govalidator.ValidateStruct(data)
	if err != nil {
		return false, err
	}
	return result, nil
}

func (h *UserHandler) Signup(w http.ResponseWriter, r *http.Request) {
	regData := &deliveryModel.RegisterData{}
	err := helpers.ReadJSON(w, r, regData)
	if err != nil {
		helpers.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	isValid, err := validateData(regData)
	if err != nil || !isValid {
		helpers.WriteErrorResponse(w, http.StatusBadRequest, ErrValidationFailed.Error(), nil)
		return
	}
	user, sessionId, err := h.usecase.CreateUser(registerToUsecaseModel(regData))
	if err != nil {
		helpers.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionId,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	helpers.WriteSuccessResponse(w, http.StatusOK, toUserToFront(user), nil)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())

	logData := &deliveryModel.LoginData{}
	err := helpers.ReadJSON(w, r, logData)
	if err != nil {
		logger.Error("failed to read login data", zap.Error(err))
		helpers.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	isValid, err := validateData(logData)
	if err != nil || !isValid {
		logger.Error("failed to validate login data", zap.Error(err))
		helpers.WriteErrorResponse(w, http.StatusBadRequest, ErrValidationFailed.Error(), nil)
		return
	}
	user, sessionId, err := h.usecase.LoginUser(loginToUsecaseModel(logData))
	if err != nil {
		logger.Error("failed to login user", zap.Error(err))
		helpers.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionId,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	helpers.WriteSuccessResponse(w, http.StatusOK, toUserToFront(user), nil)
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())
	cookie, err := r.Cookie("session_id")
	if err != nil {
		logger.Error("failed to get session id", zap.Error(err))
		helpers.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	sessionId := cookie.Value
	h.usecase.Logout(sessionId)
	cookie.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, cookie)
	helpers.WriteSuccessResponse(w, http.StatusOK, "Succesfully logged out", nil)
}

func (h *UserHandler) CheckUser(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())
	cookie, err := r.Cookie("session_id")
	if err != nil {
		logger.Error("failed to get session id", zap.Error(err))
		helpers.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	sessionId := cookie.Value
	user, err := h.usecase.GetUserBySID(sessionId)
	if err != nil {
		logger.Error("failed to get user by session id", zap.Error(err))
		helpers.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.WriteSuccessResponse(w, http.StatusOK, toUserToFront(user), nil)
}
