package model

import (
	"testing"

	protoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/user"
	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/user/model/repository"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/user/model/usecase"
	"github.com/stretchr/testify/assert"
)

func TestRegisterDataFromUsecaseToRepository(t *testing.T) {
	usecaseData := &usecaseModel.RegisterData{
		Username: "test",
		Email:    "test@test.com",
		Password: "test",
	}

	repoData := RegisterDataFromUsecaseToRepository(usecaseData)

	assert.Equal(t, usecaseData.Username, repoData.Username)
	assert.Equal(t, usecaseData.Email, repoData.Email)
	assert.Equal(t, usecaseData.Password, repoData.Password)
}

func TestUserFromRepositoryToUsecase(t *testing.T) {
	repoData := &repoModel.User{
		ID:        1,
		Username:  "test",
		Email:     "test@test.com",
		Thumbnail: "test.jpg",
	}

	usecaseData := UserFromRepositoryToUsecase(repoData)

	assert.Equal(t, repoData.Username, usecaseData.Username)
	assert.Equal(t, repoData.Email, usecaseData.Email)
	assert.Equal(t, repoData.Thumbnail, usecaseData.Thumbnail)
	assert.Equal(t, repoData.ID, usecaseData.Id)
}

func TestLoginDataFromUsecaseToRepository(t *testing.T) {
	usecaseData := &usecaseModel.LoginData{
		Username: "test",
		Email:    "test@test.com",
		Password: "test",
	}

	repoData := LoginDataFromUsecaseToRepository(usecaseData)

	assert.Equal(t, usecaseData.Username, repoData.Username)
	assert.Equal(t, usecaseData.Email, repoData.Email)
	assert.Equal(t, usecaseData.Password, repoData.Password)
}

func TestUserDeleteFromUsecaseToRepository(t *testing.T) {
	usecaseData := &usecaseModel.UserDelete{
		Username: "test",
		Email:    "test@test.com",
		Password: "test",
	}

	repoData := UserDeleteFromUsecaseToRepository(usecaseData)

	assert.Equal(t, usecaseData.Username, repoData.Username)
	assert.Equal(t, usecaseData.Email, repoData.Email)
	assert.Equal(t, usecaseData.Password, repoData.Password)
}

func TestChangeUserDataFromUsecaseToRepository(t *testing.T) {
	usecaseData := &usecaseModel.ChangeUserData{
		Password:    "test",
		NewUsername: "test",
		NewEmail:    "test@test.com",
		NewPassword: "test",
	}

	repoData := ChangeUserDataFromUsecaseToRepository(usecaseData)

	assert.Equal(t, usecaseData.Password, repoData.Password)
	assert.Equal(t, usecaseData.NewUsername, repoData.NewUsername)
	assert.Equal(t, usecaseData.NewEmail, repoData.NewEmail)
	assert.Equal(t, usecaseData.NewPassword, repoData.NewPassword)
}

func TestPrivacySettingsFromUsecaseToRepository(t *testing.T) {
	usecaseData := &usecaseModel.PrivacySettings{
		IsPublicPlaylists:       true,
		IsPublicMinutesListened: true,
		IsPublicFavoriteArtists: true,
		IsPublicTracksListened:  true,
		IsPublicFavoriteTracks:  true,
		IsPublicArtistsListened: true,
	}

	repoData := PrivacySettingsFromUsecaseToRepository(usecaseData)

	assert.Equal(t, usecaseData.IsPublicPlaylists, repoData.IsPublicPlaylists)
	assert.Equal(t, usecaseData.IsPublicMinutesListened, repoData.IsPublicMinutesListened)
	assert.Equal(t, usecaseData.IsPublicFavoriteArtists, repoData.IsPublicFavoriteArtists)
	assert.Equal(t, usecaseData.IsPublicTracksListened, repoData.IsPublicTracksListened)
	assert.Equal(t, usecaseData.IsPublicFavoriteTracks, repoData.IsPublicFavoriteTracks)
	assert.Equal(t, usecaseData.IsPublicArtistsListened, repoData.IsPublicArtistsListened)
}

func TestPrivacyFromRepositoryToUsecase(t *testing.T) {
	repoData := &repoModel.PrivacySettings{
		IsPublicPlaylists:       true,
		IsPublicMinutesListened: true,

		IsPublicFavoriteArtists: true,
		IsPublicTracksListened:  true,
		IsPublicFavoriteTracks:  true,
		IsPublicArtistsListened: true,
	}

	usecaseData := PrivacyFromRepositoryToUsecase(repoData)

	assert.Equal(t, usecaseData.IsPublicPlaylists, repoData.IsPublicPlaylists)
	assert.Equal(t, usecaseData.IsPublicMinutesListened, repoData.IsPublicMinutesListened)
	assert.Equal(t, usecaseData.IsPublicFavoriteArtists, repoData.IsPublicFavoriteArtists)
	assert.Equal(t, usecaseData.IsPublicTracksListened, repoData.IsPublicTracksListened)
	assert.Equal(t, usecaseData.IsPublicFavoriteTracks, repoData.IsPublicFavoriteTracks)
	assert.Equal(t, usecaseData.IsPublicArtistsListened, repoData.IsPublicArtistsListened)
}

