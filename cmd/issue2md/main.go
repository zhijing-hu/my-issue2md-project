package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"my-issue2md/internal/github"
	"my-issue2md/internal/markdown"
	"my-issue2md/internal/urlparser"
)

func main() {
	// Parse CLI flags
	outputFile := fmt.Sprintf("issue_%d.md", time.Now().Unix())
	flag.StringVar(&outputFile, "o", outputFile, "output file path")
	flag.StringVar(&outputFile, "output", outputFile, "output file path (alternative)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <github_url> [-o output_file.md]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Example: %s https://github.com/owner/repo/issues/123\n", os.Args[0])
	}
	flag.Parse()

	// Validation - expect exactly one argument (the URL)
	args := flag.Args()
	if len(args) != 1 {
		flag.Usage()
		os.Exit(1)
	}
	githubURL := strings.TrimSpace(args[0])

	// Validate the URL first
	_, err := urlparser.ParseURL(githubURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		logError(err, githubURL)
		os.Exit(1)
	}

	// Get auth token from environment
	token := os.Getenv("GITHUB_TOKEN")

	// Execute the full pipeline
	issue, comments, err := processIssue(githubURL, token)
	if err != nil {
		// Log error and exit with appropriate code
		exitCode := determineExitCode(err)
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		logError(err, githubURL)
		os.Exit(exitCode)
	}

	// Generate markdown content
	markdownContent := markdown.FormatAsMarkdown(issue, comments)

	// Write to file
	err = os.WriteFile(outputFile, []byte(markdownContent), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to write output file: %v\n", err)
		logError(fmt.Errorf("write output: %w", err), githubURL)
		os.Exit(4)
	}

	fmt.Printf("Successfully generated markdown to %s\n", outputFile)
}

// processIssue orchestrates all the internal packages
func processIssue(urlString, token string) (*github.Issue, []github.Comment, error) {
	// Fetch from GitHub based on resource type
	issue, comments, err := github.GetIssue(urlString, token)
	if err != nil {
		return nil, nil, fmt.Errorf("process: fetch from github: %w", err)
	}

	return issue, comments, nil
}

// logError handles error logging to file
func logError(err error, url string) {
	timestamp := time.Now().Format("20060102_150405")
	logFile := fmt.Sprintf("issue2md_%s.log", timestamp)
	logMessage := fmt.Sprintf("%s: %v\nURL: %s\n", timestamp, err, url)

	// Write to log file
	f, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		// If we can't log the error, print to stderr
		fmt.Fprintf(os.Stderr, "Failed to write log: %v\n", err)
		return
	}
	defer f.Close()

	if _, err := f.WriteString(logMessage); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write to log file: %v\n", err)
	}
}

// determineExitCode maps errors to exit codes per specification
func determineExitCode(err error) int {
	// Check for URL validation errors first
	if urlparser.IsValidationError(err) {
		return 1 // URL format error
	}

	// Check for GitHub API/auth errors
	if err == github.ErrRateLimited {
		return 2 // Rate limit error
	}
	if err == github.ErrNotFound {
		return 2 // Network/GitHub API error
	}
	if err == github.ErrUnauthorized {
		return 3 // Access/auth error
	}

	// Check for network errors by looking at the error chain
	if strings.Contains(err.Error(), "network") ||
	   strings.Contains(err.Error(), "http request") ||
	   strings.Contains(err.Error(), "create request") {
		return 2 // Network/GitHub API error
	}

	// Default: file system or other errors
	return 4
}