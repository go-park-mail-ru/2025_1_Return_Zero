package model

import (
	"testing"

	"github.com/stretchr/testify/assert"

	protoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/artist"
	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/artist/model/repository"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/artist/model/usecase"
)

func TestArtistFromRepositoryToUsecase(t *testing.T) {
	repoArtist := &repoModel.Artist{
		ID:          1,
		Title:       "Artist Title",
		Description: "Artist Description",
		Thumbnail:   "thumbnail.jpg",
		IsFavorite:  true,
	}

	usecaseArtist := ArtistFromRepositoryToUsecase(repoArtist)

	assert.Equal(t, repoArtist.ID, usecaseArtist.ID)
	assert.Equal(t, repoArtist.Title, usecaseArtist.Title)
	assert.Equal(t, repoArtist.Description, usecaseArtist.Description)
	assert.Equal(t, repoArtist.Thumbnail, usecaseArtist.Thumbnail)
	assert.Equal(t, repoArtist.IsFavorite, usecaseArtist.IsFavorite)
}

func TestArtistFromUsecaseToProto(t *testing.T) {
	usecaseArtist := &usecaseModel.Artist{
		ID:          1,
		Title:       "Artist Title",
		Description: "Artist Description",
		Thumbnail:   "thumbnail.jpg",
		IsFavorite:  true,
	}

	protoArtist := ArtistFromUsecaseToProto(usecaseArtist)

	assert.Equal(t, usecaseArtist.ID, protoArtist.Id)
	assert.Equal(t, usecaseArtist.Title, protoArtist.Title)
	assert.Equal(t, usecaseArtist.Description, protoArtist.Description)
	assert.Equal(t, usecaseArtist.Thumbnail, protoArtist.Thumbnail)
	assert.Equal(t, usecaseArtist.IsFavorite, protoArtist.IsFavorite)
}

func TestArtistDetailedFromRepositoryToUsecase(t *testing.T) {
	repoArtist := &repoModel.Artist{
		ID:          1,
		Title:       "Artist Title",
		Description: "Artist Description",
		Thumbnail:   "thumbnail.jpg",
		IsFavorite:  true,
	}

	stats := &repoModel.ArtistStats{
		ListenersCount: 1000,
		FavoritesCount: 500,
	}

	usecaseArtistDetailed := ArtistDetailedFromRepositoryToUsecase(repoArtist, stats)

	assert.Equal(t, repoArtist.ID, usecaseArtistDetailed.Artist.ID)
	assert.Equal(t, repoArtist.Title, usecaseArtistDetailed.Artist.Title)
	assert.Equal(t, repoArtist.Description, usecaseArtistDetailed.Artist.Description)
	assert.Equal(t, repoArtist.Thumbnail, usecaseArtistDetailed.Artist.Thumbnail)
	assert.Equal(t, repoArtist.IsFavorite, usecaseArtistDetailed.Artist.IsFavorite)
	assert.Equal(t, stats.ListenersCount, usecaseArtistDetailed.ListenersCount)
	assert.Equal(t, stats.FavoritesCount, usecaseArtistDetailed.FavoritesCount)
}

func TestArtistDetailedFromUsecaseToProto(t *testing.T) {
	usecaseArtist := &usecaseModel.Artist{
		ID:          1,
		Title:       "Artist Title",
		Description: "Artist Description",
		Thumbnail:   "thumbnail.jpg",
		IsFavorite:  true,
	}

	usecaseArtistDetailed := &usecaseModel.ArtistDetailed{
		Artist:         usecaseArtist,
		ListenersCount: 1000,
		FavoritesCount: 500,
	}

	protoArtistDetailed := ArtistDetailedFromUsecaseToProto(usecaseArtistDetailed)

	assert.Equal(t, usecaseArtistDetailed.Artist.ID, protoArtistDetailed.Artist.Id)
	assert.Equal(t, usecaseArtistDetailed.Artist.Title, protoArtistDetailed.Artist.Title)
	assert.Equal(t, usecaseArtistDetailed.Artist.Description, protoArtistDetailed.Artist.Description)
	assert.Equal(t, usecaseArtistDetailed.Artist.Thumbnail, protoArtistDetailed.Artist.Thumbnail)
	assert.Equal(t, usecaseArtistDetailed.Artist.IsFavorite, protoArtistDetailed.Artist.IsFavorite)
	assert.Equal(t, usecaseArtistDetailed.ListenersCount, protoArtistDetailed.ListenersCount)
	assert.Equal(t, usecaseArtistDetailed.FavoritesCount, protoArtistDetailed.FavoritesCount)
}

