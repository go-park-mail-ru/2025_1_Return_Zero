package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/gen/album"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/gen/artist"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/ctxExtractor"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/customErrors"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestGetAllAlbums(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAlbumClient := mocks.NewMockAlbumServiceClient(ctrl)
	mockArtistClient := mocks.NewMockArtistServiceClient(ctrl)
	albumUsecase := NewUsecase(mockAlbumClient, mockArtistClient)

	ctx := context.Background()
	filters := &usecaseModel.AlbumFilters{
		Pagination: &usecaseModel.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}

	protoFilters := &album.FiltersWithUserID{
		Filters: &album.Filters{
			Pagination: &album.Pagination{
				Limit:  10,
				Offset: 0,
			},
		},
		UserId: &album.UserID{Id: -1},
	}

	releaseDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	mockAlbums := &album.AlbumList{
		Albums: []*album.Album{
			{
				Id:          1,
				Title:       "Album 1",
				Type:        album.AlbumType_AlbumTypeAlbum,
				Thumbnail:   "path/to/image1",
				ReleaseDate: timestamppb.New(releaseDate),
				IsFavorite:  false,
			},
			{
				Id:          2,
				Title:       "Album 2",
				Type:        album.AlbumType_AlbumTypeEP,
				Thumbnail:   "path/to/image2",
				ReleaseDate: timestamppb.New(releaseDate),
				IsFavorite:  true,
			},
		},
	}

	artistMap := &artist.ArtistWithTitleMap{
		Artists: map[int64]*artist.ArtistWithTitleList{
			1: {
				Artists: []*artist.ArtistWithTitle{
					{
						Id:    1,
						Title: "Artist 1",
					},
				},
			},
			2: {
				Artists: []*artist.ArtistWithTitle{
					{
						Id:    2,
						Title: "Artist 2",
					},
				},
			},
		},
	}

	albumIDs := []*artist.AlbumID{
		{Id: 1},
		{Id: 2},
	}

	mockAlbumClient.EXPECT().GetAllAlbums(ctx, protoFilters).Return(mockAlbums, nil)
	mockArtistClient.EXPECT().GetArtistsByAlbumIDs(ctx, &artist.AlbumIDList{Ids: albumIDs}).Return(artistMap, nil)

	albums, err := albumUsecase.GetAllAlbums(ctx, filters)
	assert.NoError(t, err)
	assert.Len(t, albums, 2)
	assert.Equal(t, int64(1), albums[0].ID)
	assert.Equal(t, "Album 1", albums[0].Title)
	assert.Equal(t, "path/to/image1", albums[0].Thumbnail)
	assert.Equal(t, usecaseModel.AlbumTypeAlbum, albums[0].Type)
	assert.Equal(t, releaseDate.Unix(), albums[0].ReleaseDate.Unix())
	assert.False(t, albums[0].IsLiked)
	assert.Len(t, albums[0].Artists, 1)
	assert.Equal(t, int64(1), albums[0].Artists[0].ID)
	assert.Equal(t, "Artist 1", albums[0].Artists[0].Title)

	assert.Equal(t, int64(2), albums[1].ID)
	assert.Equal(t, "Album 2", albums[1].Title)
	assert.Equal(t, usecaseModel.AlbumTypeEP, albums[1].Type)
	assert.True(t, albums[1].IsLiked)
	assert.Len(t, albums[1].Artists, 1)
	assert.Equal(t, int64(2), albums[1].Artists[0].ID)
	assert.Equal(t, "Artist 2", albums[1].Artists[0].Title)
}

func TestGetAllAlbumsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAlbumClient := mocks.NewMockAlbumServiceClient(ctrl)
	mockArtistClient := mocks.NewMockArtistServiceClient(ctrl)
	albumUsecase := NewUsecase(mockAlbumClient, mockArtistClient)

	ctx := context.Background()
	filters := &usecaseModel.AlbumFilters{
		Pagination: &usecaseModel.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}

	protoFilters := &album.FiltersWithUserID{
		Filters: &album.Filters{
			Pagination: &album.Pagination{
				Limit:  10,
				Offset: 0,
			},
		},
		UserId: &album.UserID{Id: -1},
	}

	mockAlbumClient.EXPECT().GetAllAlbums(ctx, protoFilters).Return(nil, status.Error(codes.NotFound, "not found"))

	albums, err := albumUsecase.GetAllAlbums(ctx, filters)
	assert.Error(t, err)
	assert.Equal(t, customErrors.ErrAlbumNotFound, err)
	assert.Nil(t, albums)
}

