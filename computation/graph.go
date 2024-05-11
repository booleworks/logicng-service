package computation

import (
	"fmt"
	"net/http"

	"github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/graph"
	"github.com/booleworks/logicng-go/graphical"
	"github.com/booleworks/logicng-service/config"
	"github.com/booleworks/logicng-service/sio"
)

// @Summary      Compute a Constraint graph of formulas
// @Description  Takes a list of formulas. Each node represents a variable.  Two nodes are connected if the respective variables occurr in the same formula.
// @Tags         Graph
// @Param        request body	sio.FormulaInput true "Input formulas"
// @Success      200  {object}  sio.GraphResult
// @Router       /graph/constraint [post]
func HandleConstraintGraph(cfg *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fac := formula.NewFactory()
		fs, ok := parseFormulaInput(w, r, fac)
		if !ok {
			return
		}
		cg := graph.GenerateConstraintGraph(fac, fs...)

		nodeMap := make(map[formula.Formula]int32)
		nodes := make([]sio.Node, len(cg.Nodes()))
		edges := make([]sio.Edge, 0, 4)
		for i, n := range cg.Nodes() {
			nodeMap[n] = int32(i)
			nodes[i] = sio.Node{ID: int32(i), Label: n.Sprint(fac)}
			for _, neighbour := range cg.Neighbours(n) {
				if id, ok := nodeMap[neighbour]; ok {
					edges = append(edges, sio.Edge{SrcID: id, DestID: int32(i)})
				}
			}
		}
		sio.WriteGraphResult(w, r, nodes, edges)
	})
}

// @Summary      Compute a Constraint graph of formulas as a graphical representation
// @Description  Takes a list of formulas. Each node represents a variable.  Two nodes are connected if the respective variables occurr in the same formula.
// @Tags         Graph
// @Param        format query string  false "Output format" Enums(graphviz, mermaid) Default(mermaid)
// @Param        request body	sio.FormulaInput true "Input formulas"
// @Success      200  {string}  graph string
// @Router       /graph/constraint/graphical [post]
func HandleConstraintGraphGraphical(cfg *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fac := formula.NewFactory()
		fs, ok := parseFormulaInput(w, r, fac)
		if !ok {
			return
		}
		cg := graph.GenerateConstraintGraph(fac, fs...)
		representation := graph.GenerateGraphicalFormulaGraph(fac, cg, formula.DefaultFormulaGraphicalGenerator())

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
	})
}

// @Summary      Compute clusters of formulas which occurr in the same components of the constraint graph
// @Description  Takes a list of formulas. Each node represents a variable.  Two nodes are connected if the respective variables occurr in the same formula.
// @Tags         Graph
// @Param        request body	sio.FormulaInput true "Input formulas"
// @Success      200  {object}  sio.ComponentResult
// @Router       /graph/components [post]
func HandleGraphComponents(cfg *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fac := formula.NewFactory()
		ps, ok := parsePropInput(w, r, fac)
		if !ok {
			return
		}
		pMap := make(map[formula.Formula]string)
		fs := make([]formula.Formula, len(ps))
		for i, p := range ps {
			pMap[p.Formula()] = p.Description
			fs[i] = p.Formula()
		}
		cg := graph.GenerateConstraintGraph(fac, fs...)
		components := graph.ComputeConnectedComponents(cg)
		clusters := graph.SplitFormulasByComponent(fac, fs, components)

		result := make([][]sio.Formula, len(clusters))
		for i, c := range clusters {
			result[i] = make([]sio.Formula, len(c))
			for j, f := range c {
				result[i][j] = sio.Formula{Formula: f.Sprint(fac), Description: pMap[f]}
			}
		}
		sio.WriteComponentResult(w, r, result)
	})
}
