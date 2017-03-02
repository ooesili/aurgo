package mocks

type Repo struct {
	SyncCall struct {
		Received struct {
			Pkg string
		}
		Returns struct {
			Err error
		}
	}
	GetDepsCall struct {
		Received struct {
			Pkg string
		}
		Returns struct {
			Pkgs []string
			Err  error
		}
	}
	ListCall struct {
		Returns struct {
			Pkgs []string
			Err  error
		}
	}
	RemoveCall struct {
		Removed []string
		Err     error
	}
}

func (m *Repo) Sync(pkg string) error {
	m.SyncCall.Received.Pkg = pkg
	return m.SyncCall.Returns.Err
}

func (m *Repo) GetDeps(pkg string) ([]string, error) {
	m.GetDepsCall.Received.Pkg = pkg
	returns := m.GetDepsCall.Returns
	return returns.Pkgs, returns.Err
}

func (m *Repo) List() ([]string, error) {
	returns := m.ListCall.Returns
	return returns.Pkgs, returns.Err
}

func (m *Repo) Remove(pkg string) error {
	m.RemoveCall.Removed = append(m.RemoveCall.Removed, pkg)
	return m.RemoveCall.Err
}
