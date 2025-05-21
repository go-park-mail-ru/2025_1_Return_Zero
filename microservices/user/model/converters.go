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
		Id:        data.ID,
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
		Password:    data.Password,
		NewUsername: data.NewUsername,
		NewEmail:    data.NewEmail,
		NewPassword: data.NewPassword,
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

func UserFullDataFromRepositoryToUsecase(data *repoModel.UserFullData) *usecaseModel.UserFullData {
	privacyUsecase := PrivacyFromRepositoryToUsecase(data.Privacy)
	return &usecaseModel.UserFullData{
		Username:  data.Username,
		Thumbnail: data.Thumbnail,
		Email:     data.Email,
		Privacy:   privacyUsecase,
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
		Username: data.Username,
		Email:    data.Email,
		Avatar:   data.Thumbnail,
		Id:       data.Id,
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
		Password:    data.Password,
		NewUsername: data.NewUsername,
		NewEmail:    data.NewEmail,
		NewPassword: data.NewPassword,
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

func PrivacySettingsFromUsecaseToProto(data *usecaseModel.PrivacySettings) *protoModel.PrivacySettings {
	return &protoModel.PrivacySettings{
		IsPublicPlaylists:       data.IsPublicPlaylists,
		IsPublicMinutesListened: data.IsPublicMinutesListened,
		IsPublicFavoriteArtists: data.IsPublicFavoriteArtists,
		IsPublicTracksListened:  data.IsPublicTracksListened,
		IsPublicFavoriteTracks:  data.IsPublicFavoriteTracks,
		IsPublicArtistsListened: data.IsPublicArtistsListened,
	}
}

func UserFullDataFromUsecaseToProto(data *usecaseModel.UserFullData) *protoModel.UserFullData {
	privacyProto := PrivacySettingsFromUsecaseToProto(data.Privacy)
	return &protoModel.UserFullData{
		Username: data.Username,
		Avatar:   data.Thumbnail,
		Email:    data.Email,
		Privacy:  privacyProto,
	}
}

func UserIDFromUsecaseToProto(id int64) *protoModel.UserID {
	return &protoModel.UserID{
		Id: id,
	}
}

func AvatarUrlFromUsecaseToProto(url string) *protoModel.AvatarUrl {
	return &protoModel.AvatarUrl{
		Url: url,
	}
}

func FileKeyFromUsecaseToProto(fileKey string) *protoModel.FileKey {
	return &protoModel.FileKey{
		FileKey: fileKey,
	}

}

func LabelIDFromUsecaseToProto(id int64) *protoModel.LabelID {
	return &protoModel.LabelID{
		Id: id,
	}
}

func LabelFromUsecaseToProto(label *usecaseModel.Label) *protoModel.Label {
	return &protoModel.Label{
		Id:        label.ID,
		Name:      label.Name,
		Usernames: label.Members,
	}
}

func UserToFrontFromRepositoryToUsecase(user *repoModel.User) *usecaseModel.UserFront {
	return &usecaseModel.UserFront{
		Username:  user.Username,
		Email:     user.Email,
		Thumbnail: user.Thumbnail,
		Id:        user.ID,
	}
}

func UsersToFrontFromRepositoryToUsecase(users []*repoModel.User) []*usecaseModel.UserFront {
	userFronts := make([]*usecaseModel.UserFront, len(users))
	for i, user := range users {
		userFronts[i] = UserToFrontFromRepositoryToUsecase(user)
	}
	return userFronts
}

func UserToFrontFromUsecaseToProto(user *usecaseModel.UserFront) *protoModel.UserFront {
	return &protoModel.UserFront{
		Username: user.Username,
		Email:    user.Email,
		Avatar:   user.Thumbnail,
		Id:       user.Id,
	}
}

