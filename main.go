package main

import (
	"os"

	"github.com/akshatdalton/harbor-cli/cmd"
)

func main() {
	harborCLI := cmd.CreateHarborCLI()
	if err := harborCLI.Execute(); err != nil {
		os.Exit(1)
	}
}
