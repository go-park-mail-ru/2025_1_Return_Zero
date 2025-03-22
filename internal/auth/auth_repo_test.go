package auth

import (
	"testing"
	"time"
)

func TestHashPassword(t *testing.T) {
	password := "securepassword"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if !CheckPasswordHash(password, hash) {
		t.Errorf("Password hash check failed")
	}
}

func TestCreateUser(t *testing.T) {
	repo := NewAuthRepo()
	userData := &RegisterUserData{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	user := repo.CreateUser(userData)
	if user.Username != userData.Username || user.Email != userData.Email {
		t.Errorf("User data does not match")
	}
}

func TestGetUserByUsername(t *testing.T) {
	repo := NewAuthRepo()
	repo.CreateUser(&RegisterUserData{Username: "testuser", Email: "test@example.com", Password: "pass"})

	user := repo.GetUserByUsername("testuser")
	if user == nil || user.Username != "testuser" {
		t.Errorf("Failed to retrieve user by username")
	}
}

func TestGetUserByEmail(t *testing.T) {
	repo := NewAuthRepo()
	repo.CreateUser(&RegisterUserData{Username: "testuser", Email: "test@example.com", Password: "pass"})

	user := repo.GetUserByEmail("test@example.com")
	if user == nil || user.Email != "test@example.com" {
		t.Errorf("Failed to retrieve user by email")
	}
}

func TestCreateSession(t *testing.T) {
	userID := uint(1)
	session, SID := CreateSession(userID)
	if session.UserID != userID {
		t.Errorf("Session UserID mismatch")
	}
	if time.Now().After(session.ExpiresAt) {
		t.Errorf("Session expiration time incorrect")
	}
	if SID == "" {
		t.Errorf("Session ID should not be empty")
	}
}

func TestSessionManagement(t *testing.T) {
	repo := NewAuthRepo()
	userID := uint(1)
	session, SID := CreateSession(userID)

	repo.AppendSession(SID, session)
	if repo.sessions[SID] == nil {
		t.Errorf("Failed to append session")
	}

	repo.DeleteSession(SID)
	if repo.sessions[SID] != nil {
		t.Errorf("Failed to delete session")
	}
}
