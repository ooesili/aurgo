package aurgo

type Config interface {
	Packages() ([]string, error)
	AurRepoURL(pkg string) string
	SourcePath(pkg string) (string, error)
}

type Git interface {
	Clone(url, pkg string) error
}

func New(config Config, git Git) Aurgo {
	return Aurgo{
		config: config,
		git:    git,
	}
}

type Aurgo struct {
	config Config
	git    Git
}

func (a Aurgo) Sync() error {
	pkgs, err := a.config.Packages()
	if err != nil {
		return err
	}

	for _, pkg := range pkgs {
		err := a.syncPackage(pkg)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a Aurgo) syncPackage(pkg string) error {
	sourcePath, err := a.config.SourcePath(pkg)
	if err != nil {
		return err
	}
	aurURL := a.config.AurRepoURL(pkg)

	err = a.git.Clone(aurURL, sourcePath)
	if err != nil {
		return err
	}

	return nil
}