func TestUserFullDataFromRepositoryToUsecase(t *testing.T) {
	repoData := &repoModel.UserFullData{
		Username:  "test",
		Email:     "test@test.com",
		Thumbnail: "test.jpg",
		Privacy: &repoModel.PrivacySettings{
			IsPublicPlaylists:       true,
			IsPublicMinutesListened: true,
			IsPublicFavoriteArtists: true,
			IsPublicTracksListened:  true,
			IsPublicFavoriteTracks:  true,
			IsPublicArtistsListened: true,
		},
	}

	usecaseData := UserFullDataFromRepositoryToUsecase(repoData)

	assert.Equal(t, usecaseData.Username, repoData.Username)
	assert.Equal(t, usecaseData.Email, repoData.Email)
	assert.Equal(t, usecaseData.Thumbnail, repoData.Thumbnail)
	assert.Equal(t, usecaseData.Privacy, PrivacyFromRepositoryToUsecase(repoData.Privacy))
}

func TestRegisterDataFromProtoToUsecase(t *testing.T) {
	protoData := &protoModel.RegisterData{
		Username: "test",
		Email:    "test@test.com",
		Password: "test",
	}

	usecaseData := RegisterDataFromProtoToUsecase(protoData)

	assert.Equal(t, usecaseData.Username, protoData.Username)
	assert.Equal(t, usecaseData.Email, protoData.Email)
	assert.Equal(t, usecaseData.Password, protoData.Password)
}

func TestUserFrontFromUsecaseToProto(t *testing.T) {
	usecaseData := &usecaseModel.UserFront{
		Id:        1,
		Username:  "test",
		Email:     "test@test.com",
		Thumbnail: "test.jpg",
	}

	protoData := UserFrontFromUsecaseToProto(usecaseData)

	assert.Equal(t, protoData.Id, usecaseData.Id)
	assert.Equal(t, protoData.Username, usecaseData.Username)
	assert.Equal(t, protoData.Email, usecaseData.Email)
	assert.Equal(t, protoData.Avatar, usecaseData.Thumbnail)
}

func TestLoginDataFromProtoToUsecase(t *testing.T) {
	protoData := &protoModel.LoginData{
		Username: "test",
		Email:    "test@test.com",
		Password: "test",
	}

	usecaseData := LoginDataFromProtoToUsecase(protoData)

	assert.Equal(t, usecaseData.Username, protoData.Username)
	assert.Equal(t, usecaseData.Email, protoData.Email)
	assert.Equal(t, usecaseData.Password, protoData.Password)
}

func TestUserDeleteFromProtoToUsecase(t *testing.T) {
	protoData := &protoModel.UserDelete{
		Username: "test",
		Email:    "test@test.com",
		Password: "test",
	}

	usecaseData := UserDeleteFromProtoToUsecase(protoData)

	assert.Equal(t, usecaseData.Username, protoData.Username)
	assert.Equal(t, usecaseData.Email, protoData.Email)
	assert.Equal(t, usecaseData.Password, protoData.Password)
}

func TestChangeUserDataFromProtoToUsecase(t *testing.T) {
	protoData := &protoModel.ChangeUserDataMessage{
		Password:    "test",
		NewUsername: "test",
		NewEmail:    "test@test.com",
		NewPassword: "test",
	}

	usecaseData := ChangeUserDataFromProtoToUsecase(protoData)

	assert.Equal(t, usecaseData.Password, protoData.Password)
	assert.Equal(t, usecaseData.NewUsername, protoData.NewUsername)
	assert.Equal(t, usecaseData.NewEmail, protoData.NewEmail)
	assert.Equal(t, usecaseData.NewPassword, protoData.NewPassword)
}

func TestPrivacySettingsFromProtoToUsecase(t *testing.T) {
	protoData := &protoModel.PrivacySettings{
		IsPublicPlaylists:       true,
		IsPublicMinutesListened: true,
		IsPublicFavoriteArtists: true,
		IsPublicTracksListened:  true,
		IsPublicFavoriteTracks:  true,
		IsPublicArtistsListened: true,
	}

	usecaseData := PrivacySettingsFromProtoToUsecase(protoData)

	assert.Equal(t, usecaseData.IsPublicPlaylists, protoData.IsPublicPlaylists)
	assert.Equal(t, usecaseData.IsPublicMinutesListened, protoData.IsPublicMinutesListened)
	assert.Equal(t, usecaseData.IsPublicFavoriteArtists, protoData.IsPublicFavoriteArtists)
	assert.Equal(t, usecaseData.IsPublicTracksListened, protoData.IsPublicTracksListened)
	assert.Equal(t, usecaseData.IsPublicFavoriteTracks, protoData.IsPublicFavoriteTracks)
	assert.Equal(t, usecaseData.IsPublicArtistsListened, protoData.IsPublicArtistsListened)
}

