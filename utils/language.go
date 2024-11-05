package utils

import (
	"net/url"
	"regexp"
	"strings"
)

// Regex for English characters, numbers, and common URL symbols
var englishPattern = regexp.MustCompile("^[a-zA-Z0-9-.]+$")

func IsEnglishURL(inputURL string) bool {
	// Parse the URL
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return false
	}

	hostname := parsedURL.Hostname()
	hostname = strings.TrimPrefix(hostname, "www.")
	hostname = strings.TrimSuffix(hostname, ".com")
	hostname = strings.TrimSuffix(hostname, ".org")
	hostname = strings.TrimSuffix(hostname, ".net")
	hostname = strings.TrimSuffix(hostname, ".xyz")

	// Check if the hostname matches our pattern
	if !englishPattern.MatchString(hostname) {
		return false
	}

	// Check if the path also contains only English characters
	if parsedURL.Path != "" {
		pathPattern := regexp.MustCompile("^[a-zA-Z0-9-._/]+$")
		if !pathPattern.MatchString(parsedURL.Path) {
			return false
		}
	}

	return true
}

func IsEnglishWord(inputWord string) bool {

	return englishPattern.MatchString(inputWord)

}
