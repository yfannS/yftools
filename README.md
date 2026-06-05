# md2html-toolbox

一个前后端分离的在线工具箱平台，首个工具为 **Markdown → HTML 转换器**。采用 Vue 3 + Go + MySQL 架构，支持未来扩展更多工具。

## 项目结构

```
md2html/
├── backend/           # Go 后端服务 (Gin + MySQL)
├── webapp/            # Vue 3 前端应用 (Vite + TypeScript)
├── frontend/          # 原始单文件前端（保留参考）
├── docs/              # 技术文档
│   └── old/           # 早期技术方案
├── docker-compose.yml # 容器编排
└── .gitignore
```

## 快速启动

### 前置条件

- Go 1.22+
- Node.js 18+
- MySQL 8.0+（或使用 Docker）
- Docker & Docker Compose（可选）

### 方式一：Docker Compose（推荐）

```bash
# 一键启动所有服务
docker-compose up -d

# 前端：http://localhost:3000
# 后端：http://localhost:8000
# MySQL：localhost:3306
```

### 方式二：本地开发

#### 1. 启动后端

```bash
cd backend

# 配置环境变量
cp .env.example .env
# 编辑 .env 填入 MySQL 连接信息

# 安装依赖 & 启动
go mod tidy
go run cmd/server/main.go
```

后端默认运行在 `http://localhost:8000`。

#### 2. 启动前端

```bash
cd webapp

# 安装依赖
npm install

# 开发模式
npm run dev
```

前端默认运行在 `http://localhost:5173`。

#### 3. 生产构建

```bash
# 前端
cd webapp
npm run build        # 输出到 dist/

# 后端
cd backend
go build -o md2html-server cmd/server/main.go
```

## 技术栈

| 层级 | 技术 |
|------|------|
| 前端 | Vue 3 + TypeScript + Vite + Pinia + Vue Router |
| 后端 | Go 1.22+ + Gin 1.12 + Goldmark 1.7 |
| 数据库 | MySQL 8.0 |
| 部署 | Docker + Docker Compose |
| 认证 | JWT (HS256) |

## API 设计

所有工具 API 统一在 `/api/tools/` 下，认证跨工具共享在 `/api/auth/`：

| 方法 | 路径 | 认证 | 说明 |
|------|------|------|------|
| GET | `/api/health` | - | 健康检查 |
| POST | `/api/auth/register` | - | 用户注册 |
| POST | `/api/auth/login` | - | 用户登录 |
| GET | `/api/auth/profile` | JWT | 获取用户信息 |
| POST | `/api/tools/md2html/convert` | - | Markdown 转 HTML |
| GET | `/api/tools/md2html/history` | JWT | 获取转换历史 |
| POST | `/api/tools/md2html/history` | JWT | 保存转换记录 |
| DELETE | `/api/tools/md2html/history/:id` | JWT | 删除历史记录 |
| GET | `/api/tools/md2html/themes` | - | 获取主题列表 |

> 未来新增工具时，API 路径按 `/api/tools/<tool-name>/...` 扩展。

## 核心设计原则

- **本地渲染优先**：Markdown 转换优先使用浏览器端 marked.js + Web Worker，后端 Goldmark 作为兼容能力
- **工具箱架构**：每个工具独立目录（handler/service/repository），共享认证、数据库、中间件
- **渐进式迁移**：前端从单文件 HTML 迁移到 Vue 3 组件化，保留完整功能和视觉设计

## License

MIT
