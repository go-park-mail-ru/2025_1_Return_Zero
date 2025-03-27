package album

import (
	"net/http"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/album"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers"
	deliveryModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/delivery"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
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
// @Success 200 {object} delivery.APIResponse{body=[]delivery.Album} "List of albums"
// @Failure 400 {object} delivery.APIBadRequestErrorResponse{body=delivery.ErrorResponse} "Bad request - invalid filters"
// @Failure 500 {object} delivery.APIInternalServerErrorResponse{body=delivery.ErrorResponse} "Internal server error"
// @Router /albums [get]
func (h *AlbumHandler) GetAllAlbums(w http.ResponseWriter, r *http.Request) {
	pagination, err := helpers.GetPagination(r)
	if err != nil {
		helpers.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	usecaseAlbums, err := h.usecase.GetAllAlbums(&usecaseModel.AlbumFilters{
		Pagination: &usecaseModel.Pagination{
			Offset: pagination.Offset,
			Limit:  pagination.Limit,
		},
	})
	if err != nil {
		helpers.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	albums := make([]*deliveryModel.Album, 0, len(usecaseAlbums))
	for _, usecaseAlbum := range usecaseAlbums {
		albums = append(albums, &deliveryModel.Album{
			ID:        usecaseAlbum.ID,
			Title:     usecaseAlbum.Title,
			Thumbnail: usecaseAlbum.Thumbnail,
			Artist: deliveryModel.AlbumArtist{
				ID:    usecaseAlbum.Artist.ID,
				Title: usecaseAlbum.Artist.Title,
			},
		})
	}

	helpers.WriteJSON(w, http.StatusOK, albums, nil)
}
