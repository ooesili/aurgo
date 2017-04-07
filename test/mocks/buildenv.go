package mocks

type BuildEnv struct {
	CreateCall struct {
		Received struct {
			Location string
		}
		Returns struct {
			Err error
		}
	}
	ExistsCall struct {
		Received struct {
			Location string
		}
		Returns struct {
			Exists bool
			Err    error
		}
	}
}

func (b *BuildEnv) Create(location string) error {
	b.CreateCall.Received.Location = location
	return b.CreateCall.Returns.Err
}

func (b *BuildEnv) Exists(location string) (bool, error) {
	b.ExistsCall.Received.Location = location
	returns := b.ExistsCall.Returns
	return returns.Exists, returns.Err
}
