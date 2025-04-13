package repository

import (
	"bytes"
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"

	"golang.org/x/crypto/argon2"

	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/repository"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/user"
)

var (
	ErrUsernameExist = errors.New("user with this username already exists")
	ErrEmailExist    = errors.New("user with this email already exists")
	ErrUserNotFound  = errors.New("user not found")
	ErrCreateSalt    = errors.New("failed to create salt")
	ErrWrongPassword = errors.New("wrong password")
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
			SELECT u.username, u.email, u.thumbnail_url, 
				us.is_public_playlists,
				us.is_public_minutes_listened,
				us.is_public_favorite_artists,
				us.is_public_tracks_listened,
				us.is_public_favorite_tracks,
				us.is_public_artists_listened
			FROM "user" u
			INNER JOIN "user_settings" us ON u.id = us.user_id
			WHERE u.username = $1
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
			FROM STREAM
			WHERE user_id = $1
	`
	getNumUniqueArtistQuery = `
			SELECT COUNT(DISTINCT ta.artist_id) AS unique_artists_listened
			FROM stream s
			JOIN track_artist ta ON s.track_id = ta.track_id
			WHERE s.user_id = $1;
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
	row := r.db.QueryRowContext(ctx, getPasswordQuery, id)
	var storedHash string
	err := row.Scan(&storedHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", ErrUserNotFound
		}
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
	var exists bool
	err := r.db.QueryRowContext(ctx, checkUserExist, regData.Username, regData.Email).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if exists {
		return nil, ErrUsernameExist
	}

	salt := createSalt()
	if salt == nil {
		return nil, errors.New("failed to create salt")
	}
	hashedPassword := hashPassword(salt, regData.Password)

	var userID int64
	err = r.db.QueryRowContext(ctx, createUserQuery, regData.Username,
		hashedPassword, regData.Email).Scan(&userID)
	if err != nil {
		return nil, err
	}

	_, err = r.db.ExecContext(ctx, createUserSettingsQuery, userID)
	if err != nil {
		return nil, err
	}

	return &repoModel.User{
		ID:        userID,
		Username:  regData.Username,
		Email:     regData.Email,
		Thumbnail: "/default_avatar.png",
	}, nil
}

func (r *userPostgresRepository) GetUserByID(ctx context.Context, ID int64) (*repoModel.User, error) {
	row := r.db.QueryRowContext(ctx, getUserByIDQuery, ID)
	var user repoModel.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Thumbnail)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userPostgresRepository) LoginUser(ctx context.Context, logData *repoModel.User) (*repoModel.User, error) {
	var storedHash string
	row := r.db.QueryRowContext(ctx, loginUserQuery, logData.Username, logData.Email)
	var user repoModel.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &storedHash, &user.Thumbnail)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	if !checkPasswordHash(storedHash, logData.Password) {
		return nil, ErrUserNotFound
	}

	return &user, nil
}

func (r *userPostgresRepository) GetAvatar(ctx context.Context, username string) (string, error) {
	row := r.db.QueryRowContext(ctx, getAvatarQuery, username)
	var avatarUrl string
	err := row.Scan(&avatarUrl)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", ErrUserNotFound
		}
		return "", err
	}
	return avatarUrl, nil
}

func (r *userPostgresRepository) UploadAvatar(ctx context.Context, avatarUrl string, username string) error {
	_, err := r.db.ExecContext(ctx, uploadAvatarQuery, avatarUrl, username)
	if err != nil {
		return err
	}
	return nil
}

