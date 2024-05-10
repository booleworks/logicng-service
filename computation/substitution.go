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
// @Description  If a list of formulas is given, the result is computed for each formula independently.
// @Tags         Substitution
// @Param        prefix query string false "Optional prefix for the new variables" default(v)
// @Param        request body	sio.FormulaInput true "Input formulas"
// @Success      200  {object}  sio.FormulaResult
// @Router       /substitution/anonymization [post]
func handleSubstAnonymization(w http.ResponseWriter, r *http.Request) {
	input, err := sio.Unmarshal[sio.FormulaInput](r)
	if err != nil {
		sio.WriteError(w, r, err)
		return
	}

	fac := formula.NewFactory()
	ps, ok := parseProps(w, r, fac, input.Formulas)
	if !ok {
		return
	}

	prefix := "v"
	if pf := r.URL.Query().Get("prefix"); pf != "" {
		prefix = pf
	}
	anon := transformation.NewAnonymizer(fac, prefix)
	trans := func(fac formula.Factory, p *formula.StandardProposition) (formula.Formula, sio.ServiceError) {
		return anon.Anonymize(p.Formula()), nil
	}
	transformPropostions(w, r, fac, trans, ps)
}

// @Summary      Replaces variables in a formula by their given substitution formula
// @Description  If a list of formulas is given, the result is computed for each formula independently.
// @Tags         Substitution
// @Param        request body	sio.SubstitutionInput true "Input formulas and Substitution"
// @Success      200  {object}  sio.FormulaResult
// @Router       /substitution/variables [post]
func handleSubstVariable(w http.ResponseWriter, r *http.Request) {
	input, err := sio.Unmarshal[sio.SubstitutionInput](r)
	if err != nil {
		sio.WriteError(w, r, err)
		return
	}

	fac := formula.NewFactory()
	ps, ok := parseProps(w, r, fac, input.Formulas)
	if !ok {
		return
	}

	subst, ok := extractSubst(w, r, fac, input.Substitution)
	if !ok {
		return
	}
	trans := func(fac formula.Factory, p *formula.StandardProposition) (formula.Formula, sio.ServiceError) {
		res, err := transformation.Substitute(fac, p.Formula(), subst)
		if err != nil {
			return 0, sio.ErrIllegalInput(err)
		}
		return res, nil
	}
	transformPropostions(w, r, fac, trans, ps)
}

func extractSubst(
	w http.ResponseWriter,
	r *http.Request,
	fac formula.Factory,
	input map[string]string,
) (*transformation.Substitution, bool) {
	subst := transformation.NewSubstitution()
	for v, s := range input {
		replace, ok := parse(w, r, fac, sio.Formula{Formula: v})
		if !ok {
			return nil, false
		}
		if replace.Sort() != formula.SortLiteral || replace.IsNeg() {
			sio.WriteError(w, r, sio.ErrIllegalInput(fmt.Errorf("replace must be a single variable")))
			return nil, false
		}
		with, ok := parse(w, r, fac, sio.Formula{Formula: s})
		if !ok {
			return nil, false
		}
		subst.AddVar(formula.Variable(replace), with)
	}
	return subst, true
}
