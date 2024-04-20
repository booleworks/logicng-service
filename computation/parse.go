package computation

import (
	"net/http"

	"github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/parser"
	"github.com/booleworks/logicng-service/sio"
)

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
