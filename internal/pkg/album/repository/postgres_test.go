package repository

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	albumPkg "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/album"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers"
	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/repository"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func setupTest(t *testing.T) (*sql.DB, sqlmock.Sqlmock, albumPkg.Repository) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	repo := NewAlbumPostgresRepository(db)
	return db, mock, repo
}

func getTestContext() context.Context {
	logger, _ := zap.NewDevelopment()
	ctx := context.Background()
	return helpers.LoggerToContext(ctx, logger.Sugar())
}

func TestGetAllAlbums(t *testing.T) {
	releaseDate := time.Now()
	tests := []struct {
		name           string
		filters        *repoModel.AlbumFilters
		expectedAlbums []*repoModel.Album
		mockError      error
	}{
		{
			name: "Success with multiple albums",
			filters: &repoModel.AlbumFilters{
				Pagination: &repoModel.Pagination{
					Limit:  10,
					Offset: 0,
				},
			},
			expectedAlbums: []*repoModel.Album{
				{
					ID:          1,
					Title:       "Album 1",
					Type:        "album",
					Thumbnail:   "thumb1.jpg",
					ReleaseDate: releaseDate,
				},
				{
					ID:          2,
					Title:       "Album 2",
					Type:        "single",
					Thumbnail:   "thumb2.jpg",
					ReleaseDate: releaseDate,
				},
			},
		},
		{
			name: "Success with empty results",
			filters: &repoModel.AlbumFilters{
				Pagination: &repoModel.Pagination{
					Limit:  10,
					Offset: 100,
				},
			},
			expectedAlbums: []*repoModel.Album{},
		},
		{
			name: "Database error",
			filters: &repoModel.AlbumFilters{
				Pagination: &repoModel.Pagination{
					Limit:  10,
					Offset: 0,
				},
			},
			mockError: sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, repo := setupTest(t)
			defer db.Close()
			ctx := getTestContext()

			if tt.mockError != nil {
				mock.ExpectQuery(regexp.QuoteMeta(GetAllAlbumsQuery)).
					WithArgs(tt.filters.Pagination.Limit, tt.filters.Pagination.Offset).
					WillReturnError(tt.mockError)
			} else {
				rows := sqlmock.NewRows([]string{"id", "title", "type", "thumbnail_url", "release_date"})
				for _, a := range tt.expectedAlbums {
					rows.AddRow(a.ID, a.Title, a.Type, a.Thumbnail, a.ReleaseDate)
				}
				mock.ExpectQuery(regexp.QuoteMeta(GetAllAlbumsQuery)).
					WithArgs(tt.filters.Pagination.Limit, tt.filters.Pagination.Offset).
					WillReturnRows(rows)
			}

			albums, err := repo.GetAllAlbums(ctx, tt.filters)

			if tt.mockError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedAlbums, albums)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetAlbumByID(t *testing.T) {
	releaseDate := time.Now()
	tests := []struct {
		name          string
		albumID       int64
		expectedAlbum *repoModel.Album
		mockError     error
		expectedError error
	}{
		{
			name:    "Success",
			albumID: 1,
			expectedAlbum: &repoModel.Album{
				ID:          1,
				Title:       "Album 1",
				Type:        "album",
				Thumbnail:   "thumb1.jpg",
				ReleaseDate: releaseDate,
			},
		},
		{
			name:          "Not Found",
			albumID:       999,
			mockError:     sql.ErrNoRows,
			expectedError: albumPkg.ErrAlbumNotFound,
		},
		{
			name:          "Database Error",
			albumID:       2,
			mockError:     sql.ErrConnDone,
			expectedError: sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, repo := setupTest(t)
			defer db.Close()
			ctx := getTestContext()

			if tt.mockError != nil {
				mock.ExpectQuery(regexp.QuoteMeta(GetAlbumByIDQuery)).
					WithArgs(tt.albumID).
					WillReturnError(tt.mockError)
			} else {
				rows := sqlmock.NewRows([]string{"id", "title", "type", "thumbnail_url", "release_date"}).
					AddRow(tt.expectedAlbum.ID, tt.expectedAlbum.Title, tt.expectedAlbum.Type, tt.expectedAlbum.Thumbnail, tt.expectedAlbum.ReleaseDate)
				mock.ExpectQuery(regexp.QuoteMeta(GetAlbumByIDQuery)).
					WithArgs(tt.albumID).
					WillReturnRows(rows)
			}

			album, err := repo.GetAlbumByID(ctx, tt.albumID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.expectedError))
				assert.Nil(t, album)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedAlbum, album)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetAlbumTitleByIDs(t *testing.T) {
	tests := []struct {
		name           string
		albumIDs       []int64
		expectedTitles map[int64]string
		mockError      error
	}{
		{
			name:     "Success with multiple albums",
			albumIDs: []int64{1, 2},
			expectedTitles: map[int64]string{
				1: "Album 1",
				2: "Album 2",
			},
		},
		{
			name:           "Success with no albums",
			albumIDs:       []int64{999},
			expectedTitles: map[int64]string{},
		},
		{
			name:      "Database error",
			albumIDs:  []int64{3, 4},
			mockError: sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, repo := setupTest(t)
			defer db.Close()
			ctx := getTestContext()

			if tt.mockError != nil {
				mock.ExpectQuery(regexp.QuoteMeta(GetAlbumTitleByIDsQuery)).
					WithArgs(pq.Array(tt.albumIDs)).
					WillReturnError(tt.mockError)
			} else {
				rows := sqlmock.NewRows([]string{"id", "title"})
				for id, title := range tt.expectedTitles {
					rows.AddRow(id, title)
				}
				mock.ExpectQuery(regexp.QuoteMeta(GetAlbumTitleByIDsQuery)).
					WithArgs(pq.Array(tt.albumIDs)).
					WillReturnRows(rows)
			}

			titles, err := repo.GetAlbumTitleByIDs(ctx, tt.albumIDs)

			if tt.mockError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedTitles, titles)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetAlbumTitleByID(t *testing.T) {
	tests := []struct {
		name          string
		albumID       int64
		expectedTitle string
		mockError     error
		expectedError error
	}{
		{
			name:          "Success",
			albumID:       1,
			expectedTitle: "Album 1",
		},
		{
			name:          "Not Found",
			albumID:       999,
			mockError:     sql.ErrNoRows,
			expectedError: albumPkg.ErrAlbumNotFound,
		},
		{
			name:          "Database Error",
			albumID:       2,
			mockError:     sql.ErrConnDone,
			expectedError: sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, repo := setupTest(t)
			defer db.Close()
			ctx := getTestContext()

			if tt.mockError != nil {
				mock.ExpectQuery(regexp.QuoteMeta(GetAlbumTitleByIDQuery)).
					WithArgs(tt.albumID).
					WillReturnError(tt.mockError)
			} else {
				rows := sqlmock.NewRows([]string{"title"}).
					AddRow(tt.expectedTitle)
				mock.ExpectQuery(regexp.QuoteMeta(GetAlbumTitleByIDQuery)).
					WithArgs(tt.albumID).
					WillReturnRows(rows)
			}

			title, err := repo.GetAlbumTitleByID(ctx, tt.albumID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.expectedError))
				assert.Empty(t, title)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedTitle, title)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetAlbumsByArtistID(t *testing.T) {
	releaseDate := time.Now()
	tests := []struct {
		name           string
		artistID       int64
		expectedAlbums []*repoModel.Album
		mockError      error
	}{
		{
			name:     "Success with multiple albums",
			artistID: 1,
			expectedAlbums: []*repoModel.Album{
				{
					ID:          1,
					Title:       "Album 1",
					Type:        "album",
					Thumbnail:   "thumb1.jpg",
					ReleaseDate: releaseDate,
				},
				{
					ID:          2,
					Title:       "Album 2",
					Type:        "single",
					Thumbnail:   "thumb2.jpg",
					ReleaseDate: releaseDate,
				},
			},
		},
		{
			name:           "Success with no albums",
			artistID:       999,
			expectedAlbums: []*repoModel.Album{},
		},
		{
			name:      "Database error",
			artistID:  2,
			mockError: sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, repo := setupTest(t)
			defer db.Close()
			ctx := getTestContext()

			if tt.mockError != nil {
				mock.ExpectQuery(regexp.QuoteMeta(GetAlbumsByArtistIDQuery)).
					WithArgs(tt.artistID).
					WillReturnError(tt.mockError)
			} else {
				rows := sqlmock.NewRows([]string{"id", "title", "type", "thumbnail_url", "release_date"})
				for _, a := range tt.expectedAlbums {
					rows.AddRow(a.ID, a.Title, a.Type, a.Thumbnail, a.ReleaseDate)
				}
				mock.ExpectQuery(regexp.QuoteMeta(GetAlbumsByArtistIDQuery)).
					WithArgs(tt.artistID).
					WillReturnRows(rows)
			}

			albums, err := repo.GetAlbumsByArtistID(ctx, tt.artistID)

			if tt.mockError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedAlbums, albums)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
