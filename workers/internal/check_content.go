package internal

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/we-are-discussing-rest/web-crawler/workers/repository"
)

var ContentErrorDuplicateHash = errors.New("duplicate value")

func CheckContent(store repository.Repository, content string) error {
	hash := sha256.New()

	_, err := hash.Write([]byte(content))
	if err != nil {
		return err
	}

	hashBytes := hash.Sum(nil)

	hashString := hex.EncodeToString(hashBytes)

	value, err := store.Get(hashString)
	if err != nil {
		return err
	}

	if value != "" {
		return ContentErrorDuplicateHash
	}

	err = store.Insert(hashString)
	if err != nil {
		return err
	}

	return nil
}
