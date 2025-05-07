package model

import (
	"testing"

	"github.com/stretchr/testify/assert"

	protoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/playlist"
	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/playlist/model/repository"
)

func TestCreatePlaylistRequestConversions(t *testing.T) {
	protoReq := &protoModel.CreatePlaylistRequest{
		Title:     "Test Playlist",
		UserId:    123,
		Thumbnail: "thumbnail.jpg",
		IsPublic:  true,
	}

	usecaseReq := CreatePlaylistRequestFromProtoToUsecase(protoReq)
	assert.Equal(t, protoReq.GetTitle(), usecaseReq.Title)
	assert.Equal(t, protoReq.GetUserId(), usecaseReq.UserID)
	assert.Equal(t, protoReq.GetThumbnail(), usecaseReq.Thumbnail)
	assert.Equal(t, protoReq.GetIsPublic(), usecaseReq.IsPublic)

	repoReq := CreatePlaylistRequestFromUsecaseToRepository(usecaseReq)
	assert.Equal(t, usecaseReq.Title, repoReq.Title)
	assert.Equal(t, usecaseReq.UserID, repoReq.UserID)
	assert.Equal(t, usecaseReq.Thumbnail, repoReq.Thumbnail)
	assert.Equal(t, usecaseReq.IsPublic, repoReq.IsPublic)
}

func TestPlaylistConversions(t *testing.T) {
	repoPlaylist := &repoModel.Playlist{
		ID:        1,
		Title:     "Test Playlist",
		Thumbnail: "thumbnail.jpg",
		UserID:    123,
	}

	protoPlaylist := PlaylistFromRepositoryToProto(repoPlaylist)
	assert.Equal(t, repoPlaylist.ID, protoPlaylist.Id)
	assert.Equal(t, repoPlaylist.Title, protoPlaylist.Title)
	assert.Equal(t, repoPlaylist.Thumbnail, protoPlaylist.Thumbnail)

	usecasePlaylist := PlaylistFromRepositoryToUsecase(repoPlaylist)
	assert.Equal(t, repoPlaylist.ID, usecasePlaylist.ID)
	assert.Equal(t, repoPlaylist.Title, usecasePlaylist.Title)
	assert.Equal(t, repoPlaylist.Thumbnail, usecasePlaylist.Thumbnail)
	assert.Equal(t, repoPlaylist.UserID, usecasePlaylist.UserID)

	protoPlaylist2 := PlaylistFromUsecaseToProto(usecasePlaylist)
	assert.Equal(t, usecasePlaylist.ID, protoPlaylist2.Id)
	assert.Equal(t, usecasePlaylist.Title, protoPlaylist2.Title)
	assert.Equal(t, usecasePlaylist.Thumbnail, protoPlaylist2.Thumbnail)
	assert.Equal(t, usecasePlaylist.UserID, protoPlaylist2.UserId)
}

func TestPlaylistListConversions(t *testing.T) {
	repoPlaylists := []*repoModel.Playlist{
		{
			ID:        1,
			Title:     "Test Playlist 1",
			Thumbnail: "thumbnail1.jpg",
			UserID:    123,
		},
		{
			ID:        2,
			Title:     "Test Playlist 2",
			Thumbnail: "thumbnail2.jpg",
			UserID:    123,
		},
	}

	repoPlaylistList := &repoModel.PlaylistList{
		Playlists: repoPlaylists,
	}

	usecasePlaylistList := PlaylistListFromRepositoryToUsecase(repoPlaylistList)
	assert.Equal(t, len(repoPlaylistList.Playlists), len(usecasePlaylistList.Playlists))
	for i, playlist := range usecasePlaylistList.Playlists {
		assert.Equal(t, repoPlaylistList.Playlists[i].ID, playlist.ID)
		assert.Equal(t, repoPlaylistList.Playlists[i].Title, playlist.Title)
		assert.Equal(t, repoPlaylistList.Playlists[i].Thumbnail, playlist.Thumbnail)
		assert.Equal(t, repoPlaylistList.Playlists[i].UserID, playlist.UserID)
	}

	protoPlaylistList := PlaylistListFromUsecaseToProto(usecasePlaylistList)
	assert.Equal(t, len(usecasePlaylistList.Playlists), len(protoPlaylistList.Playlists))
	for i, playlist := range protoPlaylistList.Playlists {
		assert.Equal(t, usecasePlaylistList.Playlists[i].ID, playlist.Id)
		assert.Equal(t, usecasePlaylistList.Playlists[i].Title, playlist.Title)
		assert.Equal(t, usecasePlaylistList.Playlists[i].Thumbnail, playlist.Thumbnail)
		assert.Equal(t, usecasePlaylistList.Playlists[i].UserID, playlist.UserId)
	}
}

