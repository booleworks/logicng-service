package computation

import (
	"fmt"
	"net/http"

	"github.com/booleworks/logicng-go/bdd"
	"github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/graphical"
	"github.com/booleworks/logicng-go/handler"
	"github.com/booleworks/logicng-service/config"
	"github.com/booleworks/logicng-service/sio"
)

// @Summary      Compile formulas to a BDD
// @Description  If a list of formulas is given, the BDD of the conjunction of these formulas is computed.
// @Tags         BDD
// @Param        ordering query string  false "Variable ordering" Enums(bfs, dfs, min2max, max2min, force) Default(force)
// @Param        request body	sio.FormulaInput true "Input formulas"
// @Success      200  {object}  sio.GraphResult
// @Router       /bdd/compilation [post]
func HandleBDDCompilation(cfg *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fac := formula.NewFactory()
		bddRes, ok := compileBDD(w, r, cfg, fac)
		if !ok {
			return
		}
		rep := bddRes.NodeRepresentation()
		nodeMap := make(map[bdd.Node]sio.Node)
		nodes := make([]sio.Node, 0)
		edges := make([]sio.Edge, 0)
		walkBDD(rep, &nodeMap, &nodes, &edges)
		sio.WriteGraphResult(w, r, nodes, edges)
	})
}

func walkBDD(node bdd.Node, nodeMap *map[bdd.Node]sio.Node, nodes *[]sio.Node, edges *[]sio.Edge) int32 {
	if n, ok := (*nodeMap)[node]; ok {
		return n.ID
	}
	idx := int32(len(*nodes))
	newNode := sio.Node{ID: idx, Label: node.Label()}
	*nodes = append(*nodes, newNode)
	(*nodeMap)[node] = newNode
	if node.InnerNode() {
		lowNode := walkBDD(node.Low(), nodeMap, nodes, edges)
		highNode := walkBDD(node.High(), nodeMap, nodes, edges)
		*edges = append(*edges, sio.Edge{SrcID: idx, DestID: lowNode, Label: "0"})
		*edges = append(*edges, sio.Edge{SrcID: idx, DestID: highNode, Label: "1"})
	}
	return idx
}

// @Summary      Compile formulas to a BDD an return its graphical representation
// @Description  If a list of formulas is given, the BDD of the conjunction of these formulas is computed.
// @Tags         BDD
// @Param        ordering query string  false "Variable ordering" Enums(bfs, dfs, min2max, max2min, force) Default(force)
// @Param        format query string  false "Output format" Enums(graphviz, mermaid) Default(mermaid)
// @Param        request body	sio.FormulaInput true "Input formulas"
// @Success      200  {string}  graph string
// @Router       /bdd/graphical [post]
func HandleBDDGraphical(cfg *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fac := formula.NewFactory()
		bddRes, ok := compileBDD(w, r, cfg, fac)
		if !ok {
			return
		}
		representation := bdd.GenerateGraphical(bddRes, bdd.DefaultGenerator())
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

func compileBDD(w http.ResponseWriter, r *http.Request, cfg *config.Config, fac formula.Factory) (*bdd.BDD, bool) {
	fs, ok := parseFormulaInput(w, r, fac)
	if !ok {
		return nil, false
	}
	f := fac.And(fs...)

	hdl := bdd.HandlerWithTimeout(*handler.NewTimeoutWithDuration(cfg.SyncComputationTimout))
	var order []formula.Variable
	switch ordering := r.URL.Query().Get("ordering"); ordering {
	case "bfs":
		order = bdd.BFSOrder(fac, f)
	case "dfs":
		order = bdd.DFSOrder(fac, f)
	case "min2max":
		order = bdd.MinToMaxOrder(fac, f)
	case "max2min":
		order = bdd.MaxToMinOrder(fac, f)
	case "force":
		order = bdd.ForceOrder(fac, f)
	default:
		sio.WriteError(w, r, sio.ErrUnknownPath(r.URL.Path))
		return nil, false
	}
	bddRes, ok := bdd.CompileWithVarOrderAndHandler(fac, f, order, hdl)
	if !ok {
		sio.WriteError(w, r, sio.ErrTimeout())
		return nil, false
	}
	return bddRes, true
}
