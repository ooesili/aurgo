package git

import (
	"errors"
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
		err := g.execPull(path)
		if err != nil {
			return fmt.Errorf("git pull failed: %s", err)
		}

		return nil
	}

	err := g.execClone(url, path)
	if err != nil {
		return fmt.Errorf("git clone failed: %s", err)
	}

	return nil
}

func (g Git) execClone(url, path string) error {
	cmd := exec.Command("git", "clone", url, path)
	err := cmd.Run(
		exec.Stdout(g.stdout),
		exec.Stderr(g.stderr),
	)
	if err != nil {
		return err
	}

	masterRefPath := filepath.Join(path, ".git", "refs", "heads", "master")
	if _, err := os.Stat(masterRefPath); err != nil {

		// best effort cleanup, not critical if it fails
		_ = os.RemoveAll(path)

		return errors.New("repo not found")
	}

	return nil
}

func (g Git) execPull(path string) error {
	cmd := exec.Command("git", "pull")
	return cmd.Run(
		exec.Stdout(g.stdout),
		exec.Stderr(g.stderr),
		exec.Dir(path),
	)
}

func isGitRepo(path string) bool {
	return doesFileExist(filepath.Join(path, ".git"))
}

func doesFileExist(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}
