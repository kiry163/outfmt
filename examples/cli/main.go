package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kiry163/outfmt"
)

type User struct {
	ID     int    `json:"id" yaml:"id" outfmt:"ID"`
	Name   string `json:"name" yaml:"name" outfmt:"Name"`
	Email  string `json:"email" yaml:"email" outfmt:"Email"`
	Status string `json:"status" yaml:"status" outfmt:"Status"`
}

func main() {
	var output string

	flag.StringVar(&output, "output", string(outfmt.Table), "output format: table|yaml|json")
	flag.StringVar(&output, "o", string(outfmt.Table), "output format shorthand")
	flag.Parse()

	format := outfmt.Format(output)
	if err := format.Validate(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	users := []User{
		{ID: 1, Name: "alice", Email: "alice@example.com", Status: "active"},
		{ID: 2, Name: "bob", Email: "bob@example.com", Status: "inactive"},
	}

	if err := outfmt.Render(os.Stdout, users, format); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
