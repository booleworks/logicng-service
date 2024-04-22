package sio

import (
	"github.com/booleworks/logicng-service/sio/pb"
	"google.golang.org/protobuf/proto"
)

type AssignmentInput struct {
	Formula    string        `json:"formula" example:"~(A & B) => C | ~D"`
	Assignment AssignmentMap `json:"assignment"`
}

type AssignmentMap struct {
	Mapping []Assignment `json:"mapping"`
}

type Assignment struct {
	Variable string `json:"variable" example:"A"`
	Value    bool   `json:"value" example:"true"`
}

func (i AssignmentInput) ProtoBuf() (bin []byte, err error) {
	bin, err = proto.Marshal(&pb.AssignmentInput{Formula: i.Formula, Mapping: i.Assignment.pb()})
	return
}

func (AssignmentInput) DeserProtoBuf(data []byte) (AssignmentInput, error) {
	input := &pb.AssignmentInput{}
	if err := proto.Unmarshal(data, input); err != nil {
		return AssignmentInput{}, err
	}
	return AssignmentInput{Formula: input.Formula, Assignment: deserAssignmentMap(input.Mapping)}, nil
}

func (i AssignmentInput) Validate() map[string]string {
	if i.Formula == "" {
		return map[string]string{"formula": "required field is empty"}
	}
	return nil
}

func (s AssignmentMap) ProtoBuf() (bin []byte, err error) {
	bin, err = proto.Marshal(s.pb())
	return
}

func (s AssignmentMap) pb() *pb.AssignmentMap {
	substs := make([]*pb.Assignment, len(s.Mapping))
	for i, s := range s.Mapping {
		substs[i] = &pb.Assignment{Variable: s.Variable, Value: s.Value}
	}
	return &pb.AssignmentMap{Mapping: substs}
}

func (AssignmentMap) DeserProtoBuf(data []byte) (AssignmentMap, error) {
	mapping := &pb.AssignmentMap{}
	if err := proto.Unmarshal(data, mapping); err != nil {
		return AssignmentMap{}, err
	}
	return deserAssignmentMap(mapping), nil
}

func deserAssignmentMap(mapping *pb.AssignmentMap) AssignmentMap {
	substs := make([]Assignment, len(mapping.Mapping))
	for i, s := range mapping.Mapping {
		substs[i] = Assignment{Variable: s.Variable, Value: s.Value}
	}
	return AssignmentMap{Mapping: substs}
}

func (s AssignmentMap) Validate() map[string]string {
	return nil
}

func (s Assignment) ProtoBuf() (bin []byte, err error) {
	bin, err = proto.Marshal(&pb.Assignment{Variable: s.Variable, Value: s.Value})
	return
}

func (Assignment) DeserProtoBuf(data []byte) (Assignment, error) {
	ass := &pb.Assignment{}
	if err := proto.Unmarshal(data, ass); err != nil {
		return Assignment{}, err
	}
	return Assignment{Variable: ass.Variable, Value: ass.Value}, nil
}

func (s Assignment) Validate() map[string]string {
	validation := make(map[string]string)
	if s.Variable == "" {
		validation["variable"] = "required field is empty"
	}
	return validation
}
