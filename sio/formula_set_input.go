package sio

import (
	"github.com/booleworks/logicng-service/sio/pb"
	"google.golang.org/protobuf/proto"
)

type FormulaSetInput struct {
	Formulas []string `json:"formulas" example:"~(A & B) => C | ~D,~A | E"`
}

func (i FormulaSetInput) ProtoBuf() (bin []byte, err error) {
	bin, err = proto.Marshal(&pb.FormulaSetInput{Formulas: i.Formulas})
	return
}

func (FormulaSetInput) DeserProtoBuf(data []byte) (FormulaSetInput, error) {
	input := &pb.FormulaSetInput{}
	if err := proto.Unmarshal(data, input); err != nil {
		return FormulaSetInput{}, err
	}
	return FormulaSetInput{input.Formulas}, nil
}

func (i FormulaSetInput) Validate() map[string]string {
	if len(i.Formulas) == 0 {
		return map[string]string{"formulas": "empty list"}
	}
	return nil
}
