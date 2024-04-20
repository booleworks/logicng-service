package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/booleworks/logicng-service/sio"
	"github.com/booleworks/logicng-service/slog"
)

func PerformanceLogger(handler http.Handler, logger *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info(logger, r, "called %s", slog.Style(r.URL.String(), slog.StyleUnderline))
		start := time.Now()
		handler.ServeHTTP(w, r)
		elapsed := time.Since(start)
		state := r.Context().Value(sio.State{}).(*sio.ComputationState)
		if state.Error != "" {
			slog.Error(logger, r, state.Error)
			slog.Info(logger, r, "computation with error %s(%s)%s", slog.StyleGreen, elapsed, slog.StyleEnd)
		} else {
			slog.Info(logger, r, "computation with success %s(%s)%s", slog.StyleGreen, elapsed, slog.StyleEnd)
		}
	})
}
