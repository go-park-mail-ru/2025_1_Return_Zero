package http

import (
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/config"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/middleware"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers"
	deliveryModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/delivery"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/track"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type TrackHandler struct {
	usecase track.Usecase
	cfg     *config.Config
}

func NewTrackHandler(usecase track.Usecase, cfg *config.Config) *TrackHandler {
	return &TrackHandler{usecase: usecase, cfg: cfg}
}

// GetAllTracks godoc
// @Summary Get tracks
// @Description Get a list of tracks with optional pagination filters
// @Tags tracks
// @Accept json
// @Produce json
// @Param offset query integer false "Offset (default: 0)"
// @Param limit query integer false "Limit (default: 10, max: 100)"
// @Success 200 {object} delivery.APIResponse{body=[]delivery.Track} "List of tracks"
// @Failure 400 {object} delivery.APIBadRequestErrorResponse "Bad request - invalid filters"
// @Failure 500 {object} delivery.APIInternalServerErrorResponse "Internal server error"
// @Router /tracks [get]
func (h *TrackHandler) GetAllTracks(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())
	pagination, err := helpers.GetPagination(r, &h.cfg.Pagination)
	if err != nil {
		logger.Error("failed to get pagination", zap.Error(err))
		helpers.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	usecaseTracks, err := h.usecase.GetAllTracks(&usecaseModel.TrackFilters{
		Pagination: &usecaseModel.Pagination{
			Offset: pagination.Offset,
			Limit:  pagination.Limit,
		},
	})

	tracks := make([]*deliveryModel.Track, 0, len(usecaseTracks))
	for _, usecaseTrack := range usecaseTracks {
		artists := make([]*deliveryModel.TrackArtist, 0, len(usecaseTrack.Artists))
		for _, artist := range usecaseTrack.Artists {
			artists = append(artists, &deliveryModel.TrackArtist{
				ID:    artist.ID,
				Title: artist.Title,
				Role:  artist.Role,
			})
		}
		tracks = append(tracks, &deliveryModel.Track{
			ID:        usecaseTrack.ID,
			Title:     usecaseTrack.Title,
			Thumbnail: usecaseTrack.Thumbnail,
			Duration:  usecaseTrack.Duration,
			AlbumID:   usecaseTrack.AlbumID,
			Album:     usecaseTrack.Album,
			Artists:   artists,
		})
	}
	if err != nil {
		logger.Error("failed to get tracks", zap.Error(err))
		helpers.WriteErrorResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	helpers.WriteSuccessResponse(w, http.StatusOK, tracks, nil)
}

// GetTrackByID godoc
// @Summary Get track by ID
// @Description Retrieves a specific track by its ID with detailed information
// @Tags tracks
// @Accept json
// @Produce json
// @Param id path int true "Track ID"
// @Success 200 {object} delivery.APIResponse{body=delivery.TrackDetailed} "Track details"
// @Failure 400 {object} delivery.APIBadRequestErrorResponse "Bad request - invalid ID"
// @Failure 500 {object} delivery.APIInternalServerErrorResponse "Internal server error"
// @Router /tracks/{id} [get]
func (h *TrackHandler) GetTrackByID(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Error("failed to parse track ID", zap.Error(err))
		helpers.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	usecaseTrack, err := h.usecase.GetTrackByID(id)
	if err != nil {
		logger.Error("failed to get track", zap.Error(err))
		helpers.WriteErrorResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	artists := make([]*deliveryModel.TrackArtist, 0, len(usecaseTrack.Artists))
	for _, artist := range usecaseTrack.Artists {
		artists = append(artists, &deliveryModel.TrackArtist{
			ID:    artist.ID,
			Title: artist.Title,
			Role:  artist.Role,
		})
	}

	track := &deliveryModel.TrackDetailed{
		Track: deliveryModel.Track{
			ID:        usecaseTrack.ID,
			Title:     usecaseTrack.Title,
			Thumbnail: usecaseTrack.Thumbnail,
			Duration:  usecaseTrack.Duration,
			Album:     usecaseTrack.Album,
			AlbumID:   usecaseTrack.AlbumID,
			Artists:   artists,
		},
		FileUrl: usecaseTrack.FileUrl,
	}

	helpers.WriteSuccessResponse(w, http.StatusOK, track, nil)
}

// GetTracksByArtistID godoc
// @Summary Get tracks by artist ID
// @Description Get a list of tracks by a specific artist with optional pagination filters
// @Tags tracks
// @Accept json
// @Produce json
// @Param id path integer true "Artist ID"
// @Success 200 {object} delivery.APIResponse{body=[]delivery.Track} "List of tracks by artist"
// @Failure 400 {object} delivery.APIBadRequestErrorResponse "Bad request - invalid ID or filters"
// @Failure 500 {object} delivery.APIInternalServerErrorResponse "Internal server error"
// @Router /artists/{id}/tracks [get]
func (h *TrackHandler) GetTracksByArtistID(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Error("failed to parse artist ID", zap.Error(err))
		helpers.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	usecaseTracks, err := h.usecase.GetTracksByArtistID(id)
	if err != nil {
		logger.Error("failed to get tracks", zap.Error(err))
		helpers.WriteErrorResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	tracks := make([]*deliveryModel.Track, 0, len(usecaseTracks))
	for _, usecaseTrack := range usecaseTracks {
		artists := make([]*deliveryModel.TrackArtist, 0, len(usecaseTrack.Artists))
		for _, artist := range usecaseTrack.Artists {
			artists = append(artists, &deliveryModel.TrackArtist{
				ID:    artist.ID,
				Title: artist.Title,
				Role:  artist.Role,
			})
		}
		tracks = append(tracks, &deliveryModel.Track{
			ID:        usecaseTrack.ID,
			Title:     usecaseTrack.Title,
			Thumbnail: usecaseTrack.Thumbnail,
			Duration:  usecaseTrack.Duration,
			Album:     usecaseTrack.Album,
			AlbumID:   usecaseTrack.AlbumID,
			Artists:   artists,
		})
	}

	helpers.WriteSuccessResponse(w, http.StatusOK, tracks, nil)
}
