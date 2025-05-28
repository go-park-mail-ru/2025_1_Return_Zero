package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	loggerPkg "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func setupTest(t *testing.T) (*sql.DB, sqlmock.Sqlmock, context.Context) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	// Create a test logger that doesn't sync to stderr to avoid sync errors in tests
	config := zap.NewDevelopmentConfig()
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}
	logger, err := config.Build()
	if err != nil {
		// Fallback to NewNop if config fails
		logger = zap.NewNop()
	}

	ctx := loggerPkg.LoggerToContext(context.Background(), logger.Sugar())

	return db, mock, ctx
}

func TestCreateLabel(t *testing.T) {
    db, mock, ctx := setupTest(t)
    defer db.Close()

    repo := NewLabelPostgresRepository(db)
    
    rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
    
    mock.ExpectPrepare("INSERT INTO label").
        ExpectQuery().WithArgs("test_label").
        WillReturnRows(rows)
    
    id, err := repo.CreateLabel(ctx, "test_label")
    require.NoError(t, err)
    require.Equal(t, int64(1), id)
    
    err = mock.ExpectationsWereMet()
    require.NoError(t, err)
}

func TestCreateLabelError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewLabelPostgresRepository(db)

	mock.ExpectPrepare("INSERT INTO label").
		ExpectQuery().WithArgs("test_label").
		WillReturnError(sql.ErrConnDone)

	id, err := repo.CreateLabel(ctx, "test_label")
	require.Error(t, err)
	require.Equal(t, int64(-1), id)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestGetLabel(t *testing.T) {
    db, mock, ctx := setupTest(t)
    defer db.Close()

    repo := NewLabelPostgresRepository(db)

    rows := sqlmock.NewRows([]string{"name"}).AddRow("test_label")

    mock.ExpectPrepare("SELECT name FROM label").
        ExpectQuery().WithArgs(1).
        WillReturnRows(rows)

    label, err := repo.GetLabel(ctx, 1)
    require.NoError(t, err)
    require.NotNil(t, label)
    require.Equal(t, "test_label", label.Name)

    err = mock.ExpectationsWereMet()
    require.NoError(t, err)
}

func TestGetLabelNotFound(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewLabelPostgresRepository(db)

	mock.ExpectPrepare("SELECT name FROM label").
		ExpectQuery().WithArgs(1).
		WillReturnError(sql.ErrNoRows)

	label, err := repo.GetLabel(ctx, 1)
	require.Error(t, err)
	require.Nil(t, label)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestGetLabelError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewLabelPostgresRepository(db)

	mock.ExpectPrepare("SELECT name FROM label").
		ExpectQuery().WithArgs(1).
		WillReturnError(sql.ErrConnDone)

	label, err := repo.GetLabel(ctx, 1)
	require.Error(t, err)
	require.Nil(t, label)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestCheckIsLabelUnique(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewLabelPostgresRepository(db)

	mock.ExpectPrepare("SELECT 1 FROM label WHERE name = ?").
		ExpectQuery().WithArgs("unique_label").
		WillReturnRows(sqlmock.NewRows([]string{"1"}).AddRow(1))

	exists, err := repo.CheckIsLabelUnique(ctx, "unique_label")
	require.NoError(t, err)
	require.True(t, exists)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestCheckIsLabelUniqueNotFound(t *testing.T) {
    db, mock, ctx := setupTest(t)
    defer db.Close()

    repo := NewLabelPostgresRepository(db)

    mock.ExpectPrepare("SELECT 1 FROM label WHERE name = ?").
        ExpectQuery().WithArgs("non_unique_label").
        WillReturnError(sql.ErrNoRows)

    exists, err := repo.CheckIsLabelUnique(ctx, "non_unique_label")
    require.NoError(t, err)  
    require.False(t, exists)

    err = mock.ExpectationsWereMet()
    require.NoError(t, err)
}

func TestCheckIsLabelUniqueError(t *testing.T) {
	db, mock, ctx := setupTest(t)
	defer db.Close()

	repo := NewLabelPostgresRepository(db)

	mock.ExpectPrepare("SELECT 1 FROM label WHERE name = ?").
		ExpectQuery().WithArgs("error_label").
		WillReturnError(sql.ErrConnDone)

	exists, err := repo.CheckIsLabelUnique(ctx, "error_label")
	require.Error(t, err)
	require.False(t, exists)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}