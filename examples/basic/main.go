package main

import (
	"fmt"
	"os"

	"github.com/kiry163/outfmt"
)

type User struct {
	ID    int    `json:"id" yaml:"id" outfmt:"ID"`
	Name  string `json:"name" yaml:"name" outfmt:"Name"`
	Email string `json:"email" yaml:"email" outfmt:"Email"`
}

func main() {
	users := []User{
		{ID: 1, Name: "alice", Email: "alice@example.com"},
		{ID: 2, Name: "bob", Email: "bob@example.com"},
	}

	render("table", users, outfmt.Table)
	fmt.Println()
	render("yaml", users, outfmt.YAML)
	fmt.Println()
	render("json", users, outfmt.JSON)
}

func render(title string, data any, format outfmt.Format) {
	fmt.Printf("== %s ==\n", title)
	if err := outfmt.Render(os.Stdout, data, format); err != nil {
		panic(err)
	}
}
