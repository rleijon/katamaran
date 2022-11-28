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
	Marshal(io.Writer)
	Unmarshal(io.Reader)
}

type Entry struct {
	Value []byte
	Index Index
	Term  Term
}

type RequestVotesRsp struct {
	VoteGranted bool
}

func (r *RequestVotesRsp) Marshal(writer io.Writer) {
	binary.Write(writer, binary.LittleEndian, r.VoteGranted)
}

func (r *RequestVotesRsp) Unmarshal(reader io.Reader) {
	binary.Read(reader, binary.LittleEndian, &r.VoteGranted)
}

type TickReq struct {
}

func (r *TickReq) Marshal(writer io.Writer) {
}

func (r *TickReq) Unmarshal(reader io.Reader) {
}

type AddEntryReq struct {
	Value []byte
}

func (r *AddEntryReq) Marshal(writer io.Writer) {
	MarshalBytes(writer, r.Value)
}
func (r *AddEntryReq) Unmarshal(reader io.Reader) {
	r.Value = UnmarshalBytes(reader)
}

type RequestVotesReq struct {
	CTerm     Term
	Id        CandidateId
	LastIndex Index
	LastTerm  Term
}

func (r *RequestVotesReq) Marshal(writer io.Writer) {
	MarshalInt32(writer, int32(r.CTerm))
	MarshalBytes(writer, []byte(r.Id))
	MarshalInt32(writer, int32(r.LastIndex))
	MarshalInt32(writer, int32(r.LastTerm))
}
func (r *RequestVotesReq) Unmarshal(reader io.Reader) {
	r.CTerm = Term(UnmarshalInt32(reader))
	r.Id = CandidateId(UnmarshalBytes(reader))
	r.LastIndex = Index(UnmarshalInt32(reader))
	r.LastTerm = Term(UnmarshalInt32(reader))
}

type AppendEntriesRsp struct {
	CTerm   Term
	Success bool
}

func (r *AppendEntriesRsp) Marshal(writer io.Writer) {
	MarshalInt32(writer, int32(r.CTerm))
	binary.Write(writer, binary.LittleEndian, r.Success)
}
func (r *AppendEntriesRsp) Unmarshal(reader io.Reader) {
	r.CTerm = Term(UnmarshalInt32(reader))
	binary.Read(reader, binary.LittleEndian, &r.Success)
}

type AppendEntriesReq struct {
	CTerm       Term
	Id          CandidateId
	LastIndex   Index
	LastTerm    Term
	Entries     []Entry
	CommitIndex Index
}

func (r *AppendEntriesReq) Marshal(writer io.Writer) {
	MarshalInt32(writer, int32(r.CTerm))
	MarshalBytes(writer, []byte(r.Id))
	MarshalInt32(writer, int32(r.LastIndex))
	MarshalInt32(writer, int32(r.LastTerm))
	MarshalInt32(writer, int32(r.CommitIndex))
	MarshalInt32(writer, int32(len(r.Entries)))
	for _, v := range r.Entries {
		MarshalInt32(writer, int32(v.Index))
		MarshalInt32(writer, int32(v.Term))
		MarshalBytes(writer, v.Value)
	}
}
func (r *AppendEntriesReq) Unmarshal(reader io.Reader) {
	r.CTerm = Term(UnmarshalInt32(reader))
	r.Id = CandidateId(UnmarshalBytes(reader))
	r.LastIndex = Index(UnmarshalInt32(reader))
	r.LastTerm = Term(UnmarshalInt32(reader))
	r.CommitIndex = Index(UnmarshalInt32(reader))
	l := UnmarshalInt32(reader)
	r.Entries = make([]Entry, l)
	for i := int32(0); i < l; i++ {
		var entry Entry
		entry.Index = Index(UnmarshalInt32(reader))
		entry.Term = Term(UnmarshalInt32(reader))
		entry.Value = UnmarshalBytes(reader)
		r.Entries[i] = entry
	}
}

func MarshalInt32(writer io.Writer, value int32) {
	bs := binary.LittleEndian.AppendUint32(make([]byte, 0), uint32(value))
	writer.Write(bs)
}

func UnmarshalInt32(reader io.Reader) int32 {
	var l int32
	binary.Read(reader, binary.LittleEndian, &l)
	return l
}

func MarshalBytes(writer io.Writer, bytes []byte) {
	binary.Write(writer, binary.LittleEndian, int32(len(bytes)))
	writer.Write(bytes)
}

func UnmarshalBytes(reader io.Reader) []byte {
	var l int32
	binary.Read(reader, binary.LittleEndian, &l)
	bs := make([]byte, l)
	reader.Read(bs)
	return bs
}
