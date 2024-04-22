package sio

import (
	"github.com/booleworks/logicng-service/sio/pb"
	"google.golang.org/protobuf/proto"
)

type MaxSatInput struct {
	HardFormulas []string            `json:"hardFormulas" example:"~(A & B) => C | ~D"`
	SoftFormulas []FormulaWithWeight `json:"softFormulas"`
}

type FormulaWithWeight struct {
	Formula string `json:"formula" example:"~(A & B) => C | ~D"`
	Weight  int64  `json:"weight" example:"3"`
}

func (i MaxSatInput) ProtoBuf() (bin []byte, err error) {
	softFormulas := make([]*pb.FormulaWithWeight, len(i.SoftFormulas))
	for i, f := range i.SoftFormulas {
		softFormulas[i] = &pb.FormulaWithWeight{Formula: f.Formula, Weight: int64(f.Weight)}
	}
	bin, err = proto.Marshal(&pb.MaxSatInput{HardFormulas: i.HardFormulas, SoftFormulas: softFormulas})
	return
}

func (MaxSatInput) DeserProtoBuf(data []byte) (MaxSatInput, error) {
	input := &pb.MaxSatInput{}
	if err := proto.Unmarshal(data, input); err != nil {
		return MaxSatInput{}, err
	}
	softFormulas := make([]FormulaWithWeight, len(input.SoftFormulas))
	for i, f := range input.SoftFormulas {
		softFormulas[i] = FormulaWithWeight{f.Formula, f.Weight}
	}
	return MaxSatInput{input.HardFormulas, softFormulas}, nil
}

func (MaxSatInput) Validate() map[string]string {
	return nil
}

func (f FormulaWithWeight) ProtoBuf() (bin []byte, err error) {
	bin, err = proto.Marshal(&pb.FormulaWithWeight{Formula: f.Formula, Weight: int64(f.Weight)})
	return
}

func (FormulaWithWeight) DeserProtoBuf(data []byte) (FormulaWithWeight, error) {
	fww := &pb.FormulaWithWeight{}
	if err := proto.Unmarshal(data, fww); err != nil {
		return FormulaWithWeight{}, err
	}
	return FormulaWithWeight{fww.Formula, fww.Weight}, nil
}

func (f FormulaWithWeight) Validate() map[string]string {
	if f.Formula == "" {
		return map[string]string{"formula": "required field is empty"}
	}
	if f.Weight < 0 {
		return map[string]string{"weight": "weight must be >= 0"}
	}
	return nil
}
