package sio

import (
	"net/http"

	"github.com/booleworks/logicng-service/sio/pb"
	"google.golang.org/protobuf/proto"
)

type ComponentResult struct {
	State      ComputationState `json:"state"`
	Components [][]Formula      `json:"components,omitempty"`
}

func (r ComponentResult) ProtoBuf() ([]byte, error) {
	components := make([]*pb.Component, len(r.Components))
	for i, c := range r.Components {
		formulas := make([]*pb.Formula, len(c))
		for j, f := range c {
			formulas[j] = f.ProtoBuf()
		}
		components[i] = &pb.Component{Formulas: formulas}
	}
	return proto.Marshal(&pb.ComponentResult{
		State:      r.State.toPB(),
		Components: components,
	})
}

func (ComponentResult) DeserProtoBuf(data []byte) (ComponentResult, error) {
	result := &pb.ComponentResult{}
	if err := proto.Unmarshal(data, result); err != nil {
		return ComponentResult{}, err
	}
	components := make([][]Formula, len(result.Components))
	for i, c := range result.Components {
		formulas := make([]Formula, len(c.Formulas))
		for j, f := range c.Formulas {
			formulas[j] = Formula{f.Formula, f.Description}
		}
		components[i] = formulas
	}
	return ComponentResult{stateFromPB(result.State), components}, nil
}

func WriteComponentResult(w http.ResponseWriter, r *http.Request, components [][]Formula) {
	result := ComponentResult{
		State:      ComputationState{Success: true},
		Components: components,
	}
	WriteResult(w, r, result)
}
