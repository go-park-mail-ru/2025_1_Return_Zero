package album

import (
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/config"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/middleware"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/album"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers"
	deliveryModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/delivery"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type AlbumHandler struct {
	usecase album.Usecase
	cfg     *config.Config
}

func NewAlbumHandler(usecase album.Usecase, cfg *config.Config) *AlbumHandler {
	return &AlbumHandler{usecase: usecase, cfg: cfg}
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
// @Failure 400 {object} delivery.APIBadRequestErrorResponse "Bad request - invalid filters"
// @Failure 500 {object} delivery.APIInternalServerErrorResponse "Internal server error"
// @Router /albums [get]
func (h *AlbumHandler) GetAllAlbums(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())
	pagination, err := helpers.GetPagination(r, &h.cfg.Pagination)
	if err != nil {
		logger.Error("failed to get pagination", zap.Error(err))
		helpers.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	usecaseAlbums, err := h.usecase.GetAllAlbums(&usecaseModel.AlbumFilters{
		Pagination: &usecaseModel.Pagination{
			Offset: pagination.Offset,
			Limit:  pagination.Limit,
		},
	})
	if err != nil {
		logger.Error("failed to get albums", zap.Error(err))
		helpers.WriteErrorResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	albums := make([]*deliveryModel.Album, 0, len(usecaseAlbums))
	for _, usecaseAlbum := range usecaseAlbums {
		albumType := deliveryModel.AlbumType(usecaseAlbum.Type)
		albumArtists := make([]*deliveryModel.AlbumArtist, 0, len(usecaseAlbum.Artists))
		for _, usecaseArtist := range usecaseAlbum.Artists {
			albumArtists = append(albumArtists, &deliveryModel.AlbumArtist{
				ID:    usecaseArtist.ID,
				Title: usecaseArtist.Title,
			})
		}
		albums = append(albums, &deliveryModel.Album{
			ID:          usecaseAlbum.ID,
			Title:       usecaseAlbum.Title,
			Type:        albumType,
			Thumbnail:   usecaseAlbum.Thumbnail,
			Artists:     albumArtists,
			ReleaseDate: usecaseAlbum.ReleaseDate,
		})
	}

	helpers.WriteSuccessResponse(w, http.StatusOK, albums, nil)
}

// GetAlbumsByArtistID godoc
// @Summary Get albums by artist ID
// @Description Get a list of albums for a specific artist
// @Tags albums
// @Accept json
// @Produce json
// @Param id path integer true "Artist ID"
// @Success 200 {object} delivery.APIResponse{body=[]delivery.Album} "List of albums"
// @Failure 400 {object} delivery.APIBadRequestErrorResponse "Bad request - invalid artist ID"
// @Failure 500 {object} delivery.APIInternalServerErrorResponse "Internal server error"
// @Router /artists/{id}/albums [get]
func (h *AlbumHandler) GetAlbumsByArtistID(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Error("failed to parse artist ID", zap.Error(err))
		helpers.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	usecaseAlbums, err := h.usecase.GetAlbumsByArtistID(id)
	if err != nil {
		logger.Error("failed to get albums", zap.Error(err))
		helpers.WriteErrorResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	albums := make([]*deliveryModel.Album, 0, len(usecaseAlbums))
	for _, usecaseAlbum := range usecaseAlbums {
		albumArtists := make([]*deliveryModel.AlbumArtist, 0, len(usecaseAlbum.Artists))
		for _, usecaseArtist := range usecaseAlbum.Artists {
			albumArtists = append(albumArtists, &deliveryModel.AlbumArtist{
				ID:    usecaseArtist.ID,
				Title: usecaseArtist.Title,
			})
		}

		albumType := deliveryModel.AlbumType(usecaseAlbum.Type)
		albums = append(albums, &deliveryModel.Album{
			ID:          usecaseAlbum.ID,
			Title:       usecaseAlbum.Title,
			Type:        albumType,
			Thumbnail:   usecaseAlbum.Thumbnail,
			Artists:     albumArtists,
			ReleaseDate: usecaseAlbum.ReleaseDate,
		})
	}

	helpers.WriteSuccessResponse(w, http.StatusOK, albums, nil)
}
