package user

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"

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

func userDeleteToUsecaseModel(user *deliveryModel.UserDelete) *usecaseModel.User {
	return &usecaseModel.User{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	}
}

func changeDataToUsecaseModel(changeData *deliveryModel.ChangeUserData) *usecaseModel.ChangeUserData {
	return &usecaseModel.ChangeUserData{
		Username:    changeData.Username,
		Email:       changeData.Email,
		Password:    changeData.Password,
		NewUsername: changeData.NewUsername,
		NewEmail:    changeData.NewEmail,
		NewPassword: changeData.NewPassword,
	}
}

func validateData(data interface{}) (bool, error) {
	result, err := govalidator.ValidateStruct(data)
	if err != nil {
		return false, err
	}
	return result, nil
}

func createCookie(name string, value string, expiration time.Time, path string) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Expires:  expiration,
		HttpOnly: true,
		Path: 	  path,
	}
} 

// Signup godoc
// @Summary Register a new user
// @Description Creates a new user account with provided registration data
// @Tags auth
// @Accept json
// @Produce json
// @Param register body delivery.RegisterData true "User registration data"
// @Success 200 {object} delivery.APIResponse{body=delivery.UserToFront} "User successfully registered"
// @Failure 400 {object} delivery.APIBadRequestErrorResponse "Bad request - invalid registration data"
// @Router /auth/signup [post]
func (h *UserHandler) Signup(w http.ResponseWriter, r *http.Request) {
	regData := &deliveryModel.RegisterData{}
	err := helpers.ReadJSON(w, r, regData)
	if err != nil {
		helpers.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	if regData.Username == "default_avatar" {
		helpers.WriteErrorResponse(w, http.StatusBadRequest, "Wrong username", nil)
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
	cookie := createCookie("session_id", sessionId, time.Now().Add(24*time.Hour), "/")
	http.SetCookie(w, cookie)
	helpers.WriteSuccessResponse(w, http.StatusOK, toUserToFront(user), nil)
}

// Login godoc
// @Summary Authenticate user
// @Description Authenticates a user with provided login credentials and returns a session
// @Tags auth
// @Accept json
// @Produce json
// @Param login body delivery.LoginData true "User login data"
// @Success 200 {object} delivery.APIResponse{body=delivery.UserToFront} "User successfully authenticated"
// @Failure 400 {object} delivery.APIBadRequestErrorResponse "Bad request - invalid login data"
// @Router /auth/login [post]
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
	cookie := createCookie("session_id", sessionId, time.Now().Add(24*time.Hour), "/")
	http.SetCookie(w, cookie)
	helpers.WriteSuccessResponse(w, http.StatusOK, toUserToFront(user), nil)
}

// Logout godoc
// @Summary Log out user
// @Description Terminates user session and invalidates session cookie
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} delivery.APIResponse{body=delivery.Message} "Successfully logged out"
// @Failure 400 {object} delivery.APIBadRequestErrorResponse "Bad request - session not found"
// @Router /auth/logout [post]
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

	msg := &deliveryModel.Message{
		Message: "Successfully logged out",
	}

	helpers.WriteSuccessResponse(w, http.StatusOK, msg, nil)
}

// CheckUser godoc
// @Summary Check user authentication
// @Description Verifies user's session and returns user information if authenticated
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} delivery.APIResponse{body=delivery.UserToFront} "User information"
// @Failure 400 {object} delivery.APIBadRequestErrorResponse "Bad request - session not found or invalid"
// @Router /auth/check [get]
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

// GetUserAvatar godoc
// @Summary Get user avatar
// @Description Retrieves the avatar URL for a specific user
// @Tags user
// @Accept json
// @Produce json
// @Param username path string true "Username"
// @Success 200 {object} delivery.APIResponse{body=delivery.AvatarData} "Avatar URL"
// @Failure 400 {object} delivery.APIBadRequestErrorResponse "Bad request - username not found"
// @Failure 500 {object} delivery.APIInternalServerErrorResponse "Internal server error"
// @Router /user/{username}/avatar [get]
func (h *UserHandler) GetUserAvatar(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())
	vars := mux.Vars(r)
	username, ok := vars["username"]

	if !ok {
		logger.Error("username not found in URL")
		helpers.WriteErrorResponse(w, http.StatusBadRequest, "username not found in URL", nil)
		return
	}

	presignedUrl, err := h.usecase.GetAvatar(username)
	if err != nil {
		logger.Error("failed to get avatar", zap.Error(err))
		helpers.WriteErrorResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	avatar := &deliveryModel.AvatarData{
		Avatar: presignedUrl,
	}
	helpers.WriteSuccessResponse(w, http.StatusOK, avatar, nil)
}

