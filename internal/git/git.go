package git

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/exec"
)

func New() Git {
	return Git{}
}

type Git struct{}

func (g Git) Clone(url, path string) error {
	if isGitRepo(path) {
		return nil
	}

	cmd := exec.Command("git", "clone", url, path)
	err := cmd.Run()
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
