package repository

import (
	"bytes"
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"strings"

	"golang.org/x/crypto/argon2"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers"
	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/repository"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/user"
	"go.uber.org/zap"
)

type userPostgresRepository struct {
	db *sql.DB
}

const (
	checkUserExist = `
			SELECT 1 
			FROM "user"
			WHERE username = $1 OR email = $2
			`
	createUserQuery = `
			INSERT INTO "user" (username, password_hash, email) 
			VALUES ($1, $2, $3) 
			RETURNING id
			`
	getUserByIDQuery = `
			SELECT id, username, email, thumbnail_url
			FROM "user"
			WHERE id = $1
			`
	loginUserQuery = `
			SELECT id, username, email, password_hash, thumbnail_url
			FROM "user"
			WHERE username = $1 OR email = $2
			`
	uploadAvatarQuery = `
			UPDATE "user"
			SET thumbnail_url = $1
			WHERE username = $2
			`
	getAvatarQuery = `
			SELECT thumbnail_url
			FROM "user"
			WHERE username = $1
			`
	changeUsernameQuery = `
			UPDATE "user"
			SET username = $1
			WHERE id = $2
			`
	changeEmailQuery = `
			UPDATE "user"
			SET email = $1
			WHERE id = $2
			`
	changePasswordQuery = `
			UPDATE "user"
			SET password_hash = $1
			WHERE id = $2
			`
	getPasswordQuery = `
			SELECT password_hash
			FROM "user"
			WHERE id = $1
	`
	deleteUserQuery = `
			DELETE FROM "user"
			WHERE username = $1 AND email = $2
	`
	changePrivacySettingsQuery = `
			UPDATE "user_settings"
			SET is_public_playlists = $1,
				is_public_minutes_listened = $2,
				is_public_favorite_artists = $3,
				is_public_tracks_listened = $4,
				is_public_favorite_tracks = $5,
				is_public_artists_listened = $6
			WHERE user_id = $7
	`
	getIdByUsernameQuery = `
			SELECT id
			FROM "user"
			WHERE username = $1
	`
	getUserDataQuery = `
			SELECT username, email, thumbnail_url
			FROM "user"
			WHERE id = $1
	`
	createUserSettingsQuery = `
            INSERT INTO "user_settings" (
                user_id, 
                is_public_playlists, 
                is_public_minutes_listened, 
                is_public_favorite_artists, 
                is_public_tracks_listened, 
                is_public_favorite_tracks, 
                is_public_artists_listened
            ) VALUES ($1, false, false, false, false, false, false)
    `
	getNumUniqueTracksQuery = `
			SELECT COUNT(DISTINCT track_id) AS num_unique_tracks
			FROM stream
			WHERE user_id = $1
	`
	getMinutesListenedQuery = `
			SELECT COALESCE(SUM(duration) / 60, 0) AS total_minutes
			FROM stream
			WHERE user_id = $1
	`
	getNumUniqueArtistQuery = `
			SELECT COUNT(DISTINCT ta.artist_id) AS unique_artists_listened
			FROM stream s
			JOIN track_artist ta ON s.track_id = ta.track_id
			WHERE s.user_id = $1;
	`
	getUserPrivacySettingsQuery = `
			SELECT is_public_playlists, is_public_minutes_listened, is_public_favorite_artists,
				is_public_tracks_listened, is_public_favorite_tracks, is_public_artists_listened
			FROM user_settings
			WHERE user_id = $1
	`
	checkIsUsernameUniqueQuery = `
			SELECT 1 
			FROM "user"
			WHERE username = $1 AND id != $2
			`
	checkIsEmailUniqueQuery = `
			SELECT 1 
			FROM "user"
			WHERE username = $1 AND id != $2
			`
)

func hashPassword(salt []byte, password string) string {
	hashedPass := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	combined := append(salt, hashedPass...)
	return base64.StdEncoding.EncodeToString(combined)
}

func checkPasswordHash(encodedHash string, password string) bool {
	decodedHash, err := base64.StdEncoding.DecodeString(encodedHash)
	if err != nil {
		return false
	}
	salt := decodedHash[:8]
	userPassHash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	return bytes.Equal(userPassHash, decodedHash[8:])
}

