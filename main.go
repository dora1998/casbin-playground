package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

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

	fmt.Println("Casbin Playground - Enter 'sub,obj,act' format (Ctrl+D to exit)")
	fmt.Print("> ")

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			fmt.Print("> ")
			continue
		}

		parts := strings.Split(input, ",")
		if len(parts) != 3 {
			fmt.Println("Error: Please enter in 'sub,obj,act' format")
			fmt.Print("> ")
			continue
		}

		sub := strings.TrimSpace(parts[0])
		obj := strings.TrimSpace(parts[1])
		act := strings.TrimSpace(parts[2])

		ok, err := e.Enforce(sub, obj, act)
		if err != nil {
			fmt.Printf("Error: enforce: %s\n", err)
		} else {
			fmt.Printf("%s / %s / %s: %v\n", sub, obj, act, ok)
		}

		fmt.Print("> ")
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading input: %s", err)
	}

	fmt.Println("\nGoodbye!")
}
