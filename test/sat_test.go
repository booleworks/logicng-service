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
