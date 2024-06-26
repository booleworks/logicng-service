package computation

import (
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/booleworks/logicng-go/bdd"
	"github.com/booleworks/logicng-go/dnnf"
	"github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/handler"
	"github.com/booleworks/logicng-go/model"
	"github.com/booleworks/logicng-go/model/count"
	"github.com/booleworks/logicng-go/model/enum"
	"github.com/booleworks/logicng-go/model/iter"
	"github.com/booleworks/logicng-service/config"
	"github.com/booleworks/logicng-service/sio"
)

// @Summary      Count the satisfying models of a formula
// @Description  If a list of formulas is given, the result refers to the conjunction of these formulas.
// @Tags         Model
// @Param        algorithm query string  false "Counting Algorithm" Enums(bdd, dnnf, sat) Default(dnnf)
// @Param        request body	sio.FormulaInput true "Formula input"
// @Success      200  {object}  sio.StringResult
// @Router       /model/counting [post]
func HandleModelCounting(cfg *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fac := formula.NewFactory()
		formulas, ok := parseFormulaInput(w, r, fac)
		if !ok {
			return
		}
		vars := formula.Variables(fac, formulas...).Content()
		var count *big.Int
		switch algorithm := r.URL.Query().Get("algorithm"); algorithm {
		case "dnnf", "":
			count, ok = countDNNF(w, r, fac, formulas, vars, cfg.SyncComputationTimout)
		case "bdd":
			count, ok = countBDD(w, r, fac, formulas, vars, cfg.SyncComputationTimout)
		case "sat":
			count, ok = countSat(w, r, fac, formulas, vars, cfg.SyncComputationTimout)
		default:
			sio.WriteError(w, r, sio.ErrIllegalInput(fmt.Errorf("unknown model counting algorithm '%s'", algorithm)))
		}
		if ok {
			sio.WriteStringResult(w, r, count.String())
		}
	})
}

// @Summary      Count the models of a formula projected to a set of variables
// @Description  If a list of formulas is given, the result refers to the conjunction of these formulas.
// @Tags         Model
// @Param        algorithm query string  false "Counting Algorithm" Enums(sat) Default(sat)
// @Param        request body	sio.FormulaVarsInput true "Formulas and variables input"
// @Success      200  {object}  sio.StringResult
// @Router       /model/counting/projection [post]
func HandleProjectedModelCounting(cfg *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fac := formula.NewFactory()
		input, err := sio.Unmarshal[sio.FormulaVarsInput](r)
		if err != nil {
			sio.WriteError(w, r, err)
			return
		}
		formulas, ok := parseFormulas(w, r, fac, input.Formulas)
		if !ok {
			return
		}
		vars := make([]formula.Variable, len(input.Variables))
		for i, v := range input.Variables {
			vars[i] = fac.Var(v)
		}

		var count *big.Int
		switch algorithm := r.URL.Query().Get("algorithm"); algorithm {
		case "sat", "":
			count, ok = countSat(w, r, fac, formulas, vars, cfg.SyncComputationTimout)
		// case "bdd": // TODO not yet working
		// 	count, ok = countBDD(w, r, fac, formulas, vars, cfg.SyncComputationTimout)
		default:
			sio.WriteError(w, r, sio.ErrIllegalInput(fmt.Errorf("unknown projected model counting algorithm '%s'", algorithm)))
		}
		if ok {
			sio.WriteStringResult(w, r, count.String())
		}
	})
}

// @Summary      Enumerate the satisfying models of a formula
// @Description  If a list of formulas is given, the result refers to the conjunction of these formulas.
// @Tags         Model
// @Param        algorithm query string  false "Enumeration Algorithm" Enums(bdd, sat) Default(bdd)
// @Param        request body	sio.FormulaInput true "Formula input"
// @Success      200  {object}  sio.FormulaResult
// @Router       /model/enumeration [post]
func HandleModelEnumeration(cfg *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fac := formula.NewFactory()
		fs, ok := parseFormulaInput(w, r, fac)
		if !ok {
			return
		}
		vars := formula.Variables(fac, fs...).Content()
		var enumeration []*model.Model
		switch algorithm := r.URL.Query().Get("algorithm"); algorithm {
		case "bdd", "":
			enumeration, ok = enumerateBDD(w, r, fac, fs, vars, cfg.SyncComputationTimout)
		case "sat":
			enumeration, ok = enumerateSat(w, r, fac, fs, vars, cfg.SyncComputationTimout)
		default:
			sio.WriteError(w, r, sio.ErrIllegalInput(fmt.Errorf("unknown model counting algorithm '%s'", algorithm)))
		}
		if ok {
			formulas := make([]sio.Formula, len(enumeration))
			for i, m := range enumeration {
				formulas[i] = sio.Formula{Formula: m.Formula(fac).Sprint(fac)}
			}
			sio.WriteFormulaResult(w, r, formulas...)
		}
	})
}

