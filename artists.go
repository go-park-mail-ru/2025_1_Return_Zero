package main

import (
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/models"
)

// TODO: Write unit tests
// @Summary Get artists
// @Description Get a list of artists with optional pagination filters
// @Tags artists
// @Accept json
// @Produce json
// @Param offset query integer false "Offset (default: 0)"
// @Param limit query integer false "Limit (default: 10, max: 100)"
// @Success 200 {array} models.Artist "List of artists"
// @Failure 400 {string} string "Bad request - invalid filters"
// @Failure 500 {string} string "Internal server error"
// @Router /artists [get]
func (app *application) getArtists(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var filters models.Filters

	qs := r.URL.Query()
	offset, err := readInt(qs, "offset", DefaultOffset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	limit, err := readInt(qs, "limit", DefaultLimit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	filters.Offset = offset
	filters.Limit = limit

	if err := filters.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	artists, err := app.models.Artists.GetAll(filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	headers := make(http.Header)
	headers.Set("X-Total-Count", strconv.Itoa(len(artists)))

	if err := writeJSON(w, http.StatusOK, artists, headers); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
