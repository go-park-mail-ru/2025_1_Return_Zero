package model

import (
	trackProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/track"
	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/track/model/repository"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/track/model/usecase"
)

func PaginationFromUsecaseToRepository(pagination *usecaseModel.Pagination) *repoModel.Pagination {
	return &repoModel.Pagination{
		Limit:  pagination.Limit,
		Offset: pagination.Offset,
	}
}

func FiltersFromUsecaseToRepository(filters *usecaseModel.TrackFilters) *repoModel.TrackFilters {
	return &repoModel.TrackFilters{
		Pagination: PaginationFromUsecaseToRepository(filters.Pagination),
	}
}

func TrackFromRepositoryToUsecase(track *repoModel.Track) *usecaseModel.Track {
	return &usecaseModel.Track{
		ID:         track.ID,
		Title:      track.Title,
		Thumbnail:  track.Thumbnail,
		Duration:   track.Duration,
		AlbumID:    track.AlbumID,
		IsFavorite: track.IsFavorite,
	}
}

func TrackDetailedFromRepositoryToUsecase(track *repoModel.TrackWithFileKey, trackFileUrl string) *usecaseModel.TrackDetailed {
	return &usecaseModel.TrackDetailed{
		Track:   *TrackFromRepositoryToUsecase(&track.Track),
		FileUrl: trackFileUrl,
	}
}

func TrackListFromRepositoryToUsecase(tracks []*repoModel.Track) []*usecaseModel.Track {
	usecaseTracks := make([]*usecaseModel.Track, len(tracks))
	for i, track := range tracks {
		usecaseTracks[i] = TrackFromRepositoryToUsecase(track)
	}
	return usecaseTracks
}

func TrackStreamCreateDataFromUsecaseToRepository(stream *usecaseModel.TrackStreamCreateData) *repoModel.TrackStreamCreateData {
	return &repoModel.TrackStreamCreateData{
		TrackID: stream.TrackID,
		UserID:  stream.UserID,
	}
}

func TrackStreamUpdateDataFromUsecaseToRepository(stream *usecaseModel.TrackStreamUpdateData) *repoModel.TrackStreamUpdateData {
	return &repoModel.TrackStreamUpdateData{
		StreamID: stream.StreamID,
		Duration: stream.Duration,
	}
}

func TrackFromUsecaseToProto(track *usecaseModel.Track) *trackProto.Track {
	return &trackProto.Track{
		Id:         track.ID,
		Title:      track.Title,
		Thumbnail:  track.Thumbnail,
		Duration:   track.Duration,
		AlbumId:    track.AlbumID,
		IsFavorite: track.IsFavorite,
	}
}

func TrackDetailedFromUsecaseToProto(track *usecaseModel.TrackDetailed) *trackProto.TrackDetailed {
	return &trackProto.TrackDetailed{
		Track:   TrackFromUsecaseToProto(&track.Track),
		FileUrl: track.FileUrl,
	}
}

func TrackListFromUsecaseToProto(tracks []*usecaseModel.Track) *trackProto.TrackList {
	protoTracks := make([]*trackProto.Track, len(tracks))
	for i, track := range tracks {
		protoTracks[i] = TrackFromUsecaseToProto(track)
	}
	return &trackProto.TrackList{
		Tracks: protoTracks,
	}
}

func FiltersFromProtoToUsecase(filters *trackProto.Filters) *usecaseModel.TrackFilters {
	return &usecaseModel.TrackFilters{
		Pagination: &usecaseModel.Pagination{
			Limit:  filters.Pagination.Limit,
			Offset: filters.Pagination.Offset,
		},
	}
}

func StreamIDFromUsecaseToProto(streamID int64) *trackProto.StreamID {
	return &trackProto.StreamID{
		Id: streamID,
	}
}

func TrackStreamCreateDataFromProtoToUsecase(stream *trackProto.TrackStreamCreateData) *usecaseModel.TrackStreamCreateData {
	return &usecaseModel.TrackStreamCreateData{
		TrackID: stream.TrackId.Id,
		UserID:  stream.UserId.Id,
	}
}

func TrackStreamUpdateDataFromProtoToUsecase(stream *trackProto.TrackStreamUpdateData) *usecaseModel.TrackStreamUpdateData {
	return &usecaseModel.TrackStreamUpdateData{
		StreamID: stream.StreamId.Id,
		Duration: stream.Duration,
		UserID:   stream.UserId.Id,
	}
}

func TrackIDListFromProtoToUsecase(ids *trackProto.TrackIDList) ([]int64, int64) {
	usecaseIDs := make([]int64, len(ids.Ids))
	for i, id := range ids.Ids {
		usecaseIDs[i] = id.Id
	}
	return usecaseIDs, ids.UserId.Id
}

func TrackIDListWithFiltersFromProtoToUsecase(trackIdsWithFilters *trackProto.TrackIDListWithFilters) ([]int64, *usecaseModel.TrackFilters, int64) {
	usecaseIDs, userID := TrackIDListFromProtoToUsecase(trackIdsWithFilters.Ids)
	return usecaseIDs, FiltersFromProtoToUsecase(trackIdsWithFilters.Filters), userID
}

func LikeRequestFromProtoToUsecase(likeRequest *trackProto.LikeRequest) *usecaseModel.LikeRequest {
	return &usecaseModel.LikeRequest{
		TrackID: likeRequest.TrackId.Id,
		UserID:  likeRequest.UserId.Id,
		IsLike:  likeRequest.IsLike,
	}
}

func LikeRequestFromUsecaseToRepository(likeRequest *usecaseModel.LikeRequest) *repoModel.LikeRequest {
	return &repoModel.LikeRequest{
		TrackID: likeRequest.TrackID,
		UserID:  likeRequest.UserID,
	}
}

func FavoriteRequestFromProtoToUsecase(favoriteRequest *trackProto.FavoriteRequest) *usecaseModel.FavoriteRequest {
	return &usecaseModel.FavoriteRequest{
		RequestUserID: favoriteRequest.RequestUserId.Id,
		ProfileUserID: favoriteRequest.ProfileUserId.Id,
		Filters:       FiltersFromProtoToUsecase(favoriteRequest.Filters),
	}
}

func FavoriteRequestFromUsecaseToRepository(favoriteRequest *usecaseModel.FavoriteRequest) *repoModel.FavoriteRequest {
	return &repoModel.FavoriteRequest{
		RequestUserID: favoriteRequest.RequestUserID,
		ProfileUserID: favoriteRequest.ProfileUserID,
		Filters:       FiltersFromUsecaseToRepository(favoriteRequest.Filters),
	}
}
