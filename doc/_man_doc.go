package main

import (
	"fmt"
	"os"
	"time"

	cmd "github.com/goharbor/harbor-cli/cmd/harbor/root"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra/doc"
)

func ManDoc() error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	folderName := "man-docs"
	_, err = os.Stat(folderName)
	if os.IsNotExist(err) {
		err = os.Mkdir(folderName, 0755)
		if err != nil {
			log.Fatal("Error creating folder:", err)
		}
	}
	docDir := fmt.Sprintf("%s/%s", currentDir, folderName)

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
	return nil
}

func main() {
	err := ManDoc()
	if err != nil {
		log.Fatal(err)
	}
}
