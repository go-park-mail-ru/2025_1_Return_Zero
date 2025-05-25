package label

import (
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/config"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/ctxExtractor"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/customErrors"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/errorStatus"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/json"
	loggerPkg "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/pagination"
	domain "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/label/domain"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/delivery"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type LabelHandler struct {
	usecase domain.Usecase
	cfg     *config.Config
}

func NewLabelHandler(usecase domain.Usecase, cfg *config.Config) *LabelHandler {
	return &LabelHandler{usecase: usecase, cfg: cfg}
}

// CreateLabel godoc
// @Summary Create a new label
// @Description Creates a new label in the system. Only accessible by administrators.
// @Tags label
// @Accept json
// @Produce json
// @Security AdminAuth
// @Param label body delivery.Label true "Label information"
// @Success 201 {object} delivery.Label "Created label"
// @Failure 400 {string} delivery.APIBadRequestErrorResponse "Bad request - invalid input"
// @Failure 401 {string} delivery.APIUnauthorizedErrorResponse "Unauthorized - admin access required"
// @Failure 500 {string} delivery.APIInternalServerErrorResponse "Internal server error"
// @Router /api/v1/label [post]
func (h *LabelHandler) CreateLabel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := loggerPkg.LoggerFromContext(ctx)
	isAdmin := ctxExtractor.AdminFromContext(ctx)
	if !isAdmin {
		logger.Error("Unauthorized access attempt")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var label *delivery.Label

	err := json.ReadJSON(w, r, &label)
	if err != nil {
		logger.Error("Failed to read JSON", zap.Error(err))
		json.WriteErrorResponse(w, http.StatusBadRequest, "Failed to read JSON", nil)
		return
	}
	labelUsecase := model.LabelFromDeliveryToUsecase(label)
	newLabel, err := h.usecase.CreateLabel(ctx, labelUsecase)
	if err != nil {
		json.WriteErrorResponse(w, errorStatus.ErrorStatus(err), err.Error(), nil)
		return
	}

	newLabelDelivery := model.LabelFromUsecaseToDelivery(newLabel)
	json.WriteJSON(w, http.StatusCreated, newLabelDelivery, nil)
}

// GetLabel godoc
// @Summary Get label by ID
// @Description Retrieves a label by its ID. Only accessible by administrators.
// @Tags label
// @Accept json
// @Produce json
// @Security AdminAuth
// @Param id path integer true "Label ID"
// @Success 200 {object} delivery.Label "Label information"
// @Failure 400 {string} delivery.APIBadRequestErrorResponse "Bad request - invalid label ID"
// @Failure 401 {string} delivery.APIUnauthorizedErrorResponse "Unauthorized - admin access required"
// @Failure 500 {string} delivery.APIInternalServerErrorResponse "Internal server error"
// @Router /api/v1/label/{id} [get]
func (h *LabelHandler) GetLabel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := loggerPkg.LoggerFromContext(ctx)
	isAdmin := ctxExtractor.AdminFromContext(ctx)
	if !isAdmin {
		logger.Error("Unauthorized access attempt")
		json.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	vars := mux.Vars(r)
	labelID := vars["id"]
	labelIDInt, err := strconv.ParseInt(labelID, 10, 64)
	if err != nil {
		logger.Error("Invalid label ID", zap.Error(err))
		json.WriteErrorResponse(w, http.StatusBadRequest, "Invalid label ID", nil)
		return
	}

	label, err := h.usecase.GetLabel(ctx, labelIDInt)
	if err != nil {
		json.WriteErrorResponse(w, errorStatus.ErrorStatus(err), err.Error(), nil)
		return
	}

	labelDelivery := model.LabelFromUsecaseToDelivery(label)
	json.WriteJSON(w, http.StatusOK, labelDelivery, nil)
}

