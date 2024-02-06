package main

import (
	"os"

	"github.com/goharbor/harbor-cli/cmd/harbor/root"
)

func main() {
	err := root.New().Execute()
	if err != nil {
		os.Exit(1)
	}
}
