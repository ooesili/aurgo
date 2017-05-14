package mocks

import "os"

type Filesystem struct {
	ExistsCall struct {
		Received struct {
			Path string
		}
		Returns struct {
			Exists bool
			Err    error
		}
	}
	MkdirAllCall struct {
		Recieved struct {
			Path string
			Mode os.FileMode
		}
		Returns struct {
			Err error
		}
	}
	ReadFileCall struct {
		Received struct {
			Path string
		}
		Returns struct {
			Bytes []byte
			Err   error
		}
	}
	ListFilesCall struct {
		Received struct {
			Dirname string
		}
		Returns struct {
			Filenames []string
			Err       error
		}
	}
	RemoveAllCall struct {
		Received struct {
			Path string
		}
		Returns struct {
			Err error
		}
	}
}

func (f *Filesystem) Exists(path string) (bool, error) {
	f.ExistsCall.Received.Path = path
	returns := f.ExistsCall.Returns
	return returns.Exists, returns.Err
}

func (f *Filesystem) MkdirAll(path string, mode os.FileMode) error {
	f.MkdirAllCall.Recieved.Path = path
	f.MkdirAllCall.Recieved.Mode = mode
	return f.MkdirAllCall.Returns.Err
}

func (f *Filesystem) ReadFile(path string) ([]byte, error) {
	f.ReadFileCall.Received.Path = path
	returns := f.ReadFileCall.Returns
	return returns.Bytes, returns.Err
}

func (f *Filesystem) ListFiles(dirname string) ([]string, error) {
	f.ListFilesCall.Received.Dirname = dirname
	returns := f.ListFilesCall.Returns
	return returns.Filenames, returns.Err
}

func (f *Filesystem) RemoveAll(path string) error {
	f.RemoveAllCall.Received.Path = path
	return f.RemoveAllCall.Returns.Err
}
