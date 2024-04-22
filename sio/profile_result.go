package sio

import (
	"net/http"

	"github.com/booleworks/logicng-service/sio/pb"
	"google.golang.org/protobuf/proto"
)

type ProfileResult struct {
	State   ComputationState `json:"state"`
	Profile map[string]int64 `json:"profile" example:"A:3,B:4"`
}

func (r ProfileResult) ProtoBuf() (bin []byte, err error) {
	bin, err = proto.Marshal(&pb.ProfileResult{
		State:   r.State.toPB(),
		Profile: r.Profile,
	})
	return
}

func (ProfileResult) DeserProtoBuf(data []byte) (ProfileResult, error) {
	result := &pb.ProfileResult{}
	if err := proto.Unmarshal(data, result); err != nil {
		return ProfileResult{}, err
	}
	return ProfileResult{stateFromPB(result.State), result.Profile}, nil
}

func WriteProfileResult(w http.ResponseWriter, r *http.Request, profile map[string]int64) {
	result := ProfileResult{
		State:   ComputationState{Success: true},
		Profile: profile,
	}
	WriteResult(w, r, result)
}
