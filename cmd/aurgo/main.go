package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/ooesili/aurgo/internal/aur"
	"github.com/ooesili/aurgo/internal/chroot"
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
	factory, err := newFactory()
	if err != nil {
		return err
	}

	aurgo := factory.getAurgo()

	err = aurgo.SyncAll()
	if err != nil {
		return err
	}

	return nil
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
	factory, err := newFactory()
	if err != nil {
		return err
	}

	buildManager := factory.getBuildManager()

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