// CreateArtist godoc
// @Summary Create a new artist
// @Description Creates a new artist associated with the label. Only accessible by label members.
// @Tags label
// @Accept multipart/form-data
// @Produce json
// @Security LabelAuth
// @Param title formData string true "Artist name"
// @Param thumbnail formData file true "Artist profile image (max 6MB)"
// @Success 201 {object} delivery.Artist "Created artist"
// @Failure 400 {object} delivery.APIBadRequestErrorResponse "Bad request - invalid input"
// @Failure 401 {object} delivery.APIUnauthorizedErrorResponse "Unauthorized"
// @Failure 403 {object} delivery.APIForbiddenErrorResponse "Forbidden - user not in label"
// @Failure 500 {object} delivery.APIInternalServerErrorResponse "Internal server error"
// @Router /api/v1/label/artist [post]
func (h *LabelHandler) CreateArtist(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := loggerPkg.LoggerFromContext(ctx)

	labelID, isLabel := ctxExtractor.LabelFromContext(ctx)
	if !isLabel {
		logger.Error("failed to get labelID")
		json.WriteErrorResponse(w, http.StatusForbidden, "user not in label", nil)
		return
	}

	if r.ContentLength > 6<<20 {
		logger.Error("request body too large")
		json.WriteErrorResponse(w, http.StatusBadRequest, "request body too large", nil)
		return
	}

	err := r.ParseMultipartForm(6 << 20)
	if err != nil {
		logger.Error("failed to parse multipart form", zap.Error(err))
		json.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	request := &delivery.CreateArtistRequest{}

	request.Title = r.FormValue("title")
	if request.Title == "" {
		logger.Error("title is empty")
		json.WriteErrorResponse(w, http.StatusBadRequest, "title is empty", nil)
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
		json.WriteErrorResponse(w, http.StatusBadRequest, "thumbnail image is required", nil)
		return
	}

	request.Image, err = io.ReadAll(file)
	if err != nil {
		logger.Error("failed to read thumbnail", zap.Error(err))
		json.WriteErrorResponse(w, errorStatus.ErrorStatus(customErrors.ErrPlaylistImageNotUploaded), customErrors.ErrPlaylistImageNotUploaded.Error(), nil)
		return
	}

	request.LabelID = labelID
	usecaseArtist := model.ArtistLoadFromDeliveryToUsecase(request)

	artist, err := h.usecase.CreateArtist(ctx, usecaseArtist)
	if err != nil {
		logger.Error("failed to create artist", zap.Error(err))
		json.WriteErrorResponse(w, errorStatus.ErrorStatus(err), err.Error(), nil)
		return
	}
	artistResponse := model.ArtistFromUsecaseToDelivery(artist)

	json.WriteSuccessResponse(w, http.StatusCreated, artistResponse, nil)
}

// EditArtist godoc
// @Summary Edit an existing artist
// @Description Updates an artist's information (name and/or image). Only accessible by label members.
// @Tags label
// @Accept multipart/form-data
// @Produce json
// @Security LabelAuth
// @Param artist_id formData string true "Artist ID to edit"
// @Param new_title formData string false "New artist name (required if thumbnail is not provided)"
// @Param thumbnail formData file false "New artist profile image (required if new_title is not provided, max 6MB)"
// @Success 200 {object} delivery.Artist "Updated artist"
// @Failure 400 {object} delivery.APIBadRequestErrorResponse "Bad request - invalid input"
// @Failure 401 {object} delivery.APIUnauthorizedErrorResponse "Unauthorized"
// @Failure 403 {object} delivery.APIForbiddenErrorResponse "Forbidden - user not in label"
// @Failure 500 {object} delivery.APIInternalServerErrorResponse "Internal server error"
// @Router /api/v1/label/artist [put]
func (h *LabelHandler) EditArtist(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := loggerPkg.LoggerFromContext(ctx)

	labelID, isLabel := ctxExtractor.LabelFromContext(ctx)
	if !isLabel {
		logger.Error("failed to get labelID")
		json.WriteErrorResponse(w, http.StatusForbidden, "user not in label", nil)
		return
	}

	if r.ContentLength > 6<<20 {
		logger.Error("request body too large")
		json.WriteErrorResponse(w, http.StatusBadRequest, "request body too large", nil)
		return
	}

	editRequest := &delivery.EditArtistRequest{}

	artistID := r.FormValue("artist_id")
	if artistID == "" {
		logger.Error("artist_id is empty")
		json.WriteErrorResponse(w, http.StatusBadRequest, "artist_id is empty", nil)
		return
	}

	artistIdInt, err := strconv.Atoi(artistID)
	if err != nil {
		logger.Error("failed to parse artist_id", zap.Error(err))
		json.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse artist_id", nil)
		return
	}
	artistIdInt64 := int64(artistIdInt)
	editRequest.ArtistID = artistIdInt64

	editRequest.NewTitle = r.FormValue("new_title")
	if len(editRequest.NewTitle) > 100 {
		logger.Error("new title is too long")
		json.WriteErrorResponse(w, http.StatusBadRequest, "new title is too long", nil)
		return
	}

	file, _, err := r.FormFile("thumbnail")
	if err == nil {
		editRequest.Image, err = io.ReadAll(file)
		if err != nil {
			logger.Error("failed to read thumbnail", zap.Error(err))
			json.WriteErrorResponse(w, errorStatus.ErrorStatus(customErrors.ErrPlaylistImageNotUploaded), customErrors.ErrPlaylistImageNotUploaded.Error(), nil)
			return
		}
	} else {
		editRequest.Image = nil
	}
	editRequest.LabelID = labelID

	if editRequest.NewTitle == "" && editRequest.Image == nil {
		logger.Error("new title and image are empty")
		json.WriteErrorResponse(w, http.StatusBadRequest, "new title and image are empty", nil)
		return
	}

	artistEditUsecase := model.ArtistEditFromDeliveryToUsecase(editRequest)
	artist, err := h.usecase.EditArtist(ctx, artistEditUsecase)
	if err != nil {
		logger.Error("failed to edit artist", zap.Error(err))
		json.WriteErrorResponse(w, errorStatus.ErrorStatus(err), err.Error(), nil)
		return
	}
	artistResponse := model.ArtistFromUsecaseToDelivery(artist)
	json.WriteSuccessResponse(w, http.StatusOK, artistResponse, nil)
}

// GetArtists godoc
// @Summary Get artists in a label
// @Description Retrieves a list of artists associated with the user's label with optional pagination
// @Tags label
// @Accept json
// @Produce json
// @Security LabelAuth
// @Param offset query integer false "Offset (default: 0)"
// @Param limit query integer false "Limit (default: 10, max: 100)"
// @Success 200 {object} delivery.APIResponse{body=[]delivery.Artist} "List of artists in the label"
// @Failure 400 {object} delivery.APIBadRequestErrorResponse "Bad request - invalid pagination parameters"
// @Failure 401 {object} delivery.APIUnauthorizedErrorResponse "Unauthorized"
// @Failure 403 {object} delivery.APIForbiddenErrorResponse "Forbidden - user not in label"
// @Failure 500 {object} delivery.APIInternalServerErrorResponse "Internal server error"
// @Router /api/v1/label/artists [get]
func (h *LabelHandler) GetArtists(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := loggerPkg.LoggerFromContext(ctx)

	labelID, isLabel := ctxExtractor.LabelFromContext(ctx)
	if !isLabel {
		logger.Error("failed to get labelID")
		json.WriteErrorResponse(w, http.StatusForbidden, "user not in label", nil)
		return
	}

	pagination, err := pagination.GetPagination(r, &h.cfg.Pagination)
	if err != nil {
		logger.Error("failed to get pagination", zap.Error(err))
		json.WriteErrorResponse(w, errorStatus.ErrorStatus(err), err.Error(), nil)
		return
	}

	usecaseArtists, err := h.usecase.GetArtists(ctx, labelID, &usecase.ArtistFilters{
		Pagination: model.PaginationFromDeliveryToUsecase(pagination),
	})
	if err != nil {
		logger.Error("failed to get artists", zap.Error(err))
		json.WriteErrorResponse(w, errorStatus.ErrorStatus(err), err.Error(), nil)
		return
	}

	artistsDelivery := model.ArtistsFromUsecaseToDelivery(usecaseArtists)
	json.WriteSuccessResponse(w, http.StatusOK, artistsDelivery, nil)
}

// DeleteArtist godoc
// @Summary Delete an artist
// @Description Removes an artist from the label. Only accessible by label members.
// @Tags label
// @Accept json
// @Produce json
// @Security LabelAuth
// @Param request body delivery.DeleteArtistRequest true "Artist deletion request containing title"
// @Success 200 {object} delivery.DeleteArtistRequest "Artist deleted successfully"
// @Failure 400 {object} delivery.APIBadRequestErrorResponse "Bad request - invalid input"
// @Failure 401 {object} delivery.APIUnauthorizedErrorResponse "Unauthorized"
// @Failure 403 {object} delivery.APIForbiddenErrorResponse "Forbidden - user not in label"
// @Failure 500 {object} delivery.APIInternalServerErrorResponse "Internal server error"
// @Router /api/v1/label/artist [delete]
func (h *LabelHandler) DeleteArtist(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := loggerPkg.LoggerFromContext(ctx)

	labelID, isLabel := ctxExtractor.LabelFromContext(ctx)
	if !isLabel {
		logger.Error("failed to get labelID")
		json.WriteErrorResponse(w, http.StatusForbidden, "user not in label", nil)
		return
	}

	deleteArtistRequest := &delivery.DeleteArtistRequest{}

	err := json.ReadJSON(w, r, &deleteArtistRequest)
	if err != nil {
		logger.Error("failed to read JSON", zap.Error(err))
		json.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	deleteArtistRequest.LabelID = labelID

	err = h.usecase.DeleteArtist(ctx, model.ArtistDeleteFromDeliveryToUsecase(deleteArtistRequest))
	if err != nil {
		logger.Error("failed to delete artist", zap.Error(err))
		json.WriteErrorResponse(w, errorStatus.ErrorStatus(err), err.Error(), nil)
		return
	}

	json.WriteSuccessResponse(w, http.StatusOK, deleteArtistRequest, nil)
}

// CreateAlbum godoc
// @Summary Create a new album
// @Description Creates a new album with tracks. Only accessible by label members.
// @Tags label
// @Accept multipart/form-data
// @Produce json
// @Security LabelAuth
// @Param title formData string true "Album title (max 100 characters)"
// @Param type formData string true "Album type (album, single, ep, compilation)"
// @Param artists_ids formData string true "Comma-separated list of artist IDs"
// @Param thumbnail formData file true "Album cover image"
// @Param tracks[] formData file true "Array of track files"
// @Param track_titles[] formData []string true "Array of track titles corresponding to tracks[]"
// @Success 201 {object} delivery.SuccessCreateAlbum "Created album ID"
// @Failure 400 {object} delivery.APIBadRequestErrorResponse "Bad request - invalid input"
// @Failure 401 {object} delivery.APIUnauthorizedErrorResponse "Unauthorized"
// @Failure 403 {object} delivery.APIForbiddenErrorResponse "Forbidden - user not in label"
// @Failure 500 {object} delivery.APIInternalServerErrorResponse "Internal server error"
// @Router /api/v1/label/album [post]
func (h *LabelHandler) CreateAlbum(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger := loggerPkg.LoggerFromContext(ctx)

	labelID, isLabel := ctxExtractor.LabelFromContext(ctx)
	if !isLabel {
		logger.Error("failed to get labelID")
		json.WriteErrorResponse(w, http.StatusForbidden, "user not in label", nil)
		return
	}

	if r.ContentLength > 300<<20 {
		logger.Error("request body too large")
		json.WriteErrorResponse(w, http.StatusBadRequest, "request body too large", nil)
		return
	}

	err := r.ParseMultipartForm(300 << 20)
	if err != nil {
		logger.Error("failed to parse multipart form", zap.Error(err))
		json.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse multipart form", nil)
		return
	}

	request := &delivery.CreateAlbumRequest{}
	request.LabelID = labelID

	request.Title = r.FormValue("title")
	if request.Title == "" {
		logger.Error("title is empty")
		json.WriteErrorResponse(w, http.StatusBadRequest, "title is empty", nil)
		return
	}
	if len(request.Title) > 100 {
		logger.Error("title is too long")
		json.WriteErrorResponse(w, http.StatusBadRequest, "title is too long", nil)
		return
	}

	request.Type = r.FormValue("type")
	if request.Type == "" {
		logger.Error("type is empty")
		json.WriteErrorResponse(w, http.StatusBadRequest, "type is empty", nil)
		return
	}
	if request.Type != "album" && request.Type != "single" && request.Type != "ep" && request.Type != "compilation" {
		logger.Error("type is invalid")
		json.WriteErrorResponse(w, http.StatusBadRequest, "type is invalid", nil)
		return
	}

	ArtistIDs := r.FormValue("artists_ids")
	ArtistIDsSlice := strings.Split(ArtistIDs, ",")
	for _, id := range ArtistIDsSlice {
		parsedID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			logger.Error("failed to parse artists_ids", zap.Error(err))
			json.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse artists_ids", nil)
			return
		}
		request.ArtistsIDs = append(request.ArtistsIDs, parsedID)
	}

	file, _, err := r.FormFile("thumbnail")
	if err != nil {
		logger.Error("failed to get thumbnail", zap.Error(err))
		json.WriteErrorResponse(w, http.StatusBadRequest, "thumbnail image is required", nil)
		return
	}
	request.Image, err = io.ReadAll(file)
	if err != nil {
		logger.Error("failed to read thumbnail", zap.Error(err))
		json.WriteErrorResponse(w, http.StatusBadRequest, "failed to read thumbnail", nil)
		return
	}

	fileHeaders := r.MultipartForm.File["tracks[]"]
	if len(fileHeaders) == 0 {
		logger.Error("tracks are empty")
		json.WriteErrorResponse(w, http.StatusBadRequest, "tracks are empty", nil)
		return
	}

	tracks := make([]*delivery.CreateTrackRequest, 0, len(fileHeaders))
	trackTitles := r.Form["track_titles[]"]
	if len(trackTitles) != len(fileHeaders) {
		logger.Error("track titles count does not match tracks count")
		json.WriteErrorResponse(w, http.StatusBadRequest, "track titles count does not match tracks count", nil)
		return
	}

	for i, fileHeader := range fileHeaders {
		track := &delivery.CreateTrackRequest{
			Title: trackTitles[i],
		}

		if fileHeader.Size > 50<<20 {
			logger.Error("track is too large")
			json.WriteErrorResponse(w, http.StatusBadRequest, "track is too large", nil)
			return
		}

		file, err := fileHeader.Open()
		if err != nil {
			logger.Error("failed to open track file", zap.Error(err))
			json.WriteErrorResponse(w, http.StatusBadRequest, "failed to open track file", nil)
			return
		}
		defer file.Close()

		track.Track, err = io.ReadAll(file)
		if err != nil {
			logger.Error("failed to read track file", zap.Error(err))
			json.WriteErrorResponse(w, http.StatusBadRequest, "failed to read track file", nil)
			return
		}

		tracks = append(tracks, track)
	}

	request.Tracks = tracks

	albumID, albumURL, err := h.usecase.CreateAlbum(ctx, model.NewAlbumFromDeliveryToUsecase(request))
	if err != nil {
		logger.Error("failed to create album", zap.Error(err))
		json.WriteErrorResponse(w, errorStatus.ErrorStatus(err), err.Error(), nil)
		return
	}

	artists := make([]*delivery.AlbumArtist, 0, len(request.ArtistsIDs))
	for _, artistID := range request.ArtistsIDs {
		artist := &delivery.AlbumArtist{
			ID:    artistID,
			Title: "",
		}
		artists = append(artists, artist)
	}

	deliveryAlbum := delivery.Album{
		ID:          albumID,
		Title:       request.Title,
		Type:        model.AlbumTypeConverter(request.Type),
		Thumbnail:   albumURL,
		ReleaseDate: time.Now(),
		Artists:     artists,
	}

	json.WriteSuccessResponse(w, http.StatusCreated, deliveryAlbum, nil)
}

// UpdateLabel godoc
// @Summary Update a label
// @Description Updates a label by adding or removing users. Only accessible by administrators.
// @Tags label
// @Accept json
// @Produce json
// @Security AdminAuth
// @Param request body delivery.EditLabelRequest true "Label update information containing labelID, users to add, and users to remove"
// @Success 200 {object} delivery.Message "Label edited successfully"
// @Failure 400 {object} delivery.APIBadRequestErrorResponse "Bad request - invalid input"
// @Failure 401 {object} delivery.APIUnauthorizedErrorResponse "Unauthorized - admin access required"
// @Failure 500 {object} delivery.APIInternalServerErrorResponse "Internal server error"
// @Router /api/v1/label [put]
func (h *LabelHandler) UpdateLabel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := loggerPkg.LoggerFromContext(ctx)

	isAdmin := ctxExtractor.AdminFromContext(ctx)
	if !isAdmin {
		logger.Error("Unauthorized access attempt")
		json.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	var label *delivery.EditLabelRequest

	err := json.ReadJSON(w, r, &label)
	if err != nil {
		logger.Error("Failed to read JSON", zap.Error(err))
		json.WriteErrorResponse(w, http.StatusBadRequest, "Failed to read JSON", nil)
		return
	}

	err = h.usecase.UpdateLabel(ctx, label.LabelID, label.ToAdd, label.ToRemove)
	if err != nil {
		logger.Error("Failed to update label", zap.Error(err))
		json.WriteErrorResponse(w, errorStatus.ErrorStatus(err), err.Error(), nil)
		return
	}

	json.WriteJSON(w, http.StatusOK, delivery.Message{Message: "Label edited successfully"}, nil)
}

// DeleteAlbum godoc
// @Summary Delete an album
// @Description Removes an album from the label. Accessible by label members for their own label or administrators for any label.
// @Tags label
// @Accept json
// @Produce json
// @Security LabelAuth
// @Security AdminAuth
// @Param request body delivery.DeleteAlbumRequest true "Album deletion request containing labelID and albumID"
// @Success 200 {object} delivery.Message "Album deleted successfully"
// @Failure 400 {string} delivery.APIBadRequestErrorResponse "Bad request - invalid input or label ID mismatch"
// @Failure 401 {string} delivery.APIUnauthorizedErrorResponse "Unauthorized"
// @Failure 403 {string} delivery.APIForbiddenErrorResponse "Forbidden - user not in label"
// @Failure 500 {string} delivery.APIInternalServerErrorResponse "Internal server error"
// @Router /api/v1/label/album [delete]
func (h *LabelHandler) DeleteAlbum(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := loggerPkg.LoggerFromContext(ctx)

	isAdmin := ctxExtractor.AdminFromContext(ctx)
	labelID, isLabel := ctxExtractor.LabelFromContext(ctx)
	if !isLabel && !isAdmin {
		logger.Error("failed to authorize user")
		json.WriteErrorResponse(w, http.StatusForbidden, "user not in label", nil)
		return
	}

	var req delivery.DeleteAlbumRequest
	err := json.ReadJSON(w, r, &req)
	if err != nil {
		logger.Error("Failed to read JSON", zap.Error(err))
		json.WriteErrorResponse(w, http.StatusBadRequest, "Failed to read JSON", nil)
		return
	}
	if isLabel {
		if req.LabelID != labelID {
			logger.Error("Label ID mismatch", zap.Error(err))
			json.WriteErrorResponse(w, http.StatusBadRequest, "Label ID mismatch", nil)
			return
		}
	}

	err = h.usecase.DeleteAlbum(ctx, req.LabelID, req.AlbumID)
	if err != nil {
		logger.Error("Failed to delete album", zap.Error(err))
		json.WriteErrorResponse(w, errorStatus.ErrorStatus(err), err.Error(), nil)
		return
	}

	json.WriteJSON(w, http.StatusOK, delivery.Message{Message: "Album deleted successfully"}, nil)
}

// GetAlbumsByLabelID godoc
// @Summary Get albums by label ID
// @Description Retrieves a list of albums associated with the user's label with optional pagination.
// @Tags label
// @Accept json
// @Produce json
// @Security LabelAuth
// @Param offset query integer false "Offset (default: 0)"
// @Param limit query integer false "Limit (default: 10, max: 100)"
// @Success 200 {object} delivery.APIResponse{body=[]delivery.Album} "List of albums in the label"
// @Failure 400 {object} delivery.APIBadRequestErrorResponse "Bad request - invalid pagination parameters"
// @Failure 401 {object} delivery.APIUnauthorizedErrorResponse "Unauthorized"
// @Failure 403 {object} delivery.APIForbiddenErrorResponse "Forbidden - user not in label"
// @Failure 500 {object} delivery.APIInternalServerErrorResponse "Internal server error"
// @Router /api/v1/label/albums [get]
func (h *LabelHandler) GetAlbumsByLabelID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := loggerPkg.LoggerFromContext(ctx)

	labelID, isLabel := ctxExtractor.LabelFromContext(ctx)
	if !isLabel {
		logger.Error("failed to get labelID")
		json.WriteErrorResponse(w, http.StatusForbidden, "user not in label", nil)
		return
	}

	pagination, err := pagination.GetPagination(r, &h.cfg.Pagination)
	if err != nil {
		logger.Error("failed to get pagination", zap.Error(err))
		json.WriteErrorResponse(w, errorStatus.ErrorStatus(err), err.Error(), nil)
		return
	}

	albums, err := h.usecase.GetAlbumsByLabelID(ctx, labelID, &usecase.AlbumFilters{
		Pagination: model.PaginationFromDeliveryToUsecase(pagination),
	})
	if err != nil {
		logger.Error("failed to get albums", zap.Error(err))
		json.WriteErrorResponse(w, errorStatus.ErrorStatus(err), err.Error(), nil)
		return
	}

	albumsDelivery := model.AlbumsFromUsecaseToDelivery(albums)
	json.WriteSuccessResponse(w, http.StatusOK, albumsDelivery, nil)
}
