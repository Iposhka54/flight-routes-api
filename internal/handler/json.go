package handler

import (
	"encoding/json"
	apperr "flight-routes-api/internal/error"
	"net/http"
)

func writeJSON(w http.ResponseWriter, status int, data any) error {
	body, err := json.Marshal(data)
	if err != nil {
		return apperr.ErrEncodeJSON.Wrap(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err = w.Write(body); err != nil {
		return apperr.ErrEncodeJSON.Wrap(err)
	}

	return nil
}
