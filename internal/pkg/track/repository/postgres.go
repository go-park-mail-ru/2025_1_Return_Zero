package repository

import (
	"database/sql"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/repository"
)

const (
	GetAllTracksQuery = `
		SELECT id, title, thumbnail_url, duration, album_id
		FROM track
		ORDER BY listeners_count DESC, favorites_count DESC, id DESC
		LIMIT $1 OFFSET $2
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
