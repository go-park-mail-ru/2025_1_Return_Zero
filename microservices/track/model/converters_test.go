package model

import (
	"testing"

	trackProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/track"
	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/track/model/repository"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/track/model/usecase"
	"github.com/stretchr/testify/assert"
)

func TestPaginationFromUsecaseToRepository(t *testing.T) {
	usecasePagination := &usecaseModel.Pagination{
		Limit:  10,
		Offset: 20,
	}

	result := PaginationFromUsecaseToRepository(usecasePagination)

	assert.Equal(t, usecasePagination.Limit, result.Limit)
	assert.Equal(t, usecasePagination.Offset, result.Offset)
}

func TestFiltersFromUsecaseToRepository(t *testing.T) {
	usecasePagination := &usecaseModel.Pagination{
		Limit:  10,
		Offset: 20,
	}
	usecaseFilters := &usecaseModel.TrackFilters{
		Pagination: usecasePagination,
	}

	result := FiltersFromUsecaseToRepository(usecaseFilters)

	assert.Equal(t, usecaseFilters.Pagination.Limit, result.Pagination.Limit)
	assert.Equal(t, usecaseFilters.Pagination.Offset, result.Pagination.Offset)
}

func TestTrackFromRepositoryToUsecase(t *testing.T) {
	repoTrack := &repoModel.Track{
		ID:         1,
		Title:      "Test Track",
		Thumbnail:  "thumbnail.jpg",
		Duration:   180,
		AlbumID:    2,
		IsFavorite: true,
	}

	result := TrackFromRepositoryToUsecase(repoTrack)

	assert.Equal(t, repoTrack.ID, result.ID)
	assert.Equal(t, repoTrack.Title, result.Title)
	assert.Equal(t, repoTrack.Thumbnail, result.Thumbnail)
	assert.Equal(t, repoTrack.Duration, result.Duration)
	assert.Equal(t, repoTrack.AlbumID, result.AlbumID)
	assert.Equal(t, repoTrack.IsFavorite, result.IsFavorite)
}

func TestTrackDetailedFromRepositoryToUsecase(t *testing.T) {
	repoTrack := &repoModel.TrackWithFileKey{
		Track: repoModel.Track{
			ID:         1,
			Title:      "Test Track",
			Thumbnail:  "thumbnail.jpg",
			Duration:   180,
			AlbumID:    2,
			IsFavorite: true,
		},
		FileKey: "file-key",
	}
	fileUrl := "https://example.com/tracks/1.mp3"

	result := TrackDetailedFromRepositoryToUsecase(repoTrack, fileUrl)

	assert.Equal(t, repoTrack.ID, result.Track.ID)
	assert.Equal(t, repoTrack.Title, result.Track.Title)
	assert.Equal(t, repoTrack.Thumbnail, result.Track.Thumbnail)
	assert.Equal(t, repoTrack.Duration, result.Track.Duration)
	assert.Equal(t, repoTrack.AlbumID, result.Track.AlbumID)
	assert.Equal(t, repoTrack.IsFavorite, result.Track.IsFavorite)
	assert.Equal(t, fileUrl, result.FileUrl)
}

func TestTrackListFromRepositoryToUsecase(t *testing.T) {
	repoTracks := []*repoModel.Track{
		{
			ID:         1,
			Title:      "Test Track 1",
			Thumbnail:  "thumbnail1.jpg",
			Duration:   180,
			AlbumID:    2,
			IsFavorite: true,
		},
		{
			ID:         2,
			Title:      "Test Track 2",
			Thumbnail:  "thumbnail2.jpg",
			Duration:   240,
			AlbumID:    2,
			IsFavorite: false,
		},
	}

	result := TrackListFromRepositoryToUsecase(repoTracks)

	assert.Len(t, result, len(repoTracks))
	for i, track := range repoTracks {
		assert.Equal(t, track.ID, result[i].ID)
		assert.Equal(t, track.Title, result[i].Title)
		assert.Equal(t, track.Thumbnail, result[i].Thumbnail)
		assert.Equal(t, track.Duration, result[i].Duration)
		assert.Equal(t, track.AlbumID, result[i].AlbumID)
		assert.Equal(t, track.IsFavorite, result[i].IsFavorite)
	}
}

