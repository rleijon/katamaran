package node

import (
	"bytes"
	"encoding/binary"
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

func (h *TcpNodeClient) dialAndWrite(node string, header []byte) (net.Conn, error) {
	var err error
	if h.connections[node] == nil {
		h.connections[node], err = net.Dial("tcp", node)
		if err != nil {
			return nil, err
		}
	}
	_, err = h.connections[node].Write(header)
	if err != nil {
		h.connections[node].Close()
		h.connections[node] = nil
		return nil, err
	}
	return h.connections[node], nil
}

func (h *TcpNodeClient) RequestAllVotes(term Term, candidateId CandidateId, lastIndex Index, lastTerm Term) bool {
	values := RequestVotesReq{CTerm: term, Id: candidateId, LastIndex: lastIndex, LastTerm: lastTerm}
	var buf bytes.Buffer
	values.Marshal(&buf)
	bts := buf.Bytes()
	votes := 0
	for k := range h.connections {
		conn, e := h.dialAndWrite(k, RequestAllVotesHeader)
		if e != nil {
			continue
		}
		_, e0 := conn.Write(bts)
		if e0 != nil {
			fmt.Println("Error", e0)
			continue
		}
		var response RequestVotesRsp
		response.Unmarshal(conn)
		if response.VoteGranted {
			votes++
		}
	}
	fmt.Println("Got", votes, "votes")
	return votes+1 > len(h.connections)/2
}

func (h *TcpNodeClient) SendAppendEntries(node string, term Term, leaderId CandidateId, prevIndex Index, prevTerm Term, entries []Entry, leaderCommit Index) (Term, bool) {
	values := AppendEntriesReq{CTerm: term, Id: leaderId, LastIndex: prevIndex, LastTerm: prevTerm, Entries: entries, CommitIndex: leaderCommit}
	conn, e := h.dialAndWrite(node, AppendEntriesHeader)
	if e != nil {
		return 0, false
	}
	conn.Write(AppendEntriesHeader)
	values.Marshal(conn)
	var response AppendEntriesRsp
	response.Unmarshal(conn)
	return response.CTerm, response.Success
}

func StartTcpServer(address string, ch chan Message) {
	l, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, _ := l.Accept()
		go ServeTcpConnection(conn, ch)
	}
}

func ServeTcpConnection(conn net.Conn, ch chan Message) {
	header := make([]byte, 1)
	rspChan := make(chan Request)
	for {
		_, e := conn.Read(header)
		if e != nil {
			return
		}
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
			var noEntries int32
			binary.Read(conn, binary.LittleEndian, &noEntries)
			for i := int32(0); i < noEntries; i++ {
				var req AddEntryReq
				req.Unmarshal(conn)
				ch <- Message{Value: &req, RspChan: nil}
			}
			//Response
			conn.Write([]byte{0x01})
		} else {
			fmt.Println("Unknown header ", header[0])
		}
	}
}