func TestUploadPlaylistThumbnailRequestConversion(t *testing.T) {
	thumbnail := []byte("test thumbnail data")
	protoReq := &protoModel.UploadPlaylistThumbnailRequest{
		Title:     "Test Playlist",
		Thumbnail: thumbnail,
	}

	usecaseReq := UploadPlaylistThumbnailRequestFromProtoToUsecase(protoReq)
	assert.Equal(t, protoReq.GetTitle(), usecaseReq.Title)

	thumbnailData := make([]byte, len(thumbnail))
	_, err := usecaseReq.Thumbnail.Read(thumbnailData)
	assert.NoError(t, err)
	assert.Equal(t, thumbnail, thumbnailData)
}

func TestAddTrackToPlaylistRequestConversions(t *testing.T) {
	protoReq := &protoModel.AddTrackToPlaylistRequest{
		UserId:     123,
		PlaylistId: 456,
		TrackId:    789,
	}

	usecaseReq := AddTrackToPlaylistRequestFromProtoToUsecase(protoReq)
	assert.Equal(t, protoReq.GetUserId(), usecaseReq.UserID)
	assert.Equal(t, protoReq.GetPlaylistId(), usecaseReq.PlaylistID)
	assert.Equal(t, protoReq.GetTrackId(), usecaseReq.TrackID)

	repoReq := AddTrackToPlaylistRequestFromUsecaseToRepository(usecaseReq)
	assert.Equal(t, usecaseReq.UserID, repoReq.UserID)
	assert.Equal(t, usecaseReq.PlaylistID, repoReq.PlaylistID)
	assert.Equal(t, usecaseReq.TrackID, repoReq.TrackID)
}

func TestRemoveTrackFromPlaylistRequestConversions(t *testing.T) {
	protoReq := &protoModel.RemoveTrackFromPlaylistRequest{
		UserId:     123,
		PlaylistId: 456,
		TrackId:    789,
	}

	usecaseReq := RemoveTrackFromPlaylistRequestFromProtoToUsecase(protoReq)
	assert.Equal(t, protoReq.GetUserId(), usecaseReq.UserID)
	assert.Equal(t, protoReq.GetPlaylistId(), usecaseReq.PlaylistID)
	assert.Equal(t, protoReq.GetTrackId(), usecaseReq.TrackID)

	repoReq := RemoveTrackFromPlaylistRequestFromUsecaseToRepository(usecaseReq)
	assert.Equal(t, usecaseReq.UserID, repoReq.UserID)
	assert.Equal(t, usecaseReq.PlaylistID, repoReq.PlaylistID)
	assert.Equal(t, usecaseReq.TrackID, repoReq.TrackID)
}

func TestGetPlaylistTrackIdsRequestConversions(t *testing.T) {
	protoReq := &protoModel.GetPlaylistTrackIdsRequest{
		UserId:     123,
		PlaylistId: 456,
	}

	usecaseReq := GetPlaylistTrackIdsRequestFromProtoToUsecase(protoReq)
	assert.Equal(t, protoReq.GetUserId(), usecaseReq.UserID)
	assert.Equal(t, protoReq.GetPlaylistId(), usecaseReq.PlaylistID)

	repoReq := GetPlaylistTrackIdsRequestFromUsecaseToRepository(usecaseReq)
	assert.Equal(t, usecaseReq.UserID, repoReq.UserID)
	assert.Equal(t, usecaseReq.PlaylistID, repoReq.PlaylistID)
}

func TestUpdatePlaylistRequestConversions(t *testing.T) {
	protoReq := &protoModel.UpdatePlaylistRequest{
		Id:        456,
		UserId:    123,
		Title:     "Updated Playlist",
		Thumbnail: "updated-thumbnail.jpg",
	}

	usecaseReq := UpdatePlaylistRequestFromProtoToUsecase(protoReq)
	assert.Equal(t, protoReq.GetUserId(), usecaseReq.UserID)
	assert.Equal(t, protoReq.GetId(), usecaseReq.PlaylistID)
	assert.Equal(t, protoReq.GetTitle(), usecaseReq.Title)
	assert.Equal(t, protoReq.GetThumbnail(), usecaseReq.Thumbnail)

	repoReq := UpdatePlaylistRequestFromUsecaseToRepository(usecaseReq)
	assert.Equal(t, usecaseReq.UserID, repoReq.UserID)
	assert.Equal(t, usecaseReq.PlaylistID, repoReq.PlaylistID)
	assert.Equal(t, usecaseReq.Title, repoReq.Title)
	assert.Equal(t, usecaseReq.Thumbnail, repoReq.Thumbnail)
}

