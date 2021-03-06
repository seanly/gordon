// +build generate

package main

import (
	"log"

	"github.com/kyoh86/appenv/gen"
	"github.com/kyoh86/gordon/internal/env"
)

//go:generate go run -tags generate ./main.go

func main() {
	g := &gen.Generator{}

	if err := g.Do(
		"github.com/kyoh86/gordon/internal/env",
		"../",
		gen.Prop(new(env.GithubHost), gen.YAML(), gen.Envar()),
		gen.Prop(new(env.GithubUser), gen.YAML(), gen.Envar()),
		gen.Prop(new(env.Architecture), gen.YAML(), gen.Envar()),
		gen.Prop(new(env.OS), gen.YAML(), gen.Envar()),
		gen.Prop(new(env.Cache), gen.YAML(), gen.Envar()),
		gen.Prop(new(env.Bin), gen.YAML(), gen.Envar()),
		gen.Prop(new(env.Man), gen.YAML(), gen.Envar()),
		gen.Prop(new(env.Hooks), gen.YAML(), gen.Envar()),
	); err != nil {
		log.Fatalln(err)
	}
}
