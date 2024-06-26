package srv

import (
	"net/http"

	"github.com/booleworks/logicng-service/computation"
	"github.com/booleworks/logicng-service/config"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func addRoutes(
	mux *http.ServeMux,
	cfg *config.Config,
) {
	mux.Handle("POST /assignment/{ass}", computation.HandleAssignment(cfg))
	mux.Handle("POST /bdd/compilation", computation.HandleBDDCompilation(cfg))
	mux.Handle("POST /bdd/graphical", computation.HandleBDDGraphical(cfg))
	mux.Handle("POST /dnnf/compilation", computation.HandleDNNFCompilation(cfg))
	mux.Handle("POST /encoding/{enc}", computation.HandleEncoding(cfg))
	mux.Handle("POST /explanation/mus", computation.HandleMUS(cfg))
	mux.Handle("POST /explanation/smus", computation.HandleSMUS(cfg))
	mux.Handle("POST /formula/{func}", computation.HandleFormula(cfg))
	mux.Handle("POST /graph/constraint", computation.HandleConstraintGraph(cfg))
	mux.Handle("POST /graph/constraint/graphical", computation.HandleConstraintGraphGraphical(cfg))
	mux.Handle("POST /graph/components", computation.HandleGraphComponents(cfg))
	mux.Handle("POST /model/counting", computation.HandleModelCounting(cfg))
	mux.Handle("POST /model/counting/projection", computation.HandleProjectedModelCounting(cfg))
	mux.Handle("POST /model/enumeration", computation.HandleModelEnumeration(cfg))
	mux.Handle("POST /model/enumeration/projection", computation.HandleProjectedModelEnumeration(cfg))
	mux.Handle("POST /normalform/transformation/{nf}", computation.HandleNFTrans(cfg))
	mux.Handle("POST /normalform/predicate/{nf}", computation.HandleNFPred(cfg))
	mux.Handle("POST /prime/minimal-implicant", computation.HandleMinimalImplicant(cfg))
	mux.Handle("POST /prime/minimal-cover", computation.HandleMinimalImplicantCover(cfg))
	mux.Handle("POST /simplification/{simp}", computation.HandleSimplification(cfg))
	mux.Handle("POST /solver/maxsat", computation.HandleMaxSat(cfg))
	mux.Handle("POST /solver/sat", computation.HandleSat(cfg))
	mux.Handle("POST /solver/predicate/{pred}", computation.HandleSatPredicate(cfg))
	mux.Handle("POST /solver/backbone", computation.HandleSatBackbone(cfg))
	mux.Handle("POST /substitution/{subst}", computation.HandleSubstitution(cfg))

	mux.Handle("GET /randomizer/{rand}", computation.HandleRandomizer(cfg))

	// Docs
	mux.HandleFunc("GET /swagger/*", httpSwagger.Handler(httpSwagger.URL("doc.json")))
}