// UploadAvatar godoc
// @Summary Upload user avatar
// @Description Uploads a new avatar image for a specific user
// @Tags user
// @Accept multipart/form-data
// @Produce json
// @Param username path string true "Username"
// @Param avatar formData file true "Avatar image file (max 5MB, image formats only)"
// @Success 200 {object} delivery.APIResponse{body=delivery.Message} "Avatar successfully uploaded"
// @Failure 400 {object} delivery.APIBadRequestErrorResponse "Bad request - invalid file or username"
// @Router /user/{username}/avatar [post]
func (h *UserHandler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())
	vars := mux.Vars(r)
	username, ok := vars["username"]
	if !ok {
		logger.Error("username not found in URL")
		helpers.WriteErrorResponse(w, http.StatusBadRequest, "username not found in URL", nil)
		return
	}

	const maxUploadSize = 5 << 20
	err := r.ParseMultipartForm(maxUploadSize)
	if err != nil {
		logger.Error("failed to parse form", zap.Error(err))
		helpers.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	file, fileHeader, err := r.FormFile("avatar")
	if err != nil {
		logger.Error("failed to get file from form", zap.Error(err))
		helpers.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	defer file.Close()

	if fileHeader.Size > maxUploadSize {
		logger.Error("file size exceeds limit", zap.Error(err))
		helpers.WriteErrorResponse(w, http.StatusBadRequest, "file size exceeds limit", nil)
		return
	}

	contentType := fileHeader.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		logger.Error("invalid file type", zap.String("contentType", contentType))
		helpers.WriteErrorResponse(w, http.StatusBadRequest, "only image files are allowed", nil)
		return
	}

	err = h.usecase.UploadAvatar(username, file)
	if err != nil {
		logger.Error("failed to upload avatar", zap.Error(err))
		helpers.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	msg := &deliveryModel.Message{
		Message: "Avatar successfully uploaded",
	}
	helpers.WriteSuccessResponse(w, http.StatusOK, msg, nil)
}

// ChangeUserData godoc
// @Summary Change user profile data
// @Description Updates user's profile information such as username, email, or password
// @Tags user
// @Accept json
// @Produce json
// @Param data body delivery.ChangeUserData true "User data to be updated"
// @Success 200 {object} delivery.APIResponse{body=delivery.UserToFront} "User data successfully updated"
// @Failure 400 {object} delivery.APIBadRequestErrorResponse "Bad request - invalid user data or validation failure"
// @Router /user/{username} [put]
func (h *UserHandler) ChangeUserData(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())

	userAuth, exist := middleware.GetUserFromContext(r.Context())
	if !exist {
		logger.Error("user not auth")
		helpers.WriteErrorResponse(w, http.StatusBadRequest, "user not found in context", nil)
		return
	}

	changeData := &deliveryModel.ChangeUserData{}
	err := helpers.ReadJSON(w, r, changeData)
	if err != nil {
		logger.Error("failed to read change user data", zap.Error(err))
		helpers.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	isValid, err := validateData(changeData)
	if err != nil || !isValid {
		logger.Error("failed to validate change user data", zap.Error(err))
		helpers.WriteErrorResponse(w, http.StatusBadRequest, ErrValidationFailed.Error(), nil)
		return
	}
	changeDataUsecase := changeDataToUsecaseModel(changeData)
	newUser, err := h.usecase.ChangeUserData(userAuth.Username, changeDataUsecase)
	if err != nil {
		logger.Error("failed to change user data", zap.Error(err))
		helpers.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.WriteSuccessResponse(w, http.StatusOK, toUserToFront(newUser), nil)
}

// DeleteUser godoc
// @Summary Delete user account
// @Description Deletes the authenticated user's account. Requires valid session cookie and matching user credentials.
// @Tags user
// @Accept json
// @Produce json
// @Param Authorization header string true "Session ID cookie (session_id=...)"
// @Param user body delivery.UserDelete true "User credentials for deletion verification"
// @Success 200 {object} delivery.APIResponse{body=delivery.Message} "User successfully deleted"
// @Failure 400 {object} delivery.APIBadRequestErrorResponse "Possible errors: invalid request body, validation failed, credentials mismatch, session cookie missing"
// @Failure 500 {object} delivery.APIBadRequestErrorResponse "Internal server error during user deletion"
// @Router /user/{username} [delete]
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())

	userAuth, exist := middleware.GetUserFromContext(r.Context())
	if !exist {
		logger.Error("user not auth")
		helpers.WriteErrorResponse(w, http.StatusBadRequest, "user not found in context", nil)
		return
	}

	userDelete := &deliveryModel.UserDelete{}
	err := helpers.ReadJSON(w, r, userDelete)
	if err != nil {
		logger.Error("failed to read change user data", zap.Error(err))
		helpers.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	isValid, err := validateData(userDelete)
	if err != nil || !isValid {
		logger.Error("failed to validate change user data", zap.Error(err))
		helpers.WriteErrorResponse(w, http.StatusBadRequest, ErrValidationFailed.Error(), nil)
		return
	}

	if userAuth.Username != userDelete.Username || userAuth.Email != userDelete.Email {
		logger.Error("wrong user")
		helpers.WriteErrorResponse(w, http.StatusBadRequest, "user not found in context", nil)
		return
	}

	cookie, err := r.Cookie("session_id")
	if err != nil {
		logger.Error("failed to get session id", zap.Error(err))
		helpers.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	sessionId := cookie.Value

	usecaseUser := userDeleteToUsecaseModel(userDelete)
	err = h.usecase.DeleteUser(usecaseUser, sessionId)
	if err != nil {
		logger.Error("failed to delete user", zap.Error(err))
		helpers.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	cookie.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, cookie)
	msg := &deliveryModel.Message{
		Message: "User successfully deleted",
	}
	helpers.WriteSuccessResponse(w, http.StatusOK, msg, nil)
}