package test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMinimalPrimeImplicant(t *testing.T) {
	assert := assert.New(t)
	ctx := runServer(t)
	ep := endpoint("prime/minimal-implicant")
	input := jsonFormulaInput("(~(A & B) => C | ~D) & (X | Y)")
	response, err := callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	validateJSONFormulaResult(t, response, "C & X")
}

func TestMinimalPrimeCover(t *testing.T) {
	assert := assert.New(t)
	ctx := runServer(t)
	ep := endpoint("prime/minimal-cover")
	input := jsonFormulaInput("(~(A & B) => C | ~D) & (X | Y)")
	response, err := callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	body := extractJSONBody(response)

	expected := `{
  "state": {
    "success": true
  },
  "formulas": [
    "C & Y",
    "C & X",
    "~D & X",
    "A & B & X",
    "A & B & Y",
    "~D & Y"
  ]
}
`
	assert.Equal(expected, body)

	ep = endpoint("prime/minimal-cover?algorithm=min")
	input = jsonFormulaInput("(~(A & B) => C | ~D) & (X | Y)")
	response, err = callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	body = extractJSONBody(response)
	expected = `{
  "state": {
    "success": true
  },
  "formulas": [
    "~D & Y",
    "~D & X",
    "C & Y",
    "C & X",
    "A & B & X",
    "A & B & Y"
  ]
}
`
	assert.Equal(expected, body)
}
