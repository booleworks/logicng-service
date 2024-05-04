package test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDNNFCompilation(t *testing.T) {
	assert := assert.New(t)
	ctx := runServer(t)
	ep := endpoint("dnnf/compilation")
	input := jsonFormulaInput("~(A & B) => C | ~D")
	response, err := callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	validateJSONFormulaResult(t, response, "C | ~C & (D & A & B | ~D)")
}
