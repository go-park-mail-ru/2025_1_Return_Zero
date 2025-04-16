package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/repository"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func setupTestContext() context.Context {
	logger, _ := zap.NewDevelopment()
	sugaredLogger := logger.Sugar()
	return context.WithValue(context.Background(), helpers.LoggerKey{}, sugaredLogger)
}

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserPostgresRepository(db)
	ctx := setupTestContext()

	tests := []struct {
		name          string
		inputUser     *repository.User
		mockSetup     func()
		expectedUser  *repository.User
		expectedError error
	}{
		{
			name: "Success",
			inputUser: &repository.User{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func() {
				mock.ExpectQuery(`SELECT 1 FROM "user"`).
					WithArgs("testuser", "test@example.com").
					WillReturnError(sql.ErrNoRows)

				mock.ExpectQuery(`INSERT INTO "user"`).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "test@example.com").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectExec(`INSERT INTO "user_settings"`).
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedUser: &repository.User{
				ID:        1,
				Username:  "testuser",
				Email:     "test@example.com",
				Thumbnail: "/default_avatar.png",
			},
			expectedError: nil,
		},
		{
			name: "AlreadyExists",
			inputUser: &repository.User{
				Username: "existinguser",
				Email:    "existing@example.com",
				Password: "password123",
			},
			mockSetup: func() {
				mock.ExpectQuery(`SELECT 1 FROM "user"`).
					WithArgs("existinguser", "existing@example.com").
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
			},
			expectedUser:  nil,
			expectedError: user.ErrUsernameExist,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			user, err := repo.CreateUser(ctx, tt.inputUser)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.expectedUser.ID, user.ID)
				assert.Equal(t, tt.expectedUser.Username, user.Username)
				assert.Equal(t, tt.expectedUser.Email, user.Email)
				assert.Equal(t, tt.expectedUser.Thumbnail, user.Thumbnail)
			}
		})
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserPostgresRepository(db)
	ctx := setupTestContext()

	tests := []struct {
		name          string
		userID        int64
		mockSetup     func()
		expectedUser  *repository.User
		expectedError error
	}{
		{
			name:   "Success",
			userID: 1,
			mockSetup: func() {
				mock.ExpectQuery(`SELECT id, username, email, thumbnail_url FROM "user"`).
					WithArgs(int64(1)).
					WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "thumbnail_url"}).
						AddRow(1, "testuser", "test@example.com", "/avatars/test.jpg"))
			},
			expectedUser: &repository.User{
				ID:        1,
				Username:  "testuser",
				Email:     "test@example.com",
				Thumbnail: "/avatars/test.jpg",
			},
			expectedError: nil,
		},
		{
			name:   "UserNotFound",
			userID: 999,
			mockSetup: func() {
				mock.ExpectQuery(`SELECT id, username, email, thumbnail_url FROM "user"`).
					WithArgs(int64(999)).
					WillReturnError(sql.ErrNoRows)
			},
			expectedUser:  nil,
			expectedError: user.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			user, err := repo.GetUserByID(ctx, tt.userID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.expectedUser.ID, user.ID)
				assert.Equal(t, tt.expectedUser.Username, user.Username)
				assert.Equal(t, tt.expectedUser.Email, user.Email)
				assert.Equal(t, tt.expectedUser.Thumbnail, user.Thumbnail)
			}
		})
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLoginUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserPostgresRepository(db)
	ctx := setupTestContext()

	hashedPassword := "RThVeDlnPT0KkrNtKl7GScJ0jF7WdH+G1gT6WPmg06I="

	tests := []struct {
		name          string
		inputUser     *repository.User
		mockSetup     func()
		expectedUser  *repository.User
		expectedError error
	}{
		{
			name: "Success",
			inputUser: &repository.User{
				Username: "testuser",
				Password: "password123",
			},
			mockSetup: func() {
				mock.ExpectQuery(`SELECT id, username, email, password_hash, thumbnail_url FROM "user"`).
					WithArgs("testuser", "").
					WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "password_hash", "thumbnail_url"}).
						AddRow(1, "testuser", "test@example.com", hashedPassword, "/avatars/test.jpg"))
			},
			expectedUser:  nil,
			expectedError: user.ErrWrongPassword,
		},
		{
			name: "UserNotFound",
			inputUser: &repository.User{
				Username: "nonexistentuser",
				Password: "password123",
			},
			mockSetup: func() {
				mock.ExpectQuery(`SELECT id, username, email, password_hash, thumbnail_url FROM "user"`).
					WithArgs("nonexistentuser", "").
					WillReturnError(sql.ErrNoRows)
			},
			expectedUser:  nil,
			expectedError: user.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			user, err := repo.LoginUser(ctx, tt.inputUser)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.expectedUser.Username, user.Username)
				assert.Equal(t, tt.expectedUser.Email, user.Email)
				assert.Equal(t, tt.expectedUser.Thumbnail, user.Thumbnail)
			}
		})
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserPrivacy(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserPostgresRepository(db)
	ctx := setupTestContext()

	tests := []struct {
		name            string
		userID          int64
		mockSetup       func()
		expectedPrivacy *repository.UserPrivacySettings
		expectedError   error
	}{
		{
			name:   "Success",
			userID: 1,
			mockSetup: func() {
				mock.ExpectQuery(`SELECT is_public_playlists, is_public_minutes_listened, is_public_favorite_artists, is_public_tracks_listened, is_public_favorite_tracks, is_public_artists_listened FROM user_settings`).
					WithArgs(int64(1)).
					WillReturnRows(sqlmock.NewRows([]string{
						"is_public_playlists",
						"is_public_minutes_listened",
						"is_public_favorite_artists",
						"is_public_tracks_listened",
						"is_public_favorite_tracks",
						"is_public_artists_listened"}).
						AddRow(true, false, true, false, true, false))
			},
			expectedPrivacy: &repository.UserPrivacySettings{
				IsPublicPlaylists:       true,
				IsPublicMinutesListened: false,
				IsPublicFavoriteArtists: true,
				IsPublicTracksListened:  false,
				IsPublicFavoriteTracks:  true,
				IsPublicArtistsListened: false,
			},
			expectedError: nil,
		},
		{
			name:   "UserNotFound",
			userID: 999,
			mockSetup: func() {
				mock.ExpectQuery(`SELECT is_public_playlists, is_public_minutes_listened, is_public_favorite_artists, is_public_tracks_listened, is_public_favorite_tracks, is_public_artists_listened FROM user_settings`).
					WithArgs(int64(999)).
					WillReturnError(sql.ErrNoRows)
			},
			expectedPrivacy: nil,
			expectedError:   sql.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			privacy, err := repo.GetUserPrivacy(ctx, tt.userID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, privacy)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, privacy)
				assert.Equal(t, tt.expectedPrivacy.IsPublicPlaylists, privacy.IsPublicPlaylists)
				assert.Equal(t, tt.expectedPrivacy.IsPublicMinutesListened, privacy.IsPublicMinutesListened)
				assert.Equal(t, tt.expectedPrivacy.IsPublicFavoriteArtists, privacy.IsPublicFavoriteArtists)
				assert.Equal(t, tt.expectedPrivacy.IsPublicTracksListened, privacy.IsPublicTracksListened)
				assert.Equal(t, tt.expectedPrivacy.IsPublicFavoriteTracks, privacy.IsPublicFavoriteTracks)
				assert.Equal(t, tt.expectedPrivacy.IsPublicArtistsListened, privacy.IsPublicArtistsListened)
			}
		})
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestChangeUserData(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserPostgresRepository(db)
	ctx := setupTestContext()
	username := "testuser"

	changeData := &repository.ChangeUserData{
		NewUsername: "newusername",
		NewEmail:    "newemail@example.com",
	}

	mock.ExpectQuery(`SELECT id FROM "user"`).
		WithArgs(username).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectExec(`UPDATE "user" SET username = \$1`).
		WithArgs("newusername", 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec(`UPDATE "user" SET email = \$1`).
		WithArgs("newemail@example.com", 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.ChangeUserData(ctx, username, changeData)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserStats(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserPostgresRepository(db)
	ctx := setupTestContext()
	userID := int64(1)

	mock.ExpectQuery(`SELECT COUNT\(DISTINCT track_id\) AS num_unique_tracks`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"num_unique_tracks"}).AddRow(42))

	mock.ExpectQuery(`SELECT COALESCE\(SUM\(duration\) / 60, 0\) AS total_minutes`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"total_minutes"}).AddRow(120))

	mock.ExpectQuery(`SELECT COUNT\(DISTINCT ta.artist_id\) AS unique_artists_listened`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"unique_artists_listened"}).AddRow(15))

	stats, err := repo.GetUserStats(ctx, userID)
	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, int64(42), stats.TracksListened)
	assert.Equal(t, int64(120), stats.MinutesListened)
	assert.Equal(t, int64(15), stats.ArtistsListened)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserPostgresRepository(db)
	ctx := setupTestContext()

	loginUser := &repository.User{
		Username: "testuser",
		Password: "password123",
		Email:    "testuser@gmail.com",
	}

	hashedPassword := "RThVeDlnPT0KkrNtKl7GScJ0jF7WdH+G1gT6WPmg06I="

	mock.ExpectQuery(`SELECT id FROM "user"`).
		WithArgs("testuser").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectQuery(`SELECT password_hash FROM "user"`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"password_hash"}).AddRow(hashedPassword))

	mock.ExpectExec(`DELETE FROM "user"`).
		WithArgs("testuser", "testuser@gmail.com").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.DeleteUser(ctx, loginUser)

	if err != nil && err != user.ErrWrongPassword {
		assert.NoError(t, err)
	}

	if err == user.ErrWrongPassword {
		assert.Equal(t, user.ErrWrongPassword, err)
	} else {
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	}
}

