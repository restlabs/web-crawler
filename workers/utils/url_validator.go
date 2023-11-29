package utils

import (
	"fmt"
	"net/url"
)

func IsValidURL(testURL string) bool {
	parsedURL, err := url.Parse(testURL)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return false
	}

	if parsedURL.Scheme == "" {
		return false
	}

	if parsedURL.Host == "" {
		return false
	}

	return true
}