func TestGetAlbumsByArtistID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAlbumClient := mocks.NewMockAlbumServiceClient(ctrl)
	mockArtistClient := mocks.NewMockArtistServiceClient(ctrl)
	albumUsecase := NewUsecase(mockAlbumClient, mockArtistClient)

	ctx := context.Background()
	artistID := int64(1)
	filters := &usecaseModel.AlbumFilters{
		Pagination: &usecaseModel.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}

	mockAlbumIDs := &artist.AlbumIDList{
		Ids: []*artist.AlbumID{
			{Id: 1},
			{Id: 2},
		},
	}

	releaseDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	mockAlbums := &album.AlbumList{
		Albums: []*album.Album{
			{
				Id:          1,
				Title:       "Album 1",
				Type:        album.AlbumType_AlbumTypeAlbum,
				Thumbnail:   "path/to/image1",
				ReleaseDate: timestamppb.New(releaseDate),
				IsFavorite:  false,
			},
			{
				Id:          2,
				Title:       "Album 2",
				Type:        album.AlbumType_AlbumTypeEP,
				Thumbnail:   "path/to/image2",
				ReleaseDate: timestamppb.New(releaseDate),
				IsFavorite:  true,
			},
		},
	}

	albumIDsForAlbumClient := []*album.AlbumID{
		{Id: 1},
		{Id: 2},
	}

	artistMap := &artist.ArtistWithTitleMap{
		Artists: map[int64]*artist.ArtistWithTitleList{
			1: {
				Artists: []*artist.ArtistWithTitle{
					{
						Id:    1,
						Title: "Artist 1",
					},
				},
			},
			2: {
				Artists: []*artist.ArtistWithTitle{
					{
						Id:    1,
						Title: "Artist 1",
					},
				},
			},
		},
	}

	mockArtistClient.EXPECT().GetAlbumIDsByArtistID(ctx, &artist.ArtistID{Id: artistID}).Return(mockAlbumIDs, nil)
	mockAlbumClient.EXPECT().GetAlbumsByIDs(ctx, &album.AlbumIDListWithUserID{
		Ids:    &album.AlbumIDList{Ids: albumIDsForAlbumClient},
		UserId: &album.UserID{Id: -1},
	}).Return(mockAlbums, nil)
	mockArtistClient.EXPECT().GetArtistsByAlbumIDs(ctx, &artist.AlbumIDList{Ids: mockAlbumIDs.Ids}).Return(artistMap, nil)

	albums, err := albumUsecase.GetAlbumsByArtistID(ctx, artistID, filters)
	assert.NoError(t, err)
	assert.Len(t, albums, 2)
	assert.Equal(t, int64(1), albums[0].ID)
	assert.Equal(t, "Album 1", albums[0].Title)
	assert.Equal(t, int64(2), albums[1].ID)
	assert.Equal(t, "Album 2", albums[1].Title)
}

func TestGetAlbumByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAlbumClient := mocks.NewMockAlbumServiceClient(ctrl)
	mockArtistClient := mocks.NewMockArtistServiceClient(ctrl)
	albumUsecase := NewUsecase(mockAlbumClient, mockArtistClient)

	ctx := context.Background()
	albumID := int64(1)
	releaseDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	mockAlbum := &album.Album{
		Id:          albumID,
		Title:       "Album 1",
		Type:        album.AlbumType_AlbumTypeAlbum,
		Thumbnail:   "path/to/image1",
		ReleaseDate: timestamppb.New(releaseDate),
		IsFavorite:  false,
	}

	mockArtists := &artist.ArtistWithTitleList{
		Artists: []*artist.ArtistWithTitle{
			{
				Id:    1,
				Title: "Artist 1",
			},
		},
	}

	mockAlbumClient.EXPECT().GetAlbumByID(ctx, &album.AlbumIDWithUserID{
		AlbumId: &album.AlbumID{Id: albumID},
		UserId:  &album.UserID{Id: -1},
	}).Return(mockAlbum, nil)
	mockArtistClient.EXPECT().GetArtistsByAlbumID(ctx, &artist.AlbumID{Id: albumID}).Return(mockArtists, nil)

	result, err := albumUsecase.GetAlbumByID(ctx, albumID)
	assert.NoError(t, err)
	assert.Equal(t, albumID, result.ID)
	assert.Equal(t, "Album 1", result.Title)
	assert.Equal(t, "path/to/image1", result.Thumbnail)
	assert.Equal(t, usecaseModel.AlbumTypeAlbum, result.Type)
	assert.Equal(t, releaseDate.Unix(), result.ReleaseDate.Unix())
	assert.False(t, result.IsLiked)
	assert.Len(t, result.Artists, 1)
	assert.Equal(t, int64(1), result.Artists[0].ID)
	assert.Equal(t, "Artist 1", result.Artists[0].Title)
}

func TestLikeAlbum(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAlbumClient := mocks.NewMockAlbumServiceClient(ctrl)
	mockArtistClient := mocks.NewMockArtistServiceClient(ctrl)
	albumUsecase := NewUsecase(mockAlbumClient, mockArtistClient)

	ctx := context.Background()
	likeRequest := &usecaseModel.AlbumLikeRequest{
		AlbumID: 1,
		UserID:  2,
		IsLike:  true,
	}

	protoRequest := &album.LikeRequest{
		AlbumId: &album.AlbumID{Id: 1},
		UserId:  &album.UserID{Id: 2},
		IsLike:  true,
	}

	mockAlbumClient.EXPECT().LikeAlbum(ctx, protoRequest).Return(&emptypb.Empty{}, nil)

	err := albumUsecase.LikeAlbum(ctx, likeRequest)
	assert.NoError(t, err)
}

