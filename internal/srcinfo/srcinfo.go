package srcinfo

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"

	"github.com/ooesili/aurgo/internal/cache"
)

func New(targetArch string) SrcInfo {
	return SrcInfo{
		targetArch: targetArch,
	}
}

type SrcInfo struct {
	targetArch string
}

func (s SrcInfo) Parse(input []byte) (cache.Package, error) {
	fields := parseFields(input)

	seenPkgnameField := false
	foundCompatibleArch := false
	pkg := cache.Package{}

	for _, field := range fields {
		switch field.key {
		case "depends", "depends_" + s.targetArch:
			pkg.Depends = append(pkg.Depends, stripVersion(field.value))

		case "checkdepends", "checkdepends_" + s.targetArch:
			pkg.Checkdepends = append(pkg.Checkdepends, stripVersion(field.value))

		case "makedepends", "makedepends_" + s.targetArch:
			pkg.Makedepends = append(pkg.Makedepends, stripVersion(field.value))

		case "arch":
			if field.value == s.targetArch || field.value == "any" {
				foundCompatibleArch = true
			}

		case "pkgname":
			if seenPkgnameField {
				return cache.Package{}, errors.New("cannot handle split packages")
			}
			seenPkgnameField = true
		}
	}

	if !foundCompatibleArch {
		return cache.Package{}, fmt.Errorf("package does not support %s", s.targetArch)
	}

	return pkg, nil
}

func parseFields(input []byte) []field {
	buffer := bytes.NewBuffer(input)
	var fields []field

	for {
		line, err := readLine(buffer)

		field, ok := findField(line)
		if ok {
			fields = append(fields, field)
		}

		if err != nil {
			break
		}
	}

	return fields
}

func readLine(buffer *bytes.Buffer) ([]byte, error) {
	line, err := buffer.ReadBytes('\n')
	line = bytes.TrimRight(line, "\n")
	return line, err
}

func findField(line []byte) (field, bool) {
	re := regexp.MustCompile("^\\s*(\\S+) = (.+)$")
	match := re.FindSubmatch(line)
	if match == nil {
		return field{}, false
	}

	key := string(match[1])
	value := string(match[2])
	return field{key: key, value: value}, true
}

func stripVersion(dependency string) string {
	re := regexp.MustCompile("^([^=<>]*)")
	match := re.FindStringSubmatch(dependency)
	if match == nil {
		return dependency
	}

	return match[1]
}

type field struct {
	key   string
	value string
}
