package mocks

import "bytes"

type Executor struct {
	ExecuteCall struct {
		Received struct {
			Command string
			Args    []string
		}
		Returns struct {
			Err error
		}
	}
	ExecuteCaptureCall struct {
		Received struct {
			Command string
			Args    []string
		}
		Returns struct {
			Stdout *bytes.Buffer
			Err    error
		}
	}
}

func (m *Executor) Execute(command string, args ...string) error {
	m.ExecuteCall.Received.Command = command
	m.ExecuteCall.Received.Args = args
	return m.ExecuteCall.Returns.Err
}

func (e *Executor) ExecuteCapture(command string, args ...string) (*bytes.Buffer, error) {
	e.ExecuteCaptureCall.Received.Command = command
	e.ExecuteCaptureCall.Received.Args = args
	returns := e.ExecuteCaptureCall.Returns
	return returns.Stdout, returns.Err
}
