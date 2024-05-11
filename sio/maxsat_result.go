package sio

import (
	"net/http"

	"github.com/booleworks/logicng-service/sio/pb"
	"google.golang.org/protobuf/proto"
)

type MaxSatResult struct {
	State       ComputationState `json:"state"`
	Satisfiable bool             `json:"satisfiable"`
	Optimum     int64            `json:"optimum" example:"7"`
	Model       []string         `json:"model,omitempty" example:"A, ~B"`
}

func (r MaxSatResult) ProtoBuf() ([]byte, error) {
	return proto.Marshal(&pb.MaxSatResult{
		State:       r.State.toPB(),
		Satisfiable: r.Satisfiable,
		Optimum:     r.Optimum,
		Model:       r.Model,
	})
}

func (MaxSatResult) DeserProtoBuf(data []byte) (MaxSatResult, error) {
	res := &pb.MaxSatResult{}
	if err := proto.Unmarshal(data, res); err != nil {
		return MaxSatResult{}, err
	}
	return MaxSatResult{stateFromPB(res.State), res.Satisfiable, res.Optimum, res.Model}, nil
}

func WriteMaxSatResult(w http.ResponseWriter, r *http.Request, sat bool, opt int64, model []string) {
	result := MaxSatResult{
		State:       ComputationState{Success: true},
		Satisfiable: sat,
		Optimum:     opt,
		Model:       model,
	}
	WriteResult(w, r, result)
}
