package main

import (
	"flag"
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
			ch <- Message{Value: &TickReq{}, RspChan: nil}
		}
	}(n, ch)
	for {
		msg := <-ch
		switch msg.Value.(type) {
		case *RequestVotesReq:
			req := msg.Value.(*RequestVotesReq)
			_, voteGranted := n.RequestVote(req.CTerm, req.Id, req.LastIndex, req.LastTerm)
			msg.RspChan <- &RequestVotesRsp{VoteGranted: voteGranted}
		case *AppendEntriesReq:
			req := msg.Value.(*AppendEntriesReq)
			term, success := n.AppendEntries(req.CTerm, req.Id, req.LastIndex, req.LastTerm, req.Entries, req.CommitIndex)
			msg.RspChan <- &AppendEntriesRsp{CTerm: term, Success: success}
		case *AddEntryReq:
			req := msg.Value.(*AddEntryReq)
			n.AddEntry(req.Value)
		case *TickReq:
			n.Tick()
		}
	}
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
	node.StartHttpServer("localhost:"+strconv.Itoa(*port), ch)
}
