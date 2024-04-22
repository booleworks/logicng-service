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
        "mapping": [
          {
            "value": true,
            "variable": "A"
          }
        ]
      },
      "formula": "~(A & B) => C | ~D"
    }
	`
	response, err := callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	validateJSONPredicateResult(t, response, true)

	input = `
	{
      "assignment": {
        "mapping": [
          {
            "value": false,
            "variable": "A"
          }
        ]
      },
      "formula": "~(A & B) => C | D"
    }
	`
	response, err = callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	validateJSONPredicateResult(t, response, false)
}

func TestAssignmentRestriction(t *testing.T) {
	assert := assert.New(t)
	ctx := runServer(t)
	ep := endpoint("assignment/restriction")
	input := `
	{
      "assignment": {
        "mapping": [
          {
            "value": true,
            "variable": "A"
          }
        ]
      },
      "formula": "~(A & B) => C | ~D"
    }
	`
	response, err := callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	validateJSONFormulaResult(t, response, "~B => C | ~D")

	input = `
	{
      "assignment": {
        "mapping": [
          {
            "value": true,
            "variable": "A"
          },
	      {
	        "value": false,
	        "variable": "C"
	      }
        ]
      },
      "formula": "~(A & B) => C | ~D"
    }
	`
	response, err = callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	validateJSONFormulaResult(t, response, "~B => ~D")

	input = `
	{
      "assignment": {
        "mapping": [
          {
            "value": true,
            "variable": "A"
          },
	      {
	        "value": false,
	        "variable": "C"
	      },
	      {
	        "value": false,
	        "variable": "D"
	      }

        ]
      },
      "formula": "~(A & B) => C | ~D"
    }
	`
	response, err = callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	validateJSONFormulaResult(t, response, "$true")
}
