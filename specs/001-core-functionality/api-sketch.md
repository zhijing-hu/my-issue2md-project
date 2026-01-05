# API Sketch - 内部包接口说明

此文档提供了 `internal/markdown` 和 `internal/github` 包对外暴露的主要接口说明，作为后续开发的参考。

## 目录
- [internal/github 包](#internalgithub-包)
  - [类型](#类型)
  - [函数](#函数)
- [internal/markdown 包](#internalmarkdown-包)
  - [函数](#函数-1)
- [使用示例](#使用示例)
- [API 稳定性保证](#api-稳定性保证)
- [开发指南](#开发指南)

## internal/github 包

`internal/github` 包提供了对 GitHub REST API v3 的访问，用于获取 issue、pull request 和 discussion 及其关联评论。

```go
import "my-issue2md/internal/github"
```

### 类型

#### `Issue`
```go
type Issue struct {
    Title     string    // Issue/PR/Discussion 标题
    Body      string    // 主内容正文（markdown 格式）
    URL       string    // 原始 GitHub URL
    Author    string    // 创建者用户名
    CreatedAt time.Time // 创建时间戳
    Labels    []string  // 关联标签
    Milestone string    // 里程碑名称（如有设置）
    Number    int       // Issue/PR/Discussion 编号
}
```

**用途：** 存储 GitHub issue、pull request 或 discussion 的完整元数据。

#### `Comment`
```go
type Comment struct {
    Author    string    // 评论者用户名
    Body      string    // 评论内容（markdown 格式）
    CreatedAt time.Time // 评论时间戳
}
```

**用途：** 表示 issue、pull request 或 discussion 上的个别评论。

### 函数

#### `GetIssue()`
```go
func GetIssue(urlString string, token string) (*Issue, []Comment, error)
```

**参数：**
- `urlString` (string)：完整 GitHub URL（支持 issues、PRs、discussions）
- `token` (string)：可选的 GitHub 委托令牌（用于私有仓库）

**返回**：
- `*Issue`：解析后的 issue 元数据
- `[]Comment`：所有评论（按时间顺序）
- `error`：任何步骤失败时非 nil（URL 解析、网络、API、认证）

**错误类型：**
- `ErrNotFound`：资源不存在（HTTP 404）
- `ErrUnauthorized`：需要认证或认证失败（HTTP 401）
- `ErrRateLimited`：GitHub API 速率限制（HTTP 403）
- `ErrNetwork`：网络/HTTP 错误

**用法示例：**
```go
issue, comments, err := github.GetIssue(
    "https://github.com/owner/repo/issues/123",
    os.Getenv("GITHUB_TOKEN")
)
if err != nil {
    // 处理错误
}
```

## internal/markdown 包

`internal/markdown` 包提供格式化功能，将 GitHub 数据转换为规范的 markdown 输出。

```go
import "my-issue2md/internal/markdown"
```

### 函数

#### `FormatAsMarkdown()`
```go
func FormatAsMarkdown(issue *github.Issue, comments []github.Comment) string
```

**参数：**
- `issue` (*github.Issue)：来自 `github.GetIssue()` 的 issue 元数据
- `comments` ([]github.Comment)：来自 `github.GetIssue()` 的评论数据

**返回。**
- `string`：完整的 markdown 文档（符合规范格式）

**Markdown 结构：**
```markdown
# [Issue 标题]

**URL**: [original_url](original_url)
**作者**: @username
**创建日期**: YYYY-MM-DD
**标签**: tag1, tag2
**里程碑**: milestone_name

## 主内容
[Issue 主内容]

## 评论

### @commenter1 - YYYY-MM-DD
[评论内容]

### @commenter2 - YYYY-MM-DD
[评论内容]
```

**用法示例：**
```go
markdownContent := markdown.FormatAsMarkdown(issue, comments)
err := os.WriteFile("output.md", []byte(markdownContent), 0644)
```

## 使用示例

### 完整工作流
```go
package main

import (
    "fmt"
    "os"
    "my-issue2md/internal/github"
    "my-issue2md/internal/markdown"
)

func main() {
    // 从 GitHub 获取数据
    issue, comments, err := github.GetIssue(
        "https://github.com/owner/repo/issues/42",
        os.Getenv("GITHUB_TOKEN")
    )
    if err != nil {
        fmt.Printf("获取 issue 时出错: %v\n", err)
        os.Exit(determineExitCode(err))
    }

    // 生成 markdown
    markdownContent := markdown.FormatAsMarkdown(issue, comments)

    // 写入文件
    err = os.WriteFile("issue_42.md", []byte(markdownContent), 0644)
    if err != nil {
        fmt.Printf("写入文件时出错: %v\n", err)
        os.Exit(4) // 文件系统错误
    }

    fmt.Println("成功生成 markdown 输出")
}
```

### 错误处理
```go
import (
    "my-issue2md/internal/urlparser"
)

func determineExitCode(err error) int {
    if urlparser.IsValidationError(err) {
        return 1 // URL 格式错误
    }
    if github.IsNetworkOrAPIError(err) {  // 假设的辅助函数
        return 2 // 网络/API 错误
    }
    if github.IsAuthError(err) {            // 假设的辅助函数
        return 3 // 认证错误
    }
    return 4 // 文件系统/其他错误
}
```

## API 稳定性保证

### 稳定级别
- **稳定** 🟢：此文档中定义的公共接口在当前主版本中稳定
- **内部** 🟡：未记录的函数/类型可能会在未通知的情况下改变
- **实验性** 🟢：不适用（所有记录的接口都被视为稳定）

### 版本控制政策
- **主版本内** 不引入破坏性改变
- **向后兼容性**：为记录的类型/函数维护
- **废弃过程**：旧接口会用 `DEPRECATED` 注释标记后再移除

## 开发指南

### 添加新功能
- 通过新字段扩展现有类型（不破坏现有代码）
- 添加带有清晰文档的新函数
- 保持现有函数签名不变
- 遵循相同的错误包装模式

### 错误处理模式
- 使用 `fmt.Errorf("...: %w", err)` 包装错误
- 将自定义错误类型定义为包级变量
- 在函数文档中记录错误类型

### 测试
- 所有公共函数必须有全面测试
- 遵循表驱动测试模式
- 测试成功和错误两种情况
- 使用 `errors.Is()` 和 `errors.As()` 验证错误类型

## 接口演进

### 未来考虑
以下是可能的未来扩展方向，但**不属于**当前 API：

```go
// 当前未实现 - 未来可能性
// func GetPR(url string, token string) (*Issue, []Comment, error)
// func GetDiscussion(url string, token string) (*Issue, []Comment, error)
// func GetMultiple(urls []string, token string) ([]Issue, [][]Comment, error)
```

**当前设计决策**：单个 `GetIssue()` 函数内部根据 URL 解析处理所有资源类型（issues、PRs、discussions），保持更简单的 API 表面。

## 总结

此 API 文档定义了其他组件可以依赖的稳定公共接口：

- **3 种公共类型**：`Issue`、`Comment`
- **2 种公共函数**：`GetIssue()`、`FormatAsMarkdown()`
- **4 种错误类型**：`ErrNotFound`、`ErrUnauthorized`、`ErrRateLimited`，加上网络错误

所有接口遵循项目宪法原则：简洁性、显式错误处理、无全局状态。