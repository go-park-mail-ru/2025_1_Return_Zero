package http

import (
	"net/http"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/config"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/middleware"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers"
	deliveryModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/delivery"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/track"
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
// @Accept jsonW
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
		tracks = append(tracks, &deliveryModel.Track{
			ID:        usecaseTrack.ID,
			Title:     usecaseTrack.Title,
			Thumbnail: usecaseTrack.Thumbnail,
			Duration:  usecaseTrack.Duration,
			AlbumID:   usecaseTrack.AlbumID,
			Album:     usecaseTrack.Album,
			ArtistID:  usecaseTrack.ArtistID,
			Artist:    usecaseTrack.Artist,
		})
	}
	if err != nil {
		logger.Error("failed to get tracks", zap.Error(err))
		helpers.WriteErrorResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	helpers.WriteSuccessResponse(w, http.StatusOK, tracks, nil)
}
