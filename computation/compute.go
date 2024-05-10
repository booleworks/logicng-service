package computation

import (
	"net/http"

	"github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-service/sio"
)

func transform(
	w http.ResponseWriter,
	r *http.Request,
	transformation func(formula.Factory, []formula.Formula) (formula.Formula, sio.ServiceError),
) {
	fac := formula.NewFactory()
	fs, ok := parseFormulaInput(w, r, fac)
	if !ok {
		return
	}
	transformed, err := transformation(fac, fs)
	if err == nil {
		sio.WriteFormulaResult(w, r, sio.Formula{Formula: transformed.Sprint(fac)})
	} else {
		sio.WriteError(w, r, err)
	}
}

func transformPerFormula(
	w http.ResponseWriter,
	r *http.Request,
	transformation func(formula.Factory, *formula.StandardProposition) (formula.Formula, sio.ServiceError),
) {
	fac := formula.NewFactory()
	ps, ok := parsePropInput(w, r, fac)
	if !ok {
		return
	}
	transformPropostions(w, r, fac, transformation, ps)
}

func transformPropostions(
	w http.ResponseWriter,
	r *http.Request,
	fac formula.Factory,
	transformation func(formula.Factory, *formula.StandardProposition) (formula.Formula, sio.ServiceError),
	ps []*formula.StandardProposition,
) {
	result := make([]sio.Formula, len(ps))
	var err sio.ServiceError
	for i, p := range ps {
		var transformed formula.Formula
		transformed, err = transformation(fac, p)
		if err != nil {
			break
		}
		result[i] = sio.Formula{Formula: transformed.Sprint(fac), Description: p.Description}
	}

	if err == nil {
		sio.WriteFormulaResult(w, r, result...)
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
	fac := formula.NewFactory()
	formulas, ok := parseFormulaInput(w, r, fac)
	if !ok {
		return
	}
	holds := predicate(fac, fac.And(formulas...))
	sio.WriteBoolResult(w, r, holds)
}
