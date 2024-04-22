package sio

import (
	"net/http"

	"github.com/booleworks/logicng-service/sio/pb"
	"google.golang.org/protobuf/proto"
)

type IntResult struct {
	State ComputationState `json:"state"`
	Value int64            `json:"value"`
}

func (r IntResult) ProtoBuf() (bin []byte, err error) {
	bin, err = proto.Marshal(&pb.IntResult{
		State: r.State.toPB(),
		Value: r.Value,
	})
	return
}

func (IntResult) DeserProtoBuf(data []byte) (IntResult, error) {
	result := &pb.IntResult{}
	if err := proto.Unmarshal(data, result); err != nil {
		return IntResult{}, err
	}
	return IntResult{stateFromPB(result.State), result.Value}, nil
}

func WriteIntResult(w http.ResponseWriter, r *http.Request, value int64) {
	result := IntResult{
		State: ComputationState{Success: true},
		Value: value,
	}
	WriteResult(w, r, result)
}
