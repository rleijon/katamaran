package main

import (
	"encoding/binary"
	"flag"
	. "katamaran/pkg/data"
	"katamaran/pkg/node"
	"net"
	"strconv"
	"strings"
	"sync"
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
	client := flag.Bool("client", false, "Client")
	noConn := flag.Int("nc", 1, "Number of connections")
	port := flag.Int("port", 5225, "Port")
	remotes := flag.String("remote", "localhost:5226", "Remote address")
	flag.Parse()

	if *client {
		noReqs := 200_000
		address := "localhost:" + strconv.Itoa(*port)
		var wg sync.WaitGroup
		start := time.Now()
		for i := 0; i < *noConn; i++ {
			wg.Add(1)
			go func() {
				msg := []byte{0x05, 0x05, 0x10, 0x10, 0x16, 0x16, 0xA0, 0xA0, 0x00, 0x01}
				answer := make([]byte, 1)
				reqs := noReqs / *noConn
				defer wg.Done()
				conn, _ := net.Dial("tcp", address)
				conn.Write(node.AddEntryHeader)
				binary.Write(conn, binary.LittleEndian, int32(reqs))
				for j := 0; j < reqs; j++ {
					MarshalBytes(conn, msg)
				}
				conn.Read(answer)
			}()
		}
		wg.Wait()
		ms := time.Since(start).Milliseconds()
		print("Elapsed time: ", ms, ", Reqs/sec:", 1000*float64(noReqs)/float64(ms))
	} else {
		id := CandidateId(strconv.Itoa(*port))
		sender := node.MakeTcpSender(strings.Split(*remotes, ","))
		n := node.MakeNode(id, sender)
		ch := make(chan Message)

		go RunNode(n, ch)
		node.StartTcpServer("localhost:"+strconv.Itoa(*port), ch)
	}
}
