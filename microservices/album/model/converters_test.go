package model

import (
	"testing"
	"time"

	albumProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/album"
	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/album/model/repository"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/album/model/usecase"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestPaginationFromUsecaseToRepository(t *testing.T) {
	usecasePagination := &usecaseModel.Pagination{
		Offset: 10,
		Limit:  20,
	}

	repoPagination := PaginationFromUsecaseToRepository(usecasePagination)

	assert.Equal(t, usecasePagination.Offset, repoPagination.Offset)
	assert.Equal(t, usecasePagination.Limit, repoPagination.Limit)
}

func TestFiltersFromUsecaseToRepository(t *testing.T) {
	usecaseFilters := &usecaseModel.AlbumFilters{
		Pagination: &usecaseModel.Pagination{
			Offset: 10,
			Limit:  20,
		},
	}

	repoFilters := FiltersFromUsecaseToRepository(usecaseFilters)

	assert.Equal(t, usecaseFilters.Pagination.Offset, repoFilters.Pagination.Offset)
	assert.Equal(t, usecaseFilters.Pagination.Limit, repoFilters.Pagination.Limit)
}

func TestPaginationFromProtoToUsecase(t *testing.T) {
	protoPagination := &albumProto.Pagination{
		Offset: 10,
		Limit:  20,
	}

	usecasePagination := PaginationFromProtoToUsecase(protoPagination)

	assert.Equal(t, protoPagination.Offset, usecasePagination.Offset)
	assert.Equal(t, protoPagination.Limit, usecasePagination.Limit)
}

func TestAlbumFiltersFromProtoToUsecase(t *testing.T) {
	protoFilters := &albumProto.Filters{
		Pagination: &albumProto.Pagination{
			Offset: 10,
			Limit:  20,
		},
	}

	usecaseFilters := AlbumFiltersFromProtoToUsecase(protoFilters)

	assert.Equal(t, protoFilters.Pagination.Offset, usecaseFilters.Pagination.Offset)
	assert.Equal(t, protoFilters.Pagination.Limit, usecaseFilters.Pagination.Limit)
}

func TestAlbumTypeFromRepositoryToUsecase(t *testing.T) {
	testCases := []struct {
		name            string
		repoAlbumType   repoModel.AlbumType
		expectAlbumType usecaseModel.AlbumType
	}{
		{
			name:            "AlbumTypeAlbum",
			repoAlbumType:   repoModel.AlbumTypeAlbum,
			expectAlbumType: usecaseModel.AlbumTypeAlbum,
		},
		{
			name:            "AlbumTypeEP",
			repoAlbumType:   repoModel.AlbumTypeEP,
			expectAlbumType: usecaseModel.AlbumTypeEP,
		},
		{
			name:            "AlbumTypeSingle",
			repoAlbumType:   repoModel.AlbumTypeSingle,
			expectAlbumType: usecaseModel.AlbumTypeSingle,
		},
		{
			name:            "AlbumTypeCompilation",
			repoAlbumType:   repoModel.AlbumTypeCompilation,
			expectAlbumType: usecaseModel.AlbumTypeCompilation,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			usecaseAlbumType := AlbumTypeFromRepositoryToUsecase(tc.repoAlbumType)
			assert.Equal(t, tc.expectAlbumType, usecaseAlbumType)
		})
	}
}

func TestAlbumTypeFromUsecaseToProto(t *testing.T) {
	testCases := []struct {
		name             string
		usecaseAlbumType usecaseModel.AlbumType
		expectAlbumType  albumProto.AlbumType
	}{
		{
			name:             "AlbumTypeAlbum",
			usecaseAlbumType: usecaseModel.AlbumTypeAlbum,
			expectAlbumType:  albumProto.AlbumType_AlbumTypeAlbum,
		},
		{
			name:             "AlbumTypeEP",
			usecaseAlbumType: usecaseModel.AlbumTypeEP,
			expectAlbumType:  albumProto.AlbumType_AlbumTypeEP,
		},
		{
			name:             "AlbumTypeSingle",
			usecaseAlbumType: usecaseModel.AlbumTypeSingle,
			expectAlbumType:  albumProto.AlbumType_AlbumTypeSingle,
		},
		{
			name:             "AlbumTypeCompilation",
			usecaseAlbumType: usecaseModel.AlbumTypeCompilation,
			expectAlbumType:  albumProto.AlbumType_AlbumTypeCompilation,
		},
		{
			name:             "Unknown type defaults to AlbumTypeAlbum",
			usecaseAlbumType: "unknown",
			expectAlbumType:  albumProto.AlbumType_AlbumTypeAlbum,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			protoAlbumType := AlbumTypeFromUsecaseToProto(tc.usecaseAlbumType)
			assert.Equal(t, tc.expectAlbumType, protoAlbumType)
		})
	}
}

