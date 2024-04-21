package computation

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/randomizer"
	"github.com/booleworks/logicng-service/config"
	"github.com/booleworks/logicng-service/sio"
)

// @Summary      Generates a random formula
// @Tags         Randomizer
// @Param        fsort path string true "Formula sort to generate" Enums(const, var, lit, atom, not, impl, equiv, and, or, cc, amo, exo, pbc, formula) Default(formula)
// @Param        depth query int false "Formula depth"
// @Param        vars query int false "Number of variables"
// @Param        seed query int false "Seed for the randomizer"
// @Success      200  {object}  sio.FormulaResult
// @Router       /randomizer/{fsort} [get]
func HandleRandomizer(cfg *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		randCfg, ok := extractConfig(w, r)
		if !ok {
			return
		}
		depth := 3
		depthParam := r.URL.Query().Get("depth")
		if depthParam != "" {
			var err error
			depth, err = strconv.Atoi(depthParam)
			if err != nil {
				sio.WriteError(w, r, sio.ErrIllegalInput(fmt.Errorf("illegal depth value '%s'", depthParam)))
				return
			}
		}
		fac := formula.NewFactory()
		randomizer := randomizer.New(fac, randCfg)
		rand := r.PathValue("rand")
		switch rand {
		case "const":
			sio.WriteFormulaResult(w, r, randomizer.Constant().Sprint(fac))
		case "var":
			sio.WriteFormulaResult(w, r, randomizer.Variable().Sprint(fac))
		case "lit":
			sio.WriteFormulaResult(w, r, randomizer.Literal().Sprint(fac))
		case "atom":
			sio.WriteFormulaResult(w, r, randomizer.Atom().Sprint(fac))
		case "not":
			sio.WriteFormulaResult(w, r, randomizer.Not(depth).Sprint(fac))
		case "impl":
			sio.WriteFormulaResult(w, r, randomizer.Impl(depth).Sprint(fac))
		case "equiv":
			sio.WriteFormulaResult(w, r, randomizer.Equiv(depth).Sprint(fac))
		case "and":
			sio.WriteFormulaResult(w, r, randomizer.And(depth).Sprint(fac))
		case "or":
			sio.WriteFormulaResult(w, r, randomizer.Or(depth).Sprint(fac))
		case "cc":
			sio.WriteFormulaResult(w, r, randomizer.CC().Sprint(fac))
		case "amo":
			sio.WriteFormulaResult(w, r, randomizer.AMO().Sprint(fac))
		case "exo":
			sio.WriteFormulaResult(w, r, randomizer.EXO().Sprint(fac))
		case "pbc":
			sio.WriteFormulaResult(w, r, randomizer.PBC().Sprint(fac))
		case "formula":
			sio.WriteFormulaResult(w, r, randomizer.Formula(depth).Sprint(fac))
		default:
			sio.WriteError(w, r, sio.ErrUnknownPath(r.URL.Path))
		}
	})
}

func extractConfig(w http.ResponseWriter, r *http.Request) (*randomizer.Config, bool) {
	seed := r.URL.Query().Get("seed")
	numVars := r.URL.Query().Get("vars")
	randCfg := randomizer.DefaultConfig()
	randCfg.Seed = time.Now().UnixMilli()
	if seed != "" {
		s, err := strconv.Atoi(seed)
		if err != nil {
			sio.WriteError(w, r, sio.ErrIllegalInput(fmt.Errorf("illegal seed value '%s'", seed)))
			return nil, false
		}
		randCfg.Seed = int64(s)
	}
	if numVars != "" {
		n, err := strconv.Atoi(numVars)
		if err != nil {
			sio.WriteError(w, r, sio.ErrIllegalInput(fmt.Errorf("illegal num vars value '%s'", numVars)))
			return nil, false
		}
		randCfg.NumVars = n
	}
	return randCfg, true
}
