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
	GetArtistStatsQuery = `
		SELECT 
			listeners_count,
			favorites_count
		FROM artist_stats
		WHERE artist_id = $1
	`

	GetArtistsByAlbumIDQuery = `
		SELECT a.id, a.title
		FROM artist a
		JOIN album_artist aa ON a.id = aa.artist_id
		WHERE aa.album_id = $1
		ORDER BY aa.created_at, aa.id
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

	var artistObject repoModel.Artist
	err := row.Scan(&artistObject.ID, &artistObject.Title, &artistObject.Description, &artistObject.Thumbnail)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, artist.ErrArtistNotFound
		}
		return nil, err
	}

	return &artistObject, nil
}

func (r *artistPostgresRepository) GetArtistTitleByID(id int64) (string, error) {
	row := r.db.QueryRow(GetArtistTitleByIDQuery, id)

	var title string
	err := row.Scan(&title)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", artist.ErrArtistNotFound
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

func (r *artistPostgresRepository) GetArtistStats(id int64) (*repoModel.ArtistStats, error) {
	row := r.db.QueryRow(GetArtistStatsQuery, id)

	var stats repoModel.ArtistStats
	err := row.Scan(&stats.ListenersCount, &stats.FavoritesCount)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

func (r *artistPostgresRepository) GetArtistsByAlbumID(albumID int64) ([]*repoModel.ArtistWithTitle, error) {
	rows, err := r.db.Query(GetArtistsByAlbumIDQuery, albumID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	artists := make([]*repoModel.ArtistWithTitle, 0)
	for rows.Next() {
		var artist repoModel.ArtistWithTitle
		err := rows.Scan(&artist.ID, &artist.Title)
		if err != nil {
			return nil, err
		}
		artists = append(artists, &artist)
	}

	return artists, nil
}
