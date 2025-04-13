package usecase

import (
	"fmt"
	"io"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/auth"
	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/repository"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/user"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/userAvatarFile"
)

func NewUserUsecase(userRepo user.Repository, authRepo auth.Repository, userFileRepo userAvatarFile.Repository) user.Usecase {
	return userUsecase{
		userRepo: userRepo,
		authRepo: authRepo,
		userFileRepo: userFileRepo,
	}
}

type userUsecase struct {
	userRepo user.Repository
	authRepo auth.Repository
	userFileRepo userAvatarFile.Repository
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

	userUsecase := toUsecaseModel(newUser)
	sessionID := u.authRepo.CreateSession(newUser.ID)
	return userUsecase, sessionID, nil
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

	usecaseUser := toUsecaseModel(user)

	return usecaseUser, nil
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
	usecaseUser := toUsecaseModel(loginUser)

	sessionID := u.authRepo.CreateSession(loginUser.ID)
	return usecaseUser, sessionID, nil
}

func (u userUsecase) Logout(SID string) {
	u.authRepo.DeleteSession(SID)
}

func (u userUsecase) GetAvatar(username string) (string, error) {
	avatarUrl, err := u.userRepo.GetAvatar(username)
	if err != nil {
		return "", err
	}
	presignedUrl, err := u.userFileRepo.GetAvatarURL(avatarUrl)
	if err != nil {
		return "", err
	}

	return presignedUrl, nil
}

func (u userUsecase) UploadAvatar(username string, fileAvatar io.Reader) error {
	fileURL, err := u.userFileRepo.UploadUserAvatar(username, fileAvatar)
	if err != nil {
		return err
	}
	fmt.Println(fileURL)
	err = u.userRepo.UploadAvatar(fileURL, username)
	if err != nil {
		return err
	}
	return nil
}