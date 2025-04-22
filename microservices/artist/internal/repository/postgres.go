package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/customErrors"
	loggerPkg "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
	domain "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/artist/internal/domain"
	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/artist/model/repository"
	"github.com/lib/pq"
	"go.uber.org/zap"
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

	GetArtistsByTrackIDsQuery = `
		SELECT a.id, a.title, ta.role, ta.track_id
		FROM artist a
		JOIN track_artist ta ON ta.artist_id = a.id
		WHERE ta.track_id = ANY($1)
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

	GetArtistsByAlbumIDsQuery = `
		SELECT a.id, a.title, aa.album_id
		FROM artist a
		JOIN album_artist aa ON a.id = aa.artist_id
		WHERE aa.album_id = ANY($1)
		ORDER BY aa.created_at, aa.id
	`

	GetAlbumIDsByArtistIDQuery = `
		SELECT album_id
		FROM album_artist
		WHERE artist_id = $1
	`
)

type artistPostgresRepository struct {
	db *sql.DB
}

func NewArtistPostgresRepository(db *sql.DB) domain.Repository {
	return &artistPostgresRepository{db: db}
}

func (r *artistPostgresRepository) GetAllArtists(ctx context.Context, filters *repoModel.ArtistFilters) ([]*repoModel.Artist, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting all artists with filters from db", zap.Any("filters", filters), zap.String("query", GetAllArtistsQuery))
	rows, err := r.db.QueryContext(ctx, GetAllArtistsQuery, filters.Pagination.Limit, filters.Pagination.Offset)
	if err != nil {
		logger.Error("failed to get all artists", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	artists := make([]*repoModel.Artist, 0)
	for rows.Next() {
		var artist repoModel.Artist
		err = rows.Scan(&artist.ID, &artist.Title, &artist.Description, &artist.Thumbnail)
		if err != nil {
			logger.Error("failed to scan artist", zap.Error(err))
			return nil, err
		}
		artists = append(artists, &artist)
	}

	return artists, nil
}

func (r *artistPostgresRepository) GetArtistByID(ctx context.Context, id int64) (*repoModel.Artist, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting artist by id from db", zap.Int64("id", id), zap.String("query", GetArtistByIDQuery))
	row := r.db.QueryRowContext(ctx, GetArtistByIDQuery, id)

	var artistObject repoModel.Artist
	err := row.Scan(&artistObject.ID, &artistObject.Title, &artistObject.Description, &artistObject.Thumbnail)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error("artist not found", zap.Error(err))
			return nil, customErrors.ErrArtistNotFound
		}
		logger.Error("failed to get artist by id", zap.Error(err))
		return nil, err
	}

	return &artistObject, nil
}

func (r *artistPostgresRepository) GetArtistTitleByID(ctx context.Context, id int64) (string, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting artist title by id from db", zap.Int64("id", id), zap.String("query", GetArtistTitleByIDQuery))
	row := r.db.QueryRowContext(ctx, GetArtistTitleByIDQuery, id)

	var title string
	err := row.Scan(&title)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error("artist not found", zap.Error(err))
			return "", customErrors.ErrArtistNotFound
		}
		logger.Error("failed to get artist title by id", zap.Error(err))
		return "", err
	}

	return title, nil
}

func (r *artistPostgresRepository) GetArtistsByTrackID(ctx context.Context, id int64) ([]*repoModel.ArtistWithRole, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting artists by track id from db", zap.Int64("id", id), zap.String("query", GetArtistsByTrackIDQuery))
	rows, err := r.db.QueryContext(ctx, GetArtistsByTrackIDQuery, id)
	if err != nil {
		logger.Error("failed to get artists by track id", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	artists := make([]*repoModel.ArtistWithRole, 0)
	for rows.Next() {
		var artist repoModel.ArtistWithRole
		err := rows.Scan(&artist.ID, &artist.Title, &artist.Role)
		if err != nil {
			logger.Error("failed to scan artist", zap.Error(err))
			return nil, err
		}
		artists = append(artists, &artist)
	}

	return artists, nil
}

func (r *artistPostgresRepository) GetArtistsByTrackIDs(ctx context.Context, trackIDs []int64) (map[int64][]*repoModel.ArtistWithRole, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting artists by track ids from db", zap.Any("ids", trackIDs), zap.String("query", GetArtistsByTrackIDsQuery))
	rows, err := r.db.QueryContext(ctx, GetArtistsByTrackIDsQuery, pq.Array(trackIDs))
	if err != nil {
		logger.Error("failed to get artists by track ids", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	artists := make(map[int64][]*repoModel.ArtistWithRole)
	for rows.Next() {
		var artist repoModel.ArtistWithRole
		var id int64
		err := rows.Scan(&artist.ID, &artist.Title, &artist.Role, &id)
		if err != nil {
			logger.Error("failed to scan artist", zap.Error(err))
			return nil, err
		}
		artists[id] = append(artists[id], &artist)
	}

	return artists, nil
}

func (r *artistPostgresRepository) GetArtistStats(ctx context.Context, id int64) (*repoModel.ArtistStats, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting artist stats by id from db", zap.Int64("id", id), zap.String("query", GetArtistStatsQuery))
	row := r.db.QueryRowContext(ctx, GetArtistStatsQuery, id)

	var stats repoModel.ArtistStats
	err := row.Scan(&stats.ListenersCount, &stats.FavoritesCount)
	if err != nil {
		logger.Error("failed to get artist stats by id", zap.Error(err))
		return nil, err
	}

	return &stats, nil
}

func (r *artistPostgresRepository) GetArtistsByAlbumID(ctx context.Context, albumID int64) ([]*repoModel.ArtistWithTitle, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting artists by album id from db", zap.Int64("id", albumID), zap.String("query", GetArtistsByAlbumIDQuery))
	rows, err := r.db.QueryContext(ctx, GetArtistsByAlbumIDQuery, albumID)
	if err != nil {
		logger.Error("failed to get artists by album id", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	artists := make([]*repoModel.ArtistWithTitle, 0)
	for rows.Next() {
		var artist repoModel.ArtistWithTitle
		err := rows.Scan(&artist.ID, &artist.Title)
		if err != nil {
			logger.Error("failed to scan artist", zap.Error(err))
			return nil, err
		}
		artists = append(artists, &artist)
	}

	return artists, nil
}

func (r *artistPostgresRepository) GetArtistsByAlbumIDs(ctx context.Context, albumIDs []int64) (map[int64][]*repoModel.ArtistWithTitle, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting artists by album ids from db", zap.Any("ids", albumIDs), zap.String("query", GetArtistsByAlbumIDsQuery))
	rows, err := r.db.QueryContext(ctx, GetArtistsByAlbumIDsQuery, pq.Array(albumIDs))
	if err != nil {
		logger.Error("failed to get artists by album ids", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	artists := make(map[int64][]*repoModel.ArtistWithTitle)
	for rows.Next() {
		var artist repoModel.ArtistWithTitle
		var id int64
		err := rows.Scan(&artist.ID, &artist.Title, &id)
		if err != nil {
			logger.Error("failed to scan artist", zap.Error(err))
			return nil, err
		}
		artists[id] = append(artists[id], &artist)
	}

	return artists, nil
}

func (r *artistPostgresRepository) GetAlbumIDsByArtistID(ctx context.Context, id int64) ([]int64, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting album ids by artist id from db", zap.Int64("id", id), zap.String("query", GetAlbumIDsByArtistIDQuery))
	rows, err := r.db.QueryContext(ctx, GetAlbumIDsByArtistIDQuery, id)
	if err != nil {
		logger.Error("failed to get album ids by artist id", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	albumIDs := make([]int64, 0)
	for rows.Next() {
		var albumID int64
		err := rows.Scan(&albumID)
		if err != nil {
			logger.Error("failed to scan album id", zap.Error(err))
			return nil, err
		}
		albumIDs = append(albumIDs, albumID)
	}

	return albumIDs, nil
}
