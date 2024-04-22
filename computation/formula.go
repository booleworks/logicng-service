package computation

import (
	"net/http"

	"github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/function"
	"github.com/booleworks/logicng-service/config"
	"github.com/booleworks/logicng-service/sio"
)

func HandleFormula(cfg *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch function := r.PathValue("func"); function {
		case "depth":
			handleFormulaDepth(w, r)
		case "atoms":
			handleFormulaAtoms(w, r)
		case "nodes":
			handleFormulaNodes(w, r)
		case "variables":
			handleFormulaVariables(w, r)
		case "literals":
			handleFormulaLiterals(w, r)
		case "sub-formulas":
			handleFormulaSubFormulas(w, r)
		case "var-profile":
			handleFormulaVarProfile(w, r)
		case "lit-profile":
			handleFormulaLitProfile(w, r)
		// case "dag-graph":
		// case "ast-graph":
		default:
			sio.WriteError(w, r, sio.ErrUnknownPath(r.URL.Path))
		}
	})
}

// @Summary      Computes the depth of a formula's AST
// @Tags         Formula
// @Param        request body	sio.FormulaInput true "Input Formula"
// @Success      200  {object}  sio.IntResult
// @Router       /formula/depth [post]
func handleFormulaDepth(w http.ResponseWriter, r *http.Request) {
	fac := formula.NewFactory()
	formula, err := parseFormulaInput(w, r, fac)
	if !err {
		return
	}
	sio.WriteIntResult(w, r, int64(function.FormulaDepth(fac, formula)))
}

// @Summary      Computes the number of atoms of a formula
// @Tags         Formula
// @Param        request body	sio.FormulaInput true "Input Formula"
// @Success      200  {object}  sio.IntResult
// @Router       /formula/atoms [post]
func handleFormulaAtoms(w http.ResponseWriter, r *http.Request) {
	fac := formula.NewFactory()
	formula, err := parseFormulaInput(w, r, fac)
	if !err {
		return
	}
	sio.WriteIntResult(w, r, int64(function.NumberOfAtoms(fac, formula)))
}

// @Summary      Computes the number of nodes of a formula's DAG
// @Tags         Formula
// @Param        request body	sio.FormulaInput true "Input Formula"
// @Success      200  {object}  sio.IntResult
// @Router       /formula/nodes [post]
func handleFormulaNodes(w http.ResponseWriter, r *http.Request) {
	fac := formula.NewFactory()
	formula, err := parseFormulaInput(w, r, fac)
	if !err {
		return
	}
	sio.WriteIntResult(w, r, int64(function.NumberOfNodes(fac, formula)))
}

// @Summary      Computes all variables of a formula
// @Tags         Formula
// @Param        request body	sio.FormulaInput true "Input Formula"
// @Success      200  {object}  sio.StringSetResult
// @Router       /formula/variables [post]
func handleFormulaVariables(w http.ResponseWriter, r *http.Request) {
	fac := formula.NewFactory()
	f, err := parseFormulaInput(w, r, fac)
	if !err {
		return
	}
	vars := formula.Variables(fac, f).Content()
	varStrings := make([]string, len(vars))
	for i, v := range vars {
		varStrings[i] = v.Sprint(fac)
	}
	sio.WriteStringSetResult(w, r, varStrings)
}

// @Summary      Computes all literals of a formula
// @Tags         Formula
// @Param        request body	sio.FormulaInput true "Input Formula"
// @Success      200  {object}  sio.StringSetResult
// @Router       /formula/literals [post]
func handleFormulaLiterals(w http.ResponseWriter, r *http.Request) {
	fac := formula.NewFactory()
	f, err := parseFormulaInput(w, r, fac)
	if !err {
		return
	}
	lits := formula.Literals(fac, f).Content()
	litStrings := make([]string, len(lits))
	for i, l := range lits {
		litStrings[i] = l.Sprint(fac)
	}
	sio.WriteStringSetResult(w, r, litStrings)
}

// @Summary      Computes all sub-formulas of a formula
// @Tags         Formula
// @Param        request body	sio.FormulaInput true "Input Formula"
// @Success      200  {object}  sio.StringSetResult
// @Router       /formula/sub-formulas [post]
func handleFormulaSubFormulas(w http.ResponseWriter, r *http.Request) {
	fac := formula.NewFactory()
	f, err := parseFormulaInput(w, r, fac)
	if !err {
		return
	}
	sf := function.SubNodes(fac, f)
	sfStrings := make([]string, len(sf))
	for i, l := range sf {
		sfStrings[i] = l.Sprint(fac)
	}
	sio.WriteStringSetResult(w, r, sfStrings)
}

// @Summary      Computes how often each variable occurrs in a formula
// @Tags         Formula
// @Param        request body	sio.FormulaInput true "Input Formula"
// @Success      200  {object}  sio.ProfileResult
// @Router       /formula/var-profile [post]
func handleFormulaVarProfile(w http.ResponseWriter, r *http.Request) {
	fac := formula.NewFactory()
	f, err := parseFormulaInput(w, r, fac)
	if !err {
		return
	}
	pr := function.VariableProfile(fac, f)
	profile := make(map[string]int64, len(pr))
	for k, v := range pr {
		profile[k.Sprint(fac)] = int64(v)
	}
	sio.WriteProfileResult(w, r, profile)
}

// @Summary      Computes how often each literal occurrs in a formula
// @Tags         Formula
// @Param        request body	sio.FormulaInput true "Input Formula"
// @Success      200  {object}  sio.ProfileResult
// @Router       /formula/lit-profile [post]
func handleFormulaLitProfile(w http.ResponseWriter, r *http.Request) {
	fac := formula.NewFactory()
	f, err := parseFormulaInput(w, r, fac)
	if !err {
		return
	}
	pr := function.LiteralProfile(fac, f)
	profile := make(map[string]int64, len(pr))
	for k, v := range pr {
		profile[k.Sprint(fac)] = int64(v)
	}
	sio.WriteProfileResult(w, r, profile)
}
