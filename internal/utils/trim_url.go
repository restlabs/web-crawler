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

	hostWithoutScheme := strings.TrimPrefix(parsedURL.Host, "www.")
	if hostWithoutScheme == "" {
		hostWithoutScheme = parsedURL.Host
	}

	hostParts := strings.Split(hostWithoutScheme, ".")
	host := hostParts[len(hostParts)-2]

	return host, nil
}
