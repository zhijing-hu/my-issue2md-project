package urlparser

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
)

// URLInfo represents parsed GitHub URL information
type URLInfo struct {
	Owner       string
	Repo        string
	Number      int
	ResourceType string // "issue", "pr", "discussion"
}

// GitHub URL Error Types
var (
	ErrNotGitHubURL    = fmt.Errorf("urlparser: not a GitHub URL")
	ErrInvalidGitHubURL = fmt.Errorf("urlparser: invalid GitHub issue/PR/discussion URL")
)

// GitHub URL Patterns
var (
	issuePattern    = regexp.MustCompile(`^https://github\.com/([^/]+)/([^/]+)/issues/(\d+)$`)
	prPattern       = regexp.MustCompile(`^https://github\.com/([^/]+)/([^/]+)/pull/(\d+)$`)
	discussionPattern = regexp.MustCompile(`^https://github\.com/([^/]+)/([^/]+)/discussions/(\d+)$`)
)

// ParseURL parses and validates a GitHub URL
func ParseURL(rawURL string) (*URLInfo, error) {
	// Basic URL validation
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("urlparser: parse url: %w", err)
	}

	// Validate it's a GitHub URL
	if parsedURL.Host != "github.com" {
		return nil, fmt.Errorf("urlparser: %w", ErrNotGitHubURL)
	}

	// Match against patterns
	if matches := issuePattern.FindStringSubmatch(rawURL); matches != nil {
		return &URLInfo{
			Owner:       matches[1],
			Repo:        matches[2],
			Number:      atoi(matches[3]),
			ResourceType: "issue",
		}, nil
	}

	if matches := prPattern.FindStringSubmatch(rawURL); matches != nil {
		return &URLInfo{
			Owner:       matches[1],
			Repo:        matches[2],
			Number:      atoi(matches[3]),
			ResourceType: "pr",
		}, nil
	}

	if matches := discussionPattern.FindStringSubmatch(rawURL); matches != nil {
		return &URLInfo{
			Owner:       matches[1],
			Repo:        matches[2],
			Number:      atoi(matches[3]),
			ResourceType: "discussion",
		}, nil
	}

	return nil, fmt.Errorf("urlparser: %w", ErrInvalidGitHubURL)
}

// atoi converts string to int with error handling
func atoi(s string) int {
	num, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return num
}

// GetIssueAPIURL returns the GitHub API URL for the issue/PR
func (u *URLInfo) GetIssueAPIURL() string {
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/issues/%d", u.Owner, u.Repo, u.Number)
}

// GetCommentsAPIURL returns the GitHub API URL for comments
func (u *URLInfo) GetCommentsAPIURL() string {
	// PRs and Issues share the same comments endpoint
	if u.ResourceType == "pr" {
		return fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls/%d/comments", u.Owner, u.Repo, u.Number)
	}
	if u.ResourceType == "issue" {
		return fmt.Sprintf("https://api.github.com/repos/%s/%s/issues/%d/comments", u.Owner, u.Repo, u.Number)
	}
	// For discussions
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/discussions/%d/comments", u.Owner, u.Repo, u.Number)
}

// IsValidationError checks if error is a URL validation error
func IsValidationError(err error) bool {
	return err == ErrNotGitHubURL || err == ErrInvalidGitHubURL
}