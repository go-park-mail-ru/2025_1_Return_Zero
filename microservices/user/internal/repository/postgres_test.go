package repository

import (
	"context"
	"database/sql"
	"encoding/base64"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	loggerPkg "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/metrics"
	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/user/model/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"golang.org/x/crypto/argon2"
)

const (
	testUserID            = int64(1)
	testUsername          = "testuser"
	testEmail             = "test@example.com"
	testPassword          = "password123"
	testAvatarURL         = "avatar.jpg"
	testNewUsername       = "newuser"
	testNewEmail          = "newemail@example.com"
	testNewPassword       = "newpass"
	testOldPassword       = "oldpass"
	existingUsername      = "existinguser"
	existingEmail         = "existing@example.com"
	nonExistentUsername   = "nonexistentuser"
	nonExistentEmail      = "nonexistent@example.com"
	wrongPassword         = "wrongpassword"
	correctPassword       = "correctpassword"
	newAvatarURL          = "new-avatar.jpg"
	differentPasswordHash = "differentPasswordHash"
)

func setupTest(t *testing.T) (*sql.DB, sqlmock.Sqlmock, context.Context) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	logger := zap.NewNop().Sugar()
	ctx := loggerPkg.LoggerToContext(context.Background(), logger)

	return db, mock, ctx
}

