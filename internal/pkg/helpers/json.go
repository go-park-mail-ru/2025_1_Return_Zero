package helpers

import (
	"encoding/json"
	"errors"
	"io"
	"maps"
	"net/http"
)

const (
	MaxBytes      = 1024 * 1024
	DefaultStatus = http.StatusOK
)

var (
	ErrMultipleJSONValues = errors.New("body must only contain a single JSON value")
)

func ReadJSON(w http.ResponseWriter, r *http.Request, v interface{}) error {
	maxBytes := int64(MaxBytes)
	r.Body = http.MaxBytesReader(w, r.Body, maxBytes)

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(v); err != nil {
		return err
	}

	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return ErrMultipleJSONValues
	}

	return nil
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}, headers http.Header) error {
	response := struct {
		Status int         `json:"status"`
		Body   interface{} `json:"body"`
	}{
		Status: status,
		Body:   data,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		return err
	}

	maps.Copy(w.Header(), headers)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(DefaultStatus)

	_, err = w.Write(jsonData)

	return err
}

func WriteJSONError(w http.ResponseWriter, status int, message string) error {
	return WriteJSON(w, status, struct {
		Error string `json:"error"`
	}{
		Error: message,
	}, nil)
}