package computation

import (
	"net/http"

	"github.com/booleworks/logicng-go/assignment"
	"github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-service/config"
	"github.com/booleworks/logicng-service/sio"
)

func HandleAssignment(cfg *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ass := r.PathValue("ass")
		switch ass {
		case "evaluation":
			handleEvaluation(w, r)
		case "restriction":
			handleRestriction(w, r)
		default:
			sio.WriteError(w, r, sio.ErrUnknownPath(r.URL.Path))
		}
	})
}

// @Summary      Evaluates formulas with an assignment of variables
// @Description  Variables not in the assignment are assumed 'false'.  If a list of formulas is given, the result refers to the conjunction of these formulas.
// @Tags         Assignment
// @Param        request body	sio.AssignmentInput true "Input formulas and variable assignment"
// @Success      200  {object}  sio.BoolResult
// @Router       /assignment/evaluation [post]
func handleEvaluation(w http.ResponseWriter, r *http.Request) {
	input, err := sio.Unmarshal[sio.AssignmentInput](r)
	if err != nil {
		sio.WriteError(w, r, err)
		return
	}

	fac := formula.NewFactory()
	fs, ok := parseFormulas(w, r, fac, input.Formulas)
	if !ok {
		return
	}
	ass := extractAssignment(fac, input.Assignment)
	sio.WriteBoolResult(w, r, assignment.Evaluate(fac, fac.And(fs...), ass))
}

// @Summary      Restricts formulas with an assignment of variables
// @Description  If a list of formulas is given, the result is computed for each formula independently.
// @Tags         Assignment
// @Param        request body	sio.AssignmentInput true "Input formulas and variable assignment"
// @Success      200  {object}  sio.FormulaResult
// @Router       /assignment/restriction [post]
func handleRestriction(w http.ResponseWriter, r *http.Request) {
	input, err := sio.Unmarshal[sio.AssignmentInput](r)
	if err != nil {
		sio.WriteError(w, r, err)
		return
	}

	fac := formula.NewFactory()
	ps, ok := parseProps(w, r, fac, input.Formulas)
	if !ok {
		return
	}
	ass := extractAssignment(fac, input.Assignment)

	trans := func(fac formula.Factory, p *formula.StandardProposition) (formula.Formula, sio.ServiceError) {
		return assignment.Restrict(fac, p.Formula(), ass), nil
	}
	transformPropostions(w, r, fac, trans, ps)
}

func extractAssignment(fac formula.Factory, input map[string]bool) *assignment.Assignment {
	ass, _ := assignment.New(fac)
	for v, b := range input {
		ass.AddLit(fac, fac.Lit(v, b))
	}
	return ass
}
