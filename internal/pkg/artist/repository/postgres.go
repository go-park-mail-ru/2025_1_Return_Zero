package repository

import (
	"database/sql"
	"errors"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/artist"
	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/repository"
)

const (
	GetAllArtistsQuery = `
		SELECT id, title, description, thumbnail_url, listeners_count, favorites_count
		FROM artist
		ORDER BY listeners_count DESC, favorites_count DESC, id DESC
		LIMIT $1 OFFSET $2
	`
	GetArtistByIDQuery = `
		SELECT id, title, description, thumbnail_url, listeners_count, favorites_count
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
		ORDER BY a.listeners_count DESC, a.favorites_count DESC, a.id DESC
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
		err = rows.Scan(&artist.ID, &artist.Title, &artist.Description, &artist.Thumbnail, &artist.Listeners, &artist.Favorites)
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
	err := row.Scan(&artist.ID, &artist.Title, &artist.Description, &artist.Thumbnail, &artist.Listeners, &artist.Favorites)

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
