package artist

import (
	"net/http"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/artist"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model"
)

type ArtistHandler struct {
	usecase artist.Usecase
}

func NewArtistHandler(usecase artist.Usecase) *ArtistHandler {
	return &ArtistHandler{usecase: usecase}
}

func (h *ArtistHandler) GetAllArtists(w http.ResponseWriter, r *http.Request) {
	pagination := r.Context().Value(model.PaginationKey).(*model.Pagination)

	artists, err := h.usecase.GetAllArtists(&model.ArtistFilters{
		Pagination: pagination,
	})
	if err != nil {
		helpers.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.WriteJSON(w, http.StatusOK, artists, nil)
}
