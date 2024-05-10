package test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var maxSatInput = `
{
  "hardFormulas": [
	{"formula": "(A & B & ~X & ~D) | (B & C & ~A) | (C & ~D & X)"}
  ],
  "softFormulas": {
	"A": 2,
	"B": 4,
	"C": 8,
	"D": 5,
	"X": 7
	}
}
`

var expected = `{
  "state": {
    "success": true
  },
  "satisfiable": true,
  "optimum": 2,
  "model": [
    "~A",
    "B",
    "C",
    "D",
    "X"
  ]
}
`

func TestMaxSatSolver(t *testing.T) {
	assert := assert.New(t)
	ctx := runServer(t)
	ep := endpoint("solver/maxsat")
	response, err := callServiceJSON(ctx, http.MethodPost, ep, maxSatInput)
	assert.Nil(err)
	body := extractJSONBody(response)
	assert.Equal(expected, body)
}
