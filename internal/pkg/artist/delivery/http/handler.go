package artist

import (
	"net/http"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/artist"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers"
	deliveryModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/delivery"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
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
// @Success 200 {object} delivery.APIResponse{body=[]delivery.Artist} "List of artists"
// @Failure 400 {object} delivery.APIBadRequestErrorResponse{body=delivery.ErrorResponse} "Bad request - invalid filters"
// @Failure 500 {object} delivery.APIInternalServerErrorResponse{body=delivery.ErrorResponse} "Internal server error"
// @Router /artists [get]
func (h *ArtistHandler) GetAllArtists(w http.ResponseWriter, r *http.Request) {
	pagination, err := helpers.GetPagination(r)
	if err != nil {
		helpers.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	usecaseArtists, err := h.usecase.GetAllArtists(&usecaseModel.ArtistFilters{
		Pagination: &usecaseModel.Pagination{
			Offset: pagination.Offset,
			Limit:  pagination.Limit,
		},
	})

	if err != nil {
		helpers.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	artists := make([]*deliveryModel.Artist, 0, len(usecaseArtists))
	for _, usecaseArtist := range usecaseArtists {
		artists = append(artists, &deliveryModel.Artist{
			ID:        usecaseArtist.ID,
			Title:     usecaseArtist.Title,
			Thumbnail: usecaseArtist.Thumbnail,
		})
	}
	helpers.WriteJSON(w, http.StatusOK, artists, nil)
}
