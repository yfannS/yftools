# md2html Backend

Go 后端服务，基于 Gin 框架，采用 Clean Architecture（六边形架构）设计，为前端工具箱平台提供 API 支持。

## 技术栈

- **Go** 1.22+ / **Gin** 1.12
- **MySQL** 8.0+ / go-sql-driver/mysql
- **Goldmark** 1.7.4（Markdown 解析）
- **JWT** (golang-jwt/jwt/v5)
- **Docker** 部署支持

## 目录结构

```
backend/
├── cmd/server/main.go           # 入口：依赖注入 + 启动
├── internal/
│   ├── config/config.go         # 配置加载（环境变量）
│   ├── model/model.go           # 数据模型（User, Conversion, Theme）
│   ├── middleware/
│   │   ├── auth.go              # JWT 认证中间件
│   │   ├── cors.go              # CORS 跨域中间件
│   │   └── recovery.go          # Panic 恢复中间件
│   ├── handler/
│   │   ├── auth_handler.go      # 认证 Handler
│   │   ├── health_handler.go    # 健康检查
│   │   └── tools/md2html/       # md2html 工具 Handler
│   │       ├── convert_handler.go
│   │       ├── history_handler.go
│   │       └── theme_handler.go
│   ├── service/
│   │   ├── user_service.go      # 用户 Service
│   │   └── tools/md2html/       # md2html 工具 Service
│   │       ├── convert_service.go
│   │       ├── history_service.go
│   │       └── theme_service.go
│   ├── repository/
│   │   ├── user_repository.go   # 用户 Repository
│   │   └── tools/md2html/
│   │       └── history_repository.go
│   └── router/router.go         # 路由注册
├── pkg/
│   ├── converter/converter.go   # Goldmark Markdown 转换器
│   ├── jwt/jwt.go               # JWT 生成/验证工具
│   └── response/response.go     # 统一响应格式
├── migrations/                   # SQL 迁移文件
├── configs/                      # 配置文件
├── .env.example                  # 环境变量模板（已弃用）
├── Dockerfile                    # Docker 构建
└── go.mod
```

## 快速开始

### 环境要求

- Go 1.22+
- MySQL 8.0+

### 1. 配置

```bash
cp configs/config.yaml.example configs/config.yaml
```

编辑 `configs/config.yaml` 填入实际值：

```yaml
server:
  port: 8000
  mode: debug

database:
  host: 127.0.0.1
  port: 3306
  name: yftools
  user: root
  password: your_password

jwt:
  secret: your-secret-key
  expire: 168h
```

### 2. 安装依赖

```bash
go mod tidy
```

### 3. 启动服务

```bash
go run cmd/server/main.go
```

服务默认监听 `:8000`。

### 4. 构建二进制

```bash
go build -o md2html-server cmd/server/main.go
./md2html-server
```

## API 接口

### 认证（跨工具共享）

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/auth/register` | 用户注册 |
| POST | `/api/auth/login` | 用户登录 |
| GET | `/api/auth/profile` | 获取用户信息（需 JWT） |

#### 注册

```bash
curl -X POST http://localhost:8000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@example.com","password":"123456"}'
```

#### 登录

```bash
curl -X POST http://localhost:8000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"123456"}'
```

### md2html 工具

| 方法 | 路径 | 认证 | 说明 |
|------|------|------|------|
| POST | `/api/tools/md2html/convert` | - | Markdown 转 HTML |
| GET | `/api/tools/md2html/history` | JWT | 获取转换历史（分页） |
| POST | `/api/tools/md2html/history` | JWT | 保存转换记录 |
| DELETE | `/api/tools/md2html/history/:id` | JWT | 删除历史记录 |
| GET | `/api/tools/md2html/themes` | - | 获取主题列表 |

#### Markdown 转换

```bash
curl -X POST http://localhost:8000/api/tools/md2html/convert \
  -H "Content-Type: application/json" \
  -d '{"markdown":"# Hello\n\n**world**"}'
```

响应：

```json
{
  "code": 0,
  "data": {
    "html": "<h1>Hello</h1>\n<p><strong>world</strong></p>\n"
  }
}
```

### 健康检查

```bash
curl http://localhost:8000/api/health
```

## 数据库

数据库表需提前创建，启动时不会自动建表。

### 手动迁移

如需手动创建表，可执行 `migrations/` 目录下的 SQL 文件。

## 架构设计

### 依赖方向

```
Handler → Service → Repository → Database
   ↓         ↓          ↓
  请求解析  业务逻辑    数据持久化
  响应构建  领域模型    SQL 查询
```

### 扩展新工具

1. 在 `internal/handler/tools/` 下创建新目录（如 `json_formatter/`）
2. 在 `internal/service/tools/` 下创建对应 Service
3. 在 `internal/repository/tools/` 下创建对应 Repository
4. 在 `internal/router/router.go` 注册新路由组：
   ```go
   jsonTools := toolsGroup.Group("/json-formatter")
   {
       jsonTools.POST("/format", ...)
   }
   ```

## Docker 部署

```bash
# 单独构建后端镜像
docker build -t md2html-backend .

# 或使用 docker-compose（从项目根目录）
docker-compose up -d backend
```

## License

MIT
