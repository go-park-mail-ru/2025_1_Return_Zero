package main

import (
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/models"
)

// TODO: Write unit tests
// @Summary Get tracks
// @Description Get a list of tracks with optional pagination filters
// @Tags tracks
// @Accept json
// @Produce json
// @Param page query integer false "Page number (default: 1)"
// @Param page_size query integer false "Page size (default: 10, max: 100)"
// @Success 200 {array} models.Track "List of tracks"
// @Failure 400 {string} string "Bad request - invalid filters"
// @Failure 404 {string} string "No tracks found"
// @Failure 500 {string} string "Internal server error"
// @Router /tracks [get]
func getTracksHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		models.Filters
	}

	var err error

	qs := r.URL.Query()
	input.Page, err = readInt(qs, "page", 1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	input.PageSize, err = readInt(qs, "page_size", 10)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := input.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	model := models.NewTracksModel()
	tracks, err := model.GetAll(input.Filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	headers := make(http.Header)
	headers.Set("X-Total-Count", strconv.Itoa(len(tracks)))

	err = writeJSON(w, http.StatusOK, tracks, headers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
