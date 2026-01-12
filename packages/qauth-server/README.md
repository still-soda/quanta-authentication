# Quanta Authentication Server

基于 Gin + GORM + Go CDK 的现代化认证服务器

## 技术栈

- **Web 框架**: Gin
- **ORM**: GORM
- **数据库**: PostgreSQL
- **存储**: Go CDK (支持本地存储和 S3)
- **缓存**: Redis (可选)

## 项目结构

```
qauth-server/
├── main.go                 # 入口文件
├── go.mod                  # 依赖管理
├── .env.example            # 环境变量示例
├── internal/               # 内部包
│   ├── config/             # 配置管理
│   │   └── config.go
│   ├── database/           # 数据库连接
│   │   └── database.go
│   ├── storage/            # 存储服务 (Go CDK)
│   │   └── storage.go
│   ├── models/             # 数据模型
│   │   └── models.go
│   ├── handlers/           # 路由处理器
│   │   ├── health.go
│   │   ├── user.go
│   │   └── file.go
│   └── middleware/         # 中间件
│       ├── logger.go
│       ├── cors.go
│       └── recovery.go
└── uploads/                # 本地文件上传目录
```

## 快速开始

### 1. 安装依赖

```bash
go mod download
```

### 2. 配置环境变量

复制 `.env.example` 到 `.env` 并修改配置：

```bash
cp .env.example .env
```

### 3. 创建数据库

```bash
createdb qauth
```

或使用 PostgreSQL 客户端：

```sql
CREATE DATABASE qauth;
```

### 4. 运行服务

```bash
go run main.go
```

服务将在 <http://localhost:8080> 启动

## API 接口

### 健康检查

- `GET /health` - 健康检查
- `GET /ping` - Ping 检查

### 用户管理

- `POST /api/v1/users/register` - 用户注册
- `POST /api/v1/users/login` - 用户登录
- `GET /api/v1/users/:id` - 获取用户信息
- `GET /api/v1/users` - 用户列表

### 文件管理

- `POST /api/v1/files/upload` - 上传文件
- `GET /api/v1/files/:key` - 下载文件
- `DELETE /api/v1/files/:id` - 删除文件
- `GET /api/v1/files` - 文件列表

## 存储配置

### 本地存储 (默认)

在 `.env` 中设置：

```env
STORAGE_PROVIDER=local
STORAGE_LOCAL_DIR=./uploads
```

### S3 存储

在 `.env` 中设置：

```env
STORAGE_PROVIDER=s3
STORAGE_S3_BUCKET=your-bucket-name
STORAGE_S3_REGION=us-east-1
```

确保已配置 AWS 凭证（通过环境变量或 `~/.aws/credentials`）

## Go CDK 存储特性

本项目使用 [Go Cloud Development Kit (Go CDK)](https://gocloud.dev/) 实现存储抽象，具有以下优势：

- **统一接口**: 同一套代码支持多种存储后端
- **易于切换**: 通过配置即可切换存储提供商
- **支持多种后端**: 本地文件系统、S3、GCS、Azure Blob 等

## 示例请求

### 用户注册

```bash
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123",
    "nickname": "Test User"
  }'
```

### 文件上传

```bash
curl -X POST http://localhost:8080/api/v1/files/upload \
  -F "file=@/path/to/your/file.jpg"
```

## 开发

### 运行开发模式

```bash
GIN_MODE=debug go run main.go
```

### 构建生产版本

```bash
go build -o qauth-server main.go
```

## 环境变量说明

| 变量 | 说明 | 默认值 |
|------|------|--------|
| PORT | 服务端口 | 8080 |
| GIN_MODE | Gin 模式 (debug/release) | debug |
| DATABASE_URL | 数据库连接字符串 | - |
| REDIS_ADDR | Redis 地址 | localhost:6379 |
| STORAGE_PROVIDER | 存储提供商 (local/s3) | local |
| STORAGE_LOCAL_DIR | 本地存储目录 | ./uploads |
| JWT_SECRET | JWT 密钥 | - |

## 许可证

MIT
