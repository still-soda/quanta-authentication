# QAuth Server 产品需求文档 (PRD)

## 1. 产品概述

### 1.1 产品定位

**QAuth Server** 是一套基于 Go 语言开发的企业级身份认证与授权服务系统，主要面向需要统一身份管理的组织（如学校、企业）。系统实现了完整的 **OAuth 2.0** 和 **OpenID Connect (OIDC)** 协议，提供用户身份管理、单点登录（SSO）、以及基于角色的访问控制（RBAC）能力。

### 1.2 核心价值

| 价值点 | 描述 |
|--------|------|
| 统一身份认证 | 为组织内多个应用提供统一的用户认证服务 |
| 标准协议支持 | 完整实现 OAuth 2.0 / OIDC 协议，便于第三方应用接入 |
| 细粒度权限控制 | 基于 RBAC 模型的权限管理体系 |
| 安全可靠 | 支持密钥自动轮换、密码加盐哈希、令牌过期机制等安全特性 |

### 1.3 目标用户

- **学校/企业管理员**：管理用户、角色、权限、OAuth 客户端
- **开发者**：接入 OAuth2/OIDC 协议，实现第三方应用的用户认证
- **终端用户**：通过学号/工号进行身份认证

---

## 2. 核心业务流程

### 2.1 用户注册与认证流程

```
┌─────────┐    ┌──────────────┐    ┌─────────────┐    ┌──────────────┐
│  注册   │───▶│  参数校验    │───▶│  唯一性检查  │───▶│  密码加盐    │
│         │    │ (StudentID,  │    │ (StudentID, │    │  哈希存储    │
│         │    │  Email等)    │    │   Email)    │    │              │
└─────────┘    └──────────────┘    └─────────────┘    └──────────────┘
                                                              │
                                                              ▼
┌─────────┐    ┌──────────────┐    ┌─────────────┐    ┌──────────────┐
│  登录   │───▶│  查找用户    │───▶│  密码验证   │───▶│  状态检查    │
│         │    │ (StudentID)  │    │ (Salt+Hash) │    │  (ACTIVE)    │
└─────────┘    └──────────────┘    └─────────────┘    └──────────────┘
                                                              │
                                                              ▼
                                                      ┌──────────────┐
                                                      │ 生成JWT令牌  │
                                                      │ (Access +    │
                                                      │  Refresh)    │
                                                      └──────────────┘
```

### 2.2 OAuth 2.0 授权码流程

```
┌─────────┐    ┌──────────────┐    ┌─────────────┐    ┌──────────────┐
│ 客户端  │───▶│  /authorize  │───▶│  用户认证   │───▶│  授权确认    │
│ 发起    │    │  端点        │    │ (Cookie中   │    │              │
│ 请求    │    │              │    │  JWT Token) │    │              │
└─────────┘    └──────────────┘    └─────────────┘    └──────────────┘
                                                              │
                                                              ▼
┌─────────┐    ┌──────────────┐    ┌─────────────┐    ┌──────────────┐
│ 返回    │◀───│ 生成令牌     │◀───│  验证Code   │◀───│ 重定向+Code  │
│ Tokens  │    │ (Access +    │    │  & Client   │    │              │
│         │    │  ID Token)   │    │             │    │              │
└─────────┘    └──────────────┘    └─────────────┘    └──────────────┘
```

### 2.3 令牌刷新流程

```
┌─────────────────┐    ┌──────────────────┐    ┌───────────────────┐
│  提交           │───▶│  解析并验证      │───▶│  生成新的         │
│  Refresh Token  │    │  Refresh Token   │    │  Access Token +   │
│                 │    │                  │    │  Refresh Token    │
└─────────────────┘    └──────────────────┘    └───────────────────┘
```

### 2.4 权限校验流程

