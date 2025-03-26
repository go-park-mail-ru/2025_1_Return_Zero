package album

import (
	"net/http"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/album"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/middleware"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model"
)

type AlbumHandler struct {
	usecase album.Usecase
}

func NewAlbumHandler(usecase album.Usecase) *AlbumHandler {
	return &AlbumHandler{usecase: usecase}
}

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
