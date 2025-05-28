package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	loggerPkg "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/album/model/repository"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/metrics"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func setupTest(t *testing.T) (*sql.DB, sqlmock.Sqlmock, context.Context) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	logger := zap.NewNop().Sugar()
	ctx := loggerPkg.LoggerToContext(context.Background(), logger)

	return db, mock, ctx
}

func TestGetAllAlbums(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewAlbumPostgresRepository(db, metrics.NewMockMetrics())

	filters := &repoModel.AlbumFilters{
		Pagination: &repoModel.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}
	userID := int64(1)

	rows := sqlmock.NewRows([]string{"id", "title", "type", "thumbnail_url", "release_date", "is_favorite"}).
		AddRow(1, "Album 1", "album", "url1", time.Now(), true).
		AddRow(2, "Album 2", "album", "url2", time.Now(), false)

	mock.ExpectPrepare("SELECT a.id, a.title, a.type, a.thumbnail_url, a.release_date").
		ExpectQuery().
		WithArgs(filters.Pagination.Limit, filters.Pagination.Offset, userID).
		WillReturnRows(rows)

	albums, err := repo.GetAllAlbums(ctx, filters, userID)

	require.NoError(t, err)
	assert.Len(t, albums, 2)
	assert.Equal(t, int64(1), albums[0].ID)
	assert.Equal(t, "Album 1", albums[0].Title)
	assert.True(t, albums[0].IsFavorite)
	assert.Equal(t, int64(2), albums[1].ID)
	assert.Equal(t, "Album 2", albums[1].Title)
	assert.False(t, albums[1].IsFavorite)
}

func TestGetAlbumByID(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewAlbumPostgresRepository(db, metrics.NewMockMetrics())

	albumID := int64(1)
	userID := int64(1)
	releaseDate := time.Now()

	rows := sqlmock.NewRows([]string{"id", "title", "type", "thumbnail_url", "release_date", "is_favorite"}).
		AddRow(albumID, "Test Album", "album", "thumbnail_url", releaseDate, true)

	mock.ExpectPrepare("SELECT a.id, a.title, a.type, a.thumbnail_url, a.release_date").
		ExpectQuery().
		WithArgs(albumID, userID).
		WillReturnRows(rows)

	album, err := repo.GetAlbumByID(ctx, albumID, userID)

	require.NoError(t, err)
	assert.Equal(t, albumID, album.ID)
	assert.Equal(t, "Test Album", album.Title)
	assert.Equal(t, repoModel.AlbumTypeAlbum, album.Type)
	assert.Equal(t, "thumbnail_url", album.Thumbnail)
	assert.True(t, album.IsFavorite)
}

func TestGetAlbumByIDNotFound(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewAlbumPostgresRepository(db, metrics.NewMockMetrics())

	albumID := int64(999)
	userID := int64(1)

	mock.ExpectPrepare("SELECT a.id, a.title, a.type, a.thumbnail_url, a.release_date").
		ExpectQuery().
		WithArgs(albumID, userID).
		WillReturnError(sql.ErrNoRows)

	_, err := repo.GetAlbumByID(ctx, albumID, userID)

	require.Error(t, err)
}

func TestGetAlbumTitleByID(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewAlbumPostgresRepository(db, metrics.NewMockMetrics())

	albumID := int64(1)

	rows := sqlmock.NewRows([]string{"title"}).
		AddRow("Test Album")

	mock.ExpectPrepare("SELECT title").
		ExpectQuery().
		WithArgs(albumID).
		WillReturnRows(rows)

	title, err := repo.GetAlbumTitleByID(ctx, albumID)

	require.NoError(t, err)
	assert.Equal(t, "Test Album", title)
}

func TestGetAlbumTitleByIDNotFound(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewAlbumPostgresRepository(db, metrics.NewMockMetrics())

	albumID := int64(999)

	mock.ExpectPrepare("SELECT title").
		ExpectQuery().
		WithArgs(albumID).
		WillReturnError(sql.ErrNoRows)

	_, err := repo.GetAlbumTitleByID(ctx, albumID)

	require.Error(t, err)
}

