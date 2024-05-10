package test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssignmentEvaluation(t *testing.T) {
	assert := assert.New(t)
	ctx := runServer(t)
	ep := endpoint("assignment/evaluation")
	input := `
    {
      "assignment": {
        "A": true,
        "B": false
      },
      "formulas": [
        {
          "formula": "~(A & B) => C | ~D"
        }
      ]
    }
	`
	response, err := callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	validateJSONBoolResult(t, response, true)

	input = `
    {
      "assignment": {
        "A": true,
        "B": false,
        "C": false,
        "D": true
      },
      "formulas": [
        {
          "formula": "~(A & B) => C | ~D"
        }
      ]
    }
	`
	response, err = callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	validateJSONBoolResult(t, response, false)
}

func TestAssignmentRestriction(t *testing.T) {
	assert := assert.New(t)
	ctx := runServer(t)
	ep := endpoint("assignment/restriction")
	input := `
    {
      "assignment": {
        "A": true
	  },
      "formulas": [
        {
          "formula": "~(A & B) => C | ~D"
        }
      ]
    }
	`
	response, err := callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	validateJSONFormulaResult(t, response, "~B => C | ~D")

	input = `
    {
      "assignment": {
        "A": true,
        "C": false
      },
      "formulas": [
        {
          "formula": "~(A & B) => C | ~D"
        }
      ]
    }
	`
	response, err = callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	validateJSONFormulaResult(t, response, "~B => ~D")

	input = `
    {
      "assignment": {
        "A": true,
        "C": false,
        "D": false
      },
      "formulas": [
        {
          "formula": "~(A & B) => C | ~D"
        }
      ]
    }
	`
	response, err = callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	validateJSONFormulaResult(t, response, "$true")
}
