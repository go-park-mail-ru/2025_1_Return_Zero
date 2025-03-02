package main

import (
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/models"
)

type AlbumsHandler struct {
	Model *models.AlbumsModel
}

// List TODO: Write unit tests
// @Summary Get albums
// @Description Get a list of albums with optional pagination filters
// @Tags albums
// @Accept json
// @Produce json
// @Param offset query integer false "Offset (default: 0)"
// @Param limit query integer false "Limit (default: 10, max: 100)"
// @Success 200 {array} models.Album "List of albums"
// @Failure 400 {string} string "Bad request - invalid filters"
// @Failure 500 {string} string "Internal server error"
// @Router /albums [get]
func (handler *AlbumsHandler) List(w http.ResponseWriter, r *http.Request) {
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

	filters.Limit = limit
	filters.Offset = offset

	if err := filters.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	albums := handler.Model.GetAll(filters)

	headers := make(http.Header)
	headers.Set("X-Total-Count", strconv.Itoa(len(albums)))

	if err := writeJSON(w, http.StatusOK, albums, headers); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
