package tests

import (
	"github.com/we-are-discussing-rest/web-crawler/utils"
	"testing"
)

func TestTrimUrl(t *testing.T) {
	t.Run("should trim a url successfully", func(t *testing.T) {
		url := "https://wikipedia.com/test"

		want := "wikipedia"
		got, err := utils.TrimURL(url)
		if err != nil {
			t.Errorf("error thrown when trimming %v", err)
		}

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("should trim a url successfully with multiple paths", func(t *testing.T) {
		url := "https://wikipedia.com/test/this/is/a/test"

		want := "wikipedia"
		got, err := utils.TrimURL(url)
		if err != nil {
			t.Errorf("error thrown when trimming %v", err)
		}

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}
