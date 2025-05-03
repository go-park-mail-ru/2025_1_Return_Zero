package album

import (
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/config"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/album"
	ctxExtractor "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/ctxExtractor"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/customErrors"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/errorStatus"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/json"
	loggerPkg "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/pagination"
	model "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model"
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
	ctx := r.Context()
	logger := loggerPkg.LoggerFromContext(ctx)
	pagination, err := pagination.GetPagination(r, &h.cfg.Pagination)
	if err != nil {
		logger.Error("failed to get pagination", zap.Error(err))
		json.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	usecaseAlbums, err := h.usecase.GetAllAlbums(ctx, &usecaseModel.AlbumFilters{
		Pagination: model.PaginationFromDeliveryToUsecase(pagination),
	})

	if err != nil {
		logger.Error("failed to get albums", zap.Error(err))
		json.WriteErrorResponse(w, errorStatus.ErrorStatus(err), err.Error(), nil)
		return
	}

	albums := model.AlbumsFromUsecaseToDelivery(usecaseAlbums)

	json.WriteSuccessResponse(w, http.StatusOK, albums, nil)
}

// GetAlbumsByArtistID godoc
// @Summary Get albums by artist ID
// @Description Get a list of albums for a specific artist
// @Tags albums
// @Accept json
// @Produce json
// @Param id path integer true "Artist ID"
// @Param offset query integer false "Offset (default: 0)"
// @Param limit query integer false "Limit (default: 10, max: 100)"
// @Success 200 {object} delivery.APIResponse{body=[]delivery.Album} "List of albums"
// @Failure 400 {object} delivery.APIBadRequestErrorResponse "Bad request - invalid artist ID"
// @Failure 500 {object} delivery.APIInternalServerErrorResponse "Internal server error"
// @Router /artists/{id}/albums [get]
func (h *AlbumHandler) GetAlbumsByArtistID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := loggerPkg.LoggerFromContext(ctx)
	pagination, err := pagination.GetPagination(r, &h.cfg.Pagination)
	if err != nil {
		logger.Error("failed to get pagination", zap.Error(err))
		json.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Error("failed to parse artist ID", zap.Error(err))
		json.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	usecaseAlbums, err := h.usecase.GetAlbumsByArtistID(ctx, id, &usecaseModel.AlbumFilters{
		Pagination: model.PaginationFromDeliveryToUsecase(pagination),
	})
	if err != nil {
		logger.Error("failed to get albums", zap.Error(err))
		json.WriteErrorResponse(w, errorStatus.ErrorStatus(err), err.Error(), nil)
		return
	}

	albums := model.AlbumsFromUsecaseToDelivery(usecaseAlbums)
	json.WriteSuccessResponse(w, http.StatusOK, albums, nil)
}

// GetAlbumByID godoc
// @Summary Get album by ID
// @Description Get an album by its ID
// @Tags albums
// @Accept json
// @Produce json
// @Param id path integer true "Album ID"
// @Success 200 {object} delivery.APIResponse{body=delivery.Album} "Album details"
// @Failure 400 {object} delivery.APIBadRequestErrorResponse "Bad request - invalid album ID"
// @Failure 500 {object} delivery.APIInternalServerErrorResponse "Internal server error"
// @Router /albums/{id} [get]
func (h *AlbumHandler) GetAlbumByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := loggerPkg.LoggerFromContext(ctx)

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Error("failed to parse album ID", zap.Error(err))
		json.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	usecaseAlbum, err := h.usecase.GetAlbumByID(ctx, id)
	if err != nil {
		logger.Error("failed to get album", zap.Error(err))
		json.WriteErrorResponse(w, errorStatus.ErrorStatus(err), err.Error(), nil)
		return
	}

	album := model.AlbumFromUsecaseToDelivery(usecaseAlbum, usecaseAlbum.Artists)
	json.WriteSuccessResponse(w, http.StatusOK, album, nil)
}

// LikeAlbum godoc
// @Summary Like an album
// @Description Like an album for a user
// @Tags albums
// @Accept json
// @Produce json
// @Param id path integer true "Album ID"
// @Param likeRequest body delivery.AlbumLikeRequest true "Like request"
// @Success 200 {object} delivery.APIResponse{body=delivery.Message} "Album liked/unliked"
// @Failure 400 {object} delivery.APIBadRequestErrorResponse "Bad request - invalid album ID"
// @Failure 401 {object} delivery.APIUnauthorizedErrorResponse "Unauthorized"
// @Failure 404 {object} delivery.APINotFoundErrorResponse "Album not found"
// @Failure 500 {object} delivery.APIInternalServerErrorResponse "Internal server error"
// @Router /albums/{id}/like [post]
func (h *AlbumHandler) LikeAlbum(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := loggerPkg.LoggerFromContext(ctx)

	userID, exists := ctxExtractor.UserFromContext(ctx)
	if !exists {
		logger.Warn("attempt to like album for unauthorized user")
		err := customErrors.ErrLikeAlbumUnauthorized
		json.WriteErrorResponse(w, errorStatus.ErrorStatus(err), err.Error(), nil)
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Error("failed to parse album ID", zap.Error(err))
		json.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	deliveryLikeRequest := &deliveryModel.AlbumLikeRequest{
		IsLike: true,
	}

	json.ReadJSON(w, r, deliveryLikeRequest)

	usecaseLikeRequest := model.AlbumLikeRequestFromDeliveryToUsecase(deliveryLikeRequest, userID, id)

	err = h.usecase.LikeAlbum(ctx, usecaseLikeRequest)
	if err != nil {
		logger.Error("failed to like album", zap.Error(err))
		json.WriteErrorResponse(w, errorStatus.ErrorStatus(err), err.Error(), nil)
		return
	}

	json.WriteSuccessResponse(w, http.StatusOK, deliveryModel.Message{
		Message: "album liked/unliked",
	}, nil)
}
