package model

import (
	protoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/user"
	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/user/model/repository"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/user/model/usecase"
)

func RegisterDataFromUsecaseToRepository(data *usecaseModel.RegisterData) *repoModel.RegisterData {
	return &repoModel.RegisterData{
		Username: data.Username,
		Email:    data.Email,
		Password: data.Password,
	}
}

func UserFromRepositoryToUsecase(data *repoModel.User) *usecaseModel.UserFront {
	return &usecaseModel.UserFront{
		Username:  data.Username,
		Email:     data.Email,
		Thumbnail: data.Thumbnail,
	}
}

func LoginDataFromUsecaseToRepository(data *usecaseModel.LoginData) *repoModel.LoginData {
	return &repoModel.LoginData{
		Username: data.Username,
		Email:    data.Email,
		Password: data.Password,
	}
}

func UserDeleteFromUsecaseToRepository(data *usecaseModel.UserDelete) *repoModel.UserDelete {
	return &repoModel.UserDelete{
		Username: data.Username,
		Email:    data.Email,
		Password: data.Password,
	}
}

func ChangeUserDataFromUsecaseToRepository(data *usecaseModel.ChangeUserData) *repoModel.ChangeUserData {
	return &repoModel.ChangeUserData{
		Password:     data.Password,
		NewUsername:  data.NewUsername,
		NewEmail:     data.NewEmail,
		NewPassword:  data.NewPassword,
	}
}

func PrivacySettingsFromUsecaseToRepository(data *usecaseModel.PrivacySettings) *repoModel.PrivacySettings {
	return &repoModel.PrivacySettings{
		IsPublicPlaylists:       data.IsPublicPlaylists,
		IsPublicMinutesListened: data.IsPublicMinutesListened,
		IsPublicFavoriteArtists: data.IsPublicFavoriteArtists,
		IsPublicTracksListened:  data.IsPublicTracksListened,
		IsPublicFavoriteTracks:  data.IsPublicFavoriteTracks,
		IsPublicArtistsListened: data.IsPublicArtistsListened,
	}
}

func PrivacyFromRepositoryToUsecase(data *repoModel.PrivacySettings) *usecaseModel.PrivacySettings {
	return &usecaseModel.PrivacySettings{
		IsPublicPlaylists:       data.IsPublicPlaylists,
		IsPublicMinutesListened: data.IsPublicMinutesListened,
		IsPublicFavoriteArtists: data.IsPublicFavoriteArtists,
		IsPublicTracksListened:  data.IsPublicTracksListened,
		IsPublicFavoriteTracks:  data.IsPublicFavoriteTracks,
		IsPublicArtistsListened: data.IsPublicArtistsListened,
	}
}

func StatisticsFromRepositoryToUsecase(data *repoModel.Statistics) *usecaseModel.Statistics {
	return &usecaseModel.Statistics{
		MinutesListened: data.MinutesListened,
		TracksListened:  data.TracksListened,
		ArtistsListened: data.ArtistsListened,
	}
}

func UserFullDataFromRepositoryToUsecase(data *repoModel.UserFullData) *usecaseModel.UserFullData {
	privacyUsecase := PrivacyFromRepositoryToUsecase(data.Privacy)
	statisticsUsecase := StatisticsFromRepositoryToUsecase(data.Statistics)
	return &usecaseModel.UserFullData{
		Username:   data.Username,
		Thumbnail:  data.Thumbnail,
		Email:      data.Email,
		Privacy:    privacyUsecase,
		Statistics: statisticsUsecase,
	}
}

func RegisterDataFromProtoToUsecase(data *protoModel.RegisterData) *usecaseModel.RegisterData {
	return &usecaseModel.RegisterData{
		Username: data.Username,
		Email:    data.Email,
		Password: data.Password,
	}
}

func UserFrontFromUsecaseToProto(data *usecaseModel.UserFront) *protoModel.UserFront {
	return &protoModel.UserFront{
		Username:  data.Username,
		Email:     data.Email,
		Avatar: data.Thumbnail,
	}
}

func LoginDataFromProtoToUsecase(data *protoModel.LoginData) *usecaseModel.LoginData {
	return &usecaseModel.LoginData{
		Username: data.Username,
		Email:    data.Email,
		Password: data.Password,
	}
}

func UserDeleteFromProtoToUsecase(data *protoModel.UserDelete) *usecaseModel.UserDelete {
	return &usecaseModel.UserDelete{
		Username: data.Username,
		Email:    data.Email,
		Password: data.Password,
	}
}

func ChangeUserDataFromProtoToUsecase(data *protoModel.ChangeUserDataMessage) *usecaseModel.ChangeUserData {
	return &usecaseModel.ChangeUserData{
		Password:     data.Password,
		NewUsername:  data.NewUsername,
		NewEmail:     data.NewEmail,
		NewPassword:  data.NewPassword,
	}
}

func PrivacySettingsFromProtoToUsecase(data *protoModel.PrivacySettings) *usecaseModel.PrivacySettings {
	return &usecaseModel.PrivacySettings{
		IsPublicPlaylists:       data.IsPublicPlaylists,
		IsPublicMinutesListened: data.IsPublicMinutesListened,
		IsPublicFavoriteArtists: data.IsPublicFavoriteArtists,
		IsPublicTracksListened:  data.IsPublicTracksListened,
		IsPublicFavoriteTracks:  data.IsPublicFavoriteTracks,
		IsPublicArtistsListened: data.IsPublicArtistsListened,
	}
}

func UserFullDataFromUsecaseToProto(data *usecaseModel.UserFullData) *protoModel.UserFullData {
	privacyProto := &protoModel.PrivacySettings{
		IsPublicPlaylists:       data.Privacy.IsPublicPlaylists,
		IsPublicMinutesListened: data.Privacy.IsPublicMinutesListened,
		IsPublicFavoriteArtists: data.Privacy.IsPublicFavoriteArtists,
		IsPublicTracksListened:  data.Privacy.IsPublicTracksListened,
		IsPublicFavoriteTracks:  data.Privacy.IsPublicFavoriteTracks,
		IsPublicArtistsListened: data.Privacy.IsPublicArtistsListened,
	}
	statisticsProto := &protoModel.Statistics{
		MinutesListened: data.Statistics.MinutesListened,
		TracksListened:  data.Statistics.TracksListened,
		ArtistsListened: data.Statistics.ArtistsListened,
	}
	return &protoModel.UserFullData{
		Username:   data.Username,
		Avatar:  data.Thumbnail,
		Email:      data.Email,
		Privacy:    privacyProto,
		Statistics: statisticsProto,
	}
}