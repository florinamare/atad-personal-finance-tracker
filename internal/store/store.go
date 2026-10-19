// Package store keeps transactions in memory. It is safe for concurrent use.
package store

import (
	"sort"
	"strings"
	"sync"

	"atad-project/internal/money"
)

// DateLayout is the only date format the API accepts: YYYY-MM-DD.
const DateLayout = "2006-01-02"

type Transaction struct {
	ID          int          `json:"id"`
	Date        string       `json:"date"`
	Amount      money.Amount `json:"amount"`
	Description string       `json:"description"`
	Category    string       `json:"category"`
}

// Filter selects transactions. Empty fields mean "no restriction".
// From and To are inclusive and already validated as YYYY-MM-DD.
type Filter struct {
	From, To string
	Category string
	Query    string
}

type Store struct {
	mu     sync.Mutex
	nextID int
	txs    []Transaction
}

func New() *Store {
	return &Store{nextID: 1}
}

// Add assigns the next ID and stores the transaction.
func (s *Store) Add(t Transaction) Transaction {
	s.mu.Lock()
	defer s.mu.Unlock()
	t.ID = s.nextID
	s.nextID++
	s.txs = append(s.txs, t)
	return t
}

// List returns the matching transactions, newest first.
func (s *Store) List(f Filter) []Transaction {
	s.mu.Lock()
	defer s.mu.Unlock()

	q := strings.ToLower(f.Query)
	out := []Transaction{}
	for _, t := range s.txs {
		// YYYY-MM-DD strings compare correctly as text.
		if f.From != "" && t.Date < f.From {
			continue
		}
		if f.To != "" && t.Date > f.To {
			continue
		}
		if f.Category != "" && t.Category != f.Category {
			continue
		}
		if q != "" && !strings.Contains(strings.ToLower(t.Description), q) {
			continue
		}
		out = append(out, t)
	}

	sort.Slice(out, func(i, j int) bool {
		if out[i].Date != out[j].Date {
			return out[i].Date > out[j].Date
		}
		return out[i].ID > out[j].ID // same day: added later comes first
	})
	return out
}

// Reset discards everything and restarts IDs from 1.
func (s *Store) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.txs = nil
	s.nextID = 1
}