package computation

import (
	"net/http"

	"github.com/booleworks/logicng-go/bdd"
	"github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/handler"
	"github.com/booleworks/logicng-go/model/enum"
	"github.com/booleworks/logicng-go/model/iter"
	"github.com/booleworks/logicng-go/normalform"
	"github.com/booleworks/logicng-service/config"
	"github.com/booleworks/logicng-service/sio"
)

func HandleNFTrans(cfg *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		normalform := r.PathValue("nf")
		switch normalform {
		case "nnf":
			handleNFTransNNF(w, r)
		case "cnf":
			handleNFTransCNF(w, r, cfg)
		case "dnf":
			handleNFTransDNF(w, r, cfg)
		case "aig":
			handleNFTransAIG(w, r)
		default:
			sio.WriteError(w, r, sio.ErrUnknownPath(r.URL.Path))
		}
	})
}

// @Summary      Transform a formula to negation normal form
// @Tags         Normal Form
// @Param        request body	sio.FormulaInput true "Input Formula"
// @Success      200  {object}  sio.FormulaResult
// @Router       /normalform/transformation/nnf [post]
func handleNFTransNNF(w http.ResponseWriter, r *http.Request) {
	transform(w, r, func(fac formula.Factory, form formula.Formula) (formula.Formula, sio.ServiceError) {
		return normalform.NNF(fac, form), nil
	})
}

// @Summary      Transform a formula to conjunctive normal form
// @Tags         Normal Form
// @Param        algorithm query string  false "CNF Algorithm" Enums(advanced, tseitin, pg, factorization, canonical, bdd)
// @Param        request body	sio.FormulaInput true "Input Formula"
// @Success      200  {object}  sio.FormulaResult
// @Router       /normalform/transformation/cnf [post]
func handleNFTransCNF(w http.ResponseWriter, r *http.Request, cfg *config.Config) {
	algorithm := r.URL.Query().Get("algorithm")
	var method func(formula.Factory, formula.Formula) (formula.Formula, sio.ServiceError)
	switch algorithm {
	case "advanced", "":
		method = func(fac formula.Factory, f formula.Formula) (formula.Formula, sio.ServiceError) {
			return normalform.CNF(fac, f), nil
		}
	case "tseitin":
		method = func(fac formula.Factory, f formula.Formula) (formula.Formula, sio.ServiceError) {
			return normalform.TseitinCNFWithBoundary(fac, f, 3), nil
		}
	case "factorization":
		method = func(fac formula.Factory, f formula.Formula) (formula.Formula, sio.ServiceError) {
			hdl := normalform.HandlerWithTimeout(*handler.NewTimeoutWithDuration(cfg.SyncComputationTimout))
			result, ok := normalform.FactorizedCNFWithHandler(fac, f, hdl)
			return transformWithTimeout(result, ok)
		}
	case "pg":
		method = func(fac formula.Factory, f formula.Formula) (formula.Formula, sio.ServiceError) {
			return normalform.PGCNFWithBoundary(fac, f, 3), nil
		}
	case "canonical":
		method = func(fac formula.Factory, f formula.Formula) (formula.Formula, sio.ServiceError) {
			hdl := iter.HandlerWithTimeout(*handler.NewTimeoutWithDuration(cfg.SyncComputationTimout))
			result, ok := enum.CanonicalCNFWithHandler(fac, f, hdl)
			return transformWithTimeout(result, ok)
		}
	case "bdd":
		method = func(fac formula.Factory, f formula.Formula) (formula.Formula, sio.ServiceError) {
			hdl := bdd.HandlerWithTimeout(*handler.NewTimeoutWithDuration(cfg.SyncComputationTimout))
			result, ok := bdd.CNFWithHandler(fac, f, hdl)
			return transformWithTimeout(result, ok)
		}
	}
	transform(w, r, method)
}

// @Summary      Transform a formula to disjunctive normal form
// @Tags         Normal Form
// @Param        algorithm query string false "DNF Algorithm" Enums(factorization, canonical, bdd)
// @Param        request body	sio.FormulaInput true "Input Formula"
// @Success      200  {object}  sio.FormulaResult
// @Router       /normalform/transformation/dnf [post]
func handleNFTransDNF(w http.ResponseWriter, r *http.Request, cfg *config.Config) {
	algorithm := r.URL.Query().Get("algorithm")
	var method func(formula.Factory, formula.Formula) (formula.Formula, sio.ServiceError)
	switch algorithm {
	case "factorization", "":
		method = func(fac formula.Factory, f formula.Formula) (formula.Formula, sio.ServiceError) {
			hdl := normalform.HandlerWithTimeout(*handler.NewTimeoutWithDuration(cfg.SyncComputationTimout))
			result, ok := normalform.FactorizedDNFWithHandler(fac, f, hdl)
			return transformWithTimeout(result, ok)
		}
	case "canonical":
		method = func(fac formula.Factory, f formula.Formula) (formula.Formula, sio.ServiceError) {
			hdl := iter.HandlerWithTimeout(*handler.NewTimeoutWithDuration(cfg.SyncComputationTimout))
			result, ok := enum.CanonicalDNFWithHandler(fac, f, hdl)
			return transformWithTimeout(result, ok)
		}
	case "bdd":
		method = func(fac formula.Factory, f formula.Formula) (formula.Formula, sio.ServiceError) {
			hdl := bdd.HandlerWithTimeout(*handler.NewTimeoutWithDuration(cfg.SyncComputationTimout))
			result, ok := bdd.DNFWithHandler(fac, f, hdl)
			return transformWithTimeout(result, ok)
		}
	}
	transform(w, r, method)
}

// @Summary      Transform a formula to an and-inverter-graph
// @Tags         Normal Form
// @Param        request body	sio.FormulaInput true "Input Formula"
// @Success      200  {object}  sio.FormulaResult
// @Router       /normalform/transformation/aig [post]
func handleNFTransAIG(w http.ResponseWriter, r *http.Request) {
	transform(w, r, func(fac formula.Factory, form formula.Formula) (formula.Formula, sio.ServiceError) {
		return normalform.AIG(fac, form), nil
	})
}

// @Summary      Report whether a formula is an a certain normal form
// @Tags         Normal Form
// @Param        nf path string true "normal form" Enums(nnf, cnf, dnf, aig, minterm, maxterm)
// @Param        request body	sio.FormulaInput true "Input Formula"
// @Success      200  {object}  sio.BoolResult
// @Router       /normalform/predicate/{nf} [post]
func HandleNFPred(cfg *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nf := r.PathValue("nf")
		var predicate formula.Predicate
		switch nf {
		case "nnf":
			predicate = normalform.IsNNF
		case "cnf":
			predicate = normalform.IsCNF
		case "dnf":
			predicate = normalform.IsDNF
		case "aig":
			predicate = normalform.IsAIG
		case "minterm":
			predicate = normalform.IsMinterm
		case "maxterm":
			predicate = normalform.IsMaxterm
		default:
			sio.WriteError(w, r, sio.ErrUnknownPath(r.URL.Path))
			return
		}
		holds(w, r, predicate)
	})
}
