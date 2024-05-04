package computation

import (
	"net/http"

	"github.com/booleworks/logicng-go/dnnf"
	"github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/handler"
	"github.com/booleworks/logicng-service/config"
	"github.com/booleworks/logicng-service/sio"
)

// @Summary      Compiles a formula to DNNF
// @Tags         DNNF
// @Param        request body	sio.FormulaInput true "Formula input"
// @Success      200  {object}  sio.FormulaResult
// @Router       /dnnf/compilation [post]
func HandleDNNFCompilation(cfg *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fac := formula.NewFactory()
		formula, ok := parseFormulaInput(w, r, fac)
		if !ok {
			return
		}
		hdl := dnnf.HandlerWithTimeout(*handler.NewTimeoutWithDuration(cfg.SyncComputationTimout))
		compiled, ok := dnnf.CompileWithHandler(fac, formula, hdl)
		if !ok {
			sio.WriteError(w, r, sio.ErrTimeout())
		}
		sio.WriteFormulaResult(w, r, compiled.Formula.Sprint(fac))
	})
}
