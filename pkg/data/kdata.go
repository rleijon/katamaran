package data

import (
	"encoding/binary"
	"io"
)

type CandidateId string
type Term int
type Index int
type State int

type Message struct {
	Value   Request
	RspChan chan Request
}

type Request interface {
	Marshal(io.Writer) error
	Unmarshal(io.Reader) error
}

type EntryIndex struct {
	Id  Index
	Pos int32
}

func (r *EntryIndex) Marshal(writer io.Writer) error {
	e0 := MarshalInt32(writer, int32(r.Id))
	if e0 != nil {
		return e0
	}
	return MarshalInt32(writer, int32(r.Pos))
}

func (r *EntryIndex) Unmarshal(reader io.Reader) error {
	id, e0 := UnmarshalInt32(reader)
	if e0 != nil {
		return e0
	}
	pos, e1 := UnmarshalInt32(reader)
	if e1 != nil {
		return e1
	}
	r.Id = Index(id)
	r.Pos = pos
	return nil
}

type Entry struct {
	Value []byte
	Index Index
	Term  Term
}

func (r *Entry) MarshalledSize() int32 {
	return int32(12 + len(r.Value))
}

func (r *Entry) Marshal(writer io.Writer) error {
	e0 := MarshalInt32(writer, int32(r.Index))
	if e0 != nil {
		return e0
	}
	e1 := MarshalInt32(writer, int32(r.Term))
	if e1 != nil {
		return e1
	}
	return MarshalBytes(writer, r.Value)
}

func (r *Entry) Unmarshal(reader io.Reader) error {
	index, e0 := UnmarshalInt32(reader)
	if e0 != nil {
		return e0
	}
	term, e0 := UnmarshalInt32(reader)
	if e0 != nil {
		return e0
	}
	value, e0 := UnmarshalBytes(reader)
	if e0 != nil {
		return e0
	}
	r.Index = Index(index)
	r.Term = Term(term)
	r.Value = value
	return nil
}

type RequestVotesRsp struct {
	VoteGranted bool
}

func (r *RequestVotesRsp) Marshal(writer io.Writer) error {
	return binary.Write(writer, binary.LittleEndian, r.VoteGranted)
}

func (r *RequestVotesRsp) Unmarshal(reader io.Reader) error {
	return binary.Read(reader, binary.LittleEndian, &r.VoteGranted)
}

type TickReq struct {
}

func (r *TickReq) Marshal(writer io.Writer) error {
	return nil
}

func (r *TickReq) Unmarshal(reader io.Reader) error {
	return nil
}

type AddEntryReq struct {
	Value []byte
}

func (r *AddEntryReq) Marshal(writer io.Writer) error {
	return MarshalBytes(writer, r.Value)
}

func (r *AddEntryReq) Unmarshal(reader io.Reader) error {
	var err error
	r.Value, err = UnmarshalBytes(reader)
	return err
}

type AddEntryRsp struct {
	Success bool
}

func (r *AddEntryRsp) Marshal(writer io.Writer) error {
	return binary.Write(writer, binary.LittleEndian, r.Success)
}

func (r *AddEntryRsp) Unmarshal(reader io.Reader) error {
	return binary.Read(reader, binary.LittleEndian, &r.Success)
}

type RequestVotesReq struct {
	CTerm     Term
	Id        CandidateId
	LastIndex Index
	LastTerm  Term
}

func (r *RequestVotesReq) Marshal(writer io.Writer) error {
	if e := MarshalInt32(writer, int32(r.CTerm)); e != nil {
		return e
	}
	if e := MarshalBytes(writer, []byte(r.Id)); e != nil {
		return e
	}
	if e := MarshalInt32(writer, int32(r.LastIndex)); e != nil {
		return e
	}
	return MarshalInt32(writer, int32(r.LastTerm))
}
func (r *RequestVotesReq) Unmarshal(reader io.Reader) error {
	term, e0 := UnmarshalInt32(reader)
	if e0 != nil {
		return e0
	}
	value, e1 := UnmarshalBytes(reader)
	if e1 != nil {
		return e1
	}
	lastIndex, e2 := UnmarshalInt32(reader)
	if e2 != nil {
		return e2
	}
	lastTerm, e3 := UnmarshalInt32(reader)
	if e3 != nil {
		return e3
	}
	r.CTerm = Term(term)
	r.Id = CandidateId(value)
	r.LastIndex = Index(lastIndex)
	r.LastTerm = Term(lastTerm)
	return nil
}

