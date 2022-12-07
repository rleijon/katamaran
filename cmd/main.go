package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"fmt"
	. "katamaran/pkg/data"
	"katamaran/pkg/node"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime/pprof"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

func main() {
	profile := flag.Bool("profile", false, "Profile")
	client := flag.Bool("client", false, "Client")
	noReqs := flag.Int("nr", 20000, "Number of reqs")
	noConn := flag.Int("nc", 1, "Number of connections")
	port := flag.Int("port", 5225, "Port")
	remotes := flag.String("remote", "localhost:5226", "Remote address")
	flag.Parse()

	if *profile {
		f, err := os.Create("katamaran.prof")
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	if *client {
		address := "localhost:" + strconv.Itoa(*port)
		var wg sync.WaitGroup
		start := time.Now()
		for i := 0; i < *noConn; i++ {
			wg.Add(1)
			go func() {
				msg := []byte{0x05, 0x05, 0x10, 0x10, 0x16, 0x16, 0xA0, 0xA0, 0x00, 0x01}
				answer := make([]byte, 1)
				reqs := *noReqs / *noConn
				defer wg.Done()
				conn, _ := net.Dial("tcp", address)
				writer := bufio.NewWriter(conn)
				writer.Write(node.AddEntryHeader)
				binary.Write(writer, binary.LittleEndian, int32(reqs))
				for j := 0; j < reqs; j++ {
					e := MarshalBytes(writer, msg)
					if e != nil {
						fmt.Println("Couldn't write bytes:", i, e)
						break
					}
				}
				writer.Flush()
				conn.Read(answer)
			}()
		}
		wg.Wait()
		ms := time.Since(start).Milliseconds()
		print("Elapsed time: ", ms, ", Reqs/sec:", 1000*float64(*noReqs)/float64(ms))
	} else {
		id := CandidateId(strconv.Itoa(*port))
		sender := node.MakeTcpSender(strings.Split(*remotes, ","))
		n := node.MakeNode(id, sender)
		ch := make(chan Message)

		cancelChan := make(chan os.Signal, 1)
		// catch SIGETRM or SIGINTERRUPT
		signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)
		go n.RunNode(ch)
		go node.StartTcpServer("localhost:"+strconv.Itoa(*port), ch)
		sig := <-cancelChan
		log.Printf("Caught signal %v", sig)
	}
}
