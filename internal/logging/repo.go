package logging

import (
	"io"

	"github.com/ooesili/aurgo/internal/aurgo"
)

func NewRepo(repo aurgo.Repo, out io.Writer) Repo {
	return Repo{
		Repo:   repo,
		logger: newLogger(out),
	}
}

type Repo struct {
	aurgo.Repo
	logger
}

func (c Repo) Sync(pkg string) error {
	c.log("syncing package: %s", pkg)
	return c.Repo.Sync(pkg)
}

func (c Repo) Remove(pkg string) error {
	c.log("removing package: %s", pkg)
	return c.Repo.Remove(pkg)
}
