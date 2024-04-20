package sio

import (
	"net/http"

	"github.com/booleworks/logicng-service/sio/pb"
	"google.golang.org/protobuf/proto"
)

type PredicateResult struct {
	State ComputationState `json:"state"`
	Holds bool             `json:"holds"`
}

func (r PredicateResult) ProtoBuf() (bin []byte, err error) {
	bin, err = proto.Marshal(&pb.PredicateResult{
		State: r.State.toPB(),
		Holds: r.Holds,
	})
	return
}

func (PredicateResult) DeserProtoBuf(data []byte) (PredicateResult, error) {
	result := &pb.PredicateResult{}
	if err := proto.Unmarshal(data, result); err != nil {
		return PredicateResult{}, err
	}
	return PredicateResult{stateFromPB(result.State), result.Holds}, nil
}

func WritePredicateResult(w http.ResponseWriter, r *http.Request, holds bool) {
	result := PredicateResult{
		State: ComputationState{Success: true},
		Holds: holds,
	}
	WriteResult(w, r, result)
}
