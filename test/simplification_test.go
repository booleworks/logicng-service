package test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimplBackbone(t *testing.T) {
	assert := assert.New(t)
	ctx := runServer(t)
	ep := endpoint("simplification/backbone")
	input := jsonFormulaInput("A & (~A | B | C) & (~A | B | ~C)")
	response, err := callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	validateJSONFormulaResult(t, response, "A & B")
}

func TestSimplUnitProp(t *testing.T) {
	assert := assert.New(t)
	ctx := runServer(t)
	ep := endpoint("simplification/unitpropagation")
	input := jsonFormulaInput("A & (~A | B | C) & (~A | B | ~C)")
	response, err := callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	validateJSONFormulaResult(t, response, "(B | C) & (B | ~C) & A")
}

func TestSimplNegation(t *testing.T) {
	assert := assert.New(t)
	ctx := runServer(t)
	ep := endpoint("simplification/negation")
	input := jsonFormulaInput("~A & ~B & ~C & ~D")
	response, err := callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	validateJSONFormulaResult(t, response, "~(A | B | C | D)")
}

func TestSimplDistribution(t *testing.T) {
	assert := assert.New(t)
	ctx := runServer(t)
	ep := endpoint("simplification/distribution")
	input := jsonFormulaInput("(A | B) & (A | C & E) | B & C & D")
	response, err := callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	validateJSONFormulaResult(t, response, "B & C & D | A | B & C & E")
}

func TestSimplFactorOut(t *testing.T) {
	assert := assert.New(t)
	ctx := runServer(t)
	ep := endpoint("simplification/factorout")
	input := jsonFormulaInput("(A | B) & (A | C & E) | B & C & D")
	response, err := callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	validateJSONFormulaResult(t, response, "A | B & C & (E | D)")
}

func TestSimplSubsumption(t *testing.T) {
	assert := assert.New(t)
	ctx := runServer(t)
	ep := endpoint("simplification/subsumption")
	input := jsonFormulaInput("(A | B) & (D | E) & (A | B | C)")
	response, err := callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	validateJSONFormulaResult(t, response, "(A | B) & (D | E)")

	input = jsonFormulaInput("(A & B) | (D & E) | (A & B & C)")
	response, err = callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	validateJSONFormulaResult(t, response, "A & B | D & E")
}

func TestSimplQMC(t *testing.T) {
	assert := assert.New(t)
	ctx := runServer(t)
	ep := endpoint("simplification/qmc")
	input := jsonFormulaInput("(~A & ~B & ~C) | (~A & B & ~C) | (A & ~B & C) | (A & B & C)")
	response, err := callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	validateJSONFormulaResult(t, response, "~A & ~C | A & C")
}

