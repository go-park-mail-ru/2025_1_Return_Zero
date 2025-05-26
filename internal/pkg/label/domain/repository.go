package domain

import (
	"context"

	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/repository"
)

type Repository interface {
	CreateLabel(ctx context.Context, name string) (int64, error)
	GetLabel(ctx context.Context, labelID int64) (*repoModel.Label, error)
	CheckIsLabelUnique(ctx context.Context, labelName string) (bool, error)
}

type S3Repository interface {
	
}