func TestGetAlbumTitleByIDs(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewAlbumPostgresRepository(db, metrics.NewMockMetrics())

	albumIDs := []int64{1, 2}

	rows := sqlmock.NewRows([]string{"id", "title"}).
		AddRow(1, "Album 1").
		AddRow(2, "Album 2")

	mock.ExpectPrepare("SELECT id, title").
		ExpectQuery().
		WithArgs(pq.Array(albumIDs)).
		WillReturnRows(rows)

	titles, err := repo.GetAlbumTitleByIDs(ctx, albumIDs)

	require.NoError(t, err)
	assert.Len(t, titles, 2)
	assert.Equal(t, "Album 1", titles[1])
	assert.Equal(t, "Album 2", titles[2])
}

func TestGetAlbumsByIDs(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewAlbumPostgresRepository(db, metrics.NewMockMetrics())

	albumIDs := []int64{1, 2}
	userID := int64(1)
	releaseDate := time.Now()

	rows := sqlmock.NewRows([]string{"id", "title", "type", "thumbnail_url", "release_date", "is_favorite"}).
		AddRow(1, "Album 1", "album", "url1", releaseDate, true).
		AddRow(2, "Album 2", "album", "url2", releaseDate, false)

	mock.ExpectPrepare("SELECT a.id, a.title, a.type, a.thumbnail_url, a.release_date").
		ExpectQuery().
		WithArgs(pq.Array(albumIDs), userID).
		WillReturnRows(rows)

	albums, err := repo.GetAlbumsByIDs(ctx, albumIDs, userID)

	require.NoError(t, err)
	assert.Len(t, albums, 2)
	assert.Equal(t, int64(1), albums[0].ID)
	assert.Equal(t, "Album 1", albums[0].Title)
	assert.True(t, albums[0].IsFavorite)
	assert.Equal(t, int64(2), albums[1].ID)
	assert.Equal(t, "Album 2", albums[1].Title)
	assert.False(t, albums[1].IsFavorite)
}

