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
        "~(A & B) => C | ~D",
        "~A | E",
        "A"
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
        "~(A & B) => C | ~D",
        "~A | E",
        "A",
		"~E"
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
        "~(A & B) => C | ~D",
        "~A | E",
        "A",
		"~E"
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
    "~A | E",
    "A",
    "~E"
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
        "~(A & B) => C | ~D",
        "~A | E", "A", "~D"
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
