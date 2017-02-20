package pacman

import (
	"bytes"

	"github.com/pkg/exec"
)

func NewOsExecutor() OsExecutor {
	return OsExecutor{}
}

type OsExecutor struct{}

func (e OsExecutor) Execute(command string, args ...string) (*bytes.Buffer, error) {
	stdout := &bytes.Buffer{}

	cmd := exec.Command(command, args...)
	err := cmd.Run(exec.Stdout(stdout))
	if err != nil {
		return nil, err
	}

	return stdout, nil
}
