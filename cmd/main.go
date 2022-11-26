package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	. "katamaran/pkg/data"
	"katamaran/pkg/node"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func RunNode(n node.Node, ch chan Message) {
	go func(n node.Node, ch chan Message) {
		for {
			time.Sleep(time.Second)
			ch <- Message{node.TickReq{}, nil}
		}
	}(n, ch)
	for {
		msg := <-ch
		switch msg.value.(type) {
		case node.RequestVotesReq:
			req := msg.value.(node.RequestVotesReq)
			_, voteGranted := n.RequestVote(req.CTerm, req.Id, req.LastIndex, req.LastTerm)
			msg.rspChan <- node.RequestVotesRsp{VoteGranted: voteGranted}
		case node.AppendEntriesReq:
			req := msg.value.(node.AppendEntriesReq)
			//fmt.Println("Received", req)
			term, success := n.AppendEntries(req.CTerm, req.Id, req.LastIndex, req.LastTerm, req.Entries, req.CommitIndex)
			msg.rspChan <- node.AppendEntriesRsp{CTerm: term, Success: success}
		case node.AddEntryReq:
			req := msg.value.(node.AddEntryReq)
			success := n.AddEntry(req.Value)
			msg.rspChan <- success
		case node.TickReq:
			n.Tick()
		}
	}
}

type Message struct {
	value   interface{}
	rspChan chan interface{}
}

type Server struct {
	channel chan Message
}

func (h *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	rspChan := make(chan interface{})
	msg, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println("ERROR", err)
		return
	}
	if strings.Contains(req.URL.Path, "requestVote") {
		var req node.RequestVotesReq
		json.Unmarshal(msg, &req)
		h.channel <- Message{req, rspChan}
	} else if strings.Contains(req.URL.Path, "appendEntry") {
		var req node.AppendEntriesReq
		json.Unmarshal(msg, &req)
		h.channel <- Message{req, rspChan}
	} else if strings.Contains(req.URL.Path, "addEntry") {
		var req node.AddEntryReq
		json.Unmarshal(msg, &req)
		h.channel <- Message{req, rspChan}
	}
	rsp := <-rspChan
	//fmt.Println("Replying", rsp, reflect.TypeOf(rsp).String())
	rspBytes, _ := json.Marshal(rsp)
	w.Write(rspBytes)
}

func main() {
	port := flag.Int("port", 5225, "Port")
	remotes := flag.String("remote", "localhost:5226", "Remote address")
	flag.Parse()

	id := CandidateId(strconv.Itoa(*port))
	sender := node.MakeHttpSender(strings.Split(*remotes, ","))
	n := node.MakeNode(id, &sender)
	ch := make(chan Message)

	go RunNode(n, ch)
	http.ListenAndServe("localhost:"+strconv.Itoa(*port), &Server{ch})
}
