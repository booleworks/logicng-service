package computation

import (
	"fmt"
	"net/http"

	"github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/handler"
	"github.com/booleworks/logicng-go/maxsat"
	"github.com/booleworks/logicng-service/config"
	"github.com/booleworks/logicng-service/sio"
)

// @Summary      Solves a given set of hard and soft formulas with a MAX-SAT solver
// @Tags         Solver
// @Param        algorithm query string  false "MAX-SAT Algorithm" Enums(oll, msu3, wmsu3, linear-su, linear-us, wbo, inc-wbo)
// @Param        request body	sio.MaxSatInput true "MAX-SAT input"
// @Success      200  {object}  sio.MaxSatResult
// @Router       /solver/maxsat [post]
func HandleMaxSat(cfg *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fac := formula.NewFactory()
		solver, ok := extractMaxSatSolver(w, r, fac)
		if !ok {
			return
		}
		if ok := fillSolver(w, r, fac, solver); !ok {
			return
		}
		hdl := maxsat.HandlerWithTimeout(*handler.NewTimeoutWithDuration(cfg.SyncComputationTimout))
		result, ok := solver.SolveWithHandler(hdl)
		if !ok {
			sio.WriteError(w, r, sio.ErrTimeout())
		} else {
			var mdl []string
			if result.Satisfiable {
				solverModel, _ := solver.Model()
				mdl = make([]string, solverModel.Size())
				for i, l := range solverModel.Literals {
					mdl[i] = l.Sprint(fac)
				}
			}
			sio.WriteMaxSatResult(w, r, result.Satisfiable, int64(result.Optimum), mdl)
		}
	})
}

func extractMaxSatSolver(w http.ResponseWriter, r *http.Request, fac formula.Factory) (*maxsat.Solver, bool) {
	var solver *maxsat.Solver
	algorithm := r.URL.Query().Get("algorithm")
	switch algorithm {
	case "", "oll":
		solver = maxsat.OLL(fac)
	case "msu3":
		solver = maxsat.MSU3(fac)
	case "wmsu3":
		solver = maxsat.WMSU3(fac)
	case "linear-su":
		solver = maxsat.LinearSU(fac)
	case "linear-us":
		solver = maxsat.LinearUS(fac)
	case "wbo":
		solver = maxsat.WBO(fac)
	case "inc-wbo":
		solver = maxsat.IncWBO(fac)
	default:
		sio.WriteError(w, r, sio.ErrIllegalInput(fmt.Errorf("unknown maxsat algorithm '%s'", algorithm)))
		return nil, false
	}
	return solver, true
}

func fillSolver(w http.ResponseWriter, r *http.Request, fac formula.Factory, solver *maxsat.Solver) bool {
	input, err := sio.Unmarshal[sio.MaxSatInput](r)
	if err != nil {
		sio.WriteError(w, r, err)
		return false
	}
	for _, f := range input.HardFormulas {
		parsed, ok := parse(w, r, fac, f)
		if !ok {
			sio.WriteError(w, r, sio.ErrIllegalInput(fmt.Errorf("could not parse hard clause '%s'", f)))
			return false
		}
		solver.AddHardFormula(parsed)
	}
	suppWeighted := solver.SupportsWeighted()
	realWeighted := false
	for _, f := range input.SoftFormulas {
		parsed, ok := parse(w, r, fac, f.Formula)
		if !ok {
			sio.WriteError(w, r, sio.ErrIllegalInput(fmt.Errorf("could not parse soft clause '%s'", f.Formula)))
			return false
		}
		if f.Weight > 1 {
			realWeighted = true
		}
		if f.Weight > 1 && !suppWeighted {
			sio.WriteError(w, r, sio.ErrIllegalInput(fmt.Errorf("algorithm does not support weighted instances")))
			return false
		}
		solver.AddSoftFormula(parsed, int(f.Weight))
	}
	if !solver.SupportsUnweighted() && !realWeighted {
		sio.WriteError(w, r, sio.ErrIllegalInput(fmt.Errorf("algorithm does not support unweighted instances")))
		return false
	}
	return true
}
