package test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMUS(t *testing.T) {
	assert := assert.New(t)
	ctx := runServer(t)
	ep := endpoint("explanation/smus")
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
  "formulas": [
    {
      "formula": "~A | E"
    },
    {
      "formula": "A"
    },
    {
      "formula": "~E"
    }
  ]
}
`
	assert.Equal(expected, body)
}
