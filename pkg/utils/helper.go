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
	return url
}

// ValidateDomain validates subdomain, IP, or top-level domain formats
func ValidateDomain(domain string) error {
	url := FormatUrl(domain)

	err := isValidURL(url)
	if err != nil {
		return err
	}
	return nil
}

func isValidURL(rawURL string) error {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return errors.New("invaild url scheme")
	}

	if parsedURL.Host == "" {
		return errors.New("missing domain")
	}

	host, port, err := net.SplitHostPort(parsedURL.Host)
	if err != nil {
		host = parsedURL.Host
	}

	if err := isValidHost(host); err != nil {
		return err
	}

	if port != "" && !isValidPort(port) {
		return errors.New("invalid port")
	}

	return nil
}

// Validate whether the port is valid
func isValidPort(port string) bool {
	portRegex := `^([1-9][0-9]{0,4})$`
	match, _ := regexp.MatchString(portRegex, port)
	if !match {
		return false
	}
	portNumber := 0
	fmt.Sscanf(port, "%d", &portNumber)
	return portNumber >= 1 && portNumber <= 65535
}

func extractIP(host string) string {
	ipv4Regex := `(?:\d{1,3}\.){3}\d{1,3}`

	re := regexp.MustCompile(ipv4Regex)
	ipv4Matches := re.FindString(host)

	if ipv4Matches != "" {
		return ipv4Matches
	}

	return ""
}

// Validate whether the host is valid
func isValidHost(host string) error {
	ip := extractIP(host)
	if ip != "" {
		return isValidIPv4(host)
	}

	parts := strings.Split(host, ".")
	if len(parts) < 2 {
		return errors.New("invalid host: must have at least one dot")
	}

	for _, part := range parts {
		if !isValidLabel(part) {
			return fmt.Errorf("invalid host label: %s", part)
		}
	}

	if len(parts[len(parts)-1]) < 2 {
		return errors.New("invalid top-level host: must be at least 2 characters")
	}
	return nil
}

func isValidIPv4(ip string) error {
	octets := strings.Split(ip, ".")
	if len(octets) != 4 {
		return errors.New("IP: consists of more than three octets")
	}

	for _, octet := range octets {
		num, err := strconv.Atoi(octet)
		if err != nil || num < 0 || num > 255 {
			return errors.New("IP: octet exceeds range")
		}
	}
	return nil
}

// Helper function to validate individual domain labels
func isValidLabel(label string) bool {
	if len(label) == 0 || len(label) > 63 {
		return false
	}

	for i, ch := range label {
		if !(ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch >= '0' && ch <= '9' || ch == '-') {
			return false
		}
		if (i == 0 || i == len(label)-1) && ch == '-' {
			return false
		}
	}
	return true
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
