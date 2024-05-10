package computation

import (
	"fmt"
	"net/http"

	"github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/handler"
	"github.com/booleworks/logicng-go/normalform"
	"github.com/booleworks/logicng-go/sat"
	"github.com/booleworks/logicng-go/simplification"
	"github.com/booleworks/logicng-service/config"
	"github.com/booleworks/logicng-service/sio"
)

func HandleSimplification(cfg *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		simp := r.PathValue("simp")
		switch simp {
		case "backbone":
			handleSimplBackbone(w, r)
		case "unitpropagation":
			handleSimplUnitProp(w, r)
		case "negation":
			handleSimplNegation(w, r)
		case "distribution":
			handleSimplDistribution(w, r)
		case "factorout":
			handleSimplFactorOut(w, r)
		case "subsumption":
			handleSimplSubsumption(w, r)
		case "qmc":
			handleSimplQMC(w, r, cfg)
		case "advanced":
			handleSimplAdvanced(w, r, cfg)
		default:
			sio.WriteError(w, r, sio.ErrUnknownPath(r.URL.Path))
		}
	})
}

// @Summary      Simplify a formula by computing and propagating its backbone
// @Description  If a list of formulas is given, the simplification is computed for the conjunction of these formulas.  The result always contains exactly one formula.
// @Tags         Simplification
// @Param        request body	sio.FormulaInput true "Input formulas"
// @Success      200  {object}  sio.FormulaResult
// @Router       /simplification/backbone [post]
func handleSimplBackbone(w http.ResponseWriter, r *http.Request) {
	transform(w, r, func(fac formula.Factory, fs []formula.Formula) (formula.Formula, sio.ServiceError) {
		return simplification.PropagateBackbone(fac, fac.And(fs...)), nil
	})
}

// @Summary      Simplify a formula by propagating its unit literals
// @Description  If a list of formulas is given, the simplification is computed for the conjunction of these formulas.  The result always contains exactly one formula.
// @Tags         Simplification
// @Param        request body	sio.FormulaInput true "Input formulas"
// @Success      200  {object}  sio.FormulaResult
// @Router       /simplification/unitpropagation [post]
func handleSimplUnitProp(w http.ResponseWriter, r *http.Request) {
	transform(w, r, func(fac formula.Factory, fs []formula.Formula) (formula.Formula, sio.ServiceError) {
		return simplification.PropagateUnits(fac, fac.And(fs...)), nil
	})
}

// @Summary      Simplify a formula by minimizing the number of negations
// @Description  If a list of formulas is given, the simplification is computed for the conjunction of these formulas.  The result always contains exactly one formula.
// @Tags         Simplification
// @Param        request body	sio.FormulaInput true "Input formulas"
// @Success      200  {object}  sio.FormulaResult
// @Router       /simplification/negation [post]
func handleSimplNegation(w http.ResponseWriter, r *http.Request) {
	transform(w, r, func(fac formula.Factory, fs []formula.Formula) (formula.Formula, sio.ServiceError) {
		return simplification.MinimizeNegations(fac, fac.And(fs...)), nil
	})
}

// @Summary      Simplify a formula by applying the distributive laws
// @Description  If a list of formulas is given, the simplification is computed for the conjunction of these formulas.  The result always contains exactly one formula.
// @Tags         Simplification
// @Param        request body	sio.FormulaInput true "Input formulas"
// @Success      200  {object}  sio.FormulaResult
// @Router       /simplification/distribution [post]
func handleSimplDistribution(w http.ResponseWriter, r *http.Request) {
	transform(w, r, func(fac formula.Factory, fs []formula.Formula) (formula.Formula, sio.ServiceError) {
		return simplification.Distribute(fac, fac.And(fs...)), nil
	})
}

