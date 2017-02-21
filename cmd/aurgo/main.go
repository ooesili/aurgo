package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/ooesili/aurgo/internal/aur"
	"github.com/ooesili/aurgo/internal/aurgo"
	"github.com/ooesili/aurgo/internal/cache"
	"github.com/ooesili/aurgo/internal/config"
	"github.com/ooesili/aurgo/internal/git"
	"github.com/ooesili/aurgo/internal/logging"
	"github.com/ooesili/aurgo/internal/pacman"
	"github.com/ooesili/aurgo/internal/srcinfo"
)

func main() {
	if err := realMain(); err != nil {
		fmt.Fprintf(os.Stderr, "aurgo: %s\n", err)
		os.Exit(1)
	}
}

func realMain() error {
	switch os.Args[1] {
	case "sync":
		return sync()
	case "info":
		return info(os.Args[2])
	default:
		panic("unknown command")
	}
}

func sync() error {
	aurgo, err := buildAurgo()
	if err != nil {
		return err
	}

	err = aurgo.SyncAll()
	if err != nil {
		return err
	}

	return nil
}

func buildAurgo() (aurgo.Aurgo, error) {
	repoPath := os.Getenv("AURGOPATH")
	config, err := config.New(repoPath)
	if err != nil {
		return aurgo.Aurgo{}, err
	}

	executor := pacman.NewOsExecutor()
	pacman, err := pacman.New(executor)
	if err != nil {
		return aurgo.Aurgo{}, err
	}

	arch, err := srcinfo.ArchString(runtime.GOARCH)
	if err != nil {
		return aurgo.Aurgo{}, err
	}
	srcinfo := srcinfo.New(arch)

	git := git.New(
		os.Stdout,
		os.Stderr,
	)
	cache := logging.NewCache(
		cache.New(config, git, srcinfo),
		os.Stdout,
	)

	aurgo := aurgo.New(config, cache, pacman)

	return aurgo, nil
}

func info(packageName string) error {
	api, err := aur.New("https://aur.archlinux.org")
	if err != nil {
		return err
	}

	version, err := api.Version(packageName)
	if err != nil {
		return err
	}

	fmt.Printf("aur/%s %s\n", packageName, version)
	return nil
}