func TestSimplAdvanced(t *testing.T) {
	assert := assert.New(t)
	ctx := runServer(t)

	ep := endpoint("simplification/advanced")
	input := jsonFormulaInput("A&B&F&D&E&~H&~C&~G&~I&~J|A&H&D&C&G&~B&~F&~E&~I&~J|A&D&C&~B&~F&~H&~E&~G&~I&~J|" +
		"A&F&H&D&C&G&~B&~E&~I&~J|A&B&H&~F&~D&~C&~E&~G&~I&~J|A&B&H&D&G&~F&~C&~E&~I&~J|A&H&C&~B&~F&~D&~E&~G&~I&~J|" +
		"A&B&G&~F&~H&~D&~C&~E&~I&~J|A&H&C&E&G&~B&~F&~D&~I&~J|A&C&G&~B&~F&~H&~D&~E&~I&~J|A&B&H&G&~F&~D&~C&~E&~I&~J|" +
		"A&C&~B&~F&~H&~D&~E&~G&~I&~J|A&D&C&G&~B&~F&~H&~E&~I&~J|A&B&D&G&~F&~H&~C&~E&~I&~J|A&H&D&C&~B&~F&~E&~G&~I&~J|" +
		"A&H&D&C&E&G&~B&~F&~I&~J|A&B&D&E&~F&~H&~C&~G&~I&~J|A&C&E&G&~B&~F&~H&~D&~I&~J")
	response, err := callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	validateJSONFormulaResult(t, response, "A & ~I & ~J & (~F & (~B & C & (~(E | G & H) | G & (~D & E | D & H)) | "+
		"B & ~E & ~C & (G | ~D & H)) | D & (~B & ~E & H & C & G | B & E & ~H & ~C & ~G))")

	ep = endpoint("simplification/advanced?factorout=false")
	input = jsonFormulaInput("A&B&F&D&E&~H&~C&~G&~I&~J|A&H&D&C&G&~B&~F&~E&~I&~J|A&D&C&~B&~F&~H&~E&~G&~I&~J|" +
		"A&F&H&D&C&G&~B&~E&~I&~J|A&B&H&~F&~D&~C&~E&~G&~I&~J|A&B&H&D&G&~F&~C&~E&~I&~J|A&H&C&~B&~F&~D&~E&~G&~I&~J|" +
		"A&B&G&~F&~H&~D&~C&~E&~I&~J|A&H&C&E&G&~B&~F&~D&~I&~J|A&C&G&~B&~F&~H&~D&~E&~I&~J|A&B&H&G&~F&~D&~C&~E&~I&~J|" +
		"A&C&~B&~F&~H&~D&~E&~G&~I&~J|A&D&C&G&~B&~F&~H&~E&~I&~J|A&B&D&G&~F&~H&~C&~E&~I&~J|A&H&D&C&~B&~F&~E&~G&~I&~J|" +
		"A&H&D&C&E&G&~B&~F&~I&~J|A&B&D&E&~F&~H&~C&~G&~I&~J|A&C&E&G&~B&~F&~H&~D&~I&~J")
	response, err = callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	validateJSONFormulaResult(t, response, "A & ~I & ~J & (C & ~(B | F | E | G) | B & ~F & ~E & ~C & G | "+
		"C & ~(B | F | E | H) | ~B & ~F & ~D & E & C & G | ~B & D & ~E & H & C & G | ~B & ~F & D & H & C & G | "+
		"B & D & E & ~H & ~C & ~G | B & H & ~(F | D | E | C))")

	ep = endpoint("simplification/advanced?factorout=false&backbone=false")
	input = jsonFormulaInput("A&B&F&D&E&~H&~C&~G&~I&~J|A&H&D&C&G&~B&~F&~E&~I&~J|A&D&C&~B&~F&~H&~E&~G&~I&~J|" +
		"A&F&H&D&C&G&~B&~E&~I&~J|A&B&H&~F&~D&~C&~E&~G&~I&~J|A&B&H&D&G&~F&~C&~E&~I&~J|A&H&C&~B&~F&~D&~E&~G&~I&~J|" +
		"A&B&G&~F&~H&~D&~C&~E&~I&~J|A&H&C&E&G&~B&~F&~D&~I&~J|A&C&G&~B&~F&~H&~D&~E&~I&~J|A&B&H&G&~F&~D&~C&~E&~I&~J|" +
		"A&C&~B&~F&~H&~D&~E&~G&~I&~J|A&D&C&G&~B&~F&~H&~E&~I&~J|A&B&D&G&~F&~H&~C&~E&~I&~J|A&H&D&C&~B&~F&~E&~G&~I&~J|" +
		"A&H&D&C&E&G&~B&~F&~I&~J|A&B&D&E&~F&~H&~C&~G&~I&~J|A&C&E&G&~B&~F&~H&~D&~I&~J")
	response, err = callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	validateJSONFormulaResult(t, response, "A & C & ~(B | F | E | H | I | J) | A & B & G & ~(F | E | C | I | J) | "+
		"A & C & ~(B | F | E | G | I | J) | A & B & H & ~(F | D | E | C | I | J) | A & D & H & C & G & "+
		"~(B | E | I | J) | A & E & C & G & ~(B | F | D | I | J) | A & B & D & E & ~(H | C | G | I | J) | "+
		"A & D & H & C & G & ~(B | F | I | J)")

	ep = endpoint("simplification/advanced?factorout=false&backbone=false&negation=false")
	input = jsonFormulaInput("A&B&F&D&E&~H&~C&~G&~I&~J|A&H&D&C&G&~B&~F&~E&~I&~J|A&D&C&~B&~F&~H&~E&~G&~I&~J|" +
		"A&F&H&D&C&G&~B&~E&~I&~J|A&B&H&~F&~D&~C&~E&~G&~I&~J|A&B&H&D&G&~F&~C&~E&~I&~J|A&H&C&~B&~F&~D&~E&~G&~I&~J|" +
		"A&B&G&~F&~H&~D&~C&~E&~I&~J|A&H&C&E&G&~B&~F&~D&~I&~J|A&C&G&~B&~F&~H&~D&~E&~I&~J|A&B&H&G&~F&~D&~C&~E&~I&~J|" +
		"A&C&~B&~F&~H&~D&~E&~G&~I&~J|A&D&C&G&~B&~F&~H&~E&~I&~J|A&B&D&G&~F&~H&~C&~E&~I&~J|A&H&D&C&~B&~F&~E&~G&~I&~J|" +
		"A&H&D&C&E&G&~B&~F&~I&~J|A&B&D&E&~F&~H&~C&~G&~I&~J|A&C&E&G&~B&~F&~H&~D&~I&~J")
	response, err = callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	validateJSONFormulaResult(t, response, "A & C & ~(B | F | E | H | I | J) | A & B & G & ~(F | E | C | I | J) | "+
		"A & C & ~(B | F | E | G | I | J) | A & B & H & ~(F | D | E | C | I | J) | A & D & H & C & G & "+
		"~(B | E | I | J) | A & E & C & G & ~(B | F | D | I | J) | A & B & D & E & ~(H | C | G | I | J) | "+
		"A & D & H & C & G & ~(B | F | I | J)")
}
