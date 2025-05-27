package usecase

import (
	"context"

	domain "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/user/internal/domain"
	model "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/user/model"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/user/model/usecase"
)

func NewUserUsecase(userRepository domain.Repository, s3Repository domain.S3Repository) domain.Usecase {
	return &userUsecase{
		userRepo: userRepository,
		s3Repo:   s3Repository,
	}
}

type userUsecase struct {
	userRepo domain.Repository
	s3Repo   domain.S3Repository
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

func (u *userUsecase) GetIDByUsername(ctx context.Context, username string) (int64, error) {
	id, err := u.userRepo.GetIDByUsername(ctx, username)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (u *userUsecase) GetUserPrivacySettings(ctx context.Context, id int64) (*usecaseModel.PrivacySettings, error) {
	privacySettingsRepoData, err := u.userRepo.GetUserPrivacy(ctx, id)
	if err != nil {
		return nil, err
	}
	return model.PrivacyFromRepositoryToUsecase(privacySettingsRepoData), nil
}

func (u *userUsecase) GetAvatarURL(ctx context.Context, fileKey string) (string, error) {
	avatarURL, err := u.s3Repo.GetAvatarURL(ctx, fileKey)
	if err != nil {
		return "", err
	}
	return avatarURL, nil
}

func (u *userUsecase) UploadUserAvatar(ctx context.Context, username string, file []byte) (string, error) {
	avatarURL, err := u.s3Repo.UploadUserAvatar(ctx, username, file)
	if err != nil {
		return "", err
	}
	return avatarURL, nil
}

func (u *userUsecase) GetLabelIDByUserID(ctx context.Context, userID int64) (int64, error) {
	labelID, err := u.userRepo.GetLabelIDByUserID(ctx, userID)
	if err != nil {
		return 0, err
	}
	return labelID, nil
}

func (u *userUsecase) CheckUsersByUsernames(ctx context.Context, usernames []string) error {
	err := u.userRepo.CheckUsersByUsernames(ctx, usernames)
	if err != nil {
		return err
	}
	return nil
}

func (u *userUsecase) UpdateUsersLabelID(ctx context.Context, labelID int64, usernames []string) error {
	err := u.userRepo.UpdateUsersLabel(ctx, labelID, usernames)
	if err != nil {
		return err
	}
	return nil
}

func (u *userUsecase) GetUsersByLabelID(ctx context.Context, labelID int64) ([]string, error) {
	users, err := u.userRepo.GetUsersByLabelID(ctx, labelID)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u *userUsecase) RemoveUsersFromLabel(ctx context.Context, labelID int64, usernames []string) error {
	err := u.userRepo.RemoveUsersFromLabel(ctx, labelID, usernames)
	if err != nil {
		return err
	}
	return nil
}
