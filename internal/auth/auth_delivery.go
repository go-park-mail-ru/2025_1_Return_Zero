package auth

import (
	"net/http"
	"time"

	helper "github.com/go-park-mail-ru/2025_1_Return_Zero/pkg/JSONHelpers"
)

type AuthHandler struct {
	uc *AuthUserCase
}

func NewAuthHandler() *AuthHandler {
	repo := NewAuthRepo()
	usecase := NewAuthUserCase(repo)
	go usecase.CleanupSessions()
	return &AuthHandler{
		uc: usecase,
	}
}

func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var u RegisterUserData
	if err := helper.ReadJSON(w, r, &u); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	response, err := h.uc.SignupUser(&u)
	if err != nil {
		helper.WriteJSON(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	
	h.uc.mu.Lock()
	session, SID := CreateSession(response.ID)
	h.uc.mu.Unlock()
	cookie := http.Cookie{
		Name:     "session_id",
		Value:    SID,
		Expires:  session.ExpiresAt,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	h.uc.repo.AppendSession(SID, session)

	helper.WriteJSON(w, http.StatusOK, response, nil)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var u LoginUserData
	if err := helper.ReadJSON(w, r, &u); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	response, err := h.uc.LoginUser(&u)
	if err != nil {
		helper.WriteJSON(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	h.uc.mu.Lock()
	session, SID := CreateSession(response.ID)
	h.uc.mu.Unlock()
	cookie := http.Cookie{
		Name:     "session_id",
		Value:    SID,
		Expires:  session.ExpiresAt,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	h.uc.repo.AppendSession(SID, session)

	helper.WriteJSON(w, http.StatusOK, response, nil)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "No session", http.StatusBadRequest)
		return
	}

	SID := cookie.Value
	h.uc.mu.Lock()
	h.uc.repo.DeleteSession(SID)
	h.uc.mu.Unlock()

	cookie.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, cookie)

	helper.WriteJSON(w, http.StatusOK, "Logged out", nil)
}

func (h *AuthHandler) CheckUser(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "No session", http.StatusBadRequest)
		return
	}

	SID := cookie.Value
	h.uc.mu.RLock()
	_, ok := h.uc.repo.sessions[SID]
	if !ok {
		http.Error(w, "No session", http.StatusBadRequest)
		return
	}
	h.uc.mu.RUnlock()

	helper.WriteJSON(w, http.StatusOK, "User is authorized", nil)
}