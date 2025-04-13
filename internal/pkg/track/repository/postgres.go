package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/middleware"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/track"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/repository"
	"go.uber.org/zap"
)

const (
	GetAllTracksQuery = `
		SELECT id, title, thumbnail_url, duration, album_id
		FROM track
		ORDER BY created_at DESC, id DESC
		LIMIT $1 OFFSET $2
	`
	GetTrackByIDQuery = `
		SELECT id, title, thumbnail_url, duration, album_id, file_url
		FROM track
		WHERE id = $1
	`
	GetTracksByArtistIDQuery = `
		SELECT track.id, track.title, track.thumbnail_url, track.duration, track.album_id
		FROM track
		JOIN track_artist ta ON track.id = ta.track_id
		WHERE ta.artist_id = $1 AND (ta.role = 'main' OR ta.role = 'featured')
		ORDER BY track.created_at DESC, track.id DESC
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
)

type TrackPostgresRepository struct {
	db *sql.DB
}

func NewTrackPostgresRepository(db *sql.DB) *TrackPostgresRepository {
	return &TrackPostgresRepository{db: db}
}

func (r *TrackPostgresRepository) GetAllTracks(ctx context.Context, filters *repository.TrackFilters) ([]*repository.Track, error) {
	logger := middleware.LoggerFromContext(ctx)
	logger.Info("Requesting all tracks from db", zap.Any("filters", filters))
	rows, err := r.db.Query(GetAllTracksQuery, filters.Pagination.Limit, filters.Pagination.Offset)
	if err != nil {
		logger.Error("failed to get all tracks", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	tracks := make([]*repository.Track, 0)
	for rows.Next() {
		var track repository.Track
		err := rows.Scan(&track.ID, &track.Title, &track.Thumbnail, &track.Duration, &track.AlbumID)
		if err != nil {
			logger.Error("failed to scan track", zap.Error(err))
			return nil, err
		}
		tracks = append(tracks, &track)
	}

	if err := rows.Err(); err != nil {
		logger.Error("failed to get all tracks", zap.Error(err))
		return nil, err
	}

	return tracks, nil
}

func (r *TrackPostgresRepository) GetTrackByID(ctx context.Context, id int64) (*repository.TrackWithFileKey, error) {
	logger := middleware.LoggerFromContext(ctx)
	logger.Info("Requesting track by id from db", zap.Int64("id", id))
	var trackObject repository.TrackWithFileKey
	err := r.db.QueryRow(GetTrackByIDQuery, id).Scan(&trackObject.ID, &trackObject.Title, &trackObject.Thumbnail, &trackObject.Duration, &trackObject.AlbumID, &trackObject.FileKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error("track not found", zap.Error(err))
			return nil, track.ErrTrackNotFound
		}
		logger.Error("failed to get track by id", zap.Error(err))
		return nil, err
	}

	return &trackObject, nil
}

func (r *TrackPostgresRepository) GetTracksByArtistID(ctx context.Context, artistID int64) ([]*repository.Track, error) {
	logger := middleware.LoggerFromContext(ctx)
	logger.Info("Requesting tracks by artist id from db", zap.Int64("artistID", artistID))
	rows, err := r.db.Query(GetTracksByArtistIDQuery, artistID)
	if err != nil {
		logger.Error("failed to get tracks by artist id", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	tracks := make([]*repository.Track, 0)
	for rows.Next() {
		var track repository.Track
		err := rows.Scan(&track.ID, &track.Title, &track.Thumbnail, &track.Duration, &track.AlbumID)
		if err != nil {
			logger.Error("failed to scan track", zap.Error(err))
			return nil, err
		}
		tracks = append(tracks, &track)
	}

	if err := rows.Err(); err != nil {
		logger.Error("failed to get tracks by artist id", zap.Error(err))
		return nil, err
	}

	return tracks, nil
}

func (r *TrackPostgresRepository) CreateStream(ctx context.Context, createData *repository.TrackStreamCreateData) (int64, error) {
	logger := middleware.LoggerFromContext(ctx)
	logger.Info("Requesting to create stream in db", zap.Any("createData", createData))
	var streamID int64
	err := r.db.QueryRow(CreateStreamQuery, createData.TrackID, createData.UserID).Scan(&streamID)
	if err != nil {
		logger.Error("failed to create stream", zap.Error(err))
		return 0, err
	}

	return streamID, nil
}

func (r *TrackPostgresRepository) GetStreamByID(ctx context.Context, id int64) (*repository.TrackStream, error) {
	logger := middleware.LoggerFromContext(ctx)
	logger.Info("Requesting stream by id from db", zap.Int64("id", id))
	var stream repository.TrackStream
	err := r.db.QueryRow(GetStreamByIDQuery, id).Scan(&stream.ID, &stream.UserID, &stream.TrackID, &stream.Duration)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error("stream not found", zap.Error(err))
			return nil, track.ErrStreamNotFound
		}
		logger.Error("failed to get stream by id", zap.Error(err))
		return nil, err
	}

	return &stream, nil
}

func (r *TrackPostgresRepository) UpdateStreamDuration(ctx context.Context, endedStream *repository.TrackStreamUpdateData) error {
	logger := middleware.LoggerFromContext(ctx)
	logger.Info("Requesting to update stream duration in db", zap.Any("endedStream", endedStream))
	result, err := r.db.Exec(UpdateStreamDurationQuery, endedStream.Duration, endedStream.StreamID)
	if err != nil {
		logger.Error("failed to update stream duration", zap.Error(err))
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		logger.Error("failed to get rows affected", zap.Error(err))
		return err
	}

	if rows == 0 {
		logger.Error("stream not found", zap.Error(track.ErrFailedToUpdateStreamDuration))
		return track.ErrFailedToUpdateStreamDuration
	}

	return nil
}
