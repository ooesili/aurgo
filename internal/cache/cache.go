package cache

type Git interface {
	Clone(url, pkg string) error
}

type Config interface {
	SourcePath(pkg string) (string, error)
	AurRepoURL(string) string
}

func New(config Config, git Git) Cache {
	return Cache{
		config: config,
		git:    git,
	}
}

type Cache struct {
	config Config
	git    Git
}

func (c Cache) Sync(pkg string) error {
	sourcePath, err := c.config.SourcePath(pkg)
	if err != nil {
		return err
	}
	aurRepoURL := c.config.AurRepoURL(pkg)

	err = c.git.Clone(aurRepoURL, sourcePath)
	if err != nil {
		return err
	}

	return nil
}
