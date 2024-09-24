package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
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

func FormatUrl(url string) string {
	// Check if URL starts with "http://" or "https://"
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		// If not, prepend "https://"
		url = "https://" + url
	}
	return url
}

func FormatSize(size int64) string {
	mbSize := float64(size) / (1024 * 1024)
	return fmt.Sprintf("%.2fMiB", mbSize)
}

// check if the password format is vaild
func ValidatePassword(password string) (bool, error) {
	if len(password) < 8 || len(password) > 256 {
		return false, errors.New("worong! the password length must be at least 8 characters and at most 256 characters")
	}
	// checking the password has a minimum of one lower case letter
	if done, _ := regexp.MatchString("([a-z])+", password); !done {
		return false, errors.New("worong! the password doesn't have a lowercase letter")
	}

	// checking the password has a minimmum of one upper case letter
	if done, _ := regexp.MatchString("([A-Z])+", password); !done {
		return false, errors.New("worong! the password doesn't have an upppercase letter")
	}

	// checking if the password has a minimum of one digit
	if done, _ := regexp.Match("([0-9])+", []byte(password)); !done {
		return false, errors.New("worong! the password doesn't have a digit number")
	}

	return true, errors.New("")
}

// check if the tag name is valid
func VaildateTagName(tagname string) bool {
	pattern := `^[\w][\w.-]{0,127}$`

	re := regexp.MustCompile(pattern)

	return re.MatchString(tagname)
}