type AppendEntriesRsp struct {
	CTerm   Term
	Success bool
}

func (r *AppendEntriesRsp) Marshal(writer io.Writer) error {
	if e := MarshalInt32(writer, int32(r.CTerm)); e != nil {
		return e
	}
	return binary.Write(writer, binary.LittleEndian, r.Success)
}
func (r *AppendEntriesRsp) Unmarshal(reader io.Reader) error {
	term, e := UnmarshalInt32(reader)
	if e != nil {
		return e
	}
	r.CTerm = Term(term)
	return binary.Read(reader, binary.LittleEndian, &r.Success)
}

type AppendEntriesReq struct {
	CTerm       Term
	Id          CandidateId
	LastIndex   Index
	LastTerm    Term
	Entries     []Entry
	CommitIndex Index
}

func (r *AppendEntriesReq) Marshal(writer io.Writer) error {
	if e := MarshalInt32(writer, int32(r.CTerm)); e != nil {
		return e
	}
	if e := MarshalBytes(writer, []byte(r.Id)); e != nil {
		return e
	}
	if e := MarshalInt32(writer, int32(r.LastIndex)); e != nil {
		return e
	}
	if e := MarshalInt32(writer, int32(r.LastTerm)); e != nil {
		return e
	}
	if e := MarshalInt32(writer, int32(r.CommitIndex)); e != nil {
		return e
	}
	if e := MarshalInt32(writer, int32(len(r.Entries))); e != nil {
		return e
	}
	for _, v := range r.Entries {
		if e := v.Marshal(writer); e != nil {
			return e
		}
	}
	return nil
}
func (r *AppendEntriesReq) Unmarshal(reader io.Reader) error {
	term, e0 := UnmarshalInt32(reader)
	if e0 != nil {
		return e0
	}
	id, e1 := UnmarshalBytes(reader)
	if e1 != nil {
		return e1
	}
	lastIndex, e2 := UnmarshalInt32(reader)
	if e2 != nil {
		return e2
	}
	lastTerm, e3 := UnmarshalInt32(reader)
	if e3 != nil {
		return e3
	}
	commitIndex, e4 := UnmarshalInt32(reader)
	if e4 != nil {
		return e4
	}
	l, e5 := UnmarshalInt32(reader)
	if e5 != nil {
		return e5
	}
	r.CTerm = Term(term)
	r.Id = CandidateId(id)
	r.LastIndex = Index(lastIndex)
	r.LastTerm = Term(lastTerm)
	r.CommitIndex = Index(commitIndex)
	r.Entries = make([]Entry, l)
	for i := int32(0); i < l; i++ {
		var entry Entry
		e := entry.Unmarshal(reader)
		if e != nil {
			return e
		}
		r.Entries[i] = entry
	}
	return nil
}

func MarshalInt32(writer io.Writer, value int32) error {
	bs := binary.LittleEndian.AppendUint32(make([]byte, 0), uint32(value))
	_, e := writer.Write(bs)
	return e
}

func UnmarshalInt32(reader io.Reader) (int32, error) {
	var l int32
	e := binary.Read(reader, binary.LittleEndian, &l)
	return l, e
}

func MarshalBytes(writer io.Writer, bytes []byte) error {
	e := binary.Write(writer, binary.LittleEndian, int32(len(bytes)))
	if e != nil {
		return e
	}
	_, e0 := writer.Write(bytes)
	return e0
}

func UnmarshalBytes(reader io.Reader) ([]byte, error) {
	var l int32
	e := binary.Read(reader, binary.LittleEndian, &l)
	if e != nil || l < 0 {
		return nil, e
	}
	bs := make([]byte, l)
	_, e0 := reader.Read(bs)
	return bs, e0
}
