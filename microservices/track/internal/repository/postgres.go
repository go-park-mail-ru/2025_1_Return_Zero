package repository

import (
	"context"
	"database/sql"
	"errors"
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
	rows, err := r.db.Query(GetAllTracksQuery, filters.Pagination.Limit, filters.Pagination.Offset, userID)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetAllTracks").Inc()
		logger.Error("failed to get all tracks", zap.Error(err))
		return nil, trackErrors.NewInternalError("failed to get all tracks: %v", err)
	}
	defer rows.Close()

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
	var trackObject repoModel.TrackWithFileKey
	err := r.db.QueryRowContext(ctx, GetTrackByIDQuery, id, userID).Scan(&trackObject.ID, &trackObject.Title, &trackObject.Thumbnail, &trackObject.Duration, &trackObject.AlbumID, &trackObject.FileKey, &trackObject.IsFavorite)
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
	var streamID int64
	err := r.db.QueryRowContext(ctx, CreateStreamQuery, createData.TrackID, createData.UserID).Scan(&streamID)
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
	var stream repoModel.TrackStream
	err := r.db.QueryRowContext(ctx, GetStreamByIDQuery, id).Scan(&stream.ID, &stream.UserID, &stream.TrackID, &stream.Duration)
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
	result, err := r.db.Exec(UpdateStreamDurationQuery, endedStream.Duration, endedStream.StreamID)
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
	rows, err := r.db.QueryContext(ctx, GetStreamsByUserIDQuery, userID, filters.Pagination.Limit, filters.Pagination.Offset)

	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetStreamsByUserID").Inc()
		logger.Error("failed to get streams by user id", zap.Error(err))
		return nil, trackErrors.NewInternalError("failed to get streams by user id: %v", err)
	}
	defer rows.Close()

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
	rows, err := r.db.QueryContext(ctx, GetTracksByIDsQuery, pq.Array(ids), userID)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetTracksByIDs").Inc()
		logger.Error("failed to get tracks by ids", zap.Error(err))
		return nil, trackErrors.NewInternalError("failed to get tracks by ids: %v", err)
	}
	defer rows.Close()

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
	rows, err := r.db.QueryContext(ctx, GetTracksByIDsFilteredQuery, pq.Array(ids), filters.Pagination.Limit, filters.Pagination.Offset, userID)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetTracksByIDsFiltered").Inc()
		logger.Error("failed to get tracks by ids", zap.Error(err))
		return nil, trackErrors.NewInternalError("failed to get tracks by ids: %v", err)
	}
	defer rows.Close()

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
	var albumID int64
	err := r.db.QueryRowContext(ctx, GetAlbumIDByTrackIDQuery, id).Scan(&albumID)
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
	rows, err := r.db.QueryContext(ctx, GetTracksByAlbumIDQuery, id, userID)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetTracksByAlbumID").Inc()
		logger.Error("failed to get tracks by album id", zap.Error(err))
		return nil, trackErrors.NewInternalError("failed to get tracks by album id: %v", err)
	}
	defer rows.Close()

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
	var minutesListened int64
	err := r.db.QueryRowContext(ctx, GetMinutesListenedByUserIDQuery, userID).Scan(&minutesListened)
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
	var tracksListened int64
	err := r.db.QueryRowContext(ctx, GetTracksListenedByUserIDQuery, userID).Scan(&tracksListened)
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
	var exists bool
	err := r.db.QueryRowContext(ctx, CheckTrackExistsQuery, trackID).Scan(&exists)
	if err != nil {
		logger.Error("failed to check if track exists", zap.Error(err))
		return false, trackErrors.NewInternalError("failed to check if track exists: %v", err)
	}

	return exists, nil
}

func (r *TrackPostgresRepository) LikeTrack(ctx context.Context, likeRequest *repoModel.LikeRequest) error {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting to like track in db", zap.Any("likeRequest", likeRequest), zap.String("query", LikeTrackQuery))
	_, err := r.db.ExecContext(ctx, LikeTrackQuery, likeRequest.TrackID, likeRequest.UserID)
	if err != nil {
		logger.Error("failed to like track", zap.Error(err))
		return trackErrors.NewInternalError("failed to like track: %v", err)
	}

	return nil
}

func (r *TrackPostgresRepository) UnlikeTrack(ctx context.Context, likeRequest *repoModel.LikeRequest) error {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting to unlike track in db", zap.Any("likeRequest", likeRequest), zap.String("query", UnlikeTrackQuery))
	_, err := r.db.ExecContext(ctx, UnlikeTrackQuery, likeRequest.TrackID, likeRequest.UserID)
	if err != nil {
		logger.Error("failed to unlike track", zap.Error(err))
		return trackErrors.NewInternalError("failed to unlike track: %v", err)
	}

	return nil
}

func (r *TrackPostgresRepository) GetFavoriteTracks(ctx context.Context, favoriteRequest *repoModel.FavoriteRequest) ([]*repoModel.Track, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting favorite tracks from db", zap.Any("favoriteRequest", favoriteRequest), zap.String("query", GetFavoriteTracksQuery))
	rows, err := r.db.QueryContext(ctx, GetFavoriteTracksQuery, favoriteRequest.RequestUserID, favoriteRequest.ProfileUserID, favoriteRequest.Filters.Pagination.Limit, favoriteRequest.Filters.Pagination.Offset)
	if err != nil {
		logger.Error("failed to get favorite tracks", zap.Error(err))
		return nil, trackErrors.NewInternalError("failed to get favorite tracks: %v", err)
	}
	defer rows.Close()

	tracks := make([]*repoModel.Track, 0)
	for rows.Next() {
		var track repoModel.Track
		err := rows.Scan(&track.ID, &track.Title, &track.Thumbnail, &track.Duration, &track.AlbumID, &track.IsFavorite)
		if err != nil {
			logger.Error("failed to scan track", zap.Error(err))
			return nil, trackErrors.NewInternalError("failed to scan track: %v", err)
		}
		tracks = append(tracks, &track)
	}

	if err := rows.Err(); err != nil {
		logger.Error("failed to get favorite tracks", zap.Error(err))
		return nil, trackErrors.NewInternalError("failed to get favorite tracks: %v", err)
	}

	return tracks, nil
}

func (r *TrackPostgresRepository) SearchTracks(ctx context.Context, query string, userID int64) ([]*repoModel.Track, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Searching tracks in db", zap.String("search query", query), zap.String("query", SearchTracksQuery))

	words := strings.Fields(query)
	for i, word := range words {
		words[i] = word + ":*"
	}
	tsQueryString := strings.Join(words, " & ")

	rows, err := r.db.QueryContext(ctx, SearchTracksQuery, tsQueryString, userID, query)
	if err != nil {
		logger.Error("failed to search tracks", zap.Error(err))
		return nil, trackErrors.NewInternalError("failed to search tracks: %v", err)
	}
	defer rows.Close()

	var tracks []*repoModel.Track
	for rows.Next() {
		var track repoModel.Track
		err := rows.Scan(&track.ID, &track.Title, &track.Thumbnail, &track.Duration, &track.AlbumID, &track.IsFavorite)
		if err != nil {
			logger.Error("failed to scan track", zap.Error(err))
			return nil, trackErrors.NewInternalError("failed to scan track: %v", err)
		}
		tracks = append(tracks, &track)
	}

	if err := rows.Err(); err != nil {
		logger.Error("failed to search tracks", zap.Error(err))
		return nil, trackErrors.NewInternalError("failed to search tracks: %v", err)
	}

	return tracks, nil
}
