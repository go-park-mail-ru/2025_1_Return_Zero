package usecase

import (
	"context"
	"errors"
	"fmt"
	"io"

	artistProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/artist"
	authProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/auth"
	trackProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/track"
	userProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/user"
	cusstomErrors "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/customErrors"
	model "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/user"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/userAvatarFile"
)

var (
	ErrWrongUsername = errors.New("wrong username")
)

func NewUserUsecase(userClient *userProto.UserServiceClient, authClient *authProto.AuthServiceClient, userFileRepo userAvatarFile.Repository, artistClient *artistProto.ArtistServiceClient, trackClient *trackProto.TrackServiceClient) user.Usecase {
	return &userUsecase{
		userClient:   userClient,
		authClient:   authClient,
		userFileRepo: userFileRepo,
		trackClient:  trackClient,
		artistClient: artistClient,
	}
}

type userUsecase struct {
	userClient   *userProto.UserServiceClient
	authClient   *authProto.AuthServiceClient
	userFileRepo userAvatarFile.Repository
	artistClient *artistProto.ArtistServiceClient
	trackClient  *trackProto.TrackServiceClient
}

func (u *userUsecase) CreateUser(ctx context.Context, user *usecaseModel.User) (*usecaseModel.User, string, error) {
	newUser, err := (*u.userClient).CreateUser(ctx, model.RegisterDataFromUsecaseToProto(user))
	if err != nil {
		return nil, "", cusstomErrors.HandleUserGRPCError(err)
	}
	userUsecase := model.UserFromProtoToUsecase(newUser)
	avatar_url, err := u.userFileRepo.GetAvatarURL(ctx, userUsecase.AvatarUrl)
	if err != nil {
		return nil, "", err
	}
	userUsecase.AvatarUrl = avatar_url
	sessionID, err := (*u.authClient).CreateSession(ctx, model.UserIDFromUsecaseToProto(userUsecase.ID))
	if err != nil {
		return nil, "", cusstomErrors.HandleAuthGRPCError(err)
	}
	return userUsecase, model.SessionIDFromProtoToUsecase(sessionID), nil
}

func (u *userUsecase) GetUserBySID(ctx context.Context, SID string) (*usecaseModel.User, error) {
	id, err := (*u.authClient).GetSession(ctx, model.SessionIDFromUsecaseToProto(SID))
	if err != nil {
		return nil, cusstomErrors.HandleAuthGRPCError(err)
	}
	userID := model.UserIDFromProtoToUsecase(id)
	user, err := (*u.userClient).GetUserByID(ctx, model.UserIDFromUsecaseToProtoUser(userID))
	if err != nil {
		return nil, cusstomErrors.HandleUserGRPCError(err)
	}
	userUsecase := model.UserFromProtoToUsecase(user)
	avatar_url, err := u.userFileRepo.GetAvatarURL(ctx, userUsecase.AvatarUrl)
	if err != nil {
		return nil, err
	}
	userUsecase.AvatarUrl = avatar_url
	return userUsecase, nil
}

func (u *userUsecase) LoginUser(ctx context.Context, user *usecaseModel.User) (*usecaseModel.User, string, error) {
	loginUser, err := (*u.userClient).LoginUser(ctx, model.LoginDataFromUsecaseToProto(user))
	if err != nil {
		return nil, "", cusstomErrors.HandleUserGRPCError(err)
	}
	userUsecase := model.UserFromProtoToUsecase(loginUser)
	avatar_url, err := u.userFileRepo.GetAvatarURL(ctx, userUsecase.AvatarUrl)
	if err != nil {
		return nil, "", err
	}
	userUsecase.AvatarUrl = avatar_url
	sessionID, err := (*u.authClient).CreateSession(ctx, model.UserIDFromUsecaseToProto(userUsecase.ID))
	if err != nil {
		return nil, "", cusstomErrors.HandleAuthGRPCError(err)
	}
	return userUsecase, model.SessionIDFromProtoToUsecase(sessionID), nil
}

func (u *userUsecase) Logout(ctx context.Context, SID string) error {
	_, err := (*u.authClient).DeleteSession(ctx, model.SessionIDFromUsecaseToProto(SID))
	if err != nil {
		return cusstomErrors.HandleAuthGRPCError(err)
	}
	return nil
}

func (u *userUsecase) UploadAvatar(ctx context.Context, username string, fileAvatar io.Reader, ID int64) (string, error) {
	fileURL, err := u.userFileRepo.UploadUserAvatar(ctx, username, fileAvatar)
	if err != nil {
		return "", err
	}

	_, err = (*u.userClient).UploadAvatar(ctx, model.AvatarDataFromUsecaseToProto(fileURL, ID))
	if err != nil {
		return "", cusstomErrors.HandleUserGRPCError(err)
	}

	avatarURL, err := u.userFileRepo.GetAvatarURL(ctx, fileURL)
	if err != nil {
		return "", err
	}
	return avatarURL, nil
}

func (u *userUsecase) DeleteUser(ctx context.Context, user *usecaseModel.User, SID string) error {
	_, err := (*u.userClient).DeleteUser(ctx, model.DeleteUserFromUsecaseToProto(user))
	if err != nil {
		return cusstomErrors.HandleUserGRPCError(err)
	}
	_, err = (*u.authClient).DeleteSession(ctx, model.SessionIDFromUsecaseToProto(SID))
	if err != nil {
		return cusstomErrors.HandleAuthGRPCError(err)
	}
	return nil
}

