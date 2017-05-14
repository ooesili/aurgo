package sys

import (
	"bytes"
	"io"

	"github.com/pkg/exec"
)

func NewExecutor(defaultOut, defaultErr io.Writer) Executor {
	return Executor{
		defaultOut: defaultOut,
		defaultErr: defaultErr,
	}
}

type Executor struct {
	defaultOut io.Writer
	defaultErr io.Writer
}

func (e Executor) Execute(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	return cmd.Run(
		exec.Stdout(e.defaultOut),
		exec.Stderr(e.defaultErr),
	)
}

func (e Executor) ExecuteCapture(command string, args ...string) (*bytes.Buffer, error) {
	cmd := exec.Command(command, args...)
	out := &bytes.Buffer{}

	err := cmd.Run(exec.Stdout(out))
	if err != nil {
		return nil, err
	}

	return out, nil
}