func TestTrackStreamCreateDataFromUsecaseToRepository(t *testing.T) {
	usecaseStream := &usecaseModel.TrackStreamCreateData{
		TrackID: 1,
		UserID:  2,
	}

	result := TrackStreamCreateDataFromUsecaseToRepository(usecaseStream)

	assert.Equal(t, usecaseStream.TrackID, result.TrackID)
	assert.Equal(t, usecaseStream.UserID, result.UserID)
}

func TestTrackStreamUpdateDataFromUsecaseToRepository(t *testing.T) {
	usecaseStream := &usecaseModel.TrackStreamUpdateData{
		StreamID: 1,
		Duration: 180,
		UserID:   2,
	}

	result := TrackStreamUpdateDataFromUsecaseToRepository(usecaseStream)

	assert.Equal(t, usecaseStream.StreamID, result.StreamID)
	assert.Equal(t, usecaseStream.Duration, result.Duration)
}

func TestTrackFromUsecaseToProto(t *testing.T) {
	usecaseTrack := &usecaseModel.Track{
		ID:         1,
		Title:      "Test Track",
		Thumbnail:  "thumbnail.jpg",
		Duration:   180,
		AlbumID:    2,
		IsFavorite: true,
	}

	result := TrackFromUsecaseToProto(usecaseTrack)

	assert.Equal(t, usecaseTrack.ID, result.Id)
	assert.Equal(t, usecaseTrack.Title, result.Title)
	assert.Equal(t, usecaseTrack.Thumbnail, result.Thumbnail)
	assert.Equal(t, usecaseTrack.Duration, result.Duration)
	assert.Equal(t, usecaseTrack.AlbumID, result.AlbumId)
	assert.Equal(t, usecaseTrack.IsFavorite, result.IsFavorite)
}

func TestTrackDetailedFromUsecaseToProto(t *testing.T) {
	usecaseTrack := &usecaseModel.TrackDetailed{
		Track: usecaseModel.Track{
			ID:         1,
			Title:      "Test Track",
			Thumbnail:  "thumbnail.jpg",
			Duration:   180,
			AlbumID:    2,
			IsFavorite: true,
		},
		FileUrl: "https://example.com/tracks/1.mp3",
	}

	result := TrackDetailedFromUsecaseToProto(usecaseTrack)

	assert.Equal(t, usecaseTrack.Track.ID, result.Track.Id)
	assert.Equal(t, usecaseTrack.Track.Title, result.Track.Title)
	assert.Equal(t, usecaseTrack.Track.Thumbnail, result.Track.Thumbnail)
	assert.Equal(t, usecaseTrack.Track.Duration, result.Track.Duration)
	assert.Equal(t, usecaseTrack.Track.AlbumID, result.Track.AlbumId)
	assert.Equal(t, usecaseTrack.Track.IsFavorite, result.Track.IsFavorite)
	assert.Equal(t, usecaseTrack.FileUrl, result.FileUrl)
}

func TestTrackListFromUsecaseToProto(t *testing.T) {
	usecaseTracks := []*usecaseModel.Track{
		{
			ID:         1,
			Title:      "Test Track 1",
			Thumbnail:  "thumbnail1.jpg",
			Duration:   180,
			AlbumID:    2,
			IsFavorite: true,
		},
		{
			ID:         2,
			Title:      "Test Track 2",
			Thumbnail:  "thumbnail2.jpg",
			Duration:   240,
			AlbumID:    2,
			IsFavorite: false,
		},
	}

	result := TrackListFromUsecaseToProto(usecaseTracks)

	assert.Len(t, result.Tracks, len(usecaseTracks))
	for i, track := range usecaseTracks {
		assert.Equal(t, track.ID, result.Tracks[i].Id)
		assert.Equal(t, track.Title, result.Tracks[i].Title)
		assert.Equal(t, track.Thumbnail, result.Tracks[i].Thumbnail)
		assert.Equal(t, track.Duration, result.Tracks[i].Duration)
		assert.Equal(t, track.AlbumID, result.Tracks[i].AlbumId)
		assert.Equal(t, track.IsFavorite, result.Tracks[i].IsFavorite)
	}
}

