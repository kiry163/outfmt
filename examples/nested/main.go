package main

import (
	"fmt"
	"os"

	"github.com/kiry163/outfmt"
)

type Profile struct {
	City   string `json:"city" yaml:"city" outfmt:"City"`
	Active bool   `json:"active" yaml:"active" outfmt:"Active"`
}

type User struct {
	ID      int            `json:"id" yaml:"id" outfmt:"ID"`
	Name    string         `json:"name" yaml:"name" outfmt:"Name"`
	Profile *Profile       `json:"profile" yaml:"profile" outfmt:"Profile"`
	Meta    map[string]any `json:"meta" yaml:"meta" outfmt:"Meta"`
	Tags    []string       `json:"tags" yaml:"tags" outfmt:"Tags"`
}

func main() {
	users := []User{
		{
			ID:      1,
			Name:    "alice",
			Profile: &Profile{City: "shanghai", Active: true},
			Meta: map[string]any{
				"region": "cn",
				"zone":   "east",
			},
			Tags: []string{"dev", "ops"},
		},
		{
			ID:   2,
			Name: "bob",
			Meta: map[string]any{
				"region": "us",
			},
		},
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