func TestArtistListFromRepositoryToUsecase(t *testing.T) {
	repoArtists := []*repoModel.Artist{
		{
			ID:          1,
			Title:       "Artist 1",
			Description: "Description 1",
			Thumbnail:   "thumbnail1.jpg",
			IsFavorite:  true,
		},
		{
			ID:          2,
			Title:       "Artist 2",
			Description: "Description 2",
			Thumbnail:   "thumbnail2.jpg",
			IsFavorite:  false,
		},
	}

	usecaseArtistList := ArtistListFromRepositoryToUsecase(repoArtists)

	assert.Equal(t, len(repoArtists), len(usecaseArtistList.Artists))

	for i, artist := range usecaseArtistList.Artists {
		assert.Equal(t, repoArtists[i].ID, artist.ID)
		assert.Equal(t, repoArtists[i].Title, artist.Title)
		assert.Equal(t, repoArtists[i].Description, artist.Description)
		assert.Equal(t, repoArtists[i].Thumbnail, artist.Thumbnail)
		assert.Equal(t, repoArtists[i].IsFavorite, artist.IsFavorite)
	}
}

func TestArtistListFromUsecaseToProto(t *testing.T) {
	usecaseArtists := &usecaseModel.ArtistList{
		Artists: []*usecaseModel.Artist{
			{
				ID:          1,
				Title:       "Artist 1",
				Description: "Description 1",
				Thumbnail:   "thumbnail1.jpg",
				IsFavorite:  true,
			},
			{
				ID:          2,
				Title:       "Artist 2",
				Description: "Description 2",
				Thumbnail:   "thumbnail2.jpg",
				IsFavorite:  false,
			},
		},
	}

	protoArtistList := ArtistListFromUsecaseToProto(usecaseArtists)

	assert.Equal(t, len(usecaseArtists.Artists), len(protoArtistList.Artists))

	for i, artist := range protoArtistList.Artists {
		assert.Equal(t, usecaseArtists.Artists[i].ID, artist.Id)
		assert.Equal(t, usecaseArtists.Artists[i].Title, artist.Title)
		assert.Equal(t, usecaseArtists.Artists[i].Description, artist.Description)
		assert.Equal(t, usecaseArtists.Artists[i].Thumbnail, artist.Thumbnail)
		assert.Equal(t, usecaseArtists.Artists[i].IsFavorite, artist.IsFavorite)
	}
}

func TestArtistWithTitleFromRepositoryToUsecase(t *testing.T) {
	repoArtistWithTitle := &repoModel.ArtistWithTitle{
		ID:    1,
		Title: "Artist Title",
	}

	usecaseArtistWithTitle := ArtistWithTitleFromRepositoryToUsecase(repoArtistWithTitle)

	assert.Equal(t, repoArtistWithTitle.ID, usecaseArtistWithTitle.ID)
	assert.Equal(t, repoArtistWithTitle.Title, usecaseArtistWithTitle.Title)
}

func TestArtistWithTitleFromUsecaseToProto(t *testing.T) {
	usecaseArtistWithTitle := &usecaseModel.ArtistWithTitle{
		ID:    1,
		Title: "Artist Title",
	}

	protoArtistWithTitle := ArtistWithTitleFromUsecaseToProto(usecaseArtistWithTitle)

	assert.Equal(t, usecaseArtistWithTitle.ID, protoArtistWithTitle.Id)
	assert.Equal(t, usecaseArtistWithTitle.Title, protoArtistWithTitle.Title)
}

