package test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModelCount(t *testing.T) {
	assert := assert.New(t)
	ctx := runServer(t)
	ep := endpoint("model/counting")
	input := `
    {
      "formulas": [
        "~(A & B) => C | ~D",
        "~A | E",
	    "X | Y | Z",
	    "P & Q => ~V1 | V2 | V"
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
  "value": "4123"
}
`
	assert.Equal(expected, body)

	ep = endpoint("model/counting?algorithm=bdd")
	input = `
    {
      "formulas": [
        "~(A & B) => C | ~D",
        "~A | E",
	    "X | Y | Z",
	    "P & Q => ~V1 | V2 | V"
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
  "value": "4123"
}
`
	assert.Equal(expected, body)

	ep = endpoint("model/counting?algorithm=sat")
	input = `
    {
      "formulas": [
        "~(A & B) => C | ~D",
        "~A | E",
	    "X | Y | Z",
	    "P & Q => ~V1 | V2 | V"
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
  "value": "4123"
}
`
	assert.Equal(expected, body)
}

func TestProjectedModelCount(t *testing.T) {
	assert := assert.New(t)
	ctx := runServer(t)
	ep := endpoint("model/counting/projection")
	input := `
	{
	  "formulas": [
	    "~(A & B) => C | ~D",
	    "~A | E"
	  ],
	  "variables": [
	    "A",
	    "C",
	    "E"
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
  "value": "6"
}
`
	assert.Equal(expected, body)
}
