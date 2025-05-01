package repository

import (
	"context"
	"database/sql"
	"errors"

	loggerPkg "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/track/internal/domain"
	trackErrors "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/track/model/errors"
	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/track/model/repository"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

const (
	GetAllTracksQuery = `
		SELECT id, title, thumbnail_url, duration, album_id
		FROM track
		JOIN track_stats ts ON track.id = ts.track_id
		ORDER BY ts.listeners_count DESC, id DESC
		LIMIT $1 OFFSET $2
	`
	GetTrackByIDQuery = `
		SELECT id, title, thumbnail_url, duration, album_id, file_url
		FROM track
		WHERE id = $1
	`

	CreateStreamQuery = `
		INSERT INTO stream (track_id, user_id) 
		VALUES ($1, $2)
		RETURNING id
	`

	GetStreamByIDQuery = `
 		SELECT id, user_id, track_id, duration
		FROM stream
 		WHERE id = $1
	`

	UpdateStreamDurationQuery = `
		UPDATE stream
		SET duration = $1
		WHERE id = $2
	`

	GetStreamsByUserIDQuery = `
		WITH latest_streams AS (
			SELECT DISTINCT ON (track_id) id, user_id, track_id, duration, created_at
			FROM stream
			WHERE user_id = $1
			ORDER BY track_id, created_at DESC, id DESC
		)
		SELECT id, user_id, track_id, duration
		FROM latest_streams
		ORDER BY created_at DESC, id DESC
		LIMIT $2 OFFSET $3
	`

	GetTracksByIDsQuery = `
		SELECT id, title, thumbnail_url, duration, album_id
		FROM track
		WHERE id = ANY($1)
	`

	GetTracksByIDsFilteredQuery = `
		SELECT id, title, thumbnail_url, duration, album_id
		FROM track
		JOIN track_stats ts ON track.id = ts.track_id
		WHERE id = ANY($1)
		ORDER BY ts.listeners_count DESC, id DESC
		LIMIT $2 OFFSET $3
	`

	GetAlbumIDByTrackIDQuery = `
		SELECT album_id
		FROM track
		WHERE id = $1
	`

	GetTracksByAlbumIDQuery = `
		SELECT id, title, thumbnail_url, duration, album_id
		FROM track
		WHERE album_id = $1
		ORDER BY position ASC
	`

	GetMinutesListenedByUserIDQuery = `
		SELECT COALESCE(SUM(duration) / 60, 0)
		FROM track_stream
		WHERE user_id = $1
	`

	GetTracksListenedByUserIDQuery = `
		SELECT COUNT(DISTINCT track_id)
		FROM track_stream
		WHERE user_id = $1
	`
)

type TrackPostgresRepository struct {
	db *sql.DB
}

func NewTrackPostgresRepository(db *sql.DB) domain.Repository {
	return &TrackPostgresRepository{db: db}
}

func (r *TrackPostgresRepository) GetAllTracks(ctx context.Context, filters *repoModel.TrackFilters) ([]*repoModel.Track, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting all tracks from db", zap.Any("filters", filters), zap.String("query", GetAllTracksQuery))
	rows, err := r.db.Query(GetAllTracksQuery, filters.Pagination.Limit, filters.Pagination.Offset)
	if err != nil {
		logger.Error("failed to get all tracks", zap.Error(err))
		return nil, trackErrors.NewInternalError("failed to get all tracks: %v", err)
	}
	defer rows.Close()

	tracks := make([]*repoModel.Track, 0)
	for rows.Next() {
		var track repoModel.Track
		err := rows.Scan(&track.ID, &track.Title, &track.Thumbnail, &track.Duration, &track.AlbumID)
		if err != nil {
			logger.Error("failed to scan track", zap.Error(err))
			return nil, trackErrors.NewInternalError("failed to scan track: %v", err)
		}
		tracks = append(tracks, &track)
	}

	if err := rows.Err(); err != nil {
		logger.Error("failed to get all tracks", zap.Error(err))
		return nil, trackErrors.NewInternalError("failed to get all tracks: %v", err)
	}

	return tracks, nil
}

func (r *TrackPostgresRepository) GetTrackByID(ctx context.Context, id int64) (*repoModel.TrackWithFileKey, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting track by id from db", zap.Int64("id", id), zap.String("query", GetTrackByIDQuery))
	var trackObject repoModel.TrackWithFileKey
	err := r.db.QueryRow(GetTrackByIDQuery, id).Scan(&trackObject.ID, &trackObject.Title, &trackObject.Thumbnail, &trackObject.Duration, &trackObject.AlbumID, &trackObject.FileKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error("track not found", zap.Error(err))
			return nil, trackErrors.ErrTrackNotFound
		}
		logger.Error("failed to get track by id", zap.Error(err))
		return nil, trackErrors.NewInternalError("failed to get track by id: %v", err)
	}

	return &trackObject, nil
}

