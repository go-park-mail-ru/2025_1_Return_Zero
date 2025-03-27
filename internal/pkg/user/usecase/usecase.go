package usecase

import (
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/auth"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/user"
)

func NewUserUsecase(userRepo user.Repository, authRepo auth.Repository) user.Usecase {
	return UserUsecase{
		userRepo: userRepo,
		authRepo: authRepo,
	}
}

type UserUsecase struct {
	userRepo user.Repository
	authRepo auth.Repository
}

func (u UserUsecase) CreateUser(regData *model.RegisterData) (*model.UserToFront, string, error) {
	user, err := u.userRepo.CreateUser(regData)
	if err != nil {
		return nil, "", err
	}
	sessionID := u.authRepo.CreateSession(user.ID)
	return user, sessionID, nil
}

func (u UserUsecase) GetUserBySID(SID string) (*model.UserToFront, error) {
	session, err := u.authRepo.GetSession(SID)
	if err != nil {
		return nil, err
	}
	user, err := u.userRepo.GetUserByID(session.UserID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u UserUsecase) LoginUser(logData *model.LoginData) (*model.UserToFront, string, error) {
	user, err := u.userRepo.LoginUser(logData)
	if err != nil {
		return nil, "", err
	}
	sessionID := u.authRepo.CreateSession(user.ID)
	return user, sessionID, nil
}

func (u UserUsecase) Logout(SID string) {
	u.authRepo.DeleteSession(SID)
}
