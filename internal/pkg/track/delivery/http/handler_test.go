package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/config"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/delivery"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/track"
	mock_track "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/track/mocks"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func setupTest(t *testing.T) (*mock_track.MockUsecase, *TrackHandler, *httptest.ResponseRecorder) {
	ctrl := gomock.NewController(t)
	mockUsecase := mock_track.NewMockUsecase(ctrl)
	cfg := &config.Config{
		Pagination: delivery.PaginationConfig{
			DefaultLimit:  10,
			DefaultOffset: 0,
			MaxLimit:      100,
		},
	}
	handler := NewTrackHandler(mockUsecase, cfg)
	recorder := httptest.NewRecorder()
	return mockUsecase, handler, recorder
}

func TestTrackHandler_GetAllTracks(t *testing.T) {
	logger := zap.NewNop()
	sugar := logger.Sugar()
	defer logger.Sync()

	tests := []struct {
		name             string
		url              string
		setupMock        func(m *mock_track.MockUsecase)
		expectedStatus   int
		expectedResponse interface{}
	}{
		{
			name: "Success",
			url:  "/tracks?offset=0&limit=10",
			setupMock: func(m *mock_track.MockUsecase) {
				m.EXPECT().
					GetAllTracks(gomock.Any(), &usecaseModel.TrackFilters{
						Pagination: &usecaseModel.Pagination{
							Offset: 0,
							Limit:  10,
						},
					}).Return([]*usecaseModel.Track{
					{
						ID:        1,
						Title:     "Test Track",
						Duration:  200,
						Thumbnail: "test-thumbnail.jpg",
						AlbumID:   1,
						Album:     "Test Album",
						Artists: []*usecaseModel.TrackArtist{
							{
								ID:    1,
								Title: "Test Artist",
								Role:  "main",
							},
						},
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedResponse: &delivery.APIResponse{
				Status: http.StatusOK,
				Body: []*delivery.Track{
					{
						ID:        1,
						Title:     "Test Track",
						Duration:  200,
						Thumbnail: "test-thumbnail.jpg",
						AlbumID:   1,
						Album:     "Test Album",
						Artists: []*delivery.TrackArtist{
							{
								ID:    1,
								Title: "Test Artist",
								Role:  "main",
							},
						},
					},
				},
			},
		},
		{
			name: "Invalid Pagination",
			url:  "/tracks?offset=-1",
			setupMock: func(m *mock_track.MockUsecase) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedResponse: &delivery.APIErrorResponse{
				Status: http.StatusBadRequest,
				Error:  "invalid offset: should be greater than 0",
			},
		},
		{
			name: "Database Error",
			url:  "/tracks",
			setupMock: func(m *mock_track.MockUsecase) {
				m.EXPECT().
					GetAllTracks(gomock.Any(), &usecaseModel.TrackFilters{
						Pagination: &usecaseModel.Pagination{
							Offset: 0,
							Limit:  10,
						},
					}).
					Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedResponse: &delivery.APIErrorResponse{
				Status: http.StatusInternalServerError,
				Error:  "database error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase, handler, recorder := setupTest(t)
			tt.setupMock(mockUsecase)

			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			ctx := helpers.LoggerToContext(req.Context(), sugar)
			req = req.WithContext(ctx)

			handler.GetAllTracks(recorder, req)

			if tt.expectedStatus == http.StatusOK {
				var response delivery.APIResponse
				err := json.NewDecoder(recorder.Body).Decode(&response)
				assert.NoError(t, err)

				expectedResp := tt.expectedResponse.(*delivery.APIResponse)
				assert.Equal(t, expectedResp.Status, response.Status)

				expectedTracks := expectedResp.Body.([]*delivery.Track)
				responseTracks, ok := response.Body.([]interface{})
				assert.True(t, ok)
				assert.Equal(t, len(expectedTracks), len(responseTracks))

				if len(responseTracks) > 0 {
					responseTrack := responseTracks[0].(map[string]interface{})
					expectedTrack := expectedTracks[0]

					assert.Equal(t, float64(expectedTrack.ID), responseTrack["id"])
					assert.Equal(t, expectedTrack.Title, responseTrack["title"])
					assert.Equal(t, float64(expectedTrack.Duration), responseTrack["duration"])
					assert.Equal(t, expectedTrack.Thumbnail, responseTrack["thumbnail_url"])
					assert.Equal(t, expectedTrack.Album, responseTrack["album"])
					assert.Equal(t, float64(expectedTrack.AlbumID), responseTrack["album_id"])

					responseArtists, artistsOk := responseTrack["artists"].([]interface{})
					assert.True(t, artistsOk)
					assert.Equal(t, len(expectedTrack.Artists), len(responseArtists))

					if len(responseArtists) > 0 {
						responseArtist := responseArtists[0].(map[string]interface{})
						expectedArtist := expectedTrack.Artists[0]

						assert.Equal(t, float64(expectedArtist.ID), responseArtist["id"])
						assert.Equal(t, expectedArtist.Title, responseArtist["title"])
						assert.Equal(t, expectedArtist.Role, responseArtist["role"])
					}
				}
			} else {
				var response delivery.APIErrorResponse
				err := json.NewDecoder(recorder.Body).Decode(&response)
				assert.NoError(t, err)

				expectedResp := tt.expectedResponse.(*delivery.APIErrorResponse)
				assert.Equal(t, expectedResp, &response)
			}
		})
	}
}

func TestTrackHandler_GetTrackByID(t *testing.T) {
	logger := zap.NewNop()
	sugar := logger.Sugar()
	defer logger.Sync()

	tests := []struct {
		name             string
		trackID          string
		setupMock        func(m *mock_track.MockUsecase)
		expectedStatus   int
		expectedResponse interface{}
	}{
		{
			name:    "Success",
			trackID: "1",
			setupMock: func(m *mock_track.MockUsecase) {
				m.EXPECT().
					GetTrackByID(gomock.Any(), int64(1)).
					Return(&usecaseModel.TrackDetailed{
						Track: usecaseModel.Track{
							ID:        1,
							Title:     "Test Track",
							Duration:  200,
							Thumbnail: "test-thumbnail.jpg",
							AlbumID:   1,
							Album:     "Test Album",
							Artists: []*usecaseModel.TrackArtist{
								{
									ID:    1,
									Title: "Test Artist",
									Role:  "main",
								},
							},
						},
						FileUrl: "https://example.com/track.mp3",
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedResponse: &delivery.APIResponse{
				Status: http.StatusOK,
				Body: &delivery.TrackDetailed{
					Track: delivery.Track{
						ID:        1,
						Title:     "Test Track",
						Duration:  200,
						Thumbnail: "test-thumbnail.jpg",
						AlbumID:   1,
						Album:     "Test Album",
						Artists: []*delivery.TrackArtist{
							{
								ID:    1,
								Title: "Test Artist",
								Role:  "main",
							},
						},
					},
					FileUrl: "https://example.com/track.mp3",
				},
			},
		},
		{
			name:    "Invalid ID",
			trackID: "invalid",
			setupMock: func(m *mock_track.MockUsecase) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedResponse: &delivery.APIErrorResponse{
				Status: http.StatusBadRequest,
				Error:  "strconv.ParseInt: parsing \"invalid\": invalid syntax",
			},
		},
		{
			name:    "Not Found",
			trackID: "999",
			setupMock: func(m *mock_track.MockUsecase) {
				m.EXPECT().
					GetTrackByID(gomock.Any(), int64(999)).
					Return(nil, track.ErrTrackNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedResponse: &delivery.APIErrorResponse{
				Status: http.StatusNotFound,
				Error:  "track not found",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase, handler, recorder := setupTest(t)
			tt.setupMock(mockUsecase)

			req := httptest.NewRequest(http.MethodGet, "/tracks/"+tt.trackID, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.trackID})
			ctx := helpers.LoggerToContext(req.Context(), sugar)
			req = req.WithContext(ctx)

			handler.GetTrackByID(recorder, req)

			if tt.expectedStatus == http.StatusOK {
				var response delivery.APIResponse
				err := json.NewDecoder(recorder.Body).Decode(&response)
				assert.NoError(t, err)

				expectedResp := tt.expectedResponse.(*delivery.APIResponse)
				assert.Equal(t, expectedResp.Status, response.Status)

				expectedTrack := expectedResp.Body.(*delivery.TrackDetailed)
				responseTrack, ok := response.Body.(map[string]interface{})
				assert.True(t, ok)

				assert.Equal(t, float64(expectedTrack.ID), responseTrack["id"])
				assert.Equal(t, expectedTrack.Title, responseTrack["title"])
				assert.Equal(t, float64(expectedTrack.Duration), responseTrack["duration"])
				assert.Equal(t, expectedTrack.Thumbnail, responseTrack["thumbnail_url"])
				assert.Equal(t, expectedTrack.Album, responseTrack["album"])
				assert.Equal(t, float64(expectedTrack.AlbumID), responseTrack["album_id"])
				assert.Equal(t, expectedTrack.FileUrl, responseTrack["file_url"])

				responseArtists, artistsOk := responseTrack["artists"].([]interface{})
				assert.True(t, artistsOk)
				assert.Equal(t, len(expectedTrack.Artists), len(responseArtists))

				if len(responseArtists) > 0 {
					responseArtist := responseArtists[0].(map[string]interface{})
					expectedArtist := expectedTrack.Artists[0]

					assert.Equal(t, float64(expectedArtist.ID), responseArtist["id"])
					assert.Equal(t, expectedArtist.Title, responseArtist["title"])
					assert.Equal(t, expectedArtist.Role, responseArtist["role"])
				}
			} else {
				var response delivery.APIErrorResponse
				err := json.NewDecoder(recorder.Body).Decode(&response)
				assert.NoError(t, err)

				expectedResp := tt.expectedResponse.(*delivery.APIErrorResponse)
				assert.Equal(t, expectedResp, &response)
			}
		})
	}
}

func TestTrackHandler_GetTracksByArtistID(t *testing.T) {
	logger := zap.NewNop()
	sugar := logger.Sugar()
	defer logger.Sync()

	tests := []struct {
		name             string
		artistID         string
		setupMock        func(m *mock_track.MockUsecase)
		expectedStatus   int
		expectedResponse interface{}
	}{
		{
			name:     "Success",
			artistID: "1",
			setupMock: func(m *mock_track.MockUsecase) {
				m.EXPECT().
					GetTracksByArtistID(gomock.Any(), int64(1)).
					Return([]*usecaseModel.Track{
						{
							ID:        1,
							Title:     "Test Track",
							Duration:  200,
							Thumbnail: "test-thumbnail.jpg",
							AlbumID:   1,
							Album:     "Test Album",
							Artists: []*usecaseModel.TrackArtist{
								{
									ID:    1,
									Title: "Test Artist",
									Role:  "main",
								},
							},
						},
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedResponse: &delivery.APIResponse{
				Status: http.StatusOK,
				Body: []*delivery.Track{
					{
						ID:        1,
						Title:     "Test Track",
						Duration:  200,
						Thumbnail: "test-thumbnail.jpg",
						AlbumID:   1,
						Album:     "Test Album",
						Artists: []*delivery.TrackArtist{
							{
								ID:    1,
								Title: "Test Artist",
								Role:  "main",
							},
						},
					},
				},
			},
		},
		{
			name:     "Invalid ID",
			artistID: "invalid",
			setupMock: func(m *mock_track.MockUsecase) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedResponse: &delivery.APIErrorResponse{
				Status: http.StatusBadRequest,
				Error:  "strconv.ParseInt: parsing \"invalid\": invalid syntax",
			},
		},
		{
			name:     "Not Found",
			artistID: "999",
			setupMock: func(m *mock_track.MockUsecase) {
				m.EXPECT().
					GetTracksByArtistID(gomock.Any(), int64(999)).
					Return(nil, errors.New("artist not found"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedResponse: &delivery.APIErrorResponse{
				Status: http.StatusInternalServerError,
				Error:  "artist not found",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase, handler, recorder := setupTest(t)
			tt.setupMock(mockUsecase)

			req := httptest.NewRequest(http.MethodGet, "/artists/"+tt.artistID+"/tracks", nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.artistID})
			ctx := helpers.LoggerToContext(req.Context(), sugar)
			req = req.WithContext(ctx)

			handler.GetTracksByArtistID(recorder, req)

			if tt.expectedStatus == http.StatusOK {
				var response delivery.APIResponse
				err := json.NewDecoder(recorder.Body).Decode(&response)
				assert.NoError(t, err)

				expectedResp := tt.expectedResponse.(*delivery.APIResponse)
				assert.Equal(t, expectedResp.Status, response.Status)

				expectedTracks := expectedResp.Body.([]*delivery.Track)
				responseTracks, ok := response.Body.([]interface{})
				assert.True(t, ok)
				assert.Equal(t, len(expectedTracks), len(responseTracks))

				if len(responseTracks) > 0 {
					responseTrack := responseTracks[0].(map[string]interface{})
					expectedTrack := expectedTracks[0]

					assert.Equal(t, float64(expectedTrack.ID), responseTrack["id"])
					assert.Equal(t, expectedTrack.Title, responseTrack["title"])
					assert.Equal(t, float64(expectedTrack.Duration), responseTrack["duration"])
					assert.Equal(t, expectedTrack.Thumbnail, responseTrack["thumbnail_url"])
					assert.Equal(t, expectedTrack.Album, responseTrack["album"])
					assert.Equal(t, float64(expectedTrack.AlbumID), responseTrack["album_id"])

					responseArtists, artistsOk := responseTrack["artists"].([]interface{})
					assert.True(t, artistsOk)
					assert.Equal(t, len(expectedTrack.Artists), len(responseArtists))

					if len(responseArtists) > 0 {
						responseArtist := responseArtists[0].(map[string]interface{})
						expectedArtist := expectedTrack.Artists[0]

						assert.Equal(t, float64(expectedArtist.ID), responseArtist["id"])
						assert.Equal(t, expectedArtist.Title, responseArtist["title"])
						assert.Equal(t, expectedArtist.Role, responseArtist["role"])
					}
				}
			} else {
				var response delivery.APIErrorResponse
				err := json.NewDecoder(recorder.Body).Decode(&response)
				assert.NoError(t, err)

				expectedResp := tt.expectedResponse.(*delivery.APIErrorResponse)
				assert.Equal(t, expectedResp, &response)
			}
		})
	}
}

func TestTrackHandler_CreateStream(t *testing.T) {
	logger := zap.NewNop()
	sugar := logger.Sugar()
	defer logger.Sync()

	tests := []struct {
		name             string
		trackID          string
		withUser         bool
		setupMock        func(m *mock_track.MockUsecase)
		expectedStatus   int
		expectedResponse interface{}
	}{
		{
			name:     "Success",
			trackID:  "1",
			withUser: true,
			setupMock: func(m *mock_track.MockUsecase) {
				m.EXPECT().
					CreateStream(gomock.Any(), &usecaseModel.TrackStreamCreateData{
						TrackID: 1,
						UserID:  10,
					}).
					Return(int64(100), nil)
			},
			expectedStatus: http.StatusOK,
			expectedResponse: &delivery.APIResponse{
				Status: http.StatusOK,
				Body: &delivery.StreamID{
					ID: 100,
				},
			},
		},
		{
			name:     "Unauthorized",
			trackID:  "1",
			withUser: false,
			setupMock: func(m *mock_track.MockUsecase) {
			},
			expectedStatus: http.StatusUnauthorized,
			expectedResponse: &delivery.APIErrorResponse{
				Status: http.StatusUnauthorized,
				Error:  ErrUnauthorized,
			},
		},
		{
			name:     "Invalid ID",
			trackID:  "invalid",
			withUser: true,
			setupMock: func(m *mock_track.MockUsecase) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedResponse: &delivery.APIErrorResponse{
				Status: http.StatusBadRequest,
				Error:  "strconv.ParseInt: parsing \"invalid\": invalid syntax",
			},
		},
		{
			name:     "Stream Creation Error",
			trackID:  "1",
			withUser: true,
			setupMock: func(m *mock_track.MockUsecase) {
				m.EXPECT().
					CreateStream(gomock.Any(), &usecaseModel.TrackStreamCreateData{
						TrackID: 1,
						UserID:  10,
					}).
					Return(int64(0), errors.New("failed to create stream"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedResponse: &delivery.APIErrorResponse{
				Status: http.StatusInternalServerError,
				Error:  "failed to create stream",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase, handler, recorder := setupTest(t)
			tt.setupMock(mockUsecase)

			req := httptest.NewRequest(http.MethodPost, "/tracks/"+tt.trackID+"/stream", nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.trackID})
			ctx := helpers.LoggerToContext(req.Context(), sugar)

			if tt.withUser {
				userModel := &usecaseModel.User{
					ID:       10,
					Username: "testuser",
				}
				ctx = helpers.UserToContext(ctx, userModel)
			}

			req = req.WithContext(ctx)

			handler.CreateStream(recorder, req)

			if tt.expectedStatus == http.StatusOK {
				var response delivery.APIResponse
				err := json.NewDecoder(recorder.Body).Decode(&response)
				assert.NoError(t, err)

				expectedResp := tt.expectedResponse.(*delivery.APIResponse)
				assert.Equal(t, expectedResp.Status, response.Status)

				expectedStreamID := expectedResp.Body.(*delivery.StreamID)
				responseStreamID, ok := response.Body.(map[string]interface{})
				assert.True(t, ok)

				assert.Equal(t, float64(expectedStreamID.ID), responseStreamID["id"])
			} else {
				var response delivery.APIErrorResponse
				err := json.NewDecoder(recorder.Body).Decode(&response)
				assert.NoError(t, err)

				expectedResp := tt.expectedResponse.(*delivery.APIErrorResponse)
				assert.Equal(t, expectedResp, &response)
			}
		})
	}
}

func TestTrackHandler_UpdateStreamDuration(t *testing.T) {
	logger := zap.NewNop()
	sugar := logger.Sugar()
	defer logger.Sync()

	tests := []struct {
		name             string
		streamID         string
		requestBody      string
		withUser         bool
		setupMock        func(m *mock_track.MockUsecase)
		expectedStatus   int
		expectedResponse interface{}
	}{
		{
			name:        "Success",
			streamID:    "100",
			requestBody: `{"duration": 150}`,
			withUser:    true,
			setupMock: func(m *mock_track.MockUsecase) {
				m.EXPECT().
					UpdateStreamDuration(gomock.Any(), &usecaseModel.TrackStreamUpdateData{
						StreamID: 100,
						UserID:   10,
						Duration: 150,
					}).
					Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedResponse: &delivery.APIResponse{
				Status: http.StatusOK,
				Body: delivery.Message{
					Message: "stream duration was successfully updated",
				},
			},
		},
		{
			name:        "Unauthorized",
			streamID:    "100",
			requestBody: `{"duration": 150}`,
			withUser:    false,
			setupMock: func(m *mock_track.MockUsecase) {
			},
			expectedStatus: http.StatusUnauthorized,
			expectedResponse: &delivery.APIErrorResponse{
				Status: http.StatusUnauthorized,
				Error:  ErrUnauthorized,
			},
		},
		{
			name:        "Invalid ID",
			streamID:    "invalid",
			requestBody: `{"duration": 150}`,
			withUser:    true,
			setupMock: func(m *mock_track.MockUsecase) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedResponse: &delivery.APIErrorResponse{
				Status: http.StatusBadRequest,
				Error:  "strconv.ParseInt: parsing \"invalid\": invalid syntax",
			},
		},
		{
			name:        "Invalid Request Body",
			streamID:    "100",
			requestBody: `{"duration": -1}`,
			withUser:    true,
			setupMock: func(m *mock_track.MockUsecase) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedResponse: &delivery.APIErrorResponse{
				Status: http.StatusBadRequest,
				Error:  "duration: -1 does not validate as range(0|999999999)",
			},
		},
		{
			name:        "Permission Denied",
			streamID:    "100",
			requestBody: `{"duration": 150}`,
			withUser:    true,
			setupMock: func(m *mock_track.MockUsecase) {
				m.EXPECT().
					UpdateStreamDuration(gomock.Any(), &usecaseModel.TrackStreamUpdateData{
						StreamID: 100,
						UserID:   10,
						Duration: 150,
					}).
					Return(track.ErrStreamPermissionDenied)
			},
			expectedStatus: http.StatusForbidden,
			expectedResponse: &delivery.APIErrorResponse{
				Status: http.StatusForbidden,
				Error:  "user does not have permission to update this stream",
			},
		},
		{
			name:        "Stream Not Found",
			streamID:    "999",
			requestBody: `{"duration": 150}`,
			withUser:    true,
			setupMock: func(m *mock_track.MockUsecase) {
				m.EXPECT().
					UpdateStreamDuration(gomock.Any(), &usecaseModel.TrackStreamUpdateData{
						StreamID: 999,
						UserID:   10,
						Duration: 150,
					}).
					Return(track.ErrStreamNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedResponse: &delivery.APIErrorResponse{
				Status: http.StatusNotFound,
				Error:  "stream not found",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase, handler, recorder := setupTest(t)
			tt.setupMock(mockUsecase)

			req := httptest.NewRequest(http.MethodPut, "/streams/"+tt.streamID,
				strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			req = mux.SetURLVars(req, map[string]string{"id": tt.streamID})
			ctx := helpers.LoggerToContext(req.Context(), sugar)

			if tt.withUser {
				userModel := &usecaseModel.User{
					ID:       10,
					Username: "testuser",
				}
				ctx = helpers.UserToContext(ctx, userModel)
			}

			req = req.WithContext(ctx)

			handler.UpdateStreamDuration(recorder, req)

			if tt.expectedStatus == http.StatusOK {
				var response delivery.APIResponse
				err := json.NewDecoder(recorder.Body).Decode(&response)
				assert.NoError(t, err)

				expectedResp := tt.expectedResponse.(*delivery.APIResponse)
				assert.Equal(t, expectedResp.Status, response.Status)

				msgMap, ok := response.Body.(map[string]interface{})
				assert.True(t, ok, "Didn't get map", response.Body)
				message, ok := msgMap["msg"]
				assert.True(t, ok, "Expected 'msg'")
				assert.Equal(t, "stream duration was successfully updated", message)
			} else {
				var response delivery.APIErrorResponse
				err := json.NewDecoder(recorder.Body).Decode(&response)
				assert.NoError(t, err)

				expectedResp := tt.expectedResponse.(*delivery.APIErrorResponse)
				assert.Equal(t, expectedResp, &response)
			}
		})
	}
}

func TestTrackHandler_GetLastListenedTracks(t *testing.T) {
	logger := zap.NewNop()
	sugar := logger.Sugar()
	defer logger.Sync()

	tests := []struct {
		name             string
		url              string
		username         string
		setupMock        func(m *mock_track.MockUsecase)
		expectedStatus   int
		expectedResponse interface{}
	}{
		{
			name:     "Success",
			url:      "/users/testuser/history?offset=0&limit=10",
			username: "testuser",
			setupMock: func(m *mock_track.MockUsecase) {
				m.EXPECT().
					GetLastListenedTracks(gomock.Any(), "testuser", &usecaseModel.TrackFilters{
						Pagination: &usecaseModel.Pagination{
							Offset: 0,
							Limit:  10,
						},
					}).Return([]*usecaseModel.Track{
					{
						ID:        1,
						Title:     "Test Track",
						Duration:  180,
						Thumbnail: "test-thumbnail.jpg",
						AlbumID:   1,
						Album:     "Test Album",
						Artists: []*usecaseModel.TrackArtist{
							{
								ID:    1,
								Title: "Test Artist",
								Role:  "main",
							},
						},
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedResponse: &delivery.APIResponse{
				Status: http.StatusOK,
				Body: []*delivery.Track{
					{
						ID:        1,
						Title:     "Test Track",
						Duration:  180,
						Thumbnail: "test-thumbnail.jpg",
						AlbumID:   1,
						Album:     "Test Album",
						Artists: []*delivery.TrackArtist{
							{
								ID:    1,
								Title: "Test Artist",
								Role:  "main",
							},
						},
					},
				},
			},
		},
		{
			name:     "Invalid Pagination",
			url:      "/users/testuser/history?offset=-1",
			username: "testuser",
			setupMock: func(m *mock_track.MockUsecase) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedResponse: &delivery.APIErrorResponse{
				Status: http.StatusBadRequest,
				Error:  "invalid offset: should be greater than 0",
			},
		},
		{
			name:     "Empty Username",
			url:      "/users//history",
			username: "",
			setupMock: func(m *mock_track.MockUsecase) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedResponse: &delivery.APIErrorResponse{
				Status: http.StatusBadRequest,
				Error:  "username is required",
			},
		},
		{
			name:     "User Not Found",
			url:      "/users/nonexistent/history",
			username: "nonexistent",
			setupMock: func(m *mock_track.MockUsecase) {
				m.EXPECT().
					GetLastListenedTracks(gomock.Any(), "nonexistent", &usecaseModel.TrackFilters{
						Pagination: &usecaseModel.Pagination{
							Offset: 0,
							Limit:  10,
						},
					}).
					Return(nil, errors.New("user not found"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedResponse: &delivery.APIErrorResponse{
				Status: http.StatusInternalServerError,
				Error:  "user not found",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase, handler, recorder := setupTest(t)
			tt.setupMock(mockUsecase)

			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			req = mux.SetURLVars(req, map[string]string{"username": tt.username})
			ctx := helpers.LoggerToContext(req.Context(), sugar)
			req = req.WithContext(ctx)

			handler.GetLastListenedTracks(recorder, req)

			if tt.expectedStatus == http.StatusOK {
				var response delivery.APIResponse
				err := json.NewDecoder(recorder.Body).Decode(&response)
				assert.NoError(t, err)

				expectedResp := tt.expectedResponse.(*delivery.APIResponse)
				assert.Equal(t, expectedResp.Status, response.Status)

				expectedTracks := expectedResp.Body.([]*delivery.Track)
				responseTracks, ok := response.Body.([]interface{})
				assert.True(t, ok)
				assert.Equal(t, len(expectedTracks), len(responseTracks))

				if len(responseTracks) > 0 {
					responseTrack := responseTracks[0].(map[string]interface{})
					expectedTrack := expectedTracks[0]

					assert.Equal(t, float64(expectedTrack.ID), responseTrack["id"])
					assert.Equal(t, expectedTrack.Title, responseTrack["title"])
					assert.Equal(t, float64(expectedTrack.Duration), responseTrack["duration"])
					assert.Equal(t, expectedTrack.Thumbnail, responseTrack["thumbnail_url"])
					assert.Equal(t, expectedTrack.Album, responseTrack["album"])
					assert.Equal(t, float64(expectedTrack.AlbumID), responseTrack["album_id"])
				}
			} else {
				var response delivery.APIErrorResponse
				err := json.NewDecoder(recorder.Body).Decode(&response)
				assert.NoError(t, err)

				expectedResp := tt.expectedResponse.(*delivery.APIErrorResponse)
				assert.Equal(t, expectedResp, &response)
			}
		})
	}
}