func TestLikeAlbumError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAlbumClient := mocks.NewMockAlbumServiceClient(ctrl)
	mockArtistClient := mocks.NewMockArtistServiceClient(ctrl)
	albumUsecase := NewUsecase(mockAlbumClient, mockArtistClient)

	ctx := context.Background()
	likeRequest := &usecaseModel.AlbumLikeRequest{
		AlbumID: 1,
		UserID:  2,
		IsLike:  true,
	}

	protoRequest := &album.LikeRequest{
		AlbumId: &album.AlbumID{Id: 1},
		UserId:  &album.UserID{Id: 2},
		IsLike:  true,
	}

	mockAlbumClient.EXPECT().LikeAlbum(ctx, protoRequest).Return(nil, errors.New("some error"))

	err := albumUsecase.LikeAlbum(ctx, likeRequest)
	assert.Error(t, err)
}

func TestGetFavoriteAlbums(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAlbumClient := mocks.NewMockAlbumServiceClient(ctrl)
	mockArtistClient := mocks.NewMockArtistServiceClient(ctrl)
	albumUsecase := NewUsecase(mockAlbumClient, mockArtistClient)

	ctx := context.Background()
	userID := int64(1)
	filters := &usecaseModel.AlbumFilters{
		Pagination: &usecaseModel.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}

	protoFilters := &album.FiltersWithUserID{
		Filters: &album.Filters{
			Pagination: &album.Pagination{
				Limit:  10,
				Offset: 0,
			},
		},
		UserId: &album.UserID{Id: userID},
	}

	releaseDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	mockAlbums := &album.AlbumList{
		Albums: []*album.Album{
			{
				Id:          1,
				Title:       "Album 1",
				Type:        album.AlbumType_AlbumTypeAlbum,
				Thumbnail:   "path/to/image1",
				ReleaseDate: timestamppb.New(releaseDate),
				IsFavorite:  true,
			},
			{
				Id:          2,
				Title:       "Album 2",
				Type:        album.AlbumType_AlbumTypeEP,
				Thumbnail:   "path/to/image2",
				ReleaseDate: timestamppb.New(releaseDate),
				IsFavorite:  true,
			},
		},
	}

	artistMap := &artist.ArtistWithTitleMap{
		Artists: map[int64]*artist.ArtistWithTitleList{
			1: {
				Artists: []*artist.ArtistWithTitle{
					{
						Id:    1,
						Title: "Artist 1",
					},
				},
			},
			2: {
				Artists: []*artist.ArtistWithTitle{
					{
						Id:    2,
						Title: "Artist 2",
					},
				},
			},
		},
	}

	albumIDs := []*artist.AlbumID{
		{Id: 1},
		{Id: 2},
	}

	mockAlbumClient.EXPECT().GetFavoriteAlbums(ctx, protoFilters).Return(mockAlbums, nil)
	mockArtistClient.EXPECT().GetArtistsByAlbumIDs(ctx, &artist.AlbumIDList{Ids: albumIDs}).Return(artistMap, nil)

	albums, err := albumUsecase.GetFavoriteAlbums(ctx, filters, userID)
	assert.NoError(t, err)
	assert.Len(t, albums, 2)
	assert.Equal(t, int64(1), albums[0].ID)
	assert.Equal(t, "Album 1", albums[0].Title)
	assert.True(t, albums[0].IsLiked)
	assert.Equal(t, int64(2), albums[1].ID)
	assert.Equal(t, "Album 2", albums[1].Title)
	assert.True(t, albums[1].IsLiked)
}

func TestSearchAlbums(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAlbumClient := mocks.NewMockAlbumServiceClient(ctrl)
	mockArtistClient := mocks.NewMockArtistServiceClient(ctrl)
	albumUsecase := NewUsecase(mockAlbumClient, mockArtistClient)

	ctx := context.Background()
	query := "search query"

	protoRequest := &album.Query{
		Query:  query,
		UserId: &album.UserID{Id: -1},
	}

	releaseDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	mockAlbums := &album.AlbumList{
		Albums: []*album.Album{
			{
				Id:          1,
				Title:       "Album 1",
				Type:        album.AlbumType_AlbumTypeAlbum,
				Thumbnail:   "path/to/image1",
				ReleaseDate: timestamppb.New(releaseDate),
				IsFavorite:  false,
			},
			{
				Id:          2,
				Title:       "Album 2",
				Type:        album.AlbumType_AlbumTypeEP,
				Thumbnail:   "path/to/image2",
				ReleaseDate: timestamppb.New(releaseDate),
				IsFavorite:  true,
			},
		},
	}

	artistMap := &artist.ArtistWithTitleMap{
		Artists: map[int64]*artist.ArtistWithTitleList{
			1: {
				Artists: []*artist.ArtistWithTitle{
					{
						Id:    1,
						Title: "Artist 1",
					},
				},
			},
			2: {
				Artists: []*artist.ArtistWithTitle{
					{
						Id:    2,
						Title: "Artist 2",
					},
				},
			},
		},
	}

	albumIDs := []*artist.AlbumID{
		{Id: 1},
		{Id: 2},
	}

	mockAlbumClient.EXPECT().SearchAlbums(ctx, protoRequest).Return(mockAlbums, nil)
	mockArtistClient.EXPECT().GetArtistsByAlbumIDs(ctx, &artist.AlbumIDList{Ids: albumIDs}).Return(artistMap, nil)

	albums, err := albumUsecase.SearchAlbums(ctx, query)
	assert.NoError(t, err)
	assert.Len(t, albums, 2)
	assert.Equal(t, int64(1), albums[0].ID)
	assert.Equal(t, "Album 1", albums[0].Title)
	assert.Equal(t, int64(2), albums[1].ID)
	assert.Equal(t, "Album 2", albums[1].Title)
}

