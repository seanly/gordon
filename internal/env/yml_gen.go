// Code generated by main.go DO NOT EDIT.

package env

import (
	yaml "gopkg.in/yaml.v3"
	"io"
)

type YAML struct {
	GithubHost   *GithubHost   `yaml:"githubHost,omitempty"`
	GithubUser   *GithubUser   `yaml:"githubUser,omitempty"`
	Architecture *Architecture `yaml:"architecture,omitempty"`
	OS           *OS           `yaml:"os,omitempty"`
	Cache        *Cache        `yaml:"cache,omitempty"`
	Bin          *Bin          `yaml:"bin,omitempty"`
	Man          *Man          `yaml:"man,omitempty"`
	Hooks        *Hooks        `yaml:"hooks,omitempty"`
}

func saveYAML(w io.Writer, yml *YAML) error {
	return yaml.NewEncoder(w).Encode(yml)
}

var EmptyYAMLReader io.Reader = nil

func loadYAML(r io.Reader) (yml YAML, err error) {
	if r == EmptyYAMLReader {
		return
	}
	err = yaml.NewDecoder(r).Decode(&yml)
	return
}