func (r *userPostgresRepository) getPassword(ctx context.Context, id int64) (string, error) {
	logger := helpers.LoggerFromContext(ctx)
	row := r.db.QueryRowContext(ctx, getPasswordQuery, id)
	var storedHash string
	err := row.Scan(&storedHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error("user not found", zap.Error(err))
			return "", user.ErrUserNotFound
		}
		logger.Error("failed to get password hash", zap.Error(err))
		return "", err
	}
	return storedHash, nil
}

func NewUserPostgresRepository(db *sql.DB) user.Repository {
	repo := &userPostgresRepository{
		db: db,
	}

	return repo
}

func createSalt() []byte {
	salt := make([]byte, 8)
	_, err := rand.Read(salt)
	if err != nil {
		return nil
	}
	return salt
}

func (r *userPostgresRepository) CreateUser(ctx context.Context, regData *repoModel.User) (*repoModel.User, error) {
	logger := helpers.LoggerFromContext(ctx)
	var exists bool
	lowerUsername := strings.ToLower(regData.Username)
	logger.Info("Creating new user", zap.String("username: ", lowerUsername))
	err := r.db.QueryRowContext(ctx, checkUserExist, lowerUsername, regData.Email).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		logger.Error("failed to check user existence", zap.Error(err))
		return nil, err
	}
	if exists {
		logger.Error("user with this username or email already exists")
		return nil, user.ErrUsernameExist
	}

	salt := createSalt()
	if salt == nil {
		logger.Error("failed to create salt")
		return nil, user.ErrCreateSalt
	}
	hashedPassword := hashPassword(salt, regData.Password)

	var userID int64
	err = r.db.QueryRowContext(ctx, createUserQuery, lowerUsername,
		hashedPassword, regData.Email).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error("user not found", zap.Error(err))
			return nil, user.ErrUserNotFound
		}
		logger.Error("failed to create user", zap.Error(err))
		return nil, err
	}

	_, err = r.db.ExecContext(ctx, createUserSettingsQuery, userID)
	if err != nil {
		logger.Error("failed to create user settings", zap.Error(err))
		return nil, err
	}

	return &repoModel.User{
		ID:        userID,
		Username:  lowerUsername,
		Email:     regData.Email,
		Thumbnail: "/default_avatar.png",
	}, nil
}

func (r *userPostgresRepository) GetUserByID(ctx context.Context, ID int64) (*repoModel.User, error) {
	logger := helpers.LoggerFromContext(ctx)
	logger.Info("Getting user by id", zap.Int64("ID", ID))
	row := r.db.QueryRowContext(ctx, getUserByIDQuery, ID)
	var userRepo repoModel.User
	err := row.Scan(&userRepo.ID, &userRepo.Username, &userRepo.Email, &userRepo.Thumbnail)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error("user not found", zap.Error(err))
			return nil, user.ErrUserNotFound
		}
		logger.Error("failed to get user by ID", zap.Error(err))
		return nil, err
	}
	return &userRepo, nil
}

func (r *userPostgresRepository) LoginUser(ctx context.Context, logData *repoModel.User) (*repoModel.User, error) {
	logger := helpers.LoggerFromContext(ctx)
	var storedHash string
	lowerUsername := strings.ToLower(logData.Username)
	logger.Info("Loggining user", zap.String("username", lowerUsername))
	row := r.db.QueryRowContext(ctx, loginUserQuery, lowerUsername, logData.Email)
	var userRepo repoModel.User
	err := row.Scan(&userRepo.ID, &userRepo.Username, &userRepo.Email, &storedHash, &userRepo.Thumbnail)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error("user not found", zap.Error(err))
			return nil, user.ErrUserNotFound
		}
		logger.Error("failed to get user by username or email", zap.Error(err))
		return nil, err
	}

	if !checkPasswordHash(storedHash, logData.Password) {
		logger.Error("wrong password", zap.Error(err))
		return nil, user.ErrWrongPassword
	}

	return &userRepo, nil
}

func (r *userPostgresRepository) GetAvatar(ctx context.Context, username string) (string, error) {
	logger := helpers.LoggerFromContext(ctx)
	logger.Info("Getting avarat by username", zap.String("username", username))
	row := r.db.QueryRowContext(ctx, getAvatarQuery, username)
	var avatarUrl string
	err := row.Scan(&avatarUrl)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error("user not found", zap.Error(err))
			return "", user.ErrUserNotFound
		}
		logger.Error("failed to get avatar", zap.Error(err))
		return "", err
	}
	return avatarUrl, nil
}

