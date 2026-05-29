# AI PR Review Assistant

本地运行的 Go 工具，自动获取 GitHub Pull Request 的 diff，调用 AI 进行代码审查，并将结果保存为 JSON 文件与 MySQL 数据库。支持 CLI 交互模式与 REST API 两种入口。

## 功能

1. 解析 GitHub PR URL
2. 获取 PR 标题、作者、变更文件、commits、raw diff
3. 将 diff 解析为结构化 chunk（按 hunk 切分）
4. 对每个 chunk 调用 LLM 进行风险分析
5. 输出 review 结果到 `output/` 目录并写入 MySQL
6. 通过 REST API 提交审查、查询结果与历史记录

## 流程

```
PR URL → GitHub API → Diff Parser → AI Review (per chunk) → JSON Output + MySQL
```

## 前置条件

- Go 1.23+
- MySQL 5.7+ / 8.0+
- GitHub Personal Access Token（需要 `repo` 权限读取 PR）
- OpenAI 兼容 API Key（如 DeepSeek）

## 配置

复制配置模板：

```bash
cp config/config.yaml.example config/config.yaml
```

设置环境变量（PowerShell 示例）：

```powershell
$env:GITHUB_TOKEN = "ghp_xxx"
$env:OPENAI_API_KEY = "sk-xxx"
$env:MODEL = "deepseek-chat"
```

`config/config.yaml` 通过 `${VAR_NAME}` 引用环境变量，例如：

```yaml
github_token: ${GITHUB_TOKEN}
openai_api_key: ${OPENAI_API_KEY}
model: ${MODEL}
api_base: "https://api.deepseek.com/v1"
output_dir: "output"
prompt_dir: "prompts"

server:
  host: "0.0.0.0"
  port: 8080
```

`api_base` 可直接写死；敏感项请使用环境变量引用，勿写入真实密钥。

## 运行

### CLI 模式

```bash
cd ai-pr-review
go run ./cmd/cli
```

按提示输入 PR URL，例如：

```
https://github.com/org/repo/pull/123
```

### API 模式

```bash
go run ./cmd/api
```

或使用 Makefile：

```bash
make run-api
```

服务默认监听 `http://0.0.0.0:8080`。

## REST API

完整 API 文档见 [api/openapi.yaml](api/openapi.yaml)。

| 方法 | 路径 | 说明 |
|------|------|------|
| `GET` | `/health` | 健康检查 |
| `POST` | `/api/v1/reviews` | 提交 PR URL 触发审查 |
| `GET` | `/api/v1/reviews/:id` | 按 ID 查看审查结果 |
| `GET` | `/api/v1/reviews` | 分页查询历史记录 |

### 示例

提交审查：

```bash
curl -X POST http://localhost:8080/api/v1/reviews \
  -H "Content-Type: application/json" \
  -d '{"pr_url":"https://github.com/org/repo/pull/123"}'
```

查询历史：

```bash
curl "http://localhost:8080/api/v1/reviews?page=1&limit=20"
```

按 ID 查看：

```bash
curl http://localhost:8080/api/v1/reviews/1
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
├── cmd/
│   ├── cli/main.go              # CLI 入口
│   └── api/main.go              # HTTP API 入口
├── internal/
│   ├── handler/                 # HTTP 控制器
│   ├── service/                 # 业务编排层
│   ├── repository/              # 数据访问层
│   ├── middleware/              # 中间件
│   ├── model/                   # 领域数据结构
│   ├── github/                  # GitHub API 客户端
│   ├── parser/                  # Diff 解析
│   ├── prompt/                  # 提示词加载与渲染
│   ├── ai/                      # LLM 接口
│   │   └── langchain/           # langchaingo 实现
│   ├── review/                  # Review 引擎
│   └── output/                  # JSON 文件输出
├── prompts/                     # 提示词模板（与代码解耦）
│   └── review/
│       ├── system.txt
│       └── user.tmpl
├── pkg/validator/               # 输入校验
├── api/
│   ├── openapi.yaml             # API 文档
│   └── routes.go                # 路由注册
├── config/
│   ├── config.yaml              # 本地配置（不提交到 git）
│   ├── config.yaml.example
│   └── config.go                # 配置解析
├── migrations/                  # 数据库迁移
├── scripts/build.sh             # 构建脚本
├── test/api_test.go             # API 集成测试
├── Makefile
├── output/                      # Review 结果目录
├── go.mod
└── README.md
```

## 模块说明

| 模块 | 职责 |
|------|------|
| `internal/handler` | HTTP 请求处理，参数校验，JSON 响应 |
| `internal/service` | 编排 GitHub 拉取、diff 解析、AI 审查、存储 |
| `internal/repository` | MySQL 连接与 review 数据 CRUD |
| `internal/github` | 解析 PR URL，拉取 PR 信息、diff、commits |
| `internal/parser` | 解析 unified diff hunk，识别语言与 Go 函数名 |
| `internal/prompt` | 从 `prompts/` 加载模板并渲染审查提示词 |
| `internal/ai` | `LLM` 接口定义 |
| `internal/ai/langchain` | 基于 langchaingo 的 LLM 实现 |
| `prompts/review` | 系统与用户提示词模板（可直接编辑，无需改代码） |
| `internal/review` | 遍历 chunk 调用 AI，聚合 issues |
| `internal/output` | 写入 JSON 文件 |

## 构建与测试

```bash
make test          # 运行全部测试
make build         # 编译 bin/api 和 bin/cli
make run-cli       # 启动 CLI
make run-api       # 启动 API 服务
```

## 扩展点

- **LLM 接口**：`internal/ai.LLM` 由 `internal/ai/langchain`（langchaingo）实现，可替换为其他 provider
- **提示词**：编辑 `prompts/review/` 下模板即可调整审查策略，可通过 `prompt_dir` 配置切换目录
- **FunctionDetector**：`internal/parser.FunctionDetector` 可扩展 Python/JS 函数识别

## 注意事项

- 大型 PR 的 diff 可能较大，审查会按 chunk 多次调用 LLM API；API 模式当前为同步响应
- 密钥通过环境变量注入，`config.yaml` 中用 `${GITHUB_TOKEN}` 等形式引用，勿写入真实密钥
- 需要网络访问 GitHub API 与 LLM API
- API 服务启动时会自动执行 `migrations/001_create_reviews.sql`
