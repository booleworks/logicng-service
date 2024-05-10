package sio

import (
	"net/http"

	"github.com/booleworks/logicng-service/sio/pb"
	"google.golang.org/protobuf/proto"
)

type SatResult struct {
	State       ComputationState `json:"state"`
	Satisfiable bool             `json:"satisfiable"`
	Model       []string         `json:"model,omitempty" example:"A, ~B"`
	UnsatCore   []Formula        `json:"unsatCore,omitempty"`
}

func (r SatResult) ProtoBuf() (bin []byte, err error) {
	core := make([]*pb.Formula, len(r.UnsatCore))
	for i, f := range r.UnsatCore {
		core[i] = f.ProtoBuf()
	}
	return proto.Marshal(&pb.SatResult{
		State:       r.State.toPB(),
		Satisfiable: r.Satisfiable,
		Model:       r.Model,
		UnsatCore:   core,
	})
}

func (SatResult) DeserProtoBuf(data []byte) (SatResult, error) {
	res := &pb.SatResult{}
	if err := proto.Unmarshal(data, res); err != nil {
		return SatResult{}, err
	}
	core := make([]Formula, len(res.UnsatCore))
	for i, f := range res.UnsatCore {
		core[i] = Formula{f.Formula, f.Description}
	}
	return SatResult{stateFromPB(res.State), res.Satisfiable, res.Model, core}, nil
}

func WriteSatResult(w http.ResponseWriter, r *http.Request, sat bool, model []string, unsatCore []Formula) {
	result := SatResult{
		State:       ComputationState{Success: true},
		Satisfiable: sat,
		Model:       model,
		UnsatCore:   unsatCore,
	}
	WriteResult(w, r, result)
}
