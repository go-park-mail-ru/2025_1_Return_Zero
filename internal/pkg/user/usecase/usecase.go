package usecase

import (
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/auth"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/repository"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/user"
)

func NewUserUsecase(userRepo user.Repository, authRepo auth.Repository) user.Usecase {
	return userUsecase{
		userRepo: userRepo,
		authRepo: authRepo,
	}
}

type userUsecase struct {
	userRepo user.Repository
	authRepo auth.Repository
}

func toUsecaseModel(user *repoModel.User) *usecaseModel.User {
	return &usecaseModel.User{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	}
}

func (u userUsecase) CreateUser(user *usecaseModel.User) (*usecaseModel.User, string, error) {
	repoUser := &repoModel.User{
		Username: user.Username,
		Email: user.Email,
		Password: user.Password,
	}
	newUser, err := u.userRepo.CreateUser(repoUser)
	if err != nil {
		return nil, "", err
	}
	sessionID := u.authRepo.CreateSession(newUser.ID)
	return toUsecaseModel(newUser), sessionID, nil
}

func (u userUsecase) GetUserBySID(SID string) (*usecaseModel.User, error) {
	id, err := u.authRepo.GetSession(SID)
	if err != nil {
		return nil, err
	}
	user, err := u.userRepo.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	return toUsecaseModel(user), nil
}

func (u userUsecase) LoginUser(user *usecaseModel.User) (*usecaseModel.User, string, error) {
	repoUser := &repoModel.User{
		Username: user.Username,
		Email: user.Email,
		Password: user.Password,
	}
	loginUser, err := u.userRepo.LoginUser(repoUser)
	if err != nil {
		return nil, "", err
	}
	sessionID := u.authRepo.CreateSession(loginUser.ID)
	return toUsecaseModel(loginUser), sessionID, nil
}

func (u userUsecase) Logout(SID string) {
	u.authRepo.DeleteSession(SID)
}
