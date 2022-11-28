package node

import (
	"bytes"
	"fmt"
	"io/ioutil"
	. "katamaran/pkg/data"
	"net/http"
	"strings"
)

type HttpNodeClient struct {
	allNodes []string
}

func MakeHttpSender(allNodes []string) HttpNodeClient {
	return HttpNodeClient{allNodes: allNodes}
}

func (h *HttpNodeClient) GetCandidates() []CandidateId {
	candidates := make([]CandidateId, len(h.allNodes))
	for i, v := range h.allNodes {
		candidates[i] = CandidateId(v)
	}
	return candidates
}

func (h *HttpNodeClient) RequestAllVotes(term Term, candidateId CandidateId, lastIndex Index, lastTerm Term) bool {
	values := RequestVotesReq{CTerm: term, Id: candidateId, LastIndex: lastIndex, LastTerm: lastTerm}
	var buf bytes.Buffer
	values.Marshal(&buf)
	votes := 0
	for _, v := range h.allNodes {
		rsp, e0 := http.Post("http://"+v+"/requestVote", "application/json", &buf)
		if e0 != nil {
			continue
		}
		var response RequestVotesRsp
		response.Unmarshal(rsp.Body)
		if response.VoteGranted {
			votes++
		}
	}
	fmt.Println("Got", votes, "votes")
	return votes+1 > len(h.allNodes)/2
}

func (h *HttpNodeClient) SendAppendEntries(node string, term Term, leaderId CandidateId, prevIndex Index, prevTerm Term, entries []Entry, leaderCommit Index) (Term, bool) {
	values := AppendEntriesReq{CTerm: term, Id: leaderId, LastIndex: prevIndex, LastTerm: prevTerm, Entries: entries, CommitIndex: leaderCommit}
	var buf bytes.Buffer
	values.Marshal(&buf)
	rsp, e0 := http.Post("http://"+node+"/appendEntry", "application/json", &buf)
	if e0 != nil {
		return Term(0), false
	}
	var response AppendEntriesRsp
	response.Unmarshal(rsp.Body)
	return response.CTerm, response.Success
}

type HttpServer struct {
	channel chan Message
}

func StartHttpServer(address string, ch chan Message) {
	http.ListenAndServe(
		address,
		&HttpServer{channel: ch},
	)
}

func (h *HttpServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	rspChan := make(chan Request)
	body := req.Body
	var buf bytes.Buffer
	if strings.Contains(req.URL.Path, "requestVote") {
		var req RequestVotesReq
		req.Unmarshal(body)
		h.channel <- Message{Value: &req, RspChan: rspChan}
		rsp := <-rspChan
		rsp.Marshal(&buf)
		w.Write(buf.Bytes())
	} else if strings.Contains(req.URL.Path, "appendEntry") {
		var req AppendEntriesReq
		req.Unmarshal(body)
		h.channel <- Message{Value: &req, RspChan: rspChan}
		rsp := <-rspChan
		rsp.Marshal(&buf)
		w.Write(buf.Bytes())
	} else if strings.Contains(req.URL.Path, "addEntry") {
		var req AddEntryReq
		req.Value, _ = ioutil.ReadAll(body)
		h.channel <- Message{Value: &req, RspChan: rspChan}
	}
}
