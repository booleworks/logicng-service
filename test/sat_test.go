package test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSatSatisfiable(t *testing.T) {
	assert := assert.New(t)
	ctx := runServer(t)
	ep := endpoint("solver/sat")
	input := `
    {
      "formulas": [
	    {"formula": "~(A & B) => C | ~D"},
	    {"formula": "~A | E"},
	    {"formula": "A"}
      ]
    }
	`
	response, err := callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	body := extractJSONBody(response)
	expected := `{
  "state": {
    "success": true
  },
  "satisfiable": true,
  "model": [
    "A",
    "~B",
    "~C",
    "~D",
    "E"
  ]
}
`
	assert.Equal(expected, body)
}

func TestSatUnsatisfiable(t *testing.T) {
	assert := assert.New(t)
	ctx := runServer(t)
	ep := endpoint("solver/sat")
	input := `
    {
      "formulas": [
	    {"formula": "~(A & B) => C | ~D"},
	    {"formula": "~A | E"},
	    {"formula": "A"},
	    {"formula": "~E"}
      ]
    }
	`
	response, err := callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	body := extractJSONBody(response)
	expected := `{
  "state": {
    "success": true
  },
  "satisfiable": false
}
`
	assert.Equal(expected, body)

	ep = endpoint("solver/sat?core=true")
	input = `
    {
      "formulas": [
	    {"formula": "~(A & B) => C | ~D"},
	    {"formula": "~A | E"},
	    {"formula": "A"},
	    {"formula": "~E"}
      ]
    }
	`
	response, err = callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	body = extractJSONBody(response)
	expected = `{
  "state": {
    "success": true
  },
  "satisfiable": false,
  "unsatCore": [
    {
      "formula": "~E"
    },
    {
      "formula": "A"
    },
    {
      "formula": "~A | E"
    }
  ]
}
`
	assert.Equal(expected, body)
}

func TestSatBackbone(t *testing.T) {
	assert := assert.New(t)
	ctx := runServer(t)
	ep := endpoint("solver/backbone")
	input := `
    {
      "formulas": [
	    {"formula": "~(A & B) => C | ~D"},
        {"formula": "~A | E"},
	    {"formula": "A"},
	    {"formula": "~D"}
    	  ]
    }
	`
	response, err := callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	body := extractJSONBody(response)
	expected := `{
  "state": {
    "success": true
  },
  "satisfiable": true,
  "positive": [
    "A",
    "E"
  ],
  "negative": [
    "D"
  ],
  "optional": [
    "B",
    "C"
  ]
}
`
	assert.Equal(expected, body)
}

func TestTautology(t *testing.T) {
	assert := assert.New(t)
	ctx := runServer(t)
	ep := endpoint("solver/predicate/tautology")
	input := `
    {
      "formulas": [
        {"formula": "~(A & B) => C | ~D"},
        {"formula": "~A | E"},
        {"formula": "A"}
      ]
    }
    `
	response, err := callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	validateJSONBoolResult(t, response, false)

	input = `
    {
      "formulas": [
        {"formula": "(A & X) | (~A & ~X) | (~A & X) | (A & ~X)"}
      ]
    }
    `
	response, err = callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	validateJSONBoolResult(t, response, true)
}

func TestContradiction(t *testing.T) {
	assert := assert.New(t)
	ctx := runServer(t)
	ep := endpoint("solver/predicate/contradiction")
	input := `
    {
      "formulas": [
        {"formula": "~(A & B) => C | ~D"},
        {"formula": "~A | E"},
        {"formula": "A"}
      ]
    }
    `
	response, err := callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	validateJSONBoolResult(t, response, false)

	input = `
    {
      "formulas": [
        {"formula": "(A | X) & (~A | ~X)"},
        {"formula": "(~A | X) & (A | ~X)"}
      ]
    }
    `
	response, err = callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	validateJSONBoolResult(t, response, true)
}

func TestImplication(t *testing.T) {
	assert := assert.New(t)
	ctx := runServer(t)
	ep := endpoint("solver/predicate/implication")
	input := `
    {
      "formulas": [
        {"formula": "A & B"},
        {"formula": "A & B & C"}
      ]
    }
    `
	response, err := callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	validateJSONBoolResult(t, response, false)

	input = `
    {
      "formulas": [
        {"formula": "A & B & C"},
        {"formula": "A & B"}
      ]
    }
    `
	response, err = callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	validateJSONBoolResult(t, response, true)
}

func TestEquivalence(t *testing.T) {
	assert := assert.New(t)
	ctx := runServer(t)
	ep := endpoint("solver/predicate/implication")
	input := `
    {
      "formulas": [
        {"formula": "A & B"},
        {"formula": "A & B & C"}
      ]
    }
    `
	response, err := callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	validateJSONBoolResult(t, response, false)

	input = `
    {
      "formulas": [
        {"formula": "A => C"},
        {"formula": "~A | C"}
      ]
    }
    `
	response, err = callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	validateJSONBoolResult(t, response, true)
}
