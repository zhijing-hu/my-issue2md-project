# issue2md - 技术规格说明

## 概述
`issue2md` 是一个命令行工具，能够将 GitHub Issue、Pull Request 或 Discussion 的 URL 转换为 Markdown 文件，内容包含完整的主内容、评论及元数据。

## 核心功能

### 输入
- **接受格式**：单个 GitHub URL（Issue/PR/Discussion）
- **URL类型**：仅支持 GitHub URLs （https://github.com/owner/repo/issues/number、prs、discussions）
- **参数形式**：直接跟在命令后面，无需-url前缀，如 `issue2md <url>`
- **认证**：可选 GITHUB_TOKEN 环境变量，用于访问私有仓库

### 处理过程
- **内容提取**：
  - 主内容（标题、正文）
  - 所有评论，按时间顺序排列
  - 元数据：作者、创建日期、标签、里程碑
  - 原始URL引用

### 输出
- **格式**：单个 `.md` 文件
- **结构**：
  ```markdown
  # [Issue标题]

  **URL**: [original_url](original_url)
  **作者**: @username
  **创建日期**: YYYY-MM-DD
  **标签**: tag1, tag2
  **里程碑**: milestone_name (如存在)

  ## 主内容
  [主内容在此]

  ## 评论

  ### @comment_author1 - YYYY-MM-DD
  [评论内容]

  ### @comment_author2 - YYYY-MM-DD
  [评论内容]
  ```

## 命令行接口

### 用法
```bash
# 基本用法
issue2md <github_url>

# 可以指定输出文件名
issue2md <github_url> -o <output_file.md>
```

### 参数
- `<github_url>`（必需）：GitHub Issue/PR/Discussion URL，直接作为命令参数
- `-o, --output`（可选）：输出 Markdown 文件路径，默认为 `issue_<id>.md`

## 错误处理

### 错误场景
1. **URL格式错误**：非法的GitHub URL结构
2. **网络错误**：GitHub API不可用
3. **认证需要**：私有仓库但没有Token
4. **资源不存在**：Issue/PR/Discussion不存在
5. **API速率限制**：GitHub API速率限制
6. **文件系统错误**：无法写入输出文件

### 错误响应
- **日志记录**：详细错误信息保存到 `issue2md_<timestamp>.log`
- **退出码**：
  - 0：成功
  - 1：URL格式错误
  - 2：网络/GitHub API错误
  - 3：访问/授权错误
  - 4：输出文件系统错误
- **标准输出**：成功时简短的确认消息

## 技术规格

### 平台支持
- **GitHub API**：使用 REST API v3（https://docs.github.com/en/rest）
- **HTTP客户端**：标准 Go net/http 包
- **认证**：GITHUB_TOKEN 环境变量用于私有仓库

### 速率限制
- **政策**：遇到速率限制（403）立即退出与适当错误
- **重试**：不实现自动重试机制

### 性能
- **要求**：仅关注功能正确性，无性能优化需求
- **并发**：单线程、顺序处理

## 实现约束

### 不实现功能
- 批量处理（多个URLs）
- 其他平台（GitLab，Bitbucket）
- 多种输出格式
- 模板系统
- 插件架构
- 配置文件
- 速率限制重试机制

### 开发环境
- **语言**：Go 1.24+
- **错误处理**：显式错误检查，使用 fmt.Errorf("...: %w", err) 模式
- **无全局变量**：通过函数参数显式注入依赖

## 使用示例

### 公开仓库
```bash
# 基本用法 - 输出默认文件名
export GITHUB_TOKEN=""  # 公开仓库无需token
issue2md https://github.com/owner/repo/issues/123

# 指定输出文件名
issue2md https://github.com/owner/repo/issues/123 -o issue_123.md
```

### 私有仓库
```bash
# 需要token
export GITHUB_TOKEN="ghp_yourtokenhere"
issue2md https://github.com/owner/private-repo/issues/456

# 指定输出
issue2md https://github.com/owner/private-repo/issues/456 -o private_issue.md
```

## 未来考虑（暂不实现）
- 批量处理支持
- 其它输出格式
- 模板定制
- 插件体系结构
- 配置文件支持
- 速率限制重试机制