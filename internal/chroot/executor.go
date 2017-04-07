package chroot

import (
	"io"

	"github.com/pkg/exec"
)

func NewOSExecutor(stdout, stderr io.Writer) OSExecutor {
	return OSExecutor{
		stdout: stdout,
		stderr: stderr,
	}
}

type OSExecutor struct {
	stdout io.Writer
	stderr io.Writer
}

func (o OSExecutor) Execute(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	return cmd.Run(
		exec.Stdout(o.stdout),
		exec.Stderr(o.stderr),
	)
}
