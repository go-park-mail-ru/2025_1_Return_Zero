package repository

import (
	"context"
	"database/sql"
	"strings"

	loggerPkg "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
	domain "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/playlist/internal/domain"
	playlistErrors "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/playlist/model/errors"
	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/playlist/model/repository"

	"go.uber.org/zap"
)

const (
	CreatePlaylistQuery = `
		INSERT INTO playlist (title, user_id, thumbnail_url, is_public)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	GetPlaylistByIDQuery = `
		SELECT id, title, user_id, thumbnail_url, is_public
		FROM playlist
		WHERE id = $1
	`

	// Owned and favorite playlists
	GetPlaylistsByUserIDQuery = `
		SELECT DISTINCT p.id, p.title, p.user_id, p.thumbnail_url,
			CASE WHEN p.user_id = $1 THEN p.created_at ELSE fp.created_at END as effective_created_at
		FROM playlist p
		LEFT JOIN favorite_playlist fp ON p.id = fp.playlist_id AND fp.user_id = $1
		WHERE p.user_id = $1 OR (fp.user_id = $1 AND p.is_public = true)
		ORDER BY
			effective_created_at DESC
	`

	AddTrackToPlaylistQuery = `
		INSERT INTO playlist_track (playlist_id, track_id)
		VALUES ($1, $2)
	`

	RemoveTrackFromPlaylistQuery = `
		DELETE FROM playlist_track
		WHERE playlist_id = $1 AND track_id = $2
	`

	TrackExistsInPlaylistQuery = `
		SELECT EXISTS (
			SELECT 1
			FROM playlist_track
			WHERE playlist_id = $1 AND track_id = $2
		)
	`

	GetPlaylistTrackIdsQuery = `
		SELECT track_id
		FROM playlist_track
		WHERE playlist_id = $1
		ORDER BY created_at ASC, id ASC
	`

	UpdatePlaylistWithThumbnailQuery = `
		UPDATE playlist
		SET title = $2, thumbnail_url = $3
		WHERE id = $1 AND user_id = $4
		RETURNING id
	`

	UpdatePlaylistWithoutThumbnailQuery = `
		UPDATE playlist
		SET title = $2
		WHERE id = $1 AND user_id = $3
		RETURNING id
	`

	RemovePlaylistQuery = `
		DELETE FROM playlist
		WHERE id = $1 AND user_id = $2
	`

	GetPlaylistsToAddQuery = `
		SELECT p.id, p.title, p.user_id, p.thumbnail_url, 
		       EXISTS (
		           SELECT 1 
		           FROM playlist_track pt 
		           WHERE pt.playlist_id = p.id AND pt.track_id = $1
		       ) as is_included
		FROM playlist p
		WHERE p.user_id = $2
		ORDER BY p.created_at DESC
	`

	UpdatePlaylistsPublisityByUserIDQuery = `
		UPDATE playlist
		SET is_public = $2
		WHERE user_id = $1
	`

	CheckExistsPlaylistAndNotDifferentUserQuery = `
		SELECT EXISTS (
			SELECT 1
			FROM playlist
			WHERE id = $1 AND user_id != $2
		)
	`

	LikePlaylistQuery = `
		INSERT INTO favorite_playlist (user_id, playlist_id)
		VALUES ($1, $2) ON CONFLICT DO NOTHING
	`

	UnlikePlaylistQuery = `
		DELETE FROM favorite_playlist
		WHERE user_id = $1 AND playlist_id = $2
	`

	GetPlaylistWithIsLikedByIDQuery = `
		SELECT p.id, p.title, p.user_id, p.thumbnail_url, (fp.user_id IS NOT NULL) as is_liked
		FROM playlist p
		LEFT JOIN favorite_playlist fp ON p.id = fp.playlist_id AND fp.user_id = $2
		WHERE p.id = $1
	`

	GetProfilePlaylistsQuery = `
		SELECT p.id, p.title, p.user_id, p.thumbnail_url
		FROM playlist p
		WHERE p.user_id = $1
		ORDER BY p.created_at DESC
	`

	SearchPlaylistsQuery = `
		SELECT id, title, user_id, thumbnail_url
		FROM playlist
		WHERE (is_public = true OR user_id = $2) AND (search_vector @@ to_tsquery('multilingual', $1)
		   OR similarity(title_trgm, $3) > 0.3)
		ORDER BY 
		    CASE WHEN search_vector @@ to_tsquery('multilingual', $1) THEN 0 ELSE 1 END,
		    ts_rank(search_vector, to_tsquery('multilingual', $1)) DESC,
		    similarity(title_trgm, $3) DESC
	`
)

type PlaylistPostgresRepository struct {
	db *sql.DB
}

func NewPlaylistPostgresRepository(db *sql.DB) domain.Repository {
	return &PlaylistPostgresRepository{db: db}
}

func (r *PlaylistPostgresRepository) GetPlaylistByID(ctx context.Context, id int64) (*repoModel.Playlist, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Getting playlist by id", zap.Int64("id", id))

	var playlist repoModel.Playlist
	err := r.db.QueryRowContext(ctx, GetPlaylistByIDQuery, id).Scan(&playlist.ID, &playlist.Title, &playlist.UserID, &playlist.Thumbnail, &playlist.IsPublic)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Warn("Playlist not found", zap.Int64("id", id))
			return nil, playlistErrors.ErrPlaylistNotFound
		}
		logger.Error("Failed to get playlist by id", zap.Error(err))
		return nil, playlistErrors.NewInternalError("failed to get playlist by id: %v", err)
	}

	return &playlist, nil
}

func (r *PlaylistPostgresRepository) CreatePlaylist(ctx context.Context, playlistCreateRequest *repoModel.CreatePlaylistRequest) (*repoModel.Playlist, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Creating playlist", zap.Any("playlist", playlistCreateRequest))

	var id int64
	err := r.db.QueryRowContext(ctx, CreatePlaylistQuery, playlistCreateRequest.Title, playlistCreateRequest.UserID, playlistCreateRequest.Thumbnail, playlistCreateRequest.IsPublic).Scan(&id)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			logger.Warn("Failed to create playlist: duplicate title for user", zap.Error(err))
			return nil, playlistErrors.ErrPlaylistDuplicate
		}
		logger.Error("Failed to create playlist", zap.Error(err))
		return nil, playlistErrors.NewInternalError("failed to create playlist: %v", err)
	}

	playlist, err := r.GetPlaylistByID(ctx, id)
	if err != nil {
		logger.Error("Failed to get playlist by id", zap.Error(err))
		return nil, playlistErrors.NewInternalError("failed to get playlist by id: %v", err)
	}

	return playlist, nil
}

func (r *PlaylistPostgresRepository) GetCombinedPlaylistsByUserID(ctx context.Context, userID int64) (*repoModel.PlaylistList, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Getting playlists by user id", zap.Int64("user_id", userID))

	var playlists repoModel.PlaylistList
	rows, err := r.db.QueryContext(ctx, GetPlaylistsByUserIDQuery, userID)
	if err != nil {
		logger.Error("Failed to get playlists by user id", zap.Error(err))
		return nil, playlistErrors.NewInternalError("failed to get playlists by user id: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var playlist repoModel.Playlist
		err := rows.Scan(&playlist.ID, &playlist.Title, &playlist.UserID, &playlist.Thumbnail)
		if err != nil {
			logger.Error("Failed to scan playlist", zap.Error(err))
			return nil, playlistErrors.NewInternalError("failed to scan playlist: %v", err)
		}
		playlists.Playlists = append(playlists.Playlists, &playlist)
	}

	return &playlists, nil
}

func (r *PlaylistPostgresRepository) TrackExistsInPlaylist(ctx context.Context, playlistID int64, trackID int64) (bool, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Checking if track exists in playlist", "playlist_id", playlistID, "track_id", trackID)

	var exists bool
	err := r.db.QueryRowContext(ctx, TrackExistsInPlaylistQuery, playlistID, trackID).Scan(&exists)
	if err != nil {
		logger.Error("Failed to check if track exists in playlist", zap.Error(err))
		return false, playlistErrors.NewInternalError("failed to check if track exists in playlist: %v", err)
	}

	return exists, nil
}

func (r *PlaylistPostgresRepository) AddTrackToPlaylist(ctx context.Context, request *repoModel.AddTrackToPlaylistRequest) error {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Adding track to playlist", zap.Int64("playlist_id", request.PlaylistID), zap.Int64("track_id", request.TrackID))

	playlist, err := r.GetPlaylistByID(ctx, request.PlaylistID)
	if err != nil {
		logger.Error("Failed to get playlist by id", zap.Error(err))
		return err
	}

	if playlist.UserID != request.UserID {
		logger.Warn("User tryed to add track to another user's playlist", zap.Int64("playlist_id", request.PlaylistID), zap.Int64("user_id", request.UserID))
		return playlistErrors.ErrPlaylistPermissionDenied
	}

	trackExists, err := r.TrackExistsInPlaylist(ctx, request.PlaylistID, request.TrackID)
	if err != nil {
		logger.Error("Failed to check if track exists in playlist", zap.Error(err))
		return playlistErrors.NewInternalError("failed to check if track exists in playlist: %v", err)
	}

	if trackExists {
		logger.Warn("Track already in playlist", zap.Int64("playlist_id", request.PlaylistID), zap.Int64("track_id", request.TrackID))
		return playlistErrors.ErrPlaylistTrackDuplicate
	}

	_, err = r.db.ExecContext(ctx, AddTrackToPlaylistQuery, request.PlaylistID, request.TrackID)
	if err != nil {
		logger.Error("Failed to add track to playlist", zap.Error(err))
		return playlistErrors.NewInternalError("failed to add track to playlist: %v", err)
	}

	return nil
}

func (r *PlaylistPostgresRepository) RemoveTrackFromPlaylist(ctx context.Context, request *repoModel.RemoveTrackFromPlaylistRequest) error {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Removing track from playlist", zap.Int64("playlist_id", request.PlaylistID), zap.Int64("track_id", request.TrackID))

	playlist, err := r.GetPlaylistByID(ctx, request.PlaylistID)
	if err != nil {
		logger.Error("Failed to get playlist by id", zap.Error(err))
		return err
	}

	if playlist.UserID != request.UserID {
		logger.Warn("User tryed to remove track from another user's playlist", zap.Int64("playlist_id", request.PlaylistID), zap.Int64("user_id", request.UserID))
		return playlistErrors.ErrPlaylistPermissionDenied
	}

	trackExists, err := r.TrackExistsInPlaylist(ctx, request.PlaylistID, request.TrackID)
	if err != nil {
		logger.Error("Failed to check if track exists in playlist", zap.Error(err))
		return playlistErrors.NewInternalError("failed to check if track exists in playlist: %v", err)
	}

	if !trackExists {
		logger.Warn("Track does not exist in playlist", zap.Int64("playlist_id", request.PlaylistID), zap.Int64("track_id", request.TrackID))
		return playlistErrors.ErrPlaylistTrackNotFound
	}

	_, err = r.db.ExecContext(ctx, RemoveTrackFromPlaylistQuery, request.PlaylistID, request.TrackID)
	if err != nil {
		logger.Error("Failed to remove track from playlist", zap.Error(err))
		return playlistErrors.NewInternalError("failed to remove track from playlist: %v", err)
	}

	return nil
}

func (r *PlaylistPostgresRepository) GetPlaylistTrackIds(ctx context.Context, request *repoModel.GetPlaylistTrackIdsRequest) ([]int64, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Getting playlist track ids", "playlist_id", request.PlaylistID)

	var trackIds []int64
	rows, err := r.db.QueryContext(ctx, GetPlaylistTrackIdsQuery, request.PlaylistID)
	if err != nil {
		logger.Error("Failed to get playlist track ids", zap.Error(err))
		return nil, playlistErrors.NewInternalError("failed to get playlist track ids: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var trackId int64
		err := rows.Scan(&trackId)
		if err != nil {
			logger.Error("Failed to scan playlist track id", zap.Error(err))
			return nil, playlistErrors.NewInternalError("failed to scan playlist track id: %v", err)
		}
		trackIds = append(trackIds, trackId)
	}

	if err := rows.Err(); err != nil {
		logger.Error("Failed to iterate over playlist track ids", zap.Error(err))
		return nil, playlistErrors.NewInternalError("failed to iterate over playlist track ids: %v", err)
	}

	logger.Info("Playlist track ids", zap.Any("track_ids", trackIds))

	return trackIds, nil
}

func (r *PlaylistPostgresRepository) UpdatePlaylist(ctx context.Context, request *repoModel.UpdatePlaylistRequest) (*repoModel.Playlist, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Updating playlist", zap.Any("playlist", request))

	var id int64
	if request.Thumbnail != "" {
		err := r.db.QueryRowContext(ctx, UpdatePlaylistWithThumbnailQuery, request.PlaylistID, request.Title, request.Thumbnail, request.UserID).Scan(&id)
		if err != nil {
			logger.Error("Failed to update playlist", zap.Error(err))
			return nil, playlistErrors.NewInternalError("failed to update playlist: %v", err)
		}
	} else {
		err := r.db.QueryRowContext(ctx, UpdatePlaylistWithoutThumbnailQuery, request.PlaylistID, request.Title, request.UserID).Scan(&id)
		if err != nil {
			logger.Error("Failed to update playlist", zap.Error(err))
			return nil, playlistErrors.NewInternalError("failed to update playlist: %v", err)
		}
	}

	playlist, err := r.GetPlaylistByID(ctx, id)
	if err != nil {
		logger.Error("Failed to get playlist by id", zap.Error(err))
		return nil, playlistErrors.NewInternalError("failed to get playlist by id: %v", err)
	}

	return playlist, nil
}

func (r *PlaylistPostgresRepository) RemovePlaylist(ctx context.Context, request *repoModel.RemovePlaylistRequest) error {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Removing playlist", zap.Int64("playlist_id", request.PlaylistID))

	_, err := r.db.ExecContext(ctx, RemovePlaylistQuery, request.PlaylistID, request.UserID)
	if err != nil {
		logger.Error("Failed to remove playlist", zap.Error(err))
		return playlistErrors.NewInternalError("failed to remove playlist: %v", err)
	}

	return nil
}

func (r *PlaylistPostgresRepository) GetPlaylistsToAdd(ctx context.Context, request *repoModel.GetPlaylistsToAddRequest) (*repoModel.GetPlaylistsToAddResponse, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Getting playlists to add track to", zap.Int64("track_id", request.TrackID), zap.Int64("user_id", request.UserID))

	var response repoModel.GetPlaylistsToAddResponse
	rows, err := r.db.QueryContext(ctx, GetPlaylistsToAddQuery, request.TrackID, request.UserID)
	if err != nil {
		logger.Error("Failed to get playlists to add track to", zap.Error(err))
		return nil, playlistErrors.NewInternalError("failed to get playlists to add track to: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var playlist repoModel.Playlist
		var isIncluded sql.NullBool
		err := rows.Scan(&playlist.ID, &playlist.Title, &playlist.UserID, &playlist.Thumbnail, &isIncluded)
		if err != nil {
			logger.Error("Failed to scan playlist", zap.Error(err))
			return nil, playlistErrors.NewInternalError("failed to scan playlist: %v", err)
		}

		playlistWithFlag := &repoModel.PlaylistWithIsIncludedTrack{
			Playlist:   &playlist,
			IsIncluded: isIncluded.Valid && isIncluded.Bool,
		}

		response.Playlists = append(response.Playlists, playlistWithFlag)
	}

	if err := rows.Err(); err != nil {
		logger.Error("Failed to iterate over playlists", zap.Error(err))
		return nil, playlistErrors.NewInternalError("failed to iterate over playlists: %v", err)
	}

	return &response, nil
}

func (r *PlaylistPostgresRepository) UpdatePlaylistsPublisityByUserID(ctx context.Context, request *repoModel.UpdatePlaylistsPublisityByUserIDRequest) error {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Updating playlists publisity by user id", zap.Int64("user_id", request.UserID), zap.Bool("is_public", request.IsPublic))

	_, err := r.db.ExecContext(ctx, UpdatePlaylistsPublisityByUserIDQuery, request.UserID, request.IsPublic)
	if err != nil {
		logger.Error("Failed to update playlists publisity by user id", zap.Error(err))
		return playlistErrors.NewInternalError("failed to update playlists publisity by user id: %v", err)
	}

	return nil
}

func (r *PlaylistPostgresRepository) CheckExistsPlaylistAndNotDifferentUser(ctx context.Context, playlistID int64, userID int64) (bool, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Checking if playlist exists and is not different user", zap.Int64("playlist_id", playlistID), zap.Int64("user_id", userID))

	var exists bool
	err := r.db.QueryRowContext(ctx, CheckExistsPlaylistAndNotDifferentUserQuery, playlistID, userID).Scan(&exists)
	if err != nil {
		logger.Error("Failed to check if playlist exists and is not different user", zap.Error(err))
		return false, playlistErrors.NewInternalError("failed to check if playlist exists and is not different user: %v", err)
	}

	return exists, nil
}

func (r *PlaylistPostgresRepository) LikePlaylist(ctx context.Context, request *repoModel.LikePlaylistRequest) error {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Liking playlist", zap.Int64("playlist_id", request.PlaylistID), zap.Int64("user_id", request.UserID))

	exists, err := r.CheckExistsPlaylistAndNotDifferentUser(ctx, request.PlaylistID, request.UserID)
	if err != nil {
		logger.Error("Failed to check if playlist exists and is not different user", zap.Error(err))
		return playlistErrors.NewInternalError("failed to check if playlist exists and is not different user: %v", err)
	}

	if !exists {
		logger.Warn("Playlist does not exist or is different user", zap.Int64("playlist_id", request.PlaylistID), zap.Int64("user_id", request.UserID))
		return playlistErrors.ErrPlaylistNotFound
	}

	_, err = r.db.ExecContext(ctx, LikePlaylistQuery, request.UserID, request.PlaylistID)
	if err != nil {
		logger.Error("Failed to like playlist", zap.Error(err))
		return playlistErrors.NewInternalError("failed to like playlist: %v", err)
	}

	return nil
}

func (r *PlaylistPostgresRepository) UnlikePlaylist(ctx context.Context, request *repoModel.LikePlaylistRequest) error {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Unliking playlist", zap.Int64("playlist_id", request.PlaylistID), zap.Int64("user_id", request.UserID))

	_, err := r.db.ExecContext(ctx, UnlikePlaylistQuery, request.UserID, request.PlaylistID)
	if err != nil {
		logger.Error("Failed to unlike playlist", zap.Error(err))
		return playlistErrors.NewInternalError("failed to unlike playlist: %v", err)
	}

	return nil
}

func (r *PlaylistPostgresRepository) GetPlaylistWithIsLikedByID(ctx context.Context, id int64, userID int64) (*repoModel.PlaylistWithIsLiked, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Getting playlist with is liked by id", zap.Int64("playlist_id", id), zap.Int64("user_id", userID))

	var playlist repoModel.Playlist
	var isLiked sql.NullBool
	err := r.db.QueryRowContext(ctx, GetPlaylistWithIsLikedByIDQuery, id, userID).Scan(&playlist.ID, &playlist.Title, &playlist.UserID, &playlist.Thumbnail, &isLiked)
	if err != nil {
		logger.Error("Failed to get playlist with is liked by id", zap.Error(err))
		return nil, playlistErrors.NewInternalError("failed to get playlist with is liked by id: %v", err)
	}

	return &repoModel.PlaylistWithIsLiked{
		Playlist: &playlist,
		IsLiked:  isLiked.Valid && isLiked.Bool,
	}, nil
}

func (r *PlaylistPostgresRepository) GetProfilePlaylists(ctx context.Context, request *repoModel.GetProfilePlaylistsRequest) (*repoModel.GetProfilePlaylistsResponse, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Getting profile playlists", zap.Int64("user_id", request.UserID))

	var playlists repoModel.GetProfilePlaylistsResponse
	rows, err := r.db.QueryContext(ctx, GetProfilePlaylistsQuery, request.UserID)
	if err != nil {
		logger.Error("Failed to get profile playlists", zap.Error(err))
		return nil, playlistErrors.NewInternalError("failed to get profile playlists: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var playlist repoModel.Playlist
		err := rows.Scan(&playlist.ID, &playlist.Title, &playlist.UserID, &playlist.Thumbnail)
		if err != nil {
			logger.Error("Failed to scan playlist", zap.Error(err))
			return nil, playlistErrors.NewInternalError("failed to scan playlist: %v", err)
		}

		playlists.Playlists = append(playlists.Playlists, &playlist)
	}

	if err := rows.Err(); err != nil {
		logger.Error("Failed to iterate over playlists", zap.Error(err))
		return nil, playlistErrors.NewInternalError("failed to iterate over playlists: %v", err)
	}

	return &playlists, nil
}

func (r *PlaylistPostgresRepository) SearchPlaylists(ctx context.Context, request *repoModel.SearchPlaylistsRequest) (*repoModel.PlaylistList, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Searching playlists", zap.String("query", request.Query))

	words := strings.Fields(request.Query)
	for i, word := range words {
		words[i] = word + ":*"
	}
	tsQueryString := strings.Join(words, " & ")

	var playlists repoModel.PlaylistList
	rows, err := r.db.QueryContext(ctx, SearchPlaylistsQuery, tsQueryString, request.UserID, request.Query)
	if err != nil {
		logger.Error("Failed to search playlists", zap.Error(err))
		return nil, playlistErrors.NewInternalError("failed to search playlists: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var playlist repoModel.Playlist
		err := rows.Scan(&playlist.ID, &playlist.Title, &playlist.UserID, &playlist.Thumbnail)
		if err != nil {
			logger.Error("Failed to scan playlist", zap.Error(err))
			return nil, playlistErrors.NewInternalError("failed to scan playlist: %v", err)
		}

		playlists.Playlists = append(playlists.Playlists, &playlist)
	}

	if err := rows.Err(); err != nil {
		logger.Error("Failed to iterate over playlists", zap.Error(err))
		return nil, playlistErrors.NewInternalError("failed to iterate over playlists: %v", err)
	}

	return &playlists, nil
}
