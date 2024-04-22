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

// @Summary      Evaluates a formula with an assignment of variables
// @Tags         Assignment
// @Param        request body	sio.AssignmentInput true "Input formula and assignment"
// @Success      200  {object}  sio.PredicateResult
// @Router       /assignment/evaluation [post]
func handleEvaluation(w http.ResponseWriter, r *http.Request) {
	input, err := sio.Unmarshal[sio.AssignmentInput](r)
	if err != nil {
		sio.WriteError(w, r, err)
		return
	}

	fac := formula.NewFactory()
	parsed, ok := parse(w, r, fac, input.Formula)
	if !ok {
		return
	}
	ass := extractAssignment(fac, input.Assignment)
	sio.WritePredicateResult(w, r, assignment.Evaluate(fac, parsed, ass))
}

// @Summary      Restricts a formula with an assignment of variables
// @Tags         Assignment
// @Param        request body	sio.AssignmentInput true "Input formula and assignment"
// @Success      200  {object}  sio.FormulaResult
// @Router       /assignment/restriction [post]
func handleRestriction(w http.ResponseWriter, r *http.Request) {
	input, err := sio.Unmarshal[sio.AssignmentInput](r)
	if err != nil {
		sio.WriteError(w, r, err)
		return
	}

	fac := formula.NewFactory()
	parsed, ok := parse(w, r, fac, input.Formula)
	if !ok {
		return
	}
	ass := extractAssignment(fac, input.Assignment)
	sio.WriteFormulaResult(w, r, assignment.Restrict(fac, parsed, ass).Sprint(fac))
}

func extractAssignment(fac formula.Factory, input sio.AssignmentMap) *assignment.Assignment {
	ass, _ := assignment.New(fac)
	for _, a := range input.Mapping {
		ass.AddLit(fac, fac.Lit(a.Variable, a.Value))
	}
	return ass
}
