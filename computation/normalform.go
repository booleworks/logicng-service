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
// @Description  If a list of formulas is given, the normal form is computed for the conjunction of these formulas.  The result always contains exactly one formula.
// @Tags         Normal Form
// @Param        request body	sio.FormulaInput true "Input formulas"
// @Success      200  {object}  sio.FormulaResult
// @Router       /normalform/transformation/nnf [post]
func handleNFTransNNF(w http.ResponseWriter, r *http.Request) {
	transform(w, r, func(fac formula.Factory, form []formula.Formula) (formula.Formula, sio.ServiceError) {
		return normalform.NNF(fac, fac.And(form...)), nil
	})
}

// @Summary      Transform a formula to conjunctive normal form
// @Description  If a list of formulas is given, the normal form is computed for the conjunction of these formulas.  The result always contains exactly one formula.
// @Tags         Normal Form
// @Param        algorithm query string  false "CNF Algorithm" Enums(advanced, tseitin, pg, factorization, canonical, bdd)
// @Param        request body	sio.FormulaInput true "Input formulas"
// @Success      200  {object}  sio.FormulaResult
// @Router       /normalform/transformation/cnf [post]
func handleNFTransCNF(w http.ResponseWriter, r *http.Request, cfg *config.Config) {
	algorithm := r.URL.Query().Get("algorithm")
	var method func(formula.Factory, []formula.Formula) (formula.Formula, sio.ServiceError)
	switch algorithm {
	case "advanced", "":
		method = func(fac formula.Factory, f []formula.Formula) (formula.Formula, sio.ServiceError) {
			return normalform.CNF(fac, fac.And(f...)), nil
		}
	case "tseitin":
		method = func(fac formula.Factory, f []formula.Formula) (formula.Formula, sio.ServiceError) {
			return normalform.TseitinCNFWithBoundary(fac, fac.And(f...), 3), nil
		}
	case "factorization":
		method = func(fac formula.Factory, f []formula.Formula) (formula.Formula, sio.ServiceError) {
			hdl := normalform.HandlerWithTimeout(*handler.NewTimeoutWithDuration(cfg.SyncComputationTimout))
			result, ok := normalform.FactorizedCNFWithHandler(fac, fac.And(f...), hdl)
			return transformWithTimeout(result, ok)
		}
	case "pg":
		method = func(fac formula.Factory, f []formula.Formula) (formula.Formula, sio.ServiceError) {
			return normalform.PGCNFWithBoundary(fac, fac.And(f...), 3), nil
		}
	case "canonical":
		method = func(fac formula.Factory, f []formula.Formula) (formula.Formula, sio.ServiceError) {
			hdl := iter.HandlerWithTimeout(*handler.NewTimeoutWithDuration(cfg.SyncComputationTimout))
			result, ok := enum.CanonicalCNFWithHandler(fac, fac.And(f...), hdl)
			return transformWithTimeout(result, ok)
		}
	case "bdd":
		method = func(fac formula.Factory, f []formula.Formula) (formula.Formula, sio.ServiceError) {
			hdl := bdd.HandlerWithTimeout(*handler.NewTimeoutWithDuration(cfg.SyncComputationTimout))
			result, ok := bdd.CNFWithHandler(fac, fac.And(f...), hdl)
			return transformWithTimeout(result, ok)
		}
	}
	transform(w, r, method)
}

// @Summary      Transform a formula to disjunctive normal form
// @Description  If a list of formulas is given, the normal form is computed for the conjunction of these formulas.  The result always contains exactly one formula.
// @Tags         Normal Form
// @Param        algorithm query string false "DNF Algorithm" Enums(factorization, canonical, bdd)
// @Param        request body	sio.FormulaInput true "Input formulas"
// @Success      200  {object}  sio.FormulaResult
// @Router       /normalform/transformation/dnf [post]
func handleNFTransDNF(w http.ResponseWriter, r *http.Request, cfg *config.Config) {
	algorithm := r.URL.Query().Get("algorithm")
	var method func(formula.Factory, []formula.Formula) (formula.Formula, sio.ServiceError)
	switch algorithm {
	case "factorization", "":
		method = func(fac formula.Factory, f []formula.Formula) (formula.Formula, sio.ServiceError) {
			hdl := normalform.HandlerWithTimeout(*handler.NewTimeoutWithDuration(cfg.SyncComputationTimout))
			result, ok := normalform.FactorizedDNFWithHandler(fac, fac.And(f...), hdl)
			return transformWithTimeout(result, ok)
		}
	case "canonical":
		method = func(fac formula.Factory, f []formula.Formula) (formula.Formula, sio.ServiceError) {
			hdl := iter.HandlerWithTimeout(*handler.NewTimeoutWithDuration(cfg.SyncComputationTimout))
			result, ok := enum.CanonicalDNFWithHandler(fac, fac.And(f...), hdl)
			return transformWithTimeout(result, ok)
		}
	case "bdd":
		method = func(fac formula.Factory, f []formula.Formula) (formula.Formula, sio.ServiceError) {
			hdl := bdd.HandlerWithTimeout(*handler.NewTimeoutWithDuration(cfg.SyncComputationTimout))
			result, ok := bdd.DNFWithHandler(fac, fac.And(f...), hdl)
			return transformWithTimeout(result, ok)
		}
	}
	transform(w, r, method)
}

// @Summary      Transform a formula to an and-inverter-graph
// @Description  If a list of formulas is given, the normal form is computed for the conjunction of these formulas.  The result always contains exactly one formula.
// @Tags         Normal Form
// @Param        request body	sio.FormulaInput true "Input formulas"
// @Success      200  {object}  sio.FormulaResult
// @Router       /normalform/transformation/aig [post]
func handleNFTransAIG(w http.ResponseWriter, r *http.Request) {
	transform(w, r, func(fac formula.Factory, form []formula.Formula) (formula.Formula, sio.ServiceError) {
		return normalform.AIG(fac, fac.And(form...)), nil
	})
}

// @Summary      Report whether a formula is an a certain normal form
// @Description  If a list of formulas is given, the predicate is computed for the conjunction of these formulas.
// @Tags         Normal Form
// @Param        nf path string true "normal form" Enums(nnf, cnf, dnf, aig, minterm, maxterm)
// @Param        request body	sio.FormulaInput true "Input formulas"
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
