package internal

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/we-are-discussing-rest/web-crawler/internal/repository"
	"io"
)

var ContentErrorDuplicateHash = errors.New("duplicate value")

func CheckContent(store repository.Repository, content string) error {
	hash := md5.New()
	_, err := io.WriteString(hash, content)
	if err != nil {
		return fmt.Errorf("error writing hash %v", err)
	}

	value, err := store.Get(fmt.Sprintf("%x", hash))
	if err != nil {
		return err
	}

	if value != "" {
		return ContentErrorDuplicateHash
	}

	return nil
}
