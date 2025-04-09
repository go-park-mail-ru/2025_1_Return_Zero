package repository

import (
	"database/sql"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/track"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/repository"
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

func (r *TrackPostgresRepository) GetAllTracks(filters *repository.TrackFilters) ([]*repository.Track, error) {
	rows, err := r.db.Query(GetAllTracksQuery, filters.Pagination.Limit, filters.Pagination.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tracks := make([]*repository.Track, 0)
	for rows.Next() {
		var track repository.Track
		err := rows.Scan(&track.ID, &track.Title, &track.Thumbnail, &track.Duration, &track.AlbumID)
		if err != nil {
			return nil, err
		}
		tracks = append(tracks, &track)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tracks, nil
}

func (r *TrackPostgresRepository) GetTrackByID(id int64) (*repository.TrackWithFileKey, error) {
	var track repository.TrackWithFileKey
	err := r.db.QueryRow(GetTrackByIDQuery, id).Scan(&track.ID, &track.Title, &track.Thumbnail, &track.Duration, &track.AlbumID, &track.FileKey)
	if err != nil {
		return nil, err
	}

	return &track, nil
}

func (r *TrackPostgresRepository) GetTracksByArtistID(artistID int64) ([]*repository.Track, error) {
	rows, err := r.db.Query(GetTracksByArtistIDQuery, artistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tracks := make([]*repository.Track, 0)
	for rows.Next() {
		var track repository.Track
		err := rows.Scan(&track.ID, &track.Title, &track.Thumbnail, &track.Duration, &track.AlbumID)
		if err != nil {
			return nil, err
		}
		tracks = append(tracks, &track)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tracks, nil
}

func (r *TrackPostgresRepository) CreateStream(createData *repository.TrackStreamCreateData) (int64, error) {
	var streamID int64
	err := r.db.QueryRow(CreateStreamQuery, createData.TrackID, createData.UserID).Scan(&streamID)
	if err != nil {
		return 0, err
	}

	return streamID, nil
}

func (r *TrackPostgresRepository) GetStreamByID(id int64) (*repository.TrackStream, error) {
	var stream repository.TrackStream
	err := r.db.QueryRow(GetStreamByIDQuery, id).Scan(&stream.ID, &stream.UserID, &stream.TrackID, &stream.Duration)
	if err != nil {
		return nil, err
	}

	return &stream, nil
}

func (r *TrackPostgresRepository) UpdateStreamDuration(endedStream *repository.TrackStreamUpdateData) error {
	result, err := r.db.Exec(UpdateStreamDurationQuery, endedStream.Duration, endedStream.StreamID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return track.ErrFailedToUpdateStreamDuration
	}

	return nil
}
