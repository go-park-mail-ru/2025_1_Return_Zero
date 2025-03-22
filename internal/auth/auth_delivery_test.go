package auth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	model "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/models"
	"time"
)

func TestSignup(t *testing.T) {
	handler := NewAuthHandler()
	body := `{"username": "testuser", "email": "test@example.com", "password": "password123"}`
	req := httptest.NewRequest("POST", "/signup", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Signup(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK, got %v", resp.StatusCode)
	}
}

func TestLogin(t *testing.T) {
	handler := NewAuthHandler()
	password, _ := HashPassword("password123")
	handler.uc.repo.CreateUser(&RegisterUserData{Username: "testuser123", Email: "testuser123@gmail.com", Password: password})
	body := `{"username": "testuser123", "password": "password123"}`
	req := httptest.NewRequest("POST", "/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Login(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK, got %v", resp.StatusCode)
	}
}

func TestLogout(t *testing.T) {
	handler := NewAuthHandler()
	req := httptest.NewRequest("POST", "/logout", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "testsession"})
	w := httptest.NewRecorder()

	handler.Logout(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK, got %v", resp.StatusCode)
	}
}

func TestCheckUser(t *testing.T) {
	handler := NewAuthHandler()
	handler.uc.repo.AppendSession("testsession", &model.Session{UserID: 1, ExpiresAt: time.Now().Add(time.Hour)})
	req := httptest.NewRequest("GET", "/check", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "testsession"})
	w := httptest.NewRecorder()

	handler.CheckUser(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK, got %v", resp.StatusCode)
	}
}
