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

	// Remove all trailing slashes from the URL
	url = strings.TrimRight(url, "/")
	return url
}

func FormatSize(size int64) string {
	mbSize := float64(size) / (1024 * 1024)
	return fmt.Sprintf("%.2fMiB", mbSize)
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

// check if the password format is vaild
func ValidatePassword(password string) error {
	if len(password) < 8 || len(password) > 256 {
		return errors.New("worong! the password length must be at least 8 characters and at most 256 characters")
	}
	// checking the password has a minimum of one lower case letter
	if done, _ := regexp.MatchString("([a-z])+", password); !done {
		return errors.New("worong! the password doesn't have a lowercase letter")
	}

	// checking the password has a minimmum of one upper case letter
	if done, _ := regexp.MatchString("([A-Z])+", password); !done {
		return errors.New("worong! the password doesn't have an upppercase letter")
	}

	// checking if the password has a minimum of one digit
	if done, _ := regexp.Match("([0-9])+", []byte(password)); !done {
		return errors.New("worong! the password doesn't have a digit number")
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
func FormatURL(rawURL string) *string {
	rawURL = strings.TrimSpace(rawURL)
	rawURL = strings.TrimLeft(rawURL, "/")
	rawURL = strings.TrimRight(rawURL, "/")
	return &rawURL
}

func isValidDomainName(domain string) bool {
	pattern := `^(?i)[a-z0-9-]+(\.[a-z0-9-]+)*\.[a-z]{2,}$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(domain)
}

func ValidateURL(rawURL string) error {
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = "https://" + rawURL
	}

	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL format: %v", err)
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("URL must use http or https scheme")
	}

	host := parsedURL.Hostname()

	if host == "" {
		return fmt.Errorf("URL must contain a valid host")
	}
	if strings.Contains(host, "_") {
		return fmt.Errorf("domain contains invalid character '_'")
	}
	if strings.HasPrefix(host, "-") || strings.HasPrefix(host, ".") {
		return fmt.Errorf("domain cannot start with a '-' or '.'")
	}
	if strings.Contains(host, "..") {
		return fmt.Errorf("domain contains double periods")
	}
	if parts := strings.Split(host, "."); len(parts) > 1 && len(parts[len(parts)-1]) < 2 {
		return fmt.Errorf("TLD too short")
	}
	if net.ParseIP(host) == nil && !isValidDomainName(host) {
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
