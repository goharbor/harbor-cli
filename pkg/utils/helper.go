package utils

import (
	"fmt"
	"strings"
	"time"
	"unicode"
)

func FormatCreatedTime(timestamp string) (string, error) {
	t, err := time.Parse(time.RFC3339Nano, timestamp)
	if err != nil {
		return "", err
	}

	duration := time.Since(t)

	minutes := int(duration.Minutes())
	hours := int(duration.Hours())
	days := int(duration.Hours() / 24)

	if minutes < 60 {
		return fmt.Sprintf("%d minute ago", minutes), nil
	} else if hours < 24 {
		return fmt.Sprintf("%d hour ago", hours), nil
	} else {
		return fmt.Sprintf("%d day ago", days), nil
	}
}

// This function covert camelCase to Human Readable form
func CamelCaseToHR(s string) string {
	var result []string
	var word []rune

	for i, r := range s {
		if unicode.IsUpper(r) && i > 0 {
			result = append(result, string(word))
			word = []rune{r}
		} else {
			word = append(word, r)
		}
	}

	result = append(result, string(word))

	// Capitalize the first letter of each word
	for i, word := range result {
		if len(word) > 0 {
			runes := []rune(word)
			runes[0] = unicode.ToUpper(runes[0])
			result[i] = string(runes)
		}
	}

	return strings.Join(result, " ")
}
