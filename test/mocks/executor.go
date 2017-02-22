package mocks

import "bytes"

type Executor struct {
	ExecuteCall struct {
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

func (e *Executor) Execute(command string, args ...string) (*bytes.Buffer, error) {
	e.ExecuteCall.Received.Command = command
	e.ExecuteCall.Received.Args = args
	returns := e.ExecuteCall.Returns
	return returns.Stdout, returns.Err
}
