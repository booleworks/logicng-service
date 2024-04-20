package sio

import (
	"net/http"

	"github.com/booleworks/logicng-service/sio/pb"
	"google.golang.org/protobuf/proto"
)

type FormulaResult struct {
	State   ComputationState `json:"state"`
	Formula string           `json:"formula,omitempty" example:"A & B | ~D"`
}

func (r FormulaResult) ProtoBuf() (bin []byte, err error) {
	bin, err = proto.Marshal(&pb.FormulaResult{
		State:   r.State.toPB(),
		Formula: r.Formula,
	})
	return
}

func (FormulaResult) DeserProtoBuf(data []byte) (FormulaResult, error) {
	result := &pb.FormulaResult{}
	if err := proto.Unmarshal(data, result); err != nil {
		return FormulaResult{}, err
	}
	return FormulaResult{stateFromPB(result.State), result.Formula}, nil
}

func WriteFormulaResult(w http.ResponseWriter, r *http.Request, formula string) {
	result := FormulaResult{
		State:   ComputationState{Success: true},
		Formula: formula,
	}
	WriteResult(w, r, result)
}
