package cache

import (
	"path/filepath"
	"sort"
)

type Git interface {
	Clone(url, pkg string) error
}

type Config interface {
	SourcePath(pkg string) string
	AurRepoURL(string) string
	SourceBase() string
}

type SrcInfo interface {
	Parse(input []byte) (Package, error)
}

type Filesystem interface {
	ReadFile(path string) ([]byte, error)
	ListFiles(dirname string) ([]string, error)
	RemoveAll(path string) error
}

type Package struct {
	Depends      []string
	Checkdepends []string
	Makedepends  []string
}

func New(config Config, git Git, srcinfo SrcInfo, fs Filesystem) Cache {
	return Cache{
		config:  config,
		git:     git,
		srcinfo: srcinfo,
		fs:      fs,
	}
}

type Cache struct {
	config  Config
	git     Git
	srcinfo SrcInfo
	fs      Filesystem
}

func (c Cache) Sync(pkg string) error {
	sourcePath := c.config.SourcePath(pkg)
	aurRepoURL := c.config.AurRepoURL(pkg)

	err := c.git.Clone(aurRepoURL, sourcePath)
	if err != nil {
		return err
	}

	return nil
}

func (c Cache) GetDeps(pkgname string) ([]string, error) {
	sourcePath := c.config.SourcePath(pkgname)

	srcinfoPath := filepath.Join(sourcePath, ".SRCINFO")
	srcinfoBytes, err := c.fs.ReadFile(srcinfoPath)
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

func (c Cache) List() ([]string, error) {
	sourceBase := c.config.SourceBase()
	pkgs, err := c.fs.ListFiles(sourceBase)
	if err != nil {
		return nil, err
	}

	sort.Strings(pkgs)
	return pkgs, nil
}

func (c Cache) Remove(pkg string) error {
	sourcePath := c.config.SourcePath(pkg)
	return c.fs.RemoveAll(sourcePath)
}
