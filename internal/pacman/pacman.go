package pacman

import (
	"bytes"
	"regexp"
	"strings"
)

type Executor interface {
	Execute(command string, args ...string) (*bytes.Buffer, error)
}

func New(executor Executor) (Pacman, error) {
	stdout, err := executor.Execute("pacman", "-Si")
	if err != nil {
		return Pacman{}, err
	}

	availablePkgs := parsePackagesFromStdout(stdout)

	return Pacman{
		availablePkgs: availablePkgs,
	}, nil
}

type Pacman struct {
	availablePkgs []string
}

func (p Pacman) ListAvailable() []string {
	return p.availablePkgs
}

func parsePackagesFromStdout(stdout *bytes.Buffer) []string {
	var availablePkgs []string

	for {
		line, err := stdout.ReadString('\n')

		foundPkgs := parsePackagesFromLine(line)
		availablePkgs = append(availablePkgs, foundPkgs...)

		if err != nil {
			break
		}
	}

	return availablePkgs
}

func parsePackagesFromLine(line string) []string {
	key, value, ok := matchProvidesOrName(line)
	if !ok {
		return nil
	}

	switch key {
	case "Provides":
		if value == "None" {
			return nil
		}
		return splitProvides(value)

	case "Name":
		return []string{value}
	}

	panic("matched against unhandled key: " + key)
}

func matchProvidesOrName(line string) (string, string, bool) {
	re := regexp.MustCompile("^(Provides|Name)\\s*:\\s*(.*)")
	match := re.FindStringSubmatch(line)
	if match == nil {
		return "", "", false
	}

	return match[1], match[2], true
}

func splitProvides(provides string) []string {
	var pkgs []string

	constraints := splitByWhitespace(provides)

	for _, constraint := range constraints {
		pkg := stripConstraint(constraint)
		pkgs = append(pkgs, pkg)
	}

	return pkgs
}

func stripConstraint(constraint string) string {
	return strings.Split(constraint, "=")[0]
}

func splitByWhitespace(str string) []string {
	reWhitespace := regexp.MustCompile("\\s+")
	return reWhitespace.Split(str, -1)
}