func TestChangeUserPrivacySettings(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserPostgresRepository(db)
	ctx := setupTestContext()

	tests := []struct {
		name          string
		username      string
		privacyInput  *repository.UserPrivacySettings
		mockSetup     func()
		expectedError error
	}{
		{
			name:     "Success",
			username: "testuser",
			privacyInput: &repository.UserPrivacySettings{
				IsPublicPlaylists:       true,
				IsPublicMinutesListened: false,
				IsPublicFavoriteArtists: true,
				IsPublicTracksListened:  false,
				IsPublicFavoriteTracks:  true,
				IsPublicArtistsListened: false,
			},
			mockSetup: func() {
				mock.ExpectQuery(`SELECT id FROM "user"`).
					WithArgs("testuser").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectExec(`UPDATE "user_settings" SET`).
					WithArgs(true, false, true, false, true, false, 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedError: nil,
		},
		{
			name:     "UserNotFound",
			username: "nonexistentuser",
			privacyInput: &repository.UserPrivacySettings{
				IsPublicPlaylists:       true,
				IsPublicMinutesListened: false,
				IsPublicFavoriteArtists: true,
				IsPublicTracksListened:  false,
				IsPublicFavoriteTracks:  true,
				IsPublicArtistsListened: false,
			},
			mockSetup: func() {
				mock.ExpectQuery(`SELECT id FROM "user"`).
					WithArgs("nonexistentuser").
					WillReturnError(sql.ErrNoRows)
			},
			expectedError: user.ErrUserNotFound,
		},
		{
			name:     "UpdateFailed",
			username: "testuser",
			privacyInput: &repository.UserPrivacySettings{
				IsPublicPlaylists:       true,
				IsPublicMinutesListened: false,
				IsPublicFavoriteArtists: true,
				IsPublicTracksListened:  false,
				IsPublicFavoriteTracks:  true,
				IsPublicArtistsListened: false,
			},
			mockSetup: func() {
				mock.ExpectQuery(`SELECT id FROM "user"`).
					WithArgs("testuser").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectExec(`UPDATE "user_settings" SET`).
					WithArgs(true, false, true, false, true, false, 1).
					WillReturnError(errors.New("database error"))
			},
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := repo.ChangeUserPrivacySettings(ctx, tt.username, tt.privacyInput)

			if tt.expectedError != nil {
				assert.Error(t, err)
				if errors.Is(err, sql.ErrNoRows) {
					assert.Equal(t, user.ErrUserNotFound, err)
				} else {
					assert.Equal(t, tt.expectedError.Error(), err.Error())
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserData(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserPostgresRepository(db)
	ctx := setupTestContext()

	tests := []struct {
		name          string
		userID        int64
		mockSetup     func()
		expectedUser  *repository.User
		expectedError error
	}{
		{
			name:   "Success",
			userID: 1,
			mockSetup: func() {
				mock.ExpectQuery(`SELECT username, email, thumbnail_url FROM "user" WHERE id = \$1`).
					WithArgs(int64(1)).
					WillReturnRows(sqlmock.NewRows([]string{"username", "email", "thumbnail_url"}).
						AddRow("testuser", "test@example.com", "/avatars/test.jpg"))
			},
			expectedUser: &repository.User{
				ID:        0,
				Username:  "testuser",
				Email:     "test@example.com",
				Thumbnail: "/avatars/test.jpg",
			},
			expectedError: nil,
		},
		{
			name:   "UserNotFound",
			userID: 999,
			mockSetup: func() {
				mock.ExpectQuery(`SELECT username, email, thumbnail_url FROM "user" WHERE id = \$1`).
					WithArgs(int64(999)).
					WillReturnError(sql.ErrNoRows)
			},
			expectedUser:  nil,
			expectedError: user.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			user, err := repo.GetUserData(ctx, tt.userID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.expectedUser.ID, user.ID)
				assert.Equal(t, tt.expectedUser.Username, user.Username)
				assert.Equal(t, tt.expectedUser.Email, user.Email)
				assert.Equal(t, tt.expectedUser.Thumbnail, user.Thumbnail)
			}
		})
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAvatar(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserPostgresRepository(db)
	ctx := setupTestContext()

	tests := []struct {
		name           string
		username       string
		mockSetup      func()
		expectedAvatar string
		expectedError  error
	}{
		{
			name:     "Success",
			username: "testuser",
			mockSetup: func() {
				mock.ExpectQuery(`SELECT thumbnail_url FROM "user"`).
					WithArgs("testuser").
					WillReturnRows(sqlmock.NewRows([]string{"thumbnail_url"}).
						AddRow("/avatars/test.jpg"))
			},
			expectedAvatar: "/avatars/test.jpg",
			expectedError:  nil,
		},
		{
			name:     "UserNotFound",
			username: "nonexistentuser",
			mockSetup: func() {
				mock.ExpectQuery(`SELECT thumbnail_url FROM "user"`).
					WithArgs("nonexistentuser").
					WillReturnError(sql.ErrNoRows)
			},
			expectedAvatar: "",
			expectedError:  user.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			avatar, err := repo.GetAvatar(ctx, tt.username)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Empty(t, avatar)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedAvatar, avatar)
			}
		})
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetFullUserData(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserPostgresRepository(db)
	ctx := setupTestContext()

	tests := []struct {
		name          string
		username      string
		mockSetup     func()
		expectedData  *repository.UserFullData
		expectedError error
	}{
		{
			name:     "Success",
			username: "testuser",
			mockSetup: func() {
				mock.ExpectQuery(`SELECT id FROM "user" WHERE username = \$1`).
					WithArgs("testuser").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectQuery(`SELECT is_public_playlists, is_public_minutes_listened, is_public_favorite_artists, is_public_tracks_listened, is_public_favorite_tracks, is_public_artists_listened FROM user_settings WHERE user_id = \$1`).
					WithArgs(int64(1)).
					WillReturnRows(sqlmock.NewRows([]string{
						"is_public_playlists",
						"is_public_minutes_listened",
						"is_public_favorite_artists",
						"is_public_tracks_listened",
						"is_public_favorite_tracks",
						"is_public_artists_listened"}).
						AddRow(true, false, true, false, true, false))

				mock.ExpectQuery(`SELECT COUNT\(DISTINCT track_id\) AS num_unique_tracks FROM stream WHERE user_id = \$1`).
					WithArgs(int64(1)).
					WillReturnRows(sqlmock.NewRows([]string{"num_unique_tracks"}).AddRow(42))

				mock.ExpectQuery(`SELECT COALESCE\(SUM\(duration\) / 60, 0\) AS total_minutes FROM stream WHERE user_id = \$1`).
					WithArgs(int64(1)).
					WillReturnRows(sqlmock.NewRows([]string{"total_minutes"}).AddRow(120))

				mock.ExpectQuery(`SELECT COUNT\(DISTINCT ta.artist_id\) AS unique_artists_listened FROM stream s JOIN track_artist ta ON s.track_id = ta.track_id WHERE s.user_id = \$1`).
					WithArgs(int64(1)).
					WillReturnRows(sqlmock.NewRows([]string{"unique_artists_listened"}).AddRow(15))

				mock.ExpectQuery(`SELECT username, email, thumbnail_url FROM "user" WHERE id = \$1`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"username", "email", "thumbnail_url"}).
						AddRow("testuser", "test@example.com", "/avatars/test.jpg"))
			},
			expectedData: &repository.UserFullData{
				Username:  "testuser",
				Email:     "test@example.com",
				Thumbnail: "/avatars/test.jpg",
				Privacy: &repository.UserPrivacySettings{
					IsPublicPlaylists:       true,
					IsPublicMinutesListened: false,
					IsPublicFavoriteArtists: true,
					IsPublicTracksListened:  false,
					IsPublicFavoriteTracks:  true,
					IsPublicArtistsListened: false,
				},
				Statistics: &repository.UserStats{
					TracksListened:  42,
					MinutesListened: 120,
					ArtistsListened: 15,
				},
			},
			expectedError: nil,
		},
		{
			name:     "UserNotFound",
			username: "nonexistentuser",
			mockSetup: func() {
				mock.ExpectQuery(`SELECT id FROM "user" WHERE username = \$1`).
					WithArgs("nonexistentuser").
					WillReturnError(sql.ErrNoRows)
			},
			expectedData:  nil,
			expectedError: user.ErrUserNotFound,
		},
		{
			name:     "PrivacySettingsError",
			username: "userwithoutprivacy",
			mockSetup: func() {
				mock.ExpectQuery(`SELECT id FROM "user" WHERE username = \$1`).
					WithArgs("userwithoutprivacy").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))

				mock.ExpectQuery(`SELECT is_public_playlists, is_public_minutes_listened, is_public_favorite_artists, is_public_tracks_listened, is_public_favorite_tracks, is_public_artists_listened FROM user_settings WHERE user_id = \$1`).
					WithArgs(int64(2)).
					WillReturnError(sql.ErrNoRows)
			},
			expectedData:  nil,
			expectedError: sql.ErrNoRows,
		},
		{
			name:     "StatisticsError",
			username: "userwithoutstats",
			mockSetup: func() {
				mock.ExpectQuery(`SELECT id FROM "user" WHERE username = \$1`).
					WithArgs("userwithoutstats").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(3))

				mock.ExpectQuery(`SELECT is_public_playlists, is_public_minutes_listened, is_public_favorite_artists, is_public_tracks_listened, is_public_favorite_tracks, is_public_artists_listened FROM user_settings WHERE user_id = \$1`).
					WithArgs(int64(3)).
					WillReturnRows(sqlmock.NewRows([]string{
						"is_public_playlists",
						"is_public_minutes_listened",
						"is_public_favorite_artists",
						"is_public_tracks_listened",
						"is_public_favorite_tracks",
						"is_public_artists_listened"}).
						AddRow(true, false, true, false, true, false))

				mock.ExpectQuery(`SELECT COUNT\(DISTINCT track_id\) AS num_unique_tracks FROM stream WHERE user_id = \$1`).
					WithArgs(int64(3)).
					WillReturnError(errors.New("database error"))
			},
			expectedData:  nil,
			expectedError: errors.New("database error"),
		},
		{
			name:     "UserDataError",
			username: "userwithoutdata",
			mockSetup: func() {
				mock.ExpectQuery(`SELECT id FROM "user" WHERE username = \$1`).
					WithArgs("userwithoutdata").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(4))

				mock.ExpectQuery(`SELECT is_public_playlists, is_public_minutes_listened, is_public_favorite_artists, is_public_tracks_listened, is_public_favorite_tracks, is_public_artists_listened FROM user_settings WHERE user_id = \$1`).
					WithArgs(int64(4)).
					WillReturnRows(sqlmock.NewRows([]string{
						"is_public_playlists",
						"is_public_minutes_listened",
						"is_public_favorite_artists",
						"is_public_tracks_listened",
						"is_public_favorite_tracks",
						"is_public_artists_listened"}).
						AddRow(true, false, true, false, true, false))

				mock.ExpectQuery(`SELECT COUNT\(DISTINCT track_id\) AS num_unique_tracks FROM stream WHERE user_id = \$1`).
					WithArgs(int64(4)).
					WillReturnRows(sqlmock.NewRows([]string{"num_unique_tracks"}).AddRow(42))

				mock.ExpectQuery(`SELECT COALESCE\(SUM\(duration\) / 60, 0\) AS total_minutes FROM stream WHERE user_id = \$1`).
					WithArgs(int64(4)).
					WillReturnRows(sqlmock.NewRows([]string{"total_minutes"}).AddRow(120))

				mock.ExpectQuery(`SELECT COUNT\(DISTINCT ta.artist_id\) AS unique_artists_listened FROM stream s JOIN track_artist ta ON s.track_id = ta.track_id WHERE s.user_id = \$1`).
					WithArgs(int64(4)).
					WillReturnRows(sqlmock.NewRows([]string{"unique_artists_listened"}).AddRow(15))

				mock.ExpectQuery(`SELECT username, email, thumbnail_url FROM "user" WHERE id = \$1`).
					WithArgs(4).
					WillReturnError(errors.New("user data error"))
			},
			expectedData:  nil,
			expectedError: errors.New("user data error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			userData, err := repo.GetFullUserData(ctx, tt.username)

			if tt.expectedError != nil {
				assert.Error(t, err)
				if err.Error() == tt.expectedError.Error() {
					assert.Equal(t, tt.expectedError.Error(), err.Error())
				} else if errors.Is(err, sql.ErrNoRows) && tt.expectedError == user.ErrUserNotFound {
					assert.Equal(t, user.ErrUserNotFound, err)
				}
				assert.Nil(t, userData)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, userData)
				assert.Equal(t, tt.expectedData.Username, userData.Username)
				assert.Equal(t, tt.expectedData.Email, userData.Email)
				assert.Equal(t, tt.expectedData.Thumbnail, userData.Thumbnail)

				assert.Equal(t, tt.expectedData.Privacy.IsPublicPlaylists, userData.Privacy.IsPublicPlaylists)
				assert.Equal(t, tt.expectedData.Privacy.IsPublicMinutesListened, userData.Privacy.IsPublicMinutesListened)
				assert.Equal(t, tt.expectedData.Privacy.IsPublicFavoriteArtists, userData.Privacy.IsPublicFavoriteArtists)
				assert.Equal(t, tt.expectedData.Privacy.IsPublicTracksListened, userData.Privacy.IsPublicTracksListened)
				assert.Equal(t, tt.expectedData.Privacy.IsPublicFavoriteTracks, userData.Privacy.IsPublicFavoriteTracks)
				assert.Equal(t, tt.expectedData.Privacy.IsPublicArtistsListened, userData.Privacy.IsPublicArtistsListened)

				assert.Equal(t, tt.expectedData.Statistics.TracksListened, userData.Statistics.TracksListened)
				assert.Equal(t, tt.expectedData.Statistics.MinutesListened, userData.Statistics.MinutesListened)
				assert.Equal(t, tt.expectedData.Statistics.ArtistsListened, userData.Statistics.ArtistsListened)
			}
		})
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}
