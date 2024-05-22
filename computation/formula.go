package computation

import (
	"fmt"
	"net/http"

	"github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/graphical"
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
		case "graphical":
			handleFormulaGraph(w, r)
		default:
			sio.WriteError(w, r, sio.ErrUnknownPath(r.URL.Path))
		}
	})
}

// @Summary      Compute the depth of a formula's AST
// @Description  If a list of formulas is given, the result refers to the conjunction of these formulas.
// @Tags         Formula
// @Param        request body	sio.FormulaInput true "Input formulas"
// @Success      200  {object}  sio.IntResult
// @Router       /formula/depth [post]
func handleFormulaDepth(w http.ResponseWriter, r *http.Request) {
	fac := formula.NewFactory(true)
	fs, err := parseFormulaInput(w, r, fac)
	if !err {
		return
	}
	sio.WriteIntResult(w, r, int64(formula.FormulaDepth(fac, fac.And(fs...))))
}

// @Summary      Compute the number of atoms of a formula
// @Description  If a list of formulas is given, the result refers to the conjunction of these formulas.
// @Tags         Formula
// @Param        request body	sio.FormulaInput true "Input formulas"
// @Success      200  {object}  sio.IntResult
// @Router       /formula/atoms [post]
func handleFormulaAtoms(w http.ResponseWriter, r *http.Request) {
	fac := formula.NewFactory(true)
	fs, err := parseFormulaInput(w, r, fac)
	if !err {
		return
	}
	sio.WriteIntResult(w, r, int64(formula.NumberOfAtoms(fac, fac.And(fs...))))
}

// @Summary      Compute the number of nodes of a formula's DAG
// @Description  If a list of formulas is given, the result refers to the conjunction of these formulas.
// @Tags         Formula
// @Param        request body	sio.FormulaInput true "Input formulas"
// @Success      200  {object}  sio.IntResult
// @Router       /formula/nodes [post]
func handleFormulaNodes(w http.ResponseWriter, r *http.Request) {
	fac := formula.NewFactory(true)
	fs, err := parseFormulaInput(w, r, fac)
	if !err {
		return
	}
	sio.WriteIntResult(w, r, int64(formula.NumberOfNodes(fac, fac.And(fs...))))
}

// @Summary      Compute all variables of a formula
// @Description  If a list of formulas is given, the result refers to the conjunction of these formulas.
// @Tags         Formula
// @Param        request body	sio.FormulaInput true "Input formulas"
// @Success      200  {object}  sio.StringSetResult
// @Router       /formula/variables [post]
func handleFormulaVariables(w http.ResponseWriter, r *http.Request) {
	fac := formula.NewFactory(true)
	fs, err := parseFormulaInput(w, r, fac)
	if !err {
		return
	}
	vars := formula.Variables(fac, fs...).Content()
	varStrings := make([]string, len(vars))
	for i, v := range vars {
		varStrings[i] = v.Sprint(fac)
	}
	sio.WriteStringSetResult(w, r, varStrings)
}

// @Summary      Compute all literals of a formula
// @Description  If a list of formulas is given, the result refers to the conjunction of these formulas.
// @Tags         Formula
// @Param        request body	sio.FormulaInput true "Input formulas"
// @Success      200  {object}  sio.StringSetResult
// @Router       /formula/literals [post]
func handleFormulaLiterals(w http.ResponseWriter, r *http.Request) {
	fac := formula.NewFactory(true)
	fs, err := parseFormulaInput(w, r, fac)
	if !err {
		return
	}
	lits := formula.Literals(fac, fs...).Content()
	litStrings := make([]string, len(lits))
	for i, l := range lits {
		litStrings[i] = l.Sprint(fac)
	}
	sio.WriteStringSetResult(w, r, litStrings)
}

