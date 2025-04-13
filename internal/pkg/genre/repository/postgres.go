package repository

import (
	"database/sql"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/genre"
	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/repository"
)

const (
	GetGenresByAlbumIDQuery = `
		SELECT g.id, g.name
		FROM genre g
		INNER JOIN album_genre ag ON g.id = ag.genre_id
		WHERE ag.album_id = $1
		ORDER BY g.id
	`
)

type genrePostgresRepository struct {
	db *sql.DB
}

func NewGenrePostgresRepository(db *sql.DB) genre.Repository {
	return &genrePostgresRepository{db: db}
}

func (r *genrePostgresRepository) GetGenresByAlbumID(albumID int64) ([]*repoModel.Genre, error) {
	rows, err := r.db.Query(GetGenresByAlbumIDQuery, albumID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	genres := make([]*repoModel.Genre, 0)
	for rows.Next() {
		var genre repoModel.Genre
		err = rows.Scan(&genre.ID, &genre.Name)
		if err != nil {
			return nil, err
		}
		genres = append(genres, &genre)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return genres, nil
}
