package srv

import (
	"fmt"
	"net/http"

	"github.com/booleworks/logicng-service/computation"
	"github.com/booleworks/logicng-service/config"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func addRoutes(
	mux *http.ServeMux,
	cfg *config.Config,
) {
	mux.Handle("POST /normalform/transformation/{nf}", computation.HandleNFTrans(cfg))
	mux.Handle("POST /normalform/predicate/{nf}", computation.HandleNFPred(cfg))
	mux.Handle("POST /simplification/{simp}", computation.HandleSimplification(cfg))
	mux.Handle("POST /substitution/{subst}", computation.HandleSubstitution(cfg))
	mux.Handle("GET /randomizer/{rand}", computation.HandleRandomizer(cfg))

	// Docs
	mux.HandleFunc("GET /swagger/*",
		httpSwagger.Handler(httpSwagger.URL(
			fmt.Sprintf("http://%s:%s/swagger/doc.json", cfg.Host, cfg.Port))))

	mux.HandleFunc("/docs",
		func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "static/redoc.html") })
}
