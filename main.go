package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
)

type CheckResponse struct {
	OK bool `json:"ok"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

var enforcer *casbin.Enforcer

func main() {
	// コマンドライン引数の解析
	port := flag.Int("port", 8080, "Port number to listen on")
	flag.Parse()

	// Casbinの初期化
	a := fileadapter.NewAdapter("policy.csv")

	m, err := model.NewModelFromString(`
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && r.obj == p.obj && r.act == p.act
`)
	if err != nil {
		log.Fatalf("error: model: %s", err)
	}

	enforcer, err = casbin.NewEnforcer(m, a)
	if err != nil {
		log.Fatalf("error: enforcer: %s", err)
	}

	// HTTPハンドラーの設定
	http.HandleFunc("/check", checkHandler)

	addr := fmt.Sprintf(":%d", *port)
	fmt.Printf("Casbin HTTP API Server starting on %s\n", addr)
	fmt.Println("Usage: GET /check?query=sub,obj,act")
	log.Fatal(http.ListenAndServe(addr, nil))
}

func checkHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Method not allowed"})
		return
	}

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
