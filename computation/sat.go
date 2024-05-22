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

// @Summary      Compute the satisfiability of a set of formulas with a SAT solver
// @Description  If a list of formulas is given, the satisfiability is computed for the conjunction of these formulas.
// @Tags         Solver
// @Param        core query string  false "Compute an unsat core if unsatisfiable" Enums(false, true) Default(false)
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

// @Summary      Compute the backbone of a set of formulas
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
		switch predicate := r.PathValue("pred"); predicate {
		case "tautology":
			HandleTautology(w, r, cfg)
		case "contradiction":
			HandleContradiction(w, r, cfg)
		case "implication":
			HandleImplication(w, r, cfg)
		case "equivalence":
			HandleEquivalence(w, r, cfg)
		default:
			sio.WriteError(w, r, sio.ErrUnknownPath(r.URL.Path))
		}
	})
}

// @Summary      Report whether a formula is a tautology
// @Description  If a list of formulas is given it is reported of the conjunction of these formulas are a tautology.
// @Tags         Solver
// @Param        request body	sio.FormulaInput true "Input formulas"
// @Success      200  {object}  sio.BoolResult
// @Router       /solver/predicate/tautology [post]
func HandleTautology(w http.ResponseWriter, r *http.Request, cfg *config.Config) {
	handleTautCont(w, r, cfg, true)
}

// @Summary      Report whether a formula is a contradiction
// @Description  If a list of formulas is given it is reported of the conjunction of these formulas are a contradiction.
// @Tags         Solver
// @Param        request body	sio.FormulaInput true "Input formulas"
// @Success      200  {object}  sio.BoolResult
// @Router       /solver/predicate/contradiction [post]
func HandleContradiction(w http.ResponseWriter, r *http.Request, cfg *config.Config) {
	handleTautCont(w, r, cfg, false)
}

// @Summary      Report whether the first formula implies the second formula
// @Description  Must be called with exactly two formulas.
// @Tags         Solver
// @Param        request body	sio.FormulaInput true "Input formulas"
// @Success      200  {object}  sio.BoolResult
// @Router       /solver/predicate/implication [post]
func HandleImplication(w http.ResponseWriter, r *http.Request, cfg *config.Config) {
	handleImplEquiv(w, r, cfg, true)
}

// @Summary      Report whether the first formula and the second formula are equivalent
// @Description  Must be called with exactly two formulas.
// @Tags         Solver
// @Param        request body	sio.FormulaInput true "Input formulas"
// @Success      200  {object}  sio.BoolResult
// @Router       /solver/predicate/equivalence [post]
func HandleEquivalence(w http.ResponseWriter, r *http.Request, cfg *config.Config) {
	handleImplEquiv(w, r, cfg, false)
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

func handleTautCont(w http.ResponseWriter, r *http.Request, cfg *config.Config, taut bool) {
	fac := formula.NewFactory()
	fs, ok := parseFormulaInput(w, r, fac)
	if !ok {
		return
	}
	solver := sat.NewSolver(fac)
	if taut {
		solver.Add(fac.Not(fac.And(fs...)))
	} else {
		solver.Add(fac.And(fs...))
	}
	hdl := sat.HandlerWithTimeout(*handler.NewTimeoutWithDuration(cfg.SyncComputationTimout))
	result := solver.Call(sat.Params().Handler(hdl))
	if result.Aborted() {
		sio.WriteError(w, r, sio.ErrTimeout())
	} else {
		sio.WriteBoolResult(w, r, !result.Sat())
	}
}

func handleImplEquiv(w http.ResponseWriter, r *http.Request, cfg *config.Config, impl bool) {
	fac := formula.NewFactory()
	fs, ok := parseFormulaInput(w, r, fac)
	if !ok {
		return
	}
	if len(fs) != 2 {
		sio.WriteError(w, r, sio.ErrIllegalInput(fmt.Errorf("method must be called with exactly two formulas")))
		return
	}
	solver := sat.NewSolver(fac)
	if impl {
		solver.Add(fac.Not(fac.Implication(fs[0], fs[1])))
	} else {
		solver.Add(fac.Not(fac.Equivalence(fs[0], fs[1])))
	}
	hdl := sat.HandlerWithTimeout(*handler.NewTimeoutWithDuration(cfg.SyncComputationTimout))
	result := solver.Call(sat.Params().Handler(hdl))
	if result.Aborted() {
		sio.WriteError(w, r, sio.ErrTimeout())
	} else {
		sio.WriteBoolResult(w, r, !result.Sat())
	}
}
