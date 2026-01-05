package urlparser_test

import (
	"testing"

	"my-issue2md/internal/urlparser"
)

// TestParseURL tests the ParseURL function with various URL formats
func TestParseURL(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		want        *urlparser.URLInfo
		wantErr     bool
		wantErrType error
	}{
		{
			name:    "valid issue URL",
			input:   "https://github.com/owner/repo/issues/123",
			want: &urlparser.URLInfo{
				Owner:       "owner",
				Repo:        "repo",
				Number:      123,
				ResourceType: "issue",
			},
			wantErr: false,
		},
		{
			name:    "valid PR URL",
			input:   "https://github.com/owner/repo/pull/456",
			want: &urlparser.URLInfo{
				Owner:       "owner",
				Repo:        "repo",
				Number:      456,
				ResourceType: "pr",
			},
			wantErr: false,
		},
		{
			name:    "valid discussion URL",
			input:   "https://github.com/owner/repo/discussions/789",
			want: &urlparser.URLInfo{
				Owner:       "owner",
				Repo:        "repo",
				Number:      789,
				ResourceType: "discussion",
			},
			wantErr: false,
		},
		{
			name:        "non-GitHub URL",
			input:       "https://gitlab.com/owner/repo/issues/1",
			wantErr:     true,
			wantErrType: urlparser.ErrNotGitHubURL,
		},
		{
			name:        "invalid URL format",
			input:       "not-a-url",
			wantErr:     true,
			wantErrType: urlparser.ErrNotGitHubURL,
		},
		{
			name:        "invalid issue URL format",
			input:       "https://github.com/owner/repo/issues/abc",
			wantErr:     true,
			wantErrType: urlparser.ErrInvalidGitHubURL,
		},
		{
			name:        "missing issue number",
			input:       "https://github.com/owner/repo/issues/",
			wantErr:     true,
			wantErrType: urlparser.ErrInvalidGitHubURL,
		},
		{
			name:        "URL with path but no issue",
			input:       "https://github.com/owner/repo",
			wantErr:     true,
			wantErrType: urlparser.ErrInvalidGitHubURL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := urlparser.ParseURL(tt.input)

			// Check error conditions
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
					return
				}

				// Check specific error type if provided
				if tt.wantErrType != nil {
					t.Logf("error: %v", err)
				}

				return
			}

			// Check for unexpected errors
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// Compare URL info
			if got == nil {
				t.Errorf("expected URLInfo, got nil")
				return
			}

			// Compare each field
			if got.Owner != tt.want.Owner {
				t.Errorf("Owner: got %q, want %q", got.Owner, tt.want.Owner)
			}
			if got.Repo != tt.want.Repo {
				t.Errorf("Repo: got %q, want %q", got.Repo, tt.want.Repo)
			}
			if got.Number != tt.want.Number {
				t.Errorf("Number: got %d, want %d", got.Number, tt.want.Number)
			}
			if got.ResourceType != tt.want.ResourceType {
				t.Errorf("ResourceType: got %q, want %q", got.ResourceType, tt.want.ResourceType)
			}
		})
	}
}

// TestGetIssueAPIURL tests the GetIssueAPIURL method
func TestGetIssueAPIURL(t *testing.T) {
	tests := []struct {
		name     string
		urlInfo  urlparser.URLInfo
		wantURL  string
	}{
		{
			name: "issue API URL",
			urlInfo: urlparser.URLInfo{
				Owner:       "owner",
				Repo:        "repo",
				Number:      123,
				ResourceType: "issue",
			},
			wantURL: "https://api.github.com/repos/owner/repo/issues/123",
		},
		{
			name: "PR API URL (uses issues endpoint)",
			urlInfo: urlparser.URLInfo{
				Owner:       "owner",
				Repo:        "repo",
				Number:      456,
				ResourceType: "pr",
			},
			wantURL: "https://api.github.com/repos/owner/repo/issues/456",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotURL := tt.urlInfo.GetIssueAPIURL()
			if gotURL != tt.wantURL {
				t.Errorf("got %q, want %q", gotURL, tt.wantURL)
			}
		})
	}
}

// TestGetCommentsAPIURL tests the GetCommentsAPIURL method
func TestGetCommentsAPIURL(t *testing.T) {
	tests := []struct {
		name     string
		urlInfo  urlparser.URLInfo
		wantURL  string
	}{
		{
			name: "issue comments URL",
			urlInfo: urlparser.URLInfo{
				Owner:       "owner",
				Repo:        "repo",
				Number:      123,
				ResourceType: "issue",
			},
			wantURL: "https://api.github.com/repos/owner/repo/issues/123/comments",
		},
		{
			name: "PR comments URL",
			urlInfo: urlparser.URLInfo{
				Owner:       "owner",
				Repo:        "repo",
				Number:      456,
				ResourceType: "pr",
			},
			wantURL: "https://api.github.com/repos/owner/repo/pulls/456/comments",
		},
		{
			name: "discussion comments URL",
			urlInfo: urlparser.URLInfo{
				Owner:       "owner",
				Repo:        "repo",
				Number:      789,
				ResourceType: "discussion",
			},
			wantURL: "https://api.github.com/repos/owner/repo/discussions/789/comments",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotURL := tt.urlInfo.GetCommentsAPIURL()
			if gotURL != tt.wantURL {
				t.Errorf("got %q, want %q", gotURL, tt.wantURL)
			}
		})
	}
}