// @Summary      Simplify a formula by factoring out common factors repetitively
// @Description  If a list of formulas is given, the simplification is computed for the conjunction of these formulas.  The result always contains exactly one formula.
// @Tags         Simplification
// @Param        request body	sio.FormulaInput true "Input formulas"
// @Success      200  {object}  sio.FormulaResult
// @Router       /simplification/factorout [post]
func handleSimplFactorOut(w http.ResponseWriter, r *http.Request) {
	transform(w, r, func(fac formula.Factory, fs []formula.Formula) (formula.Formula, sio.ServiceError) {
		return simplification.FactorOut(fac, fac.And(fs...)), nil
	})
}

// @Summary      Simplify a CNF or DNF by applying subsumptions
// @Description  If a list of formulas is given, the simplification is computed for the conjunction of these formulas.  The result always contains exactly one formula.
// @Tags         Simplification
// @Param        request body	sio.FormulaInput true "Input formulas"
// @Success      200  {object}  sio.FormulaResult
// @Router       /simplification/subsumption [post]
func handleSimplSubsumption(w http.ResponseWriter, r *http.Request) {
	transform(w, r, func(fac formula.Factory, fs []formula.Formula) (formula.Formula, sio.ServiceError) {
		form := fac.And(fs...)
		switch {
		case normalform.IsCNF(fac, form):
			s, _ := simplification.CNFSubsumption(fac, form)
			return s, nil
		case normalform.IsDNF(fac, form):
			s, _ := simplification.DNFSubsumption(fac, form)
			return s, nil
		default:
			return 0, sio.ErrIllegalInput(fmt.Errorf("input for subsumption must be in CNF or DNF"))
		}
	})
}

// @Summary      Simplify a formula with the Quine-McCluskey algorithm
// @Description  If a list of formulas is given, the simplification is computed for the conjunction of these formulas.  The result always contains exactly one formula.
// @Tags         Simplification
// @Param        request body	sio.FormulaInput true "Input formulas"
// @Success      200  {object}  sio.FormulaResult
// @Router       /simplification/qmc [post]
func handleSimplQMC(w http.ResponseWriter, r *http.Request, cfg *config.Config) {
	transform(w, r, func(fac formula.Factory, fs []formula.Formula) (formula.Formula, sio.ServiceError) {
		hdl := sat.OptimizationHandlerWithTimeout(*handler.NewTimeoutWithDuration(cfg.SyncComputationTimout))
		result, ok := simplification.QMCWithHandler(fac, fac.And(fs...), hdl)
		return transformWithTimeout(result, ok)
	})
}

// @Summary      Simplify a formula with the advanced simplifier
// @Description  If a list of formulas is given, the simplification is computed for the conjunction of these formulas.  The result always contains exactly one formula.
// @Tags         Simplification
// @Param        backbone query string false "Simplify with backbone" Enums(true, false) default(true)
// @Param        factorout query string false "Factor out common factors" Enums(true, false) default(true)
// @Param        negations query string false "Minimize negations" Enums(true, false) default(true)
// @Param        request body	sio.FormulaInput true "Input formulas"
// @Success      200  {object}  sio.FormulaResult
// @Router       /simplification/advanced [post]
func handleSimplAdvanced(w http.ResponseWriter, r *http.Request, cfg *config.Config) {
	transform(w, r, func(fac formula.Factory, fs []formula.Formula) (formula.Formula, sio.ServiceError) {
		simpCfg := simplification.DefaultConfig()
		if r.URL.Query().Get("backbone") == "false" {
			simpCfg.RestrictBackbone = false
		}
		if r.URL.Query().Get("factorout") == "false" {
			simpCfg.FactorOut = false
		}
		if r.URL.Query().Get("negations") == "false" {
			simpCfg.SimplifyNegations = false
		}
		hdl := sat.OptimizationHandlerWithTimeout(*handler.NewTimeoutWithDuration(cfg.SyncComputationTimout))
		result, ok := simplification.AdvancedWithHandler(fac, fac.And(fs...), hdl, simpCfg)
		return transformWithTimeout(result, ok)
	})
}
