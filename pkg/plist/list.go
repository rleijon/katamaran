package plist

import (
	. "katamaran/pkg/data"
)

type Data interface {
	SetVotedFor(*CandidateId)
	GetVotedFor() *CandidateId
	SetCurrentTerm(term Term)
	GetCurrentTerm() Term
	Add(Entry)
	Truncate(Index) error
	Get(Index) Entry
	GetAllAfter(Index) []Entry
	GetNextIndex() Index
	GetLastEntry() Entry
	Flush()
}
