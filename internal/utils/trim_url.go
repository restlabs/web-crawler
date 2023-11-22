package utils

import (
	"net/url"
	"strings"
)

func TrimURL(inputURL string) (string, error) {
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return "", err
	}

	trimmedHost := strings.TrimPrefix(parsedURL.Host, "www.")
	if trimmedHost == "" {
		trimmedHost = parsedURL.Host
	}

	hostParts := strings.Split(trimmedHost, ".")
	host := hostParts[len(hostParts)-2]

	return host, nil
}
