package sio

import (
	"net/http"

	"github.com/booleworks/logicng-service/sio/pb"
	"google.golang.org/protobuf/proto"
)

type FormulaResult struct {
	State    ComputationState `json:"state"`
	Formulas []Formula        `json:"formulas,omitempty"`
}

func (r FormulaResult) ProtoBuf() (bin []byte, err error) {
	formulas := make([]*pb.Formula, len(r.Formulas))
	for i, f := range r.Formulas {
		formulas[i] = f.ProtoBuf()
	}
	bin, err = proto.Marshal(&pb.FormulaResult{
		State:    r.State.toPB(),
		Formulas: formulas,
	})
	return
}

func (FormulaResult) DeserProtoBuf(data []byte) (FormulaResult, error) {
	result := &pb.FormulaResult{}
	if err := proto.Unmarshal(data, result); err != nil {
		return FormulaResult{}, err
	}
	formulas := make([]Formula, len(result.Formulas))
	for i, f := range result.Formulas {
		formulas[i] = Formula{f.Formula, f.Description}
	}
	return FormulaResult{stateFromPB(result.State), formulas}, nil
}

func WriteFormulaResult(w http.ResponseWriter, r *http.Request, formula ...Formula) {
	result := FormulaResult{
		State:    ComputationState{Success: true},
		Formulas: formula,
	}
	WriteResult(w, r, result)
}