// Additional comprehensive test cases

func TestGetAllAlbumsWithAuthenticatedUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAlbumClient := mocks.NewMockAlbumServiceClient(ctrl)
	mockArtistClient := mocks.NewMockArtistServiceClient(ctrl)
	albumUsecase := NewUsecase(mockAlbumClient, mockArtistClient)

	userID := int64(123)
	ctx := context.WithValue(context.Background(), ctxExtractor.UserContextKey{}, userID)
	filters := &usecaseModel.AlbumFilters{
		Pagination: &usecaseModel.Pagination{
			Limit:  5,
			Offset: 10,
		},
	}

	protoFilters := &album.FiltersWithUserID{
		Filters: &album.Filters{
			Pagination: &album.Pagination{
				Limit:  5,
				Offset: 10,
			},
		},
		UserId: &album.UserID{Id: userID},
	}

	releaseDate := time.Date(2021, 6, 15, 0, 0, 0, 0, time.UTC)
	mockAlbums := &album.AlbumList{
		Albums: []*album.Album{
			{
				Id:          3,
				Title:       "Authenticated Album",
				Type:        album.AlbumType_AlbumTypeSingle,
				Thumbnail:   "path/to/auth/image",
				ReleaseDate: timestamppb.New(releaseDate),
				IsFavorite:  true,
			},
		},
	}

	artistMap := &artist.ArtistWithTitleMap{
		Artists: map[int64]*artist.ArtistWithTitleList{
			3: {
				Artists: []*artist.ArtistWithTitle{
					{
						Id:    3,
						Title: "Auth Artist",
					},
				},
			},
		},
	}

	albumIDs := []*artist.AlbumID{{Id: 3}}

	mockAlbumClient.EXPECT().GetAllAlbums(ctx, protoFilters).Return(mockAlbums, nil)
	mockArtistClient.EXPECT().GetArtistsByAlbumIDs(ctx, &artist.AlbumIDList{Ids: albumIDs}).Return(artistMap, nil)

	albums, err := albumUsecase.GetAllAlbums(ctx, filters)
	assert.NoError(t, err)
	assert.Len(t, albums, 1)
	assert.Equal(t, int64(3), albums[0].ID)
	assert.Equal(t, "Authenticated Album", albums[0].Title)
	assert.Equal(t, usecaseModel.AlbumTypeSingle, albums[0].Type)
	assert.True(t, albums[0].IsLiked)
}

func TestGetAllAlbumsArtistServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAlbumClient := mocks.NewMockAlbumServiceClient(ctrl)
	mockArtistClient := mocks.NewMockArtistServiceClient(ctrl)
	albumUsecase := NewUsecase(mockAlbumClient, mockArtistClient)

	ctx := context.Background()
	filters := &usecaseModel.AlbumFilters{
		Pagination: &usecaseModel.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}

	protoFilters := &album.FiltersWithUserID{
		Filters: &album.Filters{
			Pagination: &album.Pagination{
				Limit:  10,
				Offset: 0,
			},
		},
		UserId: &album.UserID{Id: -1},
	}

	mockAlbums := &album.AlbumList{
		Albums: []*album.Album{
			{
				Id:          1,
				Title:       "Album 1",
				Type:        album.AlbumType_AlbumTypeAlbum,
				Thumbnail:   "path/to/image1",
				ReleaseDate: timestamppb.New(time.Now()),
				IsFavorite:  false,
			},
		},
	}

	albumIDs := []*artist.AlbumID{{Id: 1}}

	mockAlbumClient.EXPECT().GetAllAlbums(ctx, protoFilters).Return(mockAlbums, nil)
	mockArtistClient.EXPECT().GetArtistsByAlbumIDs(ctx, &artist.AlbumIDList{Ids: albumIDs}).Return(nil, status.Error(codes.Internal, "artist service error"))

	albums, err := albumUsecase.GetAllAlbums(ctx, filters)
	assert.Error(t, err)
	assert.Nil(t, albums)
}

