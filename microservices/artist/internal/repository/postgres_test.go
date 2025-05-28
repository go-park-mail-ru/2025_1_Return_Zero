package repository

import (
	"context"
	"database/sql"
	stderrors "errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	loggerPkg "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
	artistErrors "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/artist/model/errors"
	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/artist/model/repository"
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

func TestGetAllArtists(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewArtistPostgresRepository(db, metrics.NewMockMetrics())
	filters := &repoModel.Filters{
		Pagination: &repoModel.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}
	userID := int64(1)

	rows := sqlmock.NewRows([]string{"id", "title", "description", "thumbnail_url", "is_favorite"}).
		AddRow(1, "Artist 1", "Description 1", "thumbnail1.jpg", true).
		AddRow(2, "Artist 2", "Description 2", "thumbnail2.jpg", false)

	mock.ExpectPrepare("SELECT artist.id, artist.title, artist.description, artist.thumbnail_url").
		ExpectQuery().
		WithArgs(filters.Pagination.Limit, filters.Pagination.Offset, userID).
		WillReturnRows(rows)

	artists, err := repo.GetAllArtists(ctx, filters, userID)
	assert.NoError(t, err)
	assert.Len(t, artists, 2)
	assert.Equal(t, int64(1), artists[0].ID)
	assert.Equal(t, "Artist 1", artists[0].Title)
	assert.Equal(t, "Description 1", artists[0].Description)
	assert.Equal(t, "thumbnail1.jpg", artists[0].Thumbnail)
	assert.True(t, artists[0].IsFavorite)

	assert.Equal(t, int64(2), artists[1].ID)
	assert.Equal(t, "Artist 2", artists[1].Title)
	assert.Equal(t, "Description 2", artists[1].Description)
	assert.Equal(t, "thumbnail2.jpg", artists[1].Thumbnail)
	assert.False(t, artists[1].IsFavorite)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllArtistsError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewArtistPostgresRepository(db, metrics.NewMockMetrics())
	filters := &repoModel.Filters{
		Pagination: &repoModel.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}
	userID := int64(1)

	mock.ExpectPrepare("SELECT artist.id, artist.title, artist.description, artist.thumbnail_url").
		ExpectQuery().
		WithArgs(filters.Pagination.Limit, filters.Pagination.Offset, userID).
		WillReturnError(stderrors.New("db error"))

	artists, err := repo.GetAllArtists(ctx, filters, userID)
	assert.Error(t, err)
	assert.Nil(t, artists)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetArtistByID(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewArtistPostgresRepository(db, metrics.NewMockMetrics())
	artistID := int64(1)
	userID := int64(1)

	rows := sqlmock.NewRows([]string{"id", "title", "description", "thumbnail_url", "is_favorite"}).
		AddRow(1, "Artist 1", "Description 1", "thumbnail1.jpg", true)

	mock.ExpectPrepare("SELECT artist.id, artist.title, artist.description, artist.thumbnail_url").
		ExpectQuery().
		WithArgs(artistID, userID).
		WillReturnRows(rows)

	artist, err := repo.GetArtistByID(ctx, artistID, userID)
	assert.NoError(t, err)
	assert.NotNil(t, artist)
	assert.Equal(t, int64(1), artist.ID)
	assert.Equal(t, "Artist 1", artist.Title)
	assert.Equal(t, "Description 1", artist.Description)
	assert.Equal(t, "thumbnail1.jpg", artist.Thumbnail)
	assert.True(t, artist.IsFavorite)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetArtistByIDNotFound(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewArtistPostgresRepository(db, metrics.NewMockMetrics())
	artistID := int64(1)
	userID := int64(1)

	mock.ExpectPrepare("SELECT artist.id, artist.title, artist.description, artist.thumbnail_url").
		ExpectQuery().
		WithArgs(artistID, userID).
		WillReturnError(sql.ErrNoRows)

	artist, err := repo.GetArtistByID(ctx, artistID, userID)
	assert.Error(t, err)
	assert.Equal(t, artistErrors.ErrArtistNotFound, err)
	assert.Nil(t, artist)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetArtistTitleByID(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewArtistPostgresRepository(db, metrics.NewMockMetrics())
	artistID := int64(1)

	rows := sqlmock.NewRows([]string{"title"}).AddRow("Artist 1")

	mock.ExpectPrepare("SELECT title").
		ExpectQuery().
		WithArgs(artistID).
		WillReturnRows(rows)

	title, err := repo.GetArtistTitleByID(ctx, artistID)
	assert.NoError(t, err)
	assert.Equal(t, "Artist 1", title)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetArtistTitleByIDNotFound(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewArtistPostgresRepository(db, metrics.NewMockMetrics())
	artistID := int64(1)

	mock.ExpectPrepare("SELECT title").
		ExpectQuery().
		WithArgs(artistID).
		WillReturnError(sql.ErrNoRows)

	title, err := repo.GetArtistTitleByID(ctx, artistID)
	assert.Error(t, err)
	assert.Equal(t, artistErrors.ErrArtistNotFound, err)
	assert.Equal(t, "", title)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetArtistsByTrackID(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewArtistPostgresRepository(db, metrics.NewMockMetrics())
	trackID := int64(1)

	rows := sqlmock.NewRows([]string{"id", "title", "role"}).
		AddRow(1, "Artist 1", "main").
		AddRow(2, "Artist 2", "featured")

	mock.ExpectPrepare("SELECT a.id, a.title, ta.role").
		ExpectQuery().
		WithArgs(trackID).
		WillReturnRows(rows)

	artists, err := repo.GetArtistsByTrackID(ctx, trackID)
	assert.NoError(t, err)
	assert.Len(t, artists, 2)
	assert.Equal(t, int64(1), artists[0].ID)
	assert.Equal(t, "Artist 1", artists[0].Title)
	assert.Equal(t, "main", artists[0].Role)
	assert.Equal(t, int64(2), artists[1].ID)
	assert.Equal(t, "Artist 2", artists[1].Title)
	assert.Equal(t, "featured", artists[1].Role)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetArtistsByTrackIDs(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewArtistPostgresRepository(db, metrics.NewMockMetrics())
	trackIDs := []int64{1, 2}

	rows := sqlmock.NewRows([]string{"id", "title", "role", "track_id"}).
		AddRow(1, "Artist 1", "main", 1).
		AddRow(2, "Artist 2", "featured", 1).
		AddRow(3, "Artist 3", "main", 2)

	mock.ExpectPrepare("SELECT a.id, a.title, ta.role, ta.track_id").
		ExpectQuery().
		WithArgs(pq.Array(trackIDs)).
		WillReturnRows(rows)

	artists, err := repo.GetArtistsByTrackIDs(ctx, trackIDs)
	assert.NoError(t, err)
	assert.Len(t, artists, 2)
	assert.Len(t, artists[1], 2)
	assert.Len(t, artists[2], 1)

	assert.Equal(t, int64(1), artists[1][0].ID)
	assert.Equal(t, "Artist 1", artists[1][0].Title)
	assert.Equal(t, "main", artists[1][0].Role)

	assert.Equal(t, int64(2), artists[1][1].ID)
	assert.Equal(t, "Artist 2", artists[1][1].Title)
	assert.Equal(t, "featured", artists[1][1].Role)

	assert.Equal(t, int64(3), artists[2][0].ID)
	assert.Equal(t, "Artist 3", artists[2][0].Title)
	assert.Equal(t, "main", artists[2][0].Role)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetArtistStats(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewArtistPostgresRepository(db, metrics.NewMockMetrics())
	artistID := int64(1)

	rows := sqlmock.NewRows([]string{"listeners_count", "favorites_count"}).
		AddRow(100, 50)

	mock.ExpectPrepare("SELECT listeners_count, favorites_count").
		ExpectQuery().
		WithArgs(artistID).
		WillReturnRows(rows)

	stats, err := repo.GetArtistStats(ctx, artistID)
	assert.NoError(t, err)
	assert.Equal(t, int64(100), stats.ListenersCount)
	assert.Equal(t, int64(50), stats.FavoritesCount)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetArtistsByAlbumID(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewArtistPostgresRepository(db, metrics.NewMockMetrics())
	albumID := int64(1)

	rows := sqlmock.NewRows([]string{"id", "title"}).
		AddRow(1, "Artist 1").
		AddRow(2, "Artist 2")

	mock.ExpectPrepare("SELECT a.id, a.title").
		ExpectQuery().
		WithArgs(albumID).
		WillReturnRows(rows)

	artists, err := repo.GetArtistsByAlbumID(ctx, albumID)
	assert.NoError(t, err)
	assert.Len(t, artists, 2)
	assert.Equal(t, int64(1), artists[0].ID)
	assert.Equal(t, "Artist 1", artists[0].Title)
	assert.Equal(t, int64(2), artists[1].ID)
	assert.Equal(t, "Artist 2", artists[1].Title)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetArtistsByAlbumIDs(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewArtistPostgresRepository(db, metrics.NewMockMetrics())
	albumIDs := []int64{1, 2}

	rows := sqlmock.NewRows([]string{"id", "title", "album_id"}).
		AddRow(1, "Artist 1", 1).
		AddRow(2, "Artist 2", 1).
		AddRow(3, "Artist 3", 2)

	mock.ExpectPrepare("SELECT a.id, a.title, aa.album_id").
		ExpectQuery().
		WithArgs(pq.Array(albumIDs)).
		WillReturnRows(rows)

	artists, err := repo.GetArtistsByAlbumIDs(ctx, albumIDs)
	assert.NoError(t, err)
	assert.Len(t, artists, 2)
	assert.Len(t, artists[1], 2)
	assert.Len(t, artists[2], 1)

	assert.Equal(t, int64(1), artists[1][0].ID)
	assert.Equal(t, "Artist 1", artists[1][0].Title)

	assert.Equal(t, int64(2), artists[1][1].ID)
	assert.Equal(t, "Artist 2", artists[1][1].Title)

	assert.Equal(t, int64(3), artists[2][0].ID)
	assert.Equal(t, "Artist 3", artists[2][0].Title)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAlbumIDsByArtistID(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewArtistPostgresRepository(db, metrics.NewMockMetrics())
	artistID := int64(1)

	rows := sqlmock.NewRows([]string{"album_id"}).
		AddRow(1).
		AddRow(2)

	mock.ExpectPrepare("SELECT album_id").
		ExpectQuery().
		WithArgs(artistID).
		WillReturnRows(rows)

	albumIDs, err := repo.GetAlbumIDsByArtistID(ctx, artistID)
	assert.NoError(t, err)
	assert.Len(t, albumIDs, 2)
	assert.Equal(t, int64(1), albumIDs[0])
	assert.Equal(t, int64(2), albumIDs[1])

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTrackIDsByArtistID(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewArtistPostgresRepository(db, metrics.NewMockMetrics())
	artistID := int64(1)

	rows := sqlmock.NewRows([]string{"track_id"}).
		AddRow(1).
		AddRow(2)

	mock.ExpectPrepare("SELECT track_id").
		ExpectQuery().
		WithArgs(artistID).
		WillReturnRows(rows)

	trackIDs, err := repo.GetTrackIDsByArtistID(ctx, artistID)
	assert.NoError(t, err)
	assert.Len(t, trackIDs, 2)
	assert.Equal(t, int64(1), trackIDs[0])
	assert.Equal(t, int64(2), trackIDs[1])

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateStreamsByArtistIDs(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewArtistPostgresRepository(db, metrics.NewMockMetrics())
	data := &repoModel.ArtistStreamCreateDataList{
		ArtistIDs: []int64{1, 2},
		UserID:    1,
	}

	mock.ExpectBegin()
	mock.ExpectPrepare("INSERT INTO artist_stream").
		ExpectExec().
		WithArgs(pq.Array(data.ArtistIDs), data.UserID).
		WillReturnResult(sqlmock.NewResult(1, 2))
	mock.ExpectCommit()

	err := repo.CreateStreamsByArtistIDs(ctx, data)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateStreamsByArtistIDsEmptyList(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewArtistPostgresRepository(db, metrics.NewMockMetrics())
	data := &repoModel.ArtistStreamCreateDataList{
		ArtistIDs: []int64{},
		UserID:    1,
	}

	err := repo.CreateStreamsByArtistIDs(ctx, data)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetArtistsListenedByUserID(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewArtistPostgresRepository(db, metrics.NewMockMetrics())
	userID := int64(1)

	rows := sqlmock.NewRows([]string{"count"}).
		AddRow(5)

	mock.ExpectPrepare("SELECT COUNT").
		ExpectQuery().
		WithArgs(userID).
		WillReturnRows(rows)

	count, err := repo.GetArtistsListenedByUserID(ctx, userID)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), count)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLikeArtist(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewArtistPostgresRepository(db, metrics.NewMockMetrics())
	request := &repoModel.LikeRequest{
		ArtistID: 1,
		UserID:   1,
	}

	mock.ExpectPrepare("INSERT INTO favorite_artist").
		ExpectExec().
		WithArgs(request.ArtistID, request.UserID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.LikeArtist(ctx, request)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUnlikeArtist(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewArtistPostgresRepository(db, metrics.NewMockMetrics())
	request := &repoModel.LikeRequest{
		ArtistID: 1,
		UserID:   1,
	}

	mock.ExpectPrepare("DELETE FROM favorite_artist").
		ExpectExec().
		WithArgs(request.ArtistID, request.UserID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.UnlikeArtist(ctx, request)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCheckArtistExists(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewArtistPostgresRepository(db, metrics.NewMockMetrics())
	artistID := int64(1)

	rows := sqlmock.NewRows([]string{"exists"}).
		AddRow(true)

	mock.ExpectPrepare("SELECT EXISTS").
		ExpectQuery().
		WithArgs(artistID).
		WillReturnRows(rows)

	exists, err := repo.CheckArtistExists(ctx, artistID)
	assert.NoError(t, err)
	assert.True(t, exists)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetFavoriteArtists(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewArtistPostgresRepository(db, metrics.NewMockMetrics())
	filters := &repoModel.Filters{
		Pagination: &repoModel.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}
	userID := int64(1)

	rows := sqlmock.NewRows([]string{"id", "title", "description", "thumbnail_url"}).
		AddRow(1, "Artist 1", "Description 1", "thumbnail1.jpg").
		AddRow(2, "Artist 2", "Description 2", "thumbnail2.jpg")

	mock.ExpectPrepare("SELECT artist.id, artist.title, artist.description, artist.thumbnail_url").
		ExpectQuery().
		WithArgs(userID, filters.Pagination.Limit, filters.Pagination.Offset).
		WillReturnRows(rows)

	artists, err := repo.GetFavoriteArtists(ctx, filters, userID)
	assert.NoError(t, err)
	assert.Len(t, artists, 2)
	assert.Equal(t, int64(1), artists[0].ID)
	assert.Equal(t, "Artist 1", artists[0].Title)
	assert.Equal(t, "Description 1", artists[0].Description)
	assert.Equal(t, "thumbnail1.jpg", artists[0].Thumbnail)
	assert.False(t, artists[0].IsFavorite)

	assert.Equal(t, int64(2), artists[1].ID)
	assert.Equal(t, "Artist 2", artists[1].Title)
	assert.Equal(t, "Description 2", artists[1].Description)
	assert.Equal(t, "thumbnail2.jpg", artists[1].Thumbnail)
	assert.False(t, artists[1].IsFavorite)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSearchArtists(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewArtistPostgresRepository(db, metrics.NewMockMetrics())
	query := "test artist"
	userID := int64(1)

	rows := sqlmock.NewRows([]string{"id", "title", "description", "thumbnail_url"}).
		AddRow(1, "Test Artist", "Description 1", "thumbnail1.jpg").
		AddRow(2, "Artist Test", "Description 2", "thumbnail2.jpg")

	mock.ExpectPrepare("SELECT a.id, a.title, a.description, a.thumbnail_url").
		ExpectQuery().
		WithArgs("test:* & artist:*", userID, query).
		WillReturnRows(rows)

	artists, err := repo.SearchArtists(ctx, query, userID)
	assert.NoError(t, err)
	assert.Len(t, artists, 2)
	assert.Equal(t, int64(1), artists[0].ID)
	assert.Equal(t, "Test Artist", artists[0].Title)
	assert.Equal(t, "Description 1", artists[0].Description)
	assert.Equal(t, "thumbnail1.jpg", artists[0].Thumbnail)

	assert.Equal(t, int64(2), artists[1].ID)
	assert.Equal(t, "Artist Test", artists[1].Title)
	assert.Equal(t, "Description 2", artists[1].Description)
	assert.Equal(t, "thumbnail2.jpg", artists[1].Thumbnail)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateArtist(t *testing.T) {
    db, mock, ctx := setupTest(t)
    defer db.Close()

    repo := NewArtistPostgresRepository(db, metrics.NewMockMetrics())
    artist := &repoModel.Artist{
        Title:     "New Artist",
        Thumbnail: "new_thumbnail.jpg",
        LabelID:   0, 
    }

    rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
    mock.ExpectPrepare("INSERT INTO artist").
        ExpectQuery().
        WithArgs(artist.Title, artist.Thumbnail, artist.LabelID).
        WillReturnRows(rows)
        
    mock.ExpectExec("REFRESH MATERIALIZED VIEW").
        WillReturnResult(sqlmock.NewResult(0, 0))

    artistNew, err := repo.CreateArtist(ctx, artist)
    assert.NoError(t, err)
    assert.NoError(t, mock.ExpectationsWereMet())
    assert.NotNil(t, artistNew)
    assert.Equal(t, int64(1), artistNew.ID)
}

func TestCreateArtistError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewArtistPostgresRepository(db, metrics.NewMockMetrics())
	artist := &repoModel.Artist{
		Title:     "New Artist",
		Thumbnail: "new_thumbnail.jpg",
		LabelID:   0,
	}

	mock.ExpectPrepare("INSERT INTO artist").
		ExpectQuery().
		WithArgs(artist.Title, artist.Thumbnail, artist.LabelID).
		WillReturnError(stderrors.New("db error"))

	artistNew, err := repo.CreateArtist(ctx, artist)
	assert.Error(t, err)
	assert.Nil(t, artistNew)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCheckArtistNameExist(t *testing.T) {
    db, mock, ctx := setupTest(t)
    defer db.Close()

    repo := NewArtistPostgresRepository(db, metrics.NewMockMetrics())
    artistID := int64(1)

    rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)
    mock.ExpectPrepare("SELECT 1").
        ExpectQuery().
        WithArgs(artistID).
        WillReturnRows(rows)

    exists, err := repo.CheckArtistNameExist(ctx, artistID)
    assert.NoError(t, err)
    assert.True(t, exists)

    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCheckArtistNameExistNotFound(t *testing.T) {
    db, mock, ctx := setupTest(t)
    defer db.Close()

    repo := NewArtistPostgresRepository(db, metrics.NewMockMetrics())
    artistID := int64(1)

    mock.ExpectPrepare("SELECT 1").
        ExpectQuery().
        WithArgs(artistID).
        WillReturnError(sql.ErrNoRows)

    exists, err := repo.CheckArtistNameExist(ctx, artistID)
    assert.NoError(t, err) 
    assert.False(t, exists)

    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCheckArtistNameExistError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewArtistPostgresRepository(db, metrics.NewMockMetrics())
	artistID := int64(1)

	mock.ExpectPrepare("SELECT 1").
		ExpectQuery().
		WithArgs(artistID).
		WillReturnError(stderrors.New("db error"))

	exists, err := repo.CheckArtistNameExist(ctx, artistID)
	assert.Error(t, err)
	assert.False(t, exists)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestChangeArtistTitle(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewArtistPostgresRepository(db, metrics.NewMockMetrics())
	artistID := int64(1)
	newTitle := "Updated Artist Title"

	mock.ExpectPrepare("UPDATE artist").
		ExpectExec().
		WithArgs(newTitle, artistID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.ChangeArtistTitle(ctx, newTitle, artistID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestChangeArtistTitleError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewArtistPostgresRepository(db, metrics.NewMockMetrics())
	artistID := int64(1)
	newTitle := "Updated Artist Title"

	mock.ExpectPrepare("UPDATE artist").
		ExpectExec().
		WithArgs(newTitle, artistID).
		WillReturnError(stderrors.New("db error"))

	err := repo.ChangeArtistTitle(ctx, newTitle, artistID)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetArtistByIDWithoutUserID(t *testing.T) {
    db, mock, ctx := setupTest(t)
    defer db.Close()

    repo := NewArtistPostgresRepository(db, metrics.NewMockMetrics())
    artistID := int64(1)

    rows := sqlmock.NewRows([]string{"id", "title", "description", "thumbnail_url"}).
        AddRow(1, "Artist 1", "Description 1", "thumbnail1.jpg")

    mock.ExpectPrepare("SELECT id, title, description, thumbnail_url").
        ExpectQuery().
        WithArgs(artistID).
        WillReturnRows(rows)

    artist, err := repo.GetArtistByIDWithoutUser(ctx, artistID)
    assert.NoError(t, err)
    assert.NotNil(t, artist)
    assert.Equal(t, int64(1), artist.ID)
    assert.Equal(t, "Artist 1", artist.Title)
    assert.Equal(t, "Description 1", artist.Description)
    assert.Equal(t, "thumbnail1.jpg", artist.Thumbnail)

    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetArtistByIDWithoutUserNotFound(t *testing.T) {
    db, mock, ctx := setupTest(t)
    defer db.Close()

    repo := NewArtistPostgresRepository(db, metrics.NewMockMetrics())
    artistID := int64(1)

    mock.ExpectPrepare("SELECT id, title, description, thumbnail_url").
        ExpectQuery().
        WithArgs(artistID).
        WillReturnError(sql.ErrNoRows)

    artist, err := repo.GetArtistByIDWithoutUser(ctx, artistID)
    assert.Error(t, err)
    assert.Nil(t, artist)

    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetArtistByIDWithoutUserError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewArtistPostgresRepository(db, metrics.NewMockMetrics())
	artistID := int64(1)

	mock.ExpectPrepare("SELECT id, title, description, thumbnail_url").
		ExpectQuery().
		WithArgs(artistID).
		WillReturnError(stderrors.New("db error"))

	artist, err := repo.GetArtistByIDWithoutUser(ctx, artistID)
	assert.Error(t, err)
	assert.Nil(t, artist)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUploadAvatar(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewArtistPostgresRepository(db, metrics.NewMockMetrics())
	artistID := int64(1)
	avatarURL := "new_avatar.jpg"

	mock.ExpectPrepare("UPDATE artist").
		ExpectExec().
		WithArgs(avatarURL, artistID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.UploadAvatar(ctx, artistID, avatarURL)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUploadAvatarError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewArtistPostgresRepository(db, metrics.NewMockMetrics())
	artistID := int64(1)
	avatarURL := "new_avatar.jpg"

	mock.ExpectPrepare("UPDATE artist").
		ExpectExec().
		WithArgs(avatarURL, artistID).
		WillReturnError(stderrors.New("db error"))

	err := repo.UploadAvatar(ctx, artistID, avatarURL)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetArtistLabelID(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewArtistPostgresRepository(db, metrics.NewMockMetrics())
	artistID := int64(1)

	rows := sqlmock.NewRows([]string{"label_id"}).AddRow(10)

	mock.ExpectPrepare("SELECT label_id").
		ExpectQuery().
		WithArgs(artistID).
		WillReturnRows(rows)

	labelID, err := repo.GetArtistLabelID(ctx, artistID)
	assert.NoError(t, err)
	assert.Equal(t, int64(10), labelID)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetArtistLabelIDNotFound(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewArtistPostgresRepository(db, metrics.NewMockMetrics())
	artistID := int64(1)

	mock.ExpectPrepare("SELECT label_id").
		ExpectQuery().
		WithArgs(artistID).
		WillReturnError(sql.ErrNoRows)

	labelID, err := repo.GetArtistLabelID(ctx, artistID)
	assert.Error(t, err)
	assert.Equal(t, int64(0), labelID)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetArtistLabelIDError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewArtistPostgresRepository(db, metrics.NewMockMetrics())
	artistID := int64(1)

	mock.ExpectPrepare("SELECT label_id").
		ExpectQuery().
		WithArgs(artistID).
		WillReturnError(stderrors.New("db error"))

	labelID, err := repo.GetArtistLabelID(ctx, artistID)
	assert.Error(t, err)
	assert.Equal(t, int64(0), labelID)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetArtistsLabelID(t *testing.T) {
    db, mock, ctx := setupTest(t)
    defer db.Close()

    repo := NewArtistPostgresRepository(db, metrics.NewMockMetrics())
    labelID := int64(1) 

    filters := &repoModel.Filters{
        Pagination: &repoModel.Pagination{
            Limit:  10,
            Offset: 0,
        },
    }

    rows := sqlmock.NewRows([]string{"id", "title", "description", "thumbnail_url", "is_favorite"}).
        AddRow(1, "Artist 1", "Description 1", "thumbnail1.jpg", false).
        AddRow(2, "Artist 2", "Description 2", "thumbnail2.jpg", false)

    mock.ExpectPrepare("SELECT artist.id, artist.title, artist.description, artist.thumbnail_url").
        ExpectQuery().
        WithArgs(filters.Pagination.Limit, filters.Pagination.Offset, labelID).
        WillReturnRows(rows)

    artists, err := repo.GetArtistsLabelID(ctx, filters, labelID)
    assert.NoError(t, err)
    assert.Len(t, artists, 2)
    assert.Equal(t, int64(1), artists[0].ID)
    assert.Equal(t, "Artist 1", artists[0].Title)
    assert.Equal(t, int64(2), artists[1].ID)
    assert.Equal(t, "Artist 2", artists[1].Title)

    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetArtistsLabelIDError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewArtistPostgresRepository(db, metrics.NewMockMetrics())
	labelID := int64(1)

	filters := &repoModel.Filters{
		Pagination: &repoModel.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}

	mock.ExpectPrepare("SELECT artist.id, artist.title, artist.description, artist.thumbnail_url").
		ExpectQuery().
		WithArgs(filters.Pagination.Limit, filters.Pagination.Offset, labelID).
		WillReturnError(stderrors.New("db error"))

	artists, err := repo.GetArtistsLabelID(ctx, filters, labelID)
	assert.Error(t, err)
	assert.Nil(t, artists)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteArtist(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewArtistPostgresRepository(db, metrics.NewMockMetrics())
	artistID := int64(1)

	mock.ExpectPrepare("DELETE FROM artist").
		ExpectExec().
		WithArgs(artistID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.DeleteArtist(ctx, artistID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteArtistError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewArtistPostgresRepository(db, metrics.NewMockMetrics())
	artistID := int64(1)

	mock.ExpectPrepare("DELETE FROM artist").
		ExpectExec().
		WithArgs(artistID).
		WillReturnError(stderrors.New("db error"))

	err := repo.DeleteArtist(ctx, artistID)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAddArtistsToAlbum(t *testing.T) {
    db, mock, ctx := setupTest(t)
    defer db.Close()

    repo := NewArtistPostgresRepository(db, metrics.NewMockMetrics())
    albumID := int64(1)
    artistIDs := []int64{1, 2}

    mock.ExpectPrepare("INSERT INTO album_artist").
        ExpectExec().
        WithArgs(pq.Array(artistIDs), albumID).
        WillReturnResult(sqlmock.NewResult(1, 2))

    err := repo.AddArtistsToAlbum(ctx, artistIDs, albumID)
    assert.NoError(t, err)
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAddArtistsToAlbumError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewArtistPostgresRepository(db, metrics.NewMockMetrics())
	albumID := int64(1)
	artistIDs := []int64{1, 2}

	mock.ExpectPrepare("INSERT INTO album_artist").
		ExpectExec().
		WithArgs(pq.Array(artistIDs), albumID).
		WillReturnError(stderrors.New("db error"))

	err := repo.AddArtistsToAlbum(ctx, artistIDs, albumID)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAddArtistsToTracks(t *testing.T) {
    db, mock, ctx := setupTest(t)
    defer db.Close()

    repo := NewArtistPostgresRepository(db, metrics.NewMockMetrics())
    artistIDs := []int64{1, 2}
    trackIDs := []int64{1, 2}

    mock.ExpectPrepare("INSERT INTO track_artist").
        ExpectExec().
        WithArgs(pq.Array(artistIDs), pq.Array(trackIDs)).
        WillReturnResult(sqlmock.NewResult(1, 2))

    err := repo.AddArtistsToTracks(ctx, artistIDs, trackIDs)
    assert.NoError(t, err)
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAddArtistsToTracksError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewArtistPostgresRepository(db, metrics.NewMockMetrics())
	artistIDs := []int64{1, 2}
	trackIDs := []int64{1, 2}

	mock.ExpectPrepare("INSERT INTO track_artist").
		ExpectExec().
		WithArgs(pq.Array(artistIDs), pq.Array(trackIDs)).
		WillReturnError(stderrors.New("db error"))

	err := repo.AddArtistsToTracks(ctx, artistIDs, trackIDs)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}