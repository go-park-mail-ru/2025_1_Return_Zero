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

// func userAndSettingsToUsecase(userAndSettings *repoModel.UserAndSettings) *usecaseModel.UserAndSettings {
// 	return &usecaseModel.UserAndSettings{
// 		Username:  userAndSettings.Username,
// 		Email:     userAndSettings.Email,
// 		AvatarUrl: userAndSettings.Thumbnail,
// 	}
// }

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
	avatar_url, err := u.userFileRepo.GetAvatarURL(newUser.Thumbnail)
	if err != nil {
		return nil, "", err
	}
	userUsecase := toUsecaseModel(newUser)
	userUsecase.AvatarUrl = avatar_url
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

func (u userUsecase) UploadAvatar(ctx context.Context, username string, fileAvatar io.Reader) (string, error) {
	fileURL, err := u.userFileRepo.UploadUserAvatar(ctx, username, fileAvatar)
	if err != nil {
		return "", err
	}

	err = u.userRepo.UploadAvatar(ctx, fileURL, username)
	if err != nil {
		return "", err
	}
	avatarURL, err := u.userFileRepo.GetAvatarURL(fileURL)
	if err != nil {
		return "", err
	}
	return avatarURL, nil
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

func (u userUsecase) GetUserData(ctx context.Context, username string) (*usecaseModel.UserAndSettings, error) {
	userAndSettings, err := u.userRepo.GetUserData(ctx, username)
	if err != nil {
		return nil, err
	}
	privacy := &usecaseModel.UserPrivacy{
		IsPublicPlaylists:       userAndSettings.IsPublicPlaylists,
		IsPublicMinutesListened: userAndSettings.IsPublicMinutesListened,
		IsPublicFavoriteArtists: userAndSettings.IsPublicFavoriteArtists,
		IsPublicTracksListened:  userAndSettings.IsPublicTracksListened,
		IsPublicFavoriteTracks:  userAndSettings.IsPublicFavoriteTracks,
		IsPublicArtistsListened: userAndSettings.IsPublicArtistsListened,
	}
	userStatistics, err := u.userRepo.GetUserStats(ctx, username)
	if err != nil {
		return nil, err
	}
	avatar_url, err := u.userFileRepo.GetAvatarURL(userAndSettings.Thumbnail)
	if err != nil {
		return nil, err
	}

	usecaseUserAndSettings := &usecaseModel.UserAndSettings{
		Username:  userAndSettings.Username,
		Email:     userAndSettings.Email,
		AvatarUrl: avatar_url,
		Privacy:   privacy,
		Statistics: &usecaseModel.UserStatistics{
			MinutesListened: userStatistics.MinutesListened,
			TracksListened:  userStatistics.TracksListened,
			ArtistsListened: userStatistics.ArtistsListened,
		},
	}
	return usecaseUserAndSettings, nil
}

func (u userUsecase) ChangeUserData(ctx context.Context, username string, userChangeData *usecaseModel.UserChangeSettings) (*usecaseModel.UserAndSettings, error) {
	privacyRepo := &repoModel.PrivacySettings{
		IsPublicPlaylists:       userChangeData.IsPublicPlaylists,
		IsPublicMinutesListened: userChangeData.IsPublicMinutesListened,
		IsPublicFavoriteArtists: userChangeData.IsPublicFavoriteArtists,
		IsPublicTracksListened:  userChangeData.IsPublicTracksListened,
		IsPublicFavoriteTracks:  userChangeData.IsPublicFavoriteTracks,
		IsPublicArtistsListened: userChangeData.IsPublicArtistsListened,
	}

	userDataRepo := &repoModel.ChangeUserData{
		Password:    userChangeData.Password,
		NewUsername: userChangeData.NewUsername,
		NewEmail:    userChangeData.NewEmail,
		NewPassword: userChangeData.NewPassword,
	}

	err := u.userRepo.ChangeUserPrivacySettings(ctx, username, privacyRepo)
	if err != nil {
		return nil, err
	}

	err = u.userRepo.ChangeUserData(ctx, username, userDataRepo)
	if err != nil {
		return nil, err
	}

	updatedUsername := username
	if userDataRepo.NewUsername != "" {
		updatedUsername = userDataRepo.NewUsername
	}

	userStats, err := u.userRepo.GetUserStats(ctx, updatedUsername)
	if err != nil {
		return nil, err
	}
	user, err := u.userRepo.GetUserData(ctx, updatedUsername)
	if err != nil {
		return nil, err
	}
	avatar_url, err := u.userFileRepo.GetAvatarURL(user.Thumbnail)
	if err != nil {
		return nil, err
	}
	privacy := &usecaseModel.UserPrivacy{
		IsPublicPlaylists:       userChangeData.IsPublicPlaylists,
		IsPublicMinutesListened: userChangeData.IsPublicMinutesListened,
		IsPublicFavoriteArtists: userChangeData.IsPublicFavoriteArtists,
		IsPublicTracksListened:  userChangeData.IsPublicTracksListened,
		IsPublicFavoriteTracks:  userChangeData.IsPublicFavoriteTracks,
		IsPublicArtistsListened: userChangeData.IsPublicArtistsListened,
	}
	userStatistics := &usecaseModel.UserStatistics{
		MinutesListened: userStats.MinutesListened,
		TracksListened:  userStats.TracksListened,
		ArtistsListened: userStats.ArtistsListened,
	}
	newUserData := &usecaseModel.UserAndSettings{
		Username:   user.Username,
		Email:      user.Email,
		AvatarUrl:  avatar_url,
		Privacy:    privacy,
		Statistics: userStatistics,
	}
	return newUserData, nil
}