func TestArtistWithTitleListFromUsecaseToProto(t *testing.T) {
	usecaseArtistWithTitleList := &usecaseModel.ArtistWithTitleList{
		Artists: []*usecaseModel.ArtistWithTitle{
			{
				ID:    1,
				Title: "Artist 1",
			},
			{
				ID:    2,
				Title: "Artist 2",
			},
		},
	}

	protoArtistWithTitleList := ArtistWithTitleListFromUsecaseToProto(usecaseArtistWithTitleList)

	assert.Equal(t, len(usecaseArtistWithTitleList.Artists), len(protoArtistWithTitleList.Artists))

	for i, artist := range protoArtistWithTitleList.Artists {
		assert.Equal(t, usecaseArtistWithTitleList.Artists[i].ID, artist.Id)
		assert.Equal(t, usecaseArtistWithTitleList.Artists[i].Title, artist.Title)
	}
}

func TestArtistWithTitleListFromRepositoryToUsecase(t *testing.T) {
	repoArtistWithTitleList := []*repoModel.ArtistWithTitle{
		{
			ID:    1,
			Title: "Artist 1",
		},
		{
			ID:    2,
			Title: "Artist 2",
		},
	}

	usecaseArtistWithTitleList := ArtistWithTitleListFromRepositoryToUsecase(repoArtistWithTitleList)

	assert.Equal(t, len(repoArtistWithTitleList), len(usecaseArtistWithTitleList.Artists))

	for i, artist := range usecaseArtistWithTitleList.Artists {
		assert.Equal(t, repoArtistWithTitleList[i].ID, artist.ID)
		assert.Equal(t, repoArtistWithTitleList[i].Title, artist.Title)
	}
}

func TestArtistWithTitleMapFromRepositoryToUsecase(t *testing.T) {
	repoArtistWithTitleMap := map[int64][]*repoModel.ArtistWithTitle{
		1: {
			{
				ID:    10,
				Title: "Artist 10",
			},
			{
				ID:    11,
				Title: "Artist 11",
			},
		},
		2: {
			{
				ID:    20,
				Title: "Artist 20",
			},
		},
	}

	usecaseArtistWithTitleMap := ArtistWithTitleMapFromRepositoryToUsecase(repoArtistWithTitleMap)

	assert.Equal(t, len(repoArtistWithTitleMap), len(usecaseArtistWithTitleMap.Artists))

	for id, artists := range usecaseArtistWithTitleMap.Artists {
		assert.Equal(t, len(repoArtistWithTitleMap[id]), len(artists.Artists))

		for i, artist := range artists.Artists {
			assert.Equal(t, repoArtistWithTitleMap[id][i].ID, artist.ID)
			assert.Equal(t, repoArtistWithTitleMap[id][i].Title, artist.Title)
		}
	}
}

func TestArtistWithTitleMapFromUsecaseToProto(t *testing.T) {
	usecaseArtistWithTitleMap := &usecaseModel.ArtistWithTitleMap{
		Artists: map[int64]*usecaseModel.ArtistWithTitleList{
			1: {
				Artists: []*usecaseModel.ArtistWithTitle{
					{
						ID:    10,
						Title: "Artist 10",
					},
					{
						ID:    11,
						Title: "Artist 11",
					},
				},
			},
			2: {
				Artists: []*usecaseModel.ArtistWithTitle{
					{
						ID:    20,
						Title: "Artist 20",
					},
				},
			},
		},
	}

	protoArtistWithTitleMap := ArtistWithTitleMapFromUsecaseToProto(usecaseArtistWithTitleMap)

	assert.Equal(t, len(usecaseArtistWithTitleMap.Artists), len(protoArtistWithTitleMap.Artists))

	for id, artists := range protoArtistWithTitleMap.Artists {
		assert.Equal(t, len(usecaseArtistWithTitleMap.Artists[id].Artists), len(artists.Artists))

		for i, artist := range artists.Artists {
			assert.Equal(t, usecaseArtistWithTitleMap.Artists[id].Artists[i].ID, artist.Id)
			assert.Equal(t, usecaseArtistWithTitleMap.Artists[id].Artists[i].Title, artist.Title)
		}
	}
}

