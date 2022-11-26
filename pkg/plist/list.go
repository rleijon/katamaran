package plist

import (
	. "katamaran/pkg/data"
	"log"
	"os"

	"github.com/genjidb/genji"
	"github.com/genjidb/genji/document"
	"github.com/genjidb/genji/types"
)

type Data interface {
	SetVotedFor(*CandidateId)
	GetVotedFor() *CandidateId
	SetCurrentTerm(term Term)
	GetCurrentTerm() Term
	Add(Entry)
	Truncate(Index) error
	Get(Index) Entry
	GetAll() []Entry
	GetNextIndex() Index
	GetLastEntry() Entry
}

type PList struct {
	db          *genji.DB
	errorLogger *log.Logger
}

func MakePList(id CandidateId) *PList {
	db, err := genji.Open(string(id))
	if err != nil {
		log.Fatal(err)
	}
	err = db.Exec(`
		CREATE TABLE IF NOT EXISTS TRANSACTION_LOG(
			idx INT PRIMARY KEY,
			val BLOB,
			term INT
		)
	`)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Exec(`
		CREATE TABLE IF NOT EXISTS CURRENT_TERM(
			id INT PRIMARY KEY,
			term INT
		)
	`)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Exec(`
		CREATE TABLE IF NOT EXISTS VOTED_FOR(
			id INT PRIMARY KEY,
			votedFor TEXT
		)
	`)
	if err != nil {
		log.Fatal(err)
	}
	return &PList{
		db:          db,
		errorLogger: log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (p *PList) SetVotedFor(votedFor *CandidateId) {
	e0 := p.db.Exec("INSERT INTO VOTED_FOR(id, votedFor) VALUES (?, ?) ON CONFLICT REPLACE", 0, votedFor)
	if e0 != nil {
		p.errorLogger.Fatal(e0)
	}
}

func (p *PList) GetVotedFor() *CandidateId {
	res, e0 := p.db.Query("SELECT votedFor FROM VOTED_FOR")
	if e0 != nil {
		p.errorLogger.Fatal(e0)
	}
	isSet := false
	var votedFor CandidateId
	res.Iterate(func(d types.Document) error {
		v, e1 := d.GetByField("votedFor")
		if e1 != nil {
			p.errorLogger.Fatal(e1)
		}
		if types.IsNull(v) {
			p.SetVotedFor(nil)
		} else {
			isSet = true
			e2 := document.Scan(d, &votedFor)
			if e2 != nil {
				p.errorLogger.Fatal(e2)
			}
		}
		return nil

	})
	if isSet {
		return &votedFor
	}
	return nil
}

func (p *PList) SetCurrentTerm(term Term) {
	e0 := p.db.Exec("INSERT INTO CURRENT_TERM(id, term) VALUES (?, ?) ON CONFLICT REPLACE", 0, int(term))
	if e0 != nil {
		p.errorLogger.Fatal(e0)
	}
}

func (p *PList) GetCurrentTerm() Term {
	res, e0 := p.db.Query("SELECT term FROM CURRENT_TERM")
	if e0 != nil {
		p.errorLogger.Fatal(e0)
	}
	isSet := false
	var term Term
	res.Iterate(func(d types.Document) error {
		isSet = true
		e2 := document.Scan(d, &term)
		if e2 != nil {
			p.errorLogger.Fatal(e2)
		}
		return nil
	})
	if !isSet {
		p.SetCurrentTerm(0)
	}
	return term
}

func (p *PList) Add(entry Entry) {
	err := p.db.Exec("INSERT INTO TRANSACTION_LOG(idx, val, term) VALUES (?, ?, ?)", entry.Index, entry.Value, entry.Term)
	if err != nil {
		p.errorLogger.Fatal("Add", "err", err)
	}
}

func (p *PList) Truncate(index Index) error {
	return p.db.Exec("DELETE FROM TRANSACTION_LOG where idx > ?", index)
}

func (p *PList) Get(index Index) Entry {
	res, e0 := p.db.QueryDocument("SELECT idx, val, term as idx FROM TRANSACTION_LOG where idx = ?", index)
	if e0 != nil {
		p.errorLogger.Fatal("Get", "e0", e0)
	}
	var entry Entry
	e1 := document.Scan(res, &entry.Index, &entry.Value, &entry.Term)
	if e1 != nil {
		p.errorLogger.Fatal("Get", "e1", e1)
	}
	return entry
}

func (p *PList) GetAll() []Entry {
	res, e0 := p.db.Query("SELECT idx, val, term as idx FROM TRANSACTION_LOG")
	if e0 != nil {
		p.errorLogger.Fatal("GetAll", "e0", e0)
	}
	defer res.Close()
	entries := make([]Entry, 0)
	res.Iterate(func(d types.Document) error {
		var entry Entry
		e := document.Scan(d, &entry.Index, &entry.Value, &entry.Term)
		if e != nil {
			return e
		}
		entries = append(entries, entry)
		return nil
	})
	return entries
}

func (p *PList) getMaxIndex() (Index, bool, error) {
	res, e0 := p.db.QueryDocument("SELECT MAX(idx) as idx FROM TRANSACTION_LOG")
	if e0 != nil {
		return 0, false, e0
	}
	v, e1 := res.GetByField("idx")
	if e1 != nil {
		return 0, false, e1
	}
	if types.IsNull(v) {
		return 0, false, nil
	}
	var maxIdx int
	e2 := document.Scan(res, &maxIdx)
	if e1 != nil {
		return 0, false, e2
	}
	return Index(maxIdx), true, e0
}

func (p *PList) GetNextIndex() Index {
	maxIdx, exists, err := p.getMaxIndex()
	if err != nil {
		p.errorLogger.Fatal("GetNextIndex", err)
	} else if !exists {
		return maxIdx
	}
	return maxIdx + 1
}

func (p *PList) GetLastEntry() Entry {
	maxIdx, success, e0 := p.getMaxIndex()
	if e0 != nil {
		p.errorLogger.Fatal("GetLastEntry", "e0", e0)
	}
	if !success {
		return Entry{Value: nil, Index: -1, Term: 0}
	}
	lastEntry, e2 := p.db.QueryDocument("SELECT idx, val, term as idx FROM TRANSACTION_LOG where idx = ?", maxIdx)
	if e2 != nil {
		p.errorLogger.Fatal("GetLastEntry", "e2", e2)
	}
	var entry Entry
	document.Scan(lastEntry, &entry.Index, &entry.Value, &entry.Term)
	return entry
}
