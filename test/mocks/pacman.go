package mocks

type Pacman struct {
	ListAvailableCall struct {
		Returns struct {
			Packages []string
		}
	}
}

func (p *Pacman) ListAvailable() []string {
	return p.ListAvailableCall.Returns.Packages
}
