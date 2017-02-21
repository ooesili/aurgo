package logging

import (
	"io"

	"github.com/ooesili/aurgo/internal/aurgo"
)

func NewCache(cache aurgo.Cache, out io.Writer) Cache {
	return Cache{
		cache:  cache,
		logger: newLogger(out),
	}
}

type Cache struct {
	cache aurgo.Cache
	logger
}

func (c Cache) GetDeps(pkg string) ([]string, error) {
	return c.cache.GetDeps(pkg)
}

func (c Cache) Sync(pkg string) error {
	c.log("syncing package: %s", pkg)
	return c.cache.Sync(pkg)
}