func (u *userUsecase) GetArtistsListened(ctx context.Context, username string) (int64, error) {
	id, err := (*u.userClient).GetIDByUsername(ctx, model.UsernameFromUsecaseToProto(username))
	if err != nil {
		return -1, cusstomErrors.HandleUserGRPCError(err)
	}
	userID := model.UserIDFromProtoToUsecaseUser(id)
	artistListened, err := (*u.artistClient).GetArtistsListenedByUserID(ctx, model.UserIDFromUsecaseToProtoArtist(userID))
	if err != nil {
		return -1, cusstomErrors.HandleArtistGRPCError(err)
	}
	artistListenedUsecase := model.ArtistsListenedFromProtoToUsecase(artistListened)
	return artistListenedUsecase, nil
}

func (u *userUsecase) GetTracksListened(ctx context.Context, username string) (int64, error) {
	id, err := (*u.userClient).GetIDByUsername(ctx, model.UsernameFromUsecaseToProto(username))
	if err != nil {
		return -1, cusstomErrors.HandleUserGRPCError(err)
	}
	userID := model.UserIDFromProtoToUsecaseUser(id)
	trackListened, err := (*u.trackClient).GetTracksListenedByUserID(ctx, model.UserIDFromUsecaseToProtoTrack(userID))
	if err != nil {
		return -1, cusstomErrors.HandleTrackGRPCError(err)
	}
	trackListenedUsecase := model.TracksListenedFromProtoToUsecase(trackListened)
	return trackListenedUsecase, nil
}

func (u *userUsecase) GetMinutesListened(ctx context.Context, username string) (int64, error) {
	id, err := (*u.userClient).GetIDByUsername(ctx, model.UsernameFromUsecaseToProto(username))
	if err != nil {
		return -1, cusstomErrors.HandleUserGRPCError(err)
	}
	userID := model.UserIDFromProtoToUsecaseUser(id)
	minutesListened, err := (*u.trackClient).GetMinutesListenedByUserID(ctx, model.UserIDFromUsecaseToProtoTrack(userID))
	if err != nil {
		return -1, cusstomErrors.HandleTrackGRPCError(err)
	}
	minutesListenedUsecase := model.MinutesListenedFromProtoToUsecase(minutesListened)
	return minutesListenedUsecase, nil
}

func (u *userUsecase) GetUserData(ctx context.Context, username string) (*usecaseModel.UserFullData, error) {
	userFullData, err := (*u.userClient).GetUserFullData(ctx, model.UsernameFromUsecaseToProto(username))
	if err != nil {
		return nil, cusstomErrors.HandleUserGRPCError(err)
	}
	fmt.Println("STEP 1")
	userFullDataUsecase := model.UserFullDataFromProtoToUsecase(userFullData)
	artistsListened, err := u.GetArtistsListened(ctx, username)
	if err != nil {
		return nil, err
	}
	fmt.Println("STEP 2")
	tracksListened, err := u.GetTracksListened(ctx, username)
	if err != nil {
		return nil, err
	}
	fmt.Println("STEP 3")
	minutesListened, err := u.GetMinutesListened(ctx, username)
	if err != nil {
		return nil, err
	}
	fmt.Println("STEP 4")
	stats := &usecaseModel.UserStatistics{
		ArtistsListened: artistsListened,
		TracksListened:  tracksListened,
		MinutesListened: minutesListened,
	}
	userFullDataUsecase.Statistics = stats
	avatarURL, err := u.userFileRepo.GetAvatarURL(ctx, userFullDataUsecase.AvatarUrl)
	if err != nil {
		return nil, err
	}
	userFullDataUsecase.AvatarUrl = avatarURL
	return userFullDataUsecase, nil
}

func (u *userUsecase) ChangeUserData(ctx context.Context, username string, userChangeData *usecaseModel.UserChangeSettings) (*usecaseModel.UserFullData, error) {
	if userChangeData.Privacy != nil {
		_, err := (*u.userClient).ChangeUserPrivacySettings(ctx, model.PrivacyFromUsecaseToProto(username, userChangeData.Privacy))
		if err != nil {
			return nil, cusstomErrors.HandleUserGRPCError(err)
		}
	}
	_, err := (*u.userClient).ChangeUserData(ctx, model.ChangeUserDataFromUsecaseToProto(username, userChangeData))
	if err != nil {
		return nil, cusstomErrors.HandleUserGRPCError(err)
	}
	updatedUsername := username
	if userChangeData.NewUsername != "" {
		updatedUsername = userChangeData.NewUsername
	}
	userFullDataUsecase, err := u.GetUserData(ctx, updatedUsername)
	if err != nil {
		return nil, err
	}
	return userFullDataUsecase, nil
}

func (u *userUsecase) GetUserByID(ctx context.Context, id int64) (*usecaseModel.User, error) {
	user, err := (*u.userClient).GetUserByID(ctx, model.UserIDFromUsecaseToProtoUser(id))
	if err != nil {
		return nil, cusstomErrors.HandleUserGRPCError(err)
	}
	userUsecase := model.UserFromProtoToUsecase(user)
	avatarURL, err := u.userFileRepo.GetAvatarURL(ctx, userUsecase.AvatarUrl)
	if err != nil {
		return nil, err
	}
	userUsecase.AvatarUrl = avatarURL
	return userUsecase, nil
}
