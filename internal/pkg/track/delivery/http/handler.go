package http

import (
	"net/http"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/middleware"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model"
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
// @Success 200 {object} model.APIResponse{body=[]model.Track} "List of tracks"
// @Failure 400 {object} model.APIBadRequestErrorResponse{body=model.ErrorResponse} "Bad request - invalid filters"
// @Failure 500 {object} model.APIInternalServerErrorResponse{body=model.ErrorResponse} "Internal server error"
// @Router /tracks [get]
func (h *TrackHandler) GetAllTracks(w http.ResponseWriter, r *http.Request) {
	pagination := middleware.PaginationFromContext(r.Context())

	tracks, err := h.usecase.GetAllTracks(&model.TrackFilters{
		Pagination: pagination,
	})
	if err != nil {
		helpers.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.WriteJSON(w, http.StatusOK, tracks, nil)
}
