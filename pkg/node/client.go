package node

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	. "katamaran/pkg/data"
	"net/http"
)

type HttpNodeClient struct {
	allNodes []string
}

func MakeHttpSender(allNodes []string) HttpNodeClient {
	return HttpNodeClient{allNodes: allNodes}
}

type RequestVotesRsp struct {
	VoteGranted bool
}

type TickReq struct {
	Value string
}

type AddEntryReq struct {
	Value []byte
}

type RequestVotesReq struct {
	CTerm     Term
	Id        CandidateId
	LastIndex Index
	LastTerm  Term
}

type AppendEntriesRsp struct {
	CTerm   Term
	Success bool
}

type AppendEntriesReq struct {
	CTerm       Term
	Id          CandidateId
	LastIndex   Index
	LastTerm    Term
	Entries     []Entry
	CommitIndex Index
}

func (h *HttpNodeClient) GetCandidates() []CandidateId {
	candidates := make([]CandidateId, len(h.allNodes))
	for i, v := range h.allNodes {
		candidates[i] = CandidateId(v)
	}
	return candidates
}

func (h *HttpNodeClient) RequestAllVotes(term Term, candidateId CandidateId, lastIndex Index, lastTerm Term) bool {
	values := RequestVotesReq{term, candidateId, lastIndex, lastTerm}
	jsonValue, _ := json.Marshal(values)
	votes := 0
	for _, v := range h.allNodes {
		rsp, e0 := http.Post("http://"+v+"/requestVote", "application/json", bytes.NewReader(jsonValue))
		if e0 != nil {
			continue
		}
		bts, e1 := ioutil.ReadAll(rsp.Body)
		if e1 != nil {
			continue
		}
		var response RequestVotesRsp
		err := json.Unmarshal(bts, &response)
		if err != nil {
			fmt.Println("Error: ", err)
			continue
		}
		if response.VoteGranted {
			votes++
		}
	}
	fmt.Println("Got", votes, "votes")
	return votes+1 > len(h.allNodes)/2
}

func (h *HttpNodeClient) SendAppendEntries(node string, term Term, leaderId CandidateId, prevIndex Index, prevTerm Term, entries []Entry, leaderCommit Index) (Term, bool) {
	values := AppendEntriesReq{term, leaderId, prevIndex, prevTerm, entries, leaderCommit}
	//fmt.Println("Appending", values)
	jsonValue, _ := json.Marshal(values)
	rsp, e0 := http.Post("http://"+node+"/appendEntry", "application/json", bytes.NewReader(jsonValue))
	if e0 != nil {
		return Term(0), false
	}
	bts, e1 := ioutil.ReadAll(rsp.Body)
	if e1 != nil {
		return Term(0), false
	}
	var response AppendEntriesRsp
	err := json.Unmarshal(bts, &response)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	return response.CTerm, response.Success
}
