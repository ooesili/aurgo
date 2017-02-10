package aurgo

type Config interface {
	Packages() ([]string, error)
}

type Cache interface {
	Sync(pkg string) error
}

func New(config Config, cache Cache) Aurgo {
	return Aurgo{
		config: config,
		cache:  cache,
	}
}

type Aurgo struct {
	config Config
	cache  Cache
}

func (a Aurgo) SyncAll() error {
	pkgs, err := a.config.Packages()
	if err != nil {
		return err
	}

	for _, pkg := range pkgs {
		err := a.cache.Sync(pkg)
		if err != nil {
			return err
		}
	}

	return nil
}
