package sio

import (
	"fmt"
	"strings"

	"github.com/booleworks/logicng-service/sio/pb"
	"google.golang.org/protobuf/proto"
)

type MaxSatInput struct {
	HardFormulas []Formula        `json:"hardFormulas"`
	SoftFormulas map[string]int64 `json:"softFormulas" example:"~A:3,~B:4,~C & D:2"`
}

func (i MaxSatInput) ProtoBuf() (bin []byte, err error) {
	hardFormulas := make([]*pb.Formula, len(i.HardFormulas))
	for i, f := range i.HardFormulas {
		hardFormulas[i] = &pb.Formula{Formula: f.Formula, Description: f.Description}
	}
	bin, err = proto.Marshal(&pb.MaxSatInput{HardFormulas: hardFormulas, SoftFormulas: i.SoftFormulas})
	return
}

func (MaxSatInput) DeserProtoBuf(data []byte) (MaxSatInput, error) {
	input := &pb.MaxSatInput{}
	if err := proto.Unmarshal(data, input); err != nil {
		return MaxSatInput{}, err
	}
	hardFormulas := make([]Formula, len(input.HardFormulas))
	for i, f := range input.HardFormulas {
		hardFormulas[i] = Formula{f.Formula, f.Description}
	}
	return MaxSatInput{hardFormulas, input.SoftFormulas}, nil
}

func (i MaxSatInput) Validate() map[string]string {
	fmt.Println(i.SoftFormulas)
	if len(i.SoftFormulas) == 0 {
		return map[string]string{"softFormulas": "required field is empty"}
	}
	for _, f := range i.HardFormulas {
		if strings.TrimSpace(f.Formula) == "" {
			return map[string]string{"hardFormulas": "contains empty formula"}
		}
	}
	for f, w := range i.SoftFormulas {
		if strings.TrimSpace(f) == "" {
			return map[string]string{"softFormulas": "contains empty formula"}
		}
		if w < 0 {
			return map[string]string{"softFormulas": "contains weight < 0"}
		}
	}
	return nil
}
