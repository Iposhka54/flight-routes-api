package middleware

import (
	"context"
	"encoding/json"
	"errors"
	apperr "flight-routes-api/internal/error"
	"net/http"

	"go.uber.org/zap"
)

type HandlerFunc func(http.ResponseWriter, *http.Request) error

type errorResponse struct {
	Error string `json:"error"`
}

func ErrorHandler(log *zap.Logger, next HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw := &responseWriter{ResponseWriter: w}

		err := next(rw, r)
		if err == nil || rw.wrote {
			return
		}

		if errors.Is(err, context.Canceled) {
			if log != nil {
				log.Info("request canceled",
					zap.Error(err),
					zap.String("method", r.Method),
					zap.String("path", r.URL.Path),
				)
			}
			return
		}

		status, message := mapError(err)
		if log != nil {
			fields := []zap.Field{
				zap.Error(err),
				zap.Int("status", status),
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
			}
			if status == http.StatusInternalServerError {
				log.Error("request failed", fields...)
			} else {
				log.Info("request failed", fields...)
			}
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(status)
		_ = json.NewEncoder(rw).Encode(errorResponse{Error: message})
	}
}

func mapError(err error) (int, string) {
	if errors.Is(err, context.DeadlineExceeded) {
		return http.StatusGatewayTimeout, "request timeout"
	}

	var appErr *apperr.Error
	if errors.As(err, &appErr) {
		return appErr.Status, appErr.Message
	}

	return apperr.ErrInternal.Status, apperr.ErrInternal.Message
}

type responseWriter struct {
	http.ResponseWriter
	wrote bool
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.wrote = true
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.wrote = true
	return w.ResponseWriter.Write(b)
}
