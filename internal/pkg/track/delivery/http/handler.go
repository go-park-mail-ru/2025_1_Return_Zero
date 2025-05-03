package http

import (
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/config"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/ctxExtractor"
	customErrors "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/customErrors"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/errorStatus"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/json"
	loggerPkg "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/pagination"
	model "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/delivery"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/track"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

const ()

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
// @Failure 404 {object} delivery.APINotFoundErrorResponse "Not found"
// @Failure 500 {object} delivery.APIInternalServerErrorResponse "Internal server error"
// @Router /tracks [get]
func (h *TrackHandler) GetAllTracks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := loggerPkg.LoggerFromContext(ctx)
	pagination, err := pagination.GetPagination(r, &h.cfg.Pagination)
	if err != nil {
		logger.Error("failed to get pagination", zap.Error(err))
		json.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	usecaseTracks, err := h.usecase.GetAllTracks(ctx, &usecaseModel.TrackFilters{
		Pagination: model.PaginationFromDeliveryToUsecase(pagination),
	})

	tracks := model.TracksFromUsecaseToDelivery(usecaseTracks)
	if err != nil {
		logger.Error("failed to get tracks", zap.Error(err))
		json.WriteErrorResponse(w, errorStatus.ErrorStatus(err), err.Error(), nil)
		return
	}

	json.WriteSuccessResponse(w, http.StatusOK, tracks, nil)
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
// @Failure 404 {object} delivery.APINotFoundErrorResponse "Not found"
// @Failure 500 {object} delivery.APIInternalServerErrorResponse "Internal server error"
// @Router /tracks/{id} [get]
func (h *TrackHandler) GetTrackByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := loggerPkg.LoggerFromContext(ctx)

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Error("failed to parse track ID", zap.Error(err))
		json.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	usecaseTrack, err := h.usecase.GetTrackByID(ctx, id)
	if err != nil {
		logger.Error("failed to get track", zap.Error(err))
		json.WriteErrorResponse(w, errorStatus.ErrorStatus(err), err.Error(), nil)
		return
	}

	trackDetailed := model.TrackDetailedFromUsecaseToDelivery(usecaseTrack)

	json.WriteSuccessResponse(w, http.StatusOK, trackDetailed, nil)
}

// GetTracksByArtistID godoc
// @Summary Get tracks by artist ID
// @Description Get a list of tracks by a specific artist with optional pagination filters
// @Tags tracks
// @Accept json
// @Produce json
// @Param id path integer true "Artist ID"
// @Param offset query integer false "Offset (default: 0)"
// @Param limit query integer false "Limit (default: 10, max: 100)"
// @Success 200 {object} delivery.APIResponse{body=[]delivery.Track} "List of tracks by artist"
// @Failure 400 {object} delivery.APIBadRequestErrorResponse "Bad request - invalid ID or filters"
// @Failure 404 {object} delivery.APINotFoundErrorResponse "Not found"
// @Failure 500 {object} delivery.APIInternalServerErrorResponse "Internal server error"
// @Router /artists/{id}/tracks [get]
func (h *TrackHandler) GetTracksByArtistID(w http.ResponseWriter, r *http.Request) {
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

	usecaseTracks, err := h.usecase.GetTracksByArtistID(ctx, id, &usecaseModel.TrackFilters{
		Pagination: model.PaginationFromDeliveryToUsecase(pagination),
	})
	if err != nil {
		logger.Error("failed to get tracks", zap.Error(err))
		json.WriteErrorResponse(w, errorStatus.ErrorStatus(err), err.Error(), nil)
		return
	}

	tracks := model.TracksFromUsecaseToDelivery(usecaseTracks)
	json.WriteSuccessResponse(w, http.StatusOK, tracks, nil)
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
	ctx := r.Context()
	logger := loggerPkg.LoggerFromContext(ctx)

	vars := mux.Vars(r)
	idStr := vars["id"]
	trackID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Error("failed to parse track ID", zap.Error(err))
		json.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	userID, exists := ctxExtractor.UserFromContext(ctx)
	if !exists {
		logger.Warn("attempt to create stream for unauthorized user")
		err := customErrors.ErrUnauthorized
		json.WriteErrorResponse(w, errorStatus.ErrorStatus(err), err.Error(), nil)
		return
	}

	trackStreamCreateData := &delivery.TrackStreamCreateData{
		TrackID: trackID,
		UserID:  userID,
	}

	streamID, err := h.usecase.CreateStream(ctx, model.TrackStreamCreateDataFromDeliveryToUsecase(trackStreamCreateData))
	if err != nil {
		logger.Error("failed to save track stream", zap.Error(err))
		json.WriteErrorResponse(w, errorStatus.ErrorStatus(err), err.Error(), nil)
		return
	}

	createResponse := &delivery.StreamID{
		ID: streamID,
	}

	json.WriteSuccessResponse(w, http.StatusOK, createResponse, nil)
}

// CreateStream godoc
// @Summary Update stream duration by id
// @Description updates listening duration at the end of stream
// @Tags tracks
// @Produce json
// @Param id path integer true "Stream ID"
// @Success 200 {object} delivery.APIResponse{body=[]delivery.Message} "Message that stream was updated"
// @Failure 400 {object} delivery.APIBadRequestErrorResponse "Bad request - invalid ID or filters"
// @Failure 401 {object} delivery.APIUnauthorizedErrorResponse "Unauthorized"
// @Failure 403 {object} delivery.APIForbiddenErrorResponse "Forbidden"
// @Failure 404 {object} delivery.APINotFoundErrorResponse "Not found"
// @Failure 500 {object} delivery.APIInternalServerErrorResponse "Internal server error"
// @Router /streams/{id} [put]
func (h *TrackHandler) UpdateStreamDuration(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := loggerPkg.LoggerFromContext(ctx)
	vars := mux.Vars(r)
	idStr := vars["id"]

	streamID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Error("failed to parse track ID", zap.Error(err))
		json.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	userID, exists := ctxExtractor.UserFromContext(ctx)
	if !exists {
		logger.Warn("attempt to update stream duration for unauthorized user")
		err := customErrors.ErrUnauthorized
		json.WriteErrorResponse(w, errorStatus.ErrorStatus(err), err.Error(), nil)
		return
	}

	var streamUpdateData delivery.TrackStreamUpdateData

	err = json.ReadJSON(w, r, &streamUpdateData)
	if err != nil {
		logger.Warn("failed to read stream duration", zap.Error(err))
		json.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	valid, err := govalidator.ValidateStruct(streamUpdateData)
	if !valid {
		logger.Warn("invalid stream duration", zap.Error(err))
		json.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	err = h.usecase.UpdateStreamDuration(ctx, model.TrackStreamUpdateDataFromDeliveryToUsecase(&streamUpdateData, userID, streamID))
	if err != nil {
		logger.Error("failed to update stream duration", zap.Error(err))
		json.WriteErrorResponse(w, errorStatus.ErrorStatus(err), err.Error(), nil)
		return
	}

	responseMessage := delivery.Message{Message: "stream duration was successfully updated"}

	json.WriteSuccessResponse(w, http.StatusOK, responseMessage, nil)
}

// GetLastListenedTracks godoc
// @Summary Get last listened tracks for a user
// @Description Retrieves a list of tracks last listened by a specific user with pagination
// @Tags tracks
// @Accept json
// @Produce json
// @Param username path string true "Username"
// @Param offset query integer false "Offset (default: 0)"
// @Param limit query integer false "Limit (default: 10, max: 100)"
// @Success 200 {object} delivery.APIResponse{body=[]delivery.Track} "List of last listened tracks"
// @Failure 400 {object} delivery.APIBadRequestErrorResponse "Bad request - invalid username or filters"
// @Failure 404 {object} delivery.APINotFoundErrorResponse "User not found"
// @Failure 500 {object} delivery.APIInternalServerErrorResponse "Internal server error"
// @Router /users/me/history [get]
func (h *TrackHandler) GetLastListenedTracks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := loggerPkg.LoggerFromContext(ctx)
	pagination, err := pagination.GetPagination(r, &h.cfg.Pagination)
	if err != nil {
		logger.Error("failed to get pagination", zap.Error(err))
		json.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	userID, exists := ctxExtractor.UserFromContext(ctx)
	if !exists {
		logger.Warn("attempt to get last listened tracks for unauthorized user")
		err := customErrors.ErrUnauthorized
		json.WriteErrorResponse(w, errorStatus.ErrorStatus(err), err.Error(), nil)
		return
	}

	usecaseTracks, err := h.usecase.GetLastListenedTracks(ctx, userID, &usecaseModel.TrackFilters{
		Pagination: model.PaginationFromDeliveryToUsecase(pagination),
	})

	if err != nil {
		logger.Error("failed to get last listened tracks", zap.Error(err))
		json.WriteErrorResponse(w, errorStatus.ErrorStatus(err), err.Error(), nil)
		return
	}

	tracks := model.TracksFromUsecaseToDelivery(usecaseTracks)
	json.WriteSuccessResponse(w, http.StatusOK, tracks, nil)
}

// GetTracksByAlbumID godoc
// @Summary Get tracks by album ID
// @Description Get a list of tracks by a specific album with optional pagination filters
// @Tags tracks
// @Accept json
// @Produce json
// @Param id path integer true "Album ID"
// @Success 200 {object} delivery.APIResponse{body=[]delivery.Track} "List of tracks by album"
// @Failure 400 {object} delivery.APIBadRequestErrorResponse "Bad request - invalid album ID or filters"
// @Failure 404 {object} delivery.APINotFoundErrorResponse "Album not found"
// @Failure 500 {object} delivery.APIInternalServerErrorResponse "Internal server error"
// @Router /albums/{id}/tracks [get]
func (h *TrackHandler) GetTracksByAlbumID(w http.ResponseWriter, r *http.Request) {
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

	usecaseTracks, err := h.usecase.GetTracksByAlbumID(ctx, id)
	if err != nil {
		logger.Error("failed to get tracks by album ID", zap.Error(err))
		json.WriteErrorResponse(w, errorStatus.ErrorStatus(err), err.Error(), nil)
		return
	}

	tracks := model.TracksFromUsecaseToDelivery(usecaseTracks)
	json.WriteSuccessResponse(w, http.StatusOK, tracks, nil)
}

// LikeTrack godoc
// @Summary Like a track
// @Description Like a track for a user
// @Tags tracks
// @Accept json
// @Produce json
// @Param id path integer true "Track ID"
// @Success 200 {object} delivery.APIResponse{body=delivery.Message} "Track liked/unliked"
// @Failure 400 {object} delivery.APIBadRequestErrorResponse "Bad request - invalid track ID"
// @Failure 401 {object} delivery.APIUnauthorizedErrorResponse "Unauthorized"
// @Failure 404 {object} delivery.APINotFoundErrorResponse "Track not found"
// @Failure 500 {object} delivery.APIInternalServerErrorResponse "Internal server error"
// @Router /tracks/{id}/like [post]
func (h *TrackHandler) LikeTrack(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := loggerPkg.LoggerFromContext(ctx)

	userID, exists := ctxExtractor.UserFromContext(ctx)
	if !exists {
		logger.Warn("attempt to like track for unauthorized user")
		err := customErrors.ErrUnauthorized
		json.WriteErrorResponse(w, errorStatus.ErrorStatus(err), err.Error(), nil)
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Error("failed to parse track ID", zap.Error(err))
		json.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	var deliveryLikeRequest delivery.TrackLikeRequest

	err = json.ReadJSON(w, r, &deliveryLikeRequest)
	if err != nil {
		logger.Warn("failed to read like request", zap.Error(err))
		json.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	usecaseLikeRequest := model.TrackLikeRequestFromDeliveryToUsecase(deliveryLikeRequest.IsLike, userID, id)

	err = h.usecase.LikeTrack(ctx, usecaseLikeRequest)
	if err != nil {
		logger.Error("failed to like track", zap.Error(err))
		json.WriteErrorResponse(w, errorStatus.ErrorStatus(err), err.Error(), nil)
		return
	}

	json.WriteSuccessResponse(w, http.StatusOK, delivery.Message{
		Message: "track liked/unliked",
	}, nil)
}

// GetPlaylistTracks godoc
// @Summary Get playlist tracks
// @Description Get a list of tracks by a specific playlist with optional pagination filters
// @Tags tracks
// @Accept json
// @Produce json
// @Param id path integer true "Playlist ID"
// @Success 200 {object} delivery.APIResponse{body=[]delivery.Track} "List of tracks by playlist"
// @Failure 400 {object} delivery.APIBadRequestErrorResponse "Bad request - invalid playlist ID or filters"
// @Failure 404 {object} delivery.APINotFoundErrorResponse "Playlist not found"
// @Failure 500 {object} delivery.APIInternalServerErrorResponse "Internal server error"
// @Router /playlists/{id}/tracks [get]
func (h *TrackHandler) GetPlaylistTracks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := loggerPkg.LoggerFromContext(ctx)

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Error("failed to parse playlist ID", zap.Error(err))
		json.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	usecaseTracks, err := h.usecase.GetPlaylistTracks(ctx, id)
	if err != nil {
		logger.Error("failed to get playlist tracks", zap.Error(err))
		json.WriteErrorResponse(w, errorStatus.ErrorStatus(err), err.Error(), nil)
		return
	}

	tracks := model.TracksFromUsecaseToDelivery(usecaseTracks)
	json.WriteSuccessResponse(w, http.StatusOK, tracks, nil)
}