```
┌─────────────────┐    ┌──────────────────┐    ┌───────────────────┐
│  请求携带       │───▶│  Auth中间件      │───▶│  解析JWT获取      │
│  Authorization  │    │  提取Token       │    │  UserInfo         │
│  Header         │    │                  │    │                   │
└─────────────────┘    └──────────────────┘    └───────────────────┘
                                                        │
                                                        ▼
┌─────────────────┐    ┌──────────────────┐    ┌───────────────────┐
│  允许/拒绝      │◀───│  检查角色是否    │◀───│  获取用户角色     │
│  访问           │    │  包含所需权限    │    │  及其权限列表     │
└─────────────────┘    └──────────────────┘    └───────────────────┘
```

---

## 3. 功能清单 (Feature List)

### 3.1 基础功能

#### 3.1.1 用户管理模块

| 功能 | API 端点 | 方法 | 认证要求 | 描述 |
|------|----------|------|----------|------|
| 用户注册 | `/_/v1/auth/register` | POST | 无 | 新用户注册，需提供 StudentID、Password、Email、Name |
| 用户登录 | `/_/v1/auth/login` | POST | 无 | 通过 StudentID + Password 认证，返回 Access Token 和 Refresh Token |
| 令牌刷新 | `/_/v1/auth/refresh-token` | POST | 无 | 使用 Refresh Token 获取新的令牌对 |

#### 3.1.2 OAuth 2.0 / OIDC 模块

| 功能 | API 端点 | 方法 | 描述 |
|------|----------|------|------|
| 授权端点 | `/v1/oauth/authorize` | GET/POST | OAuth2 授权请求入口 |
| 令牌端点 | `/v1/oauth/token` | POST | 获取 Access Token / ID Token |
| 令牌验证 | `/v1/oauth/validate` | POST | 验证 Access Token 有效性 |
| 令牌撤销 | `/v1/oauth/revoke` | POST | 撤销 Access Token 或 Refresh Token |
| 用户信息 | `/v1/oauth/userinfo` | GET | 获取当前认证用户信息（根据 scope 返回不同字段） |
| 登出 | `/v1/oauth/logout` | GET/POST | 撤销令牌并重定向 |

#### 3.1.3 OIDC 发现端点

| 功能 | API 端点 | 方法 | 描述 |
|------|----------|------|------|
| OpenID 配置 | `/.well-known/openid-configuration` | GET | OIDC 发现文档 |
| JWKS 端点 | `/.well-known/jwks.json` | GET | 公钥集合（用于验证 ID Token） |

#### 3.1.4 OAuth2 客户端管理

| 功能 | API 端点 | 方法 | 所需权限 |
|------|----------|------|----------|
| 创建客户端 | `/_/v1/clients` | POST | `oauth_client_create` |
| 获取客户端列表 | `/_/v1/clients` | GET | `oauth_client_list` |
| 获取客户端详情 | `/_/v1/clients/:id` | GET | `oauth_client_view` |
| 更新客户端 | `/_/v1/clients/:id` | PUT | `oauth_client_update` |
| 删除客户端 | `/_/v1/clients/:id` | DELETE | `oauth_client_delete` |

#### 3.1.5 角色管理

| 功能 | API 端点 | 方法 | 所需权限 |
|------|----------|------|----------|
| 获取角色列表 | `/_/v1/roles` | GET | `role_view` |
| 获取角色详情 | `/_/v1/roles/:id` | GET | `role_view` |
| 更新角色 | `/_/v1/roles/:id` | PUT | `role_update` |
| 删除角色 | `/_/v1/roles/:id` | DELETE | `role_delete` |

#### 3.1.6 权限管理

| 功能 | API 端点 | 方法 | 所需权限 |
|------|----------|------|----------|
| 授予权限给角色 | `/_/v1/permissions/grant-to-role` | POST | `role_assign_permissions` |
| 从角色撤销权限 | `/_/v1/permissions/revoke-from-role` | POST | `role_revoke_permissions` |

#### 3.1.7 密钥管理

| 功能 | API 端点 | 方法 | 描述 |
|------|----------|------|------|
| 获取密钥信息 | `/_/v1/jwks/keys` | GET | 获取所有密钥的轮换状态 |
| 强制密钥轮换 | `/_/v1/jwks/rotate` | POST | 立即触发密钥轮换 |

#### 3.1.8 资源管理

