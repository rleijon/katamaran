package plist

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	. "katamaran/pkg/data"
	"log"
	"os"
)

type MMAPList struct {
	id                CandidateId
	votedFor          *os.File
	votedForCached    *CandidateId
	currentTerm       *os.File
	currentTermCached Term
	logsFile          *os.File
	logs              *bufio.Writer
	logsCached        []Entry
	indexFile         *os.File
	index             *bufio.Writer
	indicesCached     map[Index]int32
}

func MakeMMAPList(id CandidateId) *MMAPList {
	os.MkdirAll(string(id), 0666)
	logs, e0 := os.OpenFile(fmt.Sprintf("%s/log", id), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if e0 != nil {
		log.Fatal(e0)
	}
	indices, e1 := os.OpenFile(fmt.Sprintf("%s/index", id), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if e1 != nil {
		log.Fatal(e1)
	}
	votedFor, e2 := os.OpenFile(fmt.Sprintf("%s/votedFor", id), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if e2 != nil {
		log.Fatal(e2)
	}
	currentTerm, e3 := os.OpenFile(fmt.Sprintf("%s/currentTerm", id), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if e3 != nil {
		log.Fatal(e3)
	}

	list := &MMAPList{
		id:          id,
		votedFor:    votedFor,
		currentTerm: currentTerm,
	}
	list.currentTermCached = list.getCurrentTermInternal()
	list.votedForCached = list.getVotedForInternal()
	list.logsCached = readAllEntriesInternal(logs)
	list.indicesCached = readAllIndicesInternal(logs)
	list.logs = bufio.NewWriter(logs)
	list.logsFile = logs
	list.index = bufio.NewWriter(indices)
	list.indexFile = indices
	fmt.Println("Finished loading ", len(list.logsCached), "items")
	return list
}

func readAllIndicesInternal(index *os.File) map[Index]int32 {
	indexr := bufio.NewReader(index)
	indices := make(map[Index]int32)
	c := 0
	for {
		_, e := indexr.Peek(12)
		if e != nil {
			break
		}
		var entry EntryIndex
		e = entry.Unmarshal(indexr)
		if e != nil {
			fmt.Println("Could not load index:", e)
			break
		}
		indices[entry.Id] = entry.Pos
		c++
		if c%100 == 0 {
			fmt.Println("Loaded", c, "indices", indexr.Buffered())
		}
	}
	return indices
}

func readAllEntriesInternal(logs *os.File) []Entry {
	logsr := bufio.NewReader(logs)
	entries := make([]Entry, 0)
	c := 0
	for {
		_, e := logsr.Peek(12)
		if e != nil {
			break
		}
		var entry Entry
		e = entry.Unmarshal(logsr)
		if e != nil {
			fmt.Println("Could not load entry:", e)
			break
		}
		entries = append(entries, entry)
		c++
		if c%100 == 0 {
			fmt.Println("Loaded", c, "items", logsr.Buffered())
		}
	}
	return entries
}

func (p *MMAPList) Flush() {
	p.logs.Flush()
}

func (n *MMAPList) getCurrentTermInternal() Term {
	b, _ := n.votedFor.Stat()
	if b.Size() == 0 {
		n.SetCurrentTerm(0)
		return n.GetCurrentTerm()
	}
	n.votedFor.Seek(0, io.SeekStart)
	var value int32
	binary.Read(n.currentTerm, binary.LittleEndian, &value)
	return Term(value)
}

func (n *MMAPList) getVotedForInternal() *CandidateId {
	b, err := ioutil.ReadAll(n.votedFor)
	if err != nil {
		fmt.Println("Error")
		return nil
	}
	if len(b) == 0 {
		return nil
	}
	cand := CandidateId(b)
	return &cand
}

func (n *MMAPList) SetVotedFor(votedFor *CandidateId) {
	n.votedForCached = votedFor
	n.votedFor.Truncate(0)
	if votedFor != nil {
		n.votedFor.Write([]byte(*votedFor))
	}
}

func (n *MMAPList) GetVotedFor() *CandidateId {
	return n.votedForCached
}

func (n *MMAPList) SetCurrentTerm(term Term) {
	n.currentTermCached = term
	n.currentTerm.Truncate(0)
	binary.Write(n.currentTerm, binary.LittleEndian, int32(term))
}

func (n *MMAPList) GetCurrentTerm() Term {
	return n.currentTermCached
}

func (n *MMAPList) Add(entry Entry) {
	e := entry.Marshal(n.logs)
	if e != nil {
		fmt.Println("Could not write entry:", e)
	}
	n.logsCached = append(n.logsCached, entry)
	n.addIndex(entry.Index, entry.MarshalledSize())
}

func (n *MMAPList) addIndex(index Index, size int32) {
	entrySize := size
	if len(n.indicesCached) > 0 {
		entrySize += n.indicesCached[index-1]
	}
	idx := EntryIndex{Id: index, Pos: entrySize}
	e := idx.Marshal(n.index)
	n.indicesCached[index] = entrySize
	if e != nil {
		fmt.Println("Could not write index:", e)
	}
}

func (n *MMAPList) Get(index Index) Entry {
	return n.logsCached[index]
}

func (n *MMAPList) GetAllAfter(index Index) []Entry {
	return n.logsCached[index:]
}

func (n *MMAPList) Truncate(index Index) error {
	for i := index + 1; int(i) < len(n.logsCached); i++ {
		delete(n.indicesCached, i)
	}
	n.logsCached = n.logsCached[:index]

	var e error
	e = n.logs.Flush()
	if e != nil {
		log.Fatal(e)
	}
	e = n.logsFile.Truncate(int64(n.indicesCached[index]))
	if e != nil {
		log.Fatal(e)
	}
	e = n.index.Flush()
	if e != nil {
		log.Fatal(e)
	}
	return n.indexFile.Truncate(int64(8 * index))
}

func (n *MMAPList) GetNextIndex() Index {
	return n.GetLastEntry().Index + 1
}

func (n *MMAPList) GetLastEntry() Entry {
	if len(n.logsCached) == 0 {
		return Entry{Index: -1, Term: 0, Value: nil}
	}
	return n.logsCached[len(n.logsCached)-1]
}
