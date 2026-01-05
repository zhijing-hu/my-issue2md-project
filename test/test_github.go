package test

import (
	"net/http"
	"net/http/httptest"
)

// MockGitHubServer creates a test HTTP server that mocks GitHub API responses
func MockGitHubServer() *httptest.Server {
	mux := http.NewServeMux()

	// Mock issue API response
	issueJSON := `{
		"title": "Test Issue",
		"body": "This is a test issue.",
		"html_url": "https://github.com/test/repo/issues/1",
		"user": {
			"login": "testuser"
		},
		"created_at": "2024-01-15T10:00:00Z",
		"labels": [
			{"name": "bug"},
			{"name": "good first issue"}
		],
		"milestone": {
			"title": "v1.0"
		},
		"number": 1
	}`

	// Mock comments API response
	commentsJSON := `[
		{
			"body": "First comment",
			"user": {
				"login": "commenter1"
			},
			"created_at": "2024-01-15T11:00:00Z"
		},
		{
			"body": "Second comment",
			"user": {
				"login": "commenter2"
			},
			"created_at": "2024-01-16T09:00:00Z"
		}
	]`

	mux.HandleFunc("/repos/test/repo/issues/1", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if r.Header.Get("Accept") != "application/vnd.github.v3+json" {
			http.Error(w, "invalid accept header", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(issueJSON))
	})

	mux.HandleFunc("/repos/test/repo/issues/1/comments", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(commentsJSON))
	})

	mux.HandleFunc("/repos/error/repo/issues/1", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	})

	mux.HandleFunc("/rate_limited", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "rate limit exceeded", http.StatusForbidden)
	})

	mux.HandleFunc("/unauthorized", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "requires authentication", http.StatusUnauthorized)
	})

	return httptest.NewServer(mux)
}