func TestGetAllAlbumsEmptyResult(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAlbumClient := mocks.NewMockAlbumServiceClient(ctrl)
	mockArtistClient := mocks.NewMockArtistServiceClient(ctrl)
	albumUsecase := NewUsecase(mockAlbumClient, mockArtistClient)

	ctx := context.Background()
	filters := &usecaseModel.AlbumFilters{
		Pagination: &usecaseModel.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}

	protoFilters := &album.FiltersWithUserID{
		Filters: &album.Filters{
			Pagination: &album.Pagination{
				Limit:  10,
				Offset: 0,
			},
		},
		UserId: &album.UserID{Id: -1},
	}

	mockAlbums := &album.AlbumList{Albums: []*album.Album{}}
	artistMap := &artist.ArtistWithTitleMap{Artists: map[int64]*artist.ArtistWithTitleList{}}

	mockAlbumClient.EXPECT().GetAllAlbums(ctx, protoFilters).Return(mockAlbums, nil)
	mockArtistClient.EXPECT().GetArtistsByAlbumIDs(ctx, &artist.AlbumIDList{Ids: []*artist.AlbumID{}}).Return(artistMap, nil)

	albums, err := albumUsecase.GetAllAlbums(ctx, filters)
	assert.NoError(t, err)
	assert.Len(t, albums, 0)
}

func TestGetAlbumsByArtistIDError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAlbumClient := mocks.NewMockAlbumServiceClient(ctrl)
	mockArtistClient := mocks.NewMockArtistServiceClient(ctrl)
	albumUsecase := NewUsecase(mockAlbumClient, mockArtistClient)

	ctx := context.Background()
	artistID := int64(1)
	filters := &usecaseModel.AlbumFilters{
		Pagination: &usecaseModel.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}

	mockArtistClient.EXPECT().GetAlbumIDsByArtistID(ctx, &artist.ArtistID{Id: artistID}).Return(nil, status.Error(codes.NotFound, "artist not found"))

	albums, err := albumUsecase.GetAlbumsByArtistID(ctx, artistID, filters)
	assert.Error(t, err)
	assert.Nil(t, albums)
}

func TestGetAlbumsByArtistIDAlbumServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAlbumClient := mocks.NewMockAlbumServiceClient(ctrl)
	mockArtistClient := mocks.NewMockArtistServiceClient(ctrl)
	albumUsecase := NewUsecase(mockAlbumClient, mockArtistClient)

	ctx := context.Background()
	artistID := int64(1)
	filters := &usecaseModel.AlbumFilters{
		Pagination: &usecaseModel.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}

	mockAlbumIDs := &artist.AlbumIDList{
		Ids: []*artist.AlbumID{{Id: 1}},
	}

	albumIDs := []*album.AlbumID{{Id: 1}}

	mockArtistClient.EXPECT().GetAlbumIDsByArtistID(ctx, &artist.ArtistID{Id: artistID}).Return(mockAlbumIDs, nil)
	mockAlbumClient.EXPECT().GetAlbumsByIDs(ctx, &album.AlbumIDListWithUserID{
		Ids:    &album.AlbumIDList{Ids: albumIDs},
		UserId: &album.UserID{Id: -1},
	}).Return(nil, status.Error(codes.Internal, "album service error"))

	albums, err := albumUsecase.GetAlbumsByArtistID(ctx, artistID, filters)
	assert.Error(t, err)
	assert.Nil(t, albums)
}

func TestGetAlbumsByArtistIDArtistMapError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAlbumClient := mocks.NewMockAlbumServiceClient(ctrl)
	mockArtistClient := mocks.NewMockArtistServiceClient(ctrl)
	albumUsecase := NewUsecase(mockAlbumClient, mockArtistClient)

	ctx := context.Background()
	artistID := int64(1)
	filters := &usecaseModel.AlbumFilters{
		Pagination: &usecaseModel.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}

	mockAlbumIDs := &artist.AlbumIDList{
		Ids: []*artist.AlbumID{{Id: 1}},
	}

	mockAlbums := &album.AlbumList{
		Albums: []*album.Album{
			{
				Id:          1,
				Title:       "Album 1",
				Type:        album.AlbumType_AlbumTypeAlbum,
				Thumbnail:   "path/to/image1",
				ReleaseDate: timestamppb.New(time.Now()),
				IsFavorite:  false,
			},
		},
	}

	albumIDs := []*album.AlbumID{{Id: 1}}

	mockArtistClient.EXPECT().GetAlbumIDsByArtistID(ctx, &artist.ArtistID{Id: artistID}).Return(mockAlbumIDs, nil)
	mockAlbumClient.EXPECT().GetAlbumsByIDs(ctx, &album.AlbumIDListWithUserID{
		Ids:    &album.AlbumIDList{Ids: albumIDs},
		UserId: &album.UserID{Id: -1},
	}).Return(mockAlbums, nil)
	mockArtistClient.EXPECT().GetArtistsByAlbumIDs(ctx, &artist.AlbumIDList{Ids: mockAlbumIDs.Ids}).Return(nil, status.Error(codes.Internal, "artist map error"))

	albums, err := albumUsecase.GetAlbumsByArtistID(ctx, artistID, filters)
	assert.Error(t, err)
	assert.Nil(t, albums)
}