func TestCreateStream(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewAlbumPostgresRepository(db, metrics.NewMockMetrics())

	albumID := int64(1)
	userID := int64(1)

	mock.ExpectPrepare("INSERT INTO album_stream").
		ExpectExec().
		WithArgs(albumID, userID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.CreateStream(ctx, albumID, userID)

	require.NoError(t, err)
}

func TestCheckAlbumExists(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewAlbumPostgresRepository(db, metrics.NewMockMetrics())

	albumID := int64(1)

	rows := sqlmock.NewRows([]string{"exists"}).
		AddRow(true)

	mock.ExpectPrepare("SELECT EXISTS").
		ExpectQuery().
		WithArgs(albumID).
		WillReturnRows(rows)

	exists, err := repo.CheckAlbumExists(ctx, albumID)

	require.NoError(t, err)
	assert.True(t, exists)
}

func TestLikeAlbum(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewAlbumPostgresRepository(db, metrics.NewMockMetrics())

	request := &repoModel.LikeRequest{
		AlbumID: 1,
		UserID:  1,
	}

	mock.ExpectPrepare("INSERT INTO favorite_album").
		ExpectExec().
		WithArgs(request.AlbumID, request.UserID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.LikeAlbum(ctx, request)

	require.NoError(t, err)
}

func TestUnlikeAlbum(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewAlbumPostgresRepository(db, metrics.NewMockMetrics())

	request := &repoModel.LikeRequest{
		AlbumID: 1,
		UserID:  1,
	}

	mock.ExpectPrepare("DELETE FROM favorite_album").
		ExpectExec().
		WithArgs(request.AlbumID, request.UserID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.UnlikeAlbum(ctx, request)

	require.NoError(t, err)
}

func TestGetFavoriteAlbums(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewAlbumPostgresRepository(db, metrics.NewMockMetrics())

	filters := &repoModel.AlbumFilters{
		Pagination: &repoModel.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}
	userID := int64(1)
	releaseDate := time.Now()

	rows := sqlmock.NewRows([]string{"id", "title", "type", "thumbnail_url", "release_date"}).
		AddRow(1, "Album 1", "album", "url1", releaseDate).
		AddRow(2, "Album 2", "album", "url2", releaseDate)

	mock.ExpectPrepare("SELECT a.id, a.title, a.type, a.thumbnail_url, a.release_date").
		ExpectQuery().
		WithArgs(userID, filters.Pagination.Limit, filters.Pagination.Offset).
		WillReturnRows(rows)

	albums, err := repo.GetFavoriteAlbums(ctx, filters, userID)

	require.NoError(t, err)
	assert.Len(t, albums, 2)
	assert.Equal(t, int64(1), albums[0].ID)
	assert.Equal(t, "Album 1", albums[0].Title)
	assert.True(t, albums[0].IsFavorite)
	assert.Equal(t, int64(2), albums[1].ID)
	assert.Equal(t, "Album 2", albums[1].Title)
	assert.True(t, albums[1].IsFavorite)
}

func TestSearchAlbums(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewAlbumPostgresRepository(db, metrics.NewMockMetrics())

	query := "test"
	userID := int64(1)
	releaseDate := time.Now()

	rows := sqlmock.NewRows([]string{"id", "title", "type", "thumbnail_url", "release_date", "is_favorite"}).
		AddRow(1, "Test Album", "album", "url1", releaseDate, true).
		AddRow(2, "Another Test", "album", "url2", releaseDate, false)

	mock.ExpectPrepare("SELECT a.id, a.title, a.type, a.thumbnail_url, a.release_date").
		ExpectQuery().
		WithArgs("test:*", userID, query).
		WillReturnRows(rows)

	albums, err := repo.SearchAlbums(ctx, query, userID)

	require.NoError(t, err)
	assert.Len(t, albums, 2)
	assert.Equal(t, int64(1), albums[0].ID)
	assert.Equal(t, "Test Album", albums[0].Title)
	assert.True(t, albums[0].IsFavorite)
	assert.Equal(t, int64(2), albums[1].ID)
	assert.Equal(t, "Another Test", albums[1].Title)
	assert.False(t, albums[1].IsFavorite)
}

func TestSearchAlbumsMultipleWords(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewAlbumPostgresRepository(db, metrics.NewMockMetrics())

	query := "test album"
	userID := int64(1)
	releaseDate := time.Now()

	rows := sqlmock.NewRows([]string{"id", "title", "type", "thumbnail_url", "release_date", "is_favorite"}).
		AddRow(1, "Test Album", "album", "url1", releaseDate, true)

	mock.ExpectPrepare("SELECT a.id, a.title, a.type, a.thumbnail_url, a.release_date").
		ExpectQuery().
		WithArgs("test:* & album:*", userID, query).
		WillReturnRows(rows)

	albums, err := repo.SearchAlbums(ctx, query, userID)

	require.NoError(t, err)
	assert.Len(t, albums, 1)
	assert.Equal(t, int64(1), albums[0].ID)
	assert.Equal(t, "Test Album", albums[0].Title)
	assert.True(t, albums[0].IsFavorite)
}

func TestGetAllAlbumsError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewAlbumPostgresRepository(db, metrics.NewMockMetrics())

	filters := &repoModel.AlbumFilters{
		Pagination: &repoModel.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}
	userID := int64(1)

	expectedErr := errors.New("database error")
	mock.ExpectPrepare("SELECT a.id, a.title, a.type, a.thumbnail_url, a.release_date").
		ExpectQuery().
		WithArgs(filters.Pagination.Limit, filters.Pagination.Offset, userID).
		WillReturnError(expectedErr)

	_, err := repo.GetAllAlbums(ctx, filters, userID)

	require.Error(t, err)
}

func TestGetAlbumsByIDsError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewAlbumPostgresRepository(db, metrics.NewMockMetrics())

	albumIDs := []int64{1, 2}
	userID := int64(1)

	expectedErr := errors.New("database error")
	mock.ExpectPrepare("SELECT a.id, a.title, a.type, a.thumbnail_url, a.release_date").
		ExpectQuery().
		WithArgs(pq.Array(albumIDs), userID).
		WillReturnError(expectedErr)

	_, err := repo.GetAlbumsByIDs(ctx, albumIDs, userID)

	require.Error(t, err)
}

func TestCreateStreamError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewAlbumPostgresRepository(db, metrics.NewMockMetrics())

	albumID := int64(1)
	userID := int64(1)

	expectedErr := errors.New("database error")
	mock.ExpectPrepare("INSERT INTO album_stream").
		ExpectExec().
		WithArgs(albumID, userID).
		WillReturnError(expectedErr)

	err := repo.CreateStream(ctx, albumID, userID)

	require.Error(t, err)
}

func TestCheckAlbumExistsError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewAlbumPostgresRepository(db, metrics.NewMockMetrics())

	albumID := int64(1)

	expectedErr := errors.New("database error")
	mock.ExpectPrepare("SELECT EXISTS").
		ExpectQuery().
		WithArgs(albumID).
		WillReturnError(expectedErr)

	_, err := repo.CheckAlbumExists(ctx, albumID)

	require.Error(t, err)
}

