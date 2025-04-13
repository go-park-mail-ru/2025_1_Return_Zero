package repository

import (
	"database/sql"
	"errors"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/album"
	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/repository"
)

const (
	GetAllAlbumsQuery = `
		SELECT id, title, type, thumbnail_url, release_date
		FROM album
		ORDER BY release_date DESC, id DESC
		LIMIT $1 OFFSET $2
	`
	GetAlbumByIDQuery = `
		SELECT id, title, type, thumbnail_url, release_date
		FROM album
		WHERE id = $1
	`
	GetAlbumTitleByIDQuery = `
		SELECT title
		FROM album
		WHERE id = $1
	`
	GetAlbumsByArtistIDQuery = `
		SELECT album.id, album.title, album.type, album.thumbnail_url, album.release_date
		FROM album
		JOIN album_artist aa ON album.id = aa.album_id
		WHERE aa.artist_id = $1
		ORDER BY album.release_date DESC, album.id DESC
	`
)

type albumPostgresRepository struct {
	db *sql.DB
}

func NewAlbumPostgresRepository(db *sql.DB) album.Repository {
	return &albumPostgresRepository{
		db: db,
	}
}

func (r *albumPostgresRepository) GetAllAlbums(filters *repoModel.AlbumFilters) ([]*repoModel.Album, error) {
	rows, err := r.db.Query(GetAllAlbumsQuery, filters.Pagination.Limit, filters.Pagination.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	albums := make([]*repoModel.Album, 0)
	for rows.Next() {
		var album repoModel.Album
		err = rows.Scan(&album.ID, &album.Title, &album.Type, &album.Thumbnail, &album.ReleaseDate)
		if err != nil {
			return nil, err
		}
		albums = append(albums, &album)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return albums, nil
}

func (r *albumPostgresRepository) GetAlbumByID(id int64) (*repoModel.Album, error) {
	row := r.db.QueryRow(GetAlbumByIDQuery, id)

	var albumObject repoModel.Album
	err := row.Scan(&albumObject.ID, &albumObject.Title, &albumObject.Type, &albumObject.Thumbnail, &albumObject.ReleaseDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, album.ErrAlbumNotFound
		}
		return nil, err
	}

	return &albumObject, nil
}

func (r *albumPostgresRepository) GetAlbumTitleByID(id int64) (string, error) {
	row := r.db.QueryRow(GetAlbumTitleByIDQuery, id)

	var title string
	err := row.Scan(&title)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", album.ErrAlbumNotFound
		}
		return "", err
	}

	return title, nil
}

func (r *albumPostgresRepository) GetAlbumsByArtistID(artistID int64) ([]*repoModel.Album, error) {
	rows, err := r.db.Query(GetAlbumsByArtistIDQuery, artistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	albums := make([]*repoModel.Album, 0)
	for rows.Next() {
		var album repoModel.Album
		err = rows.Scan(&album.ID, &album.Title, &album.Type, &album.Thumbnail, &album.ReleaseDate)
		if err != nil {
			return nil, err
		}
		albums = append(albums, &album)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return albums, nil
}
