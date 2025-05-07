package model

import (
	"bytes"

	protoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/playlist"
	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/playlist/model/repository"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/playlist/model/usecase"
)

func CreatePlaylistRequestFromProtoToUsecase(proto *protoModel.CreatePlaylistRequest) *usecaseModel.CreatePlaylistRequest {
	return &usecaseModel.CreatePlaylistRequest{
		Title:     proto.GetTitle(),
		UserID:    proto.GetUserId(),
		Thumbnail: proto.GetThumbnail(),
		IsPublic:  proto.GetIsPublic(),
	}
}

func CreatePlaylistRequestFromUsecaseToRepository(usecase *usecaseModel.CreatePlaylistRequest) *repoModel.CreatePlaylistRequest {
	return &repoModel.CreatePlaylistRequest{
		Title:     usecase.Title,
		UserID:    usecase.UserID,
		Thumbnail: usecase.Thumbnail,
		IsPublic:  usecase.IsPublic,
	}
}

func PlaylistFromRepositoryToProto(repo *repoModel.Playlist) *protoModel.Playlist {
	return &protoModel.Playlist{
		Id:        repo.ID,
		Title:     repo.Title,
		Thumbnail: repo.Thumbnail,
	}
}

func PlaylistFromUsecaseToProto(usecase *usecaseModel.Playlist) *protoModel.Playlist {
	return &protoModel.Playlist{
		Id:        usecase.ID,
		Title:     usecase.Title,
		Thumbnail: usecase.Thumbnail,
		UserId:    usecase.UserID,
	}
}

func PlaylistFromRepositoryToUsecase(repo *repoModel.Playlist) *usecaseModel.Playlist {
	return &usecaseModel.Playlist{
		ID:        repo.ID,
		Title:     repo.Title,
		Thumbnail: repo.Thumbnail,
		UserID:    repo.UserID,
	}
}

func UploadPlaylistThumbnailRequestFromProtoToUsecase(proto *protoModel.UploadPlaylistThumbnailRequest) *usecaseModel.UploadPlaylistThumbnailRequest {
	return &usecaseModel.UploadPlaylistThumbnailRequest{
		Title:     proto.GetTitle(),
		Thumbnail: bytes.NewBuffer(proto.GetThumbnail()),
	}
}

func GetCombinedPlaylistsByUserIDRequestFromProtoToUsecase(proto *protoModel.GetCombinedPlaylistsByUserIDRequest) *usecaseModel.GetCombinedPlaylistsByUserIDRequest {
	return &usecaseModel.GetCombinedPlaylistsByUserIDRequest{
		UserID: proto.GetUserId(),
	}
}

func PlaylistsFromRepositoryToUsecase(repo []*repoModel.Playlist) []*usecaseModel.Playlist {
	usecase := make([]*usecaseModel.Playlist, len(repo))
	for i, playlist := range repo {
		usecase[i] = PlaylistFromRepositoryToUsecase(playlist)
	}
	return usecase
}

func PlaylistListFromRepositoryToUsecase(repo *repoModel.PlaylistList) *usecaseModel.PlaylistList {
	return &usecaseModel.PlaylistList{
		Playlists: PlaylistsFromRepositoryToUsecase(repo.Playlists),
	}
}

func PlaylistsFromUsecaseToProto(usecase []*usecaseModel.Playlist) []*protoModel.Playlist {
	proto := make([]*protoModel.Playlist, len(usecase))
	for i, playlist := range usecase {
		proto[i] = PlaylistFromUsecaseToProto(playlist)
	}
	return proto
}

func PlaylistListFromUsecaseToProto(usecase *usecaseModel.PlaylistList) *protoModel.PlaylistList {
	return &protoModel.PlaylistList{
		Playlists: PlaylistsFromUsecaseToProto(usecase.Playlists),
	}
}