func (r *userPostgresRepository) UploadAvatar(ctx context.Context, avatarUrl string, username string) error {
	logger := helpers.LoggerFromContext(ctx)
	logger.Info("Loading avatar by username", zap.String("username", username))
	_, err := r.db.ExecContext(ctx, uploadAvatarQuery, avatarUrl, username)
	if err != nil {
		logger.Error("failed to upload avatar", zap.Error(err))
		return err
	}
	return nil
}

func (r *userPostgresRepository) changeUsername(ctx context.Context, id int64, newUsername string) error {
	logger := helpers.LoggerFromContext(ctx)
	newLowerUsername := strings.ToLower(newUsername)
	logger.Info("Changing username", zap.String("username", newLowerUsername))
	var isExist bool
	err := r.db.QueryRowContext(ctx, checkIsUsernameUniqueQuery, newLowerUsername, id).Scan(&isExist)
	if err != nil && err != sql.ErrNoRows {
		logger.Error("failed to check username uniqueness", zap.Error(err))
		return err
	}
	if isExist {
		logger.Error("username already occupied. username: ", zap.String("username", newLowerUsername))
		return user.ErrUsernameExist
	}
	_, err = r.db.ExecContext(ctx, changeUsernameQuery, newLowerUsername, id)
	if err != nil {
		logger.Error("failed to change username", zap.Error(err))
		return err
	}
	return nil
}

func (r *userPostgresRepository) changeEmail(ctx context.Context, id int64, newEmail string) error {
	logger := helpers.LoggerFromContext(ctx)
	logger.Info("Changing email", zap.String("email", newEmail))
	var isExist bool
	err := r.db.QueryRowContext(ctx, checkIsEmailUniqueQuery, newEmail, id).Scan(&isExist)
	if err != nil && err != sql.ErrNoRows {
		logger.Error("failed to check email uniqueness", zap.Error(err))
		return err
	}
	if isExist {
		logger.Error("email already occupied. email: ", zap.String("email", newEmail))
		return user.ErrUsernameExist
	}
	_, err = r.db.ExecContext(ctx, changeEmailQuery, newEmail, id)
	if err != nil {
		logger.Error("failed to change email", zap.Error(err))
		return err
	}
	return nil
}

func (r *userPostgresRepository) сhangePassword(ctx context.Context, password string, id int64, newPassword string) error {
	logger := helpers.LoggerFromContext(ctx)
	logger.Info("Changing password")
	storedHash, err := r.getPassword(ctx, id)
	if err != nil {
		logger.Error("failed to get password hash", zap.Error(err))
		return err
	}
	if !checkPasswordHash(storedHash, password) {
		logger.Error("wrong password", zap.Error(err))
		return user.ErrWrongPassword
	}
	salt := createSalt()
	newHashedPassword := hashPassword(salt, newPassword)
	_, err = r.db.ExecContext(ctx, changePasswordQuery, newHashedPassword, id)
	if err != nil {
		logger.Error("failed to change password", zap.Error(err))
		return err
	}
	return nil
}

func (r *userPostgresRepository) ChangeUserData(ctx context.Context, username string, changeData *repoModel.ChangeUserData) error {
	logger := helpers.LoggerFromContext(ctx)
	logger.Info("Changing data by username", zap.String("username", username))
	if changeData.NewPassword != "" && changeData.Password == "" {
		logger.Error("password is required to change password")
		return user.ErrPasswordRequired
	}
	id, err := r.GetIDByUsername(ctx, username)
	if err != nil {
		logger.Error("failed to get user ID", zap.Error(err))
		return err
	}
	if changeData.NewUsername != "" {
		err := r.changeUsername(ctx, id, changeData.NewUsername)
		if err != nil {
			logger.Error("failed to change username", zap.Error(err))
			return err
		}
	}
	if changeData.NewEmail != "" {
		err := r.changeEmail(ctx, id, changeData.NewEmail)
		if err != nil {
			logger.Error("failed to change email", zap.Error(err))
			return err
		}
	}
	if changeData.NewPassword != "" {
		err := r.сhangePassword(ctx, changeData.Password, id, changeData.NewPassword)
		if err != nil {
			logger.Error("failed to change password", zap.Error(err))
			return err
		}
	}
	return nil
}

