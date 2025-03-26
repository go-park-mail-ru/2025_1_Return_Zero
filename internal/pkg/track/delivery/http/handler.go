package http

import (
	"net/http"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/middleware"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/track"
)

type TrackHandler struct {
	usecase track.Usecase
}

func NewTrackHandler(usecase track.Usecase) *TrackHandler {
	return &TrackHandler{usecase: usecase}
}

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
