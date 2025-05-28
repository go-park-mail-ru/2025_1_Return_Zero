package repository

import (
	"context"
	"database/sql"
	stderrors "errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	loggerPkg "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/metrics"
	trackErrors "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/track/model/errors"
	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/track/model/repository"
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

func TestGetAllTracks(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewTrackPostgresRepository(db, metrics.NewMockMetrics())
	filters := &repoModel.TrackFilters{
		Pagination: &repoModel.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}
	userID := int64(1)

	rows := sqlmock.NewRows([]string{"id", "title", "thumbnail_url", "duration", "album_id", "is_favorite"}).
		AddRow(1, "Track 1", "thumbnail1.jpg", 200, 1, true).
		AddRow(2, "Track 2", "thumbnail2.jpg", 200, 1, false)

	mock.ExpectPrepare("SELECT t.id, t.title, t.thumbnail_url, t.duration, t.album_id")
	mock.ExpectQuery("SELECT t.id, t.title, t.thumbnail_url, t.duration, t.album_id").
		WithArgs(filters.Pagination.Limit, filters.Pagination.Offset, userID).
		WillReturnRows(rows)

	tracks, err := repo.GetAllTracks(ctx, filters, userID)
	assert.NoError(t, err)
	assert.Len(t, tracks, 2)
	assert.Equal(t, int64(1), tracks[0].ID)
	assert.Equal(t, "Track 1", tracks[0].Title)
	assert.Equal(t, "thumbnail1.jpg", tracks[0].Thumbnail)
	assert.Equal(t, int64(200), tracks[0].Duration)
	assert.Equal(t, int64(1), tracks[0].AlbumID)
	assert.True(t, tracks[0].IsFavorite)

	assert.Equal(t, int64(2), tracks[1].ID)
	assert.Equal(t, "Track 2", tracks[1].Title)
	assert.Equal(t, "thumbnail2.jpg", tracks[1].Thumbnail)
	assert.Equal(t, int64(200), tracks[1].Duration)
	assert.Equal(t, int64(1), tracks[1].AlbumID)
	assert.False(t, tracks[1].IsFavorite)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllTracksError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewTrackPostgresRepository(db, metrics.NewMockMetrics())
	filters := &repoModel.TrackFilters{
		Pagination: &repoModel.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}
	userID := int64(1)

	mock.ExpectPrepare("SELECT t.id, t.title, t.thumbnail_url, t.duration, t.album_id")
	mock.ExpectQuery("SELECT t.id, t.title, t.thumbnail_url, t.duration, t.album_id").
		WithArgs(filters.Pagination.Limit, filters.Pagination.Offset, userID).
		WillReturnError(stderrors.New("db error"))

	tracks, err := repo.GetAllTracks(ctx, filters, userID)
	assert.Error(t, err)
	assert.Nil(t, tracks)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTrackByID(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewTrackPostgresRepository(db, metrics.NewMockMetrics())
	trackID := int64(1)
	userID := int64(1)

	rows := sqlmock.NewRows([]string{"id", "title", "thumbnail_url", "duration", "album_id", "file_url", "is_favorite"}).
		AddRow(1, "Track 1", "thumbnail1.jpg", 200, 1, "file_key.mp3", true)

	mock.ExpectPrepare("SELECT t.id, t.title, t.thumbnail_url, t.duration, t.album_id, t.file_url")
	mock.ExpectQuery("SELECT t.id, t.title, t.thumbnail_url, t.duration, t.album_id, t.file_url").
		WithArgs(trackID, userID).
		WillReturnRows(rows)

	track, err := repo.GetTrackByID(ctx, trackID, userID)
	assert.NoError(t, err)
	assert.NotNil(t, track)
	assert.Equal(t, int64(1), track.ID)
	assert.Equal(t, "Track 1", track.Title)
	assert.Equal(t, "thumbnail1.jpg", track.Thumbnail)
	assert.Equal(t, int64(200), track.Duration)
	assert.Equal(t, int64(1), track.AlbumID)
	assert.Equal(t, "file_key.mp3", track.FileKey)
	assert.True(t, track.IsFavorite)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTrackByIDNotFound(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewTrackPostgresRepository(db, metrics.NewMockMetrics())
	trackID := int64(1)
	userID := int64(1)

	mock.ExpectPrepare("SELECT t.id, t.title, t.thumbnail_url, t.duration, t.album_id, t.file_url")
	mock.ExpectQuery("SELECT t.id, t.title, t.thumbnail_url, t.duration, t.album_id, t.file_url").
		WithArgs(trackID, userID).
		WillReturnError(sql.ErrNoRows)

	track, err := repo.GetTrackByID(ctx, trackID, userID)
	assert.Error(t, err)
	assert.Equal(t, trackErrors.ErrTrackNotFound, err)
	assert.Nil(t, track)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateStream(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewTrackPostgresRepository(db, metrics.NewMockMetrics())
	createData := &repoModel.TrackStreamCreateData{
		TrackID: 1,
		UserID:  1,
	}

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	mock.ExpectPrepare("INSERT INTO track_stream")
	mock.ExpectQuery("INSERT INTO track_stream").
		WithArgs(createData.TrackID, createData.UserID).
		WillReturnRows(rows)

	streamID, err := repo.CreateStream(ctx, createData)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), streamID)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateStreamError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewTrackPostgresRepository(db, metrics.NewMockMetrics())
	createData := &repoModel.TrackStreamCreateData{
		TrackID: 1,
		UserID:  1,
	}

	mock.ExpectPrepare("INSERT INTO track_stream")
	mock.ExpectQuery("INSERT INTO track_stream").
		WithArgs(createData.TrackID, createData.UserID).
		WillReturnError(stderrors.New("db error"))

	streamID, err := repo.CreateStream(ctx, createData)
	assert.Error(t, err)
	assert.Equal(t, int64(0), streamID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetStreamByID(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewTrackPostgresRepository(db, metrics.NewMockMetrics())
	streamID := int64(1)

	rows := sqlmock.NewRows([]string{"id", "user_id", "track_id", "duration"}).
		AddRow(1, 1, 1, 200)

	mock.ExpectPrepare("SELECT id, user_id, track_id, duration")
	mock.ExpectQuery("SELECT id, user_id, track_id, duration").
		WithArgs(streamID).
		WillReturnRows(rows)

	stream, err := repo.GetStreamByID(ctx, streamID)
	assert.NoError(t, err)
	assert.NotNil(t, stream)
	assert.Equal(t, int64(1), stream.ID)
	assert.Equal(t, int64(1), stream.UserID)
	assert.Equal(t, int64(1), stream.TrackID)
	assert.Equal(t, int64(200), stream.Duration)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetStreamByIDNotFound(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewTrackPostgresRepository(db, metrics.NewMockMetrics())
	streamID := int64(1)

	mock.ExpectPrepare("SELECT id, user_id, track_id, duration")
	mock.ExpectQuery("SELECT id, user_id, track_id, duration").
		WithArgs(streamID).
		WillReturnError(sql.ErrNoRows)

	stream, err := repo.GetStreamByID(ctx, streamID)
	assert.Error(t, err)
	assert.Equal(t, trackErrors.ErrStreamNotFound, err)
	assert.Nil(t, stream)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateStreamDuration(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewTrackPostgresRepository(db, metrics.NewMockMetrics())
	updateData := &repoModel.TrackStreamUpdateData{
		StreamID: 1,
		Duration: 200,
	}

	mock.ExpectPrepare("UPDATE track_stream")
	mock.ExpectExec("UPDATE track_stream").
		WithArgs(updateData.Duration, updateData.StreamID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.UpdateStreamDuration(ctx, updateData)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateStreamDurationNotFound(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewTrackPostgresRepository(db, metrics.NewMockMetrics())
	updateData := &repoModel.TrackStreamUpdateData{
		StreamID: 1,
		Duration: 180,
	}

	mock.ExpectPrepare("UPDATE track_stream")
	mock.ExpectExec("UPDATE track_stream").
		WithArgs(updateData.Duration, updateData.StreamID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.UpdateStreamDuration(ctx, updateData)
	assert.Error(t, err)
	assert.Equal(t, trackErrors.ErrFailedToUpdateStreamDuration, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetStreamsByUserID(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewTrackPostgresRepository(db, metrics.NewMockMetrics())
	userID := int64(1)
	filters := &repoModel.TrackFilters{
		Pagination: &repoModel.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}

	rows := sqlmock.NewRows([]string{"id", "user_id", "track_id", "duration"}).
		AddRow(1, 1, 1, 200).
		AddRow(2, 1, 2, 200)

	mock.ExpectPrepare("WITH latest_streams AS")
	mock.ExpectQuery("WITH latest_streams AS").
		WithArgs(userID, filters.Pagination.Limit, filters.Pagination.Offset).
		WillReturnRows(rows)

	streams, err := repo.GetStreamsByUserID(ctx, userID, filters)
	assert.NoError(t, err)
	assert.Len(t, streams, 2)
	assert.Equal(t, int64(1), streams[0].ID)
	assert.Equal(t, int64(1), streams[0].UserID)
	assert.Equal(t, int64(1), streams[0].TrackID)
	assert.Equal(t, int64(200), streams[0].Duration)
	assert.Equal(t, int64(2), streams[1].ID)
	assert.Equal(t, int64(1), streams[1].UserID)
	assert.Equal(t, int64(2), streams[1].TrackID)
	assert.Equal(t, int64(200), streams[1].Duration)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetStreamsByUserIDError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewTrackPostgresRepository(db, metrics.NewMockMetrics())
	userID := int64(1)
	filters := &repoModel.TrackFilters{
		Pagination: &repoModel.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}

	mock.ExpectPrepare("WITH latest_streams AS")
	mock.ExpectQuery("WITH latest_streams AS").
		WithArgs(userID, filters.Pagination.Limit, filters.Pagination.Offset).
		WillReturnError(stderrors.New("db error"))

	streams, err := repo.GetStreamsByUserID(ctx, userID, filters)
	assert.Error(t, err)
	assert.Nil(t, streams)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTracksByIDs(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewTrackPostgresRepository(db, metrics.NewMockMetrics())
	trackIDs := []int64{1, 2}
	userID := int64(1)

	rows := sqlmock.NewRows([]string{"id", "title", "thumbnail_url", "duration", "album_id", "is_favorite"}).
		AddRow(1, "Track 1", "thumbnail1.jpg", 200, 1, true).
		AddRow(2, "Track 2", "thumbnail2.jpg", 200, 1, false)

	mock.ExpectPrepare("SELECT t.id, t.title, t.thumbnail_url, t.duration, t.album_id")
	mock.ExpectQuery("SELECT t.id, t.title, t.thumbnail_url, t.duration, t.album_id").
		WithArgs(pq.Array(trackIDs), userID).
		WillReturnRows(rows)

	tracks, err := repo.GetTracksByIDs(ctx, trackIDs, userID)
	assert.NoError(t, err)
	assert.Len(t, tracks, 2)
	assert.Equal(t, int64(1), tracks[1].ID)
	assert.Equal(t, "Track 1", tracks[1].Title)
	assert.Equal(t, int64(2), tracks[2].ID)
	assert.Equal(t, "Track 2", tracks[2].Title)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTracksByIDsError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewTrackPostgresRepository(db, metrics.NewMockMetrics())
	trackIDs := []int64{1, 2}
	userID := int64(1)

	mock.ExpectPrepare("SELECT t.id, t.title, t.thumbnail_url, t.duration, t.album_id")
	mock.ExpectQuery("SELECT t.id, t.title, t.thumbnail_url, t.duration, t.album_id").
		WithArgs(pq.Array(trackIDs), userID).
		WillReturnError(stderrors.New("db error"))

	tracks, err := repo.GetTracksByIDs(ctx, trackIDs, userID)
	assert.Error(t, err)
	assert.Nil(t, tracks)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTracksByIDsFiltered(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewTrackPostgresRepository(db, metrics.NewMockMetrics())
	trackIDs := []int64{1, 2}
	filters := &repoModel.TrackFilters{
		Pagination: &repoModel.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}
	userID := int64(1)

	rows := sqlmock.NewRows([]string{"id", "title", "thumbnail_url", "duration", "album_id", "is_favorite"}).
		AddRow(1, "Track 1", "thumbnail1.jpg", 200, 1, true).
		AddRow(2, "Track 2", "thumbnail2.jpg", 200, 1, false)

	mock.ExpectPrepare("SELECT t.id, t.title, t.thumbnail_url, t.duration, t.album_id")
	mock.ExpectQuery("SELECT t.id, t.title, t.thumbnail_url, t.duration, t.album_id").
		WithArgs(pq.Array(trackIDs), filters.Pagination.Limit, filters.Pagination.Offset, userID).
		WillReturnRows(rows)

	tracks, err := repo.GetTracksByIDsFiltered(ctx, trackIDs, filters, userID)
	assert.NoError(t, err)
	assert.Len(t, tracks, 2)
	assert.Equal(t, int64(1), tracks[0].ID)
	assert.Equal(t, "Track 1", tracks[0].Title)
	assert.Equal(t, int64(2), tracks[1].ID)
	assert.Equal(t, "Track 2", tracks[1].Title)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTracksByIDsFilteredError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewTrackPostgresRepository(db, metrics.NewMockMetrics())
	trackIDs := []int64{1, 2}
	filters := &repoModel.TrackFilters{
		Pagination: &repoModel.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}
	userID := int64(1)

	mock.ExpectPrepare("SELECT t.id, t.title, t.thumbnail_url, t.duration, t.album_id")
	mock.ExpectQuery("SELECT t.id, t.title, t.thumbnail_url, t.duration, t.album_id").
		WithArgs(pq.Array(trackIDs), filters.Pagination.Limit, filters.Pagination.Offset, userID).
		WillReturnError(stderrors.New("db error"))

	tracks, err := repo.GetTracksByIDsFiltered(ctx, trackIDs, filters, userID)
	assert.Error(t, err)
	assert.Nil(t, tracks)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAlbumIDByTrackID(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewTrackPostgresRepository(db, metrics.NewMockMetrics())
	trackID := int64(1)

	rows := sqlmock.NewRows([]string{"album_id"}).AddRow(1)

	mock.ExpectPrepare("SELECT album_id")
	mock.ExpectQuery("SELECT album_id").
		WithArgs(trackID).
		WillReturnRows(rows)

	albumID, err := repo.GetAlbumIDByTrackID(ctx, trackID)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), albumID)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAlbumIDByTrackIDError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewTrackPostgresRepository(db, metrics.NewMockMetrics())
	trackID := int64(1)

	mock.ExpectPrepare("SELECT album_id")
	mock.ExpectQuery("SELECT album_id").
		WithArgs(trackID).
		WillReturnError(stderrors.New("db error"))

	albumID, err := repo.GetAlbumIDByTrackID(ctx, trackID)
	assert.Error(t, err)
	assert.Equal(t, int64(0), albumID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTracksByAlbumID(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewTrackPostgresRepository(db, metrics.NewMockMetrics())
	albumID := int64(1)
	userID := int64(1)

	rows := sqlmock.NewRows([]string{"id", "title", "thumbnail_url", "duration", "album_id", "is_favorite"}).
		AddRow(1, "Track 1", "thumbnail1.jpg", 200, 1, true).
		AddRow(2, "Track 2", "thumbnail2.jpg", 200, 1, false)

	mock.ExpectPrepare("SELECT t.id, t.title, t.thumbnail_url, t.duration, t.album_id")
	mock.ExpectQuery("SELECT t.id, t.title, t.thumbnail_url, t.duration, t.album_id").
		WithArgs(albumID, userID).
		WillReturnRows(rows)

	tracks, err := repo.GetTracksByAlbumID(ctx, albumID, userID)
	assert.NoError(t, err)
	assert.Len(t, tracks, 2)
	assert.Equal(t, int64(1), tracks[0].ID)
	assert.Equal(t, "Track 1", tracks[0].Title)
	assert.Equal(t, int64(2), tracks[1].ID)
	assert.Equal(t, "Track 2", tracks[1].Title)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTracksByAlbumIDError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewTrackPostgresRepository(db, metrics.NewMockMetrics())
	albumID := int64(1)
	userID := int64(1)

	mock.ExpectPrepare("SELECT t.id, t.title, t.thumbnail_url, t.duration, t.album_id")
	mock.ExpectQuery("SELECT t.id, t.title, t.thumbnail_url, t.duration, t.album_id").
		WithArgs(albumID, userID).
		WillReturnError(stderrors.New("db error"))

	tracks, err := repo.GetTracksByAlbumID(ctx, albumID, userID)
	assert.Error(t, err)
	assert.Nil(t, tracks)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetMinutesListenedByUserID(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewTrackPostgresRepository(db, metrics.NewMockMetrics())
	userID := int64(1)

	rows := sqlmock.NewRows([]string{"minutes"}).AddRow(1)

	mock.ExpectPrepare("SELECT COALESCE")
	mock.ExpectQuery("SELECT COALESCE").
		WithArgs(userID).
		WillReturnRows(rows)

	minutes, err := repo.GetMinutesListenedByUserID(ctx, userID)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), minutes)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetMinutesListenedByUserIDError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewTrackPostgresRepository(db, metrics.NewMockMetrics())
	userID := int64(1)

	mock.ExpectPrepare("SELECT COALESCE")
	mock.ExpectQuery("SELECT COALESCE").
		WithArgs(userID).
		WillReturnError(stderrors.New("db error"))

	minutes, err := repo.GetMinutesListenedByUserID(ctx, userID)
	assert.Error(t, err)
	assert.Equal(t, int64(0), minutes)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTracksListenedByUserID(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewTrackPostgresRepository(db, metrics.NewMockMetrics())
	userID := int64(1)

	rows := sqlmock.NewRows([]string{"count"}).AddRow(1)

	mock.ExpectPrepare("SELECT COUNT")
	mock.ExpectQuery("SELECT COUNT").
		WithArgs(userID).
		WillReturnRows(rows)

	count, err := repo.GetTracksListenedByUserID(ctx, userID)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTracksListenedByUserIDError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewTrackPostgresRepository(db, metrics.NewMockMetrics())
	userID := int64(1)

	mock.ExpectPrepare("SELECT COUNT")
	mock.ExpectQuery("SELECT COUNT").
		WithArgs(userID).
		WillReturnError(stderrors.New("db error"))

	count, err := repo.GetTracksListenedByUserID(ctx, userID)
	assert.Error(t, err)
	assert.Equal(t, int64(0), count)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCheckTrackExists(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewTrackPostgresRepository(db, metrics.NewMockMetrics())
	trackID := int64(1)

	rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)

	mock.ExpectPrepare("SELECT EXISTS")
	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(trackID).
		WillReturnRows(rows)

	exists, err := repo.CheckTrackExists(ctx, trackID)
	assert.NoError(t, err)
	assert.True(t, exists)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCheckTrackExistsError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewTrackPostgresRepository(db, metrics.NewMockMetrics())
	trackID := int64(1)

	mock.ExpectPrepare("SELECT EXISTS")
	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(trackID).
		WillReturnError(stderrors.New("db error"))

	exists, err := repo.CheckTrackExists(ctx, trackID)
	assert.Error(t, err)
	assert.False(t, exists)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLikeTrack(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewTrackPostgresRepository(db, metrics.NewMockMetrics())
	likeRequest := &repoModel.LikeRequest{
		TrackID: 1,
		UserID:  1,
	}

	mock.ExpectPrepare("INSERT INTO favorite_track")
	mock.ExpectExec("INSERT INTO favorite_track").
		WithArgs(likeRequest.TrackID, likeRequest.UserID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.LikeTrack(ctx, likeRequest)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLikeTrackError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewTrackPostgresRepository(db, metrics.NewMockMetrics())
	likeRequest := &repoModel.LikeRequest{
		TrackID: 1,
		UserID:  1,
	}

	mock.ExpectPrepare("INSERT INTO favorite_track")
	mock.ExpectExec("INSERT INTO favorite_track").
		WithArgs(likeRequest.TrackID, likeRequest.UserID).
		WillReturnError(stderrors.New("db error"))

	err := repo.LikeTrack(ctx, likeRequest)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUnlikeTrack(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewTrackPostgresRepository(db, metrics.NewMockMetrics())
	likeRequest := &repoModel.LikeRequest{
		TrackID: 1,
		UserID:  1,
	}

	mock.ExpectPrepare("DELETE FROM favorite_track")
	mock.ExpectExec("DELETE FROM favorite_track").
		WithArgs(likeRequest.TrackID, likeRequest.UserID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.UnlikeTrack(ctx, likeRequest)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUnlikeTrackError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewTrackPostgresRepository(db, metrics.NewMockMetrics())
	likeRequest := &repoModel.LikeRequest{
		TrackID: 1,
		UserID:  1,
	}

	mock.ExpectPrepare("DELETE FROM favorite_track")
	mock.ExpectExec("DELETE FROM favorite_track").
		WithArgs(likeRequest.TrackID, likeRequest.UserID).
		WillReturnError(stderrors.New("db error"))

	err := repo.UnlikeTrack(ctx, likeRequest)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetFavoriteTracks(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewTrackPostgresRepository(db, metrics.NewMockMetrics())
	favoriteRequest := &repoModel.FavoriteRequest{
		RequestUserID: 1,
		ProfileUserID: 2,
		Filters: &repoModel.TrackFilters{
			Pagination: &repoModel.Pagination{
				Limit:  10,
				Offset: 0,
			},
		},
	}

	rows := sqlmock.NewRows([]string{"id", "title", "thumbnail_url", "duration", "album_id", "is_favorite"}).
		AddRow(1, "Track 1", "thumbnail1.jpg", 200, 1, true).
		AddRow(2, "Track 2", "thumbnail2.jpg", 200, 1, false)

	mock.ExpectPrepare("SELECT t.id, t.title, t.thumbnail_url, t.duration, t.album_id")
	mock.ExpectQuery("SELECT t.id, t.title, t.thumbnail_url, t.duration, t.album_id").
		WithArgs(favoriteRequest.RequestUserID, favoriteRequest.ProfileUserID, favoriteRequest.Filters.Pagination.Limit, favoriteRequest.Filters.Pagination.Offset).
		WillReturnRows(rows)

	tracks, err := repo.GetFavoriteTracks(ctx, favoriteRequest)
	assert.NoError(t, err)
	assert.Len(t, tracks, 2)
	assert.Equal(t, int64(1), tracks[0].ID)
	assert.Equal(t, "Track 1", tracks[0].Title)
	assert.Equal(t, int64(2), tracks[1].ID)
	assert.Equal(t, "Track 2", tracks[1].Title)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetFavoriteTracksError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewTrackPostgresRepository(db, metrics.NewMockMetrics())
	favoriteRequest := &repoModel.FavoriteRequest{
		RequestUserID: 1,
		ProfileUserID: 2,
		Filters: &repoModel.TrackFilters{
			Pagination: &repoModel.Pagination{
				Limit:  10,
				Offset: 0,
			},
		},
	}

	mock.ExpectPrepare("SELECT t.id, t.title, t.thumbnail_url, t.duration, t.album_id")
	mock.ExpectQuery("SELECT t.id, t.title, t.thumbnail_url, t.duration, t.album_id").
		WithArgs(favoriteRequest.RequestUserID, favoriteRequest.ProfileUserID, favoriteRequest.Filters.Pagination.Limit, favoriteRequest.Filters.Pagination.Offset).
		WillReturnError(stderrors.New("db error"))

	tracks, err := repo.GetFavoriteTracks(ctx, favoriteRequest)
	assert.Error(t, err)
	assert.Nil(t, tracks)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSearchTracks(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewTrackPostgresRepository(db, metrics.NewMockMetrics())
	query := "test track"
	userID := int64(1)

	rows := sqlmock.NewRows([]string{"id", "title", "thumbnail_url", "duration", "album_id", "is_favorite"}).
		AddRow(1, "Test Track", "thumbnail1.jpg", 200, 1, true).
		AddRow(2, "Track Test", "thumbnail2.jpg", 200, 1, false)

	mock.ExpectPrepare("SELECT t.id, t.title, t.thumbnail_url, t.duration, t.album_id")
	mock.ExpectQuery("SELECT t.id, t.title, t.thumbnail_url, t.duration, t.album_id").
		WithArgs("test:* & track:*", userID, query).
		WillReturnRows(rows)

	tracks, err := repo.SearchTracks(ctx, query, userID)
	assert.NoError(t, err)
	assert.Len(t, tracks, 2)
	assert.Equal(t, int64(1), tracks[0].ID)
	assert.Equal(t, "Test Track", tracks[0].Title)
	assert.Equal(t, int64(2), tracks[1].ID)
	assert.Equal(t, "Track Test", tracks[1].Title)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSearchTracksError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewTrackPostgresRepository(db, metrics.NewMockMetrics())
	query := "test track"
	userID := int64(1)

	mock.ExpectPrepare("SELECT t.id, t.title, t.thumbnail_url, t.duration, t.album_id")
	mock.ExpectQuery("SELECT t.id, t.title, t.thumbnail_url, t.duration, t.album_id").
		WithArgs("test:* & track:*", userID, query).
		WillReturnError(stderrors.New("db error"))

	tracks, err := repo.SearchTracks(ctx, query, userID)
	assert.Error(t, err)
	assert.Nil(t, tracks)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTracksByIDsScanError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewTrackPostgresRepository(db, metrics.NewMockMetrics())
	trackIDs := []int64{1, 2}
	userID := int64(1)

	rows := sqlmock.NewRows([]string{"id", "title", "thumbnail_url", "duration", "album_id", "is_favorite"}).
		AddRow(1, "Track 1", "thumbnail1.jpg", "invalid_duration", 1, true)

	mock.ExpectPrepare("SELECT t.id, t.title, t.thumbnail_url, t.duration, t.album_id")
	mock.ExpectQuery("SELECT t.id, t.title, t.thumbnail_url, t.duration, t.album_id").
		WithArgs(pq.Array(trackIDs), userID).
		WillReturnRows(rows)

	tracks, err := repo.GetTracksByIDs(ctx, trackIDs, userID)
	assert.Error(t, err)
	assert.Nil(t, tracks)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTracksByIDsFilteredScanError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewTrackPostgresRepository(db, metrics.NewMockMetrics())
	trackIDs := []int64{1, 2}
	filters := &repoModel.TrackFilters{
		Pagination: &repoModel.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}
	userID := int64(1)

	rows := sqlmock.NewRows([]string{"id", "title", "thumbnail_url", "duration", "album_id", "is_favorite"}).
		AddRow(1, "Track 1", "thumbnail1.jpg", "invalid_duration", 1, true)

	mock.ExpectPrepare("SELECT t.id, t.title, t.thumbnail_url, t.duration, t.album_id")
	mock.ExpectQuery("SELECT t.id, t.title, t.thumbnail_url, t.duration, t.album_id").
		WithArgs(pq.Array(trackIDs), filters.Pagination.Limit, filters.Pagination.Offset, userID).
		WillReturnRows(rows)

	tracks, err := repo.GetTracksByIDsFiltered(ctx, trackIDs, filters, userID)
	assert.Error(t, err)
	assert.Nil(t, tracks)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllTracksScanError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewTrackPostgresRepository(db, metrics.NewMockMetrics())
	filters := &repoModel.TrackFilters{
		Pagination: &repoModel.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}
	userID := int64(1)

	rows := sqlmock.NewRows([]string{"id", "title", "thumbnail_url", "duration", "album_id", "is_favorite"}).
		AddRow(1, "Track 1", "thumbnail1.jpg", "invalid_duration", 1, true)

	mock.ExpectPrepare("SELECT t.id, t.title, t.thumbnail_url, t.duration, t.album_id")
	mock.ExpectQuery("SELECT t.id, t.title, t.thumbnail_url, t.duration, t.album_id").
		WithArgs(filters.Pagination.Limit, filters.Pagination.Offset, userID).
		WillReturnRows(rows)

	tracks, err := repo.GetAllTracks(ctx, filters, userID)
	assert.Error(t, err)
	assert.Nil(t, tracks)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetStreamsByUserIDScanError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewTrackPostgresRepository(db, metrics.NewMockMetrics())
	userID := int64(1)
	filters := &repoModel.TrackFilters{
		Pagination: &repoModel.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}

	rows := sqlmock.NewRows([]string{"id", "user_id", "track_id", "duration"}).
		AddRow(1, 1, 1, "invalid_duration")

	mock.ExpectPrepare("WITH latest_streams AS")
	mock.ExpectQuery("WITH latest_streams AS").
		WithArgs(userID, filters.Pagination.Limit, filters.Pagination.Offset).
		WillReturnRows(rows)

	streams, err := repo.GetStreamsByUserID(ctx, userID, filters)
	assert.Error(t, err)
	assert.Nil(t, streams)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAddTracksToAlbum(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewTrackPostgresRepository(db, metrics.NewMockMetrics())
	albumID := int64(1)

	tracks := []*repoModel.Track{
		{
			Title:     "Track 1",
			Thumbnail: "thumbnail1.jpg",
			Duration:  200,
			AlbumID:   albumID,
		},
		{
			Title:     "Track 2",
			Thumbnail: "thumbnail2.jpg",
			Duration:  180,
			AlbumID:   albumID,
		},
	}

	mock.ExpectBegin()

	mock.ExpectQuery("INSERT INTO track").
		WithArgs(
			tracks[0].Title, tracks[0].Thumbnail, tracks[0].Duration, tracks[0].AlbumID, "Track 1.mp3", 1,
			tracks[1].Title, tracks[1].Thumbnail, tracks[1].Duration, tracks[1].AlbumID, "Track 2.mp3", 2,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1).AddRow(2))

	mock.ExpectCommit()

	mock.ExpectExec("REFRESH MATERIALIZED VIEW CONCURRENTLY track_stats").
		WillReturnResult(sqlmock.NewResult(0, 0))

	trackIDs, err := repo.AddTracksToAlbum(ctx, tracks)
	assert.NoError(t, err)
	assert.Len(t, trackIDs, 2)
	assert.Equal(t, int64(1), trackIDs[0])
	assert.Equal(t, int64(2), trackIDs[1])

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAddTracksToAlbumError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewTrackPostgresRepository(db, metrics.NewMockMetrics())
	albumID := int64(1)

	tracks := []*repoModel.Track{
		{
			Title:     "Track 1",
			Thumbnail: "thumbnail1.jpg",
			Duration:  200,
			AlbumID:   albumID,
		},
		{
			Title:     "Track 2",
			Thumbnail: "thumbnail2.jpg",
			Duration:  180,
			AlbumID:   albumID,
		},
	}

	mock.ExpectBegin()

	mock.ExpectQuery("INSERT INTO track").
		WithArgs(
			tracks[0].Title, tracks[0].Thumbnail, tracks[0].Duration, tracks[0].AlbumID, "Track 1.mp3", 1,
			tracks[1].Title, tracks[1].Thumbnail, tracks[1].Duration, tracks[1].AlbumID, "Track 2.mp3", 2,
		).
		WillReturnError(stderrors.New("db error"))

	mock.ExpectRollback()

	trackIDs, err := repo.AddTracksToAlbum(ctx, tracks)
	assert.Error(t, err)
	assert.Nil(t, trackIDs)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteTracksByAlbumID(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewTrackPostgresRepository(db, metrics.NewMockMetrics())
	albumID := int64(1)

	mock.ExpectPrepare("DELETE FROM track")
	mock.ExpectExec("DELETE FROM track").
		WithArgs(albumID).
		WillReturnResult(sqlmock.NewResult(0, 2))

	err := repo.DeleteTracksByAlbumID(ctx, albumID)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteTracksByAlbumIDError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewTrackPostgresRepository(db, metrics.NewMockMetrics())
	albumID := int64(1)

	mock.ExpectPrepare("DELETE FROM track")
	mock.ExpectExec("DELETE FROM track").
		WithArgs(albumID).
		WillReturnError(stderrors.New("db error"))

	err := repo.DeleteTracksByAlbumID(ctx, albumID)
	assert.Error(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetMostLikedTracks(t *testing.T) {
    db, mock, ctx := setupTest(t)
    defer db.Close()

    repo := NewTrackPostgresRepository(db, metrics.NewMockMetrics())
    userID := int64(1)

    rows := sqlmock.NewRows([]string{"id", "title", "thumbnail_url", "duration", "album_id", "is_favorite"}).
        AddRow(1, "Most Liked Track", "thumbnail1.jpg", 200, 1, true)

    mock.ExpectQuery("SELECT t.id, t.title, t.thumbnail_url, t.duration, t.album_id, \\(ft.user_id IS NOT NULL\\) AS is_favorite").
        WithArgs(userID).
        WillReturnRows(rows)

    tracks, err := repo.GetMostLikedTracks(ctx, userID)
    assert.NoError(t, err)
    assert.Len(t, tracks, 1)
    assert.Equal(t, int64(1), tracks[0].ID)
    assert.Equal(t, "Most Liked Track", tracks[0].Title)

    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetMostLikedTracksError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewTrackPostgresRepository(db, metrics.NewMockMetrics())
	userID := int64(1)

	mock.ExpectQuery("SELECT t.id, t.title, t.thumbnail_url, t.duration, t.album_id, \\(ft.user_id IS NOT NULL\\) AS is_favorite").
		WithArgs(userID).
		WillReturnError(stderrors.New("db error"))

	tracks, err := repo.GetMostLikedTracks(ctx, userID)
	assert.Error(t, err)
	assert.Nil(t, tracks)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetMostRecentTracks(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewTrackPostgresRepository(db, metrics.NewMockMetrics())
	userID := int64(1)

	rows := sqlmock.NewRows([]string{"id", "title", "thumbnail_url", "duration", "album_id", "is_favorite"}).
		AddRow(1, "Most Recent Track", "thumbnail1.jpg", 200, 1, true)

	mock.ExpectQuery("SELECT t.id, t.title, t.thumbnail_url, t.duration, t.album_id, \\(ft.user_id IS NOT NULL\\) AS is_favorite").
		WithArgs(userID).
		WillReturnRows(rows)

	tracks, err := repo.GetMostRecentTracks(ctx, userID)
	assert.NoError(t, err)
	assert.Len(t, tracks, 1)
	assert.Equal(t, int64(1), tracks[0].ID)
	assert.Equal(t, "Most Recent Track", tracks[0].Title)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetMostRecentTracksError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewTrackPostgresRepository(db, metrics.NewMockMetrics())
	userID := int64(1)

	mock.ExpectQuery("SELECT t.id, t.title, t.thumbnail_url, t.duration, t.album_id, \\(ft.user_id IS NOT NULL\\) AS is_favorite").
		WithArgs(userID).
		WillReturnError(stderrors.New("db error"))

	tracks, err := repo.GetMostRecentTracks(ctx, userID)
	assert.Error(t, err)
	assert.Nil(t, tracks)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetMostListenedLastMonthTracks(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewTrackPostgresRepository(db, metrics.NewMockMetrics())
	userID := int64(1)

	rows := sqlmock.NewRows([]string{"id", "title", "thumbnail_url", "duration", "album_id", "is_favorite"}).
		AddRow(1, "Most Listened Last Month Track", "thumbnail1.jpg", 200, 1, true)

	mock.ExpectQuery("SELECT t.id, t.title, t.thumbnail_url, t.duration, t.album_id, \\(ft.user_id IS NOT NULL\\) AS is_favorite").
		WithArgs(userID).
		WillReturnRows(rows)

	tracks, err := repo.GetMostListenedLastMonthTracks(ctx, userID)
	assert.NoError(t, err)
	assert.Len(t, tracks, 1)
	assert.Equal(t, int64(1), tracks[0].ID)
	assert.Equal(t, "Most Listened Last Month Track", tracks[0].Title)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetMostListenedLastMonthTracksError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewTrackPostgresRepository(db, metrics.NewMockMetrics())
	userID := int64(1)

	mock.ExpectQuery("SELECT t.id, t.title, t.thumbnail_url, t.duration, t.album_id, \\(ft.user_id IS NOT NULL\\) AS is_favorite").
		WithArgs(userID).
		WillReturnError(stderrors.New("db error"))

	tracks, err := repo.GetMostListenedLastMonthTracks(ctx, userID)
	assert.Error(t, err)
	assert.Nil(t, tracks)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetMostLikedLastWeekTracks(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewTrackPostgresRepository(db, metrics.NewMockMetrics())
	userID := int64(1)

	rows := sqlmock.NewRows([]string{"id", "title", "thumbnail_url", "duration", "album_id", "is_favorite"}).
		AddRow(1, "Most Liked Last Week Track", "thumbnail1.jpg", 200, 1, true)

	mock.ExpectQuery("SELECT t.id, t.title, t.thumbnail_url, t.duration, t.album_id, \\(ft.user_id IS NOT NULL\\) AS is_favorite").
		WithArgs(userID).
		WillReturnRows(rows)

	tracks, err := repo.GetMostLikedLastWeekTracks(ctx, userID)
	assert.NoError(t, err)
	assert.Len(t, tracks, 1)
	assert.Equal(t, int64(1), tracks[0].ID)
	assert.Equal(t, "Most Liked Last Week Track", tracks[0].Title)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetMostLikedLastWeekTracksError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewTrackPostgresRepository(db, metrics.NewMockMetrics())
	userID := int64(1)

	mock.ExpectQuery("SELECT t.id, t.title, t.thumbnail_url, t.duration, t.album_id, \\(ft.user_id IS NOT NULL\\) AS is_favorite").
		WithArgs(userID).
		WillReturnError(stderrors.New("db error"))

	tracks, err := repo.GetMostLikedLastWeekTracks(ctx, userID)
	assert.Error(t, err)
	assert.Nil(t, tracks)

	assert.NoError(t, mock.ExpectationsWereMet())
}