func (r *userPostgresRepository) DeleteUser(ctx context.Context, userRepo *repoModel.User) error {
	logger := helpers.LoggerFromContext(ctx)
	logger.Info("Deleting user", zap.String("username", userRepo.Username))
	id, err := r.GetIDByUsername(ctx, userRepo.Username)
	if err != nil {
		logger.Error("failed to find user", zap.Error(err))
		return err
	}
	storedHash, err := r.getPassword(ctx, id)
	if err != nil {
		logger.Error("failed to get password hash", zap.Error(err))
		return err
	}
	if !checkPasswordHash(storedHash, userRepo.Password) {
		logger.Error("wrong password", zap.Error(err))
		return user.ErrWrongPassword
	}
	_, err = r.db.ExecContext(ctx, deleteUserQuery, userRepo.Username, userRepo.Email)
	if err != nil {
		logger.Error("failed to delete user", zap.Error(err))
		return err
	}
	return nil
}

func (r *userPostgresRepository) ChangeUserPrivacySettings(ctx context.Context, username string, privacySettings *repoModel.UserPrivacySettings) error {
	logger := helpers.LoggerFromContext(ctx)
	logger.Info("Changing user privacy", zap.String("username", username))
	id, err := r.GetIDByUsername(ctx, username)
	if err != nil {
		logger.Error("failed to get user ID", zap.Error(err))
		return err
	}
	_, err = r.db.ExecContext(ctx, changePrivacySettingsQuery,
		privacySettings.IsPublicPlaylists,
		privacySettings.IsPublicMinutesListened,
		privacySettings.IsPublicFavoriteArtists,
		privacySettings.IsPublicTracksListened,
		privacySettings.IsPublicFavoriteTracks,
		privacySettings.IsPublicArtistsListened,
		id,
	)
	if err != nil {
		logger.Error("failed to change privacy settings", zap.Error(err))
		return err
	}
	return nil
}

func (r *userPostgresRepository) GetIDByUsername(ctx context.Context, username string) (int64, error) {
	logger := helpers.LoggerFromContext(ctx)
	logger.Info("Getting ID by username", zap.String("username", username))
	row := r.db.QueryRowContext(ctx, getIdByUsernameQuery, username)
	var userID int64
	err := row.Scan(&userID)
	if err != nil {
		logger.Error("failed to get user ID", zap.Error(err))
		if errors.Is(err, sql.ErrNoRows) {
			return 0, user.ErrUserNotFound
		}
		logger.Error("user not found", zap.Error(err))
		return 0, err
	}
	return userID, nil
}

func (r *userPostgresRepository) GetUserData(ctx context.Context, id int64) (*repoModel.User, error) {
	logger := helpers.LoggerFromContext(ctx)
	logger.Info("Gettign user data by id", zap.Int64("ID", id))
	row := r.db.QueryRowContext(ctx, getUserDataQuery, id)
	var userRepo repoModel.User
	err := row.Scan(&userRepo.Username, &userRepo.Email, &userRepo.Thumbnail)
	if err != nil {
		logger.Error("failed to get user data", zap.Error(err))
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user.ErrUserNotFound
		}
		logger.Error("user not found", zap.Error(err))
		return nil, err
	}
	return &userRepo, nil
}

func (r *userPostgresRepository) getNumUniqueTracks(ctx context.Context, id int64) (int64, error) {
	logger := helpers.LoggerFromContext(ctx)
	logger.Info("Getting num unique tracks by id", zap.Int64("ID", id))
	row := r.db.QueryRowContext(ctx, getNumUniqueTracksQuery, id)
	var numUniqueTracks int64
	err := row.Scan(&numUniqueTracks)
	if err != nil {
		logger.Error("failed to get number of unique tracks", zap.Error(err))
		if errors.Is(err, sql.ErrNoRows) {
			return 0, user.ErrUserNotFound
		}
		logger.Error("user not found", zap.Error(err))
		return 0, err
	}
	return numUniqueTracks, nil
}

