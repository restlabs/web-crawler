package tests

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/we-are-discussing-rest/web-crawler/workers/internal"
	"reflect"
	"testing"
)

type MockRedisRepo struct {
}

type MockRDRepo struct {
	data map[string]string
}

func (m *MockRDRepo) Insert(data string) error {
	m.data[data] = data
	return nil
}

func (m *MockRDRepo) Remove(data string) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockRDRepo) Get(data string) (string, error) {
	r, ok := m.data[data]
	if !ok {
		return "", nil
	}

	return r, nil
}

func TestWorks(t *testing.T) {
	t.Run("should resolve correct IP address for host", func(t *testing.T) {
		err := internal.ResolveDns("wikipedia.org")
		if err != nil {
			t.Errorf("expected not to have an error but got one")
		}
	})

	t.Run("should throw a lookup error if no ip exists", func(t *testing.T) {
		err := internal.ResolveDns("jklajkdlafnklsg.com")
		if err == nil {
			t.Errorf("expected to return an error")
		}
	})

	t.Run("should parse html content", func(t *testing.T) {
		err := internal.ValidateHtmlContent(rawHtml)
		if err != nil {
			t.Errorf("was expecting not to receive an error while parsing html but got one")
		}
	})

	t.Run("should return an error if html is invalid", func(t *testing.T) {
		err := internal.ValidateHtmlContent(invalidHtml)
		if err == nil {
			t.Errorf("was expecting to receive an error while parsing html but didnt get one")
		}
	})

	t.Run("should check MD5 hash against content store and add to store if not exist", func(t *testing.T) {
		db := MockRDRepo{data: map[string]string{}}
		err := internal.CheckContent(&db, rawHtml)
		if errors.Is(err, internal.ContentErrorDuplicateHash) {
			t.Errorf("got a duplicate content error and was not expecting one")
		}

		hs, hsErr := hashString(rawHtml)
		if hsErr != nil {
			t.Errorf("got an error while hashing %v", hsErr)
		}
		dbVal, err := db.Get(hs)

		if dbVal != hs {
			t.Errorf("got %v, want %v", dbVal, hs)
		}

	})

	t.Run("should check MD5 hash against content store with content and throw if exists", func(t *testing.T) {
		hs, hsErr := hashString(rawHtml)
		if hsErr != nil {
			t.Errorf("got an error while hashing %v", hsErr)
		}
		db := MockRDRepo{data: map[string]string{
			hs: hs,
		}}

		err := internal.CheckContent(&db, rawHtml)
		if !errors.Is(err, internal.ContentErrorDuplicateHash) {
			t.Errorf("was expecting a duplicate content error and did not get one")
		}
	})

	t.Run("should extract links from raw HTML data", func(t *testing.T) {
		links, err := internal.ExtractHtmlLinks(rawHtml)

		if err != nil {
			t.Errorf("got an error %v", err)
		}

		if !reflect.DeepEqual(links, []string{"https://www.iana.org/domains/example"}) {
			t.Errorf("got %v, want %v", links, []string{"https://www.iana.org/domains/example"})
		}
	})
}

func hashString(content string) (string, error) {
	hash := sha256.New()

	_, err := hash.Write([]byte(content))
	if err != nil {
		return "", err
	}

	hashBytes := hash.Sum(nil)

	hashString := hex.EncodeToString(hashBytes)
	return hashString, nil
}

var rawHtml = `<!doctype html>
<html>
<head>
    <title>Example Domain</title>

    <meta charset="utf-8" />
    <meta http-equiv="Content-type" content="text/html; charset=utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <style type="text/css">
    body {
        background-color: #f0f0f2;
        margin: 0;
        padding: 0;
        font-family: -apple-system, system-ui, BlinkMacSystemFont, "Segoe UI", "Open Sans", "Helvetica Neue", Helvetica, Arial, sans-serif;
        
    }
    div {
        width: 600px;
        margin: 5em auto;
        padding: 2em;
        background-color: #fdfdff;
        border-radius: 0.5em;
        box-shadow: 2px 3px 7px 2px rgba(0,0,0,0.02);
    }
    a:link, a:visited {
        color: #38488f;
        text-decoration: none;
    }
    @media (max-width: 700px) {
        div {
            margin: 0 auto;
            width: auto;
        }
    }
    </style>    
</head>

<body>
<div>
    <h1>Example Domain</h1>
    <p>This domain is for use in illustrative examples in documents. You may use this
    domain in literature without prior coordination or asking for permission.</p>
    <p><a href="https://www.iana.org/domains/example">More information...</a></p>
</div>
</body>
</html>`

var invalidHtml = `<html><div><a>`
