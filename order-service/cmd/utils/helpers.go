package utils

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576 // one megabyte

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must have only a single JSON value")
	}
	return nil
}

func WriteJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}
	return nil
}

func ErrorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}
	var payload jsonResponse
	payload.Error = true
	payload.Message = err.Error()

	return WriteJSON(w, statusCode, payload)
}

func GetID(r *http.Request, key string) (uint, error) {
	id := chi.URLParam(r, key)
	if id == "" {
		return 0, errors.New("missing ID parameter")
	}

	idInt, err := strconv.ParseUint(id, 10, 0)
	if err != nil {
		return 0, errors.New("ID parameter must be a number")
	}

	return uint(idInt), nil
}
func GetString(r *http.Request, key string) (string, error) {
	value := chi.URLParam(r, key)
	if value == "" {
		return value, errors.New("missing  parameter")
	}

	return value, nil
}
func ItoS(i uint) string {
	return strconv.Itoa(int(i))
}

func StoUint(s string) uint {
	i, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0
	}
	return uint(i)
}
