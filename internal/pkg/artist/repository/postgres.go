package repository

import (
	"database/sql"
	"errors"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/artist"
	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/repository"
)

const (
	GetAllArtistsQuery = `
		SELECT id, title, description, thumbnail_url
		FROM artist
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	GetArtistByIDQuery = `
		SELECT id, title, description, thumbnail_url
		FROM artist
		WHERE id = $1
	`
	GetArtistTitleByIDQuery = `
		SELECT title
		FROM artist
		WHERE id = $1
	`
	GetArtistsByTrackIDQuery = `
		SELECT a.id, a.title, ta.role
		FROM artist a
		JOIN track_artist ta ON ta.artist_id = a.id
		WHERE ta.track_id = $1
		ORDER BY CASE 
			WHEN ta.role = 'main' THEN 1
			WHEN ta.role = 'featured' THEN 2
			WHEN ta.role = 'producer' THEN 3
			ELSE 4
		END ASC
	`
	GetArtistListenersCountQuery = `
		SELECT COUNT(*)
		FROM stream
		LEFT JOIN track ON stream.track_id = track.id
		LEFT JOIN track_artist ON track.id = track_artist.track_id
		WHERE track_artist.artist_id = $1
	`
	GetArtistFavoritesCountQuery = `
		SELECT COUNT(*)
		FROM favorite_artist
		WHERE artist_id = $1
	`
)

type artistPostgresRepository struct {
	db *sql.DB
}

func NewArtistPostgresRepository(db *sql.DB) artist.Repository {
	return &artistPostgresRepository{db: db}
}

func (r *artistPostgresRepository) GetAllArtists(filters *repoModel.ArtistFilters) ([]*repoModel.Artist, error) {
	rows, err := r.db.Query(GetAllArtistsQuery, filters.Pagination.Limit, filters.Pagination.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	artists := make([]*repoModel.Artist, 0)
	for rows.Next() {
		var artist repoModel.Artist
		err = rows.Scan(&artist.ID, &artist.Title, &artist.Description, &artist.Thumbnail)
		if err != nil {
			return nil, err
		}
		artists = append(artists, &artist)
	}

	return artists, nil
}

func (r *artistPostgresRepository) GetArtistByID(id int64) (*repoModel.Artist, error) {
	row := r.db.QueryRow(GetArtistByIDQuery, id)

	var artist repoModel.Artist
	err := row.Scan(&artist.ID, &artist.Title, &artist.Description, &artist.Thumbnail)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repoModel.ErrArtistNotFound
		}
		return nil, err
	}

	return &artist, nil
}

func (r *artistPostgresRepository) GetArtistTitleByID(id int64) (string, error) {
	row := r.db.QueryRow(GetArtistTitleByIDQuery, id)

	var title string
	err := row.Scan(&title)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", repoModel.ErrArtistNotFound
		}
		return "", err
	}

	return title, nil
}

func (r *artistPostgresRepository) GetArtistsByTrackID(id int64) ([]*repoModel.ArtistWithRole, error) {
	rows, err := r.db.Query(GetArtistsByTrackIDQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	artists := make([]*repoModel.ArtistWithRole, 0)
	for rows.Next() {
		var artist repoModel.ArtistWithRole
		err := rows.Scan(&artist.ID, &artist.Title, &artist.Role)
		if err != nil {
			return nil, err
		}
		artists = append(artists, &artist)
	}

	return artists, nil
}

func (r *artistPostgresRepository) GetArtistListenersCount(id int64) (int64, error) {
	row := r.db.QueryRow(GetArtistListenersCountQuery, id)

	var count int64
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *artistPostgresRepository) GetArtistFavoritesCount(id int64) (int64, error) {
	row := r.db.QueryRow(GetArtistFavoritesCountQuery, id)

	var count int64
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
