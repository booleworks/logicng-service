package computation

import (
	"fmt"
	"net/http"

	"github.com/booleworks/logicng-go/formula"
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

func transform(
	w http.ResponseWriter,
	r *http.Request,
	transformation func(formula.Factory, formula.Formula) (formula.Formula, sio.ServiceError),
) {
	fac := formula.NewFactory()
	parsed, ok := parseFormulaInput(w, r, fac)
	if !ok {
		return
	}
	transformed, err := transformation(fac, parsed)
	if err == nil {
		sio.WriteFormulaResult(w, r, transformed.Sprint(fac))
	} else {
		sio.WriteError(w, r, err)
	}
}

func transformWithTimeout(result formula.Formula, ok bool) (formula.Formula, sio.ServiceError) {
	if ok {
		return result, nil
	} else {
		return 0, sio.ErrTimeout()
	}
}

func holds(
	w http.ResponseWriter,
	r *http.Request,
	predicate func(formula.Factory, formula.Formula) bool,
) {
	input, err := sio.Unmarshal[sio.FormulaInput](r)
	if err != nil {
		sio.WriteError(w, r, err)
		return
	}

	fac := formula.NewFactory()
	parsed, ok := parse(w, r, fac, input.Formula)
	if !ok {
		return
	}
	holds := predicate(fac, parsed)
	sio.WriteBoolResult(w, r, holds)
}
