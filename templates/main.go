package main

import (
	"bytes"
	"fmt"
	"text/template"
)

type persons struct {
	Name1 string
	Name2 string
}

func main() {
	t := template.New("gotcha")
	t, _ = t.Parse("hello {{.Name1}} and {{.Name2}}!")

	// use struct to fill template
	p := persons{Name1: "Peter", Name2: "Mary"}
	var d bytes.Buffer
	t.Execute(&d, p)
	fmt.Printf("Template result: %s\n", d.String())

	// use map to fill template
	d = bytes.Buffer{}
	t.Execute(&d, map[string]interface{}{
		"Name1": "Bob", "Name2": "Alice"})
	fmt.Printf("Template result: %s\n", d.String())
}
