package chroot

import (
	"os"
	"path/filepath"
)

type Executor interface {
	Execute(command string, args ...string) error
}

type Filesystem interface {
	Exists(path string) (bool, error)
	MkdirAll(name string, mode os.FileMode) error
}

func New(executor Executor, filesystem Filesystem) Chroot {
	return Chroot{
		executor:   executor,
		filesystem: filesystem,
	}
}

type Chroot struct {
	executor   Executor
	filesystem Filesystem
}

func (c Chroot) Create(location string) error {
	err := c.filesystem.MkdirAll(location, 0755)
	if err != nil {
		return err
	}

	chrootPath := filepath.Join(location, "root")
	return c.executor.Execute("mkarchroot", chrootPath, "base-devel")
}

func (c Chroot) Exists(location string) (bool, error) {
	indicatorFile := filepath.Join(location, "root", ".arch-chroot")
	return c.filesystem.Exists(indicatorFile)
}
