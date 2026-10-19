package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Completate la build prin -ldflags (vezi Dockerfile)
var (
	version = "dev"
	commit  = "unknown"
	builtAt = "unknown"
)

var startTime = time.Now()

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{
			"status":         "ok",
			"uptime_seconds": int(time.Since(startTime).Seconds()),
		})
	})

	mux.HandleFunc("GET /version", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{
			"app":      "personal-finance-tracker",
			"version":  version,
			"commit":   commit,
			"built_at": builtAt,
		})
	})

	// {$} = doar "/", nu orice ruta
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, `<!doctype html>
<html><head><meta charset="utf-8"><title>Personal Finance Tracker</title></head>
<body>
<h1>Personal Finance Tracker</h1>
<p>Student: Florina Mare</p>
<p>Group: Group 2</p>
<p>ATAD 2026-2027 – Project 02</p>
</body></html>`)
	})

	mux.HandleFunc("POST /reset", func(w http.ResponseWriter, r *http.Request) {
		// Mai tarziu: golim store-ul si resetam ID-urile
		w.WriteHeader(http.StatusNoContent)
	})

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}