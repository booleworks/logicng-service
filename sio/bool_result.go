package sio

import (
	"net/http"

	"github.com/booleworks/logicng-service/sio/pb"
	"google.golang.org/protobuf/proto"
)

type BoolResult struct {
	State ComputationState `json:"state"`
	Value bool             `json:"value"`
}

func (r BoolResult) ProtoBuf() (bin []byte, err error) {
	bin, err = proto.Marshal(&pb.BoolResult{
		State: r.State.toPB(),
		Value: r.Value,
	})
	return
}

func (BoolResult) DeserProtoBuf(data []byte) (BoolResult, error) {
	result := &pb.BoolResult{}
	if err := proto.Unmarshal(data, result); err != nil {
		return BoolResult{}, err
	}
	return BoolResult{stateFromPB(result.State), result.Value}, nil
}

func WriteBoolResult(w http.ResponseWriter, r *http.Request, value bool) {
	result := BoolResult{
		State: ComputationState{Success: true},
		Value: value,
	}
	WriteResult(w, r, result)
}