func TestCreateUser(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewUserPostgresRepository(db, metrics.NewMockMetrics())

	regData := &repoModel.RegisterData{
		Username: testUsername,
		Email:    testEmail,
		Password: testPassword,
	}

	mock.ExpectPrepare("INSERT INTO \"user\"")

	mock.ExpectPrepare("SELECT 1").ExpectQuery().WithArgs(regData.Username, regData.Email).WillReturnError(sql.ErrNoRows)

	mock.ExpectQuery("INSERT INTO \"user\"").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(testUserID))

	mock.ExpectPrepare("INSERT INTO \"user_settings\"").ExpectExec().WithArgs(testUserID).WillReturnResult(sqlmock.NewResult(1, 1))

	user, err := repo.CreateUser(ctx, regData)

	assert.NoError(t, err)
	assert.Equal(t, testUserID, user.ID)
	assert.Equal(t, testUsername, user.Username)
	assert.Equal(t, testEmail, user.Email)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateUserAlreadyExists(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewUserPostgresRepository(db, metrics.NewMockMetrics())

	regData := &repoModel.RegisterData{
		Username: testUsername,
		Email:    testEmail,
		Password: testPassword,
	}

	mock.ExpectPrepare("INSERT INTO \"user\"")

	rows := sqlmock.NewRows([]string{"1"}).AddRow(1)
	mock.ExpectPrepare("SELECT 1").ExpectQuery().WithArgs(regData.Username, regData.Email).WillReturnRows(rows)

	user, err := repo.CreateUser(ctx, regData)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateUserFailureDuplicateUser(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewUserPostgresRepository(db, metrics.NewMockMetrics())

	regData := &repoModel.RegisterData{
		Username: existingUsername,
		Email:    existingEmail,
		Password: testPassword,
	}

	mock.ExpectPrepare("INSERT INTO \"user\"")

	rows := sqlmock.NewRows([]string{"1"}).AddRow(1)
	mock.ExpectPrepare("SELECT 1").ExpectQuery().WithArgs(regData.Username, regData.Email).WillReturnRows(rows)

	user, err := repo.CreateUser(ctx, regData)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLoginUser(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewUserPostgresRepository(db, metrics.NewMockMetrics())

	loginData := &repoModel.LoginData{
		Username: testUsername,
		Email:    testEmail,
		Password: testPassword,
	}

	salt := make([]byte, 8)
	hash := argon2.IDKey([]byte(loginData.Password), salt, 1, 64*1024, 4, 32)
	combined := append(salt, hash...)
	passwordHash := base64.StdEncoding.EncodeToString(combined)

	rows := sqlmock.NewRows([]string{"id", "username", "email", "password_hash", "thumbnail_url", "label_id"}).
		AddRow(testUserID, testUsername, testEmail, passwordHash, testAvatarURL, nil)

	lowerUsername := strings.ToLower(loginData.Username)
	mock.ExpectPrepare("SELECT id, username, email, password_hash, thumbnail_url, label_id").
		ExpectQuery().WithArgs(lowerUsername, loginData.Email).
		WillReturnRows(rows)

	user, err := repo.LoginUser(ctx, loginData)

	assert.NoError(t, err)
	assert.Equal(t, testUserID, user.ID)
	assert.Equal(t, testUsername, user.Username)
	assert.Equal(t, testEmail, user.Email)
	assert.Equal(t, testAvatarURL, user.Thumbnail)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLoginUserFailure(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewUserPostgresRepository(db, metrics.NewMockMetrics())

	loginData := &repoModel.LoginData{
		Username: nonExistentUsername,
		Email:    nonExistentEmail,
		Password: testPassword,
	}

	mock.ExpectPrepare("SELECT id, username, email, password_hash, thumbnail_url, label_id").
		ExpectQuery().WithArgs(loginData.Username, loginData.Email).
		WillReturnError(sql.ErrNoRows)

	user, err := repo.LoginUser(ctx, loginData)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByID(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewUserPostgresRepository(db, metrics.NewMockMetrics())

	userID := testUserID

	rows := sqlmock.NewRows([]string{"id", "username", "email", "thumbnail_url"}).
		AddRow(userID, testUsername, testEmail, testAvatarURL)

	mock.ExpectPrepare("SELECT id, username, email, thumbnail_url").
		ExpectQuery().WithArgs(userID).
		WillReturnRows(rows)

	user, err := repo.GetUserByID(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, testUsername, user.Username)
	assert.Equal(t, testEmail, user.Email)
	assert.Equal(t, testAvatarURL, user.Thumbnail)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByIDFailure(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewUserPostgresRepository(db, metrics.NewMockMetrics())

	userID := int64(999)

	mock.ExpectPrepare("SELECT id, username, email, thumbnail_url").
		ExpectQuery().WithArgs(userID).
		WillReturnError(sql.ErrNoRows)

	user, err := repo.GetUserByID(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetIDByUsername(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewUserPostgresRepository(db, metrics.NewMockMetrics())

	username := testUsername
	expectedID := testUserID

	rows := sqlmock.NewRows([]string{"id"}).AddRow(expectedID)

	mock.ExpectPrepare("SELECT id").
		ExpectQuery().WithArgs(username).
		WillReturnRows(rows)

	id, err := repo.GetIDByUsername(ctx, username)

	assert.NoError(t, err)
	assert.Equal(t, expectedID, id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserPrivacy(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewUserPostgresRepository(db, metrics.NewMockMetrics())

	userID := testUserID

	rows := sqlmock.NewRows([]string{
		"is_public_playlists",
		"is_public_minutes_listened",
		"is_public_favorite_artists",
		"is_public_tracks_listened",
		"is_public_favorite_tracks",
		"is_public_artists_listened",
	}).AddRow(true, false, true, false, true, false)

	mock.ExpectPrepare("SELECT is_public_playlists, is_public_minutes_listened").
		ExpectQuery().WithArgs(userID).
		WillReturnRows(rows)

	privacy, err := repo.GetUserPrivacy(ctx, userID)

	assert.NoError(t, err)
	assert.True(t, privacy.IsPublicPlaylists)
	assert.False(t, privacy.IsPublicMinutesListened)
	assert.True(t, privacy.IsPublicFavoriteArtists)
	assert.False(t, privacy.IsPublicTracksListened)
	assert.True(t, privacy.IsPublicFavoriteTracks)
	assert.False(t, privacy.IsPublicArtistsListened)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserPrivacyFailure(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewUserPostgresRepository(db, metrics.NewMockMetrics())

	userID := int64(999)

	mock.ExpectPrepare("SELECT is_public_playlists").
		ExpectQuery().WithArgs(userID).
		WillReturnError(sql.ErrNoRows)

	privacy, err := repo.GetUserPrivacy(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, privacy)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUploadAvatar(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewUserPostgresRepository(db, metrics.NewMockMetrics())

	avatarURL := newAvatarURL
	userID := testUserID

	mock.ExpectPrepare("UPDATE").
		ExpectExec().WithArgs(avatarURL, userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.UploadAvatar(ctx, avatarURL, userID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUploadAvatarFailure(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewUserPostgresRepository(db, metrics.NewMockMetrics())

	avatarURL := testAvatarURL
	userID := int64(999)

	mock.ExpectPrepare("UPDATE").
		ExpectExec().WithArgs(avatarURL, userID).
		WillReturnError(sql.ErrNoRows)

	err := repo.UploadAvatar(ctx, avatarURL, userID)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestChangeUserPrivacySettings(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewUserPostgresRepository(db, metrics.NewMockMetrics())

	username := testUsername
	userID := testUserID

	privacy := &repoModel.PrivacySettings{
		IsPublicPlaylists:       true,
		IsPublicMinutesListened: false,
		IsPublicFavoriteArtists: true,
		IsPublicTracksListened:  false,
		IsPublicFavoriteTracks:  true,
		IsPublicArtistsListened: false,
	}

	mock.ExpectPrepare("UPDATE \"user_settings\"")

	idRows := sqlmock.NewRows([]string{"id"}).AddRow(userID)
	mock.ExpectPrepare("SELECT id").ExpectQuery().WithArgs(username).WillReturnRows(idRows)

	mock.ExpectExec("UPDATE \"user_settings\"").
		WithArgs(
			privacy.IsPublicPlaylists,
			privacy.IsPublicMinutesListened,
			privacy.IsPublicFavoriteArtists,
			privacy.IsPublicTracksListened,
			privacy.IsPublicFavoriteTracks,
			privacy.IsPublicArtistsListened,
			userID,
		).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.ChangeUserPrivacySettings(ctx, username, privacy)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestChangeUserPrivacySettingsFailure(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewUserPostgresRepository(db, metrics.NewMockMetrics())

	username := nonExistentUsername

	privacy := &repoModel.PrivacySettings{
		IsPublicPlaylists:       true,
		IsPublicMinutesListened: false,
		IsPublicFavoriteArtists: true,
		IsPublicTracksListened:  false,
		IsPublicFavoriteTracks:  true,
		IsPublicArtistsListened: false,
	}

	mock.ExpectPrepare("UPDATE \"user_settings\"")

	mock.ExpectPrepare("SELECT id").
		ExpectQuery().WithArgs(username).
		WillReturnError(sql.ErrNoRows)

	err := repo.ChangeUserPrivacySettings(ctx, username, privacy)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetFullUserData(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewUserPostgresRepository(db, metrics.NewMockMetrics())

	username := testUsername
	userID := testUserID

	idRows := sqlmock.NewRows([]string{"id"}).AddRow(userID)
	mock.ExpectPrepare("SELECT id").ExpectQuery().WithArgs(username).WillReturnRows(idRows)

	privacyRows := sqlmock.NewRows([]string{
		"is_public_playlists",
		"is_public_minutes_listened",
		"is_public_favorite_artists",
		"is_public_tracks_listened",
		"is_public_favorite_tracks",
		"is_public_artists_listened",
	}).AddRow(true, false, true, false, true, false)
	mock.ExpectPrepare("SELECT is_public_playlists").
		ExpectQuery().WithArgs(userID).
		WillReturnRows(privacyRows)

	userRows := sqlmock.NewRows([]string{"id", "username", "email", "thumbnail_url"}).
		AddRow(userID, testUsername, testEmail, testAvatarURL)
	mock.ExpectPrepare("SELECT id, username, email, thumbnail_url").
		ExpectQuery().WithArgs(userID).
		WillReturnRows(userRows)

	userData, err := repo.GetFullUserData(ctx, username)

	assert.NoError(t, err)
	assert.NotNil(t, userData)
	assert.Equal(t, testUsername, userData.Username)
	assert.Equal(t, testEmail, userData.Email)
	assert.Equal(t, testAvatarURL, userData.Thumbnail)
	assert.True(t, userData.Privacy.IsPublicPlaylists)
	assert.False(t, userData.Privacy.IsPublicMinutesListened)
	assert.True(t, userData.Privacy.IsPublicFavoriteArtists)
	assert.False(t, userData.Privacy.IsPublicTracksListened)
	assert.True(t, userData.Privacy.IsPublicFavoriteTracks)
	assert.False(t, userData.Privacy.IsPublicArtistsListened)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetFullUserDataFailure(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewUserPostgresRepository(db, metrics.NewMockMetrics())

	username := nonExistentUsername

	mock.ExpectPrepare("SELECT id").
		ExpectQuery().WithArgs(username).
		WillReturnError(sql.ErrNoRows)

	userData, err := repo.GetFullUserData(ctx, username)

	assert.Error(t, err)
	assert.Nil(t, userData)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteUser(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewUserPostgresRepository(db, metrics.NewMockMetrics())

	userDelete := &repoModel.UserDelete{
		Username: testUsername,
		Email:    testEmail,
		Password: testPassword,
	}

	mock.ExpectPrepare("DELETE FROM")

	idRows := sqlmock.NewRows([]string{"id"}).AddRow(testUserID)
	mock.ExpectPrepare("SELECT id").ExpectQuery().WithArgs(userDelete.Username).WillReturnRows(idRows)

	salt := make([]byte, 8)
	hash := argon2.IDKey([]byte(userDelete.Password), salt, 1, 64*1024, 4, 32)
	combined := append(salt, hash...)
	passwordHash := base64.StdEncoding.EncodeToString(combined)

	passRows := sqlmock.NewRows([]string{"password_hash"}).AddRow(passwordHash)
	mock.ExpectPrepare("SELECT password_hash").ExpectQuery().WithArgs(testUserID).WillReturnRows(passRows)

	mock.ExpectExec("DELETE FROM").
		WithArgs(userDelete.Username, userDelete.Email).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.DeleteUser(ctx, userDelete)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteUserFailureUserNotFound(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewUserPostgresRepository(db, metrics.NewMockMetrics())

	userDelete := &repoModel.UserDelete{
		Username: nonExistentUsername,
		Email:    nonExistentEmail,
		Password: testPassword,
	}

	mock.ExpectPrepare("DELETE FROM \"user\"")
	mock.ExpectPrepare("SELECT id").ExpectQuery().WithArgs(userDelete.Username).WillReturnError(sql.ErrNoRows)

	err := repo.DeleteUser(ctx, userDelete)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteUserFailureWrongPassword(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewUserPostgresRepository(db, metrics.NewMockMetrics())

	userDelete := &repoModel.UserDelete{
		Username: testUsername,
		Email:    testEmail,
		Password: wrongPassword,
	}

	mock.ExpectPrepare("DELETE FROM \"user\"")

	idRows := sqlmock.NewRows([]string{"id"}).AddRow(testUserID)
	mock.ExpectPrepare("SELECT id").ExpectQuery().WithArgs(userDelete.Username).WillReturnRows(idRows)

	correctPasswordHash := differentPasswordHash
	passRows := sqlmock.NewRows([]string{"password_hash"}).AddRow(correctPasswordHash)
	mock.ExpectPrepare("SELECT password_hash").ExpectQuery().WithArgs(testUserID).WillReturnRows(passRows)

	err := repo.DeleteUser(ctx, userDelete)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestChangeUserData(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewUserPostgresRepository(db, metrics.NewMockMetrics())

	username := testUsername
	userID := testUserID

	changeData := &repoModel.ChangeUserData{
		Password:    testOldPassword,
		NewUsername: testNewUsername,
		NewEmail:    testNewEmail,
		NewPassword: testNewPassword,
	}

	salt := make([]byte, 8)
	hash := argon2.IDKey([]byte(changeData.Password), salt, 1, 64*1024, 4, 32)
	combined := append(salt, hash...)
	passwordHash := base64.StdEncoding.EncodeToString(combined)

	idRows := sqlmock.NewRows([]string{"id"}).AddRow(userID)
	mock.ExpectPrepare("SELECT id").ExpectQuery().WithArgs(username).WillReturnRows(idRows)

	mock.ExpectPrepare("UPDATE \"user\" SET username")

	mock.ExpectPrepare("SELECT 1").ExpectQuery().WithArgs(strings.ToLower(changeData.NewUsername)).WillReturnError(sql.ErrNoRows)

	mock.ExpectExec("UPDATE \"user\" SET username").WithArgs(strings.ToLower(changeData.NewUsername), userID).WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectPrepare("UPDATE \"user\" SET email")

	mock.ExpectPrepare("SELECT 1").ExpectQuery().WithArgs(changeData.NewEmail).WillReturnError(sql.ErrNoRows)

	mock.ExpectExec("UPDATE \"user\" SET email").WithArgs(changeData.NewEmail, userID).WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectPrepare("UPDATE \"user\" SET password_hash")

	passRows := sqlmock.NewRows([]string{"password_hash"}).AddRow(passwordHash)
	mock.ExpectPrepare("SELECT password_hash").ExpectQuery().WithArgs(userID).WillReturnRows(passRows)

	mock.ExpectExec("UPDATE \"user\" SET password_hash").WithArgs(sqlmock.AnyArg(), userID).WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.ChangeUserData(ctx, username, changeData)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	mock.ExpectPrepare("SELECT id").ExpectQuery().WithArgs("nonexistent").WillReturnError(sql.ErrNoRows)
	err = repo.ChangeUserData(ctx, "nonexistent", changeData)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestChangeUserDataFailureUsernameExists(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewUserPostgresRepository(db, metrics.NewMockMetrics())

	username := testUsername
	userID := testUserID
	changeData := &repoModel.ChangeUserData{
		Password:    testOldPassword,
		NewUsername: existingUsername,
	}

	idRows := sqlmock.NewRows([]string{"id"}).AddRow(userID)
	mock.ExpectPrepare("SELECT id").ExpectQuery().WithArgs(username).WillReturnRows(idRows)

	mock.ExpectPrepare("UPDATE \"user\" SET username")

	rows := sqlmock.NewRows([]string{"1"}).AddRow(1)
	mock.ExpectPrepare("SELECT 1").
		ExpectQuery().WithArgs(strings.ToLower(changeData.NewUsername)).
		WillReturnRows(rows)

	err := repo.ChangeUserData(ctx, username, changeData)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestChangeUserDataFailureWrongPassword(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewUserPostgresRepository(db, metrics.NewMockMetrics())

	username := testUsername
	userID := testUserID
	changeData := &repoModel.ChangeUserData{
		Password:    wrongPassword,
		NewPassword: "newpassword",
	}

	idRows := sqlmock.NewRows([]string{"id"}).AddRow(userID)
	mock.ExpectPrepare("SELECT id").ExpectQuery().WithArgs(username).WillReturnRows(idRows)

	mock.ExpectPrepare("UPDATE \"user\" SET password_hash")

	salt := make([]byte, 8)
	hash := argon2.IDKey([]byte(correctPassword), salt, 1, 64*1024, 4, 32)
	combined := append(salt, hash...)
	passwordHash := base64.StdEncoding.EncodeToString(combined)

	passRows := sqlmock.NewRows([]string{"password_hash"}).AddRow(passwordHash)
	mock.ExpectPrepare("SELECT password_hash").ExpectQuery().WithArgs(userID).WillReturnRows(passRows)

	err := repo.ChangeUserData(ctx, username, changeData)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestChangeUserDataOnlyPassword(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewUserPostgresRepository(db, metrics.NewMockMetrics())

	username := testUsername
	userID := testUserID
	changeData := &repoModel.ChangeUserData{
		Password:    testOldPassword,
		NewPassword: testNewPassword,
	}

	salt := make([]byte, 8)
	hash := argon2.IDKey([]byte(changeData.Password), salt, 1, 64*1024, 4, 32)
	combined := append(salt, hash...)
	passwordHash := base64.StdEncoding.EncodeToString(combined)

	idRows := sqlmock.NewRows([]string{"id"}).AddRow(userID)
	mock.ExpectPrepare("SELECT id").ExpectQuery().WithArgs(username).WillReturnRows(idRows)

	mock.ExpectPrepare("UPDATE \"user\" SET password_hash")

	passRows := sqlmock.NewRows([]string{"password_hash"}).AddRow(passwordHash)
	mock.ExpectPrepare("SELECT password_hash").ExpectQuery().WithArgs(userID).WillReturnRows(passRows)

	mock.ExpectExec("UPDATE \"user\" SET password_hash").WithArgs(sqlmock.AnyArg(), userID).WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.ChangeUserData(ctx, username, changeData)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestChangeUserDataOnlyEmail(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewUserPostgresRepository(db, metrics.NewMockMetrics())

	username := testUsername
	userID := testUserID

	changeData := &repoModel.ChangeUserData{
		Password: testOldPassword,
		NewEmail: testNewEmail,
	}

	idRows := sqlmock.NewRows([]string{"id"}).AddRow(userID)
	mock.ExpectPrepare("SELECT id").ExpectQuery().WithArgs(username).WillReturnRows(idRows)

	mock.ExpectPrepare("UPDATE \"user\" SET email")

	mock.ExpectPrepare("SELECT 1").
		ExpectQuery().WithArgs(changeData.NewEmail).
		WillReturnError(sql.ErrNoRows)

	mock.ExpectExec("UPDATE \"user\" SET email").WithArgs(changeData.NewEmail, userID).WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.ChangeUserData(ctx, username, changeData)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestChangeUserDataOnlyUsername(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewUserPostgresRepository(db, metrics.NewMockMetrics())

	username := testUsername
	userID := testUserID

	changeData := &repoModel.ChangeUserData{
		Password:    testOldPassword,
		NewUsername: testNewUsername,
	}

	idRows := sqlmock.NewRows([]string{"id"}).AddRow(userID)
	mock.ExpectPrepare("SELECT id").ExpectQuery().WithArgs(username).WillReturnRows(idRows)

	mock.ExpectPrepare("UPDATE \"user\" SET username")

	mock.ExpectPrepare("SELECT 1").
		ExpectQuery().WithArgs(changeData.NewUsername).
		WillReturnError(sql.ErrNoRows)

	mock.ExpectExec("UPDATE \"user\" SET username").WithArgs(changeData.NewUsername, userID).WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.ChangeUserData(ctx, username, changeData)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
