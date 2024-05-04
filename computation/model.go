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
	"github.com/booleworks/logicng-go/model/count"
	"github.com/booleworks/logicng-go/model/iter"
	"github.com/booleworks/logicng-service/config"
	"github.com/booleworks/logicng-service/sio"
)

// @Summary      Counts the models of a formula
// @Tags         Model
// @Param        algorithm query string  false "Counting Algorithm" Enums(bdd, dnnf, sat) Default(dnnf)
// @Param        request body	sio.FormulaSetInput true "Formula set input"
// @Success      200  {object}  sio.StringResult
// @Router       /model/counting [post]
func HandleModelCounting(cfg *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		algorithm := r.URL.Query().Get("algorithm")
		fac := formula.NewFactory()
		formulas, ok := parseFormulaSetInput(w, r, fac)
		if !ok {
			return
		}
		var count *big.Int
		switch algorithm {
		case "dnnf", "":
			count, ok = countDNNF(w, r, fac, formulas, cfg.SyncComputationTimout)
		case "bdd":
			count, ok = countBDD(w, r, fac, formulas, cfg.SyncComputationTimout)
		case "sat":
			count, ok = countSat(w, r, fac, formulas, cfg.SyncComputationTimout)
		default:
			sio.WriteError(w, r, sio.ErrIllegalInput(fmt.Errorf("unknown maxsat algorithm '%s'", algorithm)))
		}
		if ok {
			sio.WriteStringResult(w, r, count.String())
		}
	})
}

func countDNNF(
	w http.ResponseWriter,
	r *http.Request,
	fac formula.Factory,
	formulas []formula.Formula,
	timeout time.Duration,
) (*big.Int, bool) {
	vars := formula.Variables(fac, formulas...).Content()
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
	timeout time.Duration,
) (*big.Int, bool) {
	hdl := bdd.HandlerWithTimeout(*handler.NewTimeoutWithDuration(timeout))
	f := fac.And(formulas...)
	order := bdd.ForceOrder(fac, f)
	bdd, ok := bdd.CompileWithVarOrderAndHandler(fac, f, order, hdl)
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
	timeout time.Duration,
) (*big.Int, bool) {
	f := fac.And(formulas...)
	cfg := iter.DefaultConfig()
	cfg.Handler = iter.HandlerWithTimeout(*handler.NewTimeoutWithDuration(timeout))
	vars := formula.Variables(fac, formulas...).Content()
	cnt, ok := count.OnFormulaWithConfig(fac, f, vars, cfg)
	if !ok {
		sio.WriteError(w, r, sio.ErrTimeout())
		return nil, false
	}
	return cnt, true
}
