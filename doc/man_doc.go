package main

import (
	"fmt"
	"os"
	"time"

	cmd "github.com/goharbor/harbor-cli/cmd/harbor/root"
	"github.com/spf13/cobra/doc"
)

func main() {
	// create temporary dir in currentDir for documents.
	// Assuming you are executing from the main directory.
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	docDir := fmt.Sprintf("%s/%s", currentDir, "doc/man_docs/")
	os.RemoveAll(docDir)
	err = os.MkdirAll(docDir, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating docs directory:", err)
		os.Exit(1)
	}

	t := time.Now()

	header := &doc.GenManHeader{
		Title:   "HARBOR",
		Section: "1",
		Source:  "Habor Community",
		Manual:  "Harbor User Mannuals",
		Date:    &t,
	}

	err = doc.GenManTree(cmd.RootCmd(), header, docDir)
	if err != nil {
		fmt.Println("Error generating documentation:", err)
		os.Exit(1)
	}

	fmt.Println("Documentation generated successfully in", docDir)
}
