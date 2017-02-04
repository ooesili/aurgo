package main

import (
	"fmt"
	"os"

	"github.com/ooesili/aurgo/internal/aur"
)

func main() {
	if err := realMain(); err != nil {
		fmt.Fprintf(os.Stderr, "aurgo: %s\n", err)
		os.Exit(1)
	}
}

func realMain() error {
	packageName := os.Args[2]

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
