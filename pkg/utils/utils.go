package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"syscall"
	"unicode"

	log "github.com/sirupsen/logrus"
	"golang.org/x/term"
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

func SavePayloadJSON(filename string, payload any) {
	// Marshal the payload into a JSON string with indentation
	jsonStr, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		panic(err)
	}
	// Define the filename
	filename = filename + ".json"
	err = os.WriteFile(filename, jsonStr, 0644)
	if err != nil {
		panic(err)
	}
	fmt.Printf("JSON data has been written to %s\n", filename)
}

// Validate the secret based on the provided guidelines
func ValidatePassword(s string) error {
	const (
		minLength = 8
		maxLength = 128
	)

	var (
		hasLen    = false
		hasUpper  = false
		hasLower  = false
		hasNumber = false
	)
	if len(s) >= minLength && len(s) <= maxLength {
		hasLen = true
	}
	for _, char := range s {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		}
	}
	if hasLen && hasUpper && hasLower && hasNumber {
		return nil
	}
	return fmt.Errorf("secret should contain at least 1 uppercase, 1 lowercase and 1 number.")
}

// Get Password as Stdin
func GetSecretStdin(prompt string) (string, error) {
	fmt.Print(prompt)
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	fmt.Println() // move to the next line after input
	return strings.TrimSpace(string(bytePassword)), nil
}
