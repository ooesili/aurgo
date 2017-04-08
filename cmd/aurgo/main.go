package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/ooesili/aurgo/internal/aur"
	"github.com/ooesili/aurgo/internal/aurgo"
	"github.com/ooesili/aurgo/internal/cache"
	"github.com/ooesili/aurgo/internal/chroot"
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
	case "mkchroot":
		return mkchroot()
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

func buildConfig() (config.Config, error) {
	repoPath := os.Getenv("AURGOPATH")
	return config.New(repoPath)
}

func buildAurgo() (aurgo.Aurgo, error) {
	config, err := buildConfig()
	if err != nil {
		return aurgo.Aurgo{}, err
	}

	pacman := pacman.New(pacman.NewOsExecutor())

	arch, err := srcinfo.ArchString(runtime.GOARCH)
	if err != nil {
		return aurgo.Aurgo{}, err
	}

	repo := logging.NewRepo(
		cache.New(
			config,
			git.New(os.Stdout, os.Stderr),
			srcinfo.New(arch),
		),
		os.Stdout,
	)

	aurgo := aurgo.New(
		aurgo.NewVisitingDepWalker(
			aurgo.NewFilteringVisitor(
				aurgo.NewRepoVisitor(repo),
				pacman,
			),
		),
		aurgo.NewRepoCleaner(repo),
		config,
	)

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

func mkchroot() error {
	config, err := buildConfig()
	if err != nil {
		return err
	}

	buildManager := aurgo.NewBuildManager(
		chroot.New(
			newExecutorLogger(
				chroot.NewOSExecutor(os.Stdout, os.Stderr),
			),
			chroot.NewOSFilesystem(),
		),
		config,
	)

	return buildManager.Provision()
}

func newExecutorLogger(executor chroot.Executor) executorLogger {
	return executorLogger{
		executor: executor,
	}
}

type executorLogger struct {
	executor chroot.Executor
}

func (e executorLogger) Execute(command string, args ...string) error {
	fmt.Printf("===> $ %s %s\n", command, strings.Join(args, " "))
	return e.executor.Execute(command, args...)
}