func TestGetAlbumByIDError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAlbumClient := mocks.NewMockAlbumServiceClient(ctrl)
	mockArtistClient := mocks.NewMockArtistServiceClient(ctrl)
	albumUsecase := NewUsecase(mockAlbumClient, mockArtistClient)

	ctx := context.Background()
	albumID := int64(999)

	mockAlbumClient.EXPECT().GetAlbumByID(ctx, &album.AlbumIDWithUserID{
		AlbumId: &album.AlbumID{Id: albumID},
		UserId:  &album.UserID{Id: -1},
	}).Return(nil, status.Error(codes.NotFound, "album not found"))

	result, err := albumUsecase.GetAlbumByID(ctx, albumID)
	assert.Error(t, err)
	assert.Equal(t, customErrors.ErrAlbumNotFound, err)
	assert.Nil(t, result)
}

func TestGetAlbumByIDArtistError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAlbumClient := mocks.NewMockAlbumServiceClient(ctrl)
	mockArtistClient := mocks.NewMockArtistServiceClient(ctrl)
	albumUsecase := NewUsecase(mockAlbumClient, mockArtistClient)

	ctx := context.Background()
	albumID := int64(1)

	mockAlbum := &album.Album{
		Id:          albumID,
		Title:       "Album 1",
		Type:        album.AlbumType_AlbumTypeAlbum,
		Thumbnail:   "path/to/image1",
		ReleaseDate: timestamppb.New(time.Now()),
		IsFavorite:  false,
	}

	mockAlbumClient.EXPECT().GetAlbumByID(ctx, &album.AlbumIDWithUserID{
		AlbumId: &album.AlbumID{Id: albumID},
		UserId:  &album.UserID{Id: -1},
	}).Return(mockAlbum, nil)
	mockArtistClient.EXPECT().GetArtistsByAlbumID(ctx, &artist.AlbumID{Id: albumID}).Return(nil, status.Error(codes.Internal, "artist service error"))

	result, err := albumUsecase.GetAlbumByID(ctx, albumID)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetAlbumByIDWithAuthenticatedUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAlbumClient := mocks.NewMockAlbumServiceClient(ctrl)
	mockArtistClient := mocks.NewMockArtistServiceClient(ctrl)
	albumUsecase := NewUsecase(mockAlbumClient, mockArtistClient)

	userID := int64(456)
	ctx := context.WithValue(context.Background(), ctxExtractor.UserContextKey{}, userID)
	albumID := int64(1)
	releaseDate := time.Date(2022, 3, 10, 0, 0, 0, 0, time.UTC)

	mockAlbum := &album.Album{
		Id:          albumID,
		Title:       "User's Favorite Album",
		Type:        album.AlbumType_AlbumTypeEP,
		Thumbnail:   "path/to/favorite/image",
		ReleaseDate: timestamppb.New(releaseDate),
		IsFavorite:  true,
	}

	mockArtists := &artist.ArtistWithTitleList{
		Artists: []*artist.ArtistWithTitle{
			{
				Id:    10,
				Title: "Favorite Artist",
			},
			{
				Id:    11,
				Title: "Another Artist",
			},
		},
	}

	mockAlbumClient.EXPECT().GetAlbumByID(ctx, &album.AlbumIDWithUserID{
		AlbumId: &album.AlbumID{Id: albumID},
		UserId:  &album.UserID{Id: userID},
	}).Return(mockAlbum, nil)
	mockArtistClient.EXPECT().GetArtistsByAlbumID(ctx, &artist.AlbumID{Id: albumID}).Return(mockArtists, nil)

	result, err := albumUsecase.GetAlbumByID(ctx, albumID)
	assert.NoError(t, err)
	assert.Equal(t, albumID, result.ID)
	assert.Equal(t, "User's Favorite Album", result.Title)
	assert.Equal(t, usecaseModel.AlbumTypeEP, result.Type)
	assert.True(t, result.IsLiked)
	assert.Len(t, result.Artists, 2)
	assert.Equal(t, int64(10), result.Artists[0].ID)
	assert.Equal(t, "Favorite Artist", result.Artists[0].Title)
	assert.Equal(t, int64(11), result.Artists[1].ID)
	assert.Equal(t, "Another Artist", result.Artists[1].Title)
}

func TestLikeAlbumUnlike(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAlbumClient := mocks.NewMockAlbumServiceClient(ctrl)
	mockArtistClient := mocks.NewMockArtistServiceClient(ctrl)
	albumUsecase := NewUsecase(mockAlbumClient, mockArtistClient)

	ctx := context.Background()
	likeRequest := &usecaseModel.AlbumLikeRequest{
		AlbumID: 5,
		UserID:  10,
		IsLike:  false, // Unlike
	}

	protoRequest := &album.LikeRequest{
		AlbumId: &album.AlbumID{Id: 5},
		UserId:  &album.UserID{Id: 10},
		IsLike:  false,
	}

	mockAlbumClient.EXPECT().LikeAlbum(ctx, protoRequest).Return(&emptypb.Empty{}, nil)

	err := albumUsecase.LikeAlbum(ctx, likeRequest)
	assert.NoError(t, err)
}

func TestGetFavoriteAlbumsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAlbumClient := mocks.NewMockAlbumServiceClient(ctrl)
	mockArtistClient := mocks.NewMockArtistServiceClient(ctrl)
	albumUsecase := NewUsecase(mockAlbumClient, mockArtistClient)

	ctx := context.Background()
	userID := int64(1)
	filters := &usecaseModel.AlbumFilters{
		Pagination: &usecaseModel.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}

	protoFilters := &album.FiltersWithUserID{
		Filters: &album.Filters{
			Pagination: &album.Pagination{
				Limit:  10,
				Offset: 0,
			},
		},
		UserId: &album.UserID{Id: userID},
	}

	mockAlbumClient.EXPECT().GetFavoriteAlbums(ctx, protoFilters).Return(nil, status.Error(codes.Internal, "database error"))

	albums, err := albumUsecase.GetFavoriteAlbums(ctx, filters, userID)
	assert.Error(t, err)
	assert.Nil(t, albums)
}

func TestGetFavoriteAlbumsArtistError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAlbumClient := mocks.NewMockAlbumServiceClient(ctrl)
	mockArtistClient := mocks.NewMockArtistServiceClient(ctrl)
	albumUsecase := NewUsecase(mockAlbumClient, mockArtistClient)

	ctx := context.Background()
	userID := int64(1)
	filters := &usecaseModel.AlbumFilters{
		Pagination: &usecaseModel.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}

	protoFilters := &album.FiltersWithUserID{
		Filters: &album.Filters{
			Pagination: &album.Pagination{
				Limit:  10,
				Offset: 0,
			},
		},
		UserId: &album.UserID{Id: userID},
	}

	mockAlbums := &album.AlbumList{
		Albums: []*album.Album{
			{
				Id:          1,
				Title:       "Favorite Album",
				Type:        album.AlbumType_AlbumTypeAlbum,
				Thumbnail:   "path/to/image",
				ReleaseDate: timestamppb.New(time.Now()),
				IsFavorite:  true,
			},
		},
	}

	albumIDs := []*artist.AlbumID{{Id: 1}}

	mockAlbumClient.EXPECT().GetFavoriteAlbums(ctx, protoFilters).Return(mockAlbums, nil)
	mockArtistClient.EXPECT().GetArtistsByAlbumIDs(ctx, &artist.AlbumIDList{Ids: albumIDs}).Return(nil, status.Error(codes.Internal, "artist service error"))

	albums, err := albumUsecase.GetFavoriteAlbums(ctx, filters, userID)
	assert.Error(t, err)
	assert.Nil(t, albums)
}

func TestGetFavoriteAlbumsEmpty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAlbumClient := mocks.NewMockAlbumServiceClient(ctrl)
	mockArtistClient := mocks.NewMockArtistServiceClient(ctrl)
	albumUsecase := NewUsecase(mockAlbumClient, mockArtistClient)

	ctx := context.Background()
	userID := int64(1)
	filters := &usecaseModel.AlbumFilters{
		Pagination: &usecaseModel.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}

	protoFilters := &album.FiltersWithUserID{
		Filters: &album.Filters{
			Pagination: &album.Pagination{
				Limit:  10,
				Offset: 0,
			},
		},
		UserId: &album.UserID{Id: userID},
	}

	mockAlbums := &album.AlbumList{Albums: []*album.Album{}}
	artistMap := &artist.ArtistWithTitleMap{Artists: map[int64]*artist.ArtistWithTitleList{}}

	mockAlbumClient.EXPECT().GetFavoriteAlbums(ctx, protoFilters).Return(mockAlbums, nil)
	mockArtistClient.EXPECT().GetArtistsByAlbumIDs(ctx, &artist.AlbumIDList{Ids: []*artist.AlbumID{}}).Return(artistMap, nil)

	albums, err := albumUsecase.GetFavoriteAlbums(ctx, filters, userID)
	assert.NoError(t, err)
	assert.Len(t, albums, 0)
}

func TestSearchAlbumsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAlbumClient := mocks.NewMockAlbumServiceClient(ctrl)
	mockArtistClient := mocks.NewMockArtistServiceClient(ctrl)
	albumUsecase := NewUsecase(mockAlbumClient, mockArtistClient)

	ctx := context.Background()
	query := "search query"

	protoRequest := &album.Query{
		Query:  query,
		UserId: &album.UserID{Id: -1},
	}

	mockAlbumClient.EXPECT().SearchAlbums(ctx, protoRequest).Return(nil, status.Error(codes.Internal, "search service error"))

	albums, err := albumUsecase.SearchAlbums(ctx, query)
	assert.Error(t, err)
	assert.Nil(t, albums)
}

