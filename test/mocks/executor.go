package mocks

import "bytes"

type Executor struct {
	ExecuteCall struct {
		Recieved struct {
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
	e.ExecuteCall.Recieved.Command = command
	e.ExecuteCall.Recieved.Args = args
	returns := e.ExecuteCall.Returns
	return returns.Stdout, returns.Err
}
