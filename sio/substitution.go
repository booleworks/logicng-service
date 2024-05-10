package sio

import (
	"github.com/booleworks/logicng-service/sio/pb"
	"google.golang.org/protobuf/proto"
)

type SubstitutionInput struct {
	Formulas     []Formula         `json:"formulas"`
	Substitution map[string]string `json:"substitution" example:"A:~B,C:X & Y"`
}

func (i SubstitutionInput) ProtoBuf() (bin []byte, err error) {
	formulas := make([]*pb.Formula, len(i.Formulas))
	for i, f := range i.Formulas {
		formulas[i] = f.ProtoBuf()
	}
	return proto.Marshal(&pb.SubstitutionInput{Formulas: formulas, Substitution: i.Substitution})
}

func (SubstitutionInput) DeserProtoBuf(data []byte) (SubstitutionInput, error) {
	input := &pb.SubstitutionInput{}
	if err := proto.Unmarshal(data, input); err != nil {
		return SubstitutionInput{}, err
	}
	formulas := make([]Formula, len(input.Formulas))
	for i, f := range input.Formulas {
		formulas[i] = Formula{f.Formula, f.Description}
	}
	return SubstitutionInput{Formulas: formulas, Substitution: input.Substitution}, nil
}

func (i SubstitutionInput) Validate() map[string]string {
	if len(i.Formulas) == 0 {
		return map[string]string{"formula": "required field is empty"}
	}
	return nil
}
