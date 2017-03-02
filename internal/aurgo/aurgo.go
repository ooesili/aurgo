package aurgo

import (
	"sort"
)

type Config interface {
	Packages() []string
}

type Cache interface {
	Sync(pkg string) error
	GetDeps(pkg string) ([]string, error)
	ListExisting() ([]string, error)
	Remove(pkg string) error
}

type Pacman interface {
	ListAvailable() []string
}

func New(config Config, cache Cache, pacman Pacman) Aurgo {
	return Aurgo{
		config: config,
		cache:  cache,
		pacman: pacman,
	}
}

type Aurgo struct {
	config Config
	cache  Cache
	pacman Pacman
}

func (a Aurgo) SyncAll() error {
	initialPkgs := a.config.Packages()
	availablePackages := a.pacman.ListAvailable()

	pkgList := newPkgList(initialPkgs, availablePackages)

	for !pkgList.allProcessed() {
		batch := pkgList.nextBatch()

		for _, pkg := range batch {
			err := a.cache.Sync(pkg)
			if err != nil {
				return err
			}

			deps, err := a.cache.GetDeps(pkg)
			if err != nil {
				return err
			}

			for _, dep := range deps {
				pkgList.queueForNextBatch(dep)
			}
		}
	}

	err := a.removeOldPackages(pkgList)
	if err != nil {
		return err
	}

	return nil
}

func (a Aurgo) removeOldPackages(pkgList *pkgList) error {
	existingPkgs, err := a.cache.ListExisting()
	if err != nil {
		return err
	}

	for _, existingPkg := range existingPkgs {
		if !pkgList.wasProcessed(existingPkg) {
			err := a.cache.Remove(existingPkg)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func newPkgList(initialPkgs, availablePackages []string) *pkgList {
	pkgList := &pkgList{
		available:   make(map[string]bool),
		unprocessed: make(map[string]bool),
		processed:   make(map[string]bool),
	}

	for _, pkg := range availablePackages {
		pkgList.available[pkg] = true
	}

	for _, pkg := range initialPkgs {
		pkgList.unprocessed[pkg] = true
	}

	return pkgList
}

type pkgList struct {
	available   map[string]bool
	processed   map[string]bool
	unprocessed map[string]bool
}

func (p *pkgList) nextBatch() []string {
	var batch []string

	for pkg := range p.unprocessed {
		batch = append(batch, pkg)
		p.processed[pkg] = true
	}
	p.unprocessed = make(map[string]bool)

	sort.Strings(batch)
	return batch
}

func (p *pkgList) queueForNextBatch(pkg string) {
	if p.available[pkg] || p.processed[pkg] {
		return
	}

	p.unprocessed[pkg] = true
}

func (p *pkgList) allProcessed() bool {
	return len(p.unprocessed) == 0
}

func (p *pkgList) wasProcessed(pkg string) bool {
	return p.processed[pkg]
}
