package repository

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	artistPkg "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/artist"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers"
	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func setupTest(t *testing.T) (*sql.DB, sqlmock.Sqlmock, artistPkg.Repository) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	repo := NewArtistPostgresRepository(db)
	return db, mock, repo
}

func getTestContext() context.Context {
	logger, _ := zap.NewDevelopment()
	ctx := context.Background()
	return helpers.LoggerToContext(ctx, logger.Sugar())
}

func TestGetAllArtists(t *testing.T) {
	tests := []struct {
		name            string
		filters         *repoModel.ArtistFilters
		expectedArtists []*repoModel.Artist
		mockError       error
	}{
		{
			name: "Success with multiple artists",
			filters: &repoModel.ArtistFilters{
				Pagination: &repoModel.Pagination{
					Limit:  10,
					Offset: 0,
				},
			},
			expectedArtists: []*repoModel.Artist{
				{
					ID:          1,
					Title:       "Artist 1",
					Description: "Description 1",
					Thumbnail:   "thumb1.jpg",
				},
				{
					ID:          2,
					Title:       "Artist 2",
					Description: "Description 2",
					Thumbnail:   "thumb2.jpg",
				},
			},
		},
		{
			name: "Success with empty results",
			filters: &repoModel.ArtistFilters{
				Pagination: &repoModel.Pagination{
					Limit:  10,
					Offset: 100,
				},
			},
			expectedArtists: []*repoModel.Artist{},
		},
		{
			name: "Database error",
			filters: &repoModel.ArtistFilters{
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
				mock.ExpectQuery(regexp.QuoteMeta(GetAllArtistsQuery)).
					WithArgs(tt.filters.Pagination.Limit, tt.filters.Pagination.Offset).
					WillReturnError(tt.mockError)
			} else {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "thumbnail_url"})
				for _, a := range tt.expectedArtists {
					rows.AddRow(a.ID, a.Title, a.Description, a.Thumbnail)
				}
				mock.ExpectQuery(regexp.QuoteMeta(GetAllArtistsQuery)).
					WithArgs(tt.filters.Pagination.Limit, tt.filters.Pagination.Offset).
					WillReturnRows(rows)
			}

			artists, err := repo.GetAllArtists(ctx, tt.filters)

			if tt.mockError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedArtists, artists)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetArtistByID(t *testing.T) {
	tests := []struct {
		name           string
		artistID       int64
		expectedArtist *repoModel.Artist
		mockError      error
		expectedError  error
	}{
		{
			name:     "Success",
			artistID: 1,
			expectedArtist: &repoModel.Artist{
				ID:          1,
				Title:       "Artist 1",
				Description: "Description 1",
				Thumbnail:   "thumb1.jpg",
			},
		},
		{
			name:          "Not Found",
			artistID:      999,
			mockError:     sql.ErrNoRows,
			expectedError: artistPkg.ErrArtistNotFound,
		},
		{
			name:          "Database Error",
			artistID:      2,
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
				mock.ExpectQuery(regexp.QuoteMeta(GetArtistByIDQuery)).
					WithArgs(tt.artistID).
					WillReturnError(tt.mockError)
			} else {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "thumbnail_url"}).
					AddRow(tt.expectedArtist.ID, tt.expectedArtist.Title, tt.expectedArtist.Description, tt.expectedArtist.Thumbnail)
				mock.ExpectQuery(regexp.QuoteMeta(GetArtistByIDQuery)).
					WithArgs(tt.artistID).
					WillReturnRows(rows)
			}

			artist, err := repo.GetArtistByID(ctx, tt.artistID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.expectedError))
				assert.Nil(t, artist)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedArtist, artist)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetArtistTitleByID(t *testing.T) {
	tests := []struct {
		name          string
		artistID      int64
		expectedTitle string
		mockError     error
		expectedError error
	}{
		{
			name:          "Success",
			artistID:      1,
			expectedTitle: "Artist 1",
		},
		{
			name:          "Not Found",
			artistID:      999,
			mockError:     sql.ErrNoRows,
			expectedError: artistPkg.ErrArtistNotFound,
		},
		{
			name:          "Database Error",
			artistID:      2,
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
				mock.ExpectQuery(regexp.QuoteMeta(GetArtistTitleByIDQuery)).
					WithArgs(tt.artistID).
					WillReturnError(tt.mockError)
			} else {
				rows := sqlmock.NewRows([]string{"title"}).
					AddRow(tt.expectedTitle)
				mock.ExpectQuery(regexp.QuoteMeta(GetArtistTitleByIDQuery)).
					WithArgs(tt.artistID).
					WillReturnRows(rows)
			}

			title, err := repo.GetArtistTitleByID(ctx, tt.artistID)

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

func TestGetArtistsByTrackID(t *testing.T) {
	tests := []struct {
		name            string
		trackID         int64
		expectedArtists []*repoModel.ArtistWithRole
		mockError       error
	}{
		{
			name:    "Success with multiple artists",
			trackID: 1,
			expectedArtists: []*repoModel.ArtistWithRole{
				{
					ID:    1,
					Title: "Artist 1",
					Role:  "main",
				},
				{
					ID:    2,
					Title: "Artist 2",
					Role:  "featured",
				},
			},
		},
		{
			name:            "Success with no artists",
			trackID:         999,
			expectedArtists: []*repoModel.ArtistWithRole{},
		},
		{
			name:      "Database error",
			trackID:   2,
			mockError: sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, repo := setupTest(t)
			defer db.Close()
			ctx := getTestContext()

			if tt.mockError != nil {
				mock.ExpectQuery(regexp.QuoteMeta(GetArtistsByTrackIDQuery)).
					WithArgs(tt.trackID).
					WillReturnError(tt.mockError)
			} else {
				rows := sqlmock.NewRows([]string{"id", "title", "role"})
				for _, a := range tt.expectedArtists {
					rows.AddRow(a.ID, a.Title, a.Role)
				}
				mock.ExpectQuery(regexp.QuoteMeta(GetArtistsByTrackIDQuery)).
					WithArgs(tt.trackID).
					WillReturnRows(rows)
			}

			artists, err := repo.GetArtistsByTrackID(ctx, tt.trackID)

			if tt.mockError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedArtists, artists)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetArtistsByTrackIDs(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewArtistPostgresRepository(db)

	ctx := getTestContext()

	tests := []struct {
		name            string
		trackIDs        []int64
		expectedArtists map[int64][]*repoModel.ArtistWithRole
		mockError       error
		setupMock       func(mock sqlmock.Sqlmock, trackIDs []int64, artists map[int64][]*repoModel.ArtistWithRole, mockError error)
	}{
		{
			name:     "Success with multiple artists",
			trackIDs: []int64{1, 2},
			expectedArtists: map[int64][]*repoModel.ArtistWithRole{
				1: {
					{
						ID:    1,
						Title: "Artist 1",
						Role:  "main",
					},
				},
				2: {
					{
						ID:    2,
						Title: "Artist 2",
						Role:  "featured",
					},
				},
			},
			setupMock: func(mock sqlmock.Sqlmock, trackIDs []int64, artists map[int64][]*repoModel.ArtistWithRole, mockError error) {
				queryPattern := regexp.QuoteMeta(GetArtistsByTrackIDsQuery)

				if mockError != nil {
					mock.ExpectQuery(queryPattern).WillReturnError(mockError)
					return
				}

				rows := sqlmock.NewRows([]string{"id", "title", "role", "track_id"})
				for trackID, trackArtists := range artists {
					for _, artist := range trackArtists {
						rows.AddRow(artist.ID, artist.Title, artist.Role, trackID)
					}
				}

				mock.ExpectQuery(queryPattern).WillReturnRows(rows)
			},
		},
		{
			name:            "Success with no artists",
			trackIDs:        []int64{999},
			expectedArtists: map[int64][]*repoModel.ArtistWithRole{},
			setupMock: func(mock sqlmock.Sqlmock, trackIDs []int64, artists map[int64][]*repoModel.ArtistWithRole, mockError error) {
				queryPattern := regexp.QuoteMeta(GetArtistsByTrackIDsQuery)
				rows := sqlmock.NewRows([]string{"id", "title", "role", "track_id"})
				mock.ExpectQuery(queryPattern).WillReturnRows(rows)
			},
		},
		{
			name:      "Database error",
			trackIDs:  []int64{3, 4},
			mockError: sql.ErrConnDone,
			setupMock: func(mock sqlmock.Sqlmock, trackIDs []int64, artists map[int64][]*repoModel.ArtistWithRole, mockError error) {
				queryPattern := regexp.QuoteMeta(GetArtistsByTrackIDsQuery)
				mock.ExpectQuery(queryPattern).WillReturnError(mockError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectationsWereMet()

			tt.setupMock(mock, tt.trackIDs, tt.expectedArtists, tt.mockError)

			artists, err := repo.GetArtistsByTrackIDs(ctx, tt.trackIDs)

			if tt.mockError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if len(tt.expectedArtists) == 0 {
					assert.Empty(t, artists)
				} else {
					assert.Equal(t, tt.expectedArtists, artists)
				}
			}
		})
	}
}

func TestGetArtistStats(t *testing.T) {
	tests := []struct {
		name          string
		artistID      int64
		expectedStats *repoModel.ArtistStats
		mockError     error
	}{
		{
			name:     "Success",
			artistID: 1,
			expectedStats: &repoModel.ArtistStats{
				ListenersCount: 1000,
				FavoritesCount: 500,
			},
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
				mock.ExpectQuery(regexp.QuoteMeta(GetArtistStatsQuery)).
					WithArgs(tt.artistID).
					WillReturnError(tt.mockError)
			} else {
				rows := sqlmock.NewRows([]string{"listeners_count", "favorites_count"}).
					AddRow(tt.expectedStats.ListenersCount, tt.expectedStats.FavoritesCount)
				mock.ExpectQuery(regexp.QuoteMeta(GetArtistStatsQuery)).
					WithArgs(tt.artistID).
					WillReturnRows(rows)
			}

			stats, err := repo.GetArtistStats(ctx, tt.artistID)

			if tt.mockError != nil {
				assert.Error(t, err)
				assert.Nil(t, stats)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStats, stats)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetArtistsByAlbumID(t *testing.T) {
	tests := []struct {
		name            string
		albumID         int64
		expectedArtists []*repoModel.ArtistWithTitle
		mockError       error
	}{
		{
			name:    "Success with multiple artists",
			albumID: 1,
			expectedArtists: []*repoModel.ArtistWithTitle{
				{
					ID:    1,
					Title: "Artist 1",
				},
				{
					ID:    2,
					Title: "Artist 2",
				},
			},
		},
		{
			name:            "Success with no artists",
			albumID:         999,
			expectedArtists: []*repoModel.ArtistWithTitle{},
		},
		{
			name:      "Database error",
			albumID:   2,
			mockError: sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, repo := setupTest(t)
			defer db.Close()
			ctx := getTestContext()

			if tt.mockError != nil {
				mock.ExpectQuery(regexp.QuoteMeta(GetArtistsByAlbumIDQuery)).
					WithArgs(tt.albumID).
					WillReturnError(tt.mockError)
			} else {
				rows := sqlmock.NewRows([]string{"id", "title"})
				for _, a := range tt.expectedArtists {
					rows.AddRow(a.ID, a.Title)
				}
				mock.ExpectQuery(regexp.QuoteMeta(GetArtistsByAlbumIDQuery)).
					WithArgs(tt.albumID).
					WillReturnRows(rows)
			}

			artists, err := repo.GetArtistsByAlbumID(ctx, tt.albumID)

			if tt.mockError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedArtists, artists)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetArtistsByAlbumIDs(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewArtistPostgresRepository(db)

	ctx := getTestContext()

	tests := []struct {
		name            string
		albumIDs        []int64
		expectedArtists map[int64][]*repoModel.ArtistWithTitle
		mockError       error
		setupMock       func(mock sqlmock.Sqlmock, albumIDs []int64, artists map[int64][]*repoModel.ArtistWithTitle, mockError error)
	}{
		{
			name:     "Success with multiple artists",
			albumIDs: []int64{1, 2},
			expectedArtists: map[int64][]*repoModel.ArtistWithTitle{
				1: {
					{
						ID:    1,
						Title: "Artist 1",
					},
				},
				2: {
					{
						ID:    2,
						Title: "Artist 2",
					},
				},
			},
			setupMock: func(mock sqlmock.Sqlmock, albumIDs []int64, artists map[int64][]*repoModel.ArtistWithTitle, mockError error) {
				queryPattern := regexp.QuoteMeta(GetArtistsByAlbumIDsQuery)

				if mockError != nil {
					mock.ExpectQuery(queryPattern).WillReturnError(mockError)
					return
				}

				rows := sqlmock.NewRows([]string{"id", "title", "album_id"})
				for albumID, albumArtists := range artists {
					for _, artist := range albumArtists {
						rows.AddRow(artist.ID, artist.Title, albumID)
					}
				}

				mock.ExpectQuery(queryPattern).WillReturnRows(rows)
			},
		},
		{
			name:            "Success with no artists",
			albumIDs:        []int64{999},
			expectedArtists: map[int64][]*repoModel.ArtistWithTitle{},
			setupMock: func(mock sqlmock.Sqlmock, albumIDs []int64, artists map[int64][]*repoModel.ArtistWithTitle, mockError error) {
				queryPattern := regexp.QuoteMeta(GetArtistsByAlbumIDsQuery)
				rows := sqlmock.NewRows([]string{"id", "title", "album_id"})
				mock.ExpectQuery(queryPattern).WillReturnRows(rows)
			},
		},
		{
			name:      "Database error",
			albumIDs:  []int64{3, 4},
			mockError: sql.ErrConnDone,
			setupMock: func(mock sqlmock.Sqlmock, albumIDs []int64, artists map[int64][]*repoModel.ArtistWithTitle, mockError error) {
				queryPattern := regexp.QuoteMeta(GetArtistsByAlbumIDsQuery)
				mock.ExpectQuery(queryPattern).WillReturnError(mockError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectationsWereMet()

			tt.setupMock(mock, tt.albumIDs, tt.expectedArtists, tt.mockError)

			artists, err := repo.GetArtistsByAlbumIDs(ctx, tt.albumIDs)

			if tt.mockError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if len(tt.expectedArtists) == 0 {
					assert.Empty(t, artists)
				} else {
					assert.Equal(t, tt.expectedArtists, artists)
				}
			}
		})
	}
}
