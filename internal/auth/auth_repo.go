package auth

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	model "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/models"
	"golang.org/x/crypto/bcrypt"
)

var (
	USERS_CNT = 0
)

type AuthRepo struct {
	users    []*model.User
	sessions map[string]*model.Session
}

func NewAuthRepo() *AuthRepo {
	return &AuthRepo{
		users:    make([]*model.User, 0),
		sessions: make(map[string]*model.Session),
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (auth *AuthRepo) CreateUser(u *RegisterUserData) *model.User {
	user := &model.User{
		ID:       uint(USERS_CNT),
		Username: u.Username,
		Email:    u.Email,
		Password: u.Password,
	}
	USERS_CNT++
	auth.users = append(auth.users, user)
	return user
}

func (auth *AuthRepo) GetUserByUsername(username string) *model.User {
	for _, user := range auth.users {
		if user.Username == username {
			return user
		}
	}
	return nil
}

func (auth *AuthRepo) GetUserByEmail(email string) *model.User {
	for _, user := range auth.users {
		if user.Email == email {
			return user
		}
	}
	return nil
}

func generateSessionID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func (auth *AuthRepo) AppendSession(SID string, s *model.Session) {
	auth.sessions[SID] = s
}

func CreateSession(ID uint) (*model.Session, string) {
	SID := generateSessionID()
	expiration := time.Now().Add(24 * time.Hour)
	session := &model.Session{
		UserID:    ID,
		ExpiresAt: expiration,
	}
	return session, SID
}

func (auth *AuthRepo) DeleteSession(SID string) {
	delete(auth.sessions, SID)
}