func TestCreateAlbum(t *testing.T) {
    db, mock, ctx := setupTest(t)
    defer db.Close()

    repo := NewAlbumPostgresRepository(db, metrics.NewMockMetrics())

    album := &repoModel.CreateAlbumRequest{
        Title:     "New Album",
        Type:      repoModel.AlbumTypeAlbum,
        Thumbnail: "thumbnail_url",
        Image:     nil,
        LabelID:   1,
    }

    rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

    mock.ExpectPrepare("INSERT INTO album").
        ExpectQuery().
        WithArgs(album.Title, album.Type, album.Thumbnail, album.LabelID).
        WillReturnRows(rows)

    id, err := repo.CreateAlbum(ctx, album)

    require.NoError(t, err)
    assert.Equal(t, int64(1), id)
}

func TestCreateAlbumError(t *testing.T) {
    db, mock, ctx := setupTest(t)
    defer db.Close()

    repo := NewAlbumPostgresRepository(db, metrics.NewMockMetrics())

    album := &repoModel.CreateAlbumRequest{
        Title:     "New Album",
        Type:      repoModel.AlbumTypeAlbum,
        Thumbnail: "thumbnail_url",
        Image:     nil,
        LabelID:   1,
    }

    mock.ExpectPrepare("INSERT INTO album").
        ExpectQuery().
        WithArgs(album.Title, album.Image, album.Thumbnail, album.LabelID).
        WillReturnError(errors.New("database error"))

    id, err := repo.CreateAlbum(ctx, album)

    require.Error(t, err)
    assert.Equal(t, int64(0), id)
}

func TestDeleteAlbum(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewAlbumPostgresRepository(db, metrics.NewMockMetrics())

	albumID := int64(1)

	mock.ExpectPrepare("DELETE FROM album").
		ExpectExec().
		WithArgs(albumID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.DeleteAlbum(ctx, albumID)

	require.NoError(t, err)
}

func TestDeleteAlbumError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewAlbumPostgresRepository(db, metrics.NewMockMetrics())

	albumID := int64(1)

	mock.ExpectPrepare("DELETE FROM album").
		ExpectExec().
		WithArgs(albumID).
		WillReturnError(errors.New("database error"))

	err := repo.DeleteAlbum(ctx, albumID)

	require.Error(t, err)
}

func TestGetAlbumsLabelID(t *testing.T) {
    db, mock, ctx := setupTest(t)
    defer db.Close()

    repo := NewAlbumPostgresRepository(db, metrics.NewMockMetrics())

    labelID := int64(2)
    filters := &repoModel.AlbumFilters{
        Pagination: &repoModel.Pagination{
            Limit:  10,
            Offset: 0,
        },
    }
    releaseDate := time.Now()

    rows := sqlmock.NewRows([]string{"id", "title", "type", "thumbnail_url", "release_date", "is_favorite"}).
        AddRow(1, "Album 1", "album", "url1", releaseDate, false).
        AddRow(2, "Album 2", "album", "url2", releaseDate, false)

    mock.ExpectPrepare("SELECT a.id, a.title, a.type, a.thumbnail_url, a.release_date").
        ExpectQuery().
        WithArgs(labelID, filters.Pagination.Limit, filters.Pagination.Offset).
        WillReturnRows(rows)

    albums, err := repo.GetAlbumsLabelID(ctx, filters, labelID)

    require.NoError(t, err)
    assert.Len(t, albums, 2)
    assert.Equal(t, int64(1), albums[0].ID)
    assert.Equal(t, "Album 1", albums[0].Title)
    assert.Equal(t, int64(2), albums[1].ID)
    assert.Equal(t, "Album 2", albums[1].Title)
}

func TestGetAlbumsLabelIDError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewAlbumPostgresRepository(db, metrics.NewMockMetrics())

	labelID := int64(2)
	filters := &repoModel.AlbumFilters{
		Pagination: &repoModel.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}

	expectedErr := errors.New("database error")
	mock.ExpectPrepare("SELECT a.id, a.title, a.type, a.thumbnail_url, a.release_date").
		ExpectQuery().
		WithArgs(labelID, filters.Pagination.Limit, filters.Pagination.Offset).
		WillReturnError(expectedErr)

	_, err := repo.GetAlbumsLabelID(ctx, filters, labelID)

	require.Error(t, err)
}