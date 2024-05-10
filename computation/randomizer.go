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

// @Summary      Generate a random formula
// @Tags         Randomizer
// @Param        fsort path string true "Formula sort to generate" Enums(const, var, lit, atom, not, impl, equiv, and, or, cc, amo, exo, pbc, formula) Default(formula)
// @Param        depth query int false "Formula depth"
// @Param        vars query int false "Number of variables"
// @Param        seed query int false "Seed for the randomizer"
// @Param        formulas query int false "Number of formulas to generate"
// @Success      200  {object}  sio.FormulaResult
// @Router       /randomizer/{fsort} [get]
func HandleRandomizer(cfg *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		randCfg, ok := extractRandConfig(w, r)
		if !ok {
			return
		}
		depth, ok := extractIntParam(w, r, "depth", 3)
		if !ok {
			return
		}
		numForms, ok := extractIntParam(w, r, "formulas", 1)
		if !ok {
			return
		}

		fac := formula.NewFactory()
		randomizer := randomizer.New(fac, randCfg)
		rand := r.PathValue("rand")

		var randGen func() formula.Formula
		switch rand {
		case "const":
			randGen = randomizer.Constant
		case "var":
			randGen = func() formula.Formula { return randomizer.Variable().AsFormula() }
		case "lit":
			randGen = func() formula.Formula { return randomizer.Literal().AsFormula() }
		case "atom":
			randGen = randomizer.Atom
		case "not":
			randGen = func() formula.Formula { return randomizer.Not(depth) }
		case "impl":
			randGen = func() formula.Formula { return randomizer.Impl(depth) }
		case "equiv":
			randGen = func() formula.Formula { return randomizer.Equiv(depth) }
		case "and":
			randGen = func() formula.Formula { return randomizer.And(depth) }
		case "or":
			randGen = func() formula.Formula { return randomizer.Or(depth) }
		case "cc":
			randGen = randomizer.CC
		case "amo":
			randGen = randomizer.AMO
		case "exo":
			randGen = randomizer.EXO
		case "pbc":
			randGen = randomizer.PBC
		case "formula":
			randGen = func() formula.Formula { return randomizer.Formula(depth) }
		default:
			sio.WriteError(w, r, sio.ErrUnknownPath(r.URL.Path))
			return
		}
		res := make([]sio.Formula, numForms)
		for i := 0; i < numForms; i++ {
			res[i] = sio.Formula{Formula: randGen().Sprint(fac)}
		}
		sio.WriteFormulaResult(w, r, res...)
	})
}

func extractRandConfig(w http.ResponseWriter, r *http.Request) (*randomizer.Config, bool) {
	randCfg := randomizer.DefaultConfig()

	seed, ok := extractIntParam(w, r, "seed", int(time.Now().UnixMilli()))
	randCfg.Seed = int64(seed)
	if !ok {
		return nil, false
	}
	randCfg.Seed = int64(seed)

	numVars, ok := extractIntParam(w, r, "vars", 25)
	if !ok {
		return nil, false
	}
	randCfg.NumVars = numVars
	return randCfg, true
}

func extractIntParam(w http.ResponseWriter, r *http.Request, param string, def int) (int, bool) {
	value := def
	valueParam := r.URL.Query().Get(param)
	if valueParam != "" {
		var err error
		value, err = strconv.Atoi(valueParam)
		if err != nil {
			sio.WriteError(w, r, sio.ErrIllegalInput(fmt.Errorf("illegal %s value '%s'", param, valueParam)))
			return 0, false
		}
	}
	return value, true
}
