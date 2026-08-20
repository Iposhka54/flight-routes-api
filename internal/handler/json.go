package handler

import (
	"encoding/json"
	apperr "flight-routes-api/internal/error"
	"fmt"
	"net/http"
)

func writeJSON(w http.ResponseWriter, status int, data any) error {
	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("%w: %v", apperr.ErrEncodeJSON, err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err = w.Write(body); err != nil {
		return fmt.Errorf("%w: %v", apperr.ErrEncodeJSON, err)
	}

	return nil
}
