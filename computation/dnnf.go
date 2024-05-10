package computation

import (
	"net/http"

	"github.com/booleworks/logicng-go/dnnf"
	"github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/handler"
	"github.com/booleworks/logicng-service/config"
	"github.com/booleworks/logicng-service/sio"
)

// @Summary      Compiles formulas to DNNF
// @Description  If a list of formulas is given, the DNNF of the conjunction of these formulas is computed.  The result always contains exactly one formula.
// @Tags         DNNF
// @Param        request body	sio.FormulaInput true "Input formulas"
// @Success      200  {object}  sio.FormulaResult
// @Router       /dnnf/compilation [post]
func HandleDNNFCompilation(cfg *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fac := formula.NewFactory()
		fs, ok := parseFormulaInput(w, r, fac)
		if !ok {
			return
		}
		hdl := dnnf.HandlerWithTimeout(*handler.NewTimeoutWithDuration(cfg.SyncComputationTimout))
		compiled, ok := dnnf.CompileWithHandler(fac, fac.And(fs...), hdl)
		if !ok {
			sio.WriteError(w, r, sio.ErrTimeout())
		}
		sio.WriteFormulaResult(w, r, sio.Formula{Formula: compiled.Formula.Sprint(fac)})
	})
}
