package sio

import (
	"net/http"

	"github.com/booleworks/logicng-service/sio/pb"
	"google.golang.org/protobuf/proto"
)

type StringResult struct {
	State ComputationState `json:"state"`
	Value string           `json:"value"`
}

func (r StringResult) ProtoBuf() ([]byte, error) {
	return proto.Marshal(&pb.StringResult{
		State: r.State.toPB(),
		Value: r.Value,
	})
}

func (StringResult) DeserProtoBuf(data []byte) (StringResult, error) {
	result := &pb.StringResult{}
	if err := proto.Unmarshal(data, result); err != nil {
		return StringResult{}, err
	}
	return StringResult{stateFromPB(result.State), result.Value}, nil
}

func WriteStringResult(w http.ResponseWriter, r *http.Request, value string) {
	result := StringResult{
		State: ComputationState{Success: true},
		Value: value,
	}
	WriteResult(w, r, result)
}

func WriteStringResultAsText(w http.ResponseWriter, r *http.Request, value string) {
	WriteTextResult(w, r, value)
}
