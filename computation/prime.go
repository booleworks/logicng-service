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
// @Tags         Prime Implicant
// @Param        request body	sio.FormulaInput true "Formula input"
// @Success      200  {object}  sio.FormulaResult
// @Router       /prime/minimal-implicant [post]
func HandleMinimalImplicant(cfg *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fac := formula.NewFactory()
		form, ok := parseFormulaInput(w, r, fac)
		if !ok {
			return
		}
		implicant, err := primeimplicant.Minimum(fac, form)
		if err != nil {
			sio.WriteError(w, r, sio.ErrIllegalInput(err))
		} else {
			implicantFormula := fac.And(formula.LiteralsAsFormulas(implicant)...)
			sio.WriteFormulaResult(w, r, implicantFormula.Sprint(fac))
		}
	})
}

// @Summary      Computes a minimal prime implicant cover of a formula
// @Tags         Prime Implicant
// @Param        algorithm query string  false "min or max models" Enums(min, max) Default(max)
// @Param        request body	sio.FormulaInput true "Formula input"
// @Success      200  {object}  sio.FormulaSetResult
// @Router       /prime/minimal-cover [post]
func HandleMinimalImplicantCover(cfg *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fac := formula.NewFactory()
		form, ok := parseFormulaInput(w, r, fac)
		if !ok {
			return
		}

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
		implicants := make([]string, len(result.Implicants))
		for i, impl := range result.Implicants {
			implicants[i] = fac.And(formula.LiteralsAsFormulas(impl)...).Sprint(fac)
		}
		sio.WriteFormulaSetResult(w, r, implicants)
	})
}
