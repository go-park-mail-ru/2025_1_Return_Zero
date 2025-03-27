package repository

import (
	"errors"
	"sync"

	"github.com/asaskevich/govalidator"
	"golang.org/x/crypto/bcrypt"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/user"
)

var (
	ErrUsernameExist    = errors.New("user with this username already exists")
	ErrEmailExist       = errors.New("user with this email already exists")
	ErrUserNotFound     = errors.New("user not found")
	ErrValidationFailed = errors.New("validation failed")
)

type UserMemoryRepository struct {
	mu    sync.RWMutex
	users map[string]*model.User
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func NewUserMemoryRepository() user.Repository {
	repo := &UserMemoryRepository{
		users: make(map[string]*model.User),
	}

	return repo
}

func userToFront(user *model.User) *model.UserToFront {
	return &model.UserToFront{
		ID:       user.ID,
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

func (r *UserMemoryRepository) CreateUser(regData *model.RegisterData) (*model.UserToFront, error) {
	result, err := validateData(regData)
	if err != nil {
		return nil, err
	} else if !result {
		return nil, ErrValidationFailed
	}
	_, isExist := r.users[regData.Username]
	if isExist {
		return nil, ErrUsernameExist
	}

	for _, user := range r.users {
		if user.Email == regData.Email {
			return nil, ErrEmailExist
		}
	}

	hashedPassword, err := HashPassword(regData.Password)
	if err != nil {
		return nil, err
	}

	newUser := &model.User{
		ID:       uint(len(r.users) + 1),
		Username: regData.Username,
		Password: hashedPassword,
		Email:    regData.Email,
	}

	r.mu.Lock()
	r.users[newUser.Username] = newUser
	r.mu.Unlock()
	return userToFront(newUser), nil
}

func (r *UserMemoryRepository) GetUserByID(ID uint) (*model.UserToFront, error) {
	r.mu.RLock()
	for _, user := range r.users {
		if user.ID == uint(ID) {
			r.mu.RUnlock()
			return userToFront(user), nil
		}
	}
	r.mu.RUnlock()
	return nil, ErrUserNotFound
}

func (r *UserMemoryRepository) LoginUser(logData *model.LoginData) (*model.UserToFront, error) {
	result, err := validateData(logData)
	if err != nil {
		return nil, err
	} else if !result {
		return nil, ErrValidationFailed
	}

	var user *model.User
	var isExist bool
	r.mu.RLock()
	if logData.Username == "" {
		for _, u := range r.users {
			if u.Email == logData.Email {
				user = u
				isExist = true
				break
			}
		}
	} else {
		user, isExist = r.users[logData.Username]
	}
	r.mu.RUnlock()
	if !isExist || !CheckPasswordHash(logData.Password, user.Password) {
		return nil, ErrUserNotFound
	}
	return userToFront(user), nil
}