func TestGetPlaylistByIDRequestConversion(t *testing.T) {
	protoReq := &protoModel.GetPlaylistByIDRequest{
		Id:     456,
		UserId: 123,
	}

	usecaseReq := GetPlaylistByIDRequestFromProtoToUsecase(protoReq)
	assert.Equal(t, protoReq.GetUserId(), usecaseReq.UserID)
	assert.Equal(t, protoReq.GetId(), usecaseReq.PlaylistID)
}

func TestRemovePlaylistRequestConversions(t *testing.T) {
	protoReq := &protoModel.RemovePlaylistRequest{
		UserId:     123,
		PlaylistId: 456,
	}

	usecaseReq := RemovePlaylistRequestFromProtoToUsecase(protoReq)
	assert.Equal(t, protoReq.GetUserId(), usecaseReq.UserID)
	assert.Equal(t, protoReq.GetPlaylistId(), usecaseReq.PlaylistID)

	repoReq := RemovePlaylistRequestFromUsecaseToRepository(usecaseReq)
	assert.Equal(t, usecaseReq.UserID, repoReq.UserID)
	assert.Equal(t, usecaseReq.PlaylistID, repoReq.PlaylistID)
}

func TestPlaylistWithIsIncludedTrackConversions(t *testing.T) {
	repoPlaylist := &repoModel.Playlist{
		ID:        1,
		Title:     "Test Playlist",
		Thumbnail: "thumbnail.jpg",
		UserID:    123,
	}

	repoPlaylistWithTrack := &repoModel.PlaylistWithIsIncludedTrack{
		Playlist:   repoPlaylist,
		IsIncluded: true,
	}

	usecasePlaylistWithTrack := PlaylistWithIsIncludedTrackFromRepositoryToUsecase(repoPlaylistWithTrack)
	assert.Equal(t, repoPlaylistWithTrack.Playlist.ID, usecasePlaylistWithTrack.Playlist.ID)
	assert.Equal(t, repoPlaylistWithTrack.Playlist.Title, usecasePlaylistWithTrack.Playlist.Title)
	assert.Equal(t, repoPlaylistWithTrack.Playlist.Thumbnail, usecasePlaylistWithTrack.Playlist.Thumbnail)
	assert.Equal(t, repoPlaylistWithTrack.Playlist.UserID, usecasePlaylistWithTrack.Playlist.UserID)
	assert.Equal(t, repoPlaylistWithTrack.IsIncluded, usecasePlaylistWithTrack.IsIncluded)

	protoPlaylistWithTrack := PlaylistWithIsIncludedTrackFromUsecaseToProto(usecasePlaylistWithTrack)
	assert.Equal(t, usecasePlaylistWithTrack.Playlist.ID, protoPlaylistWithTrack.Playlist.Id)
	assert.Equal(t, usecasePlaylistWithTrack.Playlist.Title, protoPlaylistWithTrack.Playlist.Title)
	assert.Equal(t, usecasePlaylistWithTrack.Playlist.Thumbnail, protoPlaylistWithTrack.Playlist.Thumbnail)
	assert.Equal(t, usecasePlaylistWithTrack.Playlist.UserID, protoPlaylistWithTrack.Playlist.UserId)
	assert.Equal(t, usecasePlaylistWithTrack.IsIncluded, protoPlaylistWithTrack.IsIncludedTrack)
}

