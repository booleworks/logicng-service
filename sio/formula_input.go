package sio

import (
	"strings"

	"github.com/booleworks/logicng-service/sio/pb"
	"google.golang.org/protobuf/proto"
)

type Formula struct {
	Formula     string `json:"formula" example:"~(A & B) => C | ~D"`
	Description string `json:"description,omitempty" example:"description text"`
}

type FormulaInput struct {
	Formulas []Formula `json:"formulas"`
}

func (i FormulaInput) ProtoBuf() ([]byte, error) {
	formulas := make([]*pb.Formula, len(i.Formulas))
	for i, f := range i.Formulas {
		formulas[i] = f.ProtoBuf()
	}
	return proto.Marshal(&pb.FormulaInput{Formulas: formulas})
}

func (FormulaInput) DeserProtoBuf(data []byte) (FormulaInput, error) {
	input := &pb.FormulaInput{}
	if err := proto.Unmarshal(data, input); err != nil {
		return FormulaInput{}, err
	}
	formulas := make([]Formula, len(input.Formulas))
	for i, f := range input.Formulas {
		formulas[i] = Formula{f.Formula, f.Description}
	}
	return FormulaInput{formulas}, nil
}

func (f Formula) ProtoBuf() *pb.Formula {
	return &pb.Formula{Formula: f.Formula, Description: f.Description}
}

func (i FormulaInput) Validate() map[string]string {
	if len(i.Formulas) == 0 {
		return map[string]string{"formulas": "empty formula list"}
	}
	for _, f := range i.Formulas {
		if strings.TrimSpace(f.Formula) == "" {
			return map[string]string{"formulas": "contains empty formula"}
		}
	}
	return nil
}
