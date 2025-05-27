package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	loggerPkg "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
	domain "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/album/internal/domain"
	albumErrors "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/album/model/errors"
	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/album/model/repository"
	metrics "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/metrics"
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

	CreateAlbumQuery = `
		INSERT INTO album (title, type, thumbnail_url, label_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	DeleteAlbumQuery = `
		DELETE FROM album
		WHERE id = $1
	`
	GetAlbumsLabelIDQuery = `
		SELECT a.id, a.title, a.type, a.thumbnail_url, a.release_date, FALSE AS is_favorite
		FROM album a
		JOIN album_stats als ON a.id = als.album_id
		WHERE a.label_id = $1
		ORDER BY als.listeners_count DESC, a.id DESC
		LIMIT $2 OFFSET $3
	`
)

type albumPostgresRepository struct {
	db      *sql.DB
	metrics *metrics.Metrics
}

func NewAlbumPostgresRepository(db *sql.DB, metrics *metrics.Metrics) domain.Repository {
	return &albumPostgresRepository{
		db:      db,
		metrics: metrics,
	}
}

func (r *albumPostgresRepository) GetAllAlbums(ctx context.Context, filters *repoModel.AlbumFilters, userID int64) ([]*repoModel.Album, error) {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting all albums from db", zap.Any("filters", filters), zap.String("query", GetAllAlbumsQuery))

	stmt, err := r.db.PrepareContext(ctx, GetAllAlbumsQuery)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetAllAlbums").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return nil, albumErrors.NewInternalError("failed to prepare statement: %v", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			logger.Error("failed to close statement", zap.Error(err))
		}
	}()

	rows, err := stmt.QueryContext(ctx, filters.Pagination.Limit, filters.Pagination.Offset, userID)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetAllAlbums").Inc()
		logger.Error("failed to get all albums", zap.Error(err))
		return nil, albumErrors.NewInternalError("failed to get all albums: %v", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Error("failed to close rows", zap.Error(err))
		}
	}()

	albums := make([]*repoModel.Album, 0)
	for rows.Next() {
		var album repoModel.Album
		err = rows.Scan(&album.ID, &album.Title, &album.Type, &album.Thumbnail, &album.ReleaseDate, &album.IsFavorite)
		if err != nil {
			r.metrics.DatabaseErrors.WithLabelValues("GetAllAlbums").Inc()
			logger.Error("failed to scan album", zap.Error(err))
			return nil, albumErrors.NewInternalError("failed to scan album: %v", err)
		}
		albums = append(albums, &album)
	}

	if err := rows.Err(); err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetAllAlbums").Inc()
		logger.Error("failed to get all albums", zap.Error(err))
		return nil, albumErrors.NewInternalError("failed to get all albums: %v", err)
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("GetAllAlbums").Observe(duration)
	return albums, nil
}

func (r *albumPostgresRepository) GetAlbumByID(ctx context.Context, id int64, userID int64) (*repoModel.Album, error) {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting album by id from db", zap.Int64("id", id), zap.String("query", GetAlbumByIDQuery))

	stmt, err := r.db.PrepareContext(ctx, GetAlbumByIDQuery)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetAlbumByID").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return nil, albumErrors.NewInternalError("failed to prepare statement: %v", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			logger.Error("failed to close statement", zap.Error(err))
		}
	}()

	row := stmt.QueryRowContext(ctx, id, userID)

	var albumObject repoModel.Album
	err = row.Scan(&albumObject.ID, &albumObject.Title, &albumObject.Type, &albumObject.Thumbnail, &albumObject.ReleaseDate, &albumObject.IsFavorite)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetAlbumByID").Inc()
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error("album not found", zap.Error(err))
			return nil, albumErrors.ErrAlbumNotFound
		}
		logger.Error("failed to get album by id", zap.Error(err))
		return nil, albumErrors.NewInternalError("failed to get album by id: %v", err)
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("GetAlbumByID").Observe(duration)
	return &albumObject, nil
}

func (r *albumPostgresRepository) GetAlbumTitleByIDs(ctx context.Context, ids []int64) (map[int64]string, error) {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting album title by ids from db", zap.Any("ids", ids), zap.String("query", GetAlbumTitleByIDsQuery))

	stmt, err := r.db.PrepareContext(ctx, GetAlbumTitleByIDsQuery)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetAlbumTitleByIDs").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return nil, albumErrors.NewInternalError("failed to prepare statement: %v", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			logger.Error("failed to close statement", zap.Error(err))
		}
	}()

	rows, err := stmt.QueryContext(ctx, pq.Array(ids))
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetAlbumTitleByIDs").Inc()
		logger.Error("failed to get album title by ids", zap.Error(err))
		return nil, albumErrors.NewInternalError("failed to get album title by ids: %v", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Error("failed to close rows", zap.Error(err))
		}
	}()

	albums := make(map[int64]string)
	for rows.Next() {
		var id int64
		var title string
		err = rows.Scan(&id, &title)
		if err != nil {
			r.metrics.DatabaseErrors.WithLabelValues("GetAlbumTitleByIDs").Inc()
			logger.Error("failed to scan album title", zap.Error(err))
			return nil, albumErrors.NewInternalError("failed to scan album title: %v", err)
		}
		albums[id] = title
	}

	if err := rows.Err(); err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetAlbumTitleByIDs").Inc()
		logger.Error("failed to get album title by ids", zap.Error(err))
		return nil, albumErrors.NewInternalError("failed to get album title by ids: %v", err)
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("GetAlbumTitleByIDs").Observe(duration)
	return albums, nil
}

func (r *albumPostgresRepository) GetAlbumTitleByID(ctx context.Context, id int64) (string, error) {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting album title by id from db", zap.Int64("id", id), zap.String("query", GetAlbumTitleByIDQuery))

	stmt, err := r.db.PrepareContext(ctx, GetAlbumTitleByIDQuery)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetAlbumTitleByID").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return "", albumErrors.NewInternalError("failed to prepare statement: %v", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			logger.Error("failed to close statement", zap.Error(err))
		}
	}()

	row := stmt.QueryRowContext(ctx, id)

	var title string
	err = row.Scan(&title)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetAlbumTitleByID").Inc()
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error("album not found", zap.Error(err))
			return "", albumErrors.ErrAlbumNotFound
		}
		logger.Error("failed to get album title by id", zap.Error(err))
		return "", albumErrors.NewInternalError("failed to get album title by id: %v", err)
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("GetAlbumTitleByID").Observe(duration)
	return title, nil
}

func (r *albumPostgresRepository) GetAlbumsByIDs(ctx context.Context, ids []int64, userID int64) ([]*repoModel.Album, error) {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting albums by ids from db", zap.Any("ids", ids), zap.String("query", GetAlbumsByIDsQuery))

	stmt, err := r.db.PrepareContext(ctx, GetAlbumsByIDsQuery)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetAlbumsByIDs").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return nil, albumErrors.NewInternalError("failed to prepare statement: %v", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			logger.Error("failed to close statement", zap.Error(err))
		}
	}()

	rows, err := stmt.QueryContext(ctx, pq.Array(ids), userID)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetAlbumsByIDs").Inc()
		logger.Error("failed to get albums by ids", zap.Error(err))
		return nil, albumErrors.NewInternalError("failed to get albums by ids: %v", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Error("failed to close rows", zap.Error(err))
		}
	}()

	var albums []*repoModel.Album
	for rows.Next() {
		var album repoModel.Album
		err = rows.Scan(&album.ID, &album.Title, &album.Type, &album.Thumbnail, &album.ReleaseDate, &album.IsFavorite)
		if err != nil {
			r.metrics.DatabaseErrors.WithLabelValues("GetAlbumsByIDs").Inc()
			logger.Error("failed to scan album", zap.Error(err))
			return nil, albumErrors.NewInternalError("failed to scan album: %v", err)
		}
		albums = append(albums, &album)
	}

	if err := rows.Err(); err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetAlbumsByIDs").Inc()
		logger.Error("failed to get albums by ids", zap.Error(err))
		return nil, albumErrors.NewInternalError("failed to get albums by ids: %v", err)
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("GetAlbumsByIDs").Observe(duration)
	return albums, nil
}

func (r *albumPostgresRepository) CreateStream(ctx context.Context, albumID int64, userID int64) error {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Creating stream for album", zap.Int64("albumID", albumID), zap.Int64("userID", userID), zap.String("query", CreateStreamQuery))
	stmt, err := r.db.PrepareContext(ctx, CreateStreamQuery)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("CreateStream").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return albumErrors.NewInternalError("failed to prepare statement: %v", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			logger.Error("failed to close statement", zap.Error(err))
		}
	}()

	_, err = stmt.ExecContext(ctx, albumID, userID)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("CreateStream").Inc()
		logger.Error("failed to create stream", zap.Error(err))
		return albumErrors.NewInternalError("failed to create stream: %v", err)
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("CreateStream").Observe(duration)
	return nil
}

func (r *albumPostgresRepository) CheckAlbumExists(ctx context.Context, albumID int64) (bool, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Checking if album exists in db", zap.Int64("albumID", albumID), zap.String("query", CheckAlbumExistsQuery))
	stmt, err := r.db.PrepareContext(ctx, CheckAlbumExistsQuery)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("CreateStream").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return false, albumErrors.NewInternalError("failed to prepare statement: %v", err)
	}
	row := stmt.QueryRowContext(ctx, albumID)

	var exists bool
	err = row.Scan(&exists)
	if err != nil {
		logger.Error("failed to check if album exists", zap.Error(err))
		return false, albumErrors.NewInternalError("failed to check if album exists: %v", err)
	}
	return exists, nil
}

func (r *albumPostgresRepository) LikeAlbum(ctx context.Context, request *repoModel.LikeRequest) error {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Liking album", zap.Int64("albumID", request.AlbumID), zap.Int64("userID", request.UserID), zap.String("query", LikeAlbumQuery))
	stmt, err := r.db.PrepareContext(ctx, LikeAlbumQuery)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("Like album").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return albumErrors.NewInternalError("failed to prepare statement: %v", err)
	}
	_, err = stmt.ExecContext(ctx, request.AlbumID, request.UserID)
	if err != nil {
		logger.Error("failed to like album", zap.Error(err))
		return albumErrors.NewInternalError("failed to like album: %v", err)
	}
	return nil
}

func (r *albumPostgresRepository) UnlikeAlbum(ctx context.Context, request *repoModel.LikeRequest) error {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Unliking album", zap.Int64("albumID", request.AlbumID), zap.Int64("userID", request.UserID), zap.String("query", UnlikeAlbumQuery))
	stmt, err := r.db.PrepareContext(ctx, UnlikeAlbumQuery)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("Unlike album").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return albumErrors.NewInternalError("failed to prepare statement: %v", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			logger.Error("failed to close statement", zap.Error(err))
		}
	}()

	_, err = stmt.ExecContext(ctx, request.AlbumID, request.UserID)
	if err != nil {
		logger.Error("failed to unlike album", zap.Error(err))
		return albumErrors.NewInternalError("failed to unlike album: %v", err)
	}
	return nil
}

func (r *albumPostgresRepository) GetFavoriteAlbums(ctx context.Context, filters *repoModel.AlbumFilters, userID int64) ([]*repoModel.Album, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting favorite albums from db", zap.Any("filters", filters), zap.String("query", GetFavoriteAlbumsQuery))

	stmt, err := r.db.PrepareContext(ctx, GetFavoriteAlbumsQuery)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetFavoriteAlbums").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return nil, albumErrors.NewInternalError("failed to prepare statement: %v", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			logger.Error("failed to close statement", zap.Error(err))
		}
	}()
	rows, err := stmt.QueryContext(ctx, userID, filters.Pagination.Limit, filters.Pagination.Offset)
	if err != nil {
		logger.Error("failed to get favorite albums", zap.Error(err))
		return nil, albumErrors.NewInternalError("failed to get favorite albums: %v", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Error("failed to close rows", zap.Error(err))
		}
	}()

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
	stmt, err := r.db.PrepareContext(ctx, SearchAlbumsQuery)
	if err != nil {
		logger.Error("failed to prepare statement", zap.Error(err))
		return nil, albumErrors.NewInternalError("failed to prepare statement: %v", err)
	}
	words := strings.Fields(query)
	for i, word := range words {
		words[i] = word + ":*"
	}
	tsQueryString := strings.Join(words, " & ")

	rows, err := stmt.QueryContext(ctx, tsQueryString, userID, query)
	if err != nil {
		logger.Error("failed to search albums", zap.Error(err))
		return nil, albumErrors.NewInternalError("failed to search albums: %v", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Error("failed to close rows", zap.Error(err))
		}
	}()

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

func (r *albumPostgresRepository) CreateAlbum(ctx context.Context, album *repoModel.CreateAlbumRequest) (int64, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Creating album in db", zap.Any("album", album), zap.String("query", "CreateAlbum"))

	stmt, err := r.db.PrepareContext(ctx, CreateAlbumQuery)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("CreateAlbum").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return 0, albumErrors.NewInternalError("failed to prepare statement: %v", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			logger.Error("failed to close statement", zap.Error(err))
		}
	}()
	var albumID int64
	err = stmt.QueryRowContext(ctx, album.Title, album.Type, album.Thumbnail, album.LabelID).Scan(&albumID)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("CreateAlbum").Inc()
		logger.Error("failed to create album", zap.Error(err))
		return 0, albumErrors.NewInternalError("failed to create album: %v", err)
	}

	return albumID, nil
}

func (r *albumPostgresRepository) DeleteAlbum(ctx context.Context, albumID int64) error {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Deleting album from db", zap.Int64("albumID", albumID), zap.String("query", "DeleteAlbum"))

	stmt, err := r.db.PrepareContext(ctx, DeleteAlbumQuery)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("DeleteAlbum").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return albumErrors.NewInternalError("failed to prepare statement: %v", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			logger.Error("failed to close statement", zap.Error(err))
		}
	}()

	_, err = stmt.ExecContext(ctx, albumID)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("DeleteAlbum").Inc()
		logger.Error("failed to delete album", zap.Error(err))
		return albumErrors.NewInternalError("failed to delete album: %v", err)
	}

	return nil
}

func (r *albumPostgresRepository) GetAlbumsLabelID(ctx context.Context, filters *repoModel.AlbumFilters, labelID int64) ([]*repoModel.Album, error) {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting albums by label ID from db", zap.Int64("labelID", labelID), zap.Any("filters", filters), zap.String("query", "GetAlbumsLabelID"))

	stmt, err := r.db.PrepareContext(ctx, GetAlbumsLabelIDQuery)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetAlbumsLabelID").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return nil, albumErrors.NewInternalError("failed to prepare statement: %v", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			logger.Error("failed to close statement", zap.Error(err))
		}
	}()

	rows, err := stmt.QueryContext(ctx, labelID, filters.Pagination.Limit, filters.Pagination.Offset)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetAlbumsLabelID").Inc()
		logger.Error("failed to get all albums", zap.Error(err))
		return nil, albumErrors.NewInternalError("failed to get all albums: %v", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Error("failed to close rows", zap.Error(err))
		}
	}()

	albums := make([]*repoModel.Album, 0)
	for rows.Next() {
		var album repoModel.Album
		var isFavorite sql.NullBool
		err = rows.Scan(&album.ID, &album.Title, &album.Type, &album.Thumbnail, &album.ReleaseDate, &isFavorite)
		if err != nil {
			r.metrics.DatabaseErrors.WithLabelValues("GetAlbumsLabelID").Inc()
			logger.Error("failed to scan album", zap.Error(err))
			return nil, albumErrors.NewInternalError("failed to scan album: %v", err)
		}
		album.IsFavorite = isFavorite.Valid && isFavorite.Bool
		albums = append(albums, &album)
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("GetAlbumsLabelID").Observe(duration)
	return albums, nil
}
