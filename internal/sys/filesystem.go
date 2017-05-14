package sys

import (
	"io/ioutil"
	"os"
)

func NewFilesystem() Filesystem {
	return Filesystem{}
}

type Filesystem struct{}

func (f Filesystem) Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func (f Filesystem) MkdirAll(name string, perm os.FileMode) error {
	return os.MkdirAll(name, perm)
}

func (f Filesystem) ReadFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func (f Filesystem) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func (f Filesystem) ListFiles(dirname string) ([]string, error) {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}

	names := make([]string, len(files))
	for i, file := range files {
		names[i] = file.Name()
	}

	return names, nil
}
