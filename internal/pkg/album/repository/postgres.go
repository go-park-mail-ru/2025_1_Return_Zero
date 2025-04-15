package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/album"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers"
	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/repository"
	"go.uber.org/zap"
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
	GetAlbumTitleByIDsQuery = `
		SELECT id, title
		FROM album
		WHERE id = ANY($1)
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

func (r *albumPostgresRepository) GetAllAlbums(ctx context.Context, filters *repoModel.AlbumFilters) ([]*repoModel.Album, error) {
	logger := helpers.LoggerFromContext(ctx)
	logger.Info("Requesting all albums from db", zap.Any("filters", filters), zap.String("query", GetAllAlbumsQuery))
	rows, err := r.db.Query(GetAllAlbumsQuery, filters.Pagination.Limit, filters.Pagination.Offset)
	if err != nil {
		logger.Error("failed to get all albums", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	albums := make([]*repoModel.Album, 0)
	for rows.Next() {
		var album repoModel.Album
		err = rows.Scan(&album.ID, &album.Title, &album.Type, &album.Thumbnail, &album.ReleaseDate)
		if err != nil {
			logger.Error("failed to scan album", zap.Error(err))
			return nil, err
		}
		albums = append(albums, &album)
	}

	if err := rows.Err(); err != nil {
		logger.Error("failed to get all albums", zap.Error(err))
		return nil, err
	}

	return albums, nil
}

func (r *albumPostgresRepository) GetAlbumByID(ctx context.Context, id int64) (*repoModel.Album, error) {
	logger := helpers.LoggerFromContext(ctx)
	logger.Info("Requesting album by id from db", zap.Int64("id", id), zap.String("query", GetAlbumByIDQuery))
	row := r.db.QueryRow(GetAlbumByIDQuery, id)

	var albumObject repoModel.Album
	err := row.Scan(&albumObject.ID, &albumObject.Title, &albumObject.Type, &albumObject.Thumbnail, &albumObject.ReleaseDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error("album not found", zap.Error(err))
			return nil, album.ErrAlbumNotFound
		}
		logger.Error("failed to get album by id", zap.Error(err))
		return nil, err
	}

	return &albumObject, nil
}

func (r *albumPostgresRepository) GetAlbumTitleByIDs(ctx context.Context, ids []int64) (map[int64]string, error) {
	logger := helpers.LoggerFromContext(ctx)
	logger.Info("Requesting album title by ids from db", zap.Any("ids", ids), zap.String("query", GetAlbumTitleByIDsQuery))
	rows, err := r.db.Query(GetAlbumTitleByIDsQuery, ids)
	if err != nil {
		logger.Error("failed to get album title by ids", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	albums := make(map[int64]string)
	for rows.Next() {
		var id int64
		var title string
		err = rows.Scan(&id, &title)
		if err != nil {
			logger.Error("failed to scan album title", zap.Error(err))
			return nil, err
		}
		albums[id] = title
	}

	if err := rows.Err(); err != nil {
		logger.Error("failed to get album title by ids", zap.Error(err))
		return nil, err
	}

	return albums, nil
}

func (r *albumPostgresRepository) GetAlbumTitleByID(ctx context.Context, id int64) (string, error) {
	logger := helpers.LoggerFromContext(ctx)
	logger.Info("Requesting album title by id from db", zap.Int64("id", id), zap.String("query", GetAlbumTitleByIDQuery))
	row := r.db.QueryRow(GetAlbumTitleByIDQuery, id)

	var title string
	err := row.Scan(&title)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error("album not found", zap.Error(err))
			return "", album.ErrAlbumNotFound
		}
		logger.Error("failed to get album title by id", zap.Error(err))
		return "", err
	}

	return title, nil
}

func (r *albumPostgresRepository) GetAlbumsByArtistID(ctx context.Context, artistID int64) ([]*repoModel.Album, error) {
	logger := helpers.LoggerFromContext(ctx)
	logger.Info("Requesting albums by artist id from db", zap.Int64("artistID", artistID), zap.String("query", GetAlbumsByArtistIDQuery))
	rows, err := r.db.Query(GetAlbumsByArtistIDQuery, artistID)
	if err != nil {
		logger.Error("failed to get albums by artist id", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	albums := make([]*repoModel.Album, 0)
	for rows.Next() {
		var album repoModel.Album
		err = rows.Scan(&album.ID, &album.Title, &album.Type, &album.Thumbnail, &album.ReleaseDate)
		if err != nil {
			logger.Error("failed to scan album", zap.Error(err))
			return nil, err
		}
		albums = append(albums, &album)
	}

	if err := rows.Err(); err != nil {
		logger.Error("failed to get albums by artist id", zap.Error(err))
		return nil, err
	}

	return albums, nil
}
