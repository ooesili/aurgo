package main

import (
	"io"
	"os"
	"runtime"

	"github.com/ooesili/aurgo/internal/aurgo"
	"github.com/ooesili/aurgo/internal/cache"
	"github.com/ooesili/aurgo/internal/chroot"
	"github.com/ooesili/aurgo/internal/config"
	"github.com/ooesili/aurgo/internal/git"
	"github.com/ooesili/aurgo/internal/logging"
	"github.com/ooesili/aurgo/internal/pacman"
	"github.com/ooesili/aurgo/internal/srcinfo"
	"github.com/ooesili/aurgo/internal/sys"
)

func newFactory() (*factory, error) {
	repoPath := os.Getenv("AURGOPATH")
	config, err := config.New(repoPath)
	if err != nil {
		return nil, err
	}

	targetArch, err := srcinfo.ArchString(runtime.GOARCH)
	if err != nil {
		return nil, err
	}

	return &factory{
		config:     config,
		targetArch: targetArch,
	}, nil
}

type factory struct {
	config     config.Config
	targetArch string
}

func (f *factory) getStdout() io.Writer {
	return os.Stdout
}

func (f *factory) getStderr() io.Writer {
	return os.Stderr
}

func (f *factory) getFileSystem() chroot.Filesystem {
	return sys.NewFilesystem()
}

func (f *factory) getBuildEnv() aurgo.BuildEnv {
	return chroot.New(
		newExecutorLogger(f.getExecutor()),
		f.getFileSystem(),
	)
}

func (f *factory) getBuildManager() aurgo.BuildManager {
	return aurgo.NewBuildManager(f.getBuildEnv(), f.config)
}

func (f *factory) getPkgManager() aurgo.PkgManager {
	return pacman.New(f.getExecutor())
}

func (f *factory) getGit() cache.Git {
	return git.New(f.getStdout(), f.getStderr())
}

func (f *factory) getSrcInfo() cache.SrcInfo {
	return srcinfo.New(f.targetArch)
}

func (f *factory) getRepo() aurgo.Repo {
	return logging.NewRepo(
		cache.New(f.config, f.getGit(), f.getSrcInfo()),
		f.getStdout(),
	)
}

func (f *factory) getRepoCleaner() aurgo.RepoCleaner {
	return aurgo.NewRepoCleaner(f.getRepo())
}

func (f *factory) getRepoVisitor() aurgo.Visitor {
	return aurgo.NewFilteringVisitor(
		aurgo.NewRepoVisitor(f.getRepo()),
		f.getPkgManager(),
	)
}

func (f *factory) getDepWalker() aurgo.DepWalker {
	return aurgo.NewVisitingDepWalker(f.getRepoVisitor())
}

func (f *factory) getAurgo() aurgo.Aurgo {
	return aurgo.New(f.getDepWalker(), f.getRepoCleaner(), f.config)
}

func (f *factory) getExecutor() sys.Executor {
	return sys.NewExecutor(f.getStdout(), f.getStderr())
}