func (r *userPostgresRepository) changeUsername(ctx context.Context, id int64, newUsername string) error {
	_, err := r.db.ExecContext(ctx, changeUsernameQuery, newUsername, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *userPostgresRepository) changeEmail(ctx context.Context, id int64, newEmail string) error {
	_, err := r.db.ExecContext(ctx, changeEmailQuery, newEmail, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *userPostgresRepository) changePassword(ctx context.Context, password string, id int64, newPassword string) error {
	storedHash, err := r.getPassword(ctx, id)
	if err != nil {
		return err
	}
	if !checkPasswordHash(storedHash, password) {
		return ErrWrongPassword
	}
	salt := createSalt()
	newHashedPassword := hashPassword(salt, newPassword)
	_, err = r.db.ExecContext(ctx, changePasswordQuery, newHashedPassword, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *userPostgresRepository) ChangeUserData(ctx context.Context, username string, changeData *repoModel.ChangeUserData) error {
	id, err := r.GetIDByUsername(ctx, username)
	if err != nil {
		return err
	}
	if changeData.NewUsername != "" {
		err := r.changeUsername(ctx, id, changeData.NewUsername)
		if err != nil {
			return err
		}
	}
	if changeData.NewEmail != "" {
		err := r.changeEmail(ctx, id, changeData.NewEmail)
		if err != nil {
			return err
		}
	}
	if changeData.NewPassword != "" {
		err := r.changePassword(ctx, changeData.Password, id, changeData.NewPassword)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *userPostgresRepository) DeleteUser(ctx context.Context, user *repoModel.User) error {
	storedHash, err := r.getPassword(ctx, user.ID)
	if err != nil {
		return err
	}
	if !checkPasswordHash(storedHash, user.Password) {
		return ErrWrongPassword
	}
	_, err = r.db.ExecContext(ctx, deleteUserQuery, user.Username, user.Email)
	if err != nil {
		return err
	}
	return nil
}

func (r *userPostgresRepository) ChangeUserPrivacySettings(ctx context.Context, username string, privacySettings *repoModel.PrivacySettings) error {
	id, err := r.GetIDByUsername(ctx, username)
	if err != nil {
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
		return err
	}
	return nil
}

func (r *userPostgresRepository) GetIDByUsername(ctx context.Context, username string) (int64, error) {
	row := r.db.QueryRowContext(ctx, getIdByUsernameQuery, username)
	var userID int64
	err := row.Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, ErrUserNotFound
		}
		return 0, err
	}
	return userID, nil
}

func (r *userPostgresRepository) GetUserData(ctx context.Context, username string) (*repoModel.UserAndSettings, error) {
	row := r.db.QueryRowContext(ctx, getUserDataQuery, username)
	var user repoModel.UserAndSettings
	err := row.Scan(&user.Username, &user.Email, &user.Thumbnail,
		&user.IsPublicPlaylists, &user.IsPublicMinutesListened,
		&user.IsPublicFavoriteArtists, &user.IsPublicTracksListened,
		&user.IsPublicFavoriteTracks, &user.IsPublicArtistsListened,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userPostgresRepository) getNumUniqueTracks(ctx context.Context, id int64) (int64, error) {
	row := r.db.QueryRowContext(ctx, getNumUniqueTracksQuery, id)
	var numUniqueTracks int64
	err := row.Scan(&numUniqueTracks)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, ErrUserNotFound
		}
		return 0, err
	}
	return numUniqueTracks, nil
}

func (r *userPostgresRepository) getNumMinutes(ctx context.Context, id int64) (int64, error) {
	row := r.db.QueryRowContext(ctx, getMinutesListenedQuery, id)
	var numMinutes int64
	err := row.Scan(&numMinutes)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1, ErrUserNotFound
		}
		return -1, err
	}
	return numMinutes, nil
}

func (r *userPostgresRepository) getNumUniqueArtist(ctx context.Context, id int64) (int64, error) {
	row := r.db.QueryRowContext(ctx, getNumUniqueArtistQuery, id)
	var numUniqueArtist int64
	err := row.Scan(&numUniqueArtist)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1, ErrUserNotFound
		}
		return -1, err
	}
	return numUniqueArtist, nil

}

func (r *userPostgresRepository) GetUserStats(ctx context.Context, username string) (*repoModel.UserStats, error) {
	userID, err := r.GetIDByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	numUniqueTracks, err := r.getNumUniqueTracks(ctx, userID)
	if err != nil {
		return nil, err
	}
	numMinutes, err := r.getNumMinutes(ctx, userID)
	if err != nil {
		return nil, err
	}
	numUniqueArtists, err := r.getNumUniqueArtist(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &repoModel.UserStats{
		MinutesListened: numUniqueTracks,
		TracksListened:  numMinutes,
		ArtistsListened: numUniqueArtists,
	}, nil
}
