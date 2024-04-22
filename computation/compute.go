package computation

import (
	"net/http"

	"github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-service/sio"
)

func transform(
	w http.ResponseWriter,
	r *http.Request,
	transformation func(formula.Factory, formula.Formula) (formula.Formula, sio.ServiceError),
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
