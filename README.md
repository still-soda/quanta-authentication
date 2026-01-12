# Quanta Authentication

Quanta Authentication 是一个基于 OAuth 2.0/OIDC 标准的认证授权系统，使用 Turborepo 管理的 Monorepo 架构。

## 系统架构

```
┌─────────────────────────────────────────────────────────────┐
│                        Nginx Gateway                         │
│                    (Port 80/443)                             │
└────────────┬───────────────────────────────┬────────────────┘
             │                               │
             │ /api/*                        │ /
             ▼                               ▼
┌──────────────────────┐          ┌──────────────────────┐
│  qauth-server        │          │  qauth-admin-web    │
│  (Go + Gin)          │          │  (Vue 3 + PrimeVue) │
│  Port: 8080          │          │  Port: 80           │
│                      │          │                      │
│  OAuth 2.0 Server    │          │  Admin Dashboard    │
└──────────┬───────────┘          └──────────────────────┘
           │
           ├───────────────────┐
           │                   │
           ▼                   ▼
┌──────────────────┐  ┌──────────────────┐
│ PostgreSQL 16    │  │ Redis 6         │
│                 │  │                 │
│ Client Store     │  │ Token Store     │
└──────────────────┘  └──────────────────┘
```

## 技术栈

### 后端服务 (@qauth/server)
- **语言**: Go 1.25
- **框架**: Gin
- **认证**: go-oauth2/oauth2 (OAuth 2.0/OIDC)
- **数据库**: PostgreSQL + pgx/v4
- **缓存**: Redis
- **令牌存储**: go-oauth2/redis
- **配置管理**: godotenv

### 前端管理台 (@qauth/admin-web)
- **框架**: Vue 3.5
- **语言**: TypeScript
- **构建工具**: Vite (rolldown-vite)
- **UI 组件**: PrimeVue 4.5 + PrimeUx Themes
- **状态管理**: Pinia
- **路由**: Vue Router
- **样式**: TailwindCSS 4
- **HTTP 客户端**: Axios

### 基础设施
- **容器编排**: Docker Compose
- **反向代理**: Nginx
- **包管理**: pnpm
- **构建系统**: Turborepo

## 核心功能

### OAuth 2.0 支持
- ✅ Authorization Code Flow
- ✅ Access Token 管理
- ✅ Refresh Token 支持
- ✅ 客户端认证
- ✅ 用户授权

### 管理功能
- 客户端应用管理
- 用户授权管理
- 令牌审计
- 配置管理

## 项目结构

```
quanta-authentication/
├── packages/
│   ├── qauth-server/          # Go 认证服务器
│   │   ├── main.go            # 应用入口
│   │   ├── go.mod             # Go 依赖
│   │   ├── Dockerfile         # 容器镜像
│   │   └── internal/
│   │       └── config/        # 配置管理
│   │
│   └── qauth-admin-web/       # Vue 3 管理前端
│       ├── src/
│       │   ├── App.vue        # 根组件
│       │   ├── main.ts        # 应用入口
│       │   └── style.css      # 全局样式
│       ├── vite.config.ts     # Vite 配置
│       ├── tsconfig.json      # TypeScript 配置
│       ├── Dockerfile         # 容器镜像
│       └── nginx/             # Nginx 配置
│
├── docker/
│   └── docker-compose.yaml    # 服务编排
│
├── nginx/
│   └── default.conf           # Nginx 网关配置
│
├── turbo.json                 # Turborepo 配置
├── package.json               # 根包配置
└── pnpm-workspace.yaml        # pnpm 工作空间
```

## 快速开始

### 前置要求

- Node.js >= 18
- Go 1.25
- pnpm >= 9.0
- Docker & Docker Compose

### 开发模式

#### 启动所有服务（本地开发）

```bash
# 安装依赖
pnpm install

# 启动开发服务器
pnpm dev
```

#### 使用 Docker Compose（推荐）

```bash
cd docker

# 启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

### 访问应用

- **管理后台**: http://localhost
- **认证 API**: http://localhost/api
- **PostgreSQL**: localhost:15432
- **Redis**: localhost:16379

## 环境变量

在 `docker/docker-compose.yaml` 中配置或在 `packages/qauth-server/.env` 中设置：

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `DATABASE_URL` | `postgres://user:password@localhost:5432/qauth_db` | PostgreSQL 连接字符串 |
| `REDIS_URL` | `127.0.0.1:6379` | Redis 连接地址 |
| `REDIS_PASSWORD` | `` | Redis 密码 |
| `ACCESS_TOKEN_EXPIRE` | `3600` | 访问令牌有效期（秒） |
| `PORT` | `8080` | 服务器端口 |
| `API_URL` | `http://qauth-server:8080` | 前端访问的 API 地址 |

## API 端点

### OAuth 2.0 端点

| 方法 | 端点 | 说明 |
|------|------|------|
| GET | `/api/authorize` | 授权端点 |
| POST | `/api/token` | 令牌端点 |

### 健康检查

| 方法 | 端点 | 说明 |
|------|------|------|
| GET | `/health` | 服务健康状态 |

## 构建

### 构建所有包

```bash
pnpm build
```

### 构建特定包

```bash
# 构建 Go 服务
cd packages/qauth-server
go build -o dist/server cmd/server/main.go

# 构建 Vue 前端
cd packages/qauth-admin-web
pnpm build
```

### 构建镜像

```bash
# 在 docker 目录下
docker-compose build
```

## 开发

### 单独运行后端

```bash
cd packages/qauth-server
go run cmd/server/main.go
```

### 单独运行前端

```bash
cd packages/qauth-admin-web
pnpm dev
```

## 配置说明

### Nginx 网关

Nginx 作为统一网关，路由规则如下：

- `/` → qauth-admin-web (前端静态资源)
- `/api/*` → qauth-server (API 请求)
- `/health` → 健康检查端点

### 数据库

PostgreSQL 用于存储：
- OAuth 客户端信息
- 用户授权记录
- 令牌元数据

### Redis

Redis 用于存储：
- 访问令牌（Access Tokens）
- 刷新令牌（Refresh Tokens）
- 授权码（Authorization Codes）

## 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## License

本项目采用 MIT 许可证