func (r *userPostgresRepository) getNumMinutes(ctx context.Context, id int64) (int64, error) {
	logger := helpers.LoggerFromContext(ctx)
	logger.Info("Getting num listened minutes by id", zap.Int64("ID", id))
	row := r.db.QueryRowContext(ctx, getMinutesListenedQuery, id)
	var numMinutes int64
	err := row.Scan(&numMinutes)
	if err != nil {
		logger.Error("failed to get number of minutes listened", zap.Error(err))
		if errors.Is(err, sql.ErrNoRows) {
			return -1, user.ErrUserNotFound
		}
		logger.Error("user not found", zap.Error(err))
		return -1, err
	}
	return numMinutes, nil
}

func (r *userPostgresRepository) getNumUniqueArtist(ctx context.Context, id int64) (int64, error) {
	logger := helpers.LoggerFromContext(ctx)
	logger.Info("Getting num unique artists by id", zap.Int64("ID", id))
	row := r.db.QueryRowContext(ctx, getNumUniqueArtistQuery, id)
	var numUniqueArtist int64
	err := row.Scan(&numUniqueArtist)
	if err != nil {
		logger.Error("failed to get number of unique artists", zap.Error(err))
		if errors.Is(err, sql.ErrNoRows) {
			return -1, user.ErrUserNotFound
		}
		logger.Error("user not found", zap.Error(err))
		return -1, err
	}
	return numUniqueArtist, nil

}

func (r *userPostgresRepository) GetUserStats(ctx context.Context, id int64) (*repoModel.UserStats, error) {
	logger := helpers.LoggerFromContext(ctx)
	logger.Info("Getting user statistics by id", zap.Int64("ID", id))
	numUniqueTracks, err := r.getNumUniqueTracks(ctx, id)
	if err != nil {
		logger.Error("failed to get number of unique tracks", zap.Error(err))
		return nil, err
	}
	numMinutes, err := r.getNumMinutes(ctx, id)
	if err != nil {
		logger.Error("failed to get number of minutes listened", zap.Error(err))
		return nil, err
	}
	numUniqueArtists, err := r.getNumUniqueArtist(ctx, id)
	if err != nil {
		logger.Error("failed to get number of unique artists", zap.Error(err))
		return nil, err
	}
	return &repoModel.UserStats{
		MinutesListened: numMinutes,
		TracksListened:  numUniqueTracks,
		ArtistsListened: numUniqueArtists,
	}, nil
}

func (r *userPostgresRepository) GetUserPrivacy(ctx context.Context, id int64) (*repoModel.UserPrivacySettings, error) {
	logger := helpers.LoggerFromContext(ctx)
	logger.Info("Getting user privacy settings by id", zap.Int64("ID", id))
	row := r.db.QueryRowContext(ctx, getUserPrivacySettingsQuery, id)
	var privacySettings repoModel.UserPrivacySettings
	err := row.Scan(&privacySettings.IsPublicPlaylists,
		&privacySettings.IsPublicMinutesListened,
		&privacySettings.IsPublicFavoriteArtists,
		&privacySettings.IsPublicTracksListened,
		&privacySettings.IsPublicFavoriteTracks,
		&privacySettings.IsPublicArtistsListened)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error("user not found", zap.Error(err))
			return nil, err
		}
		logger.Error("user not found", zap.Error(err))
		return nil, err
	}
	return &privacySettings, nil
}

func (r *userPostgresRepository) GetFullUserData(ctx context.Context, username string) (*repoModel.UserFullData, error) {
	logger := helpers.LoggerFromContext(ctx)
	logger.Info("Get full user data by username", zap.String("username", username))
	id, err := r.GetIDByUsername(ctx, username)
	if err != nil {
		logger.Error("failed to get user ID", zap.Error(err))
		return nil, err
	}
	privacy, err := r.GetUserPrivacy(ctx, id)
	if err != nil {
		logger.Error("failed to get user privacy settings", zap.Error(err))
		return nil, err
	}
	stats, err := r.GetUserStats(ctx, id)
	if err != nil {
		logger.Error("failed to get user statistics", zap.Error(err))
		return nil, err
	}
	user, err := r.GetUserData(ctx, id)
	if err != nil {
		logger.Error("failed to get user data", zap.Error(err))
		return nil, err
	}
	return &repoModel.UserFullData{
		Username:   user.Username,
		Email:      user.Email,
		Thumbnail:  user.Thumbnail,
		Privacy:    privacy,
		Statistics: stats,
	}, nil
}
