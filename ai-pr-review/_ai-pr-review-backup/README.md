# AI PR Review Assistant

本地运行的 Go CLI 工具，自动获取 GitHub Pull Request 的 diff，调用 AI 进行代码审查，并将结果保存为本地 JSON 文件。

## 功能

1. 解析 GitHub PR URL
2. 获取 PR 标题、作者、变更文件、commits、raw diff
3. 将 diff 解析为结构化 chunk（按 hunk 切分）
4. 对每个 chunk 调用 OpenAI 进行风险分析
5. 输出 review 结果到 `output/pr_<PR号>_review.json`

## 流程

```
PR URL → GitHub API → Diff Parser → AI Review (per chunk) → JSON Output
```

## 前置条件

- Go 1.23+
- GitHub Personal Access Token（需要 `repo` 权限读取 PR）
- OpenAI API Key

## 配置

复制配置模板并填入你的密钥：

```bash
cp config/config.yaml.example config/config.yaml
```

编辑 `config/config.yaml`：

```yaml
github_token: "ghp_xxx"
openai_api_key: "sk-xxx"
model: "gpt-4o"
```

## 运行

```bash
cd ai-pr-review
go run ./cmd
```

按提示输入 PR URL，例如：

```
https://github.com/org/repo/pull/123
```

完成后输出：

```
Review completed.
Result saved to:
output/pr_123_review.json
```

## 输出示例

```json
{
  "pr_title": "fix user login",
  "pr_number": 123,
  "issues": [
    {
      "file": "service/user.go",
      "line": 45,
      "severity": "high",
      "message": "goroutine may leak",
      "suggestion": "add context cancellation"
    }
  ]
}
```

## 项目结构

```
ai-pr-review/
├── cmd/main.go              # CLI 入口
├── config/config.yaml       # 本地配置（不提交到 git）
├── internal/
│   ├── github/              # PR URL 解析与 GitHub API
│   ├── parser/              # Diff 解析、语言/函数识别
│   ├── ai/                  # LLM 接口与 OpenAI 实现
│   ├── review/              # Review 编排引擎
│   ├── output/              # JSON 输出
│   └── model/               # 共享数据结构
├── output/                  # Review 结果目录
├── go.mod
└── README.md
```

## 模块说明

| 模块 | 职责 |
|------|------|
| `internal/github` | 解析 PR URL，拉取 PR 信息、diff、commits |
| `internal/parser` | 解析 unified diff hunk，识别语言与 Go 函数名 |
| `internal/ai` | `LLM` 接口 + OpenAI 实现，prompt 独立管理 |
| `internal/review` | 遍历 chunk 调用 AI，聚合 issues |
| `internal/output` | 写入 JSON 文件 |

## 扩展点

- **LLM 接口**：`internal/ai.LLM` 可替换为 Claude、DeepSeek 等实现
- **FunctionDetector**：`internal/parser.FunctionDetector` 可扩展 Python/JS 函数识别

## 测试

```bash
go test ./...
```

## 注意事项

- 大型 PR 的 diff 可能较大，审查会按 chunk 多次调用 OpenAI API
- `config/config.yaml` 含敏感信息，已加入 `.gitignore`
- 需要网络访问 GitHub API 与 OpenAI API