func AddTrackToPlaylistRequestFromUsecaseToRepository(usecase *usecaseModel.AddTrackToPlaylistRequest) *repoModel.AddTrackToPlaylistRequest {
	return &repoModel.AddTrackToPlaylistRequest{
		UserID:     usecase.UserID,
		PlaylistID: usecase.PlaylistID,
		TrackID:    usecase.TrackID,
	}
}

func RemoveTrackFromPlaylistRequestFromUsecaseToRepository(usecase *usecaseModel.RemoveTrackFromPlaylistRequest) *repoModel.RemoveTrackFromPlaylistRequest {
	return &repoModel.RemoveTrackFromPlaylistRequest{
		UserID:     usecase.UserID,
		PlaylistID: usecase.PlaylistID,
		TrackID:    usecase.TrackID,
	}
}

func AddTrackToPlaylistRequestFromProtoToUsecase(proto *protoModel.AddTrackToPlaylistRequest) *usecaseModel.AddTrackToPlaylistRequest {
	return &usecaseModel.AddTrackToPlaylistRequest{
		UserID:     proto.GetUserId(),
		PlaylistID: proto.GetPlaylistId(),
		TrackID:    proto.GetTrackId(),
	}
}

func RemoveTrackFromPlaylistRequestFromProtoToUsecase(proto *protoModel.RemoveTrackFromPlaylistRequest) *usecaseModel.RemoveTrackFromPlaylistRequest {
	return &usecaseModel.RemoveTrackFromPlaylistRequest{
		UserID:     proto.GetUserId(),
		PlaylistID: proto.GetPlaylistId(),
		TrackID:    proto.GetTrackId(),
	}
}

func GetPlaylistTrackIdsRequestFromProtoToUsecase(proto *protoModel.GetPlaylistTrackIdsRequest) *usecaseModel.GetPlaylistTrackIdsRequest {
	return &usecaseModel.GetPlaylistTrackIdsRequest{
		UserID:     proto.GetUserId(),
		PlaylistID: proto.GetPlaylistId(),
	}
}

func GetPlaylistTrackIdsRequestFromUsecaseToRepository(usecase *usecaseModel.GetPlaylistTrackIdsRequest) *repoModel.GetPlaylistTrackIdsRequest {
	return &repoModel.GetPlaylistTrackIdsRequest{
		UserID:     usecase.UserID,
		PlaylistID: usecase.PlaylistID,
	}
}

func UpdatePlaylistRequestFromProtoToUsecase(proto *protoModel.UpdatePlaylistRequest) *usecaseModel.UpdatePlaylistRequest {
	return &usecaseModel.UpdatePlaylistRequest{
		UserID:     proto.GetUserId(),
		PlaylistID: proto.GetId(),
		Title:      proto.GetTitle(),
		Thumbnail:  proto.GetThumbnail(),
	}
}

func UpdatePlaylistRequestFromUsecaseToRepository(usecase *usecaseModel.UpdatePlaylistRequest) *repoModel.UpdatePlaylistRequest {
	return &repoModel.UpdatePlaylistRequest{
		UserID:     usecase.UserID,
		PlaylistID: usecase.PlaylistID,
		Title:      usecase.Title,
		Thumbnail:  usecase.Thumbnail,
	}
}

func GetPlaylistByIDRequestFromProtoToUsecase(proto *protoModel.GetPlaylistByIDRequest) *usecaseModel.GetPlaylistByIDRequest {
	return &usecaseModel.GetPlaylistByIDRequest{
		UserID:     proto.GetUserId(),
		PlaylistID: proto.GetId(),
	}
}

func RemovePlaylistRequestFromProtoToUsecase(proto *protoModel.RemovePlaylistRequest) *usecaseModel.RemovePlaylistRequest {
	return &usecaseModel.RemovePlaylistRequest{
		UserID:     proto.GetUserId(),
		PlaylistID: proto.GetPlaylistId(),
	}
}

