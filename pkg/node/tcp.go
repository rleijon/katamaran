package node

import (
	"bytes"
	"fmt"
	. "katamaran/pkg/data"
	"log"
	"net"
)

var RequestAllVotesHeader = []byte{0x01}
var AppendEntriesHeader = []byte{0x02}
var AddEntryHeader = []byte{0x03}

type TcpNodeClient struct {
	connections map[string]net.Conn
}

func MakeTcpSender(allNodes []string) Sender {
	connections := map[string]net.Conn{}
	for _, v := range allNodes {
		var err error
		connections[v], err = net.Dial("tcp", v)
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}
	return &TcpNodeClient{connections: connections}
}

func (h *TcpNodeClient) GetCandidates() []CandidateId {
	candidates := make([]CandidateId, len(h.connections))
	i := 0
	for k, _ := range h.connections {
		candidates[i] = CandidateId(k)
		i++
	}
	return candidates
}

func (h *TcpNodeClient) RequestAllVotes(term Term, candidateId CandidateId, lastIndex Index, lastTerm Term) bool {
	values := RequestVotesReq{CTerm: term, Id: candidateId, LastIndex: lastIndex, LastTerm: lastTerm}
	var buf bytes.Buffer
	values.Marshal(&buf)
	bts := buf.Bytes()
	votes := 0
	for _, v := range h.connections {
		v.Write(RequestAllVotesHeader)
		_, e0 := v.Write(bts)
		if e0 != nil {
			fmt.Println("Error", e0)
			continue
		}
		var response RequestVotesRsp
		response.Unmarshal(v)
		if response.VoteGranted {
			votes++
		}
	}
	fmt.Println("Got", votes, "votes")
	return votes+1 > len(h.connections)/2
}

func (h *TcpNodeClient) SendAppendEntries(node string, term Term, leaderId CandidateId, prevIndex Index, prevTerm Term, entries []Entry, leaderCommit Index) (Term, bool) {
	values := AppendEntriesReq{CTerm: term, Id: leaderId, LastIndex: prevIndex, LastTerm: prevTerm, Entries: entries, CommitIndex: leaderCommit}
	conn := h.connections[node]

	conn.Write(AppendEntriesHeader)
	values.Marshal(conn)
	var response AppendEntriesRsp
	response.Unmarshal(conn)
	return response.CTerm, response.Success
}

func StartTcpServer(address string, ch chan Message) {
	for {
		l, err := net.Listen("tcp", address)
		if err != nil {
			log.Fatal(err)
		}
		conn, _ := l.Accept()
		go ServeTcpConnection(conn, ch)
	}
}

func ServeTcpConnection(conn net.Conn, ch chan Message) {
	header := make([]byte, 1)
	rspChan := make(chan Request)
	for {
		conn.Read(header)
		if header[0] == RequestAllVotesHeader[0] {
			var req RequestVotesReq
			req.Unmarshal(conn)
			ch <- Message{Value: &req, RspChan: rspChan}
			rsp := <-rspChan
			rsp.Marshal(conn)
		} else if header[0] == AppendEntriesHeader[0] {
			rspChan := make(chan Request)
			var req AppendEntriesReq
			req.Unmarshal(conn)
			ch <- Message{Value: &req, RspChan: rspChan}
			rsp := <-rspChan
			rsp.Marshal(conn)
		} else if header[0] == AddEntryHeader[0] {
			var req AddEntryReq
			req.Unmarshal(conn)
			ch <- Message{Value: &req, RspChan: nil}
		} else {
			fmt.Println("Unknown header ", header[0])
		}
	}
}
