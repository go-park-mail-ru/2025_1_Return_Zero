package user

import (
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/user"
)

type UserHandler struct {
	usecase user.Usecase
}

func NewUserHandler(usecase user.Usecase) *UserHandler {
	return &UserHandler{
		usecase: usecase,
	}
}

func (h *UserHandler) Signup(w http.ResponseWriter, r *http.Request) {
	regData := &model.RegisterData{}
	err := helpers.ReadJSON(w, r, regData)
	if err != nil {
		helpers.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	user, sessionId, err := h.usecase.CreateUser(regData)
	if err != nil {
		helpers.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	cookie := &http.Cookie{
		Name:  "session_id",
		Value: sessionId,
		Expires: time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	err = helpers.WriteJSON(w, http.StatusOK, user, nil)
	if err != nil {
		helpers.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	logData := &model.LoginData{}
	err := helpers.ReadJSON(w, r, logData)
	if err != nil {
		helpers.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	user, sessionId, err := h.usecase.LoginUser(logData)
	if err != nil {
		helpers.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	cookie := &http.Cookie{
		Name:  "session_id",
		Value: sessionId,
		Expires: time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	err = helpers.WriteJSON(w, http.StatusOK, user, nil)
	if err != nil {
		helpers.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		helpers.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	sessionId := cookie.Value
	h.usecase.Logout(sessionId)
	cookie.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, cookie)
	err = helpers.WriteJSON(w, http.StatusOK, "Succesfully logged out", nil)
	if err != nil {
		helpers.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *UserHandler) CheckUser(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		helpers.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	sessionId := cookie.Value
	user, err := h.usecase.GetUserBySID(sessionId)
	if err != nil {
		helpers.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	err = helpers.WriteJSON(w, http.StatusOK, user, nil)
	if err != nil {
		helpers.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
