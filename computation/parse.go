package computation

import (
	"fmt"
	"net/http"

	"github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/parser"
	"github.com/booleworks/logicng-service/sio"
)

func parseFormulaInput(w http.ResponseWriter, r *http.Request, fac formula.Factory) ([]formula.Formula, bool) {
	input, err := sio.Unmarshal[sio.FormulaInput](r)
	if err != nil {
		sio.WriteError(w, r, err)
		return nil, false
	}
	return parseFormulas(w, r, fac, input.Formulas)
}

func parsePropInput(w http.ResponseWriter, r *http.Request, fac formula.Factory) ([]*formula.StandardProposition, bool) {
	input, err := sio.Unmarshal[sio.FormulaInput](r)
	if err != nil {
		sio.WriteError(w, r, err)
		return nil, false
	}
	return parseProps(w, r, fac, input.Formulas)
}

func parseFormulas(w http.ResponseWriter, r *http.Request, fac formula.Factory, strings []sio.Formula) ([]formula.Formula, bool) {
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

func parse(
	w http.ResponseWriter,
	r *http.Request,
	fac formula.Factory,
	formula sio.Formula,
) (formula.Formula, bool) {
	p := parser.New(fac)
	form, err := p.Parse(formula.Formula)
	if err != nil {
		sio.WriteError(w, r, sio.ErrIllegalInput(err))
		return 0, false
	}
	return form, true
}

func parseProps(
	w http.ResponseWriter,
	r *http.Request,
	fac formula.Factory,
	strings []sio.Formula,
) ([]*formula.StandardProposition, bool) {
	props := make([]*formula.StandardProposition, len(strings))
	for i, f := range strings {
		parsed, ok := parseProp(w, r, fac, f)
		if !ok {
			sio.WriteError(w, r, sio.ErrIllegalInput(fmt.Errorf("could not parse formula '%s'", f)))
			return nil, false
		}
		props[i] = parsed
	}
	return props, true
}

func parseProp(
	w http.ResponseWriter,
	r *http.Request,
	fac formula.Factory,
	input sio.Formula,
) (*formula.StandardProposition, bool) {
	p := parser.New(fac)
	form, err := p.Parse(input.Formula)
	if err != nil {
		sio.WriteError(w, r, sio.ErrIllegalInput(err))
		return nil, false
	}
	return formula.NewStandardProposition(form, input.Description), true
}
