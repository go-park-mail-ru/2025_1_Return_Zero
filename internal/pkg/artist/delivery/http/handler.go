package artist

import (
	"net/http"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/middleware"
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

// GetAllArtists godoc
// @Summary Get artists
// @Description Get a list of artists with optional pagination filters
// @Tags artists
// @Accept json
// @Produce json
// @Param offset query integer false "Offset (default: 0)"
// @Param limit query integer false "Limit (default: 10, max: 100)"
// @Success 200 {object} model.APIResponse{body=[]model.Artist} "List of artists"
// @Failure 400 {object} model.APIBadRequestErrorResponse{body=model.ErrorResponse} "Bad request - invalid filters"
// @Failure 500 {object} model.APIInternalServerErrorResponse{body=model.ErrorResponse} "Internal server error"
// @Router /artists [get]
func (h *ArtistHandler) GetAllArtists(w http.ResponseWriter, r *http.Request) {
	pagination := middleware.PaginationFromContext(r.Context())

	artists, err := h.usecase.GetAllArtists(&model.ArtistFilters{
		Pagination: pagination,
	})
	if err != nil {
		helpers.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.WriteJSON(w, http.StatusOK, artists, nil)
}
