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
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	mock.ExpectQuery("SELECT 1").WithArgs(regData.Username, regData.Email).WillReturnError(sql.ErrNoRows)

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectQuery("INSERT INTO").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnRows(rows)

	mock.ExpectExec("INSERT INTO").WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))

	user, err := repo.CreateUser(ctx, regData)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), user.ID)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "test@example.com", user.Email)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateUserAlreadyExists(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewUserPostgresRepository(db, metrics.NewMockMetrics())

	regData := &repoModel.RegisterData{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	rows := sqlmock.NewRows([]string{"1"}).AddRow(1)
	mock.ExpectQuery("SELECT 1").WithArgs(regData.Username, regData.Email).WillReturnRows(rows)

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
        Username: "testuser",
        Email:    "test@example.com",
        Password: "password123",
    }

    salt := make([]byte, 8)
    hash := argon2.IDKey([]byte(loginData.Password), salt, 1, 64*1024, 4, 32)
    combined := append(salt, hash...)
    passwordHash := base64.StdEncoding.EncodeToString(combined)

    rows := sqlmock.NewRows([]string{"id", "username", "email", "password_hash", "thumbnail_url"}).
        AddRow(1, "testuser", "test@example.com", passwordHash, "avatar.jpg")

    lowerUsername := strings.ToLower(loginData.Username)
    mock.ExpectQuery("SELECT id, username, email, password_hash, thumbnail_url").
        WithArgs(lowerUsername, loginData.Email).
        WillReturnRows(rows)

    user, err := repo.LoginUser(ctx, loginData)

    assert.NoError(t, err)
    assert.Equal(t, int64(1), user.ID)
    assert.Equal(t, "testuser", user.Username)
    assert.Equal(t, "test@example.com", user.Email)
    assert.Equal(t, "avatar.jpg", user.Thumbnail)
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByID(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewUserPostgresRepository(db, metrics.NewMockMetrics())

	userID := int64(1)

	rows := sqlmock.NewRows([]string{"id", "username", "email", "thumbnail_url"}).
		AddRow(userID, "testuser", "test@example.com", "avatar.jpg")

	mock.ExpectQuery("SELECT id, username, email, thumbnail_url").
		WithArgs(userID).
		WillReturnRows(rows)

	user, err := repo.GetUserByID(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "avatar.jpg", user.Thumbnail)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetIDByUsername(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewUserPostgresRepository(db, metrics.NewMockMetrics())

	username := "testuser"
	expectedID := int64(1)

	rows := sqlmock.NewRows([]string{"id"}).AddRow(expectedID)

	mock.ExpectQuery("SELECT id").
		WithArgs(username).
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

	userID := int64(1)

	rows := sqlmock.NewRows([]string{
		"is_public_playlists",
		"is_public_minutes_listened",
		"is_public_favorite_artists",
		"is_public_tracks_listened",
		"is_public_favorite_tracks",
		"is_public_artists_listened",
	}).AddRow(true, false, true, false, true, false)

	mock.ExpectQuery("SELECT is_public_playlists, is_public_minutes_listened").
		WithArgs(userID).
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

func TestUploadAvatar(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewUserPostgresRepository(db, metrics.NewMockMetrics())

	avatarURL := "new-avatar.jpg"
	userID := int64(1)

	mock.ExpectExec("UPDATE \"user\"").
		WithArgs(avatarURL, userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.UploadAvatar(ctx, avatarURL, userID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestChangeUserPrivacySettings(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewUserPostgresRepository(db, metrics.NewMockMetrics())

	username := "testuser"
	userID := int64(1)

	privacy := &repoModel.PrivacySettings{
		IsPublicPlaylists:       true,
		IsPublicMinutesListened: false,
		IsPublicFavoriteArtists: true,
		IsPublicTracksListened:  false,
		IsPublicFavoriteTracks:  true,
		IsPublicArtistsListened: false,
	}

	idRows := sqlmock.NewRows([]string{"id"}).AddRow(userID)
	mock.ExpectQuery("SELECT id").WithArgs(username).WillReturnRows(idRows)

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

func TestGetFullUserData(t *testing.T) {
    db, mock, ctx := setupTest(t)
    defer db.Close()

    repo := NewUserPostgresRepository(db, metrics.NewMockMetrics())

    username := "testuser"
    userID := int64(1)

    idRows := sqlmock.NewRows([]string{"id"}).AddRow(userID)
    mock.ExpectQuery("SELECT id").WithArgs(username).WillReturnRows(idRows)

    privacyRows := sqlmock.NewRows([]string{
        "is_public_playlists",
        "is_public_minutes_listened",
        "is_public_favorite_artists",
        "is_public_tracks_listened",
        "is_public_favorite_tracks",
        "is_public_artists_listened",
    }).AddRow(true, false, true, false, true, false)
    mock.ExpectQuery("SELECT is_public_playlists").
        WithArgs(userID).
        WillReturnRows(privacyRows)

    userRows := sqlmock.NewRows([]string{"id", "username", "email", "thumbnail_url"}).
        AddRow(userID, "testuser", "test@example.com", "avatar.jpg")
    mock.ExpectQuery("SELECT id, username, email, thumbnail_url").
        WithArgs(userID).
        WillReturnRows(userRows)

    userData, err := repo.GetFullUserData(ctx, username)

    assert.NoError(t, err)
    assert.NotNil(t, userData)
    assert.Equal(t, "testuser", userData.Username)
    assert.Equal(t, "test@example.com", userData.Email)
    assert.Equal(t, "avatar.jpg", userData.Thumbnail)
    assert.True(t, userData.Privacy.IsPublicPlaylists)
    assert.False(t, userData.Privacy.IsPublicMinutesListened)
    assert.True(t, userData.Privacy.IsPublicFavoriteArtists)
    assert.False(t, userData.Privacy.IsPublicTracksListened)
    assert.True(t, userData.Privacy.IsPublicFavoriteTracks)
    assert.False(t, userData.Privacy.IsPublicArtistsListened)
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteUser(t *testing.T) {
    db, mock, ctx := setupTest(t)
    defer db.Close()

    repo := NewUserPostgresRepository(db, metrics.NewMockMetrics())

    userDelete := &repoModel.UserDelete{
        Username: "testuser",
        Email:    "test@example.com",
        Password: "password123",
    }

    idRows := sqlmock.NewRows([]string{"id"}).AddRow(1)
    mock.ExpectQuery("SELECT id").WithArgs(userDelete.Username).WillReturnRows(idRows)

    salt := make([]byte, 8)
    hash := argon2.IDKey([]byte(userDelete.Password), salt, 1, 64*1024, 4, 32)
    combined := append(salt, hash...)
    passwordHash := base64.StdEncoding.EncodeToString(combined)
    
    passRows := sqlmock.NewRows([]string{"password_hash"}).AddRow(passwordHash)
    mock.ExpectQuery("SELECT password_hash").WithArgs(1).WillReturnRows(passRows)

    mock.ExpectExec("DELETE FROM").
        WithArgs(userDelete.Username, userDelete.Email).
        WillReturnResult(sqlmock.NewResult(0, 1))

    err := repo.DeleteUser(ctx, userDelete)

    assert.NoError(t, err)
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestChangeUserData(t *testing.T) {
    db, mock, ctx := setupTest(t)
    defer db.Close()

    repo := NewUserPostgresRepository(db, metrics.NewMockMetrics())

    username := "testuser"
    userID := int64(1)

    changeData := &repoModel.ChangeUserData{
        Password:    "oldpass",
        NewUsername: "newuser",
        NewEmail:    "newemail@example.com",
        NewPassword: "newpass",
    }

    salt := make([]byte, 8)
    hash := argon2.IDKey([]byte(changeData.Password), salt, 1, 64*1024, 4, 32)
    combined := append(salt, hash...)
    passwordHash := base64.StdEncoding.EncodeToString(combined)

    idRows := sqlmock.NewRows([]string{"id"}).AddRow(userID)
    mock.ExpectQuery("SELECT id").WithArgs(username).WillReturnRows(idRows)

    mock.ExpectQuery("SELECT 1").WithArgs(strings.ToLower(changeData.NewUsername), userID).WillReturnError(sql.ErrNoRows)

    mock.ExpectExec("UPDATE \"user\"").
        WithArgs(strings.ToLower(changeData.NewUsername), userID).
        WillReturnResult(sqlmock.NewResult(0, 1))

    mock.ExpectQuery("SELECT 1").WithArgs(changeData.NewEmail, userID).WillReturnError(sql.ErrNoRows)

    mock.ExpectExec("UPDATE \"user\"").
        WithArgs(changeData.NewEmail, userID).
        WillReturnResult(sqlmock.NewResult(0, 1))

    passRows := sqlmock.NewRows([]string{"password_hash"}).AddRow(passwordHash)
    mock.ExpectQuery("SELECT password_hash").WithArgs(userID).WillReturnRows(passRows)

    mock.ExpectExec("UPDATE \"user\"").
        WithArgs(sqlmock.AnyArg(), userID).
        WillReturnResult(sqlmock.NewResult(0, 1))

    err := repo.ChangeUserData(ctx, username, changeData)

    assert.NoError(t, err)
    assert.NoError(t, mock.ExpectationsWereMet())

    mock.ExpectQuery("SELECT id").WithArgs("nonexistent").WillReturnError(sql.ErrNoRows)
    err = repo.ChangeUserData(ctx, "nonexistent", changeData)
    assert.Error(t, err)
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestChangeUserDataOnlyPassword(t *testing.T) {
    db, mock, ctx := setupTest(t)
    defer db.Close()

    repo := NewUserPostgresRepository(db, metrics.NewMockMetrics())

    username := "testuser"
    userID := int64(1)
    changeData := &repoModel.ChangeUserData{
        Password:    "oldpass",
        NewPassword: "newpass",
    }

    salt := make([]byte, 8)
    hash := argon2.IDKey([]byte(changeData.Password), salt, 1, 64*1024, 4, 32)
    combined := append(salt, hash...)
    passwordHash := base64.StdEncoding.EncodeToString(combined)

    idRows := sqlmock.NewRows([]string{"id"}).AddRow(userID)
    mock.ExpectQuery("SELECT id").WithArgs(username).WillReturnRows(idRows)

    passRows := sqlmock.NewRows([]string{"password_hash"}).AddRow(passwordHash)
    mock.ExpectQuery("SELECT password_hash").WithArgs(userID).WillReturnRows(passRows)

    mock.ExpectExec("UPDATE \"user\"").
        WithArgs(sqlmock.AnyArg(), userID).
        WillReturnResult(sqlmock.NewResult(0, 1))

    err := repo.ChangeUserData(ctx, username, changeData)

    assert.NoError(t, err)
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestChangeUserDataOnlyEmail(t *testing.T) {
    db, mock, ctx := setupTest(t)
    defer db.Close()

    repo := NewUserPostgresRepository(db, metrics.NewMockMetrics())

    username := "testuser"
    userID := int64(1)
    
    // Изменяем только email
    changeData := &repoModel.ChangeUserData{
        Password: "oldpass",
        NewEmail: "newemail@example.com",
    }

    salt := make([]byte, 8)
    hash := argon2.IDKey([]byte(changeData.Password), salt, 1, 64*1024, 4, 32)
    combined := append(salt, hash...)
    passwordHash := base64.StdEncoding.EncodeToString(combined)

    idRows := sqlmock.NewRows([]string{"id"}).AddRow(userID)
    mock.ExpectQuery("SELECT id").WithArgs(username).WillReturnRows(idRows)
    
    mock.ExpectQuery("SELECT 1").
        WithArgs(changeData.NewEmail, userID).
        WillReturnError(sql.ErrNoRows)
    
    mock.ExpectExec("UPDATE \"user\"").
        WithArgs(changeData.NewEmail, userID).
        WillReturnResult(sqlmock.NewResult(0, 1))
        
    passRows := sqlmock.NewRows([]string{"password_hash"}).AddRow(passwordHash)
    mock.ExpectQuery("SELECT password_hash").WithArgs(userID).WillReturnRows(passRows)

    err := repo.ChangeUserData(ctx, username, changeData)

    assert.NoError(t, err)
}

func TestChangeUserDataOnlyUsername(t *testing.T) {
    db, mock, ctx := setupTest(t)
    defer db.Close()

    repo := NewUserPostgresRepository(db, metrics.NewMockMetrics())

    username := "testuser"
    userID := int64(1)
    
    // Изменяем только email
    changeData := &repoModel.ChangeUserData{
        Password: "oldpass",
        NewUsername: "newusername",
    }

    salt := make([]byte, 8)
    hash := argon2.IDKey([]byte(changeData.Password), salt, 1, 64*1024, 4, 32)
    combined := append(salt, hash...)
    passwordHash := base64.StdEncoding.EncodeToString(combined)

    idRows := sqlmock.NewRows([]string{"id"}).AddRow(userID)
    mock.ExpectQuery("SELECT id").WithArgs(username).WillReturnRows(idRows)
    
    mock.ExpectQuery("SELECT 1").
        WithArgs(changeData.NewUsername, userID).
        WillReturnError(sql.ErrNoRows)
    
    mock.ExpectExec("UPDATE \"user\"").
        WithArgs(changeData.NewUsername, userID).
        WillReturnResult(sqlmock.NewResult(0, 1))
        
    passRows := sqlmock.NewRows([]string{"password_hash"}).AddRow(passwordHash)
    mock.ExpectQuery("SELECT password_hash").WithArgs(userID).WillReturnRows(passRows)

    err := repo.ChangeUserData(ctx, username, changeData)

    assert.NoError(t, err)
}