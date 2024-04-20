package middleware

import (
	"context"
	"net/http"

	"github.com/booleworks/logicng-service/slog"
	"github.com/google/uuid"
)

const CorrIdHeader = "x-Corrrelation-id"

func CorrelationId(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		corrId := r.Header.Get(CorrIdHeader)
		if corrId == "" {
			corrId = uuid.NewString()
		}
		w.Header().Set(CorrIdHeader, corrId)
		ctx := context.WithValue(r.Context(), slog.ID{}, corrId)
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}
