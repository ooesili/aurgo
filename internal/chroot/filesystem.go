package chroot

import "os"

func NewOSFilesystem() OSFilesystem {
	return OSFilesystem{}
}

type OSFilesystem struct{}

func (o OSFilesystem) Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func (o OSFilesystem) MkdirAll(name string, perm os.FileMode) error {
	return os.MkdirAll(name, perm)
}
