package repository

import (
	"bytes"
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/argon2"

	loggerPkg "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
	metrics "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/metrics"
	domain "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/user/internal/domain"
	userErrors "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/user/model/errors"
	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/user/model/repository"
	"go.uber.org/zap"
)

const (
	getPasswordQuery = `
			SELECT password_hash
			FROM "user"
			WHERE id = $1
	`
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
	loginUserQuery = `
			SELECT id, username, email, password_hash, thumbnail_url
			FROM "user"
			WHERE username = $1 OR email = $2
	`
	getUserByIDQuery = `
			SELECT id, username, email, thumbnail_url
			FROM "user"
			WHERE id = $1
	`
	uploadAvatarQuery = `
			UPDATE "user"
			SET thumbnail_url = $1
			WHERE id = $2	
	`
	getIdByUsernameQuery = `
			SELECT id
			FROM "user"
			WHERE username = $1
	`
	deleteUserQuery = `
			DELETE FROM "user"
			WHERE username = $1 AND email = $2
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
)

type userPostgresRepository struct {
	db      *sql.DB
	metrics *metrics.Metrics
}

func NewUserPostgresRepository(db *sql.DB, metrics *metrics.Metrics) domain.Repository {
	return &userPostgresRepository{db: db, metrics: metrics}
}

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
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)

	stmt, err := r.db.PrepareContext(ctx, getPasswordQuery)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetPassword").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return "", err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, id)
	var storedHash string
	err = row.Scan(&storedHash)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetPassword").Inc()
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error("user not found", zap.Error(err))
			return "", userErrors.NewNotFoundError("user not found")
		}
		logger.Error("failed to get password hash", zap.Error(err))
		return "", err
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("GetPassword").Observe(duration)
	return storedHash, nil
}

func createSalt() []byte {
	salt := make([]byte, 8)
	_, err := rand.Read(salt)
	if err != nil {
		return nil
	}
	return salt
}

func (r *userPostgresRepository) CheckUserExist(ctx context.Context, lowerUsername, email string) (bool, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Checking user existence", zap.String("username", lowerUsername), zap.String("email", email))
	stmt, err := r.db.PrepareContext(ctx, checkUserExist)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("CheckUserExist").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return false, err
	}

	defer stmt.Close()
	var exists bool
	err = stmt.QueryRowContext(ctx, lowerUsername, email).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	return exists, nil
}

func (r *userPostgresRepository) CreateUserSettings(ctx context.Context, userID int64) error {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Creating user settings", zap.Int64("userID", userID))

	stmt, err := r.db.PrepareContext(ctx, createUserSettingsQuery)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("CreateUserSettings").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, userID)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("CreateUserSettings").Inc()
		logger.Error("failed to create user settings", zap.Error(err))
		return err
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("CreateUserSettings").Observe(duration)
	return nil
}

func (r *userPostgresRepository) CreateUser(ctx context.Context, regData *repoModel.RegisterData) (*repoModel.User, error) {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Creating new user", zap.String("username: ", regData.Username))

	stmt, err := r.db.PrepareContext(ctx, createUserQuery)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("CreateUser").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return nil, err
	}
	defer stmt.Close()

	lowerUsername := strings.ToLower(regData.Username)
	logger.Info("Creating new user", zap.String("username: ", lowerUsername))

	exists, err := r.CheckUserExist(ctx, lowerUsername, regData.Email)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("CreateUser").Inc()
		logger.Error("failed to check user existence", zap.Error(err))
		return nil, err
	}

	if exists {
		r.metrics.DatabaseErrors.WithLabelValues("CreateUser").Inc()
		logger.Error("user with this username or email already exists")
		return nil, userErrors.NewUserExistError("user with this username or email already exists %s, %s", lowerUsername, regData.Email)
	}

	salt := createSalt()
	if salt == nil {
		r.metrics.DatabaseErrors.WithLabelValues("CreateUser").Inc()
		logger.Error("failed to create salt")
		return nil, userErrors.NewCreateSaltError("failed to create salt")
	}
	hashedPassword := hashPassword(salt, regData.Password)

	var userID int64
	err = stmt.QueryRowContext(ctx, lowerUsername,
		hashedPassword, regData.Email).Scan(&userID)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("CreateUser").Inc()
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error("user not found", zap.Error(err))
			return nil, userErrors.NewNotFoundError("user not found: %s", lowerUsername)
		}
		logger.Error("failed to create user", zap.Error(err))
		return nil, err
	}

	err = r.CreateUserSettings(ctx, userID)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("CreateUser").Inc()
		logger.Error("failed to create user settings", zap.Error(err))
		return nil, err
	}

	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("CreateUser").Observe(duration)
	return &repoModel.User{
		ID:        userID,
		Username:  lowerUsername,
		Email:     regData.Email,
		Thumbnail: "/default_avatar.png",
	}, nil
}

func (r *userPostgresRepository) LoginUser(ctx context.Context, logData *repoModel.LoginData) (*repoModel.User, error) {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)

	stmt, err := r.db.PrepareContext(ctx, loginUserQuery)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("LoginUser").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return nil, err
	}
	defer stmt.Close()

	var storedHash string
	lowerUsername := strings.ToLower(logData.Username)
	logger.Info("Loggining user", zap.String("username", lowerUsername))
	row := stmt.QueryRowContext(ctx, lowerUsername, logData.Email)
	var userRepo repoModel.User
	err = row.Scan(&userRepo.ID, &userRepo.Username, &userRepo.Email, &storedHash, &userRepo.Thumbnail)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("LoginUser").Inc()
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error("user not found", zap.Error(err))
			return nil, userErrors.NewNotFoundError("user not found: %s", lowerUsername)
		}
		logger.Error("failed to get user by username or email", zap.Error(err))
		return nil, err
	}

	if !checkPasswordHash(storedHash, logData.Password) {
		r.metrics.DatabaseErrors.WithLabelValues("LoginUser").Inc()
		logger.Error("wrong password", zap.Error(err))
		return nil, userErrors.NewWrongPasswordError("wrong password")
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("LoginUser").Observe(duration)
	return &userRepo, nil
}

func (r *userPostgresRepository) GetUserByID(ctx context.Context, ID int64) (*repoModel.User, error) {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Getting user by id", zap.Int64("ID", ID))

	stmt, err := r.db.PrepareContext(ctx, getUserByIDQuery)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetUserByID").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, ID)
	var userRepo repoModel.User
	err = row.Scan(&userRepo.ID, &userRepo.Username, &userRepo.Email, &userRepo.Thumbnail)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetUserByID").Inc()
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error("user not found", zap.Error(err))
			return nil, userErrors.NewNotFoundError("user not found")
		}
		logger.Error("failed to get user by ID", zap.Error(err))
		return nil, err
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("GetUserByID").Observe(duration)
	return &userRepo, nil
}

func (r *userPostgresRepository) UploadAvatar(ctx context.Context, avatarUrl string, ID int64) error {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Loading avatar by ID", zap.Int64("ID", ID))

	stmt, err := r.db.PrepareContext(ctx, uploadAvatarQuery)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("UploadAvatar").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, avatarUrl, ID)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("UploadAvatar").Inc()
		logger.Error("failed to upload avatar", zap.Error(err))
		return err
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("UploadAvatar").Observe(duration)
	return nil
}

func (r *userPostgresRepository) GetIDByUsername(ctx context.Context, username string) (int64, error) {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Getting ID by username", zap.String("username", username))

	stmt, err := r.db.PrepareContext(ctx, getIdByUsernameQuery)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetIDByUsername").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return 0, err
	}
	defer stmt.Close()

	lowerUsername := strings.ToLower(username)
	row := stmt.QueryRowContext(ctx, lowerUsername)
	var userID int64
	err = row.Scan(&userID)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetIDByUsername").Inc()
		logger.Error("failed to get user ID", zap.Error(err))
		if errors.Is(err, sql.ErrNoRows) {
			return 0, userErrors.NewNotFoundError("user not found: %s", lowerUsername)
		}
		logger.Error("user not found", zap.Error(err))
		return 0, err
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("GetIDByUsername").Observe(duration)
	return userID, nil
}

func (r *userPostgresRepository) DeleteUser(ctx context.Context, userRepo *repoModel.UserDelete) error {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Deleting user", zap.String("username", userRepo.Username))

	stmt, err := r.db.PrepareContext(ctx, deleteUserQuery)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("DeleteUser").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return err
	}
	defer stmt.Close()

	id, err := r.GetIDByUsername(ctx, userRepo.Username)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("DeleteUser").Inc()
		logger.Error("failed to find user", zap.Error(err))
		return err
	}
	storedHash, err := r.getPassword(ctx, id)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("DeleteUser").Inc()
		logger.Error("failed to get password hash", zap.Error(err))
		return err
	}
	if !checkPasswordHash(storedHash, userRepo.Password) {
		r.metrics.DatabaseErrors.WithLabelValues("DeleteUser").Inc()
		logger.Error("wrong password", zap.Error(err))
		return userErrors.NewWrongPasswordError("wrong password")
	}
	_, err = stmt.ExecContext(ctx, userRepo.Username, userRepo.Email)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("DeleteUser").Inc()
		logger.Error("failed to delete user", zap.Error(err))
		return err
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("DeleteUser").Observe(duration)
	return nil
}

func (r *userPostgresRepository) CheckIsUsernameUnique(ctx context.Context, username string) (bool, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Checking username uniqueness", zap.String("username", username))

	stmt, err := r.db.PrepareContext(ctx, checkIsUsernameUniqueQuery)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("CheckIsUsernameUnique").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return false, err
	}
	defer stmt.Close()

	var exists bool
	err = stmt.QueryRowContext(ctx, username).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	return exists, nil
}

func (r *userPostgresRepository) changeUsername(ctx context.Context, id int64, newUsername string) error {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	newLowerUsername := strings.ToLower(newUsername)
	logger.Info("Changing username", zap.String("username", newLowerUsername))

	stmt, err := r.db.PrepareContext(ctx, changeUsernameQuery)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("СhangeUsername").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return err
	}
	defer stmt.Close()

	isExist, err := r.CheckIsUsernameUnique(ctx, newLowerUsername)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("СhangeUsername").Inc()
		logger.Error("failed to check username uniqueness", zap.Error(err))
		return err
	}

	if isExist {
		r.metrics.DatabaseErrors.WithLabelValues("СhangeUsername").Inc()
		logger.Error("username already occupied. username: ", zap.String("username", newLowerUsername))
		return userErrors.NewUserExistError("username already occupied. username: %s", newLowerUsername)
	}

	_, err = stmt.ExecContext(ctx, newLowerUsername, id)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("СhangeUsername").Inc()
		logger.Error("failed to change username", zap.Error(err))
		return err
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("СhangeUsername").Observe(duration)
	return nil
}

func (r *userPostgresRepository) CheckIsEmailUnique(ctx context.Context, email string) (bool, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Checking email uniqueness", zap.String("email", email))

	stmt, err := r.db.PrepareContext(ctx, checkIsEmailUniqueQuery)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("CheckIsEmailUnique").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return false, err
	}
	defer stmt.Close()

	var exists bool
	err = stmt.QueryRowContext(ctx, email).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	return exists, nil
}

func (r *userPostgresRepository) changeEmail(ctx context.Context, id int64, newEmail string) error {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Changing email", zap.String("email", newEmail))

	stmt, err := r.db.PrepareContext(ctx, changeEmailQuery)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("changeEmail").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return err
	}
	defer stmt.Close()

	isExist, err := r.CheckIsEmailUnique(ctx, newEmail)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("changeEmail").Inc()
		logger.Error("failed to check email uniqueness", zap.Error(err))
		return err
	}

	if isExist {
		r.metrics.DatabaseErrors.WithLabelValues("changeEmail").Inc()
		logger.Error("email already occupied. email: ", zap.String("email", newEmail))
		return userErrors.NewUserExistError("email already occupied. email: %s", newEmail)
	}
	_, err = stmt.ExecContext(ctx, newEmail, id)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("changeEmail").Inc()
		logger.Error("failed to change email", zap.Error(err))
		return err
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("changeEmail").Observe(duration)
	return nil
}

func (r *userPostgresRepository) сhangePassword(ctx context.Context, password string, id int64, newPassword string) error {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Changing password")

	stmt, err := r.db.PrepareContext(ctx, changePasswordQuery)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("changePassword").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return err
	}
	defer stmt.Close()

	storedHash, err := r.getPassword(ctx, id)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("changePassword").Inc()
		logger.Error("failed to get password hash", zap.Error(err))
		return err
	}
	if !checkPasswordHash(storedHash, password) {
		r.metrics.DatabaseErrors.WithLabelValues("changePassword").Inc()
		logger.Error("wrong password", zap.Error(err))
		return userErrors.NewWrongPasswordError("wrong password")
	}
	salt := createSalt()
	newHashedPassword := hashPassword(salt, newPassword)
	_, err = stmt.ExecContext(ctx, newHashedPassword, id)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("changePassword").Inc()
		logger.Error("failed to change password", zap.Error(err))
		return err
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("changePassword").Observe(duration)
	return nil
}

func (r *userPostgresRepository) ChangeUserData(ctx context.Context, username string, changeData *repoModel.ChangeUserData) error {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Changing data by username", zap.String("username", username))
	if changeData.NewPassword != "" && changeData.Password == "" {
		r.metrics.DatabaseErrors.WithLabelValues("changeUserData").Inc()
		logger.Error("password is required to change password")
		return userErrors.NewPasswordRequierdError("password is required to change password")
	}
	id, err := r.GetIDByUsername(ctx, username)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("changeUserData").Inc()
		logger.Error("failed to get user ID", zap.Error(err))
		return err
	}
	if changeData.NewUsername != "" {
		err := r.changeUsername(ctx, id, changeData.NewUsername)
		if err != nil {
			r.metrics.DatabaseErrors.WithLabelValues("changeUserData").Inc()
			logger.Error("failed to change username", zap.Error(err))
			return err
		}
	}
	if changeData.NewEmail != "" {
		err := r.changeEmail(ctx, id, changeData.NewEmail)
		if err != nil {
			r.metrics.DatabaseErrors.WithLabelValues("changeUserData").Inc()
			logger.Error("failed to change email", zap.Error(err))
			return err
		}
	}
	if changeData.NewPassword != "" {
		err := r.сhangePassword(ctx, changeData.Password, id, changeData.NewPassword)
		if err != nil {
			r.metrics.DatabaseErrors.WithLabelValues("changeUserData").Inc()
			logger.Error("failed to change password", zap.Error(err))
			return err
		}
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("changeUserData").Observe(duration)
	return nil
}

func (r *userPostgresRepository) ChangeUserPrivacySettings(ctx context.Context, username string, privacySettings *repoModel.PrivacySettings) error {
	start := time.Now()

	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Changing user privacy", zap.String("username", username))

	stmt, err := r.db.PrepareContext(ctx, changePrivacySettingsQuery)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("changeUserPrivacySettings").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return err
	}
	defer stmt.Close()

	id, err := r.GetIDByUsername(ctx, username)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("changeUserPrivacySettings").Inc()
		logger.Error("failed to get user ID", zap.Error(err))
		return err
	}
	_, err = stmt.ExecContext(ctx,
		privacySettings.IsPublicPlaylists,
		privacySettings.IsPublicMinutesListened,
		privacySettings.IsPublicFavoriteArtists,
		privacySettings.IsPublicTracksListened,
		privacySettings.IsPublicFavoriteTracks,
		privacySettings.IsPublicArtistsListened,
		id,
	)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("changeUserPrivacySettings").Inc()
		logger.Error("failed to change privacy settings", zap.Error(err))
		return err
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("changeUserPrivacySettings").Observe(duration)
	return nil
}

func (r *userPostgresRepository) GetUserPrivacy(ctx context.Context, id int64) (*repoModel.PrivacySettings, error) {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Getting user privacy settings by id", zap.Int64("ID", id))

	stmt, err := r.db.PrepareContext(ctx, getUserPrivacySettingsQuery)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("getUserPrivacy").Inc()
		logger.Error("failed to prepare statement", zap.Error(err))
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, id)
	var privacySettings repoModel.PrivacySettings
	err = row.Scan(&privacySettings.IsPublicPlaylists,
		&privacySettings.IsPublicMinutesListened,
		&privacySettings.IsPublicFavoriteArtists,
		&privacySettings.IsPublicTracksListened,
		&privacySettings.IsPublicFavoriteTracks,
		&privacySettings.IsPublicArtistsListened)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("getUserPrivacy").Inc()
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error("user not found", zap.Error(err))
			return nil, err
		}
		logger.Error("user not found", zap.Error(err))
		return nil, err
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("getUserPrivacy").Observe(duration)
	return &privacySettings, nil
}

func (r *userPostgresRepository) GetFullUserData(ctx context.Context, username string) (*repoModel.UserFullData, error) {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Get full user data by username", zap.String("username", username))
	id, err := r.GetIDByUsername(ctx, username)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("getFullUserData").Inc()
		logger.Error("failed to get user ID", zap.Error(err))
		return nil, err
	}
	privacy, err := r.GetUserPrivacy(ctx, id)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("getFullUserData").Inc()
		logger.Error("failed to get user privacy settings", zap.Error(err))
		return nil, err
	}
	user, err := r.GetUserByID(ctx, id)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("getFullUserData").Inc()
		logger.Error("failed to get user data", zap.Error(err))
		return nil, err
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("getFullUserData").Observe(duration)
	return &repoModel.UserFullData{
		Username:  user.Username,
		Email:     user.Email,
		Thumbnail: user.Thumbnail,
		Privacy:   privacy,
	}, nil
}