// @Summary      Enumerate the satisfying models of a formula projected to a set of variables
// @Description  If a list of formulas is given, the result refers to the conjunction of these formulas.
// @Tags         Model
// @Param        algorithm query string  false "Enumeration Algorithm" Enums(bdd, sat) Default(bdd)
// @Param        request body	sio.FormulaVarsInput true "Formulas and variables input"
// @Success      200  {object}  sio.FormulaResult
// @Router       /model/enumeration/projection [post]
func HandleProjectedModelEnumeration(cfg *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fac := formula.NewFactory()
		input, err := sio.Unmarshal[sio.FormulaVarsInput](r)
		if err != nil {
			sio.WriteError(w, r, err)
			return
		}
		fs, ok := parseFormulas(w, r, fac, input.Formulas)
		if !ok {
			return
		}
		vars := make([]formula.Variable, len(input.Variables))
		for i, v := range input.Variables {
			vars[i] = fac.Var(v)
		}

		var enumeration []*model.Model
		switch algorithm := r.URL.Query().Get("algorithm"); algorithm {
		case "bdd", "":
			enumeration, ok = enumerateBDD(w, r, fac, fs, vars, cfg.SyncComputationTimout)
		case "sat":
			enumeration, ok = enumerateSat(w, r, fac, fs, vars, cfg.SyncComputationTimout)
		default:
			sio.WriteError(w, r, sio.ErrIllegalInput(fmt.Errorf("unknown model counting algorithm '%s'", algorithm)))
		}
		if ok {
			formulas := make([]sio.Formula, len(enumeration))
			for i, m := range enumeration {
				formulas[i] = sio.Formula{Formula: m.Formula(fac).Sprint(fac)}
			}
			sio.WriteFormulaResult(w, r, formulas...)
		}
	})
}

func countDNNF(
	w http.ResponseWriter,
	r *http.Request,
	fac formula.Factory,
	formulas []formula.Formula,
	vars []formula.Variable,
	timeout time.Duration,
) (*big.Int, bool) {
	hdl := dnnf.HandlerWithTimeout(*handler.NewTimeoutWithDuration(timeout))
	cnt, err, ok := count.CountWithHandler(fac, vars, hdl, formulas...)
	if err != nil {
		sio.WriteError(w, r, sio.ErrIllegalInput(err))
		return nil, false
	}
	if !ok {
		sio.WriteError(w, r, sio.ErrTimeout())
		return nil, false
	}
	return cnt, true
}

func countBDD(
	w http.ResponseWriter,
	r *http.Request,
	fac formula.Factory,
	formulas []formula.Formula,
	vars []formula.Variable,
	timeout time.Duration,
) (*big.Int, bool) {
	hdl := bdd.HandlerWithTimeout(*handler.NewTimeoutWithDuration(timeout))
	f := fac.And(formulas...)
	order := bdd.ForceOrder(fac, f)
	bdd, ok := bdd.CompileWithVarOrderAndHandler(fac, f, order, hdl)
	allVars := formula.NewMutableVarSetCopy(formula.Variables(fac, formulas...))
	allVars.RemoveAllElements(&vars)
	if !allVars.Empty() {
		bdd = bdd.Exists(allVars.Content()...)
	}
	if !ok {
		sio.WriteError(w, r, sio.ErrTimeout())
		return nil, false
	}
	return bdd.ModelCount(), true
}

func countSat(
	w http.ResponseWriter,
	r *http.Request,
	fac formula.Factory,
	formulas []formula.Formula,
	vars []formula.Variable,
	timeout time.Duration,
) (*big.Int, bool) {
	f := fac.And(formulas...)
	cfg := iter.DefaultConfig()
	cfg.Handler = iter.HandlerWithTimeout(*handler.NewTimeoutWithDuration(timeout))
	cnt, ok := count.OnFormulaWithConfig(fac, f, vars, cfg)
	if !ok {
		sio.WriteError(w, r, sio.ErrTimeout())
		return nil, false
	}
	return cnt, true
}

func enumerateBDD(
	w http.ResponseWriter,
	r *http.Request,
	fac formula.Factory,
	formulas []formula.Formula,
	vars []formula.Variable,
	timeout time.Duration,
) ([]*model.Model, bool) {
	hdl := bdd.HandlerWithTimeout(*handler.NewTimeoutWithDuration(timeout))
	f := fac.And(formulas...)
	order := bdd.ForceOrder(fac, f)
	bdd, ok := bdd.CompileWithVarOrderAndHandler(fac, f, order, hdl)
	allVars := formula.NewMutableVarSetCopy(formula.Variables(fac, formulas...))
	allVars.RemoveAllElements(&vars)
	if !allVars.Empty() {
		bdd = bdd.Exists(allVars.Content()...)
	}
	if !ok {
		sio.WriteError(w, r, sio.ErrTimeout())
		return nil, false
	}
	return bdd.ModelEnumeration(vars...), true
}

func enumerateSat(
	w http.ResponseWriter,
	r *http.Request,
	fac formula.Factory,
	formulas []formula.Formula,
	vars []formula.Variable,
	timeout time.Duration,
) ([]*model.Model, bool) {
	f := fac.And(formulas...)
	cfg := iter.DefaultConfig()
	cfg.Handler = iter.HandlerWithTimeout(*handler.NewTimeoutWithDuration(timeout))
	enumeration, ok := enum.OnFormulaWithConfig(fac, f, vars, cfg)
	if !ok {
		sio.WriteError(w, r, sio.ErrTimeout())
		return nil, false
	}
	return enumeration, true
}
