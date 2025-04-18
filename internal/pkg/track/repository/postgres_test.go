package repository

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers"
	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/repository"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/track"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func setupTest(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *TrackPostgresRepository) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	repo := NewTrackPostgresRepository(db)
	return db, mock, repo
}

func getTestContext() context.Context {
	logger := zap.NewNop()
	ctx := context.Background()
	return helpers.LoggerToContext(ctx, logger.Sugar())
}

func TestGetAllTracks(t *testing.T) {
	tests := []struct {
		name           string
		filters        *repoModel.TrackFilters
		expectedTracks []*repoModel.Track
		mockError      error
	}{
		{
			name: "Success with multiple tracks",
			filters: &repoModel.TrackFilters{
				Pagination: &repoModel.Pagination{
					Limit:  10,
					Offset: 0,
				},
			},
			expectedTracks: []*repoModel.Track{
				{
					ID:        1,
					Title:     "Track 1",
					Thumbnail: "thumb1.jpg",
					Duration:  180,
					AlbumID:   1,
				},
				{
					ID:        2,
					Title:     "Track 2",
					Thumbnail: "thumb2.jpg",
					Duration:  240,
					AlbumID:   2,
				},
			},
		},
		{
			name: "Success with empty results",
			filters: &repoModel.TrackFilters{
				Pagination: &repoModel.Pagination{
					Limit:  10,
					Offset: 100,
				},
			},
			expectedTracks: []*repoModel.Track{},
		},
		{
			name: "Database error",
			filters: &repoModel.TrackFilters{
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
				mock.ExpectQuery(regexp.QuoteMeta(GetAllTracksQuery)).
					WithArgs(tt.filters.Pagination.Limit, tt.filters.Pagination.Offset).
					WillReturnError(tt.mockError)
			} else {
				rows := sqlmock.NewRows([]string{"id", "title", "thumbnail_url", "duration", "album_id"})
				for _, tr := range tt.expectedTracks {
					rows.AddRow(tr.ID, tr.Title, tr.Thumbnail, tr.Duration, tr.AlbumID)
				}
				mock.ExpectQuery(regexp.QuoteMeta(GetAllTracksQuery)).
					WithArgs(tt.filters.Pagination.Limit, tt.filters.Pagination.Offset).
					WillReturnRows(rows)
			}

			tracks, err := repo.GetAllTracks(ctx, tt.filters)

			if tt.mockError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedTracks, tracks)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetTrackByID(t *testing.T) {
	tests := []struct {
		name          string
		trackID       int64
		expectedTrack *repoModel.TrackWithFileKey
		mockError     error
		expectedError error
	}{
		{
			name:    "Success",
			trackID: 1,
			expectedTrack: &repoModel.TrackWithFileKey{
				Track: repoModel.Track{
					ID:        1,
					Title:     "Track 1",
					Thumbnail: "thumb1.jpg",
					Duration:  180,
					AlbumID:   1,
				},
				FileKey: "file1.mp3",
			},
		},
		{
			name:          "Not Found",
			trackID:       999,
			mockError:     sql.ErrNoRows,
			expectedError: track.ErrTrackNotFound,
		},
		{
			name:          "Database Error",
			trackID:       2,
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
				mock.ExpectQuery(regexp.QuoteMeta(GetTrackByIDQuery)).
					WithArgs(tt.trackID).
					WillReturnError(tt.mockError)
			} else {
				rows := sqlmock.NewRows([]string{"id", "title", "thumbnail_url", "duration", "album_id", "file_url"}).
					AddRow(tt.expectedTrack.Track.ID, tt.expectedTrack.Track.Title, tt.expectedTrack.Track.Thumbnail,
						tt.expectedTrack.Track.Duration, tt.expectedTrack.Track.AlbumID, tt.expectedTrack.FileKey)
				mock.ExpectQuery(regexp.QuoteMeta(GetTrackByIDQuery)).
					WithArgs(tt.trackID).
					WillReturnRows(rows)
			}

			track, err := repo.GetTrackByID(ctx, tt.trackID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.expectedError))
				assert.Nil(t, track)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedTrack, track)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetTracksByArtistID(t *testing.T) {
	tests := []struct {
		name           string
		artistID       int64
		expectedTracks []*repoModel.Track
		mockError      error
	}{
		{
			name:     "Success with multiple tracks",
			artistID: 1,
			expectedTracks: []*repoModel.Track{
				{
					ID:        1,
					Title:     "Track 1",
					Thumbnail: "thumb1.jpg",
					Duration:  180,
					AlbumID:   1,
				},
				{
					ID:        2,
					Title:     "Track 2",
					Thumbnail: "thumb2.jpg",
					Duration:  240,
					AlbumID:   2,
				},
			},
		},
		{
			name:           "Success with no tracks",
			artistID:       999,
			expectedTracks: []*repoModel.Track{},
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
				mock.ExpectQuery(regexp.QuoteMeta(GetTracksByArtistIDQuery)).
					WithArgs(tt.artistID).
					WillReturnError(tt.mockError)
			} else {
				rows := sqlmock.NewRows([]string{"id", "title", "thumbnail_url", "duration", "album_id"})
				for _, tr := range tt.expectedTracks {
					rows.AddRow(tr.ID, tr.Title, tr.Thumbnail, tr.Duration, tr.AlbumID)
				}
				mock.ExpectQuery(regexp.QuoteMeta(GetTracksByArtistIDQuery)).
					WithArgs(tt.artistID).
					WillReturnRows(rows)
			}

			tracks, err := repo.GetTracksByArtistID(ctx, tt.artistID, &repoModel.TrackFilters{
				Pagination: &repoModel.Pagination{
					Limit:  10,
					Offset: 0,
				},
			})

			if tt.mockError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedTracks, tracks)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestCreateStream(t *testing.T) {
	tests := []struct {
		name       string
		createData *repoModel.TrackStreamCreateData
		expectedID int64
		mockError  error
	}{
		{
			name: "Success",
			createData: &repoModel.TrackStreamCreateData{
				TrackID: 1,
				UserID:  1,
			},
			expectedID: 1,
		},
		{
			name: "Database error",
			createData: &repoModel.TrackStreamCreateData{
				TrackID: 2,
				UserID:  2,
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
				mock.ExpectQuery(regexp.QuoteMeta(CreateStreamQuery)).
					WithArgs(tt.createData.TrackID, tt.createData.UserID).
					WillReturnError(tt.mockError)
			} else {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(tt.expectedID)
				mock.ExpectQuery(regexp.QuoteMeta(CreateStreamQuery)).
					WithArgs(tt.createData.TrackID, tt.createData.UserID).
					WillReturnRows(rows)
			}

			id, err := repo.CreateStream(ctx, tt.createData)

			if tt.mockError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, id)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetStreamByID(t *testing.T) {
	tests := []struct {
		name           string
		streamID       int64
		expectedStream *repoModel.TrackStream
		mockError      error
		expectedError  error
	}{
		{
			name:     "Success",
			streamID: 1,
			expectedStream: &repoModel.TrackStream{
				ID:       1,
				UserID:   1,
				TrackID:  1,
				Duration: 180,
			},
		},
		{
			name:          "Not Found",
			streamID:      999,
			mockError:     sql.ErrNoRows,
			expectedError: track.ErrStreamNotFound,
		},
		{
			name:          "Database Error",
			streamID:      2,
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
				mock.ExpectQuery(regexp.QuoteMeta(GetStreamByIDQuery)).
					WithArgs(tt.streamID).
					WillReturnError(tt.mockError)
			} else {
				rows := sqlmock.NewRows([]string{"id", "user_id", "track_id", "duration"}).
					AddRow(tt.expectedStream.ID, tt.expectedStream.UserID, tt.expectedStream.TrackID, tt.expectedStream.Duration)
				mock.ExpectQuery(regexp.QuoteMeta(GetStreamByIDQuery)).
					WithArgs(tt.streamID).
					WillReturnRows(rows)
			}

			stream, err := repo.GetStreamByID(ctx, tt.streamID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.expectedError))
				assert.Nil(t, stream)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStream, stream)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUpdateStreamDuration(t *testing.T) {
	tests := []struct {
		name          string
		updateData    *repoModel.TrackStreamUpdateData
		mockResult    int64
		mockError     error
		expectedError error
	}{
		{
			name: "Success",
			updateData: &repoModel.TrackStreamUpdateData{
				StreamID: 1,
				Duration: 180,
			},
			mockResult: 1,
		},
		{
			name: "Not Found",
			updateData: &repoModel.TrackStreamUpdateData{
				StreamID: 999,
				Duration: 180,
			},
			mockResult:    0,
			expectedError: track.ErrFailedToUpdateStreamDuration,
		},
		{
			name: "Database Error",
			updateData: &repoModel.TrackStreamUpdateData{
				StreamID: 2,
				Duration: 180,
			},
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
				mock.ExpectExec(regexp.QuoteMeta(UpdateStreamDurationQuery)).
					WithArgs(tt.updateData.Duration, tt.updateData.StreamID).
					WillReturnError(tt.mockError)
			} else {
				mock.ExpectExec(regexp.QuoteMeta(UpdateStreamDurationQuery)).
					WithArgs(tt.updateData.Duration, tt.updateData.StreamID).
					WillReturnResult(sqlmock.NewResult(1, tt.mockResult))
			}

			err := repo.UpdateStreamDuration(ctx, tt.updateData)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.expectedError))
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetStreamsByUserID(t *testing.T) {
	tests := []struct {
		name            string
		userID          int64
		filters         *repoModel.TrackFilters
		expectedStreams []*repoModel.TrackStream
		mockError       error
	}{
		{
			name:   "Success with multiple streams",
			userID: 1,
			filters: &repoModel.TrackFilters{
				Pagination: &repoModel.Pagination{
					Limit:  10,
					Offset: 0,
				},
			},
			expectedStreams: []*repoModel.TrackStream{
				{
					ID:       1,
					UserID:   1,
					TrackID:  1,
					Duration: 180,
				},
				{
					ID:       2,
					UserID:   1,
					TrackID:  2,
					Duration: 240,
				},
			},
		},
		{
			name:   "Success with no streams",
			userID: 999,
			filters: &repoModel.TrackFilters{
				Pagination: &repoModel.Pagination{
					Limit:  10,
					Offset: 0,
				},
			},
			expectedStreams: []*repoModel.TrackStream{},
		},
		{
			name:   "Database error",
			userID: 2,
			filters: &repoModel.TrackFilters{
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
				mock.ExpectQuery(regexp.QuoteMeta(GetStreamsByUserIDQuery)).
					WithArgs(tt.userID, tt.filters.Pagination.Limit, tt.filters.Pagination.Offset).
					WillReturnError(tt.mockError)
			} else {
				rows := sqlmock.NewRows([]string{"id", "user_id", "track_id", "duration"})
				for _, s := range tt.expectedStreams {
					rows.AddRow(s.ID, s.UserID, s.TrackID, s.Duration)
				}
				mock.ExpectQuery(regexp.QuoteMeta(GetStreamsByUserIDQuery)).
					WithArgs(tt.userID, tt.filters.Pagination.Limit, tt.filters.Pagination.Offset).
					WillReturnRows(rows)
			}

			streams, err := repo.GetStreamsByUserID(ctx, tt.userID, tt.filters)

			if tt.mockError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStreams, streams)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetTracksByIDs(t *testing.T) {
	tests := []struct {
		name           string
		trackIDs       []int64
		expectedTracks map[int64]*repoModel.Track
		mockError      error
	}{
		{
			name:     "Success with multiple tracks",
			trackIDs: []int64{1, 2},
			expectedTracks: map[int64]*repoModel.Track{
				1: {
					ID:        1,
					Title:     "Track 1",
					Thumbnail: "thumb1.jpg",
					Duration:  180,
					AlbumID:   1,
				},
				2: {
					ID:        2,
					Title:     "Track 2",
					Thumbnail: "thumb2.jpg",
					Duration:  240,
					AlbumID:   2,
				},
			},
		},
		{
			name:           "Success with no tracks",
			trackIDs:       []int64{999},
			expectedTracks: map[int64]*repoModel.Track{},
		},
		{
			name:      "Database error",
			trackIDs:  []int64{3, 4},
			mockError: sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, repo := setupTest(t)
			defer db.Close()
			ctx := getTestContext()

			if tt.mockError != nil {
				mock.ExpectQuery(regexp.QuoteMeta(GetTracksByIDsQuery)).
					WithArgs(pq.Array(tt.trackIDs)).
					WillReturnError(tt.mockError)
			} else {
				rows := sqlmock.NewRows([]string{"id", "title", "thumbnail_url", "duration", "album_id"})
				for _, track := range tt.expectedTracks {
					rows.AddRow(track.ID, track.Title, track.Thumbnail, track.Duration, track.AlbumID)
				}
				mock.ExpectQuery(regexp.QuoteMeta(GetTracksByIDsQuery)).
					WithArgs(pq.Array(tt.trackIDs)).
					WillReturnRows(rows)
			}

			tracks, err := repo.GetTracksByIDs(ctx, tt.trackIDs)

			if tt.mockError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedTracks, tracks)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
