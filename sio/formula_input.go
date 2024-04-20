package sio

import (
	"github.com/booleworks/logicng-service/sio/pb"
	"google.golang.org/protobuf/proto"
)

type FormulaInput struct {
	Formula string `json:"formula" example:"~(A & B) => C | ~D"`
}

func (i FormulaInput) ProtoBuf() (bin []byte, err error) {
	bin, err = proto.Marshal(&pb.FormulaInput{Formula: i.Formula})
	return
}

func (FormulaInput) DeserProtoBuf(data []byte) (FormulaInput, error) {
	input := &pb.FormulaInput{}
	if err := proto.Unmarshal(data, input); err != nil {
		return FormulaInput{}, err
	}
	return FormulaInput{input.Formula}, nil
}

func (i FormulaInput) Validate() map[string]string {
	if i.Formula == "" {
		return map[string]string{"formula": "required field is empty"}
	}
	return nil
}
