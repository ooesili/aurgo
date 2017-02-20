package mocks

type Cache struct {
	SyncCall struct {
		SyncedPackages []string
		Err            error
	}
	GetDepsCall struct {
		DepMap map[string][]string
		Err    error
	}
}

func (c *Cache) Sync(pkg string) error {
	if err := c.SyncCall.Err; err != nil {
		return err
	}

	c.SyncCall.SyncedPackages = append(c.SyncCall.SyncedPackages, pkg)
	return nil
}

func (c *Cache) GetDeps(pkg string) ([]string, error) {
	if err := c.GetDepsCall.Err; err != nil {
		return nil, err
	}

	deps, ok := c.GetDepsCall.DepMap[pkg]
	if !ok {
		panic("package not listed in DepMap: " + pkg)
	}

	return deps, nil
}
