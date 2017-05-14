package sys

import "os"

func NewFilesystem() Filesystem {
	return Filesystem{}
}

type Filesystem struct{}

func (o Filesystem) Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func (o Filesystem) MkdirAll(name string, perm os.FileMode) error {
	return os.MkdirAll(name, perm)
}