func TestGetPlaylistsToAddConversions(t *testing.T) {
	protoReq := &protoModel.GetPlaylistsToAddRequest{
		UserId:  123,
		TrackId: 789,
	}

	usecaseReq := GetPlaylistsToAddRequestFromProtoToUsecase(protoReq)
	assert.Equal(t, protoReq.GetUserId(), usecaseReq.UserID)
	assert.Equal(t, protoReq.GetTrackId(), usecaseReq.TrackID)

	repoReq := GetPlaylistsToAddRequestFromUsecaseToRepository(usecaseReq)
	assert.Equal(t, usecaseReq.UserID, repoReq.UserID)
	assert.Equal(t, usecaseReq.TrackID, repoReq.TrackID)

	repoPlaylist := &repoModel.Playlist{
		ID:        1,
		Title:     "Test Playlist",
		Thumbnail: "thumbnail.jpg",
		UserID:    123,
	}

	repoPlaylistWithTrack := &repoModel.PlaylistWithIsIncludedTrack{
		Playlist:   repoPlaylist,
		IsIncluded: true,
	}

	repoResponse := &repoModel.GetPlaylistsToAddResponse{
		Playlists: []*repoModel.PlaylistWithIsIncludedTrack{repoPlaylistWithTrack},
	}

	usecaseResponse := GetPlaylistsToAddResponseFromRepositoryToUsecase(repoResponse)
	assert.Equal(t, len(repoResponse.Playlists), len(usecaseResponse.Playlists))
	assert.Equal(t, repoResponse.Playlists[0].Playlist.ID, usecaseResponse.Playlists[0].Playlist.ID)
	assert.Equal(t, repoResponse.Playlists[0].IsIncluded, usecaseResponse.Playlists[0].IsIncluded)

	protoResponse := GetPlaylistsToAddResponseFromUsecaseToProto(usecaseResponse)
	assert.Equal(t, len(usecaseResponse.Playlists), len(protoResponse.Playlists))
	assert.Equal(t, usecaseResponse.Playlists[0].Playlist.ID, protoResponse.Playlists[0].Playlist.Id)
	assert.Equal(t, usecaseResponse.Playlists[0].IsIncluded, protoResponse.Playlists[0].IsIncludedTrack)
}

func TestUpdatePlaylistsPublisityByUserIDRequestConversions(t *testing.T) {
	protoReq := &protoModel.UpdatePlaylistsPublisityByUserIDRequest{
		UserId:   123,
		IsPublic: true,
	}

	usecaseReq := UpdatePlaylistsPublisityByUserIDRequestFromProtoToUsecase(protoReq)
	assert.Equal(t, protoReq.GetUserId(), usecaseReq.UserID)
	assert.Equal(t, protoReq.GetIsPublic(), usecaseReq.IsPublic)

	repoReq := UpdatePlaylistsPublisityByUserIDRequestFromUsecaseToRepository(usecaseReq)
	assert.Equal(t, usecaseReq.UserID, repoReq.UserID)
	assert.Equal(t, usecaseReq.IsPublic, repoReq.IsPublic)
}

func TestLikePlaylistRequestConversions(t *testing.T) {
	protoReq := &protoModel.LikePlaylistRequest{
		UserId:     123,
		PlaylistId: 456,
		IsLike:     true,
	}

	usecaseReq := LikePlaylistRequestFromProtoToUsecase(protoReq)
	assert.Equal(t, protoReq.GetUserId(), usecaseReq.UserID)
	assert.Equal(t, protoReq.GetPlaylistId(), usecaseReq.PlaylistID)
	assert.Equal(t, protoReq.GetIsLike(), usecaseReq.IsLike)

	repoReq := LikePlaylistRequestFromUsecaseToRepository(usecaseReq)
	assert.Equal(t, usecaseReq.UserID, repoReq.UserID)
	assert.Equal(t, usecaseReq.PlaylistID, repoReq.PlaylistID)
}

func TestPlaylistWithIsLikedConversions(t *testing.T) {
	repoPlaylist := &repoModel.Playlist{
		ID:        1,
		Title:     "Test Playlist",
		Thumbnail: "thumbnail.jpg",
		UserID:    123,
	}

	repoPlaylistWithLiked := &repoModel.PlaylistWithIsLiked{
		Playlist: repoPlaylist,
		IsLiked:  true,
	}

	usecasePlaylistWithLiked := PlaylistWithIsLikedFromRepositoryToUsecase(repoPlaylistWithLiked)
	assert.Equal(t, repoPlaylistWithLiked.Playlist.ID, usecasePlaylistWithLiked.Playlist.ID)
	assert.Equal(t, repoPlaylistWithLiked.Playlist.Title, usecasePlaylistWithLiked.Playlist.Title)
	assert.Equal(t, repoPlaylistWithLiked.Playlist.Thumbnail, usecasePlaylistWithLiked.Playlist.Thumbnail)
	assert.Equal(t, repoPlaylistWithLiked.Playlist.UserID, usecasePlaylistWithLiked.Playlist.UserID)
	assert.Equal(t, repoPlaylistWithLiked.IsLiked, usecasePlaylistWithLiked.IsLiked)

	protoPlaylistWithLiked := PlaylistWithIsLikedFromUsecaseToProto(usecasePlaylistWithLiked)
	assert.Equal(t, usecasePlaylistWithLiked.Playlist.ID, protoPlaylistWithLiked.Playlist.Id)
	assert.Equal(t, usecasePlaylistWithLiked.Playlist.Title, protoPlaylistWithLiked.Playlist.Title)
	assert.Equal(t, usecasePlaylistWithLiked.Playlist.Thumbnail, protoPlaylistWithLiked.Playlist.Thumbnail)
	assert.Equal(t, usecasePlaylistWithLiked.Playlist.UserID, protoPlaylistWithLiked.Playlist.UserId)
	assert.Equal(t, usecasePlaylistWithLiked.IsLiked, protoPlaylistWithLiked.IsLiked)
}

