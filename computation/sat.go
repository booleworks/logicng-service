package computation

import (
	"fmt"
	"net/http"

	"github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/handler"
	"github.com/booleworks/logicng-go/sat"
	"github.com/booleworks/logicng-service/config"
	"github.com/booleworks/logicng-service/sio"
)

// @Summary      Computes the satisfiability of a set of formulas with a SAT solver
// @Description  If a list of formulas is given, the satisfiability is computed for the conjunction of these formulas.
// @Tags         Solver
// @Param        core query string  false "Compte an unsat core if unsatisfiable" Enums(false, true) Default(false)
// @Param        request body	sio.FormulaInput true "Input formulas"
// @Success      200  {object}  sio.SatResult
// @Router       /solver/sat [post]
func HandleSat(cfg *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fac := formula.NewFactory()
		core := r.URL.Query().Get("core") == "true"
		satCfg := sat.DefaultConfig().Proofs(core)
		solver := sat.NewSolver(fac, satCfg)
		vars, ok := fillSatSolver(w, r, solver)
		if !ok {
			return
		}

		hdl := sat.HandlerWithTimeout(*handler.NewTimeoutWithDuration(cfg.SyncComputationTimout))
		var call *sat.CallParams
		if core {
			call = sat.WithCore().WithModel(vars).Handler(hdl)
		} else {
			call = sat.WithModel(vars).Handler(hdl)
		}
		result := solver.Call(call)
		if result.Aborted() {
			sio.WriteError(w, r, sio.ErrTimeout())
		} else {
			var mdl []string
			if result.Sat() {
				solverModel := result.Model()
				mdl = make([]string, solverModel.Size())
				for i, l := range solverModel.Literals {
					mdl[i] = l.Sprint(fac)
				}
			}
			var unsatCore []sio.Formula
			if core && !result.Sat() {
				props := result.UnsatCore().Propositions
				unsatCore = make([]sio.Formula, len(props))
				for i, p := range props {
					prop := p.(*formula.StandardProposition)
					unsatCore[i] = sio.Formula{Formula: p.Formula().Sprint(fac), Description: prop.Description}
				}
			}
			sio.WriteSatResult(w, r, result.Sat(), mdl, unsatCore)
		}
	})
}

// @Summary      Computes the backbone of a set of formulas
// @Description  If a list of formulas is given, the backbone is computed for the conjunction of these formulas.
// @Tags         Solver
// @Param        request body	sio.FormulaInput true "Input formulas"
// @Success      200  {object}  sio.BackboneResult
// @Router       /solver/backbone [post]
func HandleSatBackbone(cfg *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fac := formula.NewFactory()
		core := r.URL.Query().Get("core") == "true"
		satCfg := sat.DefaultConfig().Proofs(core)
		solver := sat.NewSolver(fac, satCfg)
		vars, ok := fillSatSolver(w, r, solver)
		if !ok {
			return
		}
		hdl := sat.HandlerWithTimeout(*handler.NewTimeoutWithDuration(cfg.SyncComputationTimout))
		bb, ok := solver.ComputeBackboneWithHandler(fac, vars, hdl)
		if !ok {
			sio.WriteError(w, r, sio.ErrTimeout())
		} else {
			sio.WriteBackboneResult(w, r, fac, bb)
		}
	})
}

func HandleSatPredicate(cfg *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}

func fillSatSolver(w http.ResponseWriter, r *http.Request, solver *sat.Solver) ([]formula.Variable, bool) {
	input, err := sio.Unmarshal[sio.FormulaInput](r)
	if err != nil {
		sio.WriteError(w, r, err)
		return nil, false
	}
	varSet := formula.NewMutableVarSet()
	for _, f := range input.Formulas {
		prop, ok := parseProp(w, r, solver.Factory(), f)
		if !ok {
			sio.WriteError(w, r, sio.ErrIllegalInput(fmt.Errorf("could not parse formula '%s'", f)))
			return nil, false
		}
		varSet.AddAll(formula.Variables(solver.Factory(), prop.Formula()))
		solver.AddProposition(prop)
	}
	return varSet.Content(), true
}