func TestArtistWithRoleFromRepositoryToUsecase(t *testing.T) {
	repoArtistWithRole := &repoModel.ArtistWithRole{
		ID:    1,
		Title: "Artist Title",
		Role:  "main",
	}

	usecaseArtistWithRole := ArtistWithRoleFromRepositoryToUsecase(repoArtistWithRole)

	assert.Equal(t, repoArtistWithRole.ID, usecaseArtistWithRole.ID)
	assert.Equal(t, repoArtistWithRole.Title, usecaseArtistWithRole.Title)
	assert.Equal(t, repoArtistWithRole.Role, usecaseArtistWithRole.Role)
}

func TestArtistWithRoleFromUsecaseToProto(t *testing.T) {
	usecaseArtistWithRole := &usecaseModel.ArtistWithRole{
		ID:    1,
		Title: "Artist Title",
		Role:  "main",
	}

	protoArtistWithRole := ArtistWithRoleFromUsecaseToProto(usecaseArtistWithRole)

	assert.Equal(t, usecaseArtistWithRole.ID, protoArtistWithRole.Id)
	assert.Equal(t, usecaseArtistWithRole.Title, protoArtistWithRole.Title)
	assert.Equal(t, usecaseArtistWithRole.Role, protoArtistWithRole.Role)
}

func TestArtistWithRoleListFromUsecaseToProto(t *testing.T) {
	usecaseArtistWithRoleList := &usecaseModel.ArtistWithRoleList{
		Artists: []*usecaseModel.ArtistWithRole{
			{
				ID:    1,
				Title: "Artist 1",
				Role:  "main",
			},
			{
				ID:    2,
				Title: "Artist 2",
				Role:  "featured",
			},
		},
	}

	protoArtistWithRoleList := ArtistWithRoleListFromUsecaseToProto(usecaseArtistWithRoleList)

	assert.Equal(t, len(usecaseArtistWithRoleList.Artists), len(protoArtistWithRoleList.Artists))

	for i, artist := range protoArtistWithRoleList.Artists {
		assert.Equal(t, usecaseArtistWithRoleList.Artists[i].ID, artist.Id)
		assert.Equal(t, usecaseArtistWithRoleList.Artists[i].Title, artist.Title)
		assert.Equal(t, usecaseArtistWithRoleList.Artists[i].Role, artist.Role)
	}
}

func TestArtistWithRoleMapFromRepositoryToUsecase(t *testing.T) {
	repoArtistWithRoleMap := map[int64][]*repoModel.ArtistWithRole{
		1: {
			{
				ID:    10,
				Title: "Artist 10",
				Role:  "main",
			},
			{
				ID:    11,
				Title: "Artist 11",
				Role:  "featured",
			},
		},
		2: {
			{
				ID:    20,
				Title: "Artist 20",
				Role:  "main",
			},
		},
	}

	usecaseArtistWithRoleMap := ArtistWithRoleMapFromRepositoryToUsecase(repoArtistWithRoleMap)

	assert.Equal(t, len(repoArtistWithRoleMap), len(usecaseArtistWithRoleMap.Artists))

	for id, artists := range usecaseArtistWithRoleMap.Artists {
		assert.Equal(t, len(repoArtistWithRoleMap[id]), len(artists.Artists))

		for i, artist := range artists.Artists {
			assert.Equal(t, repoArtistWithRoleMap[id][i].ID, artist.ID)
			assert.Equal(t, repoArtistWithRoleMap[id][i].Title, artist.Title)
			assert.Equal(t, repoArtistWithRoleMap[id][i].Role, artist.Role)
		}
	}
}

