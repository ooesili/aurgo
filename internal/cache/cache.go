package cache

import (
	"io/ioutil"
	"path/filepath"
)

type Git interface {
	Clone(url, pkg string) error
}

type Config interface {
	SourcePath(pkg string) (string, error)
	AurRepoURL(string) string
}

type SrcInfo interface {
	Parse(input []byte) (Package, error)
}

type Package struct {
	Depends      []string
	Checkdepends []string
	Makedepends  []string
}

func New(config Config, git Git, srcinfo SrcInfo) Cache {
	return Cache{
		config:  config,
		git:     git,
		srcinfo: srcinfo,
	}
}

type Cache struct {
	config  Config
	git     Git
	srcinfo SrcInfo
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

func (c Cache) GetDeps(pkgname string) ([]string, error) {
	sourcePath, err := c.config.SourcePath(pkgname)
	if err != nil {
		return nil, err
	}

	srcinfoPath := filepath.Join(sourcePath, ".SRCINFO")
	srcinfoBytes, err := ioutil.ReadFile(srcinfoPath)
	if err != nil {
		return nil, err
	}

	pkg, err := c.srcinfo.Parse(srcinfoBytes)
	if err != nil {
		return nil, err
	}

	deps := aggregateDeps(pkg)
	return deps, nil
}

func aggregateDeps(pkg Package) []string {
	var allDeps []string

	depLists := [][]string{pkg.Depends, pkg.Makedepends, pkg.Checkdepends}
	for _, depList := range depLists {
		for _, dep := range depList {
			allDeps = append(allDeps, dep)
		}
	}

	return allDeps
}