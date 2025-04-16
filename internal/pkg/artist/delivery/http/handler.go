package artist

import (
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/config"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/artist"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers"
	model "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type ArtistHandler struct {
	usecase artist.Usecase
	cfg     *config.Config
}

func NewArtistHandler(usecase artist.Usecase, cfg *config.Config) *ArtistHandler {
	return &ArtistHandler{usecase: usecase, cfg: cfg}
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
// @Failure 400 {object} delivery.APIBadRequestErrorResponse "Bad request - invalid filters"
// @Failure 500 {object} delivery.APIInternalServerErrorResponse "Internal server error"
// @Router /artists [get]
func (h *ArtistHandler) GetAllArtists(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := helpers.LoggerFromContext(ctx)
	pagination, err := helpers.GetPagination(r, &h.cfg.Pagination)
	if err != nil {
		logger.Error("failed to get pagination", zap.Error(err))
		helpers.WriteErrorResponse(w, helpers.ErrorStatus(err), err.Error(), nil)
		return
	}

	usecaseArtists, err := h.usecase.GetAllArtists(ctx, &usecaseModel.ArtistFilters{
		Pagination: model.PaginationFromDeliveryToUsecase(pagination),
	})

	if err != nil {
		logger.Error("failed to get artists", zap.Error(err))
		helpers.WriteErrorResponse(w, helpers.ErrorStatus(err), err.Error(), nil)
		return
	}

	artists := model.ArtistsFromUsecaseToDelivery(usecaseArtists)
	helpers.WriteSuccessResponse(w, http.StatusOK, artists, nil)
}

// GetArtistByID godoc
// @Summary Get artist by ID
// @Description Get detailed information about a specific artist by their ID
// @Tags artists
// @Accept json
// @Produce json
// @Param id path integer true "Artist ID"
// @Success 200 {object} delivery.APIResponse{body=delivery.ArtistDetailed} "Artist details"
// @Failure 400 {object} delivery.APIBadRequestErrorResponse "Bad request - invalid ID"
// @Failure 500 {object} delivery.APIInternalServerErrorResponse "Internal server error"
// @Router /artists/{id} [get]
func (h *ArtistHandler) GetArtistByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := helpers.LoggerFromContext(ctx)

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Error("failed to parse artist ID", zap.Error(err))
		helpers.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	usecaseArtist, err := h.usecase.GetArtistByID(ctx, id)
	if err != nil {
		logger.Error("failed to get artist", zap.Error(err))
		helpers.WriteErrorResponse(w, helpers.ErrorStatus(err), err.Error(), nil)
		return
	}

	artistDetailed := model.ArtistDetailedFromUsecaseToDelivery(usecaseArtist)
	helpers.WriteSuccessResponse(w, http.StatusOK, artistDetailed, nil)
}