func TestArtistWithRoleMapFromUsecaseToProto(t *testing.T) {
	usecaseArtistWithRoleMap := &usecaseModel.ArtistWithRoleMap{
		Artists: map[int64]*usecaseModel.ArtistWithRoleList{
			1: {
				Artists: []*usecaseModel.ArtistWithRole{
					{
						ID:    10,
						Title: "Artist 10",
						Role:  "main",
					},
					{
						ID:    11,
						Title: "Artist 11",
						Role:  "featured",
					},
				},
			},
			2: {
				Artists: []*usecaseModel.ArtistWithRole{
					{
						ID:    20,
						Title: "Artist 20",
						Role:  "main",
					},
				},
			},
		},
	}

	protoArtistWithRoleMap := ArtistWithRoleMapFromUsecaseToProto(usecaseArtistWithRoleMap)

	assert.Equal(t, len(usecaseArtistWithRoleMap.Artists), len(protoArtistWithRoleMap.Artists))

	for id, artists := range protoArtistWithRoleMap.Artists {
		assert.Equal(t, len(usecaseArtistWithRoleMap.Artists[id].Artists), len(artists.Artists))

		for i, artist := range artists.Artists {
			assert.Equal(t, usecaseArtistWithRoleMap.Artists[id].Artists[i].ID, artist.Id)
			assert.Equal(t, usecaseArtistWithRoleMap.Artists[id].Artists[i].Title, artist.Title)
			assert.Equal(t, usecaseArtistWithRoleMap.Artists[id].Artists[i].Role, artist.Role)
		}
	}
}

func TestArtistWithRoleListFromRepositoryToUsecase(t *testing.T) {
	repoArtistWithRoleList := []*repoModel.ArtistWithRole{
		{
			ID:    1,
			Title: "Artist 1",
			Role:  "main",
		},
		{
			ID:    2,
			Title: "Artist 2",
			Role:  "featured",
		},
	}

	usecaseArtistWithRoleList := ArtistWithRoleListFromRepositoryToUsecase(repoArtistWithRoleList)

	assert.Equal(t, len(repoArtistWithRoleList), len(usecaseArtistWithRoleList.Artists))

	for i, artist := range usecaseArtistWithRoleList.Artists {
		assert.Equal(t, repoArtistWithRoleList[i].ID, artist.ID)
		assert.Equal(t, repoArtistWithRoleList[i].Title, artist.Title)
		assert.Equal(t, repoArtistWithRoleList[i].Role, artist.Role)
	}
}

func TestPaginationFromUsecaseToRepository(t *testing.T) {
	usecasePagination := &usecaseModel.Pagination{
		Offset: 10,
		Limit:  20,
	}

	repoPagination := PaginationFromUsecaseToRepository(usecasePagination)

	assert.Equal(t, usecasePagination.Offset, repoPagination.Offset)
	assert.Equal(t, usecasePagination.Limit, repoPagination.Limit)
}

func TestArtistFiltersFromUsecaseToRepository(t *testing.T) {
	usecasePagination := &usecaseModel.Pagination{
		Offset: 10,
		Limit:  20,
	}

	usecaseFilters := &usecaseModel.Filters{
		Pagination: usecasePagination,
	}

	repoFilters := ArtistFiltersFromUsecaseToRepository(usecaseFilters)

	assert.Equal(t, usecaseFilters.Pagination.Offset, repoFilters.Pagination.Offset)
	assert.Equal(t, usecaseFilters.Pagination.Limit, repoFilters.Pagination.Limit)
}

func TestTrackIDListFromProtoToUsecase(t *testing.T) {
	protoTrackIDs := []*protoModel.TrackID{
		{Id: 1},
		{Id: 2},
		{Id: 3},
	}

	trackIDs := TrackIDListFromProtoToUsecase(protoTrackIDs)

	assert.Equal(t, len(protoTrackIDs), len(trackIDs))

	for i, id := range trackIDs {
		assert.Equal(t, protoTrackIDs[i].Id, id)
	}
}