func TestGetProfilePlaylistsRequestConversions(t *testing.T) {
	protoReq := &protoModel.GetProfilePlaylistsRequest{
		UserId: 123,
	}

	usecaseReq := GetProfilePlaylistsRequestFromProtoToUsecase(protoReq)
	assert.Equal(t, protoReq.GetUserId(), usecaseReq.UserID)

	repoReq := GetProfilePlaylistsRequestFromUsecaseToRepository(usecaseReq)
	assert.Equal(t, usecaseReq.UserID, repoReq.UserID)
}

func TestGetProfilePlaylistsResponseConversions(t *testing.T) {
	repoPlaylists := []*repoModel.Playlist{
		{
			ID:        1,
			Title:     "Test Playlist 1",
			Thumbnail: "thumbnail1.jpg",
			UserID:    123,
		},
		{
			ID:        2,
			Title:     "Test Playlist 2",
			Thumbnail: "thumbnail2.jpg",
			UserID:    123,
		},
	}

	repoResponse := &repoModel.GetProfilePlaylistsResponse{
		Playlists: repoPlaylists,
	}

	usecaseResponse := GetProfilePlaylistsResponseFromRepositoryToUsecase(repoResponse)
	assert.Equal(t, len(repoResponse.Playlists), len(usecaseResponse.Playlists))
	for i, playlist := range usecaseResponse.Playlists {
		assert.Equal(t, repoResponse.Playlists[i].ID, playlist.ID)
		assert.Equal(t, repoResponse.Playlists[i].Title, playlist.Title)
		assert.Equal(t, repoResponse.Playlists[i].Thumbnail, playlist.Thumbnail)
		assert.Equal(t, repoResponse.Playlists[i].UserID, playlist.UserID)
	}

	protoResponse := GetProfilePlaylistsResponseFromUsecaseToProto(usecaseResponse)
	assert.Equal(t, len(usecaseResponse.Playlists), len(protoResponse.Playlists))
	for i, playlist := range protoResponse.Playlists {
		assert.Equal(t, usecaseResponse.Playlists[i].ID, playlist.Id)
		assert.Equal(t, usecaseResponse.Playlists[i].Title, playlist.Title)
		assert.Equal(t, usecaseResponse.Playlists[i].Thumbnail, playlist.Thumbnail)
		assert.Equal(t, usecaseResponse.Playlists[i].UserID, playlist.UserId)
	}
}

func TestSearchPlaylistsRequestConversions(t *testing.T) {
	protoReq := &protoModel.SearchPlaylistsRequest{
		UserId: 123,
		Query:  "test search",
	}

	usecaseReq := SearchPlaylistsRequestFromProtoToUsecase(protoReq)
	assert.Equal(t, protoReq.GetUserId(), usecaseReq.UserID)
	assert.Equal(t, protoReq.GetQuery(), usecaseReq.Query)

	repoReq := SearchPlaylistsRequestFromUsecaseToRepository(usecaseReq)
	assert.Equal(t, usecaseReq.UserID, repoReq.UserID)
	assert.Equal(t, usecaseReq.Query, repoReq.Query)
}

func TestGetCombinedPlaylistsByUserIDRequestConversion(t *testing.T) {
	protoReq := &protoModel.GetCombinedPlaylistsByUserIDRequest{
		UserId: 123,
	}

	usecaseReq := GetCombinedPlaylistsByUserIDRequestFromProtoToUsecase(protoReq)
	assert.Equal(t, protoReq.GetUserId(), usecaseReq.UserID)
}
