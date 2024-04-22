package test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var maxSatInput = `
{
  "hardFormulas": [
    "(A & B & ~X & ~D) | (B & C & ~A) | (C & ~D & X)"
  ],
  "softFormulas": [
    {
      "formula": "A",
      "weight": 2
    },
    {
      "formula": "B",
      "weight": 4
    },
    {
      "formula": "C",
      "weight": 8
    },
    {
      "formula": "D",
      "weight": 5
    },
    {
      "formula": "X",
      "weight": 7
    }
  ]
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