func (r *TrackPostgresRepository) CreateStream(ctx context.Context, createData *repoModel.TrackStreamCreateData) (int64, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting to create stream in db", zap.Any("createData", createData), zap.String("query", CreateStreamQuery))
	var streamID int64
	err := r.db.QueryRow(CreateStreamQuery, createData.TrackID, createData.UserID).Scan(&streamID)
	if err != nil {
		logger.Error("failed to create stream", zap.Error(err))
		return 0, trackErrors.NewInternalError("failed to create stream: %v", err)
	}

	return streamID, nil
}

func (r *TrackPostgresRepository) GetStreamByID(ctx context.Context, id int64) (*repoModel.TrackStream, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting stream by id from db", zap.Int64("id", id), zap.String("query", GetStreamByIDQuery))
	var stream repoModel.TrackStream
	err := r.db.QueryRow(GetStreamByIDQuery, id).Scan(&stream.ID, &stream.UserID, &stream.TrackID, &stream.Duration)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error("stream not found", zap.Error(err))
			return nil, trackErrors.ErrStreamNotFound
		}
		logger.Error("failed to get stream by id", zap.Error(err))
		return nil, trackErrors.NewInternalError("failed to get stream by id: %v", err)
	}

	return &stream, nil
}

func (r *TrackPostgresRepository) UpdateStreamDuration(ctx context.Context, endedStream *repoModel.TrackStreamUpdateData) error {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting to update stream duration in db", zap.Any("endedStream", endedStream), zap.String("query", UpdateStreamDurationQuery))
	result, err := r.db.Exec(UpdateStreamDurationQuery, endedStream.Duration, endedStream.StreamID)
	if err != nil {
		logger.Error("failed to update stream duration", zap.Error(err))
		return trackErrors.NewInternalError("failed to update stream duration: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		logger.Error("failed to get rows affected", zap.Error(err))
		return trackErrors.NewInternalError("failed to get rows affected: %v", err)
	}

	if rows == 0 {
		logger.Error("stream not found", zap.Error(trackErrors.ErrFailedToUpdateStreamDuration))
		return trackErrors.ErrFailedToUpdateStreamDuration
	}

	return nil
}

func (r *TrackPostgresRepository) GetStreamsByUserID(ctx context.Context, userID int64, filters *repoModel.TrackFilters) ([]*repoModel.TrackStream, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting streams by user id from db", zap.Int64("userID", userID), zap.String("query", GetStreamsByUserIDQuery))
	rows, err := r.db.Query(GetStreamsByUserIDQuery, userID, filters.Pagination.Limit, filters.Pagination.Offset)
	if err != nil {
		logger.Error("failed to get streams by user id", zap.Error(err))
		return nil, trackErrors.NewInternalError("failed to get streams by user id: %v", err)
	}
	defer rows.Close()

	streams := make([]*repoModel.TrackStream, 0)
	for rows.Next() {
		var stream repoModel.TrackStream
		err := rows.Scan(&stream.ID, &stream.UserID, &stream.TrackID, &stream.Duration)
		if err != nil {
			logger.Error("failed to scan stream", zap.Error(err))
			return nil, trackErrors.NewInternalError("failed to scan stream: %v", err)
		}
		streams = append(streams, &stream)
	}

	if err := rows.Err(); err != nil {
		logger.Error("failed to get streams by user id", zap.Error(err))
		return nil, trackErrors.NewInternalError("failed to get streams by user id: %v", err)
	}

	return streams, nil
}

func (r *TrackPostgresRepository) GetTracksByIDs(ctx context.Context, ids []int64) (map[int64]*repoModel.Track, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting tracks by ids from db", zap.Any("ids", ids), zap.String("query", GetTracksByIDsQuery))
	rows, err := r.db.Query(GetTracksByIDsQuery, pq.Array(ids))
	if err != nil {
		logger.Error("failed to get tracks by ids", zap.Error(err))
		return nil, trackErrors.NewInternalError("failed to get tracks by ids: %v", err)
	}
	defer rows.Close()

	tracks := make(map[int64]*repoModel.Track)
	for rows.Next() {
		var track repoModel.Track
		err := rows.Scan(&track.ID, &track.Title, &track.Thumbnail, &track.Duration, &track.AlbumID)
		if err != nil {
			logger.Error("failed to scan track", zap.Error(err))
			return nil, trackErrors.NewInternalError("failed to scan track: %v", err)
		}
		tracks[track.ID] = &track
	}

	if err := rows.Err(); err != nil {
		logger.Error("failed to get tracks by ids", zap.Error(err))
		return nil, trackErrors.NewInternalError("failed to get tracks by ids: %v", err)
	}

	if len(tracks) > 0 {
		for _, id := range ids {
			if _, ok := tracks[id]; !ok {
				logger.Error("track not found", zap.Int64("id", id))
				return nil, trackErrors.ErrTrackNotFound
			}
		}
	}

	return tracks, nil
}

func (r *TrackPostgresRepository) GetTracksByIDsFiltered(ctx context.Context, ids []int64, filters *repoModel.TrackFilters) ([]*repoModel.Track, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting tracks by ids from db", zap.Any("ids", ids), zap.String("query", GetTracksByIDsFilteredQuery))
	rows, err := r.db.Query(GetTracksByIDsFilteredQuery, pq.Array(ids), filters.Pagination.Limit, filters.Pagination.Offset)
	if err != nil {
		logger.Error("failed to get tracks by ids", zap.Error(err))
		return nil, trackErrors.NewInternalError("failed to get tracks by ids: %v", err)
	}
	defer rows.Close()

	tracks := make([]*repoModel.Track, 0)
	for rows.Next() {
		var track repoModel.Track
		err := rows.Scan(&track.ID, &track.Title, &track.Thumbnail, &track.Duration, &track.AlbumID)
		if err != nil {
			logger.Error("failed to scan track", zap.Error(err))
			return nil, trackErrors.NewInternalError("failed to scan track: %v", err)
		}
		tracks = append(tracks, &track)
	}

	if err := rows.Err(); err != nil {
		logger.Error("failed to get tracks by ids", zap.Error(err))
		return nil, trackErrors.NewInternalError("failed to get tracks by ids: %v", err)
	}

	return tracks, nil
}

func (r *TrackPostgresRepository) GetAlbumIDByTrackID(ctx context.Context, id int64) (int64, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting album id by track id from db", zap.Int64("id", id), zap.String("query", GetAlbumIDByTrackIDQuery))
	var albumID int64
	err := r.db.QueryRow(GetAlbumIDByTrackIDQuery, id).Scan(&albumID)
	if err != nil {
		logger.Error("failed to get album id by track id", zap.Error(err))
		return 0, trackErrors.NewInternalError("failed to get album id by track id: %v", err)
	}

	return albumID, nil
}

func (r *TrackPostgresRepository) GetTracksByAlbumID(ctx context.Context, id int64) ([]*repoModel.Track, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting tracks by album id from db", zap.Int64("id", id), zap.String("query", GetTracksByAlbumIDQuery))
	rows, err := r.db.Query(GetTracksByAlbumIDQuery, id)
	if err != nil {
		logger.Error("failed to get tracks by album id", zap.Error(err))
		return nil, trackErrors.NewInternalError("failed to get tracks by album id: %v", err)
	}
	defer rows.Close()

	tracks := make([]*repoModel.Track, 0)
	for rows.Next() {
		var track repoModel.Track
		err := rows.Scan(&track.ID, &track.Title, &track.Thumbnail, &track.Duration, &track.AlbumID)
		if err != nil {
			logger.Error("failed to scan track", zap.Error(err))
			return nil, trackErrors.NewInternalError("failed to scan track: %v", err)
		}
		tracks = append(tracks, &track)
	}

	return tracks, nil
}

func (r *TrackPostgresRepository) GetMinutesListenedByUserID(ctx context.Context, userID int64) (int64, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting minutes listened by user id from db", zap.Int64("userID", userID), zap.String("query", GetMinutesListenedByUserIDQuery))
	var minutesListened int64
	err := r.db.QueryRow(GetMinutesListenedByUserIDQuery, userID).Scan(&minutesListened)
	if err != nil {
		logger.Error("failed to get minutes listened by user id", zap.Error(err))
		return 0, trackErrors.NewInternalError("failed to get minutes listened by user id: %v", err)
	}

	return minutesListened, nil
}

func (r *TrackPostgresRepository) GetTracksListenedByUserID(ctx context.Context, userID int64) (int64, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting tracks listened by user id from db", zap.Int64("userID", userID), zap.String("query", GetTracksListenedByUserIDQuery))
	var tracksListened int64
	err := r.db.QueryRow(GetTracksListenedByUserIDQuery, userID).Scan(&tracksListened)
	if err != nil {
		logger.Error("failed to get tracks listened by user id", zap.Error(err))
		return 0, trackErrors.NewInternalError("failed to get tracks listened by user id: %v", err)
	}

	return tracksListened, nil
}
