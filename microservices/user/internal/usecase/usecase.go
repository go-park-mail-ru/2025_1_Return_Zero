package usecase

import (
	"context"

	domain "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/user/internal/domain"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/user/model/usecase"
	model "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/user/model"
)

func NewUserUsecase(userRepository domain.Repository) domain.Usecase {
	return &userUsecase{
		userRepo: userRepository,
	}
}

type userUsecase struct {
	userRepo domain.Repository
}

func (u *userUsecase) CreateUser(ctx context.Context, registerData *usecaseModel.RegisterData) (*usecaseModel.UserFront, error) {
	repoData := model.RegisterDataFromUsecaseToRepository(registerData)
	userRepoData, err := u.userRepo.CreateUser(ctx, repoData)
	if err != nil {
		return nil, err
	}
	return model.UserFromRepositoryToUsecase(userRepoData), nil
}

func (u *userUsecase) LoginUser(ctx context.Context, loginData *usecaseModel.LoginData) (*usecaseModel.UserFront, error) {
	repoData := model.LoginDataFromUsecaseToRepository(loginData)
	userRepoData, err := u.userRepo.LoginUser(ctx, repoData)
	if err != nil {
		return nil, err
	}
	return model.UserFromRepositoryToUsecase(userRepoData), nil
}

func (u *userUsecase) UploadAvatar(ctx context.Context, avatarUrl string, id int64) error {
	return u.userRepo.UploadAvatar(ctx, avatarUrl, id)
}

func (u *userUsecase) DeleteUser(ctx context.Context, deleteData *usecaseModel.UserDelete) error {
	repoData := model.UserDeleteFromUsecaseToRepository(deleteData)
	return u.userRepo.DeleteUser(ctx, repoData)
}

func (u *userUsecase) ChangeUserData(ctx context.Context, username string, changeData *usecaseModel.ChangeUserData) error {
	userData := model.ChangeUserDataFromUsecaseToRepository(changeData)
	return u.userRepo.ChangeUserData(ctx, username, userData)
}

func (u *userUsecase) ChangeUserPrivacySettings(ctx context.Context, username string, privacySettings *usecaseModel.PrivacySettings) error {
	repoData := model.PrivacySettingsFromUsecaseToRepository(privacySettings)
	return u.userRepo.ChangeUserPrivacySettings(ctx, username, repoData)
}

func (u *userUsecase) GetFullUserData(ctx context.Context, username string) (*usecaseModel.UserFullData, error) {
	userRepoData, err := u.userRepo.GetFullUserData(ctx, username)
	if err != nil {
		return nil, err
	}
	return model.UserFullDataFromRepositoryToUsecase(userRepoData), nil
}

func (u *userUsecase) GetUserByID(ctx context.Context, id int64) (*usecaseModel.UserFront, error) {
	userRepoData, err := u.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return model.UserFromRepositoryToUsecase(userRepoData), nil
}