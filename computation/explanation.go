package computation

import (
	"fmt"
	"net/http"

	"github.com/booleworks/logicng-go/explanation"
	"github.com/booleworks/logicng-go/explanation/mus"
	"github.com/booleworks/logicng-go/explanation/smus"
	"github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/handler"
	"github.com/booleworks/logicng-go/sat"
	"github.com/booleworks/logicng-service/config"
	"github.com/booleworks/logicng-service/sio"
)

// @Summary      Compute a minimal unsatisfiable set (MUS) of an unsatisfiable formula
// @Description  If a list of formulas is given, the result refers to the conjunction of these formulas.
// @Tags         Explanation
// @Param        algorithm query string  false "MUS Algorithm" Enums(deletion, insertion) Default(deletion)
// @Param        request body	sio.FormulaInput true "Formula input"
// @Success      200  {object}  sio.FormulaResult
// @Router       /explanation/mus [post]
func HandleMUS(cfg *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		algorithm := r.URL.Query().Get("algorithm")
		fac := formula.NewFactory()
		ps, ok := parsePropInput(w, r, fac)
		if !ok {
			return
		}
		props := make([]formula.Proposition, len(ps))
		for i, p := range ps {
			props[i] = p
		}

		hdl := sat.HandlerWithTimeout(*handler.NewTimeoutWithDuration(cfg.SyncComputationTimout))
		var core *explanation.UnsatCore
		var err error
		switch algorithm {
		case "deletion", "":
			core, ok, err = mus.ComputeDeletionBasedWithHandler(fac, &props, hdl)
		case "insertion":
			core, ok, err = mus.ComputeInsertionBasedWithHandler(fac, &props, hdl)
		default:
			sio.WriteError(w, r, sio.ErrIllegalInput(fmt.Errorf("unknown MUS algorithm '%s'", algorithm)))
		}
		if err != nil {
			sio.WriteError(w, r, sio.ErrIllegalInput(err))
			return
		} else if !ok {
			sio.WriteError(w, r, sio.ErrTimeout())
			return
		}
		result := make([]sio.Formula, len(core.Propositions))
		for i, p := range core.Propositions {
			prop := p.(*formula.StandardProposition)
			result[i] = sio.Formula{Formula: p.Formula().Sprint(fac), Description: prop.Description}
		}
		sio.WriteFormulaResult(w, r, result...)
	})
}

// @Summary      Compute a shortest minimal unsatisfiable set (SMUS) of an unsatisfiable formula
// @Description  If a list of formulas is given, the result refers to the conjunction of these formulas.
// @Tags         Explanation
// @Param        request body	sio.FormulaInput true "Formula input"
// @Success      200  {object}  sio.FormulaResult
// @Router       /explanation/smus [post]
func HandleSMUS(cfg *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fac := formula.NewFactory()
		ps, ok := parsePropInput(w, r, fac)
		if !ok {
			return
		}
		props := make([]formula.Proposition, len(ps))
		for i, p := range ps {
			props[i] = p
		}

		hdl := sat.OptimizationHandlerWithTimeout(*handler.NewTimeoutWithDuration(cfg.SyncComputationTimout))
		res, ok := smus.ComputeWithHandler(fac, props, hdl)
		if len(res) == 0 {
			sio.WriteError(w, r, sio.ErrIllegalInput(fmt.Errorf("bad input: formula set is satisfiable")))
			return
		} else if !ok {
			sio.WriteError(w, r, sio.ErrTimeout())
			return
		}
		result := make([]sio.Formula, len(res))
		for i, p := range res {
			prop := p.(*formula.StandardProposition)
			result[i] = sio.Formula{Formula: p.Formula().Sprint(fac), Description: prop.Description}
		}
		sio.WriteFormulaResult(w, r, result...)
	})
}
