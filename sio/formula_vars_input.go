package sio

import (
	"strings"

	"github.com/booleworks/logicng-service/sio/pb"
	"google.golang.org/protobuf/proto"
)

type FormulaVarsInput struct {
	Formulas  []Formula `json:"formulas"`
	Variables []string  `json:"variables" example:"A,C,E"`
}

func (i FormulaVarsInput) ProtoBuf() (bin []byte, err error) {
	formulas := make([]*pb.Formula, len(i.Formulas))
	for i, f := range i.Formulas {
		formulas[i] = &pb.Formula{Formula: f.Formula, Description: f.Description}
	}
	bin, err = proto.Marshal(&pb.FormulaVarsInput{Formulas: formulas, Vars: i.Variables})
	return
}

func (FormulaVarsInput) DeserProtoBuf(data []byte) (FormulaVarsInput, error) {
	input := &pb.FormulaVarsInput{}
	if err := proto.Unmarshal(data, input); err != nil {
		return FormulaVarsInput{}, err
	}
	formulas := make([]Formula, len(input.Formulas))
	for i, f := range input.Formulas {
		formulas[i] = Formula{f.Formula, f.Description}
	}
	return FormulaVarsInput{formulas, input.Vars}, nil
}

func (i FormulaVarsInput) Validate() map[string]string {
	if len(i.Formulas) == 0 {
		return map[string]string{"formulas": "empty list"}
	}
	for _, f := range i.Formulas {
		if strings.TrimSpace(f.Formula) == "" {
			return map[string]string{"formulas": "contains empty formula"}
		}
	}
	if len(i.Variables) == 0 {
		return map[string]string{"variables": "empty list"}
	}
	return nil
}
