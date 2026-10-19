package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"atad-project/internal/money"
	"atad-project/internal/store"
)

// newTransaction is the POST body. Amount is a pointer so a missing
// amount can be told apart from "0.00".
type newTransaction struct {
	Date        string        `json:"date"`
	Amount      *money.Amount `json:"amount"`
	Description string        `json:"description"`
	Category    string        `json:"category"`
}

func validDate(s string) bool {
	_, err := time.Parse(store.DateLayout, s)
	return err == nil
}

func registerTransactionRoutes(mux *http.ServeMux, st *store.Store) {
	mux.HandleFunc("POST /transactions", func(w http.ResponseWriter, r *http.Request) {
		var in newTransaction
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			writeError(w, http.StatusBadRequest, "validation", "invalid body: "+err.Error())
			return
		}
		if !validDate(in.Date) {
			writeError(w, http.StatusBadRequest, "validation", "date must be a real date in the form YYYY-MM-DD")
			return
		}
		if in.Amount == nil {
			writeError(w, http.StatusBadRequest, "validation", `amount is required, as a string like "-12.50"`)
			return
		}
		desc := strings.TrimSpace(in.Description)
		if desc == "" {
			writeError(w, http.StatusBadRequest, "validation", "description is required")
			return
		}
		cat := strings.TrimSpace(in.Category)
		if cat == "" {
			cat = "uncategorised" // replaced by the rules in the next step
		}

		t := st.Add(store.Transaction{Date: in.Date, Amount: *in.Amount, Description: desc, Category: cat})
		writeJSON(w, http.StatusCreated, t)
	})

	mux.HandleFunc("GET /transactions", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		f := store.Filter{
			From:     q.Get("from"),
			To:       q.Get("to"),
			Category: q.Get("category"),
			Query:    q.Get("q"),
		}
		if f.From != "" && !validDate(f.From) {
			writeError(w, http.StatusBadRequest, "validation", "from must be YYYY-MM-DD")
			return
		}
		if f.To != "" && !validDate(f.To) {
			writeError(w, http.StatusBadRequest, "validation", "to must be YYYY-MM-DD")
			return
		}
		if f.From != "" && f.To != "" && f.From > f.To {
			writeError(w, http.StatusBadRequest, "validation", "from must not be after to")
			return
		}
		writeJSON(w, http.StatusOK, st.List(f))
	})
}