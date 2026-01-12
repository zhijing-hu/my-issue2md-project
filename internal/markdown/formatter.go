package markdown

import (
	"fmt"
	"strings"

	"my-issue2md/internal/github"
)

// FormatAsMarkdown generates markdown content from GitHub issue data
func FormatAsMarkdown(issue *github.Issue, comments []github.Comment) string {
	var sb strings.Builder

	// Header with metadata
	sb.WriteString(fmt.Sprintf("# %s\n\n", issue.Title))
	sb.WriteString(fmt.Sprintf("**URL**: [%s](%s)\n", issue.URL, issue.URL))
	sb.WriteString(fmt.Sprintf("**作者**: @%s\n", issue.Author))
	sb.WriteString(fmt.Sprintf("**创建日期**: %s\n", issue.CreatedAt.Format("2006-01-02")))

	if len(issue.Labels) > 0 {
		sb.WriteString(fmt.Sprintf("**标签**: %s\n", strings.Join(issue.Labels, ", ")))
	}

	if issue.Milestone != "" {
		sb.WriteString(fmt.Sprintf("**里程碑**: %s\n", issue.Milestone))
	}

	// Main content
	sb.WriteString("\n## 主内容\n")
	sb.WriteString(fmt.Sprintf("%s\n\n", issue.Body))

	// Comments
	if len(comments) > 0 {
		sb.WriteString("## 评论\n\n")

		for _, comment := range comments {
			sb.WriteString(fmt.Sprintf("### @%s - %s\n", comment.Author, comment.CreatedAt.Format("2006-01-02")))
			sb.WriteString(fmt.Sprintf("%s\n\n", comment.Body))
		}
	}

	return sb.String()
}