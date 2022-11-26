package data

type CandidateId string
type Term int
type Index int
type State int

type Entry struct {
	Value []byte
	Index Index
	Term  Term
}