| 功能 | API 端点 | 方法 | 描述 |
|------|----------|------|------|
| 文件上传 | `/_/v1/resources/upload` | POST | 上传文件到存储服务 |

---

### 3.2 隐藏逻辑

#### 3.2.1 用户状态机

```
┌──────────┐
│  ACTIVE  │◀─── 默认状态（可正常登录）
└────┬─────┘
     │
     ▼
┌──────────┐
│  LOCKED  │◀─── 账户锁定（无法登录）
└────┬─────┘
     │
     ▼
┌──────────┐
│  BANNED  │◀─── 永久封禁
└──────────┘
```

**状态校验逻辑**（代码位置：`services/oauth.go`）：
- 在 OAuth2 授权时检查 `user.Status != models.UserStatusActive`
- 非 ACTIVE 状态用户无法完成 OAuth 认证

#### 3.2.2 权限校验机制

**校验流程**：
1. Auth 中间件从 Header 提取 JWT Token
2. 解析 Token 获取 `userInfo`（包含 `UserID`, `StudentID`, `Role`）
3. 将 `userInfo` 存入 Gin Context
4. 业务 Handler 调用 `VerifyPermissions()` 函数
5. 根据用户角色查询关联的权限列表
6. 判断是否包含请求所需的权限码

**权限码定义**（代码位置：`permissions/permissions.go`）：

| 权限码 | 描述 |
|--------|------|
| `oauth_client_create` | 创建 OAuth2 客户端 |
| `oauth_client_delete` | 删除 OAuth2 客户端 |
| `oauth_client_view` | 查看 OAuth2 客户端 |
| `oauth_client_update` | 更新 OAuth2 客户端 |
| `oauth_client_list` | 列出 OAuth2 客户端 |
| `role_create` | 创建角色 |
| `role_delete` | 删除角色 |
| `role_view` | 查看角色 |
| `role_update` | 更新角色 |
| `role_assign_permissions` | 授予权限 |
| `role_revoke_permissions` | 撤销权限 |

#### 3.2.3 登录状态记录

系统自动记录每次登录尝试（代码位置：`models/login_state.go`）：

| 登录类型 | 常量值 | 描述 |
|----------|--------|------|
| `PASSWORD` | 密码登录 | 系统内部登录 |
| `OAUTH2` | OAuth2 授权 | 第三方应用授权登录 |
| `REFRESH_TOKEN` | 令牌刷新 | 使用 Refresh Token 续签 |

记录字段包括：`UserID`, `ClientID`, `Time`, `IP`, `UserAgent`, `Location`, `IsSuccess`, `FailReason`

#### 3.2.4 错误记录

OAuth 服务内部错误自动记录到 `error_records` 表（代码位置：`services/oauth.go#recordError`）。

#### 3.2.5 唯一性冲突处理

用户注册时通过 PostgreSQL 错误码 `23505` 判断唯一性冲突，返回 `ErrCreateUserConflict` 错误。

#### 3.2.6 Scope 与 Claims 映射

**支持的 Scope**：

| Scope | 返回的 Claims |
|-------|---------------|
| `openid` | `sub`, `iss`, `aud`, `exp`, `iat` |
| `profile` | `name`, `student_id`, `display_name`, `picture` |
| `email` | `email`, `email_verified` |
| `roles` | `roles` |

---

## 4. 数据模型描述

### 4.1 核心实体关系图 (ER Diagram)

