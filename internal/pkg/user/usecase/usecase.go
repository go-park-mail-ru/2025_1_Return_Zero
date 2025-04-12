package usecase

import (
	"context"
	"errors"
	"io"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/auth"
	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/repository"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/user"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/userAvatarFile"
)

var (
	ErrWrongUsername = errors.New("wrong username")
)

func NewUserUsecase(userRepo user.Repository, authRepo auth.Repository, userFileRepo userAvatarFile.Repository) user.Usecase {
	return userUsecase{
		userRepo:     userRepo,
		authRepo:     authRepo,
		userFileRepo: userFileRepo,
	}
}

type userUsecase struct {
	userRepo     user.Repository
	authRepo     auth.Repository
	userFileRepo userAvatarFile.Repository
}

func toUsecaseModel(user *repoModel.User) *usecaseModel.User {
	return &usecaseModel.User{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.Password,
		AvatarUrl: user.Thumbnail,
	}
}

func userAndSettingsToUsecase(userAndSettings *repoModel.UserAndSettings) *usecaseModel.UserAndSettings {
	return &usecaseModel.UserAndSettings{
		Username:                userAndSettings.Username,
		AvatarUrl:               userAndSettings.Thumbnail,
		IsPublicPlaylists:       userAndSettings.IsPublicPlaylists,
		IsPublicMinutesListened: userAndSettings.IsPublicMinutesListened,
		IsPublicFavoriteArtists: userAndSettings.IsPublicFavoriteArtists,
		IsPublicTracksListened:  userAndSettings.IsPublicTracksListened,
		IsPublicFavoriteTracks:  userAndSettings.IsPublicFavoriteTracks,
		IsPublicArtistsListened: userAndSettings.IsPublicArtistsListened,
	}
}

func (u userUsecase) CreateUser(ctx context.Context, user *usecaseModel.User) (*usecaseModel.User, string, error) {
	repoUser := &repoModel.User{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	}
	newUser, err := u.userRepo.CreateUser(ctx, repoUser)
	if err != nil {
		return nil, "", err
	}

	userUsecase := toUsecaseModel(newUser)
	sessionID, err := u.authRepo.CreateSession(ctx, newUser.ID)
	if err != nil {
		return nil, "", err
	}
	return userUsecase, sessionID, nil
}

func (u userUsecase) GetUserBySID(ctx context.Context, SID string) (*usecaseModel.User, error) {
	id, err := u.authRepo.GetSession(ctx, SID)
	if err != nil {
		return nil, err
	}
	user, err := u.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	avatar_url, err := u.userFileRepo.GetAvatarURL(user.Thumbnail)
	if err != nil {
		return nil, err
	}
	usecaseUser := toUsecaseModel(user)
	usecaseUser.AvatarUrl = avatar_url
	return usecaseUser, nil
}

func (u userUsecase) LoginUser(ctx context.Context, user *usecaseModel.User) (*usecaseModel.User, string, error) {
	repoUser := &repoModel.User{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	}
	loginUser, err := u.userRepo.LoginUser(ctx, repoUser)
	if err != nil {
		return nil, "", err
	}
	avatar_url, err := u.userFileRepo.GetAvatarURL(loginUser.Thumbnail)
	if err != nil {
		return nil, "", err
	}
	usecaseUser := toUsecaseModel(loginUser)
	usecaseUser.AvatarUrl = avatar_url
	sessionID, err := u.authRepo.CreateSession(ctx, loginUser.ID)
	if err != nil {
		return nil, "", err
	}
	return usecaseUser, sessionID, nil
}

func (u userUsecase) Logout(ctx context.Context, SID string) error {
	err := u.authRepo.DeleteSession(ctx, SID)
	if err != nil {
		return err
	}
	return nil
}

func (u userUsecase) UploadAvatar(ctx context.Context, username string, fileAvatar io.Reader) error {
	fileURL, err := u.userFileRepo.UploadUserAvatar(ctx, username, fileAvatar)
	if err != nil {
		return err
	}

	err = u.userRepo.UploadAvatar(ctx, fileURL, username)
	if err != nil {
		return err
	}
	return nil
}

func (u userUsecase) ChangeUserData(ctx context.Context, username string, changeData *usecaseModel.ChangeUserData) (*usecaseModel.User, error) {
	if username != changeData.Username {
		return nil, ErrWrongUsername
	}
	repoChangeData := &repoModel.ChangeUserData{
		Username:    changeData.Username,
		Email:       changeData.Email,
		Password:    changeData.Password,
		NewUsername: changeData.NewUsername,
		NewEmail:    changeData.NewEmail,
		NewPassword: changeData.NewPassword,
	}
	user, err := u.userRepo.ChangeUserData(ctx, repoChangeData)
	if err != nil {
		return nil, err
	}
	usecaseUser := toUsecaseModel(user)

	return usecaseUser, nil
}

func (u userUsecase) DeleteUser(ctx context.Context, user *usecaseModel.User, SID string) error {
	repoUser := &repoModel.User{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	}
	err := u.userRepo.DeleteUser(ctx, repoUser)
	if err != nil {
		return err
	}
	err = u.authRepo.DeleteSession(ctx, SID)
	if err != nil {
		return err
	}
	err = u.userFileRepo.DeleteUserAvatar(ctx, user.Username)
	if err != nil {
		return err
	}
	return nil
}

func (u userUsecase) ChangeUserPrivacySettings(ctx context.Context, privacySettings *usecaseModel.PrivacySettings) error {
	privacy := &repoModel.PrivacySettings{
		Username:                privacySettings.Username,
		IsPublicPlaylists:       privacySettings.IsPublicPlaylists,
		IsPublicMinutesListened: privacySettings.IsPublicMinutesListened,
		IsPublicFavoriteArtists: privacySettings.IsPublicFavoriteArtists,
		IsPublicTracksListened:  privacySettings.IsPublicTracksListened,
		IsPublicFavoriteTracks:  privacySettings.IsPublicFavoriteTracks,
		IsPublicArtistsListened: privacySettings.IsPublicArtistsListened,
	}
	err := u.userRepo.ChangeUserPrivacySettings(ctx, privacy)
	if err != nil {
		return err
	}
	return nil
}

func (u userUsecase) GetUserData(ctx context.Context, username string) (*usecaseModel.UserAndSettings, error) {
	userAndSettings, err := u.userRepo.GetUserData(ctx, username)
	if err != nil {
		return nil, err
	}
	avatar_url, err := u.userFileRepo.GetAvatarURL(userAndSettings.Thumbnail)
	if err != nil {
		return nil, err
	}
	usecaseUserAndSettings := userAndSettingsToUsecase(userAndSettings)
	usecaseUserAndSettings.AvatarUrl = avatar_url
	return usecaseUserAndSettings, nil
}
