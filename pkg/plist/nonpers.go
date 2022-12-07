package plist

import . "katamaran/pkg/data"

type NList struct {
	votedFor    *CandidateId
	currentTerm Term
	entries     []Entry
}

func MakeNList(id CandidateId) *NList {
	return &NList{
		votedFor:    nil,
		currentTerm: 0,
		entries:     make([]Entry, 0),
	}
}

func (p *NList) Flush() {
}
func (n *NList) SetVotedFor(votedFor *CandidateId) {
	n.votedFor = votedFor
}
func (n *NList) GetVotedFor() *CandidateId {
	return n.votedFor
}
func (n *NList) SetCurrentTerm(term Term) {
	n.currentTerm = term
}
func (n *NList) GetCurrentTerm() Term {
	return n.currentTerm
}
func (n *NList) Add(entry Entry) {
	n.entries = append(n.entries, entry)
}
func (n *NList) Truncate(index Index) error {
	n.entries = n.entries[:index]
	return nil
}
func (n *NList) Get(index Index) Entry {
	return n.entries[index]
}
func (n *NList) GetAllAfter(id Index) []Entry {
	return n.entries[id:]
}
func (n *NList) GetNextIndex() Index {
	return Index(len(n.entries))
}
func (n *NList) GetLastEntry() Entry {
	if len(n.entries) == 0 {
		return Entry{Value: nil, Index: -1, Term: 0}
	}
	return n.entries[len(n.entries)-1]
}
