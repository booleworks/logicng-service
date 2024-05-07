package computation

import (
	"fmt"
	"net/http"

	"github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/parser"
	"github.com/booleworks/logicng-service/sio"
)

func parseFormulaInput(w http.ResponseWriter, r *http.Request, fac formula.Factory) (formula.Formula, bool) {
	input, err := sio.Unmarshal[sio.FormulaInput](r)
	if err != nil {
		sio.WriteError(w, r, err)
		return 0, false
	}

	parsed, ok := parse(w, r, fac, input.Formula)
	if !ok {
		return 0, false
	}
	return parsed, true
}

func parseFormulaSetInput(w http.ResponseWriter, r *http.Request, fac formula.Factory) ([]formula.Formula, bool) {
	input, err := sio.Unmarshal[sio.FormulaSetInput](r)
	if err != nil {
		sio.WriteError(w, r, err)
		return nil, false
	}
	return parseFormulaSet(w, r, fac, input.Formulas)
}

func parse(
	w http.ResponseWriter,
	r *http.Request,
	fac formula.Factory,
	formulaString string,
) (formula.Formula, bool) {
	p := parser.New(fac)
	form, err := p.Parse(formulaString)
	if err != nil {
		sio.WriteError(w, r, sio.ErrIllegalInput(err))
		return 0, false
	}
	return form, true
}

func parseFormulaSet(w http.ResponseWriter, r *http.Request, fac formula.Factory, strings []string) ([]formula.Formula, bool) {
	formulas := make([]formula.Formula, len(strings))
	for i, f := range strings {
		parsed, ok := parse(w, r, fac, f)
		if !ok {
			sio.WriteError(w, r, sio.ErrIllegalInput(fmt.Errorf("could not parse formula '%s'", f)))
			return nil, false
		}
		formulas[i] = parsed
	}
	return formulas, true
}
