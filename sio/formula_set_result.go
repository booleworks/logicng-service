package sio

import (
	"net/http"

	"github.com/booleworks/logicng-service/sio/pb"
	"google.golang.org/protobuf/proto"
)

type FormulaSetResult struct {
	State    ComputationState `json:"state"`
	Formulas []string         `json:"formulas,omitempty" example:"A & B | ~D,A & E"`
}

func (r FormulaSetResult) ProtoBuf() (bin []byte, err error) {
	bin, err = proto.Marshal(&pb.FormulaSetResult{
		State:    r.State.toPB(),
		Formulas: r.Formulas,
	})
	return
}

func (FormulaSetResult) DeserProtoBuf(data []byte) (FormulaSetResult, error) {
	result := &pb.FormulaSetResult{}
	if err := proto.Unmarshal(data, result); err != nil {
		return FormulaSetResult{}, err
	}
	return FormulaSetResult{stateFromPB(result.State), result.Formulas}, nil
}

func WriteFormulaSetResult(w http.ResponseWriter, r *http.Request, formulas []string) {
	result := FormulaSetResult{
		State:    ComputationState{Success: true},
		Formulas: formulas,
	}
	WriteResult(w, r, result)
}