func TestFiltersFromProtoToUsecase(t *testing.T) {
	protoPagination := &trackProto.Pagination{
		Limit:  10,
		Offset: 20,
	}
	protoFilters := &trackProto.Filters{
		Pagination: protoPagination,
	}

	result := FiltersFromProtoToUsecase(protoFilters)

	assert.Equal(t, protoFilters.Pagination.Limit, result.Pagination.Limit)
	assert.Equal(t, protoFilters.Pagination.Offset, result.Pagination.Offset)
}

func TestStreamIDFromUsecaseToProto(t *testing.T) {
	streamID := int64(123)

	result := StreamIDFromUsecaseToProto(streamID)

	assert.Equal(t, streamID, result.Id)
}

func TestTrackStreamCreateDataFromProtoToUsecase(t *testing.T) {
	protoTrackID := &trackProto.TrackID{Id: 1}
	protoUserID := &trackProto.UserID{Id: 2}
	protoStream := &trackProto.TrackStreamCreateData{
		TrackId: protoTrackID,
		UserId:  protoUserID,
	}

	result := TrackStreamCreateDataFromProtoToUsecase(protoStream)

	assert.Equal(t, protoStream.TrackId.Id, result.TrackID)
	assert.Equal(t, protoStream.UserId.Id, result.UserID)
}

func TestTrackStreamUpdateDataFromProtoToUsecase(t *testing.T) {
	protoStreamID := &trackProto.StreamID{Id: 1}
	protoUserID := &trackProto.UserID{Id: 2}
	protoStream := &trackProto.TrackStreamUpdateData{
		StreamId: protoStreamID,
		Duration: 180,
		UserId:   protoUserID,
	}

	result := TrackStreamUpdateDataFromProtoToUsecase(protoStream)

	assert.Equal(t, protoStream.StreamId.Id, result.StreamID)
	assert.Equal(t, protoStream.Duration, result.Duration)
	assert.Equal(t, protoStream.UserId.Id, result.UserID)
}

func TestTrackIDListFromProtoToUsecase(t *testing.T) {
	protoTrackIDs := []*trackProto.TrackID{
		{Id: 1},
		{Id: 2},
		{Id: 3},
	}
	protoUserID := &trackProto.UserID{Id: 10}
	protoTrackIDList := &trackProto.TrackIDList{
		Ids:    protoTrackIDs,
		UserId: protoUserID,
	}

	ids, userID := TrackIDListFromProtoToUsecase(protoTrackIDList)

	assert.Len(t, ids, len(protoTrackIDs))
	for i, id := range protoTrackIDs {
		assert.Equal(t, id.Id, ids[i])
	}
	assert.Equal(t, protoUserID.Id, userID)
}

func TestTrackIDListWithFiltersFromProtoToUsecase(t *testing.T) {
	protoTrackIDs := []*trackProto.TrackID{
		{Id: 1},
		{Id: 2},
		{Id: 3},
	}
	protoUserID := &trackProto.UserID{Id: 10}
	protoTrackIDList := &trackProto.TrackIDList{
		Ids:    protoTrackIDs,
		UserId: protoUserID,
	}
	protoPagination := &trackProto.Pagination{
		Limit:  10,
		Offset: 20,
	}
	protoFilters := &trackProto.Filters{
		Pagination: protoPagination,
	}
	protoTrackIDsWithFilters := &trackProto.TrackIDListWithFilters{
		Ids:     protoTrackIDList,
		Filters: protoFilters,
	}

	ids, filters, userID := TrackIDListWithFiltersFromProtoToUsecase(protoTrackIDsWithFilters)

	assert.Len(t, ids, len(protoTrackIDs))
	for i, id := range protoTrackIDs {
		assert.Equal(t, id.Id, ids[i])
	}
	assert.Equal(t, protoFilters.Pagination.Limit, filters.Pagination.Limit)
	assert.Equal(t, protoFilters.Pagination.Offset, filters.Pagination.Offset)
	assert.Equal(t, protoUserID.Id, userID)
}

