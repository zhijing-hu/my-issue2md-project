package markdown_test

import (
	"testing"
	"time"

	"my-issue2md/internal/github"
	"my-issue2md/internal/markdown"
)

// TestFormatAsMarkdown tests the markdown formatting functionality
func TestFormatAsMarkdown(t *testing.T) {
	// Create test data
	testTime, _ := time.Parse("2006-01-02", "2024-01-15")

	issue := &github.Issue{
		Title:     "Test Issue Title",
		Body:      "This is the main issue body with some **markdown** formatting.",
		URL:       "https://github.com/test/test-issue-repo/issues/42",
		Author:    "testuser",
		CreatedAt: testTime,
		Labels:    []string{"bug", "help wanted"},
		Milestone: "v1.0",
		Number:    42,
	}

	comments := []github.Comment{
		{
			Author:    "commenter1",
			Body:      "Great issue! Let me help.",
			CreatedAt: testTime.Add(24 * time.Hour),
		},
		{
			Author:    "commenter2",
			Body:      "I'll work on this tomorrow.",
			CreatedAt: testTime.Add(48 * time.Hour),
		},
	}

	expected := `# Test Issue Title

**URL**: [https://github.com/test/test-issue-repo/issues/42](https://github.com/test/test-issue-repo/issues/42)
**作者**: @testuser
**创建日期**: 2024-01-15
**标签**: bug, help wanted
**里程碑**: v1.0

## 主内容
This is the main issue body with some **markdown** formatting.

## 评论

### @commenter1 - 2024-01-16
Great issue! Let me help.

### @commenter2 - 2024-01-17
I'll work on this tomorrow.

`

	result := markdown.FormatAsMarkdown(issue, comments)

	if result != expected {
		t.Errorf("Markdown output mismatch:")
		t.Errorf("Expected:\n%q", expected)
		t.Errorf("Got:\n%q", result)
	}
}

// TestFormatAsMarkdown_EmptyLabels tests formatting when issue has no labels
func TestFormatAsMarkdown_EmptyLabels(t *testing.T) {
	testTime, _ := time.Parse("2006-01-02", "2024-01-15")

	issue := &github.Issue{
		Title:     "Issue Without Labels",
		Body:      "No labels here.",
		URL:       "https://github.com/test/repo/issues/1",
		Author:    "user",
		CreatedAt: testTime,
		Labels:    []string{},
		Milestone: "",
		Number:    1,
	}

	expected := `# Issue Without Labels

**URL**: [https://github.com/test/repo/issues/1](https://github.com/test/repo/issues/1)
**作者**: @user
**创建日期**: 2024-01-15

## 主内容
No labels here.

`

	result := markdown.FormatAsMarkdown(issue, []github.Comment{})

	if result != expected {
		t.Errorf("Markdown output mismatch:")
		t.Errorf("Expected:\n%q", expected)
		t.Errorf("Got:\n%q", result)
	}
}

// TestFormatAsMarkdown_NoComments tests formatting when issue has no comments
func TestFormatAsMarkdown_NoComments(t *testing.T) {
	testTime, _ := time.Parse("2006-01-02", "2024-01-15")

	issue := &github.Issue{
		Title:     "Issue Without Comments",
		Body:      "This issue has no comments.",
		URL:       "https://github.com/test/repo/issues/2",
		Author:    "user",
		CreatedAt: testTime,
		Labels:    []string{"enhancement"},
		Milestone: "next-release",
		Number:    2,
	}

	expected := `# Issue Without Comments

**URL**: [https://github.com/test/repo/issues/2](https://github.com/test/repo/issues/2)
**作者**: @user
**创建日期**: 2024-01-15
**标签**: enhancement
**里程碑**: next-release

## 主内容
This issue has no comments.

`

	result := markdown.FormatAsMarkdown(issue, []github.Comment{})

	if result != expected {
		t.Errorf("Markdown output mismatch:")
		t.Errorf("Expected:\n%q", expected)
		t.Errorf("Got:\n%q", result)
	}
}

// TestFormatAsMarkdown_MultilineContent tests issue with multiline content
func TestFormatAsMarkdown_MultilineContent(t *testing.T) {
	testTime, _ := time.Parse("2006-01-02", "2024-01-15")

	issue := &github.Issue{
		Title:     "Multiline Issue",
		Body:      "First line.\n\nSecond paragraph.\n\n\n\nLast line.",
		URL:       "https://github.com/test/repo/issues/3",
		Author:    "user",
		CreatedAt: testTime,
		Labels:    []string{},
		Milestone: "",
		Number:    3,
	}

	expected := `# Multiline Issue

**URL**: [https://github.com/test/repo/issues/3](https://github.com/test/repo/issues/3)
**作者**: @user
**创建日期**: 2024-01-15

## 主内容
First line.

Second paragraph.



Last line.

`

	result := markdown.FormatAsMarkdown(issue, []github.Comment{})

	if result != expected {
		t.Errorf("Markdown output mismatch:")
		t.Errorf("Expected:\n%q", expected)
		t.Errorf("Got:\n%q", result)
	}
}