// @Summary      Compute all sub-formulas of a formula
// @Description  If a list of formulas is given, the result refers to the conjunction of these formulas.
// @Tags         Formula
// @Param        request body	sio.FormulaInput true "Input formulas"
// @Success      200  {object}  sio.FormulaResult
// @Router       /formula/sub-formulas [post]
func handleFormulaSubFormulas(w http.ResponseWriter, r *http.Request) {
	fac := formula.NewFactory(true)
	fs, err := parseFormulaInput(w, r, fac)
	if !err {
		return
	}
	sf := formula.SubNodes(fac, fac.And(fs...))
	result := make([]sio.Formula, len(sf))
	for i, l := range sf {
		result[i] = sio.Formula{Formula: l.Sprint(fac)}
	}
	sio.WriteFormulaResult(w, r, result...)
}

// @Summary      Compute how often each variable occurrs in a formula
// @Description  If a list of formulas is given, the result refers to the conjunction of these formulas.
// @Tags         Formula
// @Param        request body	sio.FormulaInput true "Input formulas"
// @Success      200  {object}  sio.ProfileResult
// @Router       /formula/var-profile [post]
func handleFormulaVarProfile(w http.ResponseWriter, r *http.Request) {
	fac := formula.NewFactory(true)
	fs, err := parseFormulaInput(w, r, fac)
	if !err {
		return
	}
	pr := formula.VariableProfile(fac, fac.And(fs...))
	profile := make(map[string]int64, len(pr))
	for k, v := range pr {
		profile[k.Sprint(fac)] = int64(v)
	}
	sio.WriteProfileResult(w, r, profile)
}

// @Summary      Compute how often each literal occurrs in a formula
// @Description  If a list of formulas is given, the result refers to the conjunction of these formulas.
// @Tags         Formula
// @Param        request body	sio.FormulaInput true "Input formulas"
// @Success      200  {object}  sio.ProfileResult
// @Router       /formula/lit-profile [post]
func handleFormulaLitProfile(w http.ResponseWriter, r *http.Request) {
	fac := formula.NewFactory(true)
	fs, err := parseFormulaInput(w, r, fac)
	if !err {
		return
	}
	pr := formula.LiteralProfile(fac, fac.And(fs...))
	profile := make(map[string]int64, len(pr))
	for k, v := range pr {
		profile[k.Sprint(fac)] = int64(v)
	}
	sio.WriteProfileResult(w, r, profile)
}

// @Summary      Compute a graphical DAG or AST representation of formulas
// @Description  If a list of formulas is given, the result refers to the conjunction of these formulas.
// @Tags         Formula
// @Param        type query string  false "Graph type" Enums(ast, dag) Default(dag)
// @Param        format query string  false "Output format" Enums(graphviz, mermaid) Default(mermaid)
// @Param        request body	sio.FormulaInput true "Input formulas"
// @Success      200  {string}  graph string
// @Router       /formula/graphical [post]
func handleFormulaGraph(w http.ResponseWriter, r *http.Request) {
	fac := formula.NewFactory(true)
	fs, err := parseFormulaInput(w, r, fac)
	if !err {
		return
	}
	f := fac.And(fs...)

	var representation *graphical.Representation
	switch graphType := r.URL.Query().Get("type"); graphType {
	case "dag", "":
		representation = formula.GenerateGraphicalFormulaDAG(fac, f, formula.DefaultFormulaGraphicalGenerator())
	case "ast":
		representation = formula.GenerateGraphicalFormulaAST(fac, f, formula.DefaultFormulaGraphicalGenerator())
	default:
		sio.WriteError(w, r, sio.ErrIllegalInput(fmt.Errorf("unknown graph type '%s'", graphType)))
		return
	}

	var result string
	switch format := r.URL.Query().Get("format"); format {
	case "mermaid", "":
		result = graphical.WriteMermaidToString(representation)
	case "graphviz":
		result = graphical.WriteDotToString(representation)
	default:
		sio.WriteError(w, r, sio.ErrIllegalInput(fmt.Errorf("unknown output format '%s'", format)))
		return
	}
	sio.WriteStringResultAsText(w, r, result)
}