func TestLikeRequestFromProtoToUsecase(t *testing.T) {
	protoTrackID := &trackProto.TrackID{Id: 1}
	protoUserID := &trackProto.UserID{Id: 2}
	protoLikeRequest := &trackProto.LikeRequest{
		TrackId: protoTrackID,
		UserId:  protoUserID,
		IsLike:  true,
	}

	result := LikeRequestFromProtoToUsecase(protoLikeRequest)

	assert.Equal(t, protoLikeRequest.TrackId.Id, result.TrackID)
	assert.Equal(t, protoLikeRequest.UserId.Id, result.UserID)
	assert.Equal(t, protoLikeRequest.IsLike, result.IsLike)
}

func TestLikeRequestFromUsecaseToRepository(t *testing.T) {
	usecaseLikeRequest := &usecaseModel.LikeRequest{
		TrackID: 1,
		UserID:  2,
		IsLike:  true,
	}

	result := LikeRequestFromUsecaseToRepository(usecaseLikeRequest)

	assert.Equal(t, usecaseLikeRequest.TrackID, result.TrackID)
	assert.Equal(t, usecaseLikeRequest.UserID, result.UserID)
}

func TestFavoriteRequestFromProtoToUsecase(t *testing.T) {
	protoRequestUserID := &trackProto.UserID{Id: 1}
	protoProfileUserID := &trackProto.UserID{Id: 2}
	protoPagination := &trackProto.Pagination{
		Limit:  10,
		Offset: 20,
	}
	protoFilters := &trackProto.Filters{
		Pagination: protoPagination,
	}
	protoFavoriteRequest := &trackProto.FavoriteRequest{
		RequestUserId: protoRequestUserID,
		ProfileUserId: protoProfileUserID,
		Filters:       protoFilters,
	}

	result := FavoriteRequestFromProtoToUsecase(protoFavoriteRequest)

	assert.Equal(t, protoFavoriteRequest.RequestUserId.Id, result.RequestUserID)
	assert.Equal(t, protoFavoriteRequest.ProfileUserId.Id, result.ProfileUserID)
	assert.Equal(t, protoFavoriteRequest.Filters.Pagination.Limit, result.Filters.Pagination.Limit)
	assert.Equal(t, protoFavoriteRequest.Filters.Pagination.Offset, result.Filters.Pagination.Offset)
}

func TestFavoriteRequestFromUsecaseToRepository(t *testing.T) {
	usecasePagination := &usecaseModel.Pagination{
		Limit:  10,
		Offset: 20,
	}
	usecaseFilters := &usecaseModel.TrackFilters{
		Pagination: usecasePagination,
	}
	usecaseFavoriteRequest := &usecaseModel.FavoriteRequest{
		RequestUserID: 1,
		ProfileUserID: 2,
		Filters:       usecaseFilters,
	}

	result := FavoriteRequestFromUsecaseToRepository(usecaseFavoriteRequest)

	assert.Equal(t, usecaseFavoriteRequest.RequestUserID, result.RequestUserID)
	assert.Equal(t, usecaseFavoriteRequest.ProfileUserID, result.ProfileUserID)
	assert.Equal(t, usecaseFavoriteRequest.Filters.Pagination.Limit, result.Filters.Pagination.Limit)
	assert.Equal(t, usecaseFavoriteRequest.Filters.Pagination.Offset, result.Filters.Pagination.Offset)
}

func TestTracksIDWithAlbumIDFromProtoToUsecase(t *testing.T) {
	protoTracks := []*trackProto.TrackLoad{
		{Title: "Track 1", File: []byte("file1.mp3")},
		{Title: "Track 2", File: []byte("file2.mp3")},
	}

	result := TracksIDWithAlbumIDFromProtoToUsecase(protoTracks)

	assert.Len(t, result, len(protoTracks))
	for i, track := range protoTracks {
		assert.Equal(t, track.Title, result[i].Title)
		assert.Equal(t, track.File, result[i].File)
	}
}

func TestTrackLoadFromUsecaseToRepository(t *testing.T) {
	usecaseTrack := &usecaseModel.TrackLoad{
		Title:    "Test Track",
		File:     []byte("file.mp3"),
		Position: 1,
	}

	result := TrackLoadFromUsecaseToRepository(usecaseTrack)

	assert.Equal(t, usecaseTrack.Title, result.Title)
	assert.Equal(t, usecaseTrack.File, result.File)
	assert.Equal(t, usecaseTrack.Position, result.Position)
}

