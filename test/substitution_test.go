package test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubstAnonymization(t *testing.T) {
	assert := assert.New(t)
	ctx := runServer(t)
	ep := endpoint("substitution/anonymization")
	input := jsonFormulaInput("~(A & B => ~C <=> D)")
	response, err := callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	validateJSONFormulaResult(t, response, "~(v0 & v1 => ~v2 <=> v3)")

	ep = endpoint("substitution/anonymization?prefix=var")
	input = jsonFormulaInput("~(A & B => ~C <=> D)")
	response, err = callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	validateJSONFormulaResult(t, response, "~(var0 & var1 => ~var2 <=> var3)")
}

func TestSubstVariables(t *testing.T) {
	assert := assert.New(t)
	ctx := runServer(t)
	ep := endpoint("substitution/variables")
	input := jsonFormulaInput("~(A & B => ~C <=> D)")
	response, err := callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	validateJSONFormulaResult(t, response, "~(A & B => ~C <=> D)")

	input = `{
	  "formulas": [{"formula": "~(A & B) => C <=> ~D"}],
      "substitution": {
	    "A": "X & ~Y",
	    "D": "~P"
      }
	}`
	response, err = callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	validateJSONFormulaResult(t, response, "~(X & ~Y & B) => C <=> P")
}