func TestUserFullDataFromUsecaseToProto(t *testing.T) {
	usecaseData := &usecaseModel.UserFullData{
		Username:  "test",
		Email:     "test@test.com",
		Thumbnail: "test.jpg",
		Privacy: &usecaseModel.PrivacySettings{
			IsPublicPlaylists:       true,
			IsPublicMinutesListened: true,
			IsPublicFavoriteArtists: true,
			IsPublicTracksListened:  true,
			IsPublicFavoriteTracks:  true,
			IsPublicArtistsListened: true,
		},
	}

	protoData := UserFullDataFromUsecaseToProto(usecaseData)

	assert.Equal(t, protoData.Username, usecaseData.Username)
	assert.Equal(t, protoData.Email, usecaseData.Email)
	assert.Equal(t, protoData.Avatar, usecaseData.Thumbnail)
	assert.Equal(t, protoData.Privacy, PrivacySettingsFromUsecaseToProto(usecaseData.Privacy))
}

func TestUserIDFromUsecaseToProto(t *testing.T) {
	usecaseData := int64(1)

	protoData := UserIDFromUsecaseToProto(usecaseData)

	assert.Equal(t, protoData.Id, usecaseData)
}

func TestAvatarUrlFromUsecaseToProto(t *testing.T) {
	usecaseData := "test.jpg"

	protoData := AvatarUrlFromUsecaseToProto(usecaseData)

	assert.Equal(t, protoData.Url, usecaseData)
}

func TestFileKeyFromUsecaseToProto(t *testing.T) {
	usecaseData := "test.jpg"

	protoData := FileKeyFromUsecaseToProto(usecaseData)

	assert.Equal(t, protoData.FileKey, usecaseData)
}

func TestLabelIDFromUsecaseToProto(t *testing.T) {
	usecaseData := int64(1)

	protoData := LabelIDFromUsecaseToProto(usecaseData)

	assert.Equal(t, protoData.Id, usecaseData)
}

func TestLabelFromUsecaseToProto(t *testing.T) {
	usecaseData := &usecaseModel.Label{
		ID:      1,
		Name:    "Test Label",
		Members: []string{"1", "2", "3"},
	}

	protoData := LabelFromUsecaseToProto(usecaseData)

	assert.Equal(t, protoData.Id, usecaseData.ID)
	assert.Equal(t, protoData.Name, usecaseData.Name)
	assert.Equal(t, protoData.Usernames, usecaseData.Members)
}

func TestUserToFrontFromRepositoryToUsecase(t *testing.T) {
	repoData := &repoModel.User{
		ID:        1,
		Username:  "test",
		Email:     "test@test.com",
		Thumbnail: "test.jpg",
		LabelId:   1,
	}

	usecaseData := UserToFrontFromRepositoryToUsecase(repoData)

	assert.Equal(t, usecaseData.Id, repoData.ID)
	assert.Equal(t, usecaseData.Username, repoData.Username)
	assert.Equal(t, usecaseData.Email, repoData.Email)
	assert.Equal(t, usecaseData.Thumbnail, repoData.Thumbnail)
}

func TestUsersToFrontFromRepositoryToUsecase(t *testing.T) {
	repoData := []*repoModel.User{
		{
			ID:        1,
			Username:  "test",
			Email:     "test@test.com",
			Thumbnail: "test.jpg",
		},
	}

	usecaseData := UsersToFrontFromRepositoryToUsecase(repoData)

	assert.Equal(t, usecaseData[0].Id, repoData[0].ID)
	assert.Equal(t, usecaseData[0].Username, repoData[0].Username)
	assert.Equal(t, usecaseData[0].Email, repoData[0].Email)
	assert.Equal(t, usecaseData[0].Thumbnail, repoData[0].Thumbnail)
}

func TestUserFromUsecaseToProto(t *testing.T) {
	usecaseData := &usecaseModel.UserFront{
		Id:		1,
		Username:  "test",
		Email:     "test@test.com",
		Thumbnail: "test.jpg",
	}

	protoData := UserToFrontFromUsecaseToProto(usecaseData)

	assert.Equal(t, protoData.Id, usecaseData.Id)
	assert.Equal(t, protoData.Username, usecaseData.Username)
	assert.Equal(t, protoData.Email, usecaseData.Email)
	assert.Equal(t, protoData.Avatar, usecaseData.Thumbnail)
}
