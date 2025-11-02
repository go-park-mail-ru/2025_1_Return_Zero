package repository

import (
	"context"
	"database/sql"
	stderrors "errors"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	loggerPkg "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/metrics"
	playlistErrors "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/playlist/model/errors"
	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/playlist/model/repository"
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

func TestGetPlaylistByID(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println("Error closing database:", zap.Error(err))
		}
	}()

	repo := NewPlaylistPostgresRepository(db, metrics.NewMockMetrics())
	playlistID := int64(1)

	rows := sqlmock.NewRows([]string{"id", "title", "user_id", "thumbnail_url", "is_public"}).
		AddRow(1, "Playlist 1", 1, "thumbnail1.jpg", true)

	mock.ExpectPrepare("SELECT id, title, user_id, thumbnail_url, is_public")
	mock.ExpectQuery("SELECT id, title, user_id, thumbnail_url, is_public").
		WithArgs(playlistID).
		WillReturnRows(rows)

	playlist, err := repo.GetPlaylistByID(ctx, playlistID)
	assert.NoError(t, err)
	assert.NotNil(t, playlist)
	assert.Equal(t, int64(1), playlist.ID)
	assert.Equal(t, "Playlist 1", playlist.Title)
	assert.Equal(t, int64(1), playlist.UserID)
	assert.Equal(t, "thumbnail1.jpg", playlist.Thumbnail)
	assert.True(t, playlist.IsPublic)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPlaylistByIDNotFound(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println("Error closing database:", zap.Error(err))
		}
	}()

	repo := NewPlaylistPostgresRepository(db, metrics.NewMockMetrics())
	playlistID := int64(1)

	mock.ExpectPrepare("SELECT id, title, user_id, thumbnail_url, is_public")
	mock.ExpectQuery("SELECT id, title, user_id, thumbnail_url, is_public").
		WithArgs(playlistID).
		WillReturnError(sql.ErrNoRows)

	playlist, err := repo.GetPlaylistByID(ctx, playlistID)
	assert.Error(t, err)
	assert.Equal(t, playlistErrors.ErrPlaylistNotFound, err)
	assert.Nil(t, playlist)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreatePlaylist(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println("Error closing database:", zap.Error(err))
		}
	}()

	repo := NewPlaylistPostgresRepository(db, metrics.NewMockMetrics())
	request := &repoModel.CreatePlaylistRequest{
		Title:     "New Playlist",
		UserID:    1,
		Thumbnail: "thumbnail.jpg",
		IsPublic:  true,
	}

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	mock.ExpectPrepare("INSERT INTO playlist")
	mock.ExpectQuery("INSERT INTO playlist").
		WithArgs(request.Title, request.UserID, request.Thumbnail, request.IsPublic).
		WillReturnRows(rows)

	playlistRows := sqlmock.NewRows([]string{"id", "title", "user_id", "thumbnail_url", "is_public"}).
		AddRow(1, "New Playlist", 1, "thumbnail.jpg", true)

	mock.ExpectPrepare("SELECT id, title, user_id, thumbnail_url, is_public")
	mock.ExpectQuery("SELECT id, title, user_id, thumbnail_url, is_public").
		WithArgs(int64(1)).
		WillReturnRows(playlistRows)

	playlist, err := repo.CreatePlaylist(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, playlist)
	assert.Equal(t, int64(1), playlist.ID)
	assert.Equal(t, "New Playlist", playlist.Title)
	assert.Equal(t, int64(1), playlist.UserID)
	assert.Equal(t, "thumbnail.jpg", playlist.Thumbnail)
	assert.True(t, playlist.IsPublic)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreatePlaylistDuplicate(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println("Error closing database:", zap.Error(err))
		}
	}()

	repo := NewPlaylistPostgresRepository(db, metrics.NewMockMetrics())
	request := &repoModel.CreatePlaylistRequest{
		Title:     "New Playlist",
		UserID:    1,
		Thumbnail: "thumbnail.jpg",
		IsPublic:  true,
	}

	mock.ExpectPrepare("INSERT INTO playlist")
	mock.ExpectQuery("INSERT INTO playlist").
		WithArgs(request.Title, request.UserID, request.Thumbnail, request.IsPublic).
		WillReturnError(stderrors.New("duplicate key value violates unique constraint"))

	playlist, err := repo.CreatePlaylist(ctx, request)
	assert.Error(t, err)
	assert.Equal(t, playlistErrors.ErrPlaylistDuplicate, err)
	assert.Nil(t, playlist)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetCombinedPlaylistsByUserID(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println("Error closing database:", zap.Error(err))
		}
	}()

	repo := NewPlaylistPostgresRepository(db, metrics.NewMockMetrics())
	userID := int64(1)

	rows := sqlmock.NewRows([]string{"id", "title", "user_id", "thumbnail_url"}).
		AddRow(1, "Playlist 1", 1, "thumbnail1.jpg").
		AddRow(2, "Playlist 2", 2, "thumbnail2.jpg")

	mock.ExpectPrepare("SELECT p.id, p.title, p.user_id, p.thumbnail_url")
	mock.ExpectQuery("SELECT p.id, p.title, p.user_id, p.thumbnail_url").
		WithArgs(userID).
		WillReturnRows(rows)

	playlists, err := repo.GetCombinedPlaylistsByUserID(ctx, userID)
	assert.NoError(t, err)
	assert.NotNil(t, playlists)
	assert.Len(t, playlists.Playlists, 2)
	assert.Equal(t, int64(1), playlists.Playlists[0].ID)
	assert.Equal(t, "Playlist 1", playlists.Playlists[0].Title)
	assert.Equal(t, int64(1), playlists.Playlists[0].UserID)
	assert.Equal(t, "thumbnail1.jpg", playlists.Playlists[0].Thumbnail)
	assert.Equal(t, int64(2), playlists.Playlists[1].ID)
	assert.Equal(t, "Playlist 2", playlists.Playlists[1].Title)
	assert.Equal(t, int64(2), playlists.Playlists[1].UserID)
	assert.Equal(t, "thumbnail2.jpg", playlists.Playlists[1].Thumbnail)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetCombinedPlaylistsByUserIDError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println("Error closing database:", zap.Error(err))
		}
	}()

	repo := NewPlaylistPostgresRepository(db, metrics.NewMockMetrics())
	userID := int64(1)

	mock.ExpectPrepare("SELECT p.id, p.title, p.user_id, p.thumbnail_url")
	mock.ExpectQuery("SELECT p.id, p.title, p.user_id, p.thumbnail_url").
		WithArgs(userID).
		WillReturnError(stderrors.New("db error"))

	playlists, err := repo.GetCombinedPlaylistsByUserID(ctx, userID)
	assert.Error(t, err)
	assert.Nil(t, playlists)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTrackExistsInPlaylist(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println("Error closing database:", zap.Error(err))
		}
	}()

	repo := NewPlaylistPostgresRepository(db, metrics.NewMockMetrics())
	playlistID := int64(1)
	trackID := int64(2)

	rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)

	mock.ExpectPrepare("SELECT EXISTS")
	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(playlistID, trackID).
		WillReturnRows(rows)

	exists, err := repo.TrackExistsInPlaylist(ctx, playlistID, trackID)
	assert.NoError(t, err)
	assert.True(t, exists)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTrackExistsInPlaylistError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println("Error closing database:", zap.Error(err))
		}
	}()

	repo := NewPlaylistPostgresRepository(db, metrics.NewMockMetrics())
	playlistID := int64(1)
	trackID := int64(2)

	mock.ExpectPrepare("SELECT EXISTS")
	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(playlistID, trackID).
		WillReturnError(stderrors.New("db error"))

	exists, err := repo.TrackExistsInPlaylist(ctx, playlistID, trackID)
	assert.Error(t, err)
	assert.False(t, exists)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAddTrackToPlaylist(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println("Error closing database:", zap.Error(err))
		}
	}()

	repo := NewPlaylistPostgresRepository(db, metrics.NewMockMetrics())
	request := &repoModel.AddTrackToPlaylistRequest{
		PlaylistID: 1,
		TrackID:    2,
		UserID:     1,
	}

	playlistRows := sqlmock.NewRows([]string{"id", "title", "user_id", "thumbnail_url", "is_public"}).
		AddRow(1, "Playlist 1", 1, "thumbnail1.jpg", true)

	mock.ExpectPrepare("INSERT INTO playlist_track")
	mock.ExpectPrepare("SELECT id, title, user_id, thumbnail_url, is_public")
	mock.ExpectQuery("SELECT id, title, user_id, thumbnail_url, is_public").
		WithArgs(request.PlaylistID).
		WillReturnRows(playlistRows)

	existsRows := sqlmock.NewRows([]string{"exists"}).AddRow(false)

	mock.ExpectPrepare("SELECT EXISTS")
	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(request.PlaylistID, request.TrackID).
		WillReturnRows(existsRows)

	mock.ExpectExec("INSERT INTO playlist_track").
		WithArgs(request.PlaylistID, request.TrackID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.AddTrackToPlaylist(ctx, request)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAddTrackToPlaylistGetPlaylistError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println("Error closing database:", zap.Error(err))
		}
	}()

	repo := NewPlaylistPostgresRepository(db, metrics.NewMockMetrics())
	request := &repoModel.AddTrackToPlaylistRequest{
		PlaylistID: 1,
		TrackID:    2,
		UserID:     1,
	}

	mock.ExpectPrepare("INSERT INTO playlist_track")
	mock.ExpectPrepare("SELECT id, title, user_id, thumbnail_url, is_public")
	mock.ExpectQuery("SELECT id, title, user_id, thumbnail_url, is_public").
		WithArgs(request.PlaylistID).
		WillReturnError(stderrors.New("db error"))

	err := repo.AddTrackToPlaylist(ctx, request)
	assert.Error(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAddTrackToPlaylistTrackExistsError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println("Error closing database:", zap.Error(err))
		}
	}()

	repo := NewPlaylistPostgresRepository(db, metrics.NewMockMetrics())
	request := &repoModel.AddTrackToPlaylistRequest{
		PlaylistID: 1,
		TrackID:    2,
		UserID:     1,
	}

	playlistRows := sqlmock.NewRows([]string{"id", "title", "user_id", "thumbnail_url", "is_public"}).
		AddRow(1, "Playlist 1", 1, "thumbnail1.jpg", true)

	mock.ExpectPrepare("INSERT INTO playlist_track")
	mock.ExpectPrepare("SELECT id, title, user_id, thumbnail_url, is_public")
	mock.ExpectQuery("SELECT id, title, user_id, thumbnail_url, is_public").
		WithArgs(request.PlaylistID).
		WillReturnRows(playlistRows)

	mock.ExpectPrepare("SELECT EXISTS")
	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(request.PlaylistID, request.TrackID).
		WillReturnError(stderrors.New("db error"))

	err := repo.AddTrackToPlaylist(ctx, request)
	assert.Error(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAddTrackToPlaylistInsertError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println("Error closing database:", zap.Error(err))
		}
	}()

	repo := NewPlaylistPostgresRepository(db, metrics.NewMockMetrics())
	request := &repoModel.AddTrackToPlaylistRequest{
		PlaylistID: 1,
		TrackID:    2,
		UserID:     1,
	}

	playlistRows := sqlmock.NewRows([]string{"id", "title", "user_id", "thumbnail_url", "is_public"}).
		AddRow(1, "Playlist 1", 1, "thumbnail1.jpg", true)

	mock.ExpectPrepare("INSERT INTO playlist_track")
	mock.ExpectPrepare("SELECT id, title, user_id, thumbnail_url, is_public")
	mock.ExpectQuery("SELECT id, title, user_id, thumbnail_url, is_public").
		WithArgs(request.PlaylistID).
		WillReturnRows(playlistRows)

	existsRows := sqlmock.NewRows([]string{"exists"}).AddRow(false)

	mock.ExpectPrepare("SELECT EXISTS")
	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(request.PlaylistID, request.TrackID).
		WillReturnRows(existsRows)

	mock.ExpectExec("INSERT INTO playlist_track").
		WithArgs(request.PlaylistID, request.TrackID).
		WillReturnError(stderrors.New("db error"))

	err := repo.AddTrackToPlaylist(ctx, request)
	assert.Error(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAddTrackToPlaylistPermissionDenied(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println("Error closing database:", zap.Error(err))
		}
	}()

	repo := NewPlaylistPostgresRepository(db, metrics.NewMockMetrics())
	request := &repoModel.AddTrackToPlaylistRequest{
		PlaylistID: 1,
		TrackID:    2,
		UserID:     2,
	}

	playlistRows := sqlmock.NewRows([]string{"id", "title", "user_id", "thumbnail_url", "is_public"}).
		AddRow(1, "Playlist 1", 1, "thumbnail1.jpg", true)

	mock.ExpectPrepare("INSERT INTO playlist_track")
	mock.ExpectPrepare("SELECT id, title, user_id, thumbnail_url, is_public")
	mock.ExpectQuery("SELECT id, title, user_id, thumbnail_url, is_public").
		WithArgs(request.PlaylistID).
		WillReturnRows(playlistRows)

	err := repo.AddTrackToPlaylist(ctx, request)
	assert.Error(t, err)
	assert.Equal(t, playlistErrors.ErrPlaylistPermissionDenied, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAddTrackToPlaylistDuplicate(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println("Error closing database:", zap.Error(err))
		}
	}()

	repo := NewPlaylistPostgresRepository(db, metrics.NewMockMetrics())
	request := &repoModel.AddTrackToPlaylistRequest{
		PlaylistID: 1,
		TrackID:    2,
		UserID:     1,
	}

	playlistRows := sqlmock.NewRows([]string{"id", "title", "user_id", "thumbnail_url", "is_public"}).
		AddRow(1, "Playlist 1", 1, "thumbnail1.jpg", true)

	mock.ExpectPrepare("INSERT INTO playlist_track")
	mock.ExpectPrepare("SELECT id, title, user_id, thumbnail_url, is_public")
	mock.ExpectQuery("SELECT id, title, user_id, thumbnail_url, is_public").
		WithArgs(request.PlaylistID).
		WillReturnRows(playlistRows)

	existsRows := sqlmock.NewRows([]string{"exists"}).AddRow(true)

	mock.ExpectPrepare("SELECT EXISTS")
	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(request.PlaylistID, request.TrackID).
		WillReturnRows(existsRows)

	err := repo.AddTrackToPlaylist(ctx, request)
	assert.Error(t, err)
	assert.Equal(t, playlistErrors.ErrPlaylistTrackDuplicate, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRemoveTrackFromPlaylist(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println("Error closing database:", zap.Error(err))
		}
	}()

	repo := NewPlaylistPostgresRepository(db, metrics.NewMockMetrics())
	request := &repoModel.RemoveTrackFromPlaylistRequest{
		PlaylistID: 1,
		TrackID:    2,
		UserID:     1,
	}

	playlistRows := sqlmock.NewRows([]string{"id", "title", "user_id", "thumbnail_url", "is_public"}).
		AddRow(1, "Playlist 1", 1, "thumbnail1.jpg", true)

	mock.ExpectPrepare("DELETE FROM playlist_track")
	mock.ExpectPrepare("SELECT id, title, user_id, thumbnail_url, is_public")
	mock.ExpectQuery("SELECT id, title, user_id, thumbnail_url, is_public").
		WithArgs(request.PlaylistID).
		WillReturnRows(playlistRows)

	existsRows := sqlmock.NewRows([]string{"exists"}).AddRow(true)

	mock.ExpectPrepare("SELECT EXISTS")
	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(request.PlaylistID, request.TrackID).
		WillReturnRows(existsRows)

	mock.ExpectExec("DELETE FROM playlist_track").
		WithArgs(request.PlaylistID, request.TrackID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.RemoveTrackFromPlaylist(ctx, request)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRemoveTrackFromPlaylistGetPlaylistError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println("Error closing database:", zap.Error(err))
		}
	}()

	repo := NewPlaylistPostgresRepository(db, metrics.NewMockMetrics())
	request := &repoModel.RemoveTrackFromPlaylistRequest{
		PlaylistID: 1,
		TrackID:    2,
		UserID:     1,
	}

	mock.ExpectPrepare("DELETE FROM playlist_track")
	mock.ExpectPrepare("SELECT id, title, user_id, thumbnail_url, is_public")
	mock.ExpectQuery("SELECT id, title, user_id, thumbnail_url, is_public").
		WithArgs(request.PlaylistID).
		WillReturnError(stderrors.New("db error"))

	err := repo.RemoveTrackFromPlaylist(ctx, request)
	assert.Error(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRemoveTrackFromPlaylistPermissionDenied(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println("Error closing database:", zap.Error(err))
		}
	}()

	repo := NewPlaylistPostgresRepository(db, metrics.NewMockMetrics())
	request := &repoModel.RemoveTrackFromPlaylistRequest{
		PlaylistID: 1,
		TrackID:    2,
		UserID:     2,
	}

	playlistRows := sqlmock.NewRows([]string{"id", "title", "user_id", "thumbnail_url", "is_public"}).
		AddRow(1, "Playlist 1", 1, "thumbnail1.jpg", true)

	mock.ExpectPrepare("DELETE FROM playlist_track")
	mock.ExpectPrepare("SELECT id, title, user_id, thumbnail_url, is_public")
	mock.ExpectQuery("SELECT id, title, user_id, thumbnail_url, is_public").
		WithArgs(request.PlaylistID).
		WillReturnRows(playlistRows)

	err := repo.RemoveTrackFromPlaylist(ctx, request)
	assert.Error(t, err)
	assert.Equal(t, playlistErrors.ErrPlaylistPermissionDenied, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRemoveTrackFromPlaylistTrackNotFound(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println("Error closing database:", zap.Error(err))
		}
	}()

	repo := NewPlaylistPostgresRepository(db, metrics.NewMockMetrics())
	request := &repoModel.RemoveTrackFromPlaylistRequest{
		PlaylistID: 1,
		TrackID:    2,
		UserID:     1,
	}

	playlistRows := sqlmock.NewRows([]string{"id", "title", "user_id", "thumbnail_url", "is_public"}).
		AddRow(1, "Playlist 1", 1, "thumbnail1.jpg", true)

	mock.ExpectPrepare("DELETE FROM playlist_track")
	mock.ExpectPrepare("SELECT id, title, user_id, thumbnail_url, is_public")
	mock.ExpectQuery("SELECT id, title, user_id, thumbnail_url, is_public").
		WithArgs(request.PlaylistID).
		WillReturnRows(playlistRows)

	existsRows := sqlmock.NewRows([]string{"exists"}).AddRow(false)

	mock.ExpectPrepare("SELECT EXISTS")
	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(request.PlaylistID, request.TrackID).
		WillReturnRows(existsRows)

	err := repo.RemoveTrackFromPlaylist(ctx, request)
	assert.Error(t, err)
	assert.Equal(t, playlistErrors.ErrPlaylistTrackNotFound, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRemoveTrackFromPlaylistTrackExistsError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println("Error closing database:", zap.Error(err))
		}
	}()

	repo := NewPlaylistPostgresRepository(db, metrics.NewMockMetrics())
	request := &repoModel.RemoveTrackFromPlaylistRequest{
		PlaylistID: 1,
		TrackID:    2,
		UserID:     1,
	}

	playlistRows := sqlmock.NewRows([]string{"id", "title", "user_id", "thumbnail_url", "is_public"}).
		AddRow(1, "Playlist 1", 1, "thumbnail1.jpg", true)

	mock.ExpectPrepare("DELETE FROM playlist_track")
	mock.ExpectPrepare("SELECT id, title, user_id, thumbnail_url, is_public")
	mock.ExpectQuery("SELECT id, title, user_id, thumbnail_url, is_public").
		WithArgs(request.PlaylistID).
		WillReturnRows(playlistRows)

	mock.ExpectPrepare("SELECT EXISTS")
	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(request.PlaylistID, request.TrackID).
		WillReturnError(stderrors.New("db error"))

	err := repo.RemoveTrackFromPlaylist(ctx, request)
	assert.Error(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRemoveTrackFromPlaylistDeleteError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println("Error closing database:", zap.Error(err))
		}
	}()

	repo := NewPlaylistPostgresRepository(db, metrics.NewMockMetrics())
	request := &repoModel.RemoveTrackFromPlaylistRequest{
		PlaylistID: 1,
		TrackID:    2,
		UserID:     1,
	}

	playlistRows := sqlmock.NewRows([]string{"id", "title", "user_id", "thumbnail_url", "is_public"}).
		AddRow(1, "Playlist 1", 1, "thumbnail1.jpg", true)

	mock.ExpectPrepare("DELETE FROM playlist_track")
	mock.ExpectPrepare("SELECT id, title, user_id, thumbnail_url, is_public")
	mock.ExpectQuery("SELECT id, title, user_id, thumbnail_url, is_public").
		WithArgs(request.PlaylistID).
		WillReturnRows(playlistRows)

	existsRows := sqlmock.NewRows([]string{"exists"}).AddRow(true)

	mock.ExpectPrepare("SELECT EXISTS")
	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(request.PlaylistID, request.TrackID).
		WillReturnRows(existsRows)

	mock.ExpectExec("DELETE FROM playlist_track").
		WithArgs(request.PlaylistID, request.TrackID).
		WillReturnError(stderrors.New("db error"))

	err := repo.RemoveTrackFromPlaylist(ctx, request)
	assert.Error(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPlaylistTrackIds(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println("Error closing database:", zap.Error(err))
		}
	}()

	repo := NewPlaylistPostgresRepository(db, metrics.NewMockMetrics())
	request := &repoModel.GetPlaylistTrackIdsRequest{
		PlaylistID: 1,
	}

	rows := sqlmock.NewRows([]string{"track_id"}).
		AddRow(1).
		AddRow(2).
		AddRow(3)

	mock.ExpectPrepare("SELECT track_id")
	mock.ExpectQuery("SELECT track_id").
		WithArgs(request.PlaylistID).
		WillReturnRows(rows)

	trackIds, err := repo.GetPlaylistTrackIds(ctx, request)
	assert.NoError(t, err)
	assert.Len(t, trackIds, 3)
	assert.Equal(t, int64(1), trackIds[0])
	assert.Equal(t, int64(2), trackIds[1])
	assert.Equal(t, int64(3), trackIds[2])

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPlaylistTrackIdsError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println("Error closing database:", zap.Error(err))
		}
	}()

	repo := NewPlaylistPostgresRepository(db, metrics.NewMockMetrics())
	request := &repoModel.GetPlaylistTrackIdsRequest{
		PlaylistID: 1,
	}

	mock.ExpectPrepare("SELECT track_id")
	mock.ExpectQuery("SELECT track_id").
		WithArgs(request.PlaylistID).
		WillReturnError(stderrors.New("db error"))

	trackIds, err := repo.GetPlaylistTrackIds(ctx, request)
	assert.Error(t, err)
	assert.Nil(t, trackIds)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRemovePlaylist(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println("Error closing database:", zap.Error(err))
		}
	}()

	repo := NewPlaylistPostgresRepository(db, metrics.NewMockMetrics())
	request := &repoModel.RemovePlaylistRequest{
		PlaylistID: 1,
		UserID:     1,
	}

	mock.ExpectPrepare("DELETE FROM playlist")
	mock.ExpectExec("DELETE FROM playlist").
		WithArgs(request.PlaylistID, request.UserID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.RemovePlaylist(ctx, request)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRemovePlaylistError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println("Error closing database:", zap.Error(err))
		}
	}()

	repo := NewPlaylistPostgresRepository(db, metrics.NewMockMetrics())
	request := &repoModel.RemovePlaylistRequest{
		PlaylistID: 1,
		UserID:     1,
	}

	mock.ExpectPrepare("DELETE FROM playlist")
	mock.ExpectExec("DELETE FROM playlist").
		WithArgs(request.PlaylistID, request.UserID).
		WillReturnError(stderrors.New("db error"))

	err := repo.RemovePlaylist(ctx, request)
	assert.Error(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPlaylistsToAdd(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println("Error closing database:", zap.Error(err))
		}
	}()

	repo := NewPlaylistPostgresRepository(db, metrics.NewMockMetrics())
	request := &repoModel.GetPlaylistsToAddRequest{
		TrackID: 1,
		UserID:  2,
	}

	rows := sqlmock.NewRows([]string{"id", "title", "user_id", "thumbnail_url", "is_included"}).
		AddRow(1, "Playlist 1", 2, "thumbnail1.jpg", true).
		AddRow(2, "Playlist 2", 2, "thumbnail2.jpg", false)

	mock.ExpectPrepare("SELECT p.id, p.title, p.user_id, p.thumbnail_url")
	mock.ExpectQuery("SELECT p.id, p.title, p.user_id, p.thumbnail_url").
		WithArgs(request.TrackID, request.UserID).
		WillReturnRows(rows)

	response, err := repo.GetPlaylistsToAdd(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Len(t, response.Playlists, 2)
	assert.Equal(t, int64(1), response.Playlists[0].Playlist.ID)
	assert.Equal(t, "Playlist 1", response.Playlists[0].Playlist.Title)
	assert.True(t, response.Playlists[0].IsIncluded)
	assert.Equal(t, int64(2), response.Playlists[1].Playlist.ID)
	assert.Equal(t, "Playlist 2", response.Playlists[1].Playlist.Title)
	assert.False(t, response.Playlists[1].IsIncluded)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPlaylistsToAddError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println("Error closing database:", zap.Error(err))
		}
	}()

	repo := NewPlaylistPostgresRepository(db, metrics.NewMockMetrics())
	request := &repoModel.GetPlaylistsToAddRequest{
		TrackID: 1,
		UserID:  2,
	}

	mock.ExpectPrepare("SELECT p.id, p.title, p.user_id, p.thumbnail_url")
	mock.ExpectQuery("SELECT p.id, p.title, p.user_id, p.thumbnail_url").
		WithArgs(request.TrackID, request.UserID).
		WillReturnError(stderrors.New("db error"))

	response, err := repo.GetPlaylistsToAdd(ctx, request)
	assert.Error(t, err)
	assert.Nil(t, response)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPlaylistsToAddScanError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println("Error closing database:", zap.Error(err))
		}
	}()

	repo := NewPlaylistPostgresRepository(db, metrics.NewMockMetrics())
	request := &repoModel.GetPlaylistsToAddRequest{
		TrackID: 1,
		UserID:  2,
	}

	rows := sqlmock.NewRows([]string{"id", "title", "user_id", "thumbnail_url", "is_included"}).
		AddRow(1, "Playlist 1", 2, "thumbnail1.jpg", "invalid_bool")

	mock.ExpectPrepare("SELECT p.id, p.title, p.user_id, p.thumbnail_url")
	mock.ExpectQuery("SELECT p.id, p.title, p.user_id, p.thumbnail_url").
		WithArgs(request.TrackID, request.UserID).
		WillReturnRows(rows)

	response, err := repo.GetPlaylistsToAdd(ctx, request)
	assert.Error(t, err)
	assert.Nil(t, response)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCheckExistsPlaylistAndNotDifferentUser(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println("Error closing database:", zap.Error(err))
		}
	}()

	repo := NewPlaylistPostgresRepository(db, metrics.NewMockMetrics())
	playlistID := int64(1)
	userID := int64(2)

	rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)

	mock.ExpectPrepare("SELECT EXISTS")
	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(playlistID, userID).
		WillReturnRows(rows)

	exists, err := repo.CheckExistsPlaylistAndNotDifferentUser(ctx, playlistID, userID)
	assert.NoError(t, err)
	assert.True(t, exists)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCheckExistsPlaylistAndNotDifferentUserError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println("Error closing database:", zap.Error(err))
		}
	}()

	repo := NewPlaylistPostgresRepository(db, metrics.NewMockMetrics())
	playlistID := int64(1)
	userID := int64(2)

	mock.ExpectPrepare("SELECT EXISTS")
	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(playlistID, userID).
		WillReturnError(stderrors.New("db error"))

	exists, err := repo.CheckExistsPlaylistAndNotDifferentUser(ctx, playlistID, userID)
	assert.Error(t, err)
	assert.False(t, exists)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLikePlaylist(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println("Error closing database:", zap.Error(err))
		}
	}()

	repo := NewPlaylistPostgresRepository(db, metrics.NewMockMetrics())
	request := &repoModel.LikePlaylistRequest{
		PlaylistID: 1,
		UserID:     2,
	}

	existsRows := sqlmock.NewRows([]string{"exists"}).AddRow(true)

	mock.ExpectPrepare("INSERT INTO favorite_playlist")
	mock.ExpectPrepare("SELECT EXISTS")
	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(request.PlaylistID, request.UserID).
		WillReturnRows(existsRows)

	mock.ExpectExec("INSERT INTO favorite_playlist").
		WithArgs(request.UserID, request.PlaylistID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.LikePlaylist(ctx, request)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLikePlaylistCheckExistsError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println("Error closing database:", zap.Error(err))
		}
	}()

	repo := NewPlaylistPostgresRepository(db, metrics.NewMockMetrics())
	request := &repoModel.LikePlaylistRequest{
		PlaylistID: 1,
		UserID:     2,
	}

	mock.ExpectPrepare("INSERT INTO favorite_playlist")
	mock.ExpectPrepare("SELECT EXISTS")
	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(request.PlaylistID, request.UserID).
		WillReturnError(stderrors.New("db error"))

	err := repo.LikePlaylist(ctx, request)
	assert.Error(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLikePlaylistInsertError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println("Error closing database:", zap.Error(err))
		}
	}()

	repo := NewPlaylistPostgresRepository(db, metrics.NewMockMetrics())
	request := &repoModel.LikePlaylistRequest{
		PlaylistID: 1,
		UserID:     2,
	}

	existsRows := sqlmock.NewRows([]string{"exists"}).AddRow(true)

	mock.ExpectPrepare("INSERT INTO favorite_playlist")
	mock.ExpectPrepare("SELECT EXISTS")
	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(request.PlaylistID, request.UserID).
		WillReturnRows(existsRows)

	mock.ExpectExec("INSERT INTO favorite_playlist").
		WithArgs(request.UserID, request.PlaylistID).
		WillReturnError(stderrors.New("db error"))

	err := repo.LikePlaylist(ctx, request)
	assert.Error(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUnlikePlaylist(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println("Error closing database:", zap.Error(err))
		}
	}()

	repo := NewPlaylistPostgresRepository(db, metrics.NewMockMetrics())
	request := &repoModel.LikePlaylistRequest{
		PlaylistID: 1,
		UserID:     2,
	}

	mock.ExpectPrepare("DELETE FROM favorite_playlist")
	mock.ExpectExec("DELETE FROM favorite_playlist").
		WithArgs(request.UserID, request.PlaylistID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.UnlikePlaylist(ctx, request)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUnlikePlaylistError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println("Error closing database:", zap.Error(err))
		}
	}()

	repo := NewPlaylistPostgresRepository(db, metrics.NewMockMetrics())
	request := &repoModel.LikePlaylistRequest{
		PlaylistID: 1,
		UserID:     2,
	}

	mock.ExpectPrepare("DELETE FROM favorite_playlist")
	mock.ExpectExec("DELETE FROM favorite_playlist").
		WithArgs(request.UserID, request.PlaylistID).
		WillReturnError(stderrors.New("db error"))

	err := repo.UnlikePlaylist(ctx, request)
	assert.Error(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPlaylistWithIsLikedByID(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println("Error closing database:", zap.Error(err))
		}
	}()

	repo := NewPlaylistPostgresRepository(db, metrics.NewMockMetrics())
	playlistID := int64(1)
	userID := int64(2)

	rows := sqlmock.NewRows([]string{"id", "title", "user_id", "thumbnail_url", "is_liked"}).
		AddRow(1, "Playlist 1", 1, "thumbnail1.jpg", true)

	mock.ExpectPrepare("SELECT p.id, p.title, p.user_id, p.thumbnail_url")
	mock.ExpectQuery("SELECT p.id, p.title, p.user_id, p.thumbnail_url").
		WithArgs(playlistID, userID).
		WillReturnRows(rows)

	playlist, err := repo.GetPlaylistWithIsLikedByID(ctx, playlistID, userID)
	assert.NoError(t, err)
	assert.NotNil(t, playlist)
	assert.Equal(t, int64(1), playlist.Playlist.ID)
	assert.Equal(t, "Playlist 1", playlist.Playlist.Title)
	assert.Equal(t, int64(1), playlist.Playlist.UserID)
	assert.Equal(t, "thumbnail1.jpg", playlist.Playlist.Thumbnail)
	assert.True(t, playlist.IsLiked)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPlaylistWithIsLikedByIDError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println("Error closing database:", zap.Error(err))
		}
	}()

	repo := NewPlaylistPostgresRepository(db, metrics.NewMockMetrics())
	playlistID := int64(1)
	userID := int64(2)

	mock.ExpectPrepare("SELECT p.id, p.title, p.user_id, p.thumbnail_url")
	mock.ExpectQuery("SELECT p.id, p.title, p.user_id, p.thumbnail_url").
		WithArgs(playlistID, userID).
		WillReturnError(stderrors.New("db error"))

	playlist, err := repo.GetPlaylistWithIsLikedByID(ctx, playlistID, userID)
	assert.Error(t, err)
	assert.Nil(t, playlist)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetProfilePlaylists(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println("Error closing database:", zap.Error(err))
		}
	}()

	repo := NewPlaylistPostgresRepository(db, metrics.NewMockMetrics())
	request := &repoModel.GetProfilePlaylistsRequest{
		UserID: 1,
	}

	rows := sqlmock.NewRows([]string{"id", "title", "user_id", "thumbnail_url"}).
		AddRow(1, "Playlist 1", 1, "thumbnail1.jpg").
		AddRow(2, "Playlist 2", 1, "thumbnail2.jpg")

	mock.ExpectPrepare("SELECT p.id, p.title, p.user_id, p.thumbnail_url")
	mock.ExpectQuery("SELECT p.id, p.title, p.user_id, p.thumbnail_url").
		WithArgs(request.UserID).
		WillReturnRows(rows)

	response, err := repo.GetProfilePlaylists(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Len(t, response.Playlists, 2)
	assert.Equal(t, int64(1), response.Playlists[0].ID)
	assert.Equal(t, "Playlist 1", response.Playlists[0].Title)
	assert.Equal(t, int64(1), response.Playlists[0].UserID)
	assert.Equal(t, "thumbnail1.jpg", response.Playlists[0].Thumbnail)
	assert.Equal(t, int64(2), response.Playlists[1].ID)
	assert.Equal(t, "Playlist 2", response.Playlists[1].Title)
	assert.Equal(t, int64(1), response.Playlists[1].UserID)
	assert.Equal(t, "thumbnail2.jpg", response.Playlists[1].Thumbnail)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetProfilePlaylistsError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println("Error closing database:", zap.Error(err))
		}
	}()

	repo := NewPlaylistPostgresRepository(db, metrics.NewMockMetrics())
	request := &repoModel.GetProfilePlaylistsRequest{
		UserID: 1,
	}

	mock.ExpectPrepare("SELECT p.id, p.title, p.user_id, p.thumbnail_url")
	mock.ExpectQuery("SELECT p.id, p.title, p.user_id, p.thumbnail_url").
		WithArgs(request.UserID).
		WillReturnError(stderrors.New("db error"))

	response, err := repo.GetProfilePlaylists(ctx, request)
	assert.Error(t, err)
	assert.Nil(t, response)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetProfilePlaylistsScanError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println("Error closing database:", zap.Error(err))
		}
	}()

	repo := NewPlaylistPostgresRepository(db, metrics.NewMockMetrics())
	request := &repoModel.GetProfilePlaylistsRequest{
		UserID: 1,
	}

	rows := sqlmock.NewRows([]string{"id", "title", "user_id", "thumbnail_url"}).
		AddRow("invalid_id", "Playlist 1", 1, "thumbnail1.jpg")

	mock.ExpectPrepare("SELECT p.id, p.title, p.user_id, p.thumbnail_url")
	mock.ExpectQuery("SELECT p.id, p.title, p.user_id, p.thumbnail_url").
		WithArgs(request.UserID).
		WillReturnRows(rows)

	response, err := repo.GetProfilePlaylists(ctx, request)
	assert.Error(t, err)
	assert.Nil(t, response)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSearchPlaylists(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println("Error closing database:", zap.Error(err))
		}
	}()

	repo := NewPlaylistPostgresRepository(db, metrics.NewMockMetrics())
	request := &repoModel.SearchPlaylistsRequest{
		Query:  "test playlist",
		UserID: 1,
	}

	rows := sqlmock.NewRows([]string{"id", "title", "user_id", "thumbnail_url"}).
		AddRow(1, "Test Playlist", 1, "thumbnail1.jpg").
		AddRow(2, "Playlist Test", 2, "thumbnail2.jpg")

	mock.ExpectPrepare("SELECT id, title, user_id, thumbnail_url")
	mock.ExpectQuery("SELECT id, title, user_id, thumbnail_url").
		WithArgs("test:* & playlist:*", request.UserID, request.Query).
		WillReturnRows(rows)

	playlists, err := repo.SearchPlaylists(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, playlists)
	assert.Len(t, playlists.Playlists, 2)
	assert.Equal(t, int64(1), playlists.Playlists[0].ID)
	assert.Equal(t, "Test Playlist", playlists.Playlists[0].Title)
	assert.Equal(t, int64(1), playlists.Playlists[0].UserID)
	assert.Equal(t, int64(2), playlists.Playlists[1].ID)
	assert.Equal(t, "Playlist Test", playlists.Playlists[1].Title)
	assert.Equal(t, int64(2), playlists.Playlists[1].UserID)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSearchPlaylistsError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println("Error closing database:", zap.Error(err))
		}
	}()

	repo := NewPlaylistPostgresRepository(db, metrics.NewMockMetrics())
	request := &repoModel.SearchPlaylistsRequest{
		Query:  "test playlist",
		UserID: 1,
	}

	mock.ExpectPrepare("SELECT id, title, user_id, thumbnail_url")
	mock.ExpectQuery("SELECT id, title, user_id, thumbnail_url").
		WithArgs("test:* & playlist:*", request.UserID, request.Query).
		WillReturnError(stderrors.New("db error"))

	playlists, err := repo.SearchPlaylists(ctx, request)
	assert.Error(t, err)
	assert.Nil(t, playlists)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSearchPlaylistsScanError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println("Error closing database:", zap.Error(err))
		}
	}()

	repo := NewPlaylistPostgresRepository(db, metrics.NewMockMetrics())
	request := &repoModel.SearchPlaylistsRequest{
		Query:  "test playlist",
		UserID: 1,
	}

	rows := sqlmock.NewRows([]string{"id", "title", "user_id", "thumbnail_url"}).
		AddRow("invalid_id", "Test Playlist", 1, "thumbnail1.jpg")

	mock.ExpectPrepare("SELECT id, title, user_id, thumbnail_url")
	mock.ExpectQuery("SELECT id, title, user_id, thumbnail_url").
		WithArgs("test:* & playlist:*", request.UserID, request.Query).
		WillReturnRows(rows)

	playlists, err := repo.SearchPlaylists(ctx, request)
	assert.Error(t, err)
	assert.Nil(t, playlists)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreatePlaylistGetPlaylistError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println("Error closing database:", zap.Error(err))
		}
	}()

	repo := NewPlaylistPostgresRepository(db, metrics.NewMockMetrics())
	request := &repoModel.CreatePlaylistRequest{
		Title:     "New Playlist",
		UserID:    1,
		Thumbnail: "thumbnail.jpg",
		IsPublic:  true,
	}

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	mock.ExpectPrepare("INSERT INTO playlist")
	mock.ExpectQuery("INSERT INTO playlist").
		WithArgs(request.Title, request.UserID, request.Thumbnail, request.IsPublic).
		WillReturnRows(rows)

	mock.ExpectPrepare("SELECT id, title, user_id, thumbnail_url, is_public")
	mock.ExpectQuery("SELECT id, title, user_id, thumbnail_url, is_public").
		WithArgs(int64(1)).
		WillReturnError(stderrors.New("db error"))

	playlist, err := repo.CreatePlaylist(ctx, request)
	assert.Error(t, err)
	assert.Nil(t, playlist)

	assert.NoError(t, mock.ExpectationsWereMet())
}
