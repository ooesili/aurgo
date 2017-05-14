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
}

func (m *Filesystem) Exists(path string) (bool, error) {
	m.ExistsCall.Received.Path = path
	returns := m.ExistsCall.Returns
	return returns.Exists, returns.Err
}

func (m *Filesystem) MkdirAll(path string, mode os.FileMode) error {
	m.MkdirAllCall.Recieved.Path = path
	m.MkdirAllCall.Recieved.Mode = mode
	return m.MkdirAllCall.Returns.Err
}
