package track

import (
	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/repository"
)

type Repository interface {
	GetAllTracks(filters *repoModel.TrackFilters) ([]*repoModel.Track, error)
}
