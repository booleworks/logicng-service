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
	UnsatCore   []string         `json:"unsatCore,omitempty" example:"A & B, ~B, ~A"`
}

func (r SatResult) ProtoBuf() (bin []byte, err error) {
	bin, err = proto.Marshal(&pb.SatResult{
		State:       r.State.toPB(),
		Satisfiable: r.Satisfiable,
		Model:       r.Model,
		UnsatCore:   r.UnsatCore,
	})
	return
}

func (SatResult) DeserProtoBuf(data []byte) (SatResult, error) {
	res := &pb.SatResult{}
	if err := proto.Unmarshal(data, res); err != nil {
		return SatResult{}, err
	}
	return SatResult{stateFromPB(res.State), res.Satisfiable, res.Model, res.UnsatCore}, nil
}

func WriteSatResult(w http.ResponseWriter, r *http.Request, sat bool, model []string, unsatCore []string) {
	result := SatResult{
		State:       ComputationState{Success: true},
		Satisfiable: sat,
		Model:       model,
		UnsatCore:   unsatCore,
	}
	WriteResult(w, r, result)
}
