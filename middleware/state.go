package middleware

import (
	"context"
	"net/http"

	"github.com/booleworks/logicng-service/sio"
)

func AddState(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), sio.State{}, &sio.ComputationState{Success: true})
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}
