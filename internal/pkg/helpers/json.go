package helpers

import (
	"encoding/json"
	"errors"
	"io"
	"maps"
	"net/http"

	deliveryModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/delivery"
	"go.uber.org/zap"
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

func WriteJSON(w http.ResponseWriter, status int, data interface{}, headers http.Header) {
	logger := zap.L().Sugar()
	jsonData, err := json.Marshal(data)
	if err != nil {
		logger.Error("failed to marshal json", zap.Error(err))
		return
	}

	maps.Copy(w.Header(), headers)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(jsonData)
	if err != nil {
		logger.Error("failed to write json", zap.Error(err))
	}
}

func WriteSuccessResponse(w http.ResponseWriter, status int, data interface{}, headers http.Header) {
	response := deliveryModel.APIResponse{
		Status: status,
		Body:   data,
	}

	WriteJSON(w, status, response, headers)
}

func WriteErrorResponse(w http.ResponseWriter, status int, message string, headers http.Header) {
	response := deliveryModel.APIErrorResponse{
		Status: status,
		Error:  message,
	}

	WriteJSON(w, status, response, headers)
}
