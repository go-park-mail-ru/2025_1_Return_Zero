package artist

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/config"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/artist"
	mock_artist "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/artist/mocks"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/delivery"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func setupTest(t *testing.T) (*mock_artist.MockUsecase, *ArtistHandler, *httptest.ResponseRecorder) {
	ctrl := gomock.NewController(t)
	mockUsecase := mock_artist.NewMockUsecase(ctrl)
	cfg := &config.Config{
		Pagination: delivery.PaginationConfig{
			DefaultLimit:  10,
			DefaultOffset: 0,
			MaxLimit:      100,
		},
	}
	handler := NewArtistHandler(mockUsecase, cfg)
	recorder := httptest.NewRecorder()
	return mockUsecase, handler, recorder
}

func TestArtistHandler_GetAllArtists(t *testing.T) {
	logger := zap.NewNop()
	sugar := logger.Sugar()
	defer logger.Sync()

	tests := []struct {
		name             string
		url              string
		setupMock        func(m *mock_artist.MockUsecase)
		expectedStatus   int
		expectedResponse interface{}
	}{
		{
			name: "Success",
			url:  "/artists?offset=0&limit=10",
			setupMock: func(m *mock_artist.MockUsecase) {
				m.EXPECT().
					GetAllArtists(gomock.Any(), &usecaseModel.ArtistFilters{
						Pagination: &usecaseModel.Pagination{
							Offset: 0,
							Limit:  10,
						},
					}).Return([]*usecaseModel.Artist{
					{
						ID:          1,
						Title:       "Test Artist",
						Description: "Test Description",
						Thumbnail:   "test-thumbnail.jpg",
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedResponse: &delivery.APIResponse{
				Status: http.StatusOK,
				Body: []*delivery.Artist{
					{
						ID:          1,
						Title:       "Test Artist",
						Description: "Test Description",
						Thumbnail:   "test-thumbnail.jpg",
					},
				},
			},
		},
		{
			name: "Invalid Pagination",
			url:  "/artists?offset=-1",
			setupMock: func(m *mock_artist.MockUsecase) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedResponse: &delivery.APIErrorResponse{
				Status: http.StatusBadRequest,
				Error:  "invalid offset: should be greater than 0",
			},
		},
		{
			name: "Database Error",
			url:  "/artists",
			setupMock: func(m *mock_artist.MockUsecase) {
				m.EXPECT().
					GetAllArtists(gomock.Any(), &usecaseModel.ArtistFilters{
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

			handler.GetAllArtists(recorder, req)

			if tt.expectedStatus == http.StatusOK {
				var response delivery.APIResponse
				err := json.NewDecoder(recorder.Body).Decode(&response)
				assert.NoError(t, err)

				expectedResp := tt.expectedResponse.(*delivery.APIResponse)
				assert.Equal(t, expectedResp.Status, response.Status)

				expectedArtists := expectedResp.Body.([]*delivery.Artist)
				responseArtists, ok := response.Body.([]interface{})
				assert.True(t, ok)
				assert.Equal(t, len(expectedArtists), len(responseArtists))

				if len(responseArtists) > 0 {
					responseArtist := responseArtists[0].(map[string]interface{})
					expectedArtist := expectedArtists[0]

					assert.Equal(t, float64(expectedArtist.ID), responseArtist["id"])
					assert.Equal(t, expectedArtist.Title, responseArtist["title"])
					assert.Equal(t, expectedArtist.Description, responseArtist["description"])
					assert.Equal(t, expectedArtist.Thumbnail, responseArtist["thumbnail_url"])
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

func TestArtistHandler_GetArtistByID(t *testing.T) {
	logger := zap.NewNop()
	sugar := logger.Sugar()
	defer logger.Sync()

	tests := []struct {
		name             string
		artistID         string
		setupMock        func(m *mock_artist.MockUsecase)
		expectedStatus   int
		expectedResponse interface{}
	}{
		{
			name:     "Success",
			artistID: "1",
			setupMock: func(m *mock_artist.MockUsecase) {
				m.EXPECT().
					GetArtistByID(gomock.Any(), int64(1)).
					Return(&usecaseModel.ArtistDetailed{
						Artist: usecaseModel.Artist{
							ID:          1,
							Title:       "Test Artist",
							Description: "Test Description",
							Thumbnail:   "test-thumbnail.jpg",
						},
						Listeners: 1000,
						Favorites: 500,
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedResponse: &delivery.APIResponse{
				Status: http.StatusOK,
				Body: &delivery.ArtistDetailed{
					Artist: delivery.Artist{
						ID:          1,
						Title:       "Test Artist",
						Description: "Test Description",
						Thumbnail:   "test-thumbnail.jpg",
					},
					Listeners: 1000,
					Favorites: 500,
				},
			},
		},
		{
			name:     "Invalid ID",
			artistID: "invalid",
			setupMock: func(m *mock_artist.MockUsecase) {
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
			setupMock: func(m *mock_artist.MockUsecase) {
				m.EXPECT().
					GetArtistByID(gomock.Any(), int64(999)).
					Return(nil, artist.ErrArtistNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedResponse: &delivery.APIErrorResponse{
				Status: http.StatusNotFound,
				Error:  "artist not found",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase, handler, recorder := setupTest(t)
			tt.setupMock(mockUsecase)

			req := httptest.NewRequest(http.MethodGet, "/artists/"+tt.artistID, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.artistID})
			ctx := helpers.LoggerToContext(req.Context(), sugar)
			req = req.WithContext(ctx)

			handler.GetArtistByID(recorder, req)

			if tt.expectedStatus == http.StatusOK {
				var response delivery.APIResponse
				err := json.NewDecoder(recorder.Body).Decode(&response)
				assert.NoError(t, err)

				expectedResp := tt.expectedResponse.(*delivery.APIResponse)
				assert.Equal(t, expectedResp.Status, response.Status)

				expectedArtist := expectedResp.Body.(*delivery.ArtistDetailed)
				responseArtist, ok := response.Body.(map[string]interface{})
				assert.True(t, ok)

				assert.Equal(t, float64(expectedArtist.ID), responseArtist["id"])
				assert.Equal(t, expectedArtist.Title, responseArtist["title"])
				assert.Equal(t, expectedArtist.Description, responseArtist["description"])
				assert.Equal(t, expectedArtist.Thumbnail, responseArtist["thumbnail_url"])
				assert.Equal(t, float64(expectedArtist.Listeners), responseArtist["listeners_count"])
				assert.Equal(t, float64(expectedArtist.Favorites), responseArtist["favorites_count"])
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
