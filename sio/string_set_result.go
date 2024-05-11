package sio

import (
	"net/http"

	"github.com/booleworks/logicng-service/sio/pb"
	"google.golang.org/protobuf/proto"
)

type StringSetResult struct {
	State  ComputationState `json:"state"`
	Values []string         `json:"values" example:"A,B"`
}

func (r StringSetResult) ProtoBuf() ([]byte, error) {
	return proto.Marshal(&pb.StringSetResult{
		Value: r.Values,
	})
}

func (StringSetResult) DeserProtoBuf(data []byte) (StringSetResult, error) {
	result := &pb.StringSetResult{}
	if err := proto.Unmarshal(data, result); err != nil {
		return StringSetResult{}, err
	}
	return StringSetResult{stateFromPB(result.State), result.Value}, nil
}

func WriteStringSetResult(w http.ResponseWriter, r *http.Request, values []string) {
	result := StringSetResult{
		State:  ComputationState{Success: true},
		Values: values,
	}
	WriteResult(w, r, result)
}
