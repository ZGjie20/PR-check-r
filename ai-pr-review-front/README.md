# AI PR Review 前端

基于 React + Vite + TypeScript 的 AI PR 审查 Web 前端，对接 `ai-pr-review` 后端 REST API。

## 功能

- 提交 GitHub PR URL 触发 AI 代码审查
- 查看审查历史记录（分页、PR 编号筛选）
- 查看审查详情（问题分组、原始 Diff）

## 前置条件

- Node.js 18+
- 后端 API 已启动（默认 `http://localhost:8080`）

## 快速开始

```bash
# 安装依赖
cd ai-pr-review-front
npm install

# 启动开发服务器
npm run dev
```

浏览器访问 `http://localhost:5173`。

## 后端启动

在另一个终端启动后端：

```bash
cd ai-pr-review
go run ./cmd/api
```

后端需配置 `GITHUB_TOKEN`、`OPENAI_API_KEY` 及 MySQL 连接。

## 环境变量

| 文件 | 变量 | 说明 |
|------|------|------|
| `.env.development` | `VITE_API_BASE_URL=` | 留空，通过 Vite 代理转发 |
| `.env.production` | `VITE_API_BASE_URL=http://localhost:8080` | 生产环境 API 地址 |

开发环境已在 `vite.config.ts` 配置代理，解决后端未启用 CORS 的问题：

```typescript
proxy: {
  '/api': 'http://localhost:8080',
  '/health': 'http://localhost:8080',
}
```

## 构建

```bash
npm run build
npm run preview
```

## 路由

| 路径 | 说明 |
|------|------|
| `/review/new` | 新建审查 |
| `/reviews` | 历史记录 |
| `/reviews/:id` | 审查详情 |

## 技术栈

- React 18 + TypeScript
- Vite 5
- Tailwind CSS 3
- React Router v6
- Zustand
