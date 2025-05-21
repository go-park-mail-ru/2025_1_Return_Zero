package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/label/domain"
	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/repository"

	loggerPkg "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
	"go.uber.org/zap"
)

const (
	CreateLabelQuery = `
			INSERT INTO label (name)
			VALUES ($1)
			RETURNING id
	`
	GetLabelByIdQuery = `
			SELECT name
			FROM label
			WHERE id = $1
	`
	CheckIsLabelUnique = `
			SELECT 1
			FROM label
			WHERE name = $1
	`
)

type labelPostgresRepository struct {
	db *sql.DB
}

func NewLabelPostgresRepository(db *sql.DB) domain.Repository {
	return &labelPostgresRepository{
		db: db,
	}
}

func (r *labelPostgresRepository) CreateLabel(ctx context.Context, name string) (int64, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Creating label", zap.String("name", name))

	stmt, err := r.db.PrepareContext(ctx, CreateLabelQuery)
	if err != nil {
		logger.Error("failed to prepare statement", zap.Error(err))
		return -1, err
	}
	defer stmt.Close()

	var labelID int64
	err = stmt.QueryRowContext(ctx, name).Scan(&labelID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error("label not found", zap.Error(err))
			return -1, err
		}
		logger.Error("failed to create label", zap.Error(err))
		return -1, err
	}

	return labelID, nil
}

func (r *labelPostgresRepository) GetLabel(ctx context.Context, labelID int64) (*repoModel.Label, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Getting label", zap.Int64("labelID", labelID))

	stmt, err := r.db.PrepareContext(ctx, GetLabelByIdQuery)
	if err != nil {
		logger.Error("failed to prepare statement", zap.Error(err))
		return nil, err
	}
	defer stmt.Close()

	var label repoModel.Label
	err = stmt.QueryRowContext(ctx, labelID).Scan(&label.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error("label not found", zap.Error(err))
			return nil, err
		}
		logger.Error("failed to get label", zap.Error(err))
		return nil, err
	}

	return &label, nil
}

func (r *labelPostgresRepository) CheckIsLabelUnique(ctx context.Context, labelName string) (bool, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Checking if label is unique", zap.String("labelName", labelName))
	stmt, err := r.db.PrepareContext(ctx, CheckIsLabelUnique)
	if err != nil {
		logger.Error("failed to prepare statement", zap.Error(err))
		return false, err
	}
	defer stmt.Close()

	var exist bool
	err = stmt.QueryRowContext(ctx, labelName).Scan(&exist)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		logger.Error("failed to check if artist name exists", zap.Error(err))
		return false, err
	}
	return exist, nil
}