func TestAlbumFromRepositoryToUsecase(t *testing.T) {
	releaseDate := time.Now()
	repoAlbum := &repoModel.Album{
		ID:          1,
		Title:       "Test Album",
		Type:        repoModel.AlbumTypeAlbum,
		Thumbnail:   "thumbnail.jpg",
		ReleaseDate: releaseDate,
		IsFavorite:  true,
	}

	usecaseAlbum := AlbumFromRepositoryToUsecase(repoAlbum)

	assert.Equal(t, repoAlbum.ID, usecaseAlbum.ID)
	assert.Equal(t, repoAlbum.Title, usecaseAlbum.Title)
	assert.Equal(t, usecaseModel.AlbumType(repoAlbum.Type), usecaseAlbum.Type)
	assert.Equal(t, repoAlbum.Thumbnail, usecaseAlbum.Thumbnail)
	assert.Equal(t, repoAlbum.ReleaseDate, usecaseAlbum.ReleaseDate)
	assert.Equal(t, repoAlbum.IsFavorite, usecaseAlbum.IsFavorite)
}

func TestAlbumFromUsecaseToProto(t *testing.T) {
	releaseDate := time.Now()
	usecaseAlbum := &usecaseModel.Album{
		ID:          1,
		Title:       "Test Album",
		Type:        usecaseModel.AlbumTypeAlbum,
		Thumbnail:   "thumbnail.jpg",
		ReleaseDate: releaseDate,
		IsFavorite:  true,
	}

	protoAlbum := AlbumFromUsecaseToProto(usecaseAlbum)

	assert.Equal(t, usecaseAlbum.ID, protoAlbum.Id)
	assert.Equal(t, usecaseAlbum.Title, protoAlbum.Title)
	assert.Equal(t, albumProto.AlbumType_AlbumTypeAlbum, protoAlbum.Type)
	assert.Equal(t, usecaseAlbum.Thumbnail, protoAlbum.Thumbnail)
	assert.Equal(t, timestamppb.New(releaseDate).Seconds, protoAlbum.ReleaseDate.Seconds)
	assert.Equal(t, timestamppb.New(releaseDate).Nanos, protoAlbum.ReleaseDate.Nanos)
	assert.Equal(t, usecaseAlbum.IsFavorite, protoAlbum.IsFavorite)
}

func TestAlbumListFromRepositoryToUsecase(t *testing.T) {
	releaseDate := time.Now()
	repoAlbums := []*repoModel.Album{
		{
			ID:          1,
			Title:       "Test Album 1",
			Type:        repoModel.AlbumTypeAlbum,
			Thumbnail:   "thumbnail1.jpg",
			ReleaseDate: releaseDate,
			IsFavorite:  true,
		},
		{
			ID:          2,
			Title:       "Test Album 2",
			Type:        repoModel.AlbumTypeEP,
			Thumbnail:   "thumbnail2.jpg",
			ReleaseDate: releaseDate,
			IsFavorite:  false,
		},
	}

	usecaseAlbums := AlbumListFromRepositoryToUsecase(repoAlbums)

	assert.Equal(t, len(repoAlbums), len(usecaseAlbums))

	for i, repoAlbum := range repoAlbums {
		assert.Equal(t, repoAlbum.ID, usecaseAlbums[i].ID)
		assert.Equal(t, repoAlbum.Title, usecaseAlbums[i].Title)
		assert.Equal(t, usecaseModel.AlbumType(repoAlbum.Type), usecaseAlbums[i].Type)
		assert.Equal(t, repoAlbum.Thumbnail, usecaseAlbums[i].Thumbnail)
		assert.Equal(t, repoAlbum.ReleaseDate, usecaseAlbums[i].ReleaseDate)
		assert.Equal(t, repoAlbum.IsFavorite, usecaseAlbums[i].IsFavorite)
	}
}

func TestAlbumListFromUsecaseToProto(t *testing.T) {
	releaseDate := time.Now()
	usecaseAlbums := []*usecaseModel.Album{
		{
			ID:          1,
			Title:       "Test Album 1",
			Type:        usecaseModel.AlbumTypeAlbum,
			Thumbnail:   "thumbnail1.jpg",
			ReleaseDate: releaseDate,
			IsFavorite:  true,
		},
		{
			ID:          2,
			Title:       "Test Album 2",
			Type:        usecaseModel.AlbumTypeEP,
			Thumbnail:   "thumbnail2.jpg",
			ReleaseDate: releaseDate,
			IsFavorite:  false,
		},
	}

	protoAlbumList := AlbumListFromUsecaseToProto(usecaseAlbums)

	assert.Equal(t, len(usecaseAlbums), len(protoAlbumList.Albums))

	for i, usecaseAlbum := range usecaseAlbums {
		assert.Equal(t, usecaseAlbum.ID, protoAlbumList.Albums[i].Id)
		assert.Equal(t, usecaseAlbum.Title, protoAlbumList.Albums[i].Title)
		assert.Equal(t, AlbumTypeFromUsecaseToProto(usecaseAlbum.Type), protoAlbumList.Albums[i].Type)
		assert.Equal(t, usecaseAlbum.Thumbnail, protoAlbumList.Albums[i].Thumbnail)
		assert.Equal(t, timestamppb.New(usecaseAlbum.ReleaseDate).Seconds, protoAlbumList.Albums[i].ReleaseDate.Seconds)
		assert.Equal(t, timestamppb.New(usecaseAlbum.ReleaseDate).Nanos, protoAlbumList.Albums[i].ReleaseDate.Nanos)
		assert.Equal(t, usecaseAlbum.IsFavorite, protoAlbumList.Albums[i].IsFavorite)
	}
}

func TestAlbumTitleFromUsecaseToProto(t *testing.T) {
	usecaseAlbumTitle := &usecaseModel.AlbumTitle{
		Title: "Test Album",
	}

	protoAlbumTitle := AlbumTitleFromUsecaseToProto(usecaseAlbumTitle)

	assert.Equal(t, usecaseAlbumTitle.Title, protoAlbumTitle.Title)
}

func TestAlbumTitleMapFromRepositoryToUsecase(t *testing.T) {
	repoAlbumTitles := map[int64]string{
		1: "Test Album 1",
		2: "Test Album 2",
	}

	usecaseAlbumTitleMap := AlbumTitleMapFromRepositoryToUsecase(repoAlbumTitles)

	assert.Equal(t, len(repoAlbumTitles), len(usecaseAlbumTitleMap.Titles))

	for id, title := range repoAlbumTitles {
		assert.Contains(t, usecaseAlbumTitleMap.Titles, id)
		assert.Equal(t, title, usecaseAlbumTitleMap.Titles[id].Title)
	}
}

func TestAlbumTitleMapFromUsecaseToProto(t *testing.T) {
	usecaseAlbumTitleMap := &usecaseModel.AlbumTitleMap{
		Titles: map[int64]*usecaseModel.AlbumTitle{
			1: {Title: "Test Album 1"},
			2: {Title: "Test Album 2"},
		},
	}

	protoAlbumTitleMap := AlbumTitleMapFromUsecaseToProto(usecaseAlbumTitleMap)

	assert.Equal(t, len(usecaseAlbumTitleMap.Titles), len(protoAlbumTitleMap.Titles))

	for id, title := range usecaseAlbumTitleMap.Titles {
		assert.Contains(t, protoAlbumTitleMap.Titles, id)
		assert.Equal(t, title.Title, protoAlbumTitleMap.Titles[id].Title)
	}
}

func TestAlbumStreamCreateDataFromUsecaseToProto(t *testing.T) {
	usecaseStreamData := &usecaseModel.AlbumStreamCreateData{
		AlbumID: 1,
		UserID:  2,
	}

	protoStreamData := AlbumStreamCreateDataFromUsecaseToProto(usecaseStreamData)

	assert.Equal(t, usecaseStreamData.AlbumID, protoStreamData.AlbumId.Id)
	assert.Equal(t, usecaseStreamData.UserID, protoStreamData.UserId.Id)
}

func TestAlbumStreamCreateDataFromProtoToUsecase(t *testing.T) {
	protoStreamData := &albumProto.AlbumStreamCreateData{
		AlbumId: &albumProto.AlbumID{Id: 1},
		UserId:  &albumProto.UserID{Id: 2},
	}

	usecaseStreamData := AlbumStreamCreateDataFromProtoToUsecase(protoStreamData)

	assert.Equal(t, protoStreamData.AlbumId.Id, usecaseStreamData.AlbumID)
	assert.Equal(t, protoStreamData.UserId.Id, usecaseStreamData.UserID)
}

func TestLikeRequestFromProtoToUsecase(t *testing.T) {
	protoLikeRequest := &albumProto.LikeRequest{
		AlbumId: &albumProto.AlbumID{Id: 1},
		UserId:  &albumProto.UserID{Id: 2},
		IsLike:  true,
	}

	usecaseLikeRequest := LikeRequestFromProtoToUsecase(protoLikeRequest)

	assert.Equal(t, protoLikeRequest.AlbumId.Id, usecaseLikeRequest.AlbumID)
	assert.Equal(t, protoLikeRequest.UserId.Id, usecaseLikeRequest.UserID)
	assert.Equal(t, protoLikeRequest.IsLike, usecaseLikeRequest.IsLike)
}

func TestLikeRequestFromUsecaseToRepository(t *testing.T) {
	usecaseLikeRequest := &usecaseModel.LikeRequest{
		AlbumID: 1,
		UserID:  2,
		IsLike:  true,
	}

	repoLikeRequest := LikeRequestFromUsecaseToRepository(usecaseLikeRequest)

	assert.Equal(t, usecaseLikeRequest.AlbumID, repoLikeRequest.AlbumID)
	assert.Equal(t, usecaseLikeRequest.UserID, repoLikeRequest.UserID)
}