func RemovePlaylistRequestFromUsecaseToRepository(usecase *usecaseModel.RemovePlaylistRequest) *repoModel.RemovePlaylistRequest {
	return &repoModel.RemovePlaylistRequest{
		UserID:     usecase.UserID,
		PlaylistID: usecase.PlaylistID,
	}
}

func GetPlaylistsToAddRequestFromUsecaseToRepository(usecase *usecaseModel.GetPlaylistsToAddRequest) *repoModel.GetPlaylistsToAddRequest {
	return &repoModel.GetPlaylistsToAddRequest{
		UserID:  usecase.UserID,
		TrackID: usecase.TrackID,
	}
}

func PlaylistWithIsIncludedTrackFromRepositoryToUsecase(repo *repoModel.PlaylistWithIsIncludedTrack) *usecaseModel.PlaylistWithIsIncludedTrack {
	return &usecaseModel.PlaylistWithIsIncludedTrack{
		Playlist:   PlaylistFromRepositoryToUsecase(repo.Playlist),
		IsIncluded: repo.IsIncluded,
	}
}

func PlaylistsFromRepositoryToUsecaseWithIsIncludedTrack(repo []*repoModel.PlaylistWithIsIncludedTrack) []*usecaseModel.PlaylistWithIsIncludedTrack {
	usecase := make([]*usecaseModel.PlaylistWithIsIncludedTrack, len(repo))
	for i, playlist := range repo {
		usecase[i] = PlaylistWithIsIncludedTrackFromRepositoryToUsecase(playlist)
	}
	return usecase
}

func GetPlaylistsToAddResponseFromRepositoryToUsecase(repo *repoModel.GetPlaylistsToAddResponse) *usecaseModel.GetPlaylistsToAddResponse {
	return &usecaseModel.GetPlaylistsToAddResponse{
		Playlists: PlaylistsFromRepositoryToUsecaseWithIsIncludedTrack(repo.Playlists),
	}
}

func GetPlaylistsToAddResponseFromUsecaseToProto(usecase *usecaseModel.GetPlaylistsToAddResponse) *protoModel.GetPlaylistsToAddResponse {
	return &protoModel.GetPlaylistsToAddResponse{
		Playlists: PlaylistsFromUsecaseToProtoWithIsIncludedTrack(usecase.Playlists),
	}
}

func GetPlaylistsToAddRequestFromProtoToUsecase(proto *protoModel.GetPlaylistsToAddRequest) *usecaseModel.GetPlaylistsToAddRequest {
	return &usecaseModel.GetPlaylistsToAddRequest{
		UserID:  proto.GetUserId(),
		TrackID: proto.GetTrackId(),
	}
}

func PlaylistWithIsIncludedTrackFromUsecaseToProto(usecase *usecaseModel.PlaylistWithIsIncludedTrack) *protoModel.PlaylistWithIsIncludedTrack {
	return &protoModel.PlaylistWithIsIncludedTrack{
		Playlist:        PlaylistFromUsecaseToProto(usecase.Playlist),
		IsIncludedTrack: usecase.IsIncluded,
	}
}

func PlaylistsFromUsecaseToProtoWithIsIncludedTrack(usecase []*usecaseModel.PlaylistWithIsIncludedTrack) []*protoModel.PlaylistWithIsIncludedTrack {
	proto := make([]*protoModel.PlaylistWithIsIncludedTrack, len(usecase))
	for i, playlist := range usecase {
		proto[i] = PlaylistWithIsIncludedTrackFromUsecaseToProto(playlist)
	}
	return proto
}

func UpdatePlaylistsPublisityByUserIDRequestFromProtoToUsecase(proto *protoModel.UpdatePlaylistsPublisityByUserIDRequest) *usecaseModel.UpdatePlaylistsPublisityByUserIDRequest {
	return &usecaseModel.UpdatePlaylistsPublisityByUserIDRequest{
		UserID:   proto.GetUserId(),
		IsPublic: proto.GetIsPublic(),
	}
}

