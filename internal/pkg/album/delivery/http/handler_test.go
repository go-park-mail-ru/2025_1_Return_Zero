package album

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/config"
	mock_album "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/album/mocks"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/delivery"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func setupTest(t *testing.T) (*mock_album.MockUsecase, *AlbumHandler, *httptest.ResponseRecorder) {
	ctrl := gomock.NewController(t)
	mockUsecase := mock_album.NewMockUsecase(ctrl)
	cfg := &config.Config{
		Pagination: delivery.PaginationConfig{
			DefaultLimit:  10,
			DefaultOffset: 0,
			MaxLimit:      100,
		},
	}
	handler := NewAlbumHandler(mockUsecase, cfg)
	recorder := httptest.NewRecorder()
	return mockUsecase, handler, recorder
}

func TestAlbumHandler_GetAllAlbums(t *testing.T) {
	logger := zap.NewNop()
	sugar := logger.Sugar()
	defer logger.Sync()

	releaseDate := time.Date(2023, 0, 0, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name             string
		url              string
		setupMock        func(m *mock_album.MockUsecase)
		expectedStatus   int
		expectedResponse interface{}
	}{
		{
			name: "Success",
			url:  "/albums?offset=0&limit=10",
			setupMock: func(m *mock_album.MockUsecase) {
				m.EXPECT().
					GetAllAlbums(gomock.Any(), &usecaseModel.AlbumFilters{
						Pagination: &usecaseModel.Pagination{
							Offset: 0,
							Limit:  10,
						},
					}).Return([]*usecaseModel.Album{
					{
						ID:          1,
						Title:       "Test Album",
						Thumbnail:   "test-thumbnail.jpg",
						Type:        usecaseModel.AlbumTypeAlbum,
						ReleaseDate: releaseDate,
						Artists: []*usecaseModel.AlbumArtist{
							{
								ID:    1,
								Title: "Test Artist",
							},
						},
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedResponse: &delivery.APIResponse{
				Status: http.StatusOK,
				Body: []*delivery.Album{
					{
						ID:          1,
						Title:       "Test Album",
						Thumbnail:   "test-thumbnail.jpg",
						Type:        delivery.AlbumTypeAlbum,
						ReleaseDate: releaseDate,
						Artists: []*delivery.AlbumArtist{
							{
								ID:    1,
								Title: "Test Artist",
							},
						},
					},
				},
			},
		},
		{
			name: "Invalid Pagination",
			url:  "/albums?offset=-1",
			setupMock: func(m *mock_album.MockUsecase) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedResponse: &delivery.APIErrorResponse{
				Status: http.StatusBadRequest,
				Error:  "invalid offset: should be greater than 0",
			},
		},
		{
			name: "Database Error",
			url:  "/albums",
			setupMock: func(m *mock_album.MockUsecase) {
				m.EXPECT().
					GetAllAlbums(gomock.Any(), &usecaseModel.AlbumFilters{
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

			handler.GetAllAlbums(recorder, req)

			assert.Equal(t, tt.expectedStatus, recorder.Code)

			if tt.expectedStatus == http.StatusOK {
				var response delivery.APIResponse
				err := json.NewDecoder(recorder.Body).Decode(&response)
				assert.NoError(t, err)

				expectedResp := tt.expectedResponse.(*delivery.APIResponse)
				assert.Equal(t, expectedResp.Status, response.Status)

				expectedAlbums := expectedResp.Body.([]*delivery.Album)
				responseAlbums, ok := response.Body.([]interface{})
				assert.True(t, ok)
				assert.Equal(t, len(expectedAlbums), len(responseAlbums))

				if len(responseAlbums) > 0 {
					responseAlbum := responseAlbums[0].(map[string]interface{})
					expectedAlbum := expectedAlbums[0]

					assert.Equal(t, float64(expectedAlbum.ID), responseAlbum["id"])
					assert.Equal(t, expectedAlbum.Title, responseAlbum["title"])
					assert.Equal(t, expectedAlbum.Thumbnail, responseAlbum["thumbnail_url"])

					responseArtists, artistsOk := responseAlbum["artists"].([]interface{})
					assert.True(t, artistsOk)
					assert.Equal(t, len(expectedAlbum.Artists), len(responseArtists))

					if len(responseArtists) > 0 {
						responseArtist := responseArtists[0].(map[string]interface{})
						expectedArtist := expectedAlbum.Artists[0]

						assert.Equal(t, float64(expectedArtist.ID), responseArtist["id"])
						assert.Equal(t, expectedArtist.Title, responseArtist["title"])
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

func TestAlbumHandler_GetAlbumsByArtistID(t *testing.T) {
	logger := zap.NewNop()
	sugar := logger.Sugar()
	defer logger.Sync()

	releaseDate := time.Date(2023, 0, 0, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name             string
		artistID         string
		setupMock        func(m *mock_album.MockUsecase)
		expectedStatus   int
		expectedResponse interface{}
	}{
		{
			name:     "Success",
			artistID: "1",
			setupMock: func(m *mock_album.MockUsecase) {
				m.EXPECT().
					GetAlbumsByArtistID(gomock.Any(), int64(1)).
					Return([]*usecaseModel.Album{
						{
							ID:          1,
							Title:       "Test Album",
							Thumbnail:   "test-thumbnail.jpg",
							Type:        usecaseModel.AlbumTypeAlbum,
							ReleaseDate: releaseDate,
							Artists: []*usecaseModel.AlbumArtist{
								{
									ID:    1,
									Title: "Test Artist",
								},
							},
						},
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedResponse: &delivery.APIResponse{
				Status: http.StatusOK,
				Body: []*delivery.Album{
					{
						ID:          1,
						Title:       "Test Album",
						Thumbnail:   "test-thumbnail.jpg",
						Type:        delivery.AlbumTypeAlbum,
						ReleaseDate: releaseDate,
						Artists: []*delivery.AlbumArtist{
							{
								ID:    1,
								Title: "Test Artist",
							},
						},
					},
				},
			},
		},
		{
			name:     "Invalid ID",
			artistID: "invalid",
			setupMock: func(m *mock_album.MockUsecase) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedResponse: &delivery.APIErrorResponse{
				Status: http.StatusBadRequest,
				Error:  "strconv.ParseInt: parsing \"invalid\": invalid syntax",
			},
		},
		{
			name:     "Not Found",
			artistID: "1000",
			setupMock: func(m *mock_album.MockUsecase) {
				m.EXPECT().
					GetAlbumsByArtistID(gomock.Any(), int64(1000)).
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

			req := httptest.NewRequest(http.MethodGet, "/artists/"+tt.artistID+"/albums", nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.artistID})
			ctx := helpers.LoggerToContext(req.Context(), sugar)
			req = req.WithContext(ctx)

			handler.GetAlbumsByArtistID(recorder, req)

			assert.Equal(t, tt.expectedStatus, recorder.Code)

			if tt.expectedStatus == http.StatusOK {
				var response delivery.APIResponse
				err := json.NewDecoder(recorder.Body).Decode(&response)
				assert.NoError(t, err)

				expectedResp := tt.expectedResponse.(*delivery.APIResponse)
				assert.Equal(t, expectedResp.Status, response.Status)

				expectedAlbums := expectedResp.Body.([]*delivery.Album)
				responseAlbums, ok := response.Body.([]interface{})
				assert.True(t, ok)
				assert.Equal(t, len(expectedAlbums), len(responseAlbums))

				if len(responseAlbums) > 0 {
					responseAlbum := responseAlbums[0].(map[string]interface{})
					expectedAlbum := expectedAlbums[0]

					assert.Equal(t, float64(expectedAlbum.ID), responseAlbum["id"])
					assert.Equal(t, expectedAlbum.Title, responseAlbum["title"])
					assert.Equal(t, expectedAlbum.Thumbnail, responseAlbum["thumbnail_url"])

					responseArtists, artistsOk := responseAlbum["artists"].([]interface{})
					assert.True(t, artistsOk)
					assert.Equal(t, len(expectedAlbum.Artists), len(responseArtists))

					if len(responseArtists) > 0 {
						responseArtist := responseArtists[0].(map[string]interface{})
						expectedArtist := expectedAlbum.Artists[0]

						assert.Equal(t, float64(expectedArtist.ID), responseArtist["id"])
						assert.Equal(t, expectedArtist.Title, responseArtist["title"])
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
