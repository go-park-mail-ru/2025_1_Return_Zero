package auth

import (
	"errors"
	"sync"
	"time"

	validator "github.com/asaskevich/govalidator"
	model "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/models"
)

type AuthUserCase struct {
	mu   sync.RWMutex
	repo *AuthRepo
}

func NewAuthUserCase(repo *AuthRepo) *AuthUserCase {
	return &AuthUserCase{
		mu:   sync.RWMutex{},
		repo: repo,
	}
}

func (uc *AuthUserCase) CleanupSessions() {
	for {
		time.Sleep(time.Hour)
		uc.mu.Lock()
		now := time.Now()
		for sid, session := range uc.repo.sessions {
			if now.After(session.ExpiresAt) {
				delete(uc.repo.sessions, sid)
			}
		}
		uc.mu.Unlock()
	}
}

func (uc *AuthUserCase) SignupUser(u *RegisterUserData) (*model.User, error) {
	if _, err := validator.ValidateStruct(u); err != nil {
		return nil, err
	}

	uc.mu.RLock()
	user := uc.repo.GetUserByEmail(u.Email)
	uc.mu.RUnlock()
	if user != nil {
		return nil, errors.New("user with this email already exists")
	}

	uc.mu.RLock()
	user = uc.repo.GetUserByUsername(u.Username)
	uc.mu.RUnlock()
	if user != nil {
		return nil, errors.New("user with this username already exists")
	}

	hashedPassword, err := HashPassword(u.Password)
	if err != nil {
		return nil, err
	}

	u.Password = hashedPassword
	uc.mu.Lock()
	user = uc.repo.CreateUser(u)
	uc.mu.Unlock()
	if user == nil {
		return nil, errors.New("failed to create user")
	}

	response := &model.User{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	return response, nil
}

func (uc *AuthUserCase) LoginUser(u *LoginUserData) (*model.User, error) {
	uc.mu.RLock()
	user := uc.repo.GetUserByUsername(u.Username)

	if user == nil {
		user = uc.repo.GetUserByEmail(u.Email)
		if user == nil {
			return nil, errors.New("user not found")
		}
	}
	uc.mu.RUnlock()

	password := u.Password

	if !CheckPasswordHash(password, user.Password) {
		return nil, errors.New("invalid password")
	}

	response := &model.User{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	return response, nil
}
