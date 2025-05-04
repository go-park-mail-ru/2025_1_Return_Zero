package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	loggerPkg "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
	domain "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/album/internal/domain"
	albumErrors "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/album/model/errors"
	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/album/model/repository"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

const (
	GetAllAlbumsQuery = `
		SELECT a.id, a.title, a.type, a.thumbnail_url, a.release_date, (fa.user_id IS NOT NULL) AS is_favorite
		FROM album a
		LEFT JOIN album_stats als ON a.id = als.album_id
		LEFT JOIN favorite_album fa ON a.id = fa.album_id AND fa.user_id = $3
		ORDER BY als.listeners_count DESC, a.id DESC
		LIMIT $1 OFFSET $2
	`
	GetAlbumByIDQuery = `
		SELECT a.id, a.title, a.type, a.thumbnail_url, a.release_date, (fa.user_id IS NOT NULL) AS is_favorite
		FROM album a
		LEFT JOIN favorite_album fa ON a.id = fa.album_id AND fa.user_id = $2
		WHERE a.id = $1
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
	GetAlbumsByIDsQuery = `
		SELECT a.id, a.title, a.type, a.thumbnail_url, a.release_date, (fa.user_id IS NOT NULL) AS is_favorite
		FROM album a
		LEFT JOIN album_stats als ON a.id = als.album_id
		LEFT JOIN favorite_album fa ON a.id = fa.album_id AND fa.user_id = $2
		WHERE a.id = ANY($1)
		ORDER BY als.listeners_count DESC, a.id DESC
	`

	CreateStreamQuery = `
		INSERT INTO album_stream (album_id, user_id)
		VALUES ($1, $2)
	`

	CheckAlbumExistsQuery = `
		SELECT EXISTS (
			SELECT 1
			FROM album
			WHERE id = $1
		)
	`

	LikeAlbumQuery = `
		INSERT INTO favorite_album (album_id, user_id)
		VALUES ($1, $2) ON CONFLICT DO NOTHING
	`

	UnlikeAlbumQuery = `
		DELETE FROM favorite_album
		WHERE album_id = $1 AND user_id = $2
	`

	GetFavoriteAlbumsQuery = `
		SELECT a.id, a.title, a.type, a.thumbnail_url, a.release_date
		FROM album a
		JOIN favorite_album fa ON a.id = fa.album_id
		WHERE fa.user_id = $1
		ORDER BY fa.created_at DESC, a.id DESC
		LIMIT $2 OFFSET $3
	`

	SearchAlbumsQuery = `
		SELECT a.id, a.title, a.type, a.thumbnail_url, a.release_date, (fa.user_id IS NOT NULL) AS is_favorite
		FROM album a
		LEFT JOIN favorite_album fa ON a.id = fa.album_id AND fa.user_id = $2
		WHERE a.search_vector @@ to_tsquery('multilingual', $1)
		   OR similarity(a.title_trgm, $3) > 0.3
		ORDER BY 
		    CASE WHEN a.search_vector @@ to_tsquery('multilingual', $1) THEN 0 ELSE 1 END,
		    ts_rank(a.search_vector, to_tsquery('multilingual', $1)) DESC,
		    similarity(a.title_trgm, $3) DESC
	`
)

type albumPostgresRepository struct {
	db *sql.DB
}

func NewAlbumPostgresRepository(db *sql.DB) domain.Repository {
	return &albumPostgresRepository{
		db: db,
	}
}

func (r *albumPostgresRepository) GetAllAlbums(ctx context.Context, filters *repoModel.AlbumFilters, userID int64) ([]*repoModel.Album, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting all albums from db", zap.Any("filters", filters), zap.String("query", GetAllAlbumsQuery))
	rows, err := r.db.Query(GetAllAlbumsQuery, filters.Pagination.Limit, filters.Pagination.Offset, userID)
	if err != nil {
		logger.Error("failed to get all albums", zap.Error(err))
		return nil, albumErrors.NewInternalError("failed to get all albums: %v", err)
	}
	defer rows.Close()

	albums := make([]*repoModel.Album, 0)
	for rows.Next() {
		var album repoModel.Album
		err = rows.Scan(&album.ID, &album.Title, &album.Type, &album.Thumbnail, &album.ReleaseDate, &album.IsFavorite)
		if err != nil {
			logger.Error("failed to scan album", zap.Error(err))
			return nil, albumErrors.NewInternalError("failed to scan album: %v", err)
		}
		albums = append(albums, &album)
	}

	if err := rows.Err(); err != nil {
		logger.Error("failed to get all albums", zap.Error(err))
		return nil, albumErrors.NewInternalError("failed to get all albums: %v", err)
	}

	return albums, nil
}

func (r *albumPostgresRepository) GetAlbumByID(ctx context.Context, id int64, userID int64) (*repoModel.Album, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting album by id from db", zap.Int64("id", id), zap.String("query", GetAlbumByIDQuery))
	row := r.db.QueryRow(GetAlbumByIDQuery, id, userID)

	var albumObject repoModel.Album
	err := row.Scan(&albumObject.ID, &albumObject.Title, &albumObject.Type, &albumObject.Thumbnail, &albumObject.ReleaseDate, &albumObject.IsFavorite)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error("album not found", zap.Error(err))
			return nil, albumErrors.ErrAlbumNotFound
		}
		logger.Error("failed to get album by id", zap.Error(err))
		return nil, albumErrors.NewInternalError("failed to get album by id: %v", err)
	}

	return &albumObject, nil
}

func (r *albumPostgresRepository) GetAlbumTitleByIDs(ctx context.Context, ids []int64) (map[int64]string, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting album title by ids from db", zap.Any("ids", ids), zap.String("query", GetAlbumTitleByIDsQuery))
	rows, err := r.db.Query(GetAlbumTitleByIDsQuery, pq.Array(ids))
	if err != nil {
		logger.Error("failed to get album title by ids", zap.Error(err))
		return nil, albumErrors.NewInternalError("failed to get album title by ids: %v", err)
	}
	defer rows.Close()

	albums := make(map[int64]string)
	for rows.Next() {
		var id int64
		var title string
		err = rows.Scan(&id, &title)
		if err != nil {
			logger.Error("failed to scan album title", zap.Error(err))
			return nil, albumErrors.NewInternalError("failed to scan album title: %v", err)
		}
		albums[id] = title
	}

	if err := rows.Err(); err != nil {
		logger.Error("failed to get album title by ids", zap.Error(err))
		return nil, albumErrors.NewInternalError("failed to get album title by ids: %v", err)
	}

	return albums, nil
}

func (r *albumPostgresRepository) GetAlbumTitleByID(ctx context.Context, id int64) (string, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting album title by id from db", zap.Int64("id", id), zap.String("query", GetAlbumTitleByIDQuery))
	row := r.db.QueryRow(GetAlbumTitleByIDQuery, id)

	var title string
	err := row.Scan(&title)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error("album not found", zap.Error(err))
			return "", albumErrors.ErrAlbumNotFound
		}
		logger.Error("failed to get album title by id", zap.Error(err))
		return "", albumErrors.NewInternalError("failed to get album title by id: %v", err)
	}

	return title, nil
}

func (r *albumPostgresRepository) GetAlbumsByIDs(ctx context.Context, ids []int64, userID int64) ([]*repoModel.Album, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting albums by ids from db", zap.Any("ids", ids), zap.String("query", GetAlbumsByIDsQuery))
	rows, err := r.db.Query(GetAlbumsByIDsQuery, pq.Array(ids), userID)
	if err != nil {
		logger.Error("failed to get albums by ids", zap.Error(err))
		return nil, albumErrors.NewInternalError("failed to get albums by ids: %v", err)
	}
	defer rows.Close()

	var albums []*repoModel.Album
	for rows.Next() {
		var album repoModel.Album
		err = rows.Scan(&album.ID, &album.Title, &album.Type, &album.Thumbnail, &album.ReleaseDate, &album.IsFavorite)
		if err != nil {
			logger.Error("failed to scan album", zap.Error(err))
			return nil, albumErrors.NewInternalError("failed to scan album: %v", err)
		}
		albums = append(albums, &album)
	}

	if err := rows.Err(); err != nil {
		logger.Error("failed to get albums by ids", zap.Error(err))
		return nil, albumErrors.NewInternalError("failed to get albums by ids: %v", err)
	}

	return albums, nil
}

func (r *albumPostgresRepository) CreateStream(ctx context.Context, albumID int64, userID int64) error {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Creating stream for album", zap.Int64("albumID", albumID), zap.Int64("userID", userID), zap.String("query", CreateStreamQuery))
	_, err := r.db.Exec(CreateStreamQuery, albumID, userID)
	if err != nil {
		logger.Error("failed to create stream", zap.Error(err))
		return albumErrors.NewInternalError("failed to create stream: %v", err)
	}
	return nil
}

func (r *albumPostgresRepository) CheckAlbumExists(ctx context.Context, albumID int64) (bool, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Checking if album exists in db", zap.Int64("albumID", albumID), zap.String("query", CheckAlbumExistsQuery))
	row := r.db.QueryRow(CheckAlbumExistsQuery, albumID)

	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		logger.Error("failed to check if album exists", zap.Error(err))
		return false, albumErrors.NewInternalError("failed to check if album exists: %v", err)
	}
	return exists, nil
}

func (r *albumPostgresRepository) LikeAlbum(ctx context.Context, request *repoModel.LikeRequest) error {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Liking album", zap.Int64("albumID", request.AlbumID), zap.Int64("userID", request.UserID), zap.String("query", LikeAlbumQuery))
	_, err := r.db.Exec(LikeAlbumQuery, request.AlbumID, request.UserID)
	if err != nil {
		logger.Error("failed to like album", zap.Error(err))
		return albumErrors.NewInternalError("failed to like album: %v", err)
	}
	return nil
}

func (r *albumPostgresRepository) UnlikeAlbum(ctx context.Context, request *repoModel.LikeRequest) error {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Unliking album", zap.Int64("albumID", request.AlbumID), zap.Int64("userID", request.UserID), zap.String("query", UnlikeAlbumQuery))
	_, err := r.db.Exec(UnlikeAlbumQuery, request.AlbumID, request.UserID)
	if err != nil {
		logger.Error("failed to unlike album", zap.Error(err))
		return albumErrors.NewInternalError("failed to unlike album: %v", err)
	}
	return nil
}

func (r *albumPostgresRepository) GetFavoriteAlbums(ctx context.Context, filters *repoModel.AlbumFilters, userID int64) ([]*repoModel.Album, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting favorite albums from db", zap.Any("filters", filters), zap.String("query", GetFavoriteAlbumsQuery))
	rows, err := r.db.Query(GetFavoriteAlbumsQuery, userID, filters.Pagination.Limit, filters.Pagination.Offset)
	if err != nil {
		logger.Error("failed to get favorite albums", zap.Error(err))
		return nil, albumErrors.NewInternalError("failed to get favorite albums: %v", err)
	}
	defer rows.Close()

	var albums []*repoModel.Album
	for rows.Next() {
		var album repoModel.Album
		// Ставим по дефолту, так как запрашивашиваются избранные, то есть заведомо известно, что они лайкнуты
		album.IsFavorite = true
		err = rows.Scan(&album.ID, &album.Title, &album.Type, &album.Thumbnail, &album.ReleaseDate)
		if err != nil {
			logger.Error("failed to scan album", zap.Error(err))
			return nil, albumErrors.NewInternalError("failed to scan album: %v", err)
		}
		albums = append(albums, &album)
	}

	if err := rows.Err(); err != nil {
		logger.Error("failed to get favorite albums", zap.Error(err))
		return nil, albumErrors.NewInternalError("failed to get favorite albums: %v", err)
	}

	return albums, nil
}

func (r *albumPostgresRepository) SearchAlbums(ctx context.Context, query string, userID int64) ([]*repoModel.Album, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Searching albums by query", zap.String("query", query), zap.String("query", SearchAlbumsQuery))

	words := strings.Fields(query)
	for i, word := range words {
		words[i] = word + ":*"
	}
	tsQueryString := strings.Join(words, " & ")

	rows, err := r.db.Query(SearchAlbumsQuery, tsQueryString, userID, query)
	if err != nil {
		logger.Error("failed to search albums", zap.Error(err))
		return nil, albumErrors.NewInternalError("failed to search albums: %v", err)
	}
	defer rows.Close()

	var albums []*repoModel.Album
	for rows.Next() {
		var album repoModel.Album
		err = rows.Scan(&album.ID, &album.Title, &album.Type, &album.Thumbnail, &album.ReleaseDate, &album.IsFavorite)
		if err != nil {
			logger.Error("failed to scan album", zap.Error(err))
			return nil, albumErrors.NewInternalError("failed to scan album: %v", err)
		}
		albums = append(albums, &album)
	}

	if err := rows.Err(); err != nil {
		logger.Error("failed to search albums", zap.Error(err))
		return nil, albumErrors.NewInternalError("failed to search albums: %v", err)
	}

	return albums, nil
}
