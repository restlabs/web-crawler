package internal

import (
	"fmt"
	"io"
	"net/http"
)

func DownloadRawHtml(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error downloading html %v", err)
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("error reading html %v", err)
	}

	return fmt.Sprintf("%s", b), nil
}
