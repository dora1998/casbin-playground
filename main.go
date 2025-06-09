package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/casbin/casbin/v2"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
)

type CheckResponse struct {
	OK bool `json:"ok"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

var enforcer *casbin.Enforcer
var watcherCh chan string

func main() {
	port := flag.Int("port", 8080, "Port number to listen on")
	flag.Parse()

	watcherCh = make(chan string)
	adapter := fileadapter.NewAdapter("policy.csv")
	watcher, err := NewWatcher(watcherCh)
	if err != nil {
		log.Fatalf("error: watcher: %s", err)
	}

	enforcer, err = NewEnforcer()
	if err != nil {
		log.Fatalf("error: enforcer: %s", err)
	}

	enforcer.SetAdapter(adapter)
	enforcer.SetWatcher(watcher)
	enforcer.LoadPolicy()

	http.HandleFunc("GET /check", checkHandler)
	http.HandleFunc("POST /update", updateHandler)

	addr := fmt.Sprintf(":%d", *port)
	fmt.Printf("Casbin HTTP API Server starting on %s\n", addr)
	fmt.Println("Usage: \n- GET /check?query=sub,obj,act\n- POST /update")
	log.Fatal(http.ListenAndServe(addr, nil))
}

func checkHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "query parameter is required"})
		return
	}

	parts := strings.Split(query, ",")
	if len(parts) != 3 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "query must be in 'sub,obj,act' format"})
		return
	}

	sub := strings.TrimSpace(parts[0])
	obj := strings.TrimSpace(parts[1])
	act := strings.TrimSpace(parts[2])

	ok, err := enforcer.Enforce(sub, obj, act)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: fmt.Sprintf("enforce error: %s", err)})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(CheckResponse{OK: ok})
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	watcherCh <- "update"

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
