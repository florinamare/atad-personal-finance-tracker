package store

import (
	"testing"

	"atad-project/internal/money"
)

func add(s *Store, date, amount, desc, cat string) Transaction {
	a, err := money.Parse(amount)
	if err != nil {
		panic(err)
	}
	return s.Add(Transaction{Date: date, Amount: a, Description: desc, Category: cat})
}

func TestIDsStartAtOne(t *testing.T) {
	s := New()
	if got := add(s, "2026-10-01", "-1.00", "a", "x").ID; got != 1 {
		t.Errorf("first ID = %d, want 1", got)
	}
	if got := add(s, "2026-10-01", "-1.00", "b", "x").ID; got != 2 {
		t.Errorf("second ID = %d, want 2", got)
	}
}

func TestListNewestFirst(t *testing.T) {
	s := New()
	add(s, "2026-10-01", "-1.00", "old", "x")
	add(s, "2026-10-05", "-1.00", "new", "x")
	add(s, "2026-10-05", "-1.00", "newer same day", "x")

	got := s.List(Filter{})
	want := []string{"newer same day", "new", "old"}
	for i, w := range want {
		if got[i].Description != w {
			t.Errorf("position %d = %q, want %q", i, got[i].Description, w)
		}
	}
}

func TestListFilters(t *testing.T) {
	s := New()
	add(s, "2026-09-30", "-5.00", "Coffee Shop", "food")
	add(s, "2026-10-01", "-9.00", "Lidl groceries", "food")
	add(s, "2026-10-15", "1000.00", "Salary", "income")
	add(s, "2026-11-01", "-3.00", "coffee beans", "food")

	cases := []struct {
		name string
		f    Filter
		want int
	}{
		{"no filter", Filter{}, 4},
		{"from inclusive", Filter{From: "2026-10-01"}, 3},
		{"to inclusive", Filter{To: "2026-10-01"}, 2},
		{"range", Filter{From: "2026-10-01", To: "2026-10-31"}, 2},
		{"category", Filter{Category: "food"}, 3},
		{"text, case-insensitive", Filter{Query: "COFFEE"}, 2},
		{"combined", Filter{Category: "food", Query: "coffee", From: "2026-10-01"}, 1},
	}
	for _, c := range cases {
		if got := len(s.List(c.f)); got != c.want {
			t.Errorf("%s: got %d results, want %d", c.name, got, c.want)
		}
	}
}

func TestResetRestartsIDs(t *testing.T) {
	s := New()
	add(s, "2026-10-01", "-1.00", "a", "x")
	s.Reset()
	if n := len(s.List(Filter{})); n != 0 {
		t.Errorf("after reset: %d transactions, want 0", n)
	}
	if got := add(s, "2026-10-01", "-1.00", "b", "x").ID; got != 1 {
		t.Errorf("ID after reset = %d, want 1", got)
	}
}