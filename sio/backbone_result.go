package sio

import (
	"net/http"

	"github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/sat"
	"github.com/booleworks/logicng-service/sio/pb"
	"google.golang.org/protobuf/proto"
)

type BackboneResult struct {
	State       ComputationState `json:"state"`
	Satisfiable bool             `json:"satisfiable"`
	Positive    []string         `json:"positive,omitempty" example:"A, B"`
	Negative    []string         `json:"negative,omitempty" example:"C, D"`
	Optional    []string         `json:"optional,omitempty" example:"X, Y"`
}

func (r BackboneResult) ProtoBuf() ([]byte, error) {
	return proto.Marshal(&pb.BackboneResult{
		State:       r.State.toPB(),
		Satisfiable: r.Satisfiable,
		Positive:    r.Positive,
		Negative:    r.Negative,
		Optional:    r.Optional,
	})
}

func (BackboneResult) DeserProtoBuf(data []byte) (BackboneResult, error) {
	res := &pb.BackboneResult{}
	if err := proto.Unmarshal(data, res); err != nil {
		return BackboneResult{}, err
	}
	return BackboneResult{stateFromPB(res.State), res.Satisfiable, res.Positive, res.Negative, res.Optional}, nil
}

func WriteBackboneResult(w http.ResponseWriter, r *http.Request, fac formula.Factory, bb *sat.Backbone) {
	result := BackboneResult{
		State:       ComputationState{Success: true},
		Satisfiable: bb.Sat,
		Positive:    extractVarList(fac, bb.Positive),
		Negative:    extractVarList(fac, bb.Negative),
		Optional:    extractVarList(fac, bb.Optional),
	}
	WriteResult(w, r, result)
}

func extractVarList(fac formula.Factory, vars []formula.Variable) []string {
	var strings []string
	if len(vars) > 0 {
		strings = make([]string, len(vars))
		for i, v := range vars {
			strings[i] = v.Sprint(fac)
		}
	}
	return strings
}
