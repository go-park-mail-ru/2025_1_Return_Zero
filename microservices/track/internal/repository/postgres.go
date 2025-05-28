package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	loggerPkg "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/metrics"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/track/internal/domain"
	trackErrors "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/track/model/errors"
	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/track/model/repository"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

const (
	GetAllTracksQuery = `
		SELECT t.id, t.title, t.thumbnail_url, t.duration, t.album_id, (ft.user_id IS NOT NULL) AS is_favorite
		FROM track t
		LEFT JOIN track_stats ts ON t.id = ts.track_id
		LEFT JOIN favorite_track ft ON t.id = ft.track_id AND ft.user_id = $3
		ORDER BY ts.listeners_count DESC, t.id DESC
		LIMIT $1 OFFSET $2
	`
	GetTrackByIDQuery = `
		SELECT t.id, t.title, t.thumbnail_url, t.duration, t.album_id, t.file_url, (ft.user_id IS NOT NULL) AS is_favorite
		FROM track t
		LEFT JOIN favorite_track ft ON t.id = ft.track_id AND ft.user_id = $2
		WHERE t.id = $1
	`

	CreateStreamQuery = `
		INSERT INTO track_stream (track_id, user_id) 
		VALUES ($1, $2)
		RETURNING id
	`

	GetStreamByIDQuery = `
 		SELECT id, user_id, track_id, duration
		FROM track_stream
 		WHERE id = $1
	`

	UpdateStreamDurationQuery = `
		UPDATE track_stream
		SET duration = $1
		WHERE id = $2
	`

	GetStreamsByUserIDQuery = `
		WITH latest_streams AS (
			SELECT DISTINCT ON (track_id) id, user_id, track_id, duration, created_at
			FROM track_stream
			WHERE user_id = $1
			ORDER BY track_id, created_at DESC, id DESC
		)
		SELECT id, user_id, track_id, duration
		FROM latest_streams
		ORDER BY created_at DESC, id DESC
		LIMIT $2 OFFSET $3
	`

	GetTracksByIDsQuery = `
		SELECT t.id, t.title, t.thumbnail_url, t.duration, t.album_id, (ft.user_id IS NOT NULL) AS is_favorite
		FROM track t
		LEFT JOIN favorite_track ft ON t.id = ft.track_id AND ft.user_id = $2
		WHERE t.id = ANY($1)
	`

	GetTracksByIDsFilteredQuery = `
		SELECT t.id, t.title, t.thumbnail_url, t.duration, t.album_id, (ft.user_id IS NOT NULL) AS is_favorite
		FROM track t
		JOIN track_stats ts ON t.id = ts.track_id
		LEFT JOIN favorite_track ft ON t.id = ft.track_id AND ft.user_id = $4
		WHERE t.id = ANY($1)
		ORDER BY ts.listeners_count DESC, t.id DESC
		LIMIT $2 OFFSET $3
	`

	GetAlbumIDByTrackIDQuery = `
		SELECT album_id
		FROM track
		WHERE id = $1
	`

	GetTracksByAlbumIDQuery = `
		SELECT t.id, t.title, t.thumbnail_url, t.duration, t.album_id, (ft.user_id IS NOT NULL) AS is_favorite
		FROM track t
		LEFT JOIN favorite_track ft ON t.id = ft.track_id AND ft.user_id = $2
		WHERE t.album_id = $1
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

	CheckTrackExistsQuery = `
		SELECT EXISTS(SELECT 1 FROM track WHERE id = $1)
	`

	LikeTrackQuery = `
		INSERT INTO favorite_track (track_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING
	`

	UnlikeTrackQuery = `
		DELETE FROM favorite_track WHERE track_id = $1 AND user_id = $2
	`

	GetFavoriteTracksQuery = `
		SELECT t.id, t.title, t.thumbnail_url, t.duration, t.album_id, (ft_req.user_id IS NOT NULL) AS is_favorite
		FROM track t
		JOIN favorite_track ft_prof ON t.id = ft_prof.track_id
		LEFT JOIN favorite_track ft_req ON t.id = ft_req.track_id AND ft_req.user_id = $1
		WHERE ft_prof.user_id = $2
		ORDER BY ft_prof.created_at DESC, t.id DESC
		LIMIT $3 OFFSET $4
	`

	SearchTracksQuery = `
		SELECT t.id, t.title, t.thumbnail_url, t.duration, t.album_id, (ft.user_id IS NOT NULL) AS is_favorite
		FROM track t
		LEFT JOIN favorite_track ft ON t.id = ft.track_id AND ft.user_id = $2
		WHERE t.search_vector @@ to_tsquery('multilingual', $1)
		   OR similarity(t.title_trgm, $3) > 0.3
		ORDER BY 
		    CASE WHEN t.search_vector @@ to_tsquery('multilingual', $1) THEN 0 ELSE 1 END,
		    ts_rank(t.search_vector, to_tsquery('multilingual', $1)) DESC,
		    similarity(t.title_trgm, $3) DESC
	`

	AddTracksToAlbumQuery = `
		INSERT INTO track (title, thumbnail_url, duration, album_id, file_url)
		VALUES ($1, $2, $3, $4, $5)
	`
	DeleteTracksByAlbumIDQuery = `
		DELETE FROM track
		WHERE album_id = $1
`
	GetMostLikedTracksQuery = `
		SELECT t.id, t.title, t.thumbnail_url, t.duration, t.album_id, (ft.user_id IS NOT NULL) AS is_favorite
		FROM track t
		LEFT JOIN track_stats ts ON t.id = ts.track_id
		LEFT JOIN favorite_track ft ON t.id = ft.track_id AND ft.user_id = $1
		ORDER BY ts.favorites_count DESC, t.id DESC
		LIMIT 20
	`

	GetMostRecentTracksQuery = `
		SELECT t.id, t.title, t.thumbnail_url, t.duration, t.album_id, (ft.user_id IS NOT NULL) AS is_favorite
		FROM track t
		LEFT JOIN favorite_track ft ON t.id = ft.track_id AND ft.user_id = $1
		ORDER BY t.created_at DESC, t.id DESC
		LIMIT 20
	`

	GetMostListenedLastMonthTracksQuery = `
		SELECT t.id, t.title, t.thumbnail_url, t.duration, t.album_id, (ft.user_id IS NOT NULL) AS is_favorite
		FROM track t
		LEFT JOIN track_stats ts ON t.id = ts.track_id
		LEFT JOIN favorite_track ft ON t.id = ft.track_id AND ft.user_id = $1
		ORDER BY ts.listeners_count_last_month DESC, t.id DESC
		LIMIT 20
	`

	GetMostLikedLastWeekTracksQuery = `
		SELECT t.id, t.title, t.thumbnail_url, t.duration, t.album_id, (ft.user_id IS NOT NULL) AS is_favorite
		FROM track t
		LEFT JOIN track_stats ts ON t.id = ts.track_id
		LEFT JOIN favorite_track ft ON t.id = ft.track_id AND ft.user_id = $1
		ORDER BY ts.favorites_count_last_week DESC, t.id DESC
		LIMIT 20
	`
)

type TrackPostgresRepository struct {
	db      *sql.DB
	metrics *metrics.Metrics
}

func NewTrackPostgresRepository(db *sql.DB, metrics *metrics.Metrics) domain.Repository {
	return &TrackPostgresRepository{db: db, metrics: metrics}
}

func (r *TrackPostgresRepository) GetAllTracks(ctx context.Context, filters *repoModel.TrackFilters, userID int64) ([]*repoModel.Track, error) {
	start := time.Now()

	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting all tracks from db", zap.Any("filters", filters), zap.String("query", GetAllTracksQuery))

	stmt, err := r.db.PrepareContext(ctx, GetAllTracksQuery)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetAllTracks").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return nil, trackErrors.NewInternalError("failed to prepare statement: %v", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			logger.Error("Error closing statement:", zap.Error(err))
		}
	}()

	rows, err := stmt.QueryContext(ctx, filters.Pagination.Limit, filters.Pagination.Offset, userID)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetAllTracks").Inc()
		logger.Error("failed to get all tracks", zap.Error(err))
		return nil, trackErrors.NewInternalError("failed to get all tracks: %v", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Error("Error closing rows:", zap.Error(err))
		}
	}()

	tracks := make([]*repoModel.Track, 0)
	for rows.Next() {
		var track repoModel.Track
		err := rows.Scan(&track.ID, &track.Title, &track.Thumbnail, &track.Duration, &track.AlbumID, &track.IsFavorite)
		if err != nil {
			r.metrics.DatabaseErrors.WithLabelValues("GetAllTracks").Inc()
			logger.Error("failed to scan track", zap.Error(err))
			return nil, trackErrors.NewInternalError("failed to scan track: %v", err)
		}
		tracks = append(tracks, &track)
	}

	if err := rows.Err(); err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetAllTracks").Inc()
		logger.Error("failed to get all tracks", zap.Error(err))
		return nil, trackErrors.NewInternalError("failed to get all tracks: %v", err)
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("GetAllTracks").Observe(duration)
	return tracks, nil
}

func (r *TrackPostgresRepository) GetTrackByID(ctx context.Context, id int64, userID int64) (*repoModel.TrackWithFileKey, error) {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting track by id from db", zap.Int64("id", id), zap.String("query", GetTrackByIDQuery))

	stmt, err := r.db.PrepareContext(ctx, GetTrackByIDQuery)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetTrackByID").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return nil, trackErrors.NewInternalError("failed to prepare statement: %v", err)
	}

	var trackObject repoModel.TrackWithFileKey
	err = stmt.QueryRowContext(ctx, id, userID).Scan(&trackObject.ID, &trackObject.Title, &trackObject.Thumbnail, &trackObject.Duration, &trackObject.AlbumID, &trackObject.FileKey, &trackObject.IsFavorite)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetTrackByID").Inc()
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error("track not found", zap.Error(err))
			return nil, trackErrors.ErrTrackNotFound
		}
		logger.Error("failed to get track by id", zap.Error(err))
		return nil, trackErrors.NewInternalError("failed to get track by id: %v", err)
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("GetTrackByID").Observe(duration)
	return &trackObject, nil
}

func (r *TrackPostgresRepository) CreateStream(ctx context.Context, createData *repoModel.TrackStreamCreateData) (int64, error) {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting to create stream in db", zap.Any("createData", createData), zap.String("query", CreateStreamQuery))

	stmt, err := r.db.PrepareContext(ctx, CreateStreamQuery)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("CreateStream").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return 0, trackErrors.NewInternalError("failed to prepare statement: %v", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			logger.Error("Error closing statement:", zap.Error(err))
		}
	}()

	var streamID int64
	err = stmt.QueryRowContext(ctx, createData.TrackID, createData.UserID).Scan(&streamID)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("CreateStream").Inc()
		logger.Error("failed to create stream", zap.Error(err))
		return 0, trackErrors.NewInternalError("failed to create stream: %v", err)
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("CreateStream").Observe(duration)
	return streamID, nil
}

func (r *TrackPostgresRepository) GetStreamByID(ctx context.Context, id int64) (*repoModel.TrackStream, error) {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting stream by id from db", zap.Int64("id", id), zap.String("query", GetStreamByIDQuery))

	stmt, err := r.db.PrepareContext(ctx, GetStreamByIDQuery)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetStreamByID").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return nil, trackErrors.NewInternalError("failed to prepare statement: %v", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			logger.Error("Error closing statement:", zap.Error(err))
		}
	}()

	var stream repoModel.TrackStream
	err = stmt.QueryRowContext(ctx, id).Scan(&stream.ID, &stream.UserID, &stream.TrackID, &stream.Duration)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetStreamByID").Inc()
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error("stream not found", zap.Error(err))
			return nil, trackErrors.ErrStreamNotFound
		}
		logger.Error("failed to get stream by id", zap.Error(err))
		return nil, trackErrors.NewInternalError("failed to get stream by id: %v", err)
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("GetStreamByID").Observe(duration)
	return &stream, nil
}

func (r *TrackPostgresRepository) UpdateStreamDuration(ctx context.Context, endedStream *repoModel.TrackStreamUpdateData) error {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting to update stream duration in db", zap.Any("endedStream", endedStream), zap.String("query", UpdateStreamDurationQuery))

	stmt, err := r.db.PrepareContext(ctx, UpdateStreamDurationQuery)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("UpdateStreamDuration").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return trackErrors.NewInternalError("failed to prepare statement: %v", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			logger.Error("Error closing statement:", zap.Error(err))
		}
	}()

	result, err := stmt.ExecContext(ctx, endedStream.Duration, endedStream.StreamID)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("UpdateStreamDuration").Inc()
		logger.Error("failed to update stream duration", zap.Error(err))
		return trackErrors.NewInternalError("failed to update stream duration: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("UpdateStreamDuration").Inc()
		logger.Error("failed to get rows affected", zap.Error(err))
		return trackErrors.NewInternalError("failed to get rows affected: %v", err)
	}

	if rows == 0 {
		r.metrics.DatabaseErrors.WithLabelValues("UpdateStreamDuration").Inc()
		logger.Error("stream not found", zap.Error(trackErrors.ErrFailedToUpdateStreamDuration))
		return trackErrors.ErrFailedToUpdateStreamDuration
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("UpdateStreamDuration").Observe(duration)
	return nil
}

func (r *TrackPostgresRepository) GetStreamsByUserID(ctx context.Context, userID int64, filters *repoModel.TrackFilters) ([]*repoModel.TrackStream, error) {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting streams by user id from db", zap.Int64("userID", userID), zap.String("query", GetStreamsByUserIDQuery))

	stmt, err := r.db.PrepareContext(ctx, GetStreamsByUserIDQuery)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetStreamsByUserID").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return nil, trackErrors.NewInternalError("failed to prepare statement: %v", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			logger.Error("Error closing statement:", zap.Error(err))
		}
	}()

	rows, err := stmt.QueryContext(ctx, userID, filters.Pagination.Limit, filters.Pagination.Offset)

	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetStreamsByUserID").Inc()
		logger.Error("failed to get streams by user id", zap.Error(err))
		return nil, trackErrors.NewInternalError("failed to get streams by user id: %v", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Error("Error closing rows:", zap.Error(err))
		}
	}()

	streams := make([]*repoModel.TrackStream, 0)
	for rows.Next() {
		var stream repoModel.TrackStream
		err := rows.Scan(&stream.ID, &stream.UserID, &stream.TrackID, &stream.Duration)
		if err != nil {
			r.metrics.DatabaseErrors.WithLabelValues("GetStreamsByUserID").Inc()
			logger.Error("failed to scan stream", zap.Error(err))
			return nil, trackErrors.NewInternalError("failed to scan stream: %v", err)
		}
		streams = append(streams, &stream)
	}

	if err := rows.Err(); err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetStreamsByUserID").Inc()
		logger.Error("failed to get streams by user id", zap.Error(err))
		return nil, trackErrors.NewInternalError("failed to get streams by user id: %v", err)
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("GetStreamsByUserID").Observe(duration)
	return streams, nil
}

func (r *TrackPostgresRepository) GetTracksByIDs(ctx context.Context, ids []int64, userID int64) (map[int64]*repoModel.Track, error) {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting tracks by ids from db", zap.Any("ids", ids), zap.String("query", GetTracksByIDsQuery))

	stmt, err := r.db.PrepareContext(ctx, GetTracksByIDsQuery)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetTracksByIDs").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return nil, trackErrors.NewInternalError("failed to prepare statement: %v", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			logger.Error("Error closing statement:", zap.Error(err))
		}
	}()

	rows, err := stmt.QueryContext(ctx, pq.Array(ids), userID)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetTracksByIDs").Inc()
		logger.Error("failed to get tracks by ids", zap.Error(err))
		return nil, trackErrors.NewInternalError("failed to get tracks by ids: %v", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Error("Error closing rows:", zap.Error(err))
		}
	}()

	tracks := make(map[int64]*repoModel.Track)
	for rows.Next() {
		var track repoModel.Track
		err := rows.Scan(&track.ID, &track.Title, &track.Thumbnail, &track.Duration, &track.AlbumID, &track.IsFavorite)
		if err != nil {
			r.metrics.DatabaseErrors.WithLabelValues("GetTracksByIDs").Inc()
			logger.Error("failed to scan track", zap.Error(err))
			return nil, trackErrors.NewInternalError("failed to scan track: %v", err)
		}
		tracks[track.ID] = &track
	}

	if err := rows.Err(); err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetTracksByIDs").Inc()
		logger.Error("failed to get tracks by ids", zap.Error(err))
		return nil, trackErrors.NewInternalError("failed to get tracks by ids: %v", err)
	}

	if len(tracks) > 0 && len(tracks) < len(ids) {
		missingIDs := make([]int64, 0)
		for _, id := range ids {
			if _, ok := tracks[id]; !ok {
				r.metrics.DatabaseErrors.WithLabelValues("GetTracksByIDs").Inc()
				missingIDs = append(missingIDs, id)
				logger.Error("failed to get tracks by ids", zap.Int64("id", id))
			}

		}
		if len(missingIDs) > 0 {
			logger.Warn("some tracks were not found", zap.Int64s("missing_ids", missingIDs))
		}
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("GetTracksByIDs").Observe(duration)
	return tracks, nil
}

func (r *TrackPostgresRepository) GetTracksByIDsFiltered(ctx context.Context, ids []int64, filters *repoModel.TrackFilters, userID int64) ([]*repoModel.Track, error) {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting tracks by ids from db", zap.Any("ids", ids), zap.String("query", GetTracksByIDsFilteredQuery))

	stmt, err := r.db.PrepareContext(ctx, GetTracksByIDsFilteredQuery)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetTracksByIDsFiltered").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return nil, trackErrors.NewInternalError("failed to prepare statement: %v", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			logger.Error("Error closing statement:", zap.Error(err))
		}
	}()

	rows, err := stmt.QueryContext(ctx, pq.Array(ids), filters.Pagination.Limit, filters.Pagination.Offset, userID)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetTracksByIDsFiltered").Inc()
		logger.Error("failed to get tracks by ids", zap.Error(err))
		return nil, trackErrors.NewInternalError("failed to get tracks by ids: %v", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Error("Error closing rows:", zap.Error(err))
		}
	}()

	tracks := make([]*repoModel.Track, 0)
	for rows.Next() {
		var track repoModel.Track
		err := rows.Scan(&track.ID, &track.Title, &track.Thumbnail, &track.Duration, &track.AlbumID, &track.IsFavorite)
		if err != nil {
			r.metrics.DatabaseErrors.WithLabelValues("GetTracksByIDsFiltered").Inc()
			logger.Error("failed to scan track", zap.Error(err))
			return nil, trackErrors.NewInternalError("failed to scan track: %v", err)
		}
		tracks = append(tracks, &track)
	}

	if err := rows.Err(); err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetTracksByIDsFiltered").Inc()
		logger.Error("failed to get tracks by ids", zap.Error(err))
		return nil, trackErrors.NewInternalError("failed to get tracks by ids: %v", err)
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("GetTracksByIDsFiltered").Observe(duration)
	return tracks, nil
}

func (r *TrackPostgresRepository) GetAlbumIDByTrackID(ctx context.Context, id int64) (int64, error) {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting album id by track id from db", zap.Int64("id", id), zap.String("query", GetAlbumIDByTrackIDQuery))

	stmt, err := r.db.PrepareContext(ctx, GetAlbumIDByTrackIDQuery)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetAlbumIDByTrackID").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return 0, trackErrors.NewInternalError("failed to prepare statement: %v", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			logger.Error("Error closing statement:", zap.Error(err))
		}
	}()

	var albumID int64
	err = stmt.QueryRowContext(ctx, id).Scan(&albumID)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetAlbumIDByTrackID").Inc()
		logger.Error("failed to get album id by track id", zap.Error(err))
		return 0, trackErrors.NewInternalError("failed to get album id by track id: %v", err)
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("GetAlbumIDByTrackID").Observe(duration)
	return albumID, nil
}

func (r *TrackPostgresRepository) GetTracksByAlbumID(ctx context.Context, id int64, userID int64) ([]*repoModel.Track, error) {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting tracks by album id from db", zap.Int64("id", id), zap.String("query", GetTracksByAlbumIDQuery))

	stmt, err := r.db.PrepareContext(ctx, GetTracksByAlbumIDQuery)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetTracksByAlbumID").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return nil, trackErrors.NewInternalError("failed to prepare statement: %v", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			logger.Error("Error closing statement:", zap.Error(err))
		}
	}()

	rows, err := stmt.QueryContext(ctx, id, userID)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetTracksByAlbumID").Inc()
		logger.Error("failed to get tracks by album id", zap.Error(err))
		return nil, trackErrors.NewInternalError("failed to get tracks by album id: %v", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Error("Error closing rows:", zap.Error(err))
		}
	}()

	tracks := make([]*repoModel.Track, 0)
	for rows.Next() {
		var track repoModel.Track
		err := rows.Scan(&track.ID, &track.Title, &track.Thumbnail, &track.Duration, &track.AlbumID, &track.IsFavorite)
		if err != nil {
			r.metrics.DatabaseErrors.WithLabelValues("GetTracksByAlbumID").Inc()
			logger.Error("failed to scan track", zap.Error(err))
			return nil, trackErrors.NewInternalError("failed to scan track: %v", err)
		}
		tracks = append(tracks, &track)
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("GetTracksByAlbumID").Observe(duration)
	return tracks, nil
}

func (r *TrackPostgresRepository) GetMinutesListenedByUserID(ctx context.Context, userID int64) (int64, error) {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting minutes listened by user id from db", zap.Int64("userID", userID), zap.String("query", GetMinutesListenedByUserIDQuery))

	stmt, err := r.db.PrepareContext(ctx, GetMinutesListenedByUserIDQuery)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetMinutesListenedByUserID").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return 0, trackErrors.NewInternalError("failed to prepare statement: %v", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			logger.Error("Error closing statement:", zap.Error(err))
		}
	}()

	var minutesListened int64
	err = stmt.QueryRowContext(ctx, userID).Scan(&minutesListened)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetMinutesListenedByUserID").Inc()
		logger.Error("failed to get minutes listened by user id", zap.Error(err))
		return 0, trackErrors.NewInternalError("failed to get minutes listened by user id: %v", err)
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("GetMinutesListenedByUserID").Observe(duration)
	return minutesListened, nil
}

func (r *TrackPostgresRepository) GetTracksListenedByUserID(ctx context.Context, userID int64) (int64, error) {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting tracks listened by user id from db", zap.Int64("userID", userID), zap.String("query", GetTracksListenedByUserIDQuery))

	stmt, err := r.db.PrepareContext(ctx, GetTracksListenedByUserIDQuery)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetTracksListenedByUserID").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return 0, trackErrors.NewInternalError("failed to prepare statement: %v", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			logger.Error("Error closing statement:", zap.Error(err))
		}
	}()

	var tracksListened int64
	err = stmt.QueryRowContext(ctx, userID).Scan(&tracksListened)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetTracksListenedByUserID").Inc()
		logger.Error("failed to get tracks listened by user id", zap.Error(err))
		return 0, trackErrors.NewInternalError("failed to get tracks listened by user id: %v", err)
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("GetTracksListenedByUserID").Observe(duration)
	return tracksListened, nil
}

func (r *TrackPostgresRepository) CheckTrackExists(ctx context.Context, trackID int64) (bool, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting to check if track exists in db", zap.Int64("trackID", trackID), zap.String("query", CheckTrackExistsQuery))

	start := time.Now()

	stmt, err := r.db.PrepareContext(ctx, CheckTrackExistsQuery)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("CheckTrackExists").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return false, trackErrors.NewInternalError("failed to prepare statement: %v", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			logger.Error("Error closing statement:", zap.Error(err))
		}
	}()

	var exists bool
	err = stmt.QueryRowContext(ctx, trackID).Scan(&exists)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("CheckTrackExists").Inc()
		logger.Error("failed to check if track exists", zap.Error(err))
		return false, trackErrors.NewInternalError("failed to check if track exists: %v", err)
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("CheckTrackExists").Observe(duration)
	return exists, nil
}

func (r *TrackPostgresRepository) LikeTrack(ctx context.Context, likeRequest *repoModel.LikeRequest) error {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting to like track in db", zap.Any("likeRequest", likeRequest), zap.String("query", LikeTrackQuery))

	start := time.Now()

	stmt, err := r.db.PrepareContext(ctx, LikeTrackQuery)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("LikeTrack").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return trackErrors.NewInternalError("failed to prepare statement: %v", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			logger.Error("Error closing statement:", zap.Error(err))
		}
	}()

	_, err = stmt.ExecContext(ctx, likeRequest.TrackID, likeRequest.UserID)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("LikeTrack").Inc()
		logger.Error("failed to like track", zap.Error(err))
		return trackErrors.NewInternalError("failed to like track: %v", err)
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("LikeTrack").Observe(duration)
	return nil
}

func (r *TrackPostgresRepository) UnlikeTrack(ctx context.Context, likeRequest *repoModel.LikeRequest) error {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting to unlike track in db", zap.Any("likeRequest", likeRequest), zap.String("query", UnlikeTrackQuery))

	start := time.Now()

	stmt, err := r.db.PrepareContext(ctx, UnlikeTrackQuery)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("UnlikeTrack").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return trackErrors.NewInternalError("failed to prepare statement: %v", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			logger.Error("Error closing statement:", zap.Error(err))
		}
	}()

	_, err = stmt.ExecContext(ctx, likeRequest.TrackID, likeRequest.UserID)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("UnlikeTrack").Inc()
		logger.Error("failed to unlike track", zap.Error(err))
		return trackErrors.NewInternalError("failed to unlike track: %v", err)
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("UnlikeTrack").Observe(duration)
	return nil
}

func (r *TrackPostgresRepository) GetFavoriteTracks(ctx context.Context, favoriteRequest *repoModel.FavoriteRequest) ([]*repoModel.Track, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting favorite tracks from db", zap.Any("favoriteRequest", favoriteRequest), zap.String("query", GetFavoriteTracksQuery))

	start := time.Now()

	stmt, err := r.db.PrepareContext(ctx, GetFavoriteTracksQuery)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetFavoriteTracks").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return nil, trackErrors.NewInternalError("failed to prepare statement: %v", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			logger.Error("Error closing statement:", zap.Error(err))
		}
	}()

	rows, err := stmt.QueryContext(ctx, favoriteRequest.RequestUserID, favoriteRequest.ProfileUserID, favoriteRequest.Filters.Pagination.Limit, favoriteRequest.Filters.Pagination.Offset)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetFavoriteTracks").Inc()
		logger.Error("failed to get favorite tracks", zap.Error(err))
		return nil, trackErrors.NewInternalError("failed to get favorite tracks: %v", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Error("Error closing rows:", zap.Error(err))
		}
	}()

	tracks := make([]*repoModel.Track, 0)
	for rows.Next() {
		var track repoModel.Track
		err := rows.Scan(&track.ID, &track.Title, &track.Thumbnail, &track.Duration, &track.AlbumID, &track.IsFavorite)
		if err != nil {
			r.metrics.DatabaseErrors.WithLabelValues("GetFavoriteTracks").Inc()
			logger.Error("failed to scan track", zap.Error(err))
			return nil, trackErrors.NewInternalError("failed to scan track: %v", err)
		}
		tracks = append(tracks, &track)
	}

	if err := rows.Err(); err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetFavoriteTracks").Inc()
		logger.Error("failed to get favorite tracks", zap.Error(err))
		return nil, trackErrors.NewInternalError("failed to get favorite tracks: %v", err)
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("GetFavoriteTracks").Observe(duration)
	return tracks, nil
}

func (r *TrackPostgresRepository) SearchTracks(ctx context.Context, query string, userID int64) ([]*repoModel.Track, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Searching tracks in db", zap.String("search query", query), zap.String("query", SearchTracksQuery))

	start := time.Now()

	stmt, err := r.db.PrepareContext(ctx, SearchTracksQuery)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("SearchTracks").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return nil, trackErrors.NewInternalError("failed to prepare statement: %v", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			logger.Error("Error closing statement:", zap.Error(err))
		}
	}()

	words := strings.Fields(query)
	for i, word := range words {
		words[i] = word + ":*"
	}
	tsQueryString := strings.Join(words, " & ")

	rows, err := stmt.QueryContext(ctx, tsQueryString, userID, query)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("SearchTracks").Inc()
		logger.Error("failed to search tracks", zap.Error(err))
		return nil, trackErrors.NewInternalError("failed to search tracks: %v", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Error("Error closing rows:", zap.Error(err))
		}
	}()

	var tracks []*repoModel.Track
	for rows.Next() {
		var track repoModel.Track
		err := rows.Scan(&track.ID, &track.Title, &track.Thumbnail, &track.Duration, &track.AlbumID, &track.IsFavorite)
		if err != nil {
			r.metrics.DatabaseErrors.WithLabelValues("SearchTracks").Inc()
			logger.Error("failed to scan track", zap.Error(err))
			return nil, trackErrors.NewInternalError("failed to scan track: %v", err)
		}
		tracks = append(tracks, &track)
	}

	if err := rows.Err(); err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("SearchTracks").Inc()
		logger.Error("failed to search tracks", zap.Error(err))
		return nil, trackErrors.NewInternalError("failed to search tracks: %v", err)
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("SearchTracks").Observe(duration)

	return tracks, nil
}

func (r *TrackPostgresRepository) AddTracksToAlbum(ctx context.Context, tracksList []*repoModel.Track) ([]int64, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Adding tracks to album in db", zap.Any("tracksList", tracksList))

	start := time.Now()

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("AddTracksToAlbum").Inc()
		logger.Error("failed to begin transaction", zap.Error(err))
		return nil, trackErrors.NewInternalError("failed to begin transaction: %v", err)
	}

	valueStrings := make([]string, 0, len(tracksList))
	valueArgs := make([]interface{}, 0, len(tracksList)*6)

	for i, track := range tracksList {
		baseIndex := i * 6
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)",
			baseIndex+1, baseIndex+2, baseIndex+3, baseIndex+4, baseIndex+5, baseIndex+6))

		valueArgs = append(valueArgs,
			track.Title,
			track.Thumbnail,
			track.Duration,
			track.AlbumID,
			fmt.Sprintf("%s.mp3", track.Title),
			i+1)
	}

	stmt := fmt.Sprintf("INSERT INTO track (title, thumbnail_url, duration, album_id, file_url, position) VALUES %s RETURNING id",
		strings.Join(valueStrings, ","))

	rows, err := tx.QueryContext(ctx, stmt, valueArgs...)
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			logger.Error("failed to rollback transaction", zap.Error(err))
		}
		r.metrics.DatabaseErrors.WithLabelValues("AddTracksToAlbum").Inc()
		logger.Error("failed to insert tracks", zap.Error(err))
		return nil, trackErrors.NewInternalError("failed to insert tracks: %v", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Error("Error closing rows:", zap.Error(err))
		}
	}()

	var trackIDs []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			err = tx.Rollback()
			if err != nil {
				logger.Error("failed to rollback transaction", zap.Error(err))
			}
			r.metrics.DatabaseErrors.WithLabelValues("AddTracksToAlbum").Inc()
			logger.Error("failed to scan track id", zap.Error(err))
			return nil, trackErrors.NewInternalError("failed to scan track id: %v", err)
		}
		trackIDs = append(trackIDs, id)
	}

	if err = rows.Err(); err != nil {
		err = tx.Rollback()
		if err != nil {
			logger.Error("failed to rollback transaction", zap.Error(err))
		}
		r.metrics.DatabaseErrors.WithLabelValues("AddTracksToAlbum").Inc()
		logger.Error("failed to iterate over result rows", zap.Error(err))
		return nil, trackErrors.NewInternalError("failed to iterate over result rows: %v", err)
	}

	if err = tx.Commit(); err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("AddTracksToAlbum").Inc()
		logger.Error("failed to commit transaction", zap.Error(err))
		return nil, trackErrors.NewInternalError("failed to commit transaction: %v", err)
	}

	_, err = r.db.ExecContext(ctx, "REFRESH MATERIALIZED VIEW CONCURRENTLY track_stats")
	if err != nil {
		logger.Warn("failed to refresh track_stats view, new tracks may not be visible in statistics for up to 1 hour", zap.Error(err))
	}

	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("AddTracksToAlbum").Observe(duration)

	return trackIDs, nil
}

func (r *TrackPostgresRepository) DeleteTracksByAlbumID(ctx context.Context, albumID int64) error {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Deleting tracks by album id from db", zap.Int64("albumID", albumID))

	start := time.Now()

	stmt, err := r.db.PrepareContext(ctx, DeleteTracksByAlbumIDQuery)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("DeleteTracksByAlbumID").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return trackErrors.NewInternalError("failed to prepare statement: %v", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			logger.Error("Error closing statement:", zap.Error(err))
		}
	}()

	_, err = stmt.ExecContext(ctx, albumID)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("DeleteTracksByAlbumID").Inc()
		logger.Error("failed to delete tracks by album id", zap.Error(err))
		return trackErrors.NewInternalError("failed to delete tracks by album id: %v", err)
	}

	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("DeleteTracksByAlbumID").Observe(duration)

	return nil
}
func (r *TrackPostgresRepository) GetMostLikedTracks(ctx context.Context, userID int64) ([]*repoModel.Track, error) {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting most liked tracks from db", zap.Int64("userID", userID), zap.String("query", GetMostLikedTracksQuery))
	rows, err := r.db.Query(GetMostLikedTracksQuery, userID)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetMostLikedTracks").Inc()
		logger.Error("failed to get most liked tracks", zap.Error(err))
		return nil, trackErrors.NewInternalError("failed to get most liked tracks: %v", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Error("Error closing rows:", zap.Error(err))
		}
	}()

	tracks := make([]*repoModel.Track, 0)
	for rows.Next() {
		var track repoModel.Track
		err := rows.Scan(&track.ID, &track.Title, &track.Thumbnail, &track.Duration, &track.AlbumID, &track.IsFavorite)
		if err != nil {
			r.metrics.DatabaseErrors.WithLabelValues("GetMostLikedTracks").Inc()
			logger.Error("failed to scan track", zap.Error(err))
			return nil, trackErrors.NewInternalError("failed to scan track: %v", err)
		}
		tracks = append(tracks, &track)
	}

	if err := rows.Err(); err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetMostLikedTracks").Inc()
		logger.Error("failed to get most liked tracks", zap.Error(err))
		return nil, trackErrors.NewInternalError("failed to get most liked tracks: %v", err)
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("GetMostLikedTracks").Observe(duration)
	return tracks, nil
}

func (r *TrackPostgresRepository) GetMostRecentTracks(ctx context.Context, userID int64) ([]*repoModel.Track, error) {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting most recent tracks from db", zap.Int64("userID", userID), zap.String("query", GetMostRecentTracksQuery))
	rows, err := r.db.Query(GetMostRecentTracksQuery, userID)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetMostRecentTracks").Inc()
		logger.Error("failed to get most recent tracks", zap.Error(err))
		return nil, trackErrors.NewInternalError("failed to get most recent tracks: %v", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Error("Error closing rows:", zap.Error(err))
		}
	}()

	tracks := make([]*repoModel.Track, 0)
	for rows.Next() {
		var track repoModel.Track
		err := rows.Scan(&track.ID, &track.Title, &track.Thumbnail, &track.Duration, &track.AlbumID, &track.IsFavorite)
		if err != nil {
			r.metrics.DatabaseErrors.WithLabelValues("GetMostRecentTracks").Inc()
			logger.Error("failed to scan track", zap.Error(err))
			return nil, trackErrors.NewInternalError("failed to scan track: %v", err)
		}
		tracks = append(tracks, &track)
	}

	if err := rows.Err(); err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetMostRecentTracks").Inc()
		logger.Error("failed to get most recent tracks", zap.Error(err))
		return nil, trackErrors.NewInternalError("failed to get most recent tracks: %v", err)
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("GetMostRecentTracks").Observe(duration)
	return tracks, nil
}

func (r *TrackPostgresRepository) GetMostListenedLastMonthTracks(ctx context.Context, userID int64) ([]*repoModel.Track, error) {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting most listened last month tracks from db", zap.Int64("userID", userID), zap.String("query", GetMostListenedLastMonthTracksQuery))
	rows, err := r.db.Query(GetMostListenedLastMonthTracksQuery, userID)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetMostListenedLastMonthTracks").Inc()
		logger.Error("failed to get most listened last month tracks", zap.Error(err))
		return nil, trackErrors.NewInternalError("failed to get most listened last month tracks: %v", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Error("Error closing rows:", zap.Error(err))
		}
	}()

	tracks := make([]*repoModel.Track, 0)
	for rows.Next() {
		var track repoModel.Track
		err := rows.Scan(&track.ID, &track.Title, &track.Thumbnail, &track.Duration, &track.AlbumID, &track.IsFavorite)
		if err != nil {
			r.metrics.DatabaseErrors.WithLabelValues("GetMostListenedLastMonthTracks").Inc()
			logger.Error("failed to scan track", zap.Error(err))
			return nil, trackErrors.NewInternalError("failed to scan track: %v", err)
		}
		tracks = append(tracks, &track)
	}

	if err := rows.Err(); err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetMostListenedLastMonthTracks").Inc()
		logger.Error("failed to get most listened last month tracks", zap.Error(err))
		return nil, trackErrors.NewInternalError("failed to get most listened last month tracks: %v", err)
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("GetMostListenedLastMonthTracks").Observe(duration)
	return tracks, nil
}

func (r *TrackPostgresRepository) GetMostLikedLastWeekTracks(ctx context.Context, userID int64) ([]*repoModel.Track, error) {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting most liked last week tracks from db", zap.Int64("userID", userID), zap.String("query", GetMostLikedLastWeekTracksQuery))
	rows, err := r.db.Query(GetMostLikedLastWeekTracksQuery, userID)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetMostLikedLastWeekTracks").Inc()
		logger.Error("failed to get most liked last week tracks", zap.Error(err))
		return nil, trackErrors.NewInternalError("failed to get most liked last week tracks: %v", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Error("Error closing rows:", zap.Error(err))
		}
	}()

	tracks := make([]*repoModel.Track, 0)
	for rows.Next() {
		var track repoModel.Track
		err := rows.Scan(&track.ID, &track.Title, &track.Thumbnail, &track.Duration, &track.AlbumID, &track.IsFavorite)
		if err != nil {
			r.metrics.DatabaseErrors.WithLabelValues("GetMostLikedLastWeekTracks").Inc()
			logger.Error("failed to scan track", zap.Error(err))
			return nil, trackErrors.NewInternalError("failed to scan track: %v", err)
		}
		tracks = append(tracks, &track)
	}

	if err := rows.Err(); err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetMostLikedLastWeekTracks").Inc()
		logger.Error("failed to get most liked last week tracks", zap.Error(err))
		return nil, trackErrors.NewInternalError("failed to get most liked last week tracks: %v", err)
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("GetMostLikedLastWeekTracks").Observe(duration)
	return tracks, nil
}