func TestAlbumIDListFromProtoToUsecase(t *testing.T) {
	protoAlbumIDs := []*protoModel.AlbumID{
		{Id: 1},
		{Id: 2},
		{Id: 3},
	}

	albumIDs := AlbumIDListFromProtoToUsecase(protoAlbumIDs)

	assert.Equal(t, len(protoAlbumIDs), len(albumIDs))

	for i, id := range albumIDs {
		assert.Equal(t, protoAlbumIDs[i].Id, id)
	}
}

func TestPaginationFromProtoToUsecase(t *testing.T) {
	protoPagination := &protoModel.Pagination{
		Offset: 10,
		Limit:  20,
	}

	usecasePagination := PaginationFromProtoToUsecase(protoPagination)

	assert.Equal(t, protoPagination.Offset, usecasePagination.Offset)
	assert.Equal(t, protoPagination.Limit, usecasePagination.Limit)
}

func TestArtistFiltersFromProtoToUsecase(t *testing.T) {
	protoPagination := &protoModel.Pagination{
		Offset: 10,
		Limit:  20,
	}

	protoFilters := &protoModel.Filters{
		Pagination: protoPagination,
	}

	usecaseFilters := ArtistFiltersFromProtoToUsecase(protoFilters)

	assert.Equal(t, protoFilters.Pagination.Offset, usecaseFilters.Pagination.Offset)
	assert.Equal(t, protoFilters.Pagination.Limit, usecaseFilters.Pagination.Limit)
}

func TestArtistStreamCreateDataFromProtoToUsecase(t *testing.T) {
	protoStreamData := &protoModel.ArtistStreamCreateDataList{
		ArtistIds: &protoModel.ArtistIDList{
			Ids: []*protoModel.ArtistID{
				{Id: 1},
				{Id: 2},
			},
		},
		UserId: &protoModel.UserID{Id: 10},
	}

	usecaseStreamData := ArtistStreamCreateDataFromProtoToUsecase(protoStreamData)

	assert.Equal(t, len(protoStreamData.ArtistIds.Ids), len(usecaseStreamData.ArtistIDs))

	for i, id := range usecaseStreamData.ArtistIDs {
		assert.Equal(t, protoStreamData.ArtistIds.Ids[i].Id, id)
	}

	assert.Equal(t, protoStreamData.UserId.Id, usecaseStreamData.UserID)
}

func TestArtistStreamCreateDataFromUsecaseToRepository(t *testing.T) {
	usecaseStreamData := &usecaseModel.ArtistStreamCreateDataList{
		ArtistIDs: []int64{1, 2},
		UserID:    10,
	}

	repoStreamData := ArtistStreamCreateDataFromUsecaseToRepository(usecaseStreamData)

	assert.Equal(t, usecaseStreamData.ArtistIDs, repoStreamData.ArtistIDs)
	assert.Equal(t, usecaseStreamData.UserID, repoStreamData.UserID)
}

func TestLikeRequestFromProtoToUsecase(t *testing.T) {
	protoLikeRequest := &protoModel.LikeRequest{
		ArtistId: &protoModel.ArtistID{Id: 1},
		UserId:   &protoModel.UserID{Id: 2},
		IsLike:   true,
	}

	usecaseLikeRequest := LikeRequestFromProtoToUsecase(protoLikeRequest)

	assert.Equal(t, protoLikeRequest.ArtistId.Id, usecaseLikeRequest.ArtistID)
	assert.Equal(t, protoLikeRequest.UserId.Id, usecaseLikeRequest.UserID)
	assert.Equal(t, protoLikeRequest.IsLike, usecaseLikeRequest.IsLike)
}

func TestLikeRequestFromUsecaseToRepository(t *testing.T) {
	usecaseLikeRequest := &usecaseModel.LikeRequest{
		ArtistID: 1,
		UserID:   2,
		IsLike:   true,
	}

	repoLikeRequest := LikeRequestFromUsecaseToRepository(usecaseLikeRequest)

	assert.Equal(t, usecaseLikeRequest.ArtistID, repoLikeRequest.ArtistID)
	assert.Equal(t, usecaseLikeRequest.UserID, repoLikeRequest.UserID)
}
