package srcinfo

import "fmt"

var archMap = map[string]string{
	"amd64": "x86_64",
}

func ArchString(goarch string) (string, error) {
	arch, ok := archMap[goarch]
	if !ok {
		return "", fmt.Errorf("unsupported architecture: %s", goarch)
	}

	return arch, nil
}
