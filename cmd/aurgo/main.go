package main

import (
	"fmt"
	"os"

	"github.com/ooesili/aurgo/internal/aur"
	"github.com/ooesili/aurgo/internal/aurgo"
	"github.com/ooesili/aurgo/internal/config"
	"github.com/ooesili/aurgo/internal/git"
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
	repoPath := os.Getenv("AURGOPATH")
	config, err := config.New(repoPath)
	if err != nil {
		return err
	}

	git := git.New()
	aurgo := aurgo.New(config, git)

	err = aurgo.Sync()
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
