package repository

import (
	"context"
	"database/sql"
	"errors"
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
		SELECT id, title, type, thumbnail_url, release_date
		FROM album
		LEFT JOIN album_stats ON album.id = album_stats.album_id
		ORDER BY album_stats.listeners_count DESC, id DESC
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
	GetAlbumsByIDsQuery = `
		SELECT id, title, type, thumbnail_url, release_date
		FROM album
		LEFT JOIN album_stats ON album.id = album_stats.album_id
		WHERE id = ANY($1)
		ORDER BY album_stats.listeners_count DESC, id DESC
	`

	CreateStreamQuery = `
		INSERT INTO album_stream (album_id, user_id)
		VALUES ($1, $2)
	`
)

type albumPostgresRepository struct {
	db *sql.DB
	metrics *metrics.Metrics
}

func NewAlbumPostgresRepository(db *sql.DB, metrics *metrics.Metrics) domain.Repository {
	return &albumPostgresRepository{
		db: db,
		metrics: metrics,
	}
}

func (r *albumPostgresRepository) GetAllAlbums(ctx context.Context, filters *repoModel.AlbumFilters) ([]*repoModel.Album, error) {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting all albums from db", zap.Any("filters", filters), zap.String("query", GetAllAlbumsQuery))
	rows, err := r.db.Query(GetAllAlbumsQuery, filters.Pagination.Limit, filters.Pagination.Offset)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetAllAlbums").Inc()
		logger.Error("failed to get all albums", zap.Error(err))
		return nil, albumErrors.NewInternalError("failed to get all albums: %v", err)
	}
	defer rows.Close()

	albums := make([]*repoModel.Album, 0)
	for rows.Next() {
		var album repoModel.Album
		err = rows.Scan(&album.ID, &album.Title, &album.Type, &album.Thumbnail, &album.ReleaseDate)
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

func (r *albumPostgresRepository) GetAlbumByID(ctx context.Context, id int64) (*repoModel.Album, error) {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting album by id from db", zap.Int64("id", id), zap.String("query", GetAlbumByIDQuery))
	row := r.db.QueryRow(GetAlbumByIDQuery, id)

	var albumObject repoModel.Album
	err := row.Scan(&albumObject.ID, &albumObject.Title, &albumObject.Type, &albumObject.Thumbnail, &albumObject.ReleaseDate)
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
	rows, err := r.db.Query(GetAlbumTitleByIDsQuery, pq.Array(ids))
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetAlbumTitleByIDs").Inc()
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
	row := r.db.QueryRow(GetAlbumTitleByIDQuery, id)

	var title string
	err := row.Scan(&title)
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

func (r *albumPostgresRepository) GetAlbumsByIDs(ctx context.Context, ids []int64) ([]*repoModel.Album, error) {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting albums by ids from db", zap.Any("ids", ids), zap.String("query", GetAlbumsByIDsQuery))
	rows, err := r.db.Query(GetAlbumsByIDsQuery, pq.Array(ids))
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetAlbumsByIDs").Inc()
		logger.Error("failed to get albums by ids", zap.Error(err))
		return nil, albumErrors.NewInternalError("failed to get albums by ids: %v", err)
	}
	defer rows.Close()

	albums := make([]*repoModel.Album, 0)
	for rows.Next() {
		var album repoModel.Album
		err = rows.Scan(&album.ID, &album.Title, &album.Type, &album.Thumbnail, &album.ReleaseDate)
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
	_, err := r.db.Exec(CreateStreamQuery, albumID, userID)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("CreateStream").Inc()
		logger.Error("failed to create stream", zap.Error(err))
		return albumErrors.NewInternalError("failed to create stream: %v", err)
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("CreateStream").Observe(duration)
	return nil
}
