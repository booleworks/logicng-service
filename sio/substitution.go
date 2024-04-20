package sio

import (
	"github.com/booleworks/logicng-service/sio/pb"
	"google.golang.org/protobuf/proto"
)

type SubstitutionInput struct {
	Formula      string          `json:"formula" example:"~(A & B) => C | ~D"`
	Substitution SubstitutionMap `json:"substitution"`
}

type SubstitutionMap struct {
	Mapping []Substitution `json:"mapping"`
}

type Substitution struct {
	Replace string `json:"replace" example:"A"`
	With    string `json:"with" example:"X & Y"`
}

func (i SubstitutionInput) ProtoBuf() (bin []byte, err error) {
	bin, err = proto.Marshal(&pb.SubstitutionInput{Formula: i.Formula, Mapping: i.Substitution.pb()})
	return
}

func (SubstitutionInput) DeserProtoBuf(data []byte) (SubstitutionInput, error) {
	input := &pb.SubstitutionInput{}
	if err := proto.Unmarshal(data, input); err != nil {
		return SubstitutionInput{}, err
	}
	return SubstitutionInput{Formula: input.Formula, Substitution: deserSubstitutionMap(input.Mapping)}, nil
}

func (i SubstitutionInput) Validate() map[string]string {
	if i.Formula == "" {
		return map[string]string{"formula": "required field is empty"}
	}
	return nil
}

func (s SubstitutionMap) ProtoBuf() (bin []byte, err error) {
	bin, err = proto.Marshal(s.pb())
	return
}

func (s SubstitutionMap) pb() *pb.SubstitutionMap {
	substs := make([]*pb.Substitution, len(s.Mapping))
	for i, s := range s.Mapping {
		substs[i] = &pb.Substitution{Replace: s.Replace, With: s.With}
	}
	return &pb.SubstitutionMap{Mapping: substs}
}

func (SubstitutionMap) DeserProtoBuf(data []byte) (SubstitutionMap, error) {
	mapping := &pb.SubstitutionMap{}
	if err := proto.Unmarshal(data, mapping); err != nil {
		return SubstitutionMap{}, err
	}
	return deserSubstitutionMap(mapping), nil
}

func deserSubstitutionMap(mapping *pb.SubstitutionMap) SubstitutionMap {
	substs := make([]Substitution, len(mapping.Mapping))
	for i, s := range mapping.Mapping {
		substs[i] = Substitution{Replace: s.Replace, With: s.With}
	}
	return SubstitutionMap{Mapping: substs}
}

func (s SubstitutionMap) Validate() map[string]string {
	return nil
}

func (s Substitution) ProtoBuf() (bin []byte, err error) {
	bin, err = proto.Marshal(&pb.Substitution{Replace: s.Replace, With: s.With})
	return
}

func (Substitution) DeserProtoBuf(data []byte) (Substitution, error) {
	subst := &pb.Substitution{}
	if err := proto.Unmarshal(data, subst); err != nil {
		return Substitution{}, err
	}
	return Substitution{Replace: subst.Replace, With: subst.Replace}, nil
}

func (s Substitution) Validate() map[string]string {
	validation := make(map[string]string)
	if s.Replace == "" {
		validation["replace"] = "required field is empty"
	}
	if s.With == "" {
		validation["with"] = "required field is empty"
	}
	return validation
}
