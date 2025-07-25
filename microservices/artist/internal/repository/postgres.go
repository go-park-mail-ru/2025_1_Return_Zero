package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	loggerPkg "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
	domain "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/artist/internal/domain"
	artistErrors "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/artist/model/errors"
	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/artist/model/repository"
	metrics "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/metrics"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

const (
	GetAllArtistsQuery = `
		SELECT artist.id, artist.title, artist.description, artist.thumbnail_url, (favorite_artist.user_id IS NOT NULL) AS is_favorite
		FROM artist
		JOIN artist_stats ON artist.id = artist_stats.artist_id
		LEFT JOIN favorite_artist ON artist.id = favorite_artist.artist_id AND favorite_artist.user_id = $3
		ORDER BY artist_stats.listeners_count DESC, id DESC
		LIMIT $1 OFFSET $2
	`
	GetArtistByIDQuery = `
		SELECT artist.id, artist.title, artist.description, artist.thumbnail_url, (favorite_artist.user_id IS NOT NULL) AS is_favorite
		FROM artist
		LEFT JOIN favorite_artist ON artist.id = favorite_artist.artist_id AND favorite_artist.user_id = $2
		WHERE artist.id = $1
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
		ORDER BY aa.created_at DESC, aa.id DESC
	`

	GetArtistsByAlbumIDsQuery = `
		SELECT a.id, a.title, aa.album_id
		FROM artist a
		JOIN album_artist aa ON a.id = aa.artist_id
		WHERE aa.album_id = ANY($1)
		ORDER BY aa.created_at DESC, aa.id DESC
	`

	GetAlbumIDsByArtistIDQuery = `
		SELECT album_id
		FROM album_artist
		WHERE artist_id = $1
	`

	GetTrackIDsByArtistID = `
		SELECT track_id
		FROM track_artist
		WHERE artist_id = $1
	`

	GetArtistsListenedByUserIDQuery = `
		SELECT COUNT(DISTINCT artist_id)
		FROM artist_stream
		WHERE user_id = $1
	`

	LikeArtistByUserIDQuery = `
		INSERT INTO favorite_artist (artist_id, user_id) VALUES ($1, $2)
		ON CONFLICT (artist_id, user_id) DO NOTHING
	`

	UnlikeArtistByUserIDQuery = `
		DELETE FROM favorite_artist WHERE artist_id = $1 AND user_id = $2
	`

	CheckArtistExistsQuery = `
		SELECT EXISTS (SELECT 1 FROM artist WHERE id = $1)
	`

	GetFavoriteArtistsQuery = `
		SELECT artist.id, artist.title, artist.description, artist.thumbnail_url
		FROM artist
		JOIN favorite_artist ON artist.id = favorite_artist.artist_id
		WHERE favorite_artist.user_id = $1
		ORDER BY favorite_artist.created_at DESC, artist.id DESC
		LIMIT $2 OFFSET $3
	`

	SearchArtistsQuery = `
		SELECT a.id, a.title, a.description, a.thumbnail_url
		FROM artist a
		LEFT JOIN favorite_artist fa ON a.id = fa.artist_id AND fa.user_id = $2
		WHERE a.search_vector @@ to_tsquery('multilingual', $1)
		   OR similarity(a.title_trgm, $3) > 0.3
		ORDER BY 
		    CASE WHEN a.search_vector @@ to_tsquery('multilingual', $1) THEN 0 ELSE 1 END,
		    ts_rank(a.search_vector, to_tsquery('multilingual', $1)) DESC,
		    similarity(a.title_trgm, $3) DESC
	`

	CreateArtistQuery = `
		INSERT INTO artist (title, thumbnail_url, label_id)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	CheckArtistNameExist = `
		SELECT 1 
		FROM artist
		WHERE id = $1
	`

	UpdateArtistNameQuery = `
		UPDATE artist
		SET title = $1
		WHERE id = $2
	`

	GetArtistByIdWithoutUserQuery = `
		SELECT id, title, description, thumbnail_url
		FROM artist
		WHERE id = $1
	`

	UpdateArtistAvatarQuery = `
		UPDATE artist
		SET thumbnail_url = $1
		WHERE id = $2
	`

	GetArtistLabelIdByIDQuery = `
		SELECT label_id
		FROM artist
		WHERE id = $1
	`

	GetArtistsLabelIDQuery = `
        SELECT artist.id, artist.title, artist.description, artist.thumbnail_url, FALSE AS is_favorite
        FROM artist
        JOIN artist_stats ON artist.id = artist_stats.artist_id
        WHERE artist.label_id = $3
        ORDER BY artist_stats.listeners_count DESC, id DESC
        LIMIT $1 OFFSET $2
    `

	DeleteArtistQuery = `
		DELETE 
		FROM artist
		WHERE id = $1
	`

	AddArtistsToAlbum = `
		INSERT INTO album_artist (artist_id, album_id) 
		SELECT unnest($1::bigint[]), $2
	`

	AddArtistsToTracks = `
		INSERT INTO track_artist (artist_id, track_id, role)
        SELECT a, t, 'main'
        FROM unnest($1::bigint[]) AS a
        CROSS JOIN unnest($2::bigint[]) AS t
        ON CONFLICT (track_id, artist_id, role) DO NOTHING
	`
)

type artistPostgresRepository struct {
	db      *sql.DB
	metrics *metrics.Metrics
}

func NewArtistPostgresRepository(db *sql.DB, metrics *metrics.Metrics) domain.Repository {
	return &artistPostgresRepository{db: db, metrics: metrics}
}

func (r *artistPostgresRepository) GetAllArtists(ctx context.Context, filters *repoModel.Filters, userID int64) ([]*repoModel.Artist, error) {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting all artists with filters from db", zap.Any("filters", filters), zap.String("query", GetAllArtistsQuery))
	stmt, err := r.db.PrepareContext(ctx, GetAllArtistsQuery)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetAllArtitst").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return nil, artistErrors.NewInternalError("failed to prepare statement: %v", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			logger.Error("failed to close statement", zap.Error(err))
		}
	}()
	rows, err := stmt.QueryContext(ctx, filters.Pagination.Limit, filters.Pagination.Offset, userID)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetAllArtitst").Inc()
		logger.Error("failed to get all artists", zap.Error(err))
		return nil, artistErrors.NewInternalError("failed to get all artists: %v", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Error("failed to close rows", zap.Error(err))
		}
	}()

	artists := make([]*repoModel.Artist, 0)
	for rows.Next() {
		var artist repoModel.Artist
		var isFavorite sql.NullBool
		err = rows.Scan(&artist.ID, &artist.Title, &artist.Description, &artist.Thumbnail, &isFavorite)
		if err != nil {
			r.metrics.DatabaseErrors.WithLabelValues("GetAllArtitst").Inc()
			logger.Error("failed to scan artist", zap.Error(err))
			return nil, artistErrors.NewInternalError("failed to scan artist: %v", err)
		}
		artist.IsFavorite = isFavorite.Valid && isFavorite.Bool
		artists = append(artists, &artist)
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("GetAllArtists").Observe(duration)
	return artists, nil
}

func (r *artistPostgresRepository) GetArtistByID(ctx context.Context, id int64, userID int64) (*repoModel.Artist, error) {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting artist by id from db", zap.Int64("id", id), zap.String("query", GetArtistByIDQuery))
	stmt, err := r.db.PrepareContext(ctx, GetArtistByIDQuery)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetArtistByID").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return nil, artistErrors.NewInternalError("failed to prepare statement: %v", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			logger.Error("failed to close statement", zap.Error(err))
		}
	}()
	row := stmt.QueryRowContext(ctx, id, userID)

	var artistObject repoModel.Artist
	var isFavorite sql.NullBool
	err = row.Scan(&artistObject.ID, &artistObject.Title, &artistObject.Description, &artistObject.Thumbnail, &isFavorite)

	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetArtistByID").Inc()
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error("artist not found", zap.Error(err))
			return nil, artistErrors.ErrArtistNotFound
		}
		logger.Error("failed to get artist by id", zap.Error(err))
		return nil, artistErrors.NewInternalError("failed to get artist by id: %v", err)
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("GetArtistByID").Observe(duration)

	artistObject.IsFavorite = isFavorite.Valid && isFavorite.Bool

	return &artistObject, nil
}

func (r *artistPostgresRepository) GetArtistTitleByID(ctx context.Context, id int64) (string, error) {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting artist title by id from db", zap.Int64("id", id), zap.String("query", GetArtistTitleByIDQuery))
	stmt, err := r.db.PrepareContext(ctx, GetArtistTitleByIDQuery)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetArtistTitleByID").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return "", artistErrors.NewInternalError("failed to prepare statement: %v", err)
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
		r.metrics.DatabaseErrors.WithLabelValues("GetArtistTitleByID").Inc()
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error("artist not found", zap.Error(err))
			return "", artistErrors.ErrArtistNotFound
		}
		logger.Error("failed to get artist title by id", zap.Error(err))
		return "", artistErrors.NewInternalError("failed to get artist title by id: %v", err)
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("GetArtistTitleByID").Observe(duration)
	return title, nil
}

func (r *artistPostgresRepository) GetArtistsByTrackID(ctx context.Context, id int64) ([]*repoModel.ArtistWithRole, error) {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting artists by track id from db", zap.Int64("id", id), zap.String("query", GetArtistsByTrackIDQuery))
	stmt, err := r.db.PrepareContext(ctx, GetArtistsByTrackIDQuery)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetArtistsByTrackID").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return nil, artistErrors.NewInternalError("failed to prepare statement: %v", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			logger.Error("failed to close statement", zap.Error(err))
		}
	}()
	rows, err := stmt.QueryContext(ctx, id)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetArtistsByTrackID").Inc()
		logger.Error("failed to get artists by track id", zap.Error(err))
		return nil, artistErrors.NewInternalError("failed to get artists by track id: %v", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Error("failed to close rows", zap.Error(err))
		}
	}()

	artists := make([]*repoModel.ArtistWithRole, 0)
	for rows.Next() {
		var artist repoModel.ArtistWithRole
		err := rows.Scan(&artist.ID, &artist.Title, &artist.Role)
		if err != nil {
			r.metrics.DatabaseErrors.WithLabelValues("GetArtistsByTrackID").Inc()
			logger.Error("failed to scan artist", zap.Error(err))
			return nil, artistErrors.NewInternalError("failed to scan artist: %v", err)
		}
		artists = append(artists, &artist)
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("GetArtistsByTrackID").Observe(duration)
	return artists, nil
}

func (r *artistPostgresRepository) GetArtistsByTrackIDs(ctx context.Context, trackIDs []int64) (map[int64][]*repoModel.ArtistWithRole, error) {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting artists by track ids from db", zap.Any("ids", trackIDs), zap.String("query", GetArtistsByTrackIDsQuery))
	stmt, err := r.db.PrepareContext(ctx, GetArtistsByTrackIDsQuery)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetArtistsByTrackIDs").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return nil, artistErrors.NewInternalError("failed to prepare statement: %v", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			logger.Error("failed to close statement", zap.Error(err))
		}
	}()
	rows, err := stmt.QueryContext(ctx, pq.Array(trackIDs))
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetArtistsByTrackIDs").Inc()
		logger.Error("failed to get artists by track ids", zap.Error(err))
		return nil, artistErrors.NewInternalError("failed to get artists by track ids: %v", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Error("failed to close rows", zap.Error(err))
		}
	}()

	artists := make(map[int64][]*repoModel.ArtistWithRole)
	for rows.Next() {
		var artist repoModel.ArtistWithRole
		var id int64
		err := rows.Scan(&artist.ID, &artist.Title, &artist.Role, &id)
		if err != nil {
			r.metrics.DatabaseErrors.WithLabelValues("GetArtistsByTrackIDs").Inc()
			logger.Error("failed to scan artist", zap.Error(err))
			return nil, artistErrors.NewInternalError("failed to scan artist: %v", err)
		}
		artists[id] = append(artists[id], &artist)
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("GetArtistsByTrackIDs").Observe(duration)
	return artists, nil
}

func (r *artistPostgresRepository) GetArtistStats(ctx context.Context, id int64) (*repoModel.ArtistStats, error) {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting artist stats by id from db", zap.Int64("id", id), zap.String("query", GetArtistStatsQuery))
	stmt, err := r.db.PrepareContext(ctx, GetArtistStatsQuery)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetArtistStats").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return nil, artistErrors.NewInternalError("failed to prepare statement: %v", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			logger.Error("failed to close statement", zap.Error(err))
		}
	}()
	row := stmt.QueryRowContext(ctx, id)

	var stats repoModel.ArtistStats
	err = row.Scan(&stats.ListenersCount, &stats.FavoritesCount)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetArtistStats").Inc()
		logger.Error("failed to get artist stats by id", zap.Error(err))
		return nil, artistErrors.NewInternalError("failed to get artist stats by id: %v", err)
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("GetArtistStats").Observe(duration)
	return &stats, nil
}

func (r *artistPostgresRepository) GetArtistsByAlbumID(ctx context.Context, albumID int64) ([]*repoModel.ArtistWithTitle, error) {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting artists by album id from db", zap.Int64("id", albumID), zap.String("query", GetArtistsByAlbumIDQuery))
	stmt, err := r.db.PrepareContext(ctx, GetArtistsByAlbumIDQuery)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetArtistsByAlbumID").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return nil, artistErrors.NewInternalError("failed to prepare statement: %v", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			logger.Error("failed to close statement", zap.Error(err))
		}
	}()
	rows, err := stmt.QueryContext(ctx, albumID)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetArtistsByAlbumID").Inc()
		logger.Error("failed to get artists by album id", zap.Error(err))
		return nil, artistErrors.NewInternalError("failed to get artists by album id: %v", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Error("failed to close rows", zap.Error(err))
		}
	}()

	artists := make([]*repoModel.ArtistWithTitle, 0)
	for rows.Next() {
		var artist repoModel.ArtistWithTitle
		err := rows.Scan(&artist.ID, &artist.Title)
		if err != nil {
			r.metrics.DatabaseErrors.WithLabelValues("GetArtistsByAlbumID").Inc()
			logger.Error("failed to scan artist", zap.Error(err))
			return nil, artistErrors.NewInternalError("failed to scan artist: %v", err)
		}
		artists = append(artists, &artist)
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("GetArtistsByAlbumID").Observe(duration)
	return artists, nil
}

func (r *artistPostgresRepository) GetArtistsByAlbumIDs(ctx context.Context, albumIDs []int64) (map[int64][]*repoModel.ArtistWithTitle, error) {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting artists by album ids from db", zap.Any("ids", albumIDs), zap.String("query", GetArtistsByAlbumIDsQuery))
	stmt, err := r.db.PrepareContext(ctx, GetArtistsByAlbumIDsQuery)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetArtistsByAlbumIDs").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return nil, artistErrors.NewInternalError("failed to prepare statement: %v", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			logger.Error("failed to close statement", zap.Error(err))
		}
	}()
	rows, err := stmt.QueryContext(ctx, pq.Array(albumIDs))
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetArtistsByAlbumIDs").Inc()
		logger.Error("failed to get artists by album ids", zap.Error(err))
		return nil, artistErrors.NewInternalError("failed to get artists by album ids: %v", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Error("failed to close rows", zap.Error(err))
		}
	}()

	artists := make(map[int64][]*repoModel.ArtistWithTitle)
	for rows.Next() {
		var artist repoModel.ArtistWithTitle
		var albumID int64
		err := rows.Scan(&artist.ID, &artist.Title, &albumID)
		if err != nil {
			r.metrics.DatabaseErrors.WithLabelValues("GetArtistsByAlbumIDs").Inc()
			logger.Error("failed to scan artist", zap.Error(err))
			return nil, artistErrors.NewInternalError("failed to scan artist: %v", err)
		}
		artists[albumID] = append(artists[albumID], &artist)
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("GetArtistsByAlbumIDs").Observe(duration)
	return artists, nil
}

func (r *artistPostgresRepository) GetAlbumIDsByArtistID(ctx context.Context, id int64) ([]int64, error) {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting album ids by artist id from db", zap.Int64("id", id), zap.String("query", GetAlbumIDsByArtistIDQuery))
	stmt, err := r.db.PrepareContext(ctx, GetAlbumIDsByArtistIDQuery)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetAlbumIDsByArtistID").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return nil, artistErrors.NewInternalError("failed to prepare statement: %v", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			logger.Error("failed to close statement", zap.Error(err))
		}
	}()
	rows, err := stmt.QueryContext(ctx, id)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetAlbumIDsByArtistID").Inc()
		logger.Error("failed to get album ids by artist id", zap.Error(err))
		return nil, artistErrors.NewInternalError("failed to get album ids by artist id: %v", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Error("failed to close rows", zap.Error(err))
		}
	}()

	albumIDs := make([]int64, 0)
	for rows.Next() {
		var albumID int64
		err := rows.Scan(&albumID)
		if err != nil {
			r.metrics.DatabaseErrors.WithLabelValues("GetAlbumIDsByArtistID").Inc()
			logger.Error("failed to scan album id", zap.Error(err))
			return nil, artistErrors.NewInternalError("failed to scan album id: %v", err)
		}
		albumIDs = append(albumIDs, albumID)
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("GetAlbumIDsByArtistID").Observe(duration)
	return albumIDs, nil
}

func (r *artistPostgresRepository) GetTrackIDsByArtistID(ctx context.Context, id int64) ([]int64, error) {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting track ids by artist id from db", zap.Int64("id", id), zap.String("query", GetTrackIDsByArtistID))
	stmt, err := r.db.PrepareContext(ctx, GetTrackIDsByArtistID)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetTrackIDsByArtistID").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return nil, artistErrors.NewInternalError("failed to prepare statement: %v", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			logger.Error("failed to close statement", zap.Error(err))
		}
	}()
	rows, err := stmt.QueryContext(ctx, id)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetTrackIDsByArtistID").Inc()
		logger.Error("failed to get track ids by artist id", zap.Error(err))
		return nil, artistErrors.NewInternalError("failed to get track ids by artist id: %v", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Error("failed to close rows", zap.Error(err))
		}
	}()

	trackIDs := make([]int64, 0)
	for rows.Next() {
		var trackID int64
		err := rows.Scan(&trackID)
		if err != nil {
			r.metrics.DatabaseErrors.WithLabelValues("GetTrackIDsByArtistID").Inc()
			logger.Error("failed to scan track id", zap.Error(err))
			return nil, artistErrors.NewInternalError("failed to scan track id: %v", err)
		}
		trackIDs = append(trackIDs, trackID)
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("GetTrackIDsByArtistID").Observe(duration)
	return trackIDs, nil
}

func (r *artistPostgresRepository) CreateStreamsByArtistIDs(ctx context.Context, data *repoModel.ArtistStreamCreateDataList) error {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Creating streams for artists", zap.Any("data", data))

	if len(data.ArtistIDs) == 0 {
		r.metrics.DatabaseErrors.WithLabelValues("CreateStreamsByArtistIDs").Inc()
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("CreateStreamsByArtistIDs").Inc()
		logger.Error("failed to begin transaction", zap.Error(err))
		return artistErrors.NewInternalError("failed to begin transaction: %v", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			logger.Error("failed to rollback transaction", zap.Error(err))
		}
	}()

	query := `INSERT INTO artist_stream (artist_id, user_id) 
              SELECT unnest($1::bigint[]), $2`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("CreateStreamsByArtistIDs").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return artistErrors.NewInternalError("failed to prepare statement: %v", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			logger.Error("failed to close statement", zap.Error(err))
		}
	}()

	_, err = stmt.ExecContext(ctx, pq.Array(data.ArtistIDs), data.UserID)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("CreateStreamsByArtistIDs").Inc()
		logger.Error("failed to create streams for artists", zap.Error(err))
		return artistErrors.NewInternalError("failed to create streams for artists: %v", err)
	}

	if err := tx.Commit(); err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("CreateStreamsByArtistIDs").Inc()
		logger.Error("failed to commit transaction", zap.Error(err))
		return artistErrors.NewInternalError("failed to commit transaction: %v", err)
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("CreateStreamsByArtistIDs").Observe(duration)
	return nil
}

func (r *artistPostgresRepository) GetArtistsListenedByUserID(ctx context.Context, userID int64) (int64, error) {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting artists listened by user id from db", zap.Int64("userID", userID), zap.String("query", GetArtistsListenedByUserIDQuery))
	stmt, err := r.db.PrepareContext(ctx, GetArtistsListenedByUserIDQuery)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetArtistsListenedByUserID").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return 0, artistErrors.NewInternalError("failed to prepare statement: %v", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			logger.Error("failed to close statement", zap.Error(err))
		}
	}()

	row := stmt.QueryRowContext(ctx, userID)

	var artistsListened int64
	err = row.Scan(&artistsListened)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetArtistsListenedByUserID").Inc()
		logger.Error("failed to get artists listened by user id", zap.Error(err))
		return 0, artistErrors.NewInternalError("failed to get artists listened by user id: %v", err)
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("GetArtistsListenedByUserID").Observe(duration)
	return artistsListened, nil
}

func (r *artistPostgresRepository) CheckArtistExists(ctx context.Context, id int64) (bool, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Checking if artist exists", zap.Int64("id", id), zap.String("query", CheckArtistExistsQuery))
	stmt, err := r.db.PrepareContext(ctx, CheckArtistExistsQuery)
	if err != nil {
		logger.Error("failed to prepare statement", zap.Error(err))
		return false, artistErrors.NewInternalError("failed to prepare statement: %v", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			logger.Error("failed to close statement", zap.Error(err))
		}
	}()
	row := stmt.QueryRowContext(ctx, id)

	var exists bool
	err = row.Scan(&exists)
	if err != nil {
		logger.Error("failed to check if artist exists", zap.Error(err))
		return false, artistErrors.NewInternalError("failed to check if artist exists: %v", err)
	}

	return exists, nil
}

// Мы не проверяем, какое значение было у зафаворченного исполнителя, а просто задаем его новое значение игнорируя предидущее. Такое подход по идее должен избавить нас от лишних проверок и запросов в бд.
func (r *artistPostgresRepository) LikeArtist(ctx context.Context, request *repoModel.LikeRequest) error {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting to like artist", zap.Any("request", request), zap.Int64("artistID", request.ArtistID), zap.Int64("userID", request.UserID))
	stmt, err := r.db.PrepareContext(ctx, LikeArtistByUserIDQuery)
	if err != nil {
		logger.Error("failed to prepare statement", zap.Error(err))
		return artistErrors.NewInternalError("failed to prepare statement: %v", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			logger.Error("failed to close statement", zap.Error(err))
		}
	}()

	_, err = stmt.ExecContext(ctx, request.ArtistID, request.UserID)
	if err != nil {
		logger.Error("failed to like artist", zap.Error(err))
		return artistErrors.NewInternalError("failed to like artist: %v", err)
	}

	return nil
}

func (r *artistPostgresRepository) UnlikeArtist(ctx context.Context, request *repoModel.LikeRequest) error {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting to unlike artist", zap.Any("request", request), zap.Int64("artistID", request.ArtistID), zap.Int64("userID", request.UserID))
	stmt, err := r.db.PrepareContext(ctx, UnlikeArtistByUserIDQuery)
	if err != nil {
		logger.Error("failed to prepare statement", zap.Error(err))
		return artistErrors.NewInternalError("failed to prepare statement: %v", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			logger.Error("failed to close statement", zap.Error(err))
		}
	}()

	_, err = stmt.ExecContext(ctx, request.ArtistID, request.UserID)
	if err != nil {
		logger.Error("failed to unlike artist", zap.Error(err))
		return artistErrors.NewInternalError("failed to unlike artist: %v", err)
	}

	return nil
}

func (r *artistPostgresRepository) GetFavoriteArtists(ctx context.Context, filters *repoModel.Filters, userID int64) ([]*repoModel.Artist, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting favorite artists by user id from db", zap.Int64("userID", userID), zap.String("query", GetFavoriteArtistsQuery))
	stmt, err := r.db.PrepareContext(ctx, GetFavoriteArtistsQuery)
	if err != nil {
		logger.Error("failed to prepare statement", zap.Error(err))
		return nil, artistErrors.NewInternalError("failed to prepare statement: %v", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			logger.Error("failed to close statement", zap.Error(err))
		}
	}()

	rows, err := stmt.QueryContext(ctx, userID, filters.Pagination.Limit, filters.Pagination.Offset)
	if err != nil {
		logger.Error("failed to get favorite artists", zap.Error(err))
		return nil, artistErrors.NewInternalError("failed to get favorite artists: %v", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Error("failed to close rows", zap.Error(err))
		}
	}()

	var artists []*repoModel.Artist
	for rows.Next() {
		var artist repoModel.Artist
		err := rows.Scan(&artist.ID, &artist.Title, &artist.Description, &artist.Thumbnail)
		// Так как это мы не отображаем в списке, то можно не делать лишнюю проверку
		// В идеале поменяем к рк4
		artist.IsFavorite = false
		if err != nil {
			logger.Error("failed to scan artist", zap.Error(err))
			return nil, artistErrors.NewInternalError("failed to scan artist: %v", err)
		}
		artists = append(artists, &artist)
	}

	return artists, nil
}

func (r *artistPostgresRepository) SearchArtists(ctx context.Context, query string, userID int64) ([]*repoModel.Artist, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting search artists by query from db", zap.String("query", query), zap.Int64("userID", userID), zap.String("query", SearchArtistsQuery))
	stmt, err := r.db.PrepareContext(ctx, SearchArtistsQuery)
	if err != nil {
		logger.Error("failed to prepare statement", zap.Error(err))
		return nil, artistErrors.NewInternalError("failed to prepare statement: %v", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			logger.Error("failed to close statement", zap.Error(err))
		}
	}()

	words := strings.Fields(query)
	for i, word := range words {
		words[i] = word + ":*"
	}
	tsQueryString := strings.Join(words, " & ")

	rows, err := stmt.QueryContext(ctx, tsQueryString, userID, query)
	if err != nil {
		logger.Error("failed to search artists", zap.Error(err))
		return nil, artistErrors.NewInternalError("failed to search artists: %v", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Error("failed to close rows", zap.Error(err))
		}
	}()

	var artists []*repoModel.Artist
	for rows.Next() {
		var artist repoModel.Artist
		err := rows.Scan(&artist.ID, &artist.Title, &artist.Description, &artist.Thumbnail)
		if err != nil {
			logger.Error("failed to scan artist", zap.Error(err))
			return nil, artistErrors.NewInternalError("failed to scan artist: %v", err)
		}
		artists = append(artists, &artist)
	}
	return artists, nil
}

func (r *artistPostgresRepository) CreateArtist(ctx context.Context, artist *repoModel.Artist) (*repoModel.Artist, error) {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Creating artist", zap.Any("artist", artist))
	stmt, err := r.db.PrepareContext(ctx, CreateArtistQuery)
	if err != nil {
		logger.Error("failed to prepare statement", zap.Error(err))
		r.metrics.DatabaseErrors.WithLabelValues("CreateArtist").Inc()
		return nil, artistErrors.NewInternalError("failed to prepare statement: %v", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			logger.Error("failed to close statement", zap.Error(err))
		}
	}()

	var artistID int64
	err = stmt.QueryRowContext(ctx, artist.Title, artist.Thumbnail, artist.LabelID).Scan(&artistID)
	if err != nil {
		logger.Error("failed to create artist", zap.Error(err))
		r.metrics.DatabaseErrors.WithLabelValues("CreateArtist").Inc()
		return nil, artistErrors.NewInternalError("failed to create artist: %v", err)
	}
	artist.ID = artistID
	duration := time.Since(start).Seconds()

	_, err = r.db.ExecContext(ctx, "REFRESH MATERIALIZED VIEW CONCURRENTLY artist_stats")
	if err != nil {
		logger.Warn("failed to refresh artist_stats view, artist may not be visible for up to 1 minute", zap.Error(err))
	}

	r.metrics.DatabaseDuration.WithLabelValues("CreateArtist").Observe(duration)
	return artist, nil
}

func (r *artistPostgresRepository) CheckArtistNameExist(ctx context.Context, id int64) (bool, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Checking if artist name exists", zap.Int64("id", id))
	stmt, err := r.db.PrepareContext(ctx, CheckArtistNameExist)
	if err != nil {
		logger.Error("failed to prepare statement", zap.Error(err))
		return false, artistErrors.NewInternalError("failed to prepare statement: %v", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			logger.Error("failed to close statement", zap.Error(err))
		}
	}()

	var exists bool
	err = stmt.QueryRowContext(ctx, id).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		logger.Error("failed to check if artist name exists", zap.Error(err))
		return false, artistErrors.NewInternalError("failed to check if artist name exists: %v", err)
	}

	return exists, nil
}

func (r *artistPostgresRepository) ChangeArtistTitle(ctx context.Context, newTitle string, id int64) error {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Changing artist title", zap.Int64("id", id))
	stmt, err := r.db.PrepareContext(ctx, UpdateArtistNameQuery)
	if err != nil {
		logger.Error("failed to prepare statement", zap.Error(err))
		return err
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			logger.Error("failed to close statement", zap.Error(err))
		}
	}()

	_, err = stmt.ExecContext(ctx, newTitle, id)
	if err != nil {
		logger.Error("failed to change artist title", zap.Error(err))
		return err
	}
	return nil
}

func (r *artistPostgresRepository) GetArtistByIDWithoutUser(ctx context.Context, id int64) (*repoModel.Artist, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting artist by id from db", zap.Int64("id", id))
	stmt, err := r.db.PrepareContext(ctx, GetArtistByIdWithoutUserQuery)
	if err != nil {
		logger.Error("failed to prepare statement", zap.Error(err))
		return nil, artistErrors.NewInternalError("failed to prepare statement: %v", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			logger.Error("failed to close statement", zap.Error(err))
		}
	}()

	row := stmt.QueryRowContext(ctx, id)

	var artistObject repoModel.Artist
	err = row.Scan(&artistObject.ID, &artistObject.Title, &artistObject.Description, &artistObject.Thumbnail)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error("artist not found", zap.Error(err))
			return nil, artistErrors.ErrArtistNotFound
		}
		logger.Error("failed to get artist by title", zap.Error(err))
		return nil, artistErrors.NewInternalError("failed to get artist by title: %v", err)
	}

	return &artistObject, nil
}

func (r *artistPostgresRepository) UploadAvatar(ctx context.Context, artistID int64, avatarURL string) error {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Uploading artist avatar", zap.Int64("artistID", artistID), zap.String("avatarURL", avatarURL))
	stmt, err := r.db.PrepareContext(ctx, UpdateArtistAvatarQuery)
	if err != nil {
		logger.Error("failed to prepare statement", zap.Error(err))
		return err
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			logger.Error("failed to close statement", zap.Error(err))
		}
	}()

	_, err = stmt.ExecContext(ctx, avatarURL, artistID)
	if err != nil {
		logger.Error("failed to upload artist avatar", zap.Error(err))
		return err
	}
	return nil
}

func (r *artistPostgresRepository) GetArtistLabelID(ctx context.Context, artistID int64) (int64, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting artist label id by id from db", zap.Int64("artistID", artistID))
	stmt, err := r.db.PrepareContext(ctx, GetArtistLabelIdByIDQuery)
	if err != nil {
		logger.Error("failed to prepare statement", zap.Error(err))
		return 0, artistErrors.NewInternalError("failed to prepare statement: %v", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			logger.Error("failed to close statement", zap.Error(err))
		}
	}()

	row := stmt.QueryRowContext(ctx, artistID)

	var labelID int64
	err = row.Scan(&labelID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error("artist not found", zap.Error(err))
			return 0, artistErrors.ErrArtistNotFound
		}
		logger.Error("failed to get artist label id by title", zap.Error(err))
		return 0, artistErrors.NewInternalError("failed to get artist label id by title: %v", err)
	}

	return labelID, nil
}

func (r *artistPostgresRepository) GetArtistsLabelID(ctx context.Context, filters *repoModel.Filters, labelID int64) ([]*repoModel.Artist, error) {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Requesting all artists with filters from db", zap.Any("filters", filters), zap.String("query", GetAllArtistsQuery))
	stmt, err := r.db.PrepareContext(ctx, GetArtistsLabelIDQuery)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetArtistsLabelID").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return nil, artistErrors.NewInternalError("failed to prepare statement: %v", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			logger.Error("failed to close statement", zap.Error(err))
		}
	}()
	rows, err := stmt.QueryContext(ctx, filters.Pagination.Limit, filters.Pagination.Offset, labelID)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetArtistsLabelID").Inc()
		logger.Error("failed to get all artists", zap.Error(err))
		return nil, artistErrors.NewInternalError("failed to get all artists: %v", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Error("failed to close rows", zap.Error(err))
		}
	}()

	artists := make([]*repoModel.Artist, 0)
	for rows.Next() {
		var artist repoModel.Artist
		var isFavorite sql.NullBool
		err = rows.Scan(&artist.ID, &artist.Title, &artist.Description, &artist.Thumbnail, &isFavorite)
		if err != nil {
			r.metrics.DatabaseErrors.WithLabelValues("GetArtistsLabelID").Inc()
			logger.Error("failed to scan artist", zap.Error(err))
			return nil, artistErrors.NewInternalError("failed to scan artist: %v", err)
		}
		artist.IsFavorite = isFavorite.Valid && isFavorite.Bool
		artists = append(artists, &artist)
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("GetArtistsLabelID").Observe(duration)
	return artists, nil
}

func (r *artistPostgresRepository) DeleteArtist(ctx context.Context, artistID int64) error {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("deleting user", zap.Int64("artist id", artistID))

	stmt, err := r.db.PrepareContext(ctx, DeleteArtistQuery)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("DeleteArtist").Inc()
		logger.Error("failed to delete artist", zap.Error(err))
		return err
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			logger.Error("failed to close statement", zap.Error(err))
		}
	}()

	_, err = stmt.ExecContext(ctx, artistID)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("DeleteArtist").Inc()
		logger.Error("failed to delete artist", zap.Error(err))
		return err
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("DeleteArtist").Observe(duration)
	return nil
}

func (r *artistPostgresRepository) AddArtistsToAlbum(ctx context.Context, artistsIDs []int64, albumID int64) error {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Adding artists to album", zap.Any("artistsIds", artistsIDs), zap.Int64("albumID", albumID))

	if len(artistsIDs) == 0 {
		return nil
	}

	stmt, err := r.db.PrepareContext(ctx, AddArtistsToAlbum)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("AddArtistsToAlbum").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return artistErrors.NewInternalError("failed to prepare statement: %v", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			logger.Error("failed to close statement", zap.Error(err))
		}
	}()

	_, err = stmt.ExecContext(ctx, pq.Array(artistsIDs), albumID)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("AddArtistsToAlbum").Inc()
		logger.Error("failed to add artists to album", zap.Error(err))
		return artistErrors.NewInternalError("failed to add artists to album: %v", err)
	}

	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("AddArtistsToAlbum").Observe(duration)
	return nil
}

func (r *artistPostgresRepository) AddArtistsToTracks(ctx context.Context, artistsIDs []int64, trackIDs []int64) error {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Adding artists to tracks", zap.Any("artistsIds", artistsIDs), zap.Any("trackIDs", trackIDs))

	if len(artistsIDs) == 0 || len(trackIDs) == 0 {
		return nil
	}

	stmt, err := r.db.PrepareContext(ctx, AddArtistsToTracks)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("AddArtistsToTracks").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return artistErrors.NewInternalError("failed to prepare statement: %v", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			logger.Error("failed to close statement", zap.Error(err))
		}
	}()

	_, err = stmt.ExecContext(ctx, pq.Array(artistsIDs), pq.Array(trackIDs))
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("AddArtistsToTracks").Inc()
		logger.Error("failed to add artists to tracks", zap.Error(err))
		return artistErrors.NewInternalError("failed to add artists to tracks: %v", err)
	}

	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("AddArtistsToTracks").Observe(duration)
	return nil
}
