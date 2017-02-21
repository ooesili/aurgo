package git

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/exec"
)

func New(stdout, stderr io.Writer) Git {
	return Git{
		stdout: stdout,
		stderr: stderr,
	}
}

type Git struct {
	stdout io.Writer
	stderr io.Writer
}

func (g Git) Clone(url, path string) error {
	if isGitRepo(path) {
		return nil
	}

	cmd := exec.Command("git", "clone", url, path)
	err := cmd.Run(
		exec.Stdout(g.stdout),
		exec.Stderr(g.stderr),
	)
	if err != nil {
		return fmt.Errorf("git clone failed: %s", err)
	}

	return nil
}

func isGitRepo(path string) bool {
	return doesFileExist(filepath.Join(path, ".git"))
}

func doesFileExist(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}
