package logging

import (
	"io"

	"github.com/ooesili/aurgo/internal/aurgo"
)

func NewCache(cache aurgo.Cache, out io.Writer) Cache {
	return Cache{
		Cache:  cache,
		logger: newLogger(out),
	}
}

type Cache struct {
	aurgo.Cache
	logger
}

func (c Cache) Sync(pkg string) error {
	c.log("syncing package: %s", pkg)
	return c.Cache.Sync(pkg)
}

func (c Cache) Remove(pkg string) error {
	c.log("removing package: %s", pkg)
	return c.Cache.Remove(pkg)
}
