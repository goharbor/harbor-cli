package utils

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	log "github.com/sirupsen/logrus"
)

// Returns Harbor v2 client for given clientConfig

func PrintPayloadInJSONFormat(payload any) {
	if payload == nil {
		return
	}

	jsonStr, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(jsonStr))
}

func ParseProjectRepo(projectRepo string) (string, string) {
	split := strings.Split(projectRepo, "/")
	if len(split) != 2 {
		log.Fatalf("invalid project/repository format: %s", projectRepo)
	}
	return split[0], split[1]
}

func ParseProjectRepoReference(projectRepoReference string) (string, string, string) {
	split := strings.Split(projectRepoReference, "/")
	if len(split) != 3 {
		log.Fatalf("invalid project/repository/reference format: %s", projectRepoReference)
	}
	return split[0], split[1], split[2]
}

// RemoveColumns removes columns with specified titles from the given columns array.
func RemoveColumns(columns []table.Column, colsToRemove []string) []table.Column {
	titleMap := make(map[string]bool)
	for _, title := range colsToRemove {
		titleMap[title] = true
	}

	var filteredColumns []table.Column
	for _, column := range columns {
		if !titleMap[column.Title] {
			filteredColumns = append(filteredColumns, column)
		}
	}

	return filteredColumns
}

