package aurgo

func NewRepoCleaner(repo Repo) RepoCleaner {
	return RepoCleaner{
		repo: repo,
	}
}

type RepoCleaner struct {
	repo Repo
}

func (c RepoCleaner) Clean(usedPkgs []string) error {
	existingPkgs, err := c.repo.List()
	if err != nil {
		return err
	}

	used := make(map[string]bool)
	for _, pkg := range usedPkgs {
		used[pkg] = true
	}

	for _, pkg := range existingPkgs {
		if !used[pkg] {
			err := c.repo.Remove(pkg)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
