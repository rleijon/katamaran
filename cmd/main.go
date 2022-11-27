package main

import (
	"context"
	"flag"
	. "katamaran/pkg/data"
	msg "katamaran/pkg/msg/pkg"
	"katamaran/pkg/node"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func RunNode(n node.Node, ch chan Message) {
	go func(n node.Node, ch chan Message) {
		for {
			time.Sleep(time.Second)
			ch <- Message{node.TickReq{}, nil}
		}
	}(n, ch)
	for {
		m := <-ch
		switch m.value.(type) {
		case *msg.RequestVotesReq:
			req := m.value.(*msg.RequestVotesReq)
			_, voteGranted := n.RequestVote(Term(req.Term), CandidateId(req.Id), Index(req.LastIndex), Term(req.LastTerm))
			m.rspChan <- &msg.RequestVotesRsp{VoteGranted: voteGranted}
		case *msg.AppendEntriesReq:
			req := m.value.(*msg.AppendEntriesReq)
			//fmt.Println("Received", req)
			term, success := n.AppendEntries(Term(req.Term), CandidateId(req.Id), Index(req.LastIndex), Term(req.LastTerm),
				req.Entries, Index(req.CommitIndex))
			m.rspChan <- &msg.AppendEntriesRsp{Term: int32(term), Success: success}
		case *msg.AddEntryReq:
			req := m.value.(*msg.AddEntryReq)
			n.AddEntry(req.Value)
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
	msg.UnimplementedKatamaranServer
	channel chan Message
}

func (h *Server) AddEntry(ctxt context.Context, req *msg.AddEntryReq) (*msg.Empty, error) {
	h.channel <- Message{req, nil}
	return &msg.Empty{}, nil
}

func (h *Server) AppendEntry(ctxt context.Context, req *msg.AppendEntriesReq) (*msg.AppendEntriesRsp, error) {
	rspChan := make(chan interface{})
	h.channel <- Message{req, rspChan}
	rsp := <-rspChan
	typedRsp := rsp.(*msg.AppendEntriesRsp)
	return typedRsp, nil
}

func (h *Server) RequestAllVotes(ctxt context.Context, req *msg.RequestVotesReq) (*msg.RequestVotesRsp, error) {
	rspChan := make(chan interface{})
	h.channel <- Message{req, rspChan}
	rsp := <-rspChan
	typedRsp := rsp.(*msg.RequestVotesRsp)
	return typedRsp, nil
}

func main() {
	debug := flag.Bool("client", false, "Client")
	port := flag.Int("port", 5225, "Port")
	remotes := flag.String("remote", "localhost:5226", "Remote address")
	flag.Parse()
	if *debug {
		v, _ := grpc.Dial("localhost:"+strconv.Itoa(*port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		client := msg.NewKatamaranClient(v)
		client.AddEntry(context.Background(), &msg.AddEntryReq{Value: []byte("abcdefg")})
		return
	}

	id := CandidateId(strconv.Itoa(*port))
	sender := node.MakeHttpSender(strings.Split(*remotes, ","))
	n := node.MakeNode(id, &sender)
	ch := make(chan Message)

	go RunNode(n, ch)
	lis, err := net.Listen("tcp", "localhost:"+strconv.Itoa(*port))
	if err != nil {
		log.Fatal(err)
	}
	server := grpc.NewServer()
	msg.RegisterKatamaranServer(server, &Server{channel: ch})
	server.Serve(lis)
}
