package node

import (
	. "katamaran/pkg/data"
	"katamaran/pkg/plist"
	"testing"
	"time"
)

func TestNode_AddEntry(t *testing.T) {
	type fields struct {
		id            CandidateId
		leaderId      *CandidateId
		state         State
		list          plist.Data
		commitIndex   Index
		lastApplied   Index
		nextIndex     map[CandidateId]Index
		matchIndex    map[CandidateId]Index
		nextHeartBeat time.Time
		lastSeen      time.Time
		sender        Sender
	}
	type args struct {
		value []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{name: "Test1", },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &Node{
				id:            tt.fields.id,
				leaderId:      tt.fields.leaderId,
				state:         tt.fields.state,
				list:          tt.fields.list,
				commitIndex:   tt.fields.commitIndex,
				lastApplied:   tt.fields.lastApplied,
				nextIndex:     tt.fields.nextIndex,
				matchIndex:    tt.fields.matchIndex,
				nextHeartBeat: tt.fields.nextHeartBeat,
				lastSeen:      tt.fields.lastSeen,
				sender:        tt.fields.sender,
			}
			if got := n.AddEntry(tt.args.value); got != tt.want {
				t.Errorf("Node.AddEntry() = %v, want %v", got, tt.want)
			}
		})
	}
}
