package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Issue represents a GitHub issue with metadata
type Issue struct {
	Title     string
	Body      string
	URL       string
	Author    string
	CreatedAt time.Time
	Labels    []string
	Milestone string
	Number    int
}

// Comment represents a comment on an issue/PR/discussion
type Comment struct {
	Author    string
	Body      string
	CreatedAt time.Time
}

// GitHub Error Types (following principle 3.1)
var (
	ErrRateLimited = fmt.Errorf("github: rate limited")
	ErrNotFound    = fmt.Errorf("github: resource not found")
	ErrUnauthorized = fmt.Errorf("github: unauthorized")
)

func fetchURL(url string, token string) ([]byte, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("github: create request: %w", err)
	}

	// Add authentication if token provided
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("github: http request: %w", err)
	}
	defer resp.Body.Close()

	// Handle rate limiting
	if resp.StatusCode == http.StatusForbidden {
		return nil, fmt.Errorf("github: %w", ErrRateLimited)
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("github: %w", ErrNotFound)
		}
		if resp.StatusCode == http.StatusUnauthorized {
			return nil, fmt.Errorf("github: %w", ErrUnauthorized)
		}
		return nil, fmt.Errorf("github: unexpected status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("github: read response: %w", err)
	}

	return body, nil
}

// GetIssue fetches an issue and its comments from GitHub API
func GetIssue(urlString string, token string) (*Issue, []Comment, error) {
	// This would use urlparser package once we create it
	// For now, we'll simulate the URL parsing
	parsedURL := parseSimpleURL(urlString)
	if parsedURL == nil {
		return nil, nil, fmt.Errorf("github: invalid URL format")
	}

	// Construct API URLs
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues/%d", parsedURL.Owner, parsedURL.Repo, parsedURL.Number)
	commentsURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues/%d/comments", parsedURL.Owner, parsedURL.Repo, parsedURL.Number)

	// Fetch issue data
	issueData, err := fetchURL(apiURL, token)
	if err != nil {
		return nil, nil, fmt.Errorf("github: fetch issue: %w", err)
	}

	// Parse issue response
	var apiResponse struct {
		Title     string    `json:"title"`
		Body      string    `json:"body"`
		HTMLURL   string    `json:"html_url"`
		User      struct {
			Login string `json:"login"`
		} `json:"user"`
		CreatedAt string    `json:"created_at"`
		Labels    []struct {
			Name string `json:"name"`
		} `json:"labels"`
		Milestone *struct {
			Title string `json:"title"`
		} `json:"milestone"`
		Number int `json:"number"`
	}

	if err := json.Unmarshal(issueData, &apiResponse); err != nil {
		return nil, nil, fmt.Errorf("github: parse issue json: %w", err)
	}

	// Parse created_at time
	createdAt, err := time.Parse(time.RFC3339, apiResponse.CreatedAt)
	if err != nil {
		return nil, nil, fmt.Errorf("github: parse created_at: %w", err)
	}

	// Build issue struct
	issue := &Issue{
		Title:     apiResponse.Title,
		Body:      apiResponse.Body,
		URL:       apiResponse.HTMLURL,
		Author:    apiResponse.User.Login,
		CreatedAt: createdAt,
		Number:    apiResponse.Number,
	}

	// Extract labels
	for _, label := range apiResponse.Labels {
		issue.Labels = append(issue.Labels, label.Name)
	}

	// Extract milestone if present
	if apiResponse.Milestone != nil {
		issue.Milestone = apiResponse.Milestone.Title
	}

	// Fetch comments
	commentsData, err := fetchURL(commentsURL, token)
	if err != nil {
		return nil, nil, fmt.Errorf("github: fetch comments: %w", err)
	}

	var comments []Comment
	var commentsResponse []struct {
		Body      string `json:"body"`
		User      struct {
			Login string `json:"login"`
		} `json:"user"`
		CreatedAt string `json:"created_at"`
	}

	if err := json.Unmarshal(commentsData, &commentsResponse); err != nil {
		return nil, nil, fmt.Errorf("github: parse comments json: %w", err)
	}

	// Build comments array
	for _, c := range commentsResponse {
		createdAt, err := time.Parse(time.RFC3339, c.CreatedAt)
		if err != nil {
			// Log and skip individual comment with bad time format
			continue
		}
		comments = append(comments, Comment{
			Author:    c.User.Login,
			Body:      c.Body,
			CreatedAt: createdAt,
		})
	}

	return issue, comments, nil
}

// simpleURL represents a simplified parsed URL for initial implementation
type simpleURL struct {
	Owner string
	Repo  string
	Number int
}

// parseSimpleURL provides basic URL parsing until we implement urlparser
func parseSimpleURL(urlString string) *simpleURL {
	// Simple parsing for basic GitHub URLs
	// This will be replaced by the urlparser package
	parts := strings.Split(urlString, "/")
	if len(parts) < 7 || !strings.Contains(urlString, "github.com") {
		return nil
	}

	// Extract owner, repo, and issue number
	owner := parts[3]
	repo := parts[4]
	number := 0

	// Handle both issues and pull requests
	if parts[5] == "issues" || parts[5] == "pull" {
		if _, err := fmt.Sscanf(parts[6], "%d", &number); err != nil || number == 0 {
			return nil
		}
	} else if parts[5] == "discussions" {
		if _, err := fmt.Sscanf(parts[6], "%d", &number); err != nil || number == 0 {
			return nil
		}
	}

	return &simpleURL{
		Owner: owner,
		Repo:  repo,
		Number: number,
	}
}