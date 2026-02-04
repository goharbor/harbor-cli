// Copyright Project Harbor Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package utils

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strconv"
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

func FormatUrl(url string) string {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}

	url = strings.TrimRight(url, "/")
	url = strings.TrimLeft(url, "/")

	return url
}

func FormatSize(size int64) string {
	const (
		_         = iota
		KiB int64 = 1 << (10 * iota)
		MiB
		GiB
		TiB
	)

	switch {
	case size >= TiB:
		return fmt.Sprintf("%.2fTiB", float64(size)/float64(TiB))
	case size >= GiB:
		return fmt.Sprintf("%.2fGiB", float64(size)/float64(GiB))
	case size >= MiB:
		return fmt.Sprintf("%.2fMiB", float64(size)/float64(MiB))
	case size >= KiB:
		return fmt.Sprintf("%.2fKiB", float64(size)/float64(KiB))
	default:
		return fmt.Sprintf("%dB", size)
	}
}

// ValidateUserName checks if the username is valid by length and allowed characters.
func ValidateUserName(username string) bool {
	username = strings.TrimSpace(username)
	return len(username) >= 1 && len(username) <= 255 && !strings.ContainsAny(username, `,"~#%$`)
}

func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(email)
}

func ValidateConfigPath(configPath string) bool {
	pattern := `^[\w./-]{1,255}\.(yaml|yml)$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(configPath)
}

func ValidateFL(name string) bool {
	pattern := `^[A-Za-z]{1,20}\s[A-Za-z]{1,20}$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(name)
}

// check if the password format is valid
func ValidatePassword(password string) error {
	password = strings.TrimSpace(password)
	if password == "" {
		return errors.New("password cannot be empty or only spaces")
	}

	if len(password) < 8 || len(password) > 256 {
		return errors.New("wrong! the password length must be at least 8 characters and at most 256 characters")
	}
	// checking the password has a minimum of one lower case letter
	if done, _ := regexp.MatchString("([a-z])+", password); !done {
		return errors.New("wrong! the password doesn't have a lowercase letter")
	}

	// checking the password has a minimum of one upper case letter
	if done, _ := regexp.MatchString("([A-Z])+", password); !done {
		return errors.New("wrong! the password doesn't have an uppercase letter")
	}

	// checking if the password has a minimum of one digit
	if done, _ := regexp.Match("([0-9])+", []byte(password)); !done {
		return errors.New("wrong! the password doesn't have a digit number")
	}

	return nil
}

// check if the tag name is valid
func ValidateTagName(tagName string) bool {
	pattern := `^[\w][\w.-]{0,127}$`

	re := regexp.MustCompile(pattern)

	return re.MatchString(tagName)
}

// check if the project name is valid
func ValidateProjectName(projectName string) bool {
	pattern := `^[a-z0-9][a-z0-9._-]{0,254}$`

	re := regexp.MustCompile(pattern)

	return re.MatchString(projectName)
}

func ValidateStorageLimit(sl string) error {
	storageLimit, err := strconv.Atoi(sl)
	if err != nil {
		return errors.New("the storage limit only takes integer values")
	}

	if storageLimit < -1 || (storageLimit > -1 && storageLimit < 0) || storageLimit > 1024 {
		return errors.New("the maximum value for the storage cannot exceed 1024 terabytes and -1 for no limit")
	}
	return nil
}

func ValidateRegistryName(rn string) bool {
	pattern := `^[\w][\w.-]{0,63}$`

	re := regexp.MustCompile(pattern)

	return re.MatchString(rn)
}

// ValidateURL checks if the URL has valid format, non-empty host, and host is a valid IP or domain.
// Domain regex: labels must start/end with alphanumeric, can contain hyphens, max 63 chars, TLD min 2 letters.
func ValidateURL(rawURL string) error {
	var domainNameRegex = regexp.MustCompile(`^(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`)

	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL format: %v", err)
	}

	host := parsedURL.Hostname()
	if host == "" {
		return fmt.Errorf("URL must contain a valid host")
	}

	if net.ParseIP(host) != nil {
		return nil
	}

	if !domainNameRegex.MatchString(host) {
		return fmt.Errorf("invalid host: must be a valid IP address or domain name")
	}

	return nil
}

func PrintFormat[T any](resp T, format string) error {
	if format == "json" {
		PrintPayloadInJSONFormat(resp)
		return nil
	}
	if format == "yaml" {
		PrintPayloadInYAMLFormat(resp)
		return nil
	}
	return fmt.Errorf("unable to output in the specified '%s' format", format)
}

func EmptyStringValidator(variable string) func(string) error {
	return func(str string) error {
		if str == "" {
			return fmt.Errorf("%s cannot be empty", variable)
		}
		return nil
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