```
┌───────────────────┐       ┌───────────────────┐       ┌───────────────────┐
│      Users        │       │    UsersRoles     │       │      Roles        │
├───────────────────┤       ├───────────────────┤       ├───────────────────┤
│ ID (UUID, PK)     │◀──────│ UsersID (FK)      │──────▶│ ID (UUID, PK)     │
│ StudentID (UQ)    │       │ RolesID (FK)      │       │ Code (UQ)         │
│ Email (UQ)        │       │ AssignedAt        │       │ Name              │
│ PasswordHash      │       └───────────────────┘       │ Description       │
│ Salt              │                                   │ IsSystem          │
│ Name              │                                   └─────────┬─────────┘
│ Phone (UQ)        │                                             │
│ DisplayName       │                                             │
│ AvatarID (FK)     │──────────────────┐       ┌──────────────────┘
│ Status            │                  │       │
│ EmailVerified     │                  │       ▼
└─────────┬─────────┘       ┌───────────────────┐       ┌───────────────────┐
          │                 │ RolesPermissions  │       │   Permissions     │
          │                 ├───────────────────┤       ├───────────────────┤
          │                 │ RolesID (FK)      │◀──────│ ID (UUID, PK)     │
          │                 │ PermissionsID (FK)│       │ Resource          │
          │                 └───────────────────┘       │ Action            │
          │                                             │ Code (UQ)         │
          │                                             │ Description       │
          ▼                                             └───────────────────┘
┌───────────────────┐       ┌───────────────────┐
│     Images        │       │    LoginState     │
├───────────────────┤       ├───────────────────┤
│ ID (UUID, PK)     │       │ ID (UUID, PK)     │
│ Width             │       │ UserID (FK)       │◀──────┐
│ Height            │       │ ClientID (FK)     │       │
│ FileID (FK)       │       │ Time              │       │
│ CreatorID (FK)    │       │ Type              │       │
└───────────────────┘       │ IP                │       │
                            │ UserAgent         │       │
                            │ Location          │       │
                            │ IsSuccess         │       │
                            │ FailReason        │       │
                            └───────────────────┘       │
                                                        │
┌───────────────────┐       ┌───────────────────┐       │
│   OAuth2Client    │       │   Organization    │       │
├───────────────────┤       ├───────────────────┤       │
│ ID (UUID, PK)     │◀──────│ ID (UUID, PK)     │       │
│ Secret            │       │ UserID (FK)       │───────┘
│ Name              │       │ SuperiorID (FK)   │◀──┐
│ Domain            │       │ AncestorPath      │   │
│ Data (JSONB)      │       │ Depth             │   │
└───────────────────┘       │ OrgRole           │───┘
                            │ Class             │ (自引用)
                            └───────────────────┘
```

### 4.2 实体详细定义

#### 4.2.1 Users（用户）

| 字段 | 类型 | 约束 | 描述 |
|------|------|------|------|
| ID | UUID | PK, 自动生成 | 用户唯一标识 |
| StudentID | VARCHAR(11) | UQ, NOT NULL | 学号（业务主键） |
| Email | VARCHAR(100) | UQ, NOT NULL | 邮箱 |
| PasswordHash | VARCHAR(255) | NOT NULL | 密码哈希值 |
| Salt | VARCHAR(255) | - | 密码盐值 |
| Name | VARCHAR(50) | NOT NULL | 用户姓名 |
| Phone | VARCHAR(20) | UQ | 手机号 |
| DisplayName | VARCHAR(50) | - | 显示名称 |
| AvatarID | UUID | FK → Images | 头像图片 |
| Status | ENUM | 默认 ACTIVE | 用户状态 |
| EmailVerified | BOOLEAN | 默认 false | 邮箱是否验证 |

#### 4.2.2 Roles（角色）

| 字段 | 类型 | 约束 | 描述 |
|------|------|------|------|
| ID | UUID | PK | 角色唯一标识 |
| Code | VARCHAR(50) | UQ, NOT NULL | 角色代码（如 `system_super_admin`） |
| Name | VARCHAR(50) | NOT NULL | 角色名称 |
| Description | VARCHAR(255) | - | 角色描述 |
| IsSystem | BOOLEAN | 默认 false | 是否为系统内置角色 |

**预设角色**（种子数据）：

| Code | Name | 描述 |
|------|------|------|
| `system_super_admin` | 系统超级管理员 | 拥有系统内所有权限 |
| `system_admin` | 系统管理员 | 拥有系统内大部分权限 |
| `system_user` | 系统普通用户 | 拥有系统内基本权限 |

#### 4.2.3 Permissions（权限）

