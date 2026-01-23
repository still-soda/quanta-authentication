# QAuth Server

基于 Go 的 OAuth 2.0 / OpenID Connect 认证服务。

## 技术栈

| 类别 | 技术 |
|:---|:---|
| 语言 | Go 1.25 |
| Web 框架 | Gin |
| ORM | GORM |
| 数据库 | PostgreSQL 16 |
| 缓存 | Redis |
| 存储 | Go CDK（支持本地 / S3） |

## 项目结构

```
qauth-server/
├── main.go                    # 入口文件
├── go.mod                     # 依赖管理
├── internal/
│   ├── config/                # 配置管理
│   ├── database/              # 数据库连接与迁移
│   ├── errors/                # 错误定义
│   ├── handlers/              # 请求处理器
│   │   ├── auth.go            # 认证处理
│   │   ├── oauth.go           # OAuth2 处理
│   │   ├── oidc.go            # OIDC 处理
│   │   ├── role.go            # 角色管理
│   │   └── file.go            # 文件上传
│   ├── middleware/            # 中间件
│   ├── models/                # 数据模型
│   ├── permissions/           # 权限定义
│   ├── routes/                # 路由注册
│   ├── services/              # 业务逻辑
│   └── utilities/             # 工具函数
├── pkg/
│   ├── jwks/                  # JWKS 密钥管理
│   ├── jwt/                   # JWT 工具
│   └── response/              # 响应封装
└── docs/                      # 文档
```

## 快速开始

### 1. 安装依赖

```bash
go mod download
```

### 2. 配置环境变量

复制 `.env.example` 到 `.env` 并修改配置。

### 3. 启动服务

```bash
go run main.go
```

服务将在 http://localhost:8080 启动。

## 环境变量

| 变量 | 默认值 | 说明 |
|:---|:---|:---|
| `PORT` | `8080` | 服务端口 |
| `GIN_MODE` | `debug` | 运行模式 |
| `DATABASE_URL` | - | PostgreSQL 连接字符串 |
| `REDIS_ADDR` | `localhost:6379` | Redis 地址 |
| `REDIS_PASSWORD` | - | Redis 密码 |
| `JWT_SECRET` | - | JWT 签名密钥 |
| `ACCESS_TOKEN_EXPIRE` | `3600` | Access Token 有效期（秒） |
| `REFRESH_TOKEN_EXPIRE` | `604800` | Refresh Token 有效期（秒） |
| `OIDC_ISSUER` | `http://localhost:8080` | OIDC 发行者 URL |
| `OIDC_KEY_ROTATION_INTERVAL` | `86400` | JWKS 密钥轮换间隔（秒） |
| `STORAGE_LOCAL_DIR` | `./uploads` | 本地存储目录 |

## 构建

```bash
go build -o qauth-server main.go
```
