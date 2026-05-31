# PR-check-r — AI PR Review Assistant

基于 AI 的 GitHub Pull Request 自动代码审查工具，包含 Go 后端与 React Web 前端。

---

## 仓库分支说明

开发过程中，代码一直直接提交到 **`main`** 分支，按天迭代，忘记同步走 Pull Request 流程。后续pr补充在master分支

**`master`** 分支上的 PR 是在 **2025 年 5 月 31 日** 集中补提的，用于将已有改动按功能拆成可 review 的 PR 记录；其内容与 `main` 上对应功能一致，属于 **`main` 分支日常提交的补充说明**，而非独立开发线。

| 分支 | 提交方式 | 说明 |
|------|----------|------|
| `main` | 每日直接 commit | 实际开发主线，包含完整功能演进 |
| `master` | 5.31 日通过 PR 合并 | 对 `main` 已有改动的 PR 化归档，便于 code review 展示 |

查看完整功能请以 **`main`** 为准；`master` 上的 PR 编号（如 `#14`、`#19` 等）对应各功能模块的拆分记录。

---

## 第三方库与框架

### 后端（`ai-pr-review`）

| 类别 | 名称 | 用途 |
|------|------|------|
| 语言 / 运行时 | Go 1.25+ | 后端实现语言 |
| Web 框架 | [Gin](https://github.com/gin-gonic/gin) | REST API HTTP 服务 |
| GitHub 客户端 | [go-github](https://github.com/google/go-github) | 拉取 PR 信息、diff、提交 review / merge |
| LLM 集成 | [langchaingo](https://github.com/tmc/langchaingo) | 调用 OpenAI 兼容 API 进行代码审查 |
| 数据库驱动 | [go-sql-driver/mysql](https://github.com/go-sql-driver/mysql) | MySQL 持久化审查记录 |
| OAuth | [golang.org/x/oauth2](https://pkg.go.dev/golang.org/x/oauth2) | GitHub Token 认证 |
| 配置解析 | [gopkg.in/yaml.v3](https://gopkg.in/yaml.v3) | 读取 `config.yaml` |

### 前端（`ai-pr-review-front`）

> **说明：** 前端页面所使用的图片均来源于网络。

| 类别 | 名称 | 用途 |
|------|------|------|
| UI 框架 | [React 18](https://react.dev/) | 页面组件 |
| 语言 | [TypeScript](https://www.typescriptlang.org/) | 类型安全 |
| 构建工具 | [Vite 5](https://vitejs.dev/) | 开发与生产构建 |
| 路由 | [React Router v6](https://reactrouter.com/) | 页面路由 |
| 状态管理 | [Zustand](https://github.com/pmndrs/zustand) | 全局状态 |
| 样式 | [Tailwind CSS 3](https://tailwindcss.com/) | UI 样式 |
| CSS 工具链 | PostCSS、Autoprefixer | Tailwind 编译与浏览器前缀 |

### 外部服务

| 服务 | 用途 |
|------|------|
| [GitHub API](https://docs.github.com/en/rest) | 获取 PR、diff、commits；Approve / Merge / 打回 |
| MySQL 5.7+ / 8.0+ | 存储审查历史与结果 |
| OpenAI 兼容 LLM API（如 [DeepSeek](https://platform.deepseek.com/)） | AI 代码风险分析 |

---

## Demo 演示视频

<!-- 在此填写演示视频链接 -->

**视频链接：** _（待填写）_

### 启动前环境变量

启动本项目前，需先复制后端配置模板并按 [`ai-pr-review/config/config.yaml.example`](ai-pr-review/config/config.yaml.example) 设置对应环境变量：

```bash
cd ai-pr-review
cp config/config.yaml.example config/config.yaml
```

`config.yaml` 中通过 `${VAR_NAME}` 引用环境变量，需配置以下项：

| 环境变量 | 配置项 | 说明 | 示例 |
|----------|--------|------|------|
| `GITHUB_TOKEN` | `github_token` | GitHub Personal Access Token（需 `repo` 权限） | `ghp_xxx` |
| `OPENAI_API_KEY` | `openai_api_key` | OpenAI 兼容 API Key | `sk-xxx` |
| `MODEL` | `model` | LLM 模型名称 | `deepseek-chat` |
| `DB_HOST` | `database.host` | MySQL 主机地址 | `127.0.0.1` |
| `DB_PORT` | `database.port` | MySQL 端口 | `3306` |
| `DB_USER` | `database.user` | MySQL 用户名 | `root` |
| `DB_PASSWORD` | `database.password` | MySQL 密码 | `your_password` |
| `DB_NAME` | `database.name` | MySQL 数据库名 | `ai_pr_review` |

PowerShell 示例：

```powershell
$env:GITHUB_TOKEN = "ghp_xxx"
$env:OPENAI_API_KEY = "sk-xxx"
$env:MODEL = "deepseek-chat"
$env:DB_HOST = "127.0.0.1"
$env:DB_PORT = "3306"
$env:DB_USER = "root"
$env:DB_PASSWORD = "your_password"
$env:DB_NAME = "ai_pr_review"
```


配置完成后，分别启动后端（`go run ./cmd/api`）与前端（`npm run dev`）。`api_base`、`output_dir`、`prompt_dir` 及 `server` 等项已在配置模板中写死，无需额外环境变量。

---

## 项目概览

```
PR URL → GitHub API → Diff Parser → AI Review (per chunk) → JSON + MySQL
                                                              ↓
                                                    React Web 前端展示
```

本仓库为 monorepo，包含两个子项目：

| 目录 | 说明 |
|------|------|
| [`ai-pr-review/`](ai-pr-review/) | Go 后端：CLI 交互审查 + REST API |
| [`ai-pr-review-front/`](ai-pr-review-front/) | React 前端：Web 界面，对接后端 API |

### 核心能力

1. 解析 GitHub PR URL，拉取标题、作者、变更文件、commits、raw diff
2. 将 diff 按 hunk 切分为 chunk，逐块调用 LLM 进行风险分析
3. 输出分级 issues（high / medium / low），保存 JSON 并写入 MySQL
4. **CLI 模式**：审查完成后可交互 Approve / Merge 或打回 PR 并评论
5. **API 模式**：REST 接口提交审查、查询结果与历史记录
6. **Web 前端**：新建审查、历史列表（分页 / 筛选）、详情页（问题分组、原始 Diff 风险高亮、PR 变更总结、Approve / Merge / 打回）

---

## 快速开始

### 1. 后端

```bash
cd ai-pr-review
cp config/config.yaml.example config/config.yaml
```

设置环境变量（PowerShell 示例）：

```powershell
$env:GITHUB_TOKEN = "ghp_xxx"
$env:OPENAI_API_KEY = "sk-xxx"
$env:MODEL = "deepseek-chat"
```

启动 API 服务：

```bash
go run ./cmd/api
# 默认 http://0.0.0.0:8080
```

或 CLI 模式：

```bash
go run ./cmd/cli
# 输入 PR URL，如 https://github.com/org/repo/pull/123
```

### 2. 前端

```bash
cd ai-pr-review-front
npm install
npm run dev
# 浏览器访问 http://localhost:5173
```

开发环境通过 Vite 代理将 `/api`、`/health` 转发到 `http://localhost:8080`。

### 3. 前置条件

- Go 1.23+（项目 `go.mod` 要求 1.25+）
- Node.js 18+
- MySQL 5.7+ / 8.0+
- GitHub Personal Access Token（`repo` 权限；CLI / 前端 Approve / Merge 还需写入权限）
- OpenAI 兼容 API Key

---

## REST API 概览

完整文档见 [`ai-pr-review/api/openapi.yaml`](ai-pr-review/api/openapi.yaml)。

| 方法 | 路径 | 说明 |
|------|------|------|
| `GET` | `/health` | 健康检查 |
| `POST` | `/api/v1/reviews` | 提交 PR URL 触发审查 |
| `GET` | `/api/v1/reviews/:id` | 按 ID 查看审查结果 |
| `GET` | `/api/v1/reviews` | 分页查询历史记录 |

示例：

```bash
curl -X POST http://localhost:8080/api/v1/reviews \
  -H "Content-Type: application/json" \
  -d '{"pr_url":"https://github.com/org/repo/pull/123"}'
```

---

## 前端路由

| 路径 | 说明 |
|------|------|
| `/review/new` | 新建审查 |
| `/reviews` | 历史记录 |
| `/reviews/:id` | 审查详情 |

---

## 项目结构

```
PR-check-r/
├── ai-pr-review/              # Go 后端
│   ├── cmd/cli/               # CLI 入口
│   ├── cmd/api/               # HTTP API 入口
│   ├── internal/              # 业务逻辑（github、parser、ai、review…）
│   ├── prompts/review/        # AI 提示词模板
│   ├── api/openapi.yaml       # API 文档
│   ├── config/                # 配置
│   ├── migrations/            # 数据库迁移
│   └── output/                # JSON 审查结果
├── ai-pr-review-front/        # React 前端
│   └── src/                   # 页面、组件、状态、API 客户端
└── README.md
```

---

## 构建与测试

**后端：**

```bash
cd ai-pr-review
make test          # 运行测试
make build         # 编译 bin/api 和 bin/cli
make run-api       # 启动 API
make run-cli       # 启动 CLI
```

**前端：**

```bash
cd ai-pr-review-front
npm run build      # 生产构建
npm run preview    # 预览构建产物
```

---

## 配置说明

`ai-pr-review/config/config.yaml` 通过 `${VAR_NAME}` 引用环境变量：

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

敏感信息请使用环境变量，勿写入配置文件并提交到 git。

前端环境变量：

| 文件 | 变量 | 说明 |
|------|------|------|
| `.env.development` | `VITE_API_BASE_URL=` | 留空，走 Vite 代理 |
| `.env.production` | `VITE_API_BASE_URL=http://localhost:8080` | 生产 API 地址 |

---

## 输出示例

```json
{
  "pr_title": "config+main",
  "pr_number": 1,
  "total_issues": 2,
  "high_issues": 1,
  "medium_issues": 0,
  "low_issues": 1,
  "review_result": {
    "high": [
      {
        "file": "config.yaml",
        "line": 42,
        "message": "硬编码敏感凭据：在配置文件中发现明文 token，存在泄露风险。",
        "suggestion": "移除硬编码的 token，改为通过环境变量注入。"
      }
    ],
    "low": [
      {
        "file": "main.go",
        "line": 6,
        "message": "无限循环没有退出条件，可能导致程序无法正常退出。",
        "suggestion": "考虑添加退出条件以便优雅关闭。"
      }
    ]
  }
}
```

---

## 扩展点

- **LLM 接口**：`internal/ai.LLM` 由 langchaingo 实现，可替换为其他 provider
- **提示词**：编辑 `ai-pr-review/prompts/review/` 下模板即可调整审查策略
- **FunctionDetector**：`internal/parser.FunctionDetector` 可扩展 Python / JS 函数识别

---

## 注意事项

- 大型 PR 的 diff 会按 chunk 多次调用 LLM，API 模式当前为同步响应
- 需要网络访问 GitHub API 与 LLM API
- CLI / API 的 Merge 受仓库 branch protection 限制，可能返回 405 / 422

---

## 子项目文档

- 后端详细说明：[`ai-pr-review/README.md`](ai-pr-review/README.md)
- 前端详细说明：[`ai-pr-review-front/README.md`](ai-pr-review-front/README.md)
