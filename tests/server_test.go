package tests

import (
	"fmt"
	"github.com/we-are-discussing-rest/web-crawler/cmd/server"
	"github.com/we-are-discussing-rest/web-crawler/logger"
	"github.com/we-are-discussing-rest/web-crawler/utils"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

type MockRepository struct {
	queue map[string][]string
}

func (m *MockRepository) Insert(data string) error {
	queueName, err := utils.TrimURL(data)
	if err != nil {
		fmt.Errorf("error trimming url: %v", err)
		return err
	}

	m.queue[queueName] = append(m.queue[queueName], data)
	return nil
}

func (m *MockRepository) Remove(data string) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockRepository) Get(data string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func NewMockRepo() *MockRepository {
	m := make(map[string][]string)
	return &MockRepository{queue: m}
}

func TestServer(t *testing.T) {
	r := NewMockRepo()
	l := logger.NewLogger()
	s := server.NewServer(r, l)

	t.Run("healthcheck should return a 200", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/healthcheck", nil)
		res := httptest.NewRecorder()
		s.ServeHTTP(res, req)

		got := res.Code
		want := http.StatusOK

		assertResponseStatus(t, got, want)
	})

	t.Run("should add correct queue and return a 201", func(t *testing.T) {
		rawBody := `{
	"seedUrls": [
		"https://wikipedia.com/test",
		"https://pkg.go.dev/encoding/json",
		"https://pkg.go.dev/net/http#Redirect",
		"https://en.wikipedia.org/wiki/Breadth-first_search"
	]
}`
		expectedQueue := map[string][]string{
			"go": {
				"https://pkg.go.dev/encoding/json",
				"https://pkg.go.dev/net/http#Redirect",
			},
			"wikipedia": {
				"https://wikipedia.com/test",
				"https://en.wikipedia.org/wiki/Breadth-first_search",
			},
		}

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/seed", strings.NewReader(rawBody))
		res := httptest.NewRecorder()
		s.ServeHTTP(res, req)

		assertResponseStatus(t, res.Code, http.StatusCreated)
		asserUrlFrontierMQ(t, r, 2, expectedQueue)
	})
}

func assertResponseStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("incorrect status code, got %v, want %v", got, want)
	}
}

func asserUrlFrontierMQ(t *testing.T, q *MockRepository, expectedLen int, expectedQueue map[string][]string) {
	t.Helper()

	if len(q.queue) != expectedLen {
		t.Errorf("queue length is incorrect, expect %v, got %v", expectedLen, len(q.queue))
	}

	if !reflect.DeepEqual(q.queue, expectedQueue) {
		t.Errorf("expect queue %v, got queue %v", expectedQueue, q.queue)
	}
}
