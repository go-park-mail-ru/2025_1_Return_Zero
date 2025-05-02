package http

import (
	"io"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/config"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/ctxExtractor"
	customErrors "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/customErrors"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/errorStatus"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/json"
	loggerPkg "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/query"
	model "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/delivery"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/playlist"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

const ()

type PlaylistHandler struct {
	usecase playlist.Usecase
	cfg     *config.Config
}

func NewPlaylistHandler(usecase playlist.Usecase, cfg *config.Config) *PlaylistHandler {
	return &PlaylistHandler{usecase: usecase, cfg: cfg}
}

// CreatePlaylist godoc
// @Summary Create a new playlist
// @Description Create a new playlist with a title and thumbnail image
// @Tags playlists
// @Accept multipart/form-data
// @Produce json
// @Param title formData string true "Playlist title"
// @Param thumbnail formData file true "Playlist thumbnail image"
// @Success 201 {object} delivery.APIResponse{body=delivery.Playlist} "Created playlist"
// @Failure 400 {object} delivery.APIBadRequestErrorResponse "Bad request - invalid form data"
// @Failure 401 {object} delivery.APIUnauthorizedErrorResponse "Unauthorized"
// @Failure 413 {object} delivery.APIRequestEntityTooLargeErrorResponse "Request entity too large"
// @Failure 500 {object} delivery.APIInternalServerErrorResponse "Internal server error"
// @Router /playlists [post]
func (h *PlaylistHandler) CreatePlaylist(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := loggerPkg.LoggerFromContext(ctx)

	user, exists := ctxExtractor.UserFromContext(ctx)
	if !exists {
		logger.Warn("attempt to create playlist for unauthorized user")
		err := customErrors.ErrPlaylistUnauthorized
		json.WriteErrorResponse(w, errorStatus.ErrorStatus(err), err.Error(), nil)
		return
	}

	if r.ContentLength > 6<<20 {
		logger.Error("content length is too large")
		json.WriteErrorResponse(w, http.StatusRequestEntityTooLarge, "content length is too large", nil)
		return
	}

	err := r.ParseMultipartForm(6 << 20) // 6 MB
	if err != nil {
		logger.Error("failed to parse multipart form", zap.Error(err))
		json.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	request := &delivery.CreatePlaylistRequest{}

	request.Title = r.FormValue("title")
	if request.Title == "" {
		logger.Error("title is required")
		json.WriteErrorResponse(w, http.StatusBadRequest, "title is required", nil)
		return
	}

	if len(request.Title) > 100 {
		logger.Error("title is too long")
		json.WriteErrorResponse(w, http.StatusBadRequest, "title is too long", nil)
		return
	}

	file, _, err := r.FormFile("thumbnail")
	if err != nil {
		logger.Error("failed to get thumbnail", zap.Error(err))
		json.WriteErrorResponse(w, errorStatus.ErrorStatus(customErrors.ErrPlaylistImageNotUploaded), customErrors.ErrPlaylistImageNotUploaded.Error(), nil)
		return
	}

	request.Thumbnail, err = io.ReadAll(file)
	if err != nil {
		logger.Error("failed to read thumbnail", zap.Error(err))
		json.WriteErrorResponse(w, errorStatus.ErrorStatus(customErrors.ErrPlaylistImageNotUploaded), customErrors.ErrPlaylistImageNotUploaded.Error(), nil)
		return
	}

	usecaseRequest := model.CreatePlaylistRequestFromDeliveryToUsecase(request, user.ID)

	playlist, err := h.usecase.CreatePlaylist(ctx, usecaseRequest)
	if err != nil {
		logger.Error("failed to create playlist", zap.Error(err))
		json.WriteErrorResponse(w, errorStatus.ErrorStatus(err), err.Error(), nil)
		return
	}

	deliveryPlaylist := model.PlaylistFromUsecaseToDelivery(playlist)

	json.WriteSuccessResponse(w, http.StatusCreated, deliveryPlaylist, nil)
}

// GetCombinedPlaylistsForCurrentUser godoc
// @Summary Get combined playlists for current user
// @Description Retrieves all playlists accessible to the current user
// @Tags playlists
// @Accept json
// @Produce json
// @Success 200 {object} delivery.APIResponse{body=[]delivery.Playlist} "List of playlists"
// @Failure 401 {object} delivery.APIUnauthorizedErrorResponse "Unauthorized"
// @Failure 500 {object} delivery.APIInternalServerErrorResponse "Internal server error"
// @Router /playlists [get]
func (h *PlaylistHandler) GetCombinedPlaylistsForCurrentUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := loggerPkg.LoggerFromContext(ctx)

	user, exists := ctxExtractor.UserFromContext(ctx)
	if !exists {
		logger.Warn("attempt to get combined playlists for unauthorized user")
		err := customErrors.ErrPlaylistUnauthorized
		json.WriteErrorResponse(w, errorStatus.ErrorStatus(err), err.Error(), nil)
		return
	}

	playlists, err := h.usecase.GetCombinedPlaylistsForCurrentUser(ctx, user.ID)
	if err != nil {
		logger.Error("failed to get combined playlists", zap.Error(err))
		json.WriteErrorResponse(w, errorStatus.ErrorStatus(err), err.Error(), nil)
		return
	}

	deliveryPlaylists := model.PlaylistsFromUsecaseToDelivery(playlists)

	json.WriteSuccessResponse(w, http.StatusOK, deliveryPlaylists, nil)
}

// AddTrackToPlaylist godoc
// @Summary Add a track to a playlist
// @Description Adds a track to a specific playlist
// @Tags playlists
// @Accept json
// @Produce json
// @Param id path integer true "Playlist ID"
// @Param request body delivery.AddTrackToPlaylistRequest true "Track information"
// @Success 200 {object} delivery.APIResponse{} "Track added successfully"
// @Failure 400 {object} delivery.APIBadRequestErrorResponse "Bad request - invalid ID or request body"
// @Failure 401 {object} delivery.APIUnauthorizedErrorResponse "Unauthorized"
// @Failure 404 {object} delivery.APINotFoundErrorResponse "Playlist or track not found"
// @Failure 500 {object} delivery.APIInternalServerErrorResponse "Internal server error"
// @Router /playlists/{id}/tracks [post]
func (h *PlaylistHandler) AddTrackToPlaylist(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := loggerPkg.LoggerFromContext(ctx)

	user, exists := ctxExtractor.UserFromContext(ctx)
	if !exists {
		logger.Warn("attempt to add track to playlist for unauthorized user")
		err := customErrors.ErrPlaylistUnauthorized
		json.WriteErrorResponse(w, errorStatus.ErrorStatus(err), err.Error(), nil)
		return
	}

	vars := mux.Vars(r)
	playlistID := vars["id"]
	if playlistID == "" {
		logger.Error("playlist_id is required")
		json.WriteErrorResponse(w, http.StatusBadRequest, "playlist_id is required", nil)
		return
	}

	playlistIDInt, err := strconv.ParseInt(playlistID, 10, 64)
	if err != nil {
		logger.Error("failed to parse playlist_id", zap.Error(err))
		json.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	request := &delivery.AddTrackToPlaylistRequest{}
	err = json.ReadJSON(w, r, request)
	if err != nil {
		logger.Error("failed to read add track to playlist request", zap.Error(err))
		json.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	usecaseRequest := model.AddTrackToPlaylistRequestFromDeliveryToUsecase(request, user.ID, playlistIDInt)

	err = h.usecase.AddTrackToPlaylist(ctx, usecaseRequest)
	if err != nil {
		logger.Error("failed to add track to playlist", zap.Error(err))
		json.WriteErrorResponse(w, errorStatus.ErrorStatus(err), err.Error(), nil)
		return
	}

	json.WriteSuccessResponse(w, http.StatusOK, delivery.Message{Message: "Track added successfully"}, nil)
}

// RemoveTrackFromPlaylist godoc
// @Summary Remove a track from a playlist
// @Description Removes a track from a specific playlist
// @Tags playlists
// @Accept json
// @Produce json
// @Param id path integer true "Playlist ID"
// @Param trackId path integer true "Track ID"
// @Success 200 {object} delivery.APIResponse{} "Track removed successfully"
// @Failure 400 {object} delivery.APIBadRequestErrorResponse "Bad request - invalid IDs"
// @Failure 401 {object} delivery.APIUnauthorizedErrorResponse "Unauthorized"
// @Failure 404 {object} delivery.APINotFoundErrorResponse "Playlist or track not found"
// @Failure 500 {object} delivery.APIInternalServerErrorResponse "Internal server error"
// @Router /playlists/{id}/tracks/{trackId} [delete]
func (h *PlaylistHandler) RemoveTrackFromPlaylist(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := loggerPkg.LoggerFromContext(ctx)

	user, exists := ctxExtractor.UserFromContext(ctx)
	if !exists {
		logger.Warn("attempt to remove track from playlist for unauthorized user")
		err := customErrors.ErrPlaylistUnauthorized
		json.WriteErrorResponse(w, errorStatus.ErrorStatus(err), err.Error(), nil)
		return
	}

	vars := mux.Vars(r)
	playlistID := vars["id"]
	if playlistID == "" {
		logger.Error("playlist_id is required")
		json.WriteErrorResponse(w, http.StatusBadRequest, "playlist_id is required", nil)
		return
	}

	playlistIDInt, err := strconv.ParseInt(playlistID, 10, 64)
	if err != nil {
		logger.Error("failed to parse playlist_id", zap.Error(err))
		json.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	trackID := vars["trackId"]
	if trackID == "" {
		logger.Error("trackId is required")
		json.WriteErrorResponse(w, http.StatusBadRequest, "trackId is required", nil)
		return
	}

	trackIDInt, err := strconv.ParseInt(trackID, 10, 64)
	if err != nil {
		logger.Error("failed to parse trackId", zap.Error(err))
		json.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	usecaseRequest := model.RemoveTrackFromPlaylistRequestFromDeliveryToUsecase(trackIDInt, user.ID, playlistIDInt)

	err = h.usecase.RemoveTrackFromPlaylist(ctx, usecaseRequest)
	if err != nil {
		logger.Error("failed to remove track from playlist", zap.Error(err))
		json.WriteErrorResponse(w, errorStatus.ErrorStatus(err), err.Error(), nil)
		return
	}

	json.WriteSuccessResponse(w, http.StatusOK, delivery.Message{Message: "Track removed successfully"}, nil)
}

// UpdatePlaylist godoc
// @Summary Update a playlist
// @Description Update a playlist's title and/or thumbnail
// @Tags playlists
// @Accept multipart/form-data
// @Produce json
// @Param id path integer true "Playlist ID"
// @Param title formData string true "New playlist title"
// @Param thumbnail formData file false "New playlist thumbnail image"
// @Success 200 {object} delivery.APIResponse{body=delivery.Playlist} "Updated playlist"
// @Failure 400 {object} delivery.APIBadRequestErrorResponse "Bad request - invalid ID or form data"
// @Failure 401 {object} delivery.APIUnauthorizedErrorResponse "Unauthorized"
// @Failure 404 {object} delivery.APINotFoundErrorResponse "Playlist not found"
// @Failure 500 {object} delivery.APIInternalServerErrorResponse "Internal server error"
// @Router /playlists/{id} [put]
func (h *PlaylistHandler) UpdatePlaylist(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := loggerPkg.LoggerFromContext(ctx)

	user, exists := ctxExtractor.UserFromContext(ctx)
	if !exists {
		logger.Warn("attempt to update playlist for unauthorized user")
		err := customErrors.ErrPlaylistUnauthorized
		json.WriteErrorResponse(w, errorStatus.ErrorStatus(err), err.Error(), nil)
		return
	}

	vars := mux.Vars(r)
	playlistID := vars["id"]
	if playlistID == "" {
		logger.Error("playlist_id is required")
		json.WriteErrorResponse(w, http.StatusBadRequest, "playlist_id is required", nil)
		return
	}

	playlistIDInt, err := strconv.ParseInt(playlistID, 10, 64)
	if err != nil {
		logger.Error("failed to parse playlist_id", zap.Error(err))
		json.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	request := &delivery.UpdatePlaylistRequest{}
	request.Title = r.FormValue("title")
	if request.Title == "" {
		logger.Error("title is required")
		json.WriteErrorResponse(w, http.StatusBadRequest, "title is required", nil)
		return
	}

	file, _, err := r.FormFile("thumbnail")
	if err != nil {
		request.Thumbnail = nil
	} else {
		request.Thumbnail, err = io.ReadAll(file)
		if err != nil {
			logger.Error("failed to read thumbnail", zap.Error(err))
			json.WriteErrorResponse(w, errorStatus.ErrorStatus(customErrors.ErrPlaylistImageNotUploaded), customErrors.ErrPlaylistImageNotUploaded.Error(), nil)
			return
		}
	}

	usecaseRequest := model.UpdatePlaylistRequestFromDeliveryToUsecase(request, user.ID, playlistIDInt)

	playlist, err := h.usecase.UpdatePlaylist(ctx, usecaseRequest)
	if err != nil {
		logger.Error("failed to update playlist", zap.Error(err))
		json.WriteErrorResponse(w, errorStatus.ErrorStatus(err), err.Error(), nil)
		return
	}

	deliveryPlaylist := model.PlaylistFromUsecaseToDelivery(playlist)

	json.WriteSuccessResponse(w, http.StatusOK, deliveryPlaylist, nil)
}

// GetPlaylistByID godoc
// @Summary Get a playlist by ID
// @Description Retrieves a specific playlist by its ID with all its details
// @Tags playlists
// @Accept json
// @Produce json
// @Param id path integer true "Playlist ID"
// @Success 200 {object} delivery.APIResponse{body=delivery.Playlist} "Playlist details"
// @Failure 400 {object} delivery.APIBadRequestErrorResponse "Bad request - invalid ID"
// @Failure 404 {object} delivery.APINotFoundErrorResponse "Playlist not found"
// @Failure 500 {object} delivery.APIInternalServerErrorResponse "Internal server error"
// @Router /playlists/{id} [get]
func (h *PlaylistHandler) GetPlaylistByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := loggerPkg.LoggerFromContext(ctx)

	vars := mux.Vars(r)
	playlistID := vars["id"]
	if playlistID == "" {
		logger.Error("playlist_id is required")
		json.WriteErrorResponse(w, http.StatusBadRequest, "playlist_id is required", nil)
		return
	}

	playlistIDInt, err := strconv.ParseInt(playlistID, 10, 64)
	if err != nil {
		logger.Error("failed to parse playlist_id", zap.Error(err))
		json.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	playlist, err := h.usecase.GetPlaylistByID(ctx, playlistIDInt)
	if err != nil {
		logger.Error("failed to get playlist by id", zap.Error(err))
		json.WriteErrorResponse(w, errorStatus.ErrorStatus(err), err.Error(), nil)
		return
	}

	deliveryPlaylist := model.PlaylistFromUsecaseToDelivery(playlist)

	json.WriteSuccessResponse(w, http.StatusOK, deliveryPlaylist, nil)
}

// RemovePlaylist godoc
// @Summary Remove a playlist
// @Description Deletes a playlist by its ID (only available to the playlist owner)
// @Tags playlists
// @Accept json
// @Produce json
// @Param id path integer true "Playlist ID"
// @Success 200 {object} delivery.APIResponse{body=delivery.Message} "Playlist removed successfully"
// @Failure 400 {object} delivery.APIBadRequestErrorResponse "Bad request - invalid ID"
// @Failure 401 {object} delivery.APIUnauthorizedErrorResponse "Unauthorized"
// @Failure 403 {object} delivery.APIForbiddenErrorResponse "Forbidden - user is not the playlist owner"
// @Failure 404 {object} delivery.APINotFoundErrorResponse "Playlist not found"
// @Failure 500 {object} delivery.APIInternalServerErrorResponse "Internal server error"
// @Router /playlists/{id} [delete]
func (h *PlaylistHandler) RemovePlaylist(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := loggerPkg.LoggerFromContext(ctx)

	user, exists := ctxExtractor.UserFromContext(ctx)
	if !exists {
		logger.Warn("attempt to remove playlist for unauthorized user")
		err := customErrors.ErrPlaylistUnauthorized
		json.WriteErrorResponse(w, errorStatus.ErrorStatus(err), err.Error(), nil)
		return
	}

	vars := mux.Vars(r)
	playlistID := vars["id"]
	if playlistID == "" {
		logger.Error("playlist_id is required")
		json.WriteErrorResponse(w, http.StatusBadRequest, "playlist_id is required", nil)
		return
	}

	playlistIDInt, err := strconv.ParseInt(playlistID, 10, 64)
	if err != nil {
		logger.Error("failed to parse playlist_id", zap.Error(err))
		json.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	usecaseRequest := model.RemovePlaylistRequestFromDeliveryToUsecase(playlistIDInt, user.ID)

	err = h.usecase.RemovePlaylist(ctx, usecaseRequest)
	if err != nil {
		logger.Error("failed to remove playlist", zap.Error(err))
		json.WriteErrorResponse(w, errorStatus.ErrorStatus(err), err.Error(), nil)
		return
	}

	json.WriteSuccessResponse(w, http.StatusOK, delivery.Message{Message: "Playlist removed successfully"}, nil)
}

// GetPlaylistsToAdd godoc
// @Summary Get playlists available for adding a track
// @Description Retrieves all playlists owned by the current user with information about whether the track is already included
// @Tags playlists
// @Accept json
// @Produce json
// @Param trackId query integer true "Track ID to check inclusion status"
// @Success 200 {object} delivery.APIResponse{body=[]delivery.PlaylistWithIsIncludedTrack} "List of playlists with track inclusion status"
// @Failure 400 {object} delivery.APIBadRequestErrorResponse "Bad request - invalid track ID"
// @Failure 401 {object} delivery.APIUnauthorizedErrorResponse "Unauthorized"
// @Failure 500 {object} delivery.APIInternalServerErrorResponse "Internal server error"
// @Router /playlists/to-add [get]
func (h *PlaylistHandler) GetPlaylistsToAdd(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := loggerPkg.LoggerFromContext(ctx)

	user, exists := ctxExtractor.UserFromContext(ctx)
	if !exists {
		logger.Warn("attempt to get playlists to add for unauthorized user")
		err := customErrors.ErrPlaylistUnauthorized
		json.WriteErrorResponse(w, errorStatus.ErrorStatus(err), err.Error(), nil)
		return
	}

	trackID, err := query.ReadInt(r.URL.Query(), "trackId", -1)
	if err != nil {
		logger.Error("failed to read track_id", zap.Error(err))
		json.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	usecaseRequest := model.GetPlaylistsToAddRequestFromDeliveryToUsecase(int64(trackID), user.ID)

	playlists, err := h.usecase.GetPlaylistsToAdd(ctx, usecaseRequest)
	if err != nil {
		logger.Error("failed to get playlists to add", zap.Error(err))
		json.WriteErrorResponse(w, errorStatus.ErrorStatus(err), err.Error(), nil)
		return
	}

	deliveryPlaylists := model.PlaylistsWithIsIncludedTrackFromUsecaseToDelivery(playlists)

	json.WriteSuccessResponse(w, http.StatusOK, deliveryPlaylists, nil)
}
