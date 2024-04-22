package test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodingCC(t *testing.T) {
	assert := assert.New(t)
	ctx := runServer(t)

	ep := endpoint("encoding/cc")
	input := jsonFormulaInput("A + B + C + D <= 1")
	response, err := callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	validateJSONFormulaResult(t, response, "(~A | ~B) & (~A | ~C) & (~A | ~D) & (~B | ~C) & (~B | ~D) & (~C | ~D)")

	ep = endpoint("encoding/cc?algorithm=ladder")
	input = jsonFormulaInput("A + B + C + D <= 1")
	response, err = callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	validateJSONFormulaResult(t, response, "(~A | @RESERVED_CC_0) & (~B | @RESERVED_CC_1) & "+
		"(~@RESERVED_CC_0 | @RESERVED_CC_1) & (~B | ~@RESERVED_CC_0) & (~C | @RESERVED_CC_2) & "+
		"(~@RESERVED_CC_1 | @RESERVED_CC_2) & (~C | ~@RESERVED_CC_1) & (~D | ~@RESERVED_CC_2)")

	ep = endpoint("encoding/cc?algorithm=totalizer")
	input = jsonFormulaInput("A + B + C + D <= 2")
	response, err = callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	validateJSONFormulaResult(t, response, "(~@RESERVED_CC_6 | @RESERVED_CC_0) & (~@RESERVED_CC_7 | @RESERVED_CC_1) "+
		"& (~@RESERVED_CC_4 | @RESERVED_CC_0) & (~@RESERVED_CC_4 | ~@RESERVED_CC_6 | @RESERVED_CC_1) & "+
		"(~@RESERVED_CC_4 | ~@RESERVED_CC_7 | @RESERVED_CC_2) & (~@RESERVED_CC_5 | @RESERVED_CC_1) & "+
		"(~@RESERVED_CC_5 | ~@RESERVED_CC_6 | @RESERVED_CC_2) & (~C | @RESERVED_CC_4) & (~D | @RESERVED_CC_4) & "+
		"(~D | ~C | @RESERVED_CC_5) & (~A | @RESERVED_CC_6) & (~B | @RESERVED_CC_6) & (~B | ~A | @RESERVED_CC_7) & "+
		"~@RESERVED_CC_2 & ~@RESERVED_CC_3")
}

func TestEncodingPBC(t *testing.T) {
	assert := assert.New(t)
	ctx := runServer(t)

	ep := endpoint("encoding/pbc")
	input := jsonFormulaInput("2*A + 3*B >= 2")
	response, err := callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	validateJSONFormulaResult(t, response, "(A | @RESERVED_PBC_0) & (A | @RESERVED_PBC_1) & (~@RESERVED_PBC_0 | "+
		"@RESERVED_PBC_3) & (B | @RESERVED_PBC_3) & (~@RESERVED_PBC_1 | @RESERVED_PBC_4) & (B | @RESERVED_PBC_4) & "+
		"(~@RESERVED_PBC_2 | @RESERVED_PBC_5) & (B | @RESERVED_PBC_5) & (~@RESERVED_PBC_0 | B)")

	ep = endpoint("encoding/pbc?algorithm=adder_networks")
	input = jsonFormulaInput("2*A + 3*B >= 2")
	response, err = callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	validateJSONFormulaResult(t, response, "(A | B | ~@RESERVED_PBC_0) & (~A | ~B | ~@RESERVED_PBC_0) & "+
		"(A | ~B | @RESERVED_PBC_0) & (~A | B | @RESERVED_PBC_0) & (~A | ~@RESERVED_PBC_1) & (~B | ~@RESERVED_PBC_1) "+
		"& (A | B | @RESERVED_PBC_1) & ~@RESERVED_PBC_1")
}
