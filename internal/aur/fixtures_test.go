package aur_test

import (
	"io/ioutil"
	"path/filepath"
)

var fixtures map[string]string

func init() {
	fixtures = map[string]string{}

	files, err := filepath.Glob("fixtures/*")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		contents, err := ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}

		fixtures[filepath.Base(file)] = string(contents)
	}
}
