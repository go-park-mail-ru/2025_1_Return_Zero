package http

import (
	"net/http"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers"
	deliveryModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/delivery"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/track"
)

type TrackHandler struct {
	usecase track.Usecase
}

func NewTrackHandler(usecase track.Usecase) *TrackHandler {
	return &TrackHandler{usecase: usecase}
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
// @Failure 400 {object} delivery.APIBadRequestErrorResponse{body=delivery.ErrorResponse} "Bad request - invalid filters"
// @Failure 500 {object} delivery.APIInternalServerErrorResponse{body=delivery.ErrorResponse} "Internal server error"
// @Router /tracks [get]
func (h *TrackHandler) GetAllTracks(w http.ResponseWriter, r *http.Request) {
	pagination, err := helpers.GetPagination(r)
	if err != nil {
		helpers.WriteJSONError(w, http.StatusBadRequest, err.Error())
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
			Album: deliveryModel.TrackAlbum{
				ID:    usecaseTrack.Album.ID,
				Title: usecaseTrack.Album.Title,
			},
			Artist: deliveryModel.TrackArtist{
				ID:    usecaseTrack.Artist.ID,
				Title: usecaseTrack.Artist.Title,
			},
		})
	}
	if err != nil {
		helpers.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.WriteJSON(w, http.StatusOK, tracks, nil)
}
