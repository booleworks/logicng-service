package sio

import (
	"github.com/booleworks/logicng-service/sio/pb"
	"google.golang.org/protobuf/proto"
)

type FormulaSetVarsInput struct {
	Formulas  []string `json:"formulas" example:"~(A & B) => C | ~D,~A | E"`
	Variables []string `json:"variables" example:"A,C,E"`
}

func (i FormulaSetVarsInput) ProtoBuf() (bin []byte, err error) {
	bin, err = proto.Marshal(&pb.FormulaSetVarsInput{Formulas: i.Formulas, Vars: i.Variables})
	return
}

func (FormulaSetVarsInput) DeserProtoBuf(data []byte) (FormulaSetVarsInput, error) {
	input := &pb.FormulaSetVarsInput{}
	if err := proto.Unmarshal(data, input); err != nil {
		return FormulaSetVarsInput{}, err
	}
	return FormulaSetVarsInput{input.Formulas, input.Vars}, nil
}

func (i FormulaSetVarsInput) Validate() map[string]string {
	if len(i.Formulas) == 0 {
		return map[string]string{"formulas": "empty list"}
	}
	if len(i.Variables) == 0 {
		return map[string]string{"variables": "empty list"}
	}
	return nil
}
