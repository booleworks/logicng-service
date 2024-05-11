package sio

import (
	"net/http"

	"github.com/booleworks/logicng-service/sio/pb"
	"google.golang.org/protobuf/proto"
)

type GraphResult struct {
	State ComputationState `json:"state"`
	Nodes []Node           `json:"nodes,omitempty"`
	Edges []Edge           `json:"edges,omitempty"`
}

type Node struct {
	ID    int32  `json:"id" example:"42"`
	Label string `json:"label,omitempty" example:"A & B"`
}

type Edge struct {
	SrcID  int32  `json:"srcID" example:"42"`
	DestID int32  `json:"destID" example:"32"`
	Label  string `json:"label,omitempty" example:"1"`
}

func (r GraphResult) ProtoBuf() ([]byte, error) {
	nodes := make([]*pb.Node, len(r.Nodes))
	for i, n := range r.Nodes {
		nodes[i] = &pb.Node{Id: n.ID, Label: n.Label}
	}
	edges := make([]*pb.Edge, len(r.Edges))
	for i, e := range r.Edges {
		edges[i] = &pb.Edge{SrcId: e.SrcID, DestId: e.DestID, Label: e.Label}
	}
	return proto.Marshal(&pb.GraphResult{
		State: r.State.toPB(),
		Nodes: nodes,
		Edges: edges,
	})
}

func (GraphResult) DeserProtoBuf(data []byte) (GraphResult, error) {
	result := &pb.GraphResult{}
	if err := proto.Unmarshal(data, result); err != nil {
		return GraphResult{}, err
	}
	nodes := make([]Node, len(result.Nodes))
	for i, n := range result.Nodes {
		nodes[i] = Node{n.Id, n.Label}
	}
	edges := make([]Edge, len(result.Edges))
	for i, e := range result.Edges {
		edges[i] = Edge{e.SrcId, e.DestId, e.Label}
	}
	return GraphResult{stateFromPB(result.State), nodes, edges}, nil
}

func WriteGraphResult(w http.ResponseWriter, r *http.Request, nodes []Node, edges []Edge) {
	result := GraphResult{
		State: ComputationState{Success: true},
		Nodes: nodes,
		Edges: edges,
	}
	WriteResult(w, r, result)
}
