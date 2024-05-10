package sio

import (
	"github.com/booleworks/logicng-service/sio/pb"
	"google.golang.org/protobuf/proto"
)

type AssignmentInput struct {
	Formulas   []Formula       `json:"formulas"`
	Assignment map[string]bool `json:"assignment" example:"A:true,B:false"`
}

func (i AssignmentInput) ProtoBuf() (bin []byte, err error) {
	formulas := make([]*pb.Formula, len(i.Formulas))
	for i, f := range i.Formulas {
		formulas[i] = f.ProtoBuf()
	}
	return proto.Marshal(&pb.AssignmentInput{Formulas: formulas, Mapping: i.Assignment})
}

func (AssignmentInput) DeserProtoBuf(data []byte) (AssignmentInput, error) {
	input := &pb.AssignmentInput{}
	if err := proto.Unmarshal(data, input); err != nil {
		return AssignmentInput{}, err
	}
	formulas := make([]Formula, len(input.Formulas))
	for i, f := range input.Formulas {
		formulas[i] = Formula{f.Formula, f.Description}
	}
	return AssignmentInput{Formulas: formulas, Assignment: input.Mapping}, nil
}

func (i AssignmentInput) Validate() map[string]string {
	if len(i.Formulas) == 0 {
		return map[string]string{"formulas": "required field is empty"}
	}
	return nil
}
