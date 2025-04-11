package http

import (
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/config"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/middleware"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers"
	model "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/delivery"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/track"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

const (
	unauthorizedError = "unauthorized users can't save to stream history"
)

type TrackHandler struct {
	usecase track.Usecase
	cfg     *config.Config
}

func NewTrackHandler(usecase track.Usecase, cfg *config.Config) *TrackHandler {
	return &TrackHandler{usecase: usecase, cfg: cfg}
}

// GetAllTracks godoc
// @Summary Get tracks
// @Description Get a list of tracks with optional pagination filters
// @Tags tracks
// @Accept json
// @Produce json
// @Param offset query integer false "Offset (default: 0)"
// @Param limit query integer false "Limit (default: 10, max: 100)"
// @Success 200 {object} delivery.APIResponse{body=[]delivery.Track} "List of tracks"
// @Failure 400 {object} delivery.APIBadRequestErrorResponse "Bad request - invalid filters"
// @Failure 500 {object} delivery.APIInternalServerErrorResponse "Internal server error"
// @Router /tracks [get]
func (h *TrackHandler) GetAllTracks(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())
	pagination, err := helpers.GetPagination(r, &h.cfg.Pagination)
	if err != nil {
		logger.Error("failed to get pagination", zap.Error(err))
		helpers.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	usecaseTracks, err := h.usecase.GetAllTracks(&usecaseModel.TrackFilters{
		Pagination: model.PaginationFromDeliveryToUsecase(pagination),
	})

	tracks := model.TracksFromUsecaseToDelivery(usecaseTracks)
	if err != nil {
		logger.Error("failed to get tracks", zap.Error(err))
		helpers.WriteErrorResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	helpers.WriteSuccessResponse(w, http.StatusOK, tracks, nil)
}

// GetTrackByID godoc
// @Summary Get track by ID
// @Description Retrieves a specific track by its ID with detailed information
// @Tags tracks
// @Accept json
// @Produce json
// @Param id path int true "Track ID"
// @Success 200 {object} delivery.APIResponse{body=delivery.TrackDetailed} "Track details"
// @Failure 400 {object} delivery.APIBadRequestErrorResponse "Bad request - invalid ID"
// @Failure 500 {object} delivery.APIInternalServerErrorResponse "Internal server error"
// @Router /tracks/{id} [get]
func (h *TrackHandler) GetTrackByID(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Error("failed to parse track ID", zap.Error(err))
		helpers.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	usecaseTrack, err := h.usecase.GetTrackByID(id)
	if err != nil {
		logger.Error("failed to get track", zap.Error(err))
		helpers.WriteErrorResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	trackDetailed := model.TrackDetailedFromUsecaseToDelivery(usecaseTrack)

	helpers.WriteSuccessResponse(w, http.StatusOK, trackDetailed, nil)
}

// GetTracksByArtistID godoc
// @Summary Get tracks by artist ID
// @Description Get a list of tracks by a specific artist with optional pagination filters
// @Tags tracks
// @Accept json
// @Produce json
// @Param id path integer true "Artist ID"
// @Success 200 {object} delivery.APIResponse{body=[]delivery.Track} "List of tracks by artist"
// @Failure 400 {object} delivery.APIBadRequestErrorResponse "Bad request - invalid ID or filters"
// @Failure 500 {object} delivery.APIInternalServerErrorResponse "Internal server error"
// @Router /artists/{id}/tracks [get]
func (h *TrackHandler) GetTracksByArtistID(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Error("failed to parse artist ID", zap.Error(err))
		helpers.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	usecaseTracks, err := h.usecase.GetTracksByArtistID(id)
	if err != nil {
		logger.Error("failed to get tracks", zap.Error(err))
		helpers.WriteErrorResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	tracks := model.TracksFromUsecaseToDelivery(usecaseTracks)
	helpers.WriteSuccessResponse(w, http.StatusOK, tracks, nil)
}

// CreateStream godoc
// @Summary Create stream for track by id
// @Description Creates stream for track by id, essentially it means saving track to listening history
// @Tags tracks
// @Produce json
// @Param id path integer true "Track ID"
// @Success 200 {object} delivery.APIResponse{body=[]delivery.StreamID} "ID of created stream"
// @Failure 400 {object} delivery.APIBadRequestErrorResponse "Bad request - invalid ID or filters"
// @Failure 500 {object} delivery.APIInternalServerErrorResponse "Internal server error"
// @Router /tracks/{id}/stream [post]
func (h *TrackHandler) CreateStream(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())

	vars := mux.Vars(r)
	idStr := vars["id"]
	trackID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Error("failed to parse track ID", zap.Error(err))
		helpers.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
	}

	user, exists := middleware.GetUserFromContext(r.Context())
	if !exists {
		logger.Warn("attempt to create stream for unauthorized user")
		helpers.WriteErrorResponse(w, http.StatusUnauthorized, unauthorizedError, nil)
		return
	}
	userID := user.ID

	trackStreamCreateData := &delivery.TrackStreamCreateData{
		TrackID: trackID,
		UserID:  userID,
	}

	streamID, err := h.usecase.CreateStream(model.TrackStreamCreateDataFromDeliveryToUsecase(trackStreamCreateData))
	if err != nil {
		logger.Error("failed to save track stream", zap.Error(err))
		helpers.WriteErrorResponse(w, http.StatusInternalServerError, err.Error(), nil)
	}

	createResponse := &delivery.StreamID{
		ID: streamID,
	}

	helpers.WriteSuccessResponse(w, http.StatusOK, createResponse, nil)
}

// CreateStream godoc
// @Summary Update stream duration by id
// @Description updates listening duration at the end of stream
// @Tags tracks
// @Produce json
// @Param id path integer true "Stream ID"
// @Success 200 {object} delivery.APIResponse{body=[]delivery.Message} "Message that stream was updated"
// @Failure 400 {object} delivery.APIBadRequestErrorResponse "Bad request - invalid ID or filters"
// @Failure 500 {object} delivery.APIInternalServerErrorResponse "Internal server error"
// @Router /streams/{id} [put]
func (h *TrackHandler) UpdateStreamDuration(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())
	vars := mux.Vars(r)
	idStr := vars["id"]

	streamID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Error("failed to parse track ID", zap.Error(err))
		helpers.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
	}

	user, exists := middleware.GetUserFromContext(r.Context())
	if !exists {
		logger.Warn("attempt to update stream duration for unauthorized user")
		helpers.WriteErrorResponse(w, http.StatusUnauthorized, unauthorizedError, nil)
	}

	userID := user.ID

	var streamUpdateData delivery.TrackStreamUpdateData

	err = helpers.ReadJSON(w, r, &streamUpdateData)
	if err != nil {
		logger.Warn("failed to read stream duration", zap.Error(err))
		helpers.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	valid, err := govalidator.ValidateStruct(streamUpdateData)
	if !valid {
		logger.Warn("invalid stream duration", zap.Error(err))
		helpers.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	err = h.usecase.UpdateStreamDuration(model.TrackStreamUpdateDataFromDeliveryToUsecase(&streamUpdateData, userID, streamID))
	if err != nil {
		logger.Error("failed to update stream duration", zap.Error(err))
		helpers.WriteErrorResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	responseMessage := delivery.Message{Message: "stream duration was successfully updated"}

	helpers.WriteSuccessResponse(w, http.StatusOK, responseMessage, nil)
}
