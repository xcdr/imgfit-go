package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestStatus(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/server-status")

	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Server didn’t respond 200 OK: %s", resp.Status)
	}
}

func TestBadRequest(t *testing.T) {
	invalidURLs := []string{
		"http://localhost:8080/1/01/sample1.gif",
		"http://localhost:8080/11/01/sample1.jpg",
		"http://localhost:8080/1/01/sample1.jpg?v=abc",
	}

	for _, url := range invalidURLs {
		resp, err := http.Get(url)

		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Server didn’t respond 400 Bad Request: %s for: %s", resp.Status, url)
		}

		resp.Body.Close()
	}
}

func TestNotFound(t *testing.T) {
	invalidURLs := []string{
		"http://localhost:8080/1/01/not-found.jpg",
	}

	for _, url := range invalidURLs {
		resp, err := http.Get(url)

		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Server didn’t respond 404 Not Found: %s for: %s", resp.Status, url)
		}

		resp.Body.Close()
	}
}

func TestFound(t *testing.T) {
	validURLs := []string{
		"http://localhost:8080/1/01/sample1.jpg",
		"http://localhost:8080/1/01/sample1.jpg?v=123",
		"http://lvh.me:8080/1/02/sample2.jpg",
		"http://LVH.ME:8080/1/02/sample2.jpg?v=123",
	}

	for _, url := range validURLs {
		resp, err := http.Get(url)

		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Server didn’t respond 200 OK: %s for: %s", resp.Status, url)
		}

		resp.Body.Close()
	}
}

func BenchmarkStatus(b *testing.B) {
	for i := 0; i < b.N; i++ {
		resp, err := http.Get("http://localhost:8080/server-status")

		if err != nil {
			b.Error(err)
		}

		if resp.StatusCode != http.StatusOK {
			b.Errorf("Server didn’t respond 200 OK: %s", resp.Status)
		}

		resp.Body.Close()
	}
}

func BenchmarkNotFound(b *testing.B) {
	for i := 0; i < b.N; i++ {
		resp, err := http.Get("http://localhost:8080/1/01/not-found.jpg?v=123")

		if err != nil {
			b.Error(err)
		}

		if resp.StatusCode != http.StatusNotFound {
			b.Errorf("Server didn’t respond 404 Not Found: %s", resp.Status)
		}

		resp.Body.Close()
	}
}

func BenchmarkNotCached(b *testing.B) {
	for i := 0; i < b.N; i++ {
		resp, err := http.Get(fmt.Sprintf("http://localhost:8080/1/01/sample1.jpg?v=%d", i))

		if err != nil {
			b.Error(err)
		}

		if resp.StatusCode != http.StatusOK {
			b.Errorf("Server didn’t respond 200 OK: %s", resp.Status)
		}

		resp.Body.Close()
	}
}

func BenchmarkCached(b *testing.B) {
	for i := 0; i < b.N; i++ {
		resp, err := http.Get("http://localhost:8080/1/01/sample1.jpg?v=123")

		if err != nil {
			b.Error(err)
		}

		if resp.StatusCode != http.StatusOK {
			b.Errorf("Server didn’t respond 200 OK: %s", resp.Status)
		}

		resp.Body.Close()
	}
}
