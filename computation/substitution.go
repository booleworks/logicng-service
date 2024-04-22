package computation

import (
	"fmt"
	"net/http"

	"github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/transformation"
	"github.com/booleworks/logicng-service/config"
	"github.com/booleworks/logicng-service/sio"
)

func HandleSubstitution(cfg *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		subst := r.PathValue("subst")
		switch subst {
		case "anonymization":
			handleSubstAnonymization(w, r)
		case "variables":
			handleSubstVariable(w, r)
		default:
			sio.WriteError(w, r, sio.ErrUnknownPath(r.URL.Path))
		}
	})
}

// @Summary      Replaces all variables in a formula with anonymous ones
// @Tags         Substitution
// @Param        prefix query string false "Optional prefix for the new variables" default(v)
// @Param        request body	sio.FormulaInput true "Input Formula"
// @Success      200  {object}  sio.FormulaResult
// @Router       /substitution/anonymization [post]
func handleSubstAnonymization(w http.ResponseWriter, r *http.Request) {
	transform(w, r, func(fac formula.Factory, form formula.Formula) (formula.Formula, sio.ServiceError) {
		prefix := "v"
		if pf := r.URL.Query().Get("prefix"); pf != "" {
			prefix = pf
		}
		anon := transformation.NewAnonymizer(fac, prefix)
		return anon.Anonymize(form), nil
	})
}

// @Summary      Replaces variables in a formula by their given substitution formula
// @Tags         Substitution
// @Param        request body	sio.SubstitutionInput true "Input Formula and Substitution"
// @Success      200  {object}  sio.FormulaResult
// @Router       /substitution/variables [post]
func handleSubstVariable(w http.ResponseWriter, r *http.Request) {
	input, err := sio.Unmarshal[sio.SubstitutionInput](r)
	if err != nil {
		sio.WriteError(w, r, err)
		return
	}

	fac := formula.NewFactory()
	parsed, ok := parse(w, r, fac, input.Formula)
	if !ok {
		return
	}
	subst := transformation.NewSubstitution()
	for _, s := range input.Substitution.Mapping {
		replace, ok := parse(w, r, fac, s.Replace)
		if !ok {
			return
		}
		if replace.Sort() != formula.SortLiteral || replace.IsNeg() {
			sio.WriteError(w, r, sio.ErrIllegalInput(fmt.Errorf("replace must be a single variable")))
			return
		}
		with, ok := parse(w, r, fac, s.With)
		if !ok {
			return
		}
		subst.AddVar(formula.Variable(replace), with)
	}

	substituted, substErr := transformation.Substitute(fac, parsed, subst)
	if substErr == nil {
		sio.WriteFormulaResult(w, r, substituted.Sprint(fac))
	} else {
		sio.WriteError(w, r, sio.ErrIllegalInput(substErr))
	}
}
