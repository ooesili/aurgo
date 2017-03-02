package aurgo

type Config interface {
	Packages() []string
}

type DepWalker interface {
	Walk(pkgs []string) ([]string, error)
}

type Cleaner interface {
	Clean(pkgs []string) error
}

func New(depWalker DepWalker, cleaner Cleaner, config Config) Aurgo {
	return Aurgo{
		config:    config,
		depWalker: depWalker,
		cleaner:   cleaner,
	}
}

type Aurgo struct {
	config    Config
	depWalker DepWalker
	cleaner   Cleaner
}

func (a Aurgo) SyncAll() error {
	requiredPackages := a.config.Packages()

	allPackages, err := a.depWalker.Walk(requiredPackages)
	if err != nil {
		return err
	}

	return a.cleaner.Clean(allPackages)
}
