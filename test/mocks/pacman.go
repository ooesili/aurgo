package mocks

type Pacman struct {
	ListAvailableCall struct {
		Returns struct {
			Packages []string
			Err      error
		}
	}
}

func (p *Pacman) ListAvailable() ([]string, error) {
	returns := p.ListAvailableCall.Returns
	return returns.Packages, returns.Err
}
