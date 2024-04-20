package sio

import (
	"github.com/booleworks/logicng-service/sio/pb"
	"google.golang.org/protobuf/proto"
)

type State struct{}

type ServiceInput[T any] interface {
	ProtoBuf() ([]byte, error)
	DeserProtoBuf([]byte) (T, error)
	Validate() map[string]string
}

type ServiceOutput[T any] interface {
	ProtoBuf() ([]byte, error)
	DeserProtoBuf([]byte) (T, error)
}

type ComputationResult struct {
	State ComputationState `json:"state"`
}

func (r ComputationResult) ProtoBuf() (bin []byte, err error) {
	bin, err = proto.Marshal(&pb.FormulaResult{
		State: r.State.toPB(),
	})
	return
}

func (ComputationResult) DeserProtoBuf(data []byte) (ComputationResult, error) {
	result := &pb.FormulaResult{}
	if err := proto.Unmarshal(data, result); err != nil {
		return ComputationResult{}, err
	}
	return ComputationResult{stateFromPB(result.State)}, nil
}

type ComputationState struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty" example:""`
}

func (c ComputationState) toPB() *pb.ComputationState {
	return &pb.ComputationState{
		Success: c.Success,
		Error:   c.Error,
	}
}

func stateFromPB(bin *pb.ComputationState) ComputationState {
	return ComputationState{bin.Success, bin.Error}
}
