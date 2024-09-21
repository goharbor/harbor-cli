package utils

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"unicode"

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

func ValidEmail(email string) bool {
	return regexp.MustCompile(`[a-z0-9]+@[a-z]+\.[a-z]{2,3}`).MatchString(email)
}

func ValidatePassword(s string) bool {
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
		return true
	}
	return false
}
