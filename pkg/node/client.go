package node

import (
	"context"
	"fmt"
	. "katamaran/pkg/data"
	msg "katamaran/pkg/msg/pkg"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type HttpNodeClient struct {
	allNodes    []string
	connections map[string]msg.KatamaranClient
}

func MakeHttpSender(allNodes []string) HttpNodeClient {
	return HttpNodeClient{allNodes: allNodes, connections: map[string]msg.KatamaranClient{}}
}

type RequestVotesRsp struct {
	VoteGranted bool
}

type TickReq struct {
	Value string
}

func (h *HttpNodeClient) redial(v string) bool {
	if h.connections[v] == nil {
		conn, e := grpc.Dial(v, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if e != nil {
			fmt.Println("Failing to dial", v, e)
			return false
		}
		h.connections[v] = msg.NewKatamaranClient(conn)
	}
	return true
}

func (h *HttpNodeClient) GetCandidates() []CandidateId {
	candidates := make([]CandidateId, len(h.allNodes))
	for i, v := range h.allNodes {
		candidates[i] = CandidateId(v)
	}
	return candidates
}

func (h *HttpNodeClient) RequestAllVotes(term Term, candidateId CandidateId, lastIndex Index, lastTerm Term) bool {
	values := msg.RequestVotesReq{Term: int32(term), Id: string(candidateId), LastIndex: int32(lastIndex), LastTerm: int32(lastTerm)}
	votes := 0
	for _, v := range h.allNodes {
		if !h.redial(v) {
			continue
		}
		rsp, e0 := h.connections[v].RequestAllVotes(context.Background(), &values)
		if e0 != nil {
			fmt.Println("Error req votes", e0)
			continue
		}
		if rsp.VoteGranted {
			votes++
		}
	}
	fmt.Println("Got", votes, "votes")
	return votes+1 > len(h.allNodes)/2
}

func (h *HttpNodeClient) SendAppendEntries(node string, term Term, leaderId CandidateId, prevIndex Index, prevTerm Term, entries []Entry, leaderCommit Index) (Term, bool) {
	if !h.redial(node) {
		return Term(0), false
	}
	reqEntries := make([]*msg.Entry, len(entries))
	for i, v := range entries {
		reqEntries[i] = &msg.Entry{Term: int32(v.Term), Index: int32(v.Index), Value: v.Value}
	}
	values := msg.AppendEntriesReq{Term: int32(term), Id: string(leaderId), LastIndex: int32(prevIndex), LastTerm: int32(prevTerm), Entries: reqEntries, CommitIndex: int32(leaderCommit)}
	rsp, e0 := h.connections[node].AppendEntry(context.Background(), &values)
	if e0 != nil {
		fmt.Println("Error app entr", e0)
		return Term(0), false
	}
	return Term(rsp.Term), rsp.Success
}
