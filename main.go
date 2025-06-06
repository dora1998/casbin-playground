package main

import (
	"fmt"
	"log"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
)

func main() {
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

	e, err := casbin.NewEnforcer(m, a)
	if err != nil {
		log.Fatalf("error: enforcer: %s", err)
	}

	ok, err := e.Enforce("alice", "data1", "read")
	if err != nil {
		log.Fatalf("error: enforce: %s", err)
	}

	fmt.Printf("alice / data1 / read: %v\n", ok)

	ok, err = e.Enforce("bob", "data1", "read")
	if err != nil {
		log.Fatalf("error: enforce: %s", err)
	}

	fmt.Printf("bob / data1 / read: %v\n", ok)
}