| 字段 | 类型 | 约束 | 描述 |
|------|------|------|------|
| ID | UUID | PK | 权限唯一标识 |
| Resource | VARCHAR(50) | NOT NULL | 资源名称（如 `oauth_clients`） |
| Action | INT8 | NOT NULL | 操作类型（0=Create, 1=Read, 2=Update, 3=Delete） |
| Code | VARCHAR(100) | UQ, NOT NULL | 权限代码 |
| Description | VARCHAR(255) | - | 权限描述 |

#### 4.2.4 OAuth2Client（OAuth2 客户端）

| 字段 | 类型 | 约束 | 描述 |
|------|------|------|------|
| ID | UUID | PK | 客户端 ID |
| Secret | VARCHAR(255) | NOT NULL | 客户端密钥 |
| Name | VARCHAR(100) | NOT NULL | 客户端名称 |
| Domain | VARCHAR(500) | - | 允许的回调域名 |
| Data | JSONB | - | 扩展数据（包含 Public、UserID 等） |

#### 4.2.5 Organization（组织架构）

| 字段 | 类型 | 约束 | 描述 |
|------|------|------|------|
| ID | UUID | PK | 组织节点 ID |
| UserID | UUID | FK, NOT NULL | 关联用户 |
| SuperiorID | UUID | FK, 自引用 | 上级节点 |
| AncestorPath | VARCHAR(500) | - | 祖先路径（用于快速查询） |
| Depth | INT | 默认 0 | 层级深度 |
| OrgRole | VARCHAR(50) | - | 组织角色 |
| Class | VARCHAR(50) | - | 班级/部门 |

#### 4.2.6 Files（文件）

| 字段 | 类型 | 描述 |
|------|------|------|
| ID | UUID | 文件唯一标识 |
| StorageKey | VARCHAR(500) | 存储路径/键 |
| Bucket | VARCHAR(100) | 存储桶名称 |
| MimeType | VARCHAR(100) | MIME 类型 |
| SizeBytes | BIGINT | 文件大小（字节） |
| CreatorID | UUID | 上传者 |
| IsTemporary | BOOLEAN | 是否为临时文件 |

---

## 5. 非功能性需求

### 5.1 安全机制

#### 5.1.1 密码安全

| 机制 | 实现方式 | 代码位置 |
|------|----------|----------|
| 密码加盐 | 16 字节随机盐值 | `utilities/encrypt.go#GenerateSalt` |
| 密码哈希 | SHA256(salt + password) | `utilities/encrypt.go#HashPassword` |
| 密码验证 | 重新计算哈希并比对 | `utilities/encrypt.go#VerifyPassword` |

#### 5.1.2 JWT 令牌

| 配置项 | 默认值 | 描述 |
|--------|--------|------|
| Access Token 有效期 | 1 小时 (3600s) | 短期令牌 |
| Refresh Token 有效期 | 7 天 (604800s) | 长期令牌，用于刷新 |
| 签名算法 | HS256 (内部令牌) | 对称加密 |
| 签名算法 | RS256 (ID Token) | 非对称加密 |

#### 5.1.3 JWKS 密钥管理

| 配置项 | 默认值 | 描述 |
|--------|--------|------|
| RSA 密钥大小 | 2048 bits | 安全强度 |
| 密钥轮换间隔 | 24 小时 | 自动轮换周期 |
| 旧密钥保留期 | 1 小时 | Grace Period，确保旧令牌仍可验证 |

**密钥状态机**：
```
Active (当前使用) → Rotating (轮换中) → Expired (已过期，删除)
```

#### 5.1.4 认证中间件

**Authorization Header 格式**：
```
Authorization: Bearer <access_token>
```

**Token 提取与验证流程**：
1. 从 Header 提取 Bearer Token
2. 使用密钥解析 JWT
3. 验证签名和过期时间
4. 将用户信息注入 Context

### 5.2 性能考虑

#### 5.2.1 数据库连接池

| 配置项 | 值 | 描述 |
|--------|-----|------|
| MaxIdleConns | 10 | 最大空闲连接数 |
| MaxOpenConns | 100 | 最大打开连接数 |
| ConnMaxLifetime | 1 小时 | 连接最大存活时间 |

#### 5.2.2 Token 存储

