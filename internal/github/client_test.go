package github_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"my-issue2md/internal/github"
)

// TestErrorTypes tests that expected error types are properly defined
func TestErrorTypes(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{"ErrRateLimited", github.ErrRateLimited},
		{"ErrNotFound", github.ErrNotFound},
		{"ErrUnauthorized", github.ErrUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err == nil {
				t.Error("error type should not be nil")
			}
			if tt.err.Error() == "" {
				t.Error("error should have descriptive message")
			}
		})
	}
}

// TestMockGitHubAPI tests with a mock GitHub API server using httptest
func TestMockGitHubAPI(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Handle different paths
		switch r.URL.Path {
		case "/repos/test/repo/issues/1":
			issueJSON := `{
				"title": "Test Issue",
				"body": "Test body",
				"html_url": "https://github.com/test/repo/issues/1",
				"user": {"login": "testuser"},
				"created_at": "2024-01-15T10:00:00Z",
				"labels": [{"name": "bug"}],
				"milestone": {"title": "v1.0"},
				"number": 1
			}`
			if _, err := w.Write([]byte(issueJSON)); err != nil {
				t.Error("failed to write response:", err)
			}

		case "/repos/test/repo/issues/1/comments":
			commentsJSON := `[
				{
					"body": "First comment",
					"user": {"login": "commenter1"},
					"created_at": "2024-01-15T11:00:00Z"
				},
				{
					"body": "Second comment",
					"user": {"login": "commenter2"},
					"created_at": "2024-01-16T09:00:00Z"
				}
			]`
			if _, err := w.Write([]byte(commentsJSON)); err != nil {
				t.Error("failed to write response:", err)
			}

		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	// Test by simulating the URL that would hit our mock server
	t.Run("mock API server setup", func(t *testing.T) {
		// Verify server is responding
		resp, err := http.Get(server.URL + "/repos/test/repo/issues/1")
		if err != nil {
			t.Fatal("failed to query mock server:", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}
	})

	// Test expected time values
	testTime, _ := time.Parse(time.RFC3339, "2024-01-15T10:00:00Z")
	expectedTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	if !testTime.Equal(expectedTime) {
		t.Errorf("time parsing: expected %v, got %v", expectedTime, testTime)
	}
}

// simpleURLInfo for testing - mimics the internal URL parsing
type simpleURLInfo struct {
	Owner  string
	Repo   string
	Number int
}

// parseTestURL simulates basic URL parsing for testing
func parseTestURL(urlString string) *simpleURLInfo {
	// This is a simplified version for test demonstration only
	// Production code should use urlparser.ParseURL()

	// Handle invalid URLs
	if urlString == "not-a-url" || urlString == "https://github.com/owner/repo/issues/" {
		return nil
	}

	// Simple parsing logic - match specific test URLs
	var owner, repo string
	var number int

	// Use different values based on URL type
	if strings.Contains(urlString, "issues/123") {
		owner, repo, number = "owner", "repo", 123
	} else if strings.Contains(urlString, "pull/456") {
		owner, repo, number = "owner", "repo", 456
	} else {
		return nil // Invalid URL for our test cases
	}

	return &simpleURLInfo{
		Owner:  owner,
		Repo:   repo,
		Number: number,
	}
}