func TestSearchAlbumsArtistError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAlbumClient := mocks.NewMockAlbumServiceClient(ctrl)
	mockArtistClient := mocks.NewMockArtistServiceClient(ctrl)
	albumUsecase := NewUsecase(mockAlbumClient, mockArtistClient)

	ctx := context.Background()
	query := "search query"

	protoRequest := &album.Query{
		Query:  query,
		UserId: &album.UserID{Id: -1},
	}

	mockAlbums := &album.AlbumList{
		Albums: []*album.Album{
			{
				Id:          1,
				Title:       "Search Result",
				Type:        album.AlbumType_AlbumTypeAlbum,
				Thumbnail:   "path/to/image",
				ReleaseDate: timestamppb.New(time.Now()),
				IsFavorite:  false,
			},
		},
	}

	albumIDs := []*artist.AlbumID{{Id: 1}}

	mockAlbumClient.EXPECT().SearchAlbums(ctx, protoRequest).Return(mockAlbums, nil)
	mockArtistClient.EXPECT().GetArtistsByAlbumIDs(ctx, &artist.AlbumIDList{Ids: albumIDs}).Return(nil, status.Error(codes.Internal, "artist service error"))

	albums, err := albumUsecase.SearchAlbums(ctx, query)
	assert.Error(t, err)
	assert.Nil(t, albums)
}

func TestSearchAlbumsWithAuthenticatedUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAlbumClient := mocks.NewMockAlbumServiceClient(ctrl)
	mockArtistClient := mocks.NewMockArtistServiceClient(ctrl)
	albumUsecase := NewUsecase(mockAlbumClient, mockArtistClient)

	userID := int64(789)
	ctx := context.WithValue(context.Background(), ctxExtractor.UserContextKey{}, userID)
	query := "authenticated search"

	protoRequest := &album.Query{
		Query:  query,
		UserId: &album.UserID{Id: userID},
	}

	releaseDate := time.Date(2023, 8, 20, 0, 0, 0, 0, time.UTC)
	mockAlbums := &album.AlbumList{
		Albums: []*album.Album{
			{
				Id:          100,
				Title:       "Authenticated Search Result",
				Type:        album.AlbumType_AlbumTypeCompilation,
				Thumbnail:   "path/to/auth/search/image",
				ReleaseDate: timestamppb.New(releaseDate),
				IsFavorite:  true,
			},
		},
	}

	artistMap := &artist.ArtistWithTitleMap{
		Artists: map[int64]*artist.ArtistWithTitleList{
			100: {
				Artists: []*artist.ArtistWithTitle{
					{
						Id:    50,
						Title: "Search Artist",
					},
				},
			},
		},
	}

	albumIDs := []*artist.AlbumID{{Id: 100}}

	mockAlbumClient.EXPECT().SearchAlbums(ctx, protoRequest).Return(mockAlbums, nil)
	mockArtistClient.EXPECT().GetArtistsByAlbumIDs(ctx, &artist.AlbumIDList{Ids: albumIDs}).Return(artistMap, nil)

	albums, err := albumUsecase.SearchAlbums(ctx, query)
	assert.NoError(t, err)
	assert.Len(t, albums, 1)
	assert.Equal(t, int64(100), albums[0].ID)
	assert.Equal(t, "Authenticated Search Result", albums[0].Title)
	assert.Equal(t, usecaseModel.AlbumTypeCompilation, albums[0].Type)
	assert.True(t, albums[0].IsLiked)
	assert.Len(t, albums[0].Artists, 1)
	assert.Equal(t, int64(50), albums[0].Artists[0].ID)
	assert.Equal(t, "Search Artist", albums[0].Artists[0].Title)
}

func TestSearchAlbumsEmpty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAlbumClient := mocks.NewMockAlbumServiceClient(ctrl)
	mockArtistClient := mocks.NewMockArtistServiceClient(ctrl)
	albumUsecase := NewUsecase(mockAlbumClient, mockArtistClient)

	ctx := context.Background()
	query := "no results query"

	protoRequest := &album.Query{
		Query:  query,
		UserId: &album.UserID{Id: -1},
	}

	mockAlbums := &album.AlbumList{Albums: []*album.Album{}}
	artistMap := &artist.ArtistWithTitleMap{Artists: map[int64]*artist.ArtistWithTitleList{}}

	mockAlbumClient.EXPECT().SearchAlbums(ctx, protoRequest).Return(mockAlbums, nil)
	mockArtistClient.EXPECT().GetArtistsByAlbumIDs(ctx, &artist.AlbumIDList{Ids: []*artist.AlbumID{}}).Return(artistMap, nil)

	albums, err := albumUsecase.SearchAlbums(ctx, query)
	assert.NoError(t, err)
	assert.Len(t, albums, 0)
}

func TestNewUsecase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAlbumClient := mocks.NewMockAlbumServiceClient(ctrl)
	mockArtistClient := mocks.NewMockArtistServiceClient(ctrl)

	usecase := NewUsecase(mockAlbumClient, mockArtistClient)
	assert.NotNil(t, usecase)

	// Verify that the usecase implements the expected interface
	albumUsecase, ok := usecase.(*albumUsecase)
	assert.True(t, ok)
	assert.Equal(t, mockAlbumClient, albumUsecase.albumClient)
	assert.Equal(t, mockArtistClient, albumUsecase.artistClient)
}