- **Token 存储后端**：Redis（外部依赖）
- **存储库**：`github.com/go-oauth2/redis/v4`

#### 5.2.3 缓存策略

- JWKS 端点设置缓存头：`Cache-Control: public, max-age=3600`

### 5.3 可观测性

#### 5.3.1 日志记录

- 使用自定义 Logger（`utilities/logger.go`）
- 中间件记录请求日志（`middleware/logger.go`）
- 关键业务操作记录到 `operation_record` 表

#### 5.3.2 错误追踪

- OAuth 服务内部错误记录到 `error_records` 表
- 包含：UserID、ClientID、ErrorType、Message、Timestamp

#### 5.3.3 健康检查

- 端点：`GET /health`
- 用于负载均衡器或 Kubernetes 探针

### 5.4 优雅停机

系统支持优雅停机（Graceful Shutdown）：
- 监听 `SIGINT`、`SIGTERM` 信号
- 10 秒超时等待请求处理完成
- 关闭数据库连接、停止密钥轮换

---

## 6. 外部依赖

| 依赖 | 用途 | 状态 |
|------|------|------|
| PostgreSQL | 主数据库 | 必须 |
| Redis | OAuth2 Token 存储 | 必须 |
| 本地文件系统 | 文件上传存储 | 可选（可扩展为 S3） |

---

## 7. 配置项清单

| 环境变量 | 默认值 | 描述 |
|----------|--------|------|
| `PORT` | 8080 | 服务端口 |
| `GIN_MODE` | debug | 运行模式 (debug/release/test) |
| `DATABASE_URL` | - | PostgreSQL 连接字符串 |
| `REDIS_ADDR` | localhost:6379 | Redis 地址 |
| `REDIS_PASSWORD` | - | Redis 密码 |
| `REDIS_DB` | 0 | Redis 数据库编号 |
| `JWT_SECRET` | your-secret-key | JWT 签名密钥 |
| `ACCESS_TOKEN_EXPIRE` | 3600 | Access Token 过期时间（秒） |
| `REFRESH_TOKEN_EXPIRE` | 604800 | Refresh Token 过期时间（秒） |
| `OIDC_ISSUER` | http://localhost:8080 | OIDC 发行者 URL |
| `OIDC_KEY_ROTATION_INTERVAL` | 86400 | JWKS 密钥轮换间隔（秒） |
| `OIDC_KEY_SIZE` | 2048 | RSA 密钥大小 |
| `STORAGE_LOCAL_DIR` | ./uploads | 本地存储目录 |

---

## 8. 业务术语表 (Glossary)

| 术语 | 英文 | 描述 |
|------|------|------|
| StudentID | 学号 | 用户的唯一业务标识 |
| Access Token | 访问令牌 | 短期令牌，用于 API 认证 |
| Refresh Token | 刷新令牌 | 长期令牌，用于获取新的 Access Token |
| ID Token | 身份令牌 | OIDC 标准令牌，包含用户身份信息 |
| JWKS | JSON Web Key Set | 公钥集合，用于验证 ID Token |
| Scope | 作用域 | OAuth2 授权范围，决定返回的用户信息 |
| Grant Type | 授权类型 | OAuth2 授权模式（如 authorization_code） |
| RBAC | Role-Based Access Control | 基于角色的访问控制 |
| Issuer | 发行者 | OIDC 令牌的签发方标识 |

---

## 9. 待确认/外部依赖项

| 项目 | 状态 | 说明 |
|------|------|------|
| 头像上传功能 | TODO | 代码中有 `picture` claim，但标注为待实现 |
| 邮箱验证流程 | 未实现 | `EmailVerified` 字段存在，但无验证接口 |
| 密码重置功能 | 未实现 | 无相关 API |
| 用户状态切换 | 未暴露 API | 仅在数据模型中定义状态 |
| 操作记录查询 | 未暴露 API | 仅有数据模型，无查询接口 |
| S3 存储支持 | 外部依赖/待确认 | StorageService 可能支持扩展 |

---

*文档版本: 1.0*
*生成时间: 2026-01-21*
*基于代码分析自动生成*