func UpdatePlaylistsPublisityByUserIDRequestFromUsecaseToRepository(usecase *usecaseModel.UpdatePlaylistsPublisityByUserIDRequest) *repoModel.UpdatePlaylistsPublisityByUserIDRequest {
	return &repoModel.UpdatePlaylistsPublisityByUserIDRequest{
		UserID:   usecase.UserID,
		IsPublic: usecase.IsPublic,
	}
}

func LikePlaylistRequestFromProtoToUsecase(proto *protoModel.LikePlaylistRequest) *usecaseModel.LikePlaylistRequest {
	return &usecaseModel.LikePlaylistRequest{
		UserID:     proto.GetUserId(),
		PlaylistID: proto.GetPlaylistId(),
		IsLike:     proto.GetIsLike(),
	}
}

func LikePlaylistRequestFromUsecaseToRepository(usecase *usecaseModel.LikePlaylistRequest) *repoModel.LikePlaylistRequest {
	return &repoModel.LikePlaylistRequest{
		UserID:     usecase.UserID,
		PlaylistID: usecase.PlaylistID,
	}
}

func PlaylistWithIsLikedFromRepositoryToUsecase(repo *repoModel.PlaylistWithIsLiked) *usecaseModel.PlaylistWithIsLiked {
	return &usecaseModel.PlaylistWithIsLiked{
		Playlist: PlaylistFromRepositoryToUsecase(repo.Playlist),
		IsLiked:  repo.IsLiked,
	}
}

func PlaylistWithIsLikedFromUsecaseToProto(usecase *usecaseModel.PlaylistWithIsLiked) *protoModel.PlaylistWithIsLiked {
	return &protoModel.PlaylistWithIsLiked{
		Playlist: PlaylistFromUsecaseToProto(usecase.Playlist),
		IsLiked:  usecase.IsLiked,
	}
}

func GetProfilePlaylistsRequestFromUsecaseToRepository(usecase *usecaseModel.GetProfilePlaylistsRequest) *repoModel.GetProfilePlaylistsRequest {
	return &repoModel.GetProfilePlaylistsRequest{
		UserID: usecase.UserID,
	}
}

func GetProfilePlaylistsResponseFromRepositoryToUsecase(repo *repoModel.GetProfilePlaylistsResponse) *usecaseModel.GetProfilePlaylistsResponse {
	return &usecaseModel.GetProfilePlaylistsResponse{
		Playlists: PlaylistsFromRepositoryToUsecase(repo.Playlists),
	}
}

func GetProfilePlaylistsRequestFromProtoToUsecase(proto *protoModel.GetProfilePlaylistsRequest) *usecaseModel.GetProfilePlaylistsRequest {
	return &usecaseModel.GetProfilePlaylistsRequest{
		UserID: proto.GetUserId(),
	}
}

func GetProfilePlaylistsResponseFromUsecaseToProto(usecase *usecaseModel.GetProfilePlaylistsResponse) *protoModel.GetProfilePlaylistsResponse {
	return &protoModel.GetProfilePlaylistsResponse{
		Playlists: PlaylistsFromUsecaseToProto(usecase.Playlists),
	}
}

func SearchPlaylistsRequestFromProtoToUsecase(proto *protoModel.SearchPlaylistsRequest) *usecaseModel.SearchPlaylistsRequest {
	return &usecaseModel.SearchPlaylistsRequest{
		UserID: proto.GetUserId(),
		Query:  proto.GetQuery(),
	}
}

func SearchPlaylistsRequestFromUsecaseToRepository(usecase *usecaseModel.SearchPlaylistsRequest) *repoModel.SearchPlaylistsRequest {
	return &repoModel.SearchPlaylistsRequest{
		UserID: usecase.UserID,
		Query:  usecase.Query,
	}
}
