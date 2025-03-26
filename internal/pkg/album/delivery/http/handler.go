package album

import (
	"net/http"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/middleware"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/album"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model"
)

type AlbumHandler struct {
	usecase album.Usecase
}

func NewAlbumHandler(usecase album.Usecase) *AlbumHandler {
	return &AlbumHandler{usecase: usecase}
}

// GetAllAlbums godoc
// @Summary Get albums
// @Description Get a list of albums with optional pagination filters
// @Tags albums
// @Accept json
// @Produce json
// @Param offset query integer false "Offset (default: 0)"
// @Param limit query integer false "Limit (default: 10, max: 100)"
// @Success 200 {object} model.APIResponse{body=[]model.Album} "List of albums"
// @Failure 400 {object} model.APIBadRequestErrorResponse "Bad request - invalid filters"
// @Failure 500 {object} model.APIInternalServerErrorResponse "Internal server error"
// @Router /albums [get]
func (h *AlbumHandler) GetAllAlbums(w http.ResponseWriter, r *http.Request) {
	pagination := middleware.PaginationFromContext(r.Context())

	albums, err := h.usecase.GetAllAlbums(&model.AlbumFilters{
		Pagination: pagination,
	})
	if err != nil {
		helpers.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.WriteJSON(w, http.StatusOK, albums, nil)
}
