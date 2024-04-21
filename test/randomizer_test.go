package test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomizer(t *testing.T) {
	assert := assert.New(t)
	ctx := runServer(t)
	ep := endpoint("randomizer/const?seed=42")
	response, err := callServiceJSON(ctx, http.MethodGet, ep, "")
	assert.Nil(err)
	validateJSONFormulaResult(t, response, "$true")

	ep = endpoint("randomizer/formula?vars=5&seed=42")
	response, err = callServiceJSON(ctx, http.MethodGet, ep, "")
	assert.Nil(err)
	validateJSONFormulaResult(t, response, "(v3 | v1) & v4 & v0 | ~v1 & v4 & ~v2 & ~v0 & (~v3 | v1 | ~v2) & "+
		"(~v4 | ~v3 | ~v1) & (v3 | ~v4)")

	ep = endpoint("randomizer/formula?seed=42")
	response, err = callServiceJSON(ctx, http.MethodGet, ep, "")
	assert.Nil(err)
	validateJSONFormulaResult(t, response, "(v01 | v04 | ~v03) & (~v10 | v09 | v19) | ~v04 & v17 & ~v10 | ~v12 "+
		"| ~v04 | ~v17 | v22 | ~v05 | v09 | ~v13 & ~v04 | v03 & ~v18 | ~v18 | ~v19 | ~v08 & v04 | v14 & v12 & v23 & ~v05")
}
