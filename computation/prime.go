package computation

import (
	"fmt"
	"net/http"

	"github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/handler"
	"github.com/booleworks/logicng-go/primeimplicant"
	"github.com/booleworks/logicng-go/sat"
	"github.com/booleworks/logicng-service/config"
	"github.com/booleworks/logicng-service/sio"
)

// @Summary      Computes a minimal prime implicant of a formula
// @Description  If a list of formulas is given, the prime implicant is computed for the conjunction of these formulas.  The result always contains exactly one formula.
// @Tags         Prime Implicant
// @Param        request body	sio.FormulaInput true "Input formulas"
// @Success      200  {object}  sio.FormulaResult
// @Router       /prime/minimal-implicant [post]
func HandleMinimalImplicant(cfg *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fac := formula.NewFactory()
		formulas, ok := parseFormulaInput(w, r, fac)
		if !ok {
			return
		}
		implicant, err := primeimplicant.Minimum(fac, fac.And(formulas...))
		if err != nil {
			sio.WriteError(w, r, sio.ErrIllegalInput(err))
		} else {
			implicantFormula := fac.And(formula.LiteralsAsFormulas(implicant)...)
			sio.WriteFormulaResult(w, r, sio.Formula{Formula: implicantFormula.Sprint(fac)})
		}
	})
}

// @Summary      Computes a minimal prime implicant cover of a formula
// @Description  If a list of formulas is given, the prime implicant cover is computed for the conjunction of these formulas.
// @Tags         Prime Implicant
// @Param        algorithm query string  false "min or max models" Enums(min, max) Default(max)
// @Param        request body	sio.FormulaInput true "Input formulas"
// @Success      200  {object}  sio.FormulaResult
// @Router       /prime/minimal-cover [post]
func HandleMinimalImplicantCover(cfg *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fac := formula.NewFactory()
		formulas, ok := parseFormulaInput(w, r, fac)
		if !ok {
			return
		}
		form := fac.And(formulas...)

		hdl := sat.OptimizationHandlerWithTimeout(*handler.NewTimeoutWithDuration(cfg.SyncComputationTimout))
		var result *primeimplicant.PrimeResult
		switch algorithm := r.URL.Query().Get("algorithm"); algorithm {
		case "max", "":
			result, ok = primeimplicant.CoverMaxWithHandler(fac, form, primeimplicant.CoverImplicants, hdl)
		case "min":
			result, ok = primeimplicant.CoverMinWithHandler(fac, form, primeimplicant.CoverImplicants, hdl)
		default:
			sio.WriteError(w, r, sio.ErrIllegalInput(fmt.Errorf("unknown prime implicant cover algorithm '%s'", algorithm)))
			return
		}
		if !ok {
			sio.WriteError(w, r, sio.ErrTimeout())
			return
		}
		implicants := make([]sio.Formula, len(result.Implicants))
		for i, impl := range result.Implicants {
			f := fac.And(formula.LiteralsAsFormulas(impl)...).Sprint(fac)
			implicants[i] = sio.Formula{Formula: f}
		}
		sio.WriteFormulaResult(w, r, implicants...)
	})
}
