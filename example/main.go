package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gumpen/jsonpath"
)

func main() {
	b, err := os.ReadFile("example.json")
	if err != nil {
		log.Fatal(err)
	}

	var d interface{}
	err = json.Unmarshal(b, &d)
	if err != nil {
		log.Fatal(err)
	}

	q := "$.Author.ID"
	p := jsonpath.NewPath(q)
	err = p.Parse()
	if err != nil {
		log.Fatal(err)
	}

	out, err := p.Execute(d)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%#v\n", out)
}
