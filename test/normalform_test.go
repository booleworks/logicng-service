package test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNFTransAIG(t *testing.T) {
	assert := assert.New(t)
	ctx := runServer(t)
	ep := endpoint("normalform/transformation/aig")
	input := jsonFormulaInput("~(A & B => ~C <=> D)")
	response, err := callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	validateJSONFormulaResult(t, response, "~(~(~(A & B & C) & ~D) & ~(~(A & B => ~C) & D))")
}

func TestNFTransNNF(t *testing.T) {
	assert := assert.New(t)
	ctx := runServer(t)
	ep := endpoint("normalform/transformation/nnf")
	input := jsonFormulaInput("~(A & B => ~C <=> D)")
	response, err := callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	validateJSONFormulaResult(t, response, "(A & B & C | ~D) & (~A | ~B | ~C | D)")
}

func TestNFTransCNF(t *testing.T) {
	assert := assert.New(t)
	ctx := runServer(t)

	ep := endpoint("normalform/transformation/cnf")
	input := jsonFormulaInput("~(A & B => ~C <=> D)")
	response, err := callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	validateJSONFormulaResult(t, response, "(A | ~D) & (B | ~D) & (C | ~D) & (~A | ~B | ~C | D)")

	ep = endpoint("normalform/transformation/cnf?algorithm=factorization")
	input = jsonFormulaInput("~(A & B => ~C <=> D)")
	response, err = callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	validateJSONFormulaResult(t, response, "(A | ~D) & (B | ~D) & (C | ~D) & (~A | ~B | ~C | D)")

	ep = endpoint("normalform/transformation/cnf?algorithm=advanced")
	input = jsonFormulaInput("~(A & B => ~C <=> D)")
	response, err = callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	validateJSONFormulaResult(t, response, "(A | ~D) & (B | ~D) & (C | ~D) & (~A | ~B | ~C | D)")

	ep = endpoint("normalform/transformation/cnf?algorithm=bdd")
	input = jsonFormulaInput("~(A & B => ~C <=> D)")
	response, err = callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	validateJSONFormulaResult(t, response, "(A | ~D) & (~A | D | ~B | ~C) & (~A | ~D | B) & (~A | ~D | ~B | C)")

	ep = endpoint("normalform/transformation/cnf?algorithm=pg")
	input = jsonFormulaInput("~(A & B => ~C <=> D)")
	response, err = callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	validateJSONFormulaResult(t, response, "@RESERVED_CNF_1 & @RESERVED_CNF_2 & "+
		"(~@RESERVED_CNF_1 | @RESERVED_CNF_3 | ~D) & (~@RESERVED_CNF_3 | A) & (~@RESERVED_CNF_3 | B) & "+
		"(~@RESERVED_CNF_3 | C) & (~@RESERVED_CNF_2 | ~A | ~B | ~C | D)")

	ep = endpoint("normalform/transformation/cnf?algorithm=tseitin")
	input = jsonFormulaInput("~(A & B => ~C <=> D)")
	response, err = callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	validateJSONFormulaResult(t, response, "(~@RESERVED_CNF_0 | A) & (~@RESERVED_CNF_0 | B) & "+
		"(~@RESERVED_CNF_0 | C) & (@RESERVED_CNF_0 | ~A | ~B | ~C) & (@RESERVED_CNF_1 | ~@RESERVED_CNF_0) & "+
		"(@RESERVED_CNF_1 | D) & (~@RESERVED_CNF_1 | @RESERVED_CNF_0 | ~D) & (@RESERVED_CNF_2 | A) & "+
		"(@RESERVED_CNF_2 | B) & (@RESERVED_CNF_2 | C) & (@RESERVED_CNF_2 | ~D) & "+
		"(~@RESERVED_CNF_2 | ~A | ~B | ~C | D) & @RESERVED_CNF_1 & @RESERVED_CNF_2")

	ep = endpoint("normalform/transformation/cnf?algorithm=canonical")
	input = jsonFormulaInput("~(A & B => ~C <=> D)")
	response, err = callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	validateJSONFormulaResult(t, response, "(A | B | C | ~D) & (A | ~B | C | ~D) & (A | ~B | ~C | ~D) & "+
		"(~A | ~B | ~C | D) & (~A | ~B | C | ~D) & (~A | B | C | ~D) & (~A | B | ~C | ~D) & (A | B | ~C | ~D)")
}

func TestNFTransDNF(t *testing.T) {
	assert := assert.New(t)
	ctx := runServer(t)

	ep := endpoint("normalform/transformation/dnf")
	input := jsonFormulaInput("~(A & B => ~C <=> D)")
	response, err := callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	validateJSONFormulaResult(t, response, "D & A & B & C | ~A & ~D | ~B & ~D | ~C & ~D")

	ep = endpoint("normalform/transformation/dnf?algorithm=factorization")
	input = jsonFormulaInput("~(A & B => ~C <=> D)")
	response, err = callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	validateJSONFormulaResult(t, response, "D & A & B & C | ~A & ~D | ~B & ~D | ~C & ~D")

	ep = endpoint("normalform/transformation/dnf?algorithm=canonical")
	input = jsonFormulaInput("~(A & B => ~C <=> D)")
	response, err = callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	validateJSONFormulaResult(t, response, "~A & ~B & ~C & ~D | ~A & B & ~C & ~D | ~A & B & C & ~D | "+
		"A & B & C & D | A & B & ~C & ~D | A & ~B & ~C & ~D | A & ~B & C & ~D | ~A & ~B & C & ~D")

	ep = endpoint("normalform/transformation/dnf?algorithm=bdd")
	input = jsonFormulaInput("~(A & B => ~C <=> D)")
	response, err = callServiceJSON(ctx, http.MethodPost, ep, input)
	assert.Nil(err)
	validateJSONFormulaResult(t, response, "~A & ~D | A & ~D & ~B | A & ~D & B & ~C | A & D & B & C")
}

func TestNFTransProtoBuf(t *testing.T) {
	assert := assert.New(t)
	ctx := runServer(t)
	ep := endpoint("normalform/transformation/nnf")
	body := pbFormulaInput("~(A & B => ~C <=> D)")
	response, err := callServiceProtoBuf(ctx, http.MethodPost, ep, body)
	assert.Nil(err)
	validateProtoBufFormulaResult(t, response, "(A & B & C | ~D) & (~A | ~B | ~C | D)")
}
