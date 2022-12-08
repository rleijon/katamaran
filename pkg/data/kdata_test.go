package data

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"testing"
)

func TestEntryIndex_Marshal(t *testing.T) {
	type fields struct {
		Id  Index
		Pos int32
	}
	tests := []struct {
		name       string
		fields     fields
		wantWriter string
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &EntryIndex{
				Id:  tt.fields.Id,
				Pos: tt.fields.Pos,
			}
			writer := &bytes.Buffer{}
			if err := r.Marshal(writer); (err != nil) != tt.wantErr {
				t.Errorf("EntryIndex.Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWriter := writer.String(); gotWriter != tt.wantWriter {
				t.Errorf("EntryIndex.Marshal() = %v, want %v", gotWriter, tt.wantWriter)
			}
		})
	}
}

func TestEntryIndex_Unmarshal(t *testing.T) {
	type fields struct {
		Id  Index
		Pos int32
	}
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &EntryIndex{
				Id:  tt.fields.Id,
				Pos: tt.fields.Pos,
			}
			if err := r.Unmarshal(tt.args.reader); (err != nil) != tt.wantErr {
				t.Errorf("EntryIndex.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEntry_MarshalledSize(t *testing.T) {
	type fields struct {
		Value []byte
		Index Index
		Term  Term
	}
	tests := []struct {
		name   string
		fields fields
		want   int32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Entry{
				Value: tt.fields.Value,
				Index: tt.fields.Index,
				Term:  tt.fields.Term,
			}
			if got := r.MarshalledSize(); got != tt.want {
				t.Errorf("Entry.MarshalledSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEntry_Marshal(t *testing.T) {
	type fields struct {
		Value []byte
		Index Index
		Term  Term
	}
	tests := []struct {
		name       string
		fields     fields
		wantWriter string
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Entry{
				Value: tt.fields.Value,
				Index: tt.fields.Index,
				Term:  tt.fields.Term,
			}
			writer := &bytes.Buffer{}
			if err := r.Marshal(writer); (err != nil) != tt.wantErr {
				t.Errorf("Entry.Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWriter := writer.String(); gotWriter != tt.wantWriter {
				t.Errorf("Entry.Marshal() = %v, want %v", gotWriter, tt.wantWriter)
			}
		})
	}
}

func TestEntry_Unmarshal(t *testing.T) {
	type fields struct {
		Value []byte
		Index Index
		Term  Term
	}
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Entry{
				Value: tt.fields.Value,
				Index: tt.fields.Index,
				Term:  tt.fields.Term,
			}
			if err := r.Unmarshal(tt.args.reader); (err != nil) != tt.wantErr {
				t.Errorf("Entry.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRequestVotesRsp_Marshal(t *testing.T) {
	type fields struct {
		VoteGranted bool
	}
	tests := []struct {
		name       string
		fields     fields
		wantWriter string
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RequestVotesRsp{
				VoteGranted: tt.fields.VoteGranted,
			}
			writer := &bytes.Buffer{}
			if err := r.Marshal(writer); (err != nil) != tt.wantErr {
				t.Errorf("RequestVotesRsp.Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWriter := writer.String(); gotWriter != tt.wantWriter {
				t.Errorf("RequestVotesRsp.Marshal() = %v, want %v", gotWriter, tt.wantWriter)
			}
		})
	}
}

func TestRequestVotesRsp_Unmarshal(t *testing.T) {
	type fields struct {
		VoteGranted bool
	}
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RequestVotesRsp{
				VoteGranted: tt.fields.VoteGranted,
			}
			if err := r.Unmarshal(tt.args.reader); (err != nil) != tt.wantErr {
				t.Errorf("RequestVotesRsp.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTickReq_Marshal(t *testing.T) {
	tests := []struct {
		name       string
		r          *TickReq
		wantWriter string
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &TickReq{}
			writer := &bytes.Buffer{}
			if err := r.Marshal(writer); (err != nil) != tt.wantErr {
				t.Errorf("TickReq.Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWriter := writer.String(); gotWriter != tt.wantWriter {
				t.Errorf("TickReq.Marshal() = %v, want %v", gotWriter, tt.wantWriter)
			}
		})
	}
}

func TestTickReq_Unmarshal(t *testing.T) {
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name    string
		r       *TickReq
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &TickReq{}
			if err := r.Unmarshal(tt.args.reader); (err != nil) != tt.wantErr {
				t.Errorf("TickReq.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAddEntryReq_Marshal(t *testing.T) {
	type fields struct {
		Value []byte
	}
	tests := []struct {
		name       string
		fields     fields
		wantWriter string
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &AddEntryReq{
				Value: tt.fields.Value,
			}
			writer := &bytes.Buffer{}
			if err := r.Marshal(writer); (err != nil) != tt.wantErr {
				t.Errorf("AddEntryReq.Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWriter := writer.String(); gotWriter != tt.wantWriter {
				t.Errorf("AddEntryReq.Marshal() = %v, want %v", gotWriter, tt.wantWriter)
			}
		})
	}
}

func TestAddEntryReq_Unmarshal(t *testing.T) {
	type fields struct {
		Value []byte
	}
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &AddEntryReq{
				Value: tt.fields.Value,
			}
			if err := r.Unmarshal(tt.args.reader); (err != nil) != tt.wantErr {
				t.Errorf("AddEntryReq.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAddEntryRsp_Marshal(t *testing.T) {
	type fields struct {
		Success bool
	}
	tests := []struct {
		name       string
		fields     fields
		wantWriter string
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &AddEntryRsp{
				Success: tt.fields.Success,
			}
			writer := &bytes.Buffer{}
			if err := r.Marshal(writer); (err != nil) != tt.wantErr {
				t.Errorf("AddEntryRsp.Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWriter := writer.String(); gotWriter != tt.wantWriter {
				t.Errorf("AddEntryRsp.Marshal() = %v, want %v", gotWriter, tt.wantWriter)
			}
		})
	}
}

func TestAddEntryRsp_Unmarshal(t *testing.T) {
	type fields struct {
		Success bool
	}
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &AddEntryRsp{
				Success: tt.fields.Success,
			}
			if err := r.Unmarshal(tt.args.reader); (err != nil) != tt.wantErr {
				t.Errorf("AddEntryRsp.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRequestVotesReq_Marshal(t *testing.T) {
	type fields struct {
		CTerm     Term
		Id        CandidateId
		LastIndex Index
		LastTerm  Term
	}
	tests := []struct {
		name       string
		fields     fields
		wantWriter string
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RequestVotesReq{
				CTerm:     tt.fields.CTerm,
				Id:        tt.fields.Id,
				LastIndex: tt.fields.LastIndex,
				LastTerm:  tt.fields.LastTerm,
			}
			writer := &bytes.Buffer{}
			if err := r.Marshal(writer); (err != nil) != tt.wantErr {
				t.Errorf("RequestVotesReq.Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWriter := writer.String(); gotWriter != tt.wantWriter {
				t.Errorf("RequestVotesReq.Marshal() = %v, want %v", gotWriter, tt.wantWriter)
			}
		})
	}
}

func TestRequestVotesReq_Unmarshal(t *testing.T) {
	type fields struct {
		CTerm     Term
		Id        CandidateId
		LastIndex Index
		LastTerm  Term
	}
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RequestVotesReq{
				CTerm:     tt.fields.CTerm,
				Id:        tt.fields.Id,
				LastIndex: tt.fields.LastIndex,
				LastTerm:  tt.fields.LastTerm,
			}
			if err := r.Unmarshal(tt.args.reader); (err != nil) != tt.wantErr {
				t.Errorf("RequestVotesReq.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAppendEntriesRsp_Marshal(t *testing.T) {
	type fields struct {
		CTerm   Term
		Success bool
	}
	tests := []struct {
		name       string
		fields     fields
		wantWriter string
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &AppendEntriesRsp{
				CTerm:   tt.fields.CTerm,
				Success: tt.fields.Success,
			}
			writer := &bytes.Buffer{}
			if err := r.Marshal(writer); (err != nil) != tt.wantErr {
				t.Errorf("AppendEntriesRsp.Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWriter := writer.String(); gotWriter != tt.wantWriter {
				t.Errorf("AppendEntriesRsp.Marshal() = %v, want %v", gotWriter, tt.wantWriter)
			}
		})
	}
}

func TestAppendEntriesRsp_Unmarshal(t *testing.T) {
	type fields struct {
		CTerm   Term
		Success bool
	}
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &AppendEntriesRsp{
				CTerm:   tt.fields.CTerm,
				Success: tt.fields.Success,
			}
			if err := r.Unmarshal(tt.args.reader); (err != nil) != tt.wantErr {
				t.Errorf("AppendEntriesRsp.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAppendEntriesReq_Marshal(t *testing.T) {
	type fields struct {
		CTerm       Term
		Id          CandidateId
		LastIndex   Index
		LastTerm    Term
		Entries     []Entry
		CommitIndex Index
	}
	tests := []struct {
		name       string
		fields     fields
		wantWriter string
		wantErr    bool
	}{
		{name: "BasicTest", fields: fields{
			CTerm: 1, Id: "EXAMPLE", LastIndex: 2,
			LastTerm: 3, Entries: nil, CommitIndex: 1},
			wantWriter: "",
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &AppendEntriesReq{
				CTerm:       tt.fields.CTerm,
				Id:          tt.fields.Id,
				LastIndex:   tt.fields.LastIndex,
				LastTerm:    tt.fields.LastTerm,
				Entries:     tt.fields.Entries,
				CommitIndex: tt.fields.CommitIndex,
			}
			writer := &bytes.Buffer{}
			if err := r.Marshal(writer); (err != nil) != tt.wantErr {
				t.Errorf("AppendEntriesReq.Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			var rsp AppendEntriesReq
			if err := rsp.Unmarshal(writer); (err != nil) != tt.wantErr {
				t.Errorf("AppendEntriesReq.Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println("What?!", rsp)
			t.Log(rsp)
		})
	}
}

func TestAppendEntriesReq_Unmarshal(t *testing.T) {
	type fields struct {
		CTerm       Term
		Id          CandidateId
		LastIndex   Index
		LastTerm    Term
		Entries     []Entry
		CommitIndex Index
	}
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &AppendEntriesReq{
				CTerm:       tt.fields.CTerm,
				Id:          tt.fields.Id,
				LastIndex:   tt.fields.LastIndex,
				LastTerm:    tt.fields.LastTerm,
				Entries:     tt.fields.Entries,
				CommitIndex: tt.fields.CommitIndex,
			}
			if err := r.Unmarshal(tt.args.reader); (err != nil) != tt.wantErr {
				t.Errorf("AppendEntriesReq.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMarshalInt32(t *testing.T) {
	type args struct {
		value int32
	}
	tests := []struct {
		name       string
		args       args
		wantWriter string
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &bytes.Buffer{}
			if err := MarshalInt32(writer, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("MarshalInt32() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWriter := writer.String(); gotWriter != tt.wantWriter {
				t.Errorf("MarshalInt32() = %v, want %v", gotWriter, tt.wantWriter)
			}
		})
	}
}

func TestUnmarshalInt32(t *testing.T) {
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    int32
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnmarshalInt32(tt.args.reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalInt32() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UnmarshalInt32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMarshalBytes(t *testing.T) {
	type args struct {
		bytes []byte
	}
	tests := []struct {
		name       string
		args       args
		wantWriter string
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &bytes.Buffer{}
			if err := MarshalBytes(writer, tt.args.bytes); (err != nil) != tt.wantErr {
				t.Errorf("MarshalBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWriter := writer.String(); gotWriter != tt.wantWriter {
				t.Errorf("MarshalBytes() = %v, want %v", gotWriter, tt.wantWriter)
			}
		})
	}
}

func TestUnmarshalBytes(t *testing.T) {
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnmarshalBytes(tt.args.reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnmarshalBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
