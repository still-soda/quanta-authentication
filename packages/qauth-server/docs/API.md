# Quanta Authentication API 文档

## 概述

Quanta Authentication 是一个基于 OAuth2.0 和 OpenID Connect (OIDC) 的身份认证授权服务。

**Base URL**: `http://localhost:8080` (开发环境)

**协议版本**: OAuth 2.0, OpenID Connect 1.0

---

## 目录

- [健康检查](#健康检查)
- [文件上传](#文件上传)
- [认证接口](#认证接口)
- [OAuth2.0 接口](#oauth20-接口)
- [OIDC 接口](#oidc-接口)
- [管理员接口](#管理员接口)
- [错误码说明](#错误码说明)

---

## 健康检查

### 检查服务状态

检查服务是否正常运行。

**请求**

```http
GET /health
```

**响应**

```json
{
  "status": "ok",
  "message": "Server is running"
}
```

**状态码**: `200 OK`

---

## 文件上传

### 上传文件

上传文件到服务器。

**请求**

```http
POST /upload
Content-Type: multipart/form-data
```

**参数**

| 名称 | 类型 | 必填 | 描述 |
|------|------|------|------|
| file | File | 是 | 要上传的文件 |

**响应**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "fileName": "uploaded_file_name.ext"
  }
}
```

---

## 认证接口

### 用户注册

创建新用户账号。

**请求**

```http
POST /auth/register
Content-Type: application/json
```

**请求体**

```json
{
  "student_id": "20240001",
  "password": "secure_password",
  "email": "student@example.com",
  "name": "张三"
}
```

**参数说明**

| 字段 | 类型 | 必填 | 描述 |
|------|------|------|------|
| student_id | string | 是 | 学号，唯一标识 |
| password | string | 是 | 用户密码 |
| email | string | 是 | 用户邮箱，必须是有效的邮箱格式 |
| name | string | 是 | 用户姓名 |

**响应**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": "uuid-string",
    "student_id": "20240001",
    "email": "student@example.com",
    "name": "张三",
    "created_at": "2026-01-17T10:00:00Z"
  }
}
```

**错误响应**

- `400` - 请求参数错误
- `409` - 用户已存在（学号或邮箱冲突）

---

### 用户登录

使用学号和密码登录，获取访问令牌。

**请求**

```http
POST /auth/login
Content-Type: application/json
```

**请求体**

```json
{
  "student_id": "20240001",
  "password": "secure_password"
}
```

**参数说明**

| 字段 | 类型 | 必填 | 描述 |
|------|------|------|------|
| student_id | string | 是 | 学号 |
| password | string | 是 | 密码 |

**响应**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "user": {
      "id": "uuid-string",
      "student_id": "20240001",
      "email": "student@example.com",
      "name": "张三"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

**错误响应**

- `400` - 请求参数错误
- `401` - 认证失败（学号或密码错误）

---

### 刷新令牌

使用刷新令牌获取新的访问令牌和刷新令牌。

**请求**

```http
POST /auth/refresh-token
Content-Type: application/json
```

**请求体**

```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**参数说明**

| 字段 | 类型 | 必填 | 描述 |
|------|------|------|------|
| refresh_token | string | 是 | 刷新令牌 |

**响应**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

**错误响应**

- `400` - 请求参数错误
- `401` - 刷新令牌无效或已过期

---

## OAuth2.0 接口

### 授权端点

处理 OAuth2.0 授权请求（授权码模式）。

**请求**

```http
GET /oauth/authorize
```

或

```http
POST /oauth/authorize
```

**查询参数**

| 参数 | 类型 | 必填 | 描述 |
|------|------|------|------|
| client_id | string | 是 | 客户端 ID |
| redirect_uri | string | 是 | 回调 URI |
| response_type | string | 是 | 响应类型，通常为 `code` |
| scope | string | 否 | 请求的权限范围 |
| state | string | 推荐 | 防止 CSRF 攻击的随机字符串 |

**成功响应**

重定向到 `redirect_uri`，并附带授权码：

```
{redirect_uri}?code={authorization_code}&state={state}
```

**错误响应**

重定向到 `redirect_uri`，并附带错误信息：

```
{redirect_uri}?error={error_code}&error_description={description}&state={state}
```

---

### 令牌端点

使用授权码换取访问令牌。

**请求**

```http
POST /oauth/token
Content-Type: application/x-www-form-urlencoded
```

**请求体参数**

| 参数 | 类型 | 必填 | 描述 |
|------|------|------|------|
| grant_type | string | 是 | 授权类型：`authorization_code`、`refresh_token`、`client_credentials` |
| code | string | 条件 | 授权码（grant_type 为 authorization_code 时必填） |
| redirect_uri | string | 条件 | 回调 URI（grant_type 为 authorization_code 时必填） |
| client_id | string | 是 | 客户端 ID |
| client_secret | string | 是 | 客户端密钥 |
| refresh_token | string | 条件 | 刷新令牌（grant_type 为 refresh_token 时必填） |

**响应**

```json
{
  "access_token": "eyJhbGciOiJSUzI1NiIs...",
  "token_type": "Bearer",
  "expires_in": 3600,
  "refresh_token": "eyJhbGciOiJSUzI1NiIs...",
  "scope": "openid profile email"
}
```

---

### 验证令牌

验证访问令牌的有效性。

**请求**

```http
POST /oauth/validate
Authorization: Bearer {access_token}
```

**响应**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "active": true,
    "client_id": "client_id_string",
    "user_id": "user_id_string",
    "expires_in": 3600,
    "scope": "openid profile"
  }
}
```

**错误响应**

- `401` - 令牌无效或已过期

---

### 撤销令牌

撤销访问令牌或刷新令牌。

**请求**

```http
POST /oauth/revoke
Content-Type: application/json
```

**请求体**

```json
{
  "token": "token_to_revoke",
  "token_type_hint": "access_token"
}
```

**参数说明**

| 字段 | 类型 | 必填 | 描述 |
|------|------|------|------|
| token | string | 是 | 要撤销的令牌 |
| token_type_hint | string | 否 | 令牌类型提示：`access_token` 或 `refresh_token` |

**响应**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "revoked": true
  }
}
```

---

### OAuth2 客户端管理

#### 创建客户端

创建新的 OAuth2 客户端。

**请求**

```http
POST /oauth/clients
Content-Type: application/json
```

**请求体**

```json
{
  "name": "My Application",
  "domain": "https://example.com",
  "secret": "client_secret_string",
  "is_public": false,
  "user_id": "owner_user_id"
}
```

**参数说明**

| 字段 | 类型 | 必填 | 描述 |
|------|------|------|------|
| name | string | 是 | 客户端名称 |
| domain | string | 是 | 客户端域名 |
| secret | string | 是 | 客户端密钥 |
| is_public | boolean | 否 | 是否为公开客户端 |
| user_id | string | 否 | 所有者用户 ID |

**响应**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "client_id": "generated_client_id",
    "client_name": "My Application",
    "client_domain": "https://example.com"
  }
}
```

---

#### 获取客户端信息

获取指定客户端的详细信息。

**请求**

```http
GET /oauth/clients/:id
```

**路径参数**

| 参数 | 类型 | 必填 | 描述 |
|------|------|------|------|
| id | string | 是 | 客户端 ID |

**响应**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "client_id": "client_id_string",
    "client_name": "My Application",
    "client_domain": "https://example.com",
    "created_at": "2026-01-17T10:00:00Z"
  }
}
```

**错误响应**

- `404` - 客户端不存在

---

#### 更新客户端信息

更新指定客户端的信息。

**请求**

```http
PUT /oauth/clients/:id
Content-Type: application/json
```

**路径参数**

| 参数 | 类型 | 必填 | 描述 |
|------|------|------|------|
| id | string | 是 | 客户端 ID |

**请求体**

```json
{
  "name": "Updated Application Name",
  "domain": "https://new-domain.com"
}
```

**响应**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "updated": true
  }
}
```

---

#### 删除客户端

删除指定的 OAuth2 客户端。

**请求**

```http
DELETE /oauth/clients/:id
```

**路径参数**

| 参数 | 类型 | 必填 | 描述 |
|------|------|------|------|
| id | string | 是 | 客户端 ID |

**响应**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "deleted": true
  }
}
```

---

#### 获取客户端列表

获取 OAuth2 客户端列表（分页）。

**请求**

```http
GET /oauth/clients
```

**查询参数**

| 参数 | 类型 | 必填 | 默认值 | 描述 |
|------|------|------|--------|------|
| page | integer | 否 | 1 | 页码 |
| page_size | integer | 否 | 10 | 每页数量（最大 100） |

**响应**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "items": [
      {
        "client_id": "client_id_1",
        "client_name": "App 1",
        "client_domain": "https://app1.com"
      },
      {
        "client_id": "client_id_2",
        "client_name": "App 2",
        "client_domain": "https://app2.com"
      }
    ],
    "total": 25,
    "page": 1,
    "size": 10
  }
}
```

---

## OIDC 接口

### OpenID Connect 配置发现

获取 OpenID Connect 提供者的配置信息。

**请求**

```http
GET /.well-known/openid-configuration
```

**响应**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "issuer": "http://localhost:8080",
    "authorization_endpoint": "http://localhost:8080/oauth/authorize",
    "token_endpoint": "http://localhost:8080/oauth/token",
    "userinfo_endpoint": "http://localhost:8080/oauth/userinfo",
    "jwks_uri": "http://localhost:8080/.well-known/jwks.json",
    "response_types_supported": ["code", "token", "id_token"],
    "subject_types_supported": ["public"],
    "id_token_signing_alg_values_supported": ["RS256"],
    "scopes_supported": ["openid", "profile", "email"],
    "token_endpoint_auth_methods_supported": ["client_secret_basic", "client_secret_post"],
    "claims_supported": ["sub", "name", "email", "email_verified"]
  }
}
```

---

### JWKS 端点

获取用于验证 JWT 签名的公钥集合。

**请求**

```http
GET /.well-known/jwks.json
```

**响应头**

```
Cache-Control: public, max-age=3600
Content-Type: application/json
```

**响应**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "keys": [
      {
        "kty": "RSA",
        "use": "sig",
        "kid": "key_id_1",
        "n": "modulus_value",
        "e": "AQAB",
        "alg": "RS256"
      }
    ]
  }
}
```

---

### 用户信息端点

获取当前用户的详细信息（需要有效的访问令牌）。

**请求**

```http
GET /oauth/userinfo
Authorization: Bearer {access_token}
```

**响应**

```json
{
  "code": 501,
  "message": "userinfo endpoint is not implemented yet"
}
```

**注意**: 此端点当前未实现，将在后续版本中提供。

---

### 登出端点

用户登出并可选择重定向到指定 URI。

**请求**

```http
GET /oauth/logout?post_logout_redirect_uri=https://example.com
```

或

```http
POST /oauth/logout
```

**查询参数**

| 参数 | 类型 | 必填 | 描述 |
|------|------|------|------|
| post_logout_redirect_uri | string | 否 | 登出后重定向的 URI，默认为 `/` |

**响应**

重定向到指定的 `post_logout_redirect_uri`。

**状态码**: `302 Found`

---

## 管理员接口

### 获取密钥轮换信息

获取当前所有 JWKS 密钥的状态信息。

**请求**

```http
GET /admin/jwks/keys
```

**响应**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "active_key_id": "current_active_key_id",
    "keys": [
      {
        "id": "key_id_1",
        "status": "active",
        "created_at": "2026-01-17T10:00:00Z",
        "expires_at": "2026-02-17T10:00:00Z"
      },
      {
        "id": "key_id_2",
        "status": "retired",
        "created_at": "2025-12-17T10:00:00Z",
        "expires_at": "2026-01-17T10:00:00Z"
      }
    ],
    "total": 2
  }
}
```

---

### 强制密钥轮换

立即执行密钥轮换操作。

**请求**

```http
POST /admin/jwks/rotate
```

**响应**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "message": "key rotation completed",
    "active_key_id": "new_active_key_id"
  }
}
```

---

## 错误码说明

### 标准 HTTP 状态码

| 状态码 | 说明 |
|--------|------|
| 200 | 请求成功 |
| 400 | 请求参数错误 |
| 401 | 未授权或认证失败 |
| 404 | 资源不存在 |
| 409 | 资源冲突 |
| 500 | 服务器内部错误 |
| 501 | 功能未实现 |

### 错误响应格式

```json
{
  "code": 400,
  "message": "具体的错误描述",
  "data": null
}
```

### 常见错误码

| 错误码 | 错误信息 | 描述 |
|--------|----------|------|
| 400 | Bad Request | 请求参数错误 |
| 401 | Unauthorized | 未授权或认证失败 |
| 404 | Not Found | 资源不存在 |
| 409 | Conflict | 用户已存在或资源冲突 |
| 500 | Internal Server Error | 服务器内部错误 |

---

## 认证流程示例

### 授权码模式流程

1. **用户授权请求**

   应用将用户重定向到授权端点：

   ```
   GET /oauth/authorize?client_id={client_id}&redirect_uri={redirect_uri}&response_type=code&scope=openid%20profile&state={random_state}
   ```

2. **用户登录并授权**

   用户在认证服务器上登录并同意授权。

3. **获取授权码**

   认证服务器重定向回应用，附带授权码：

   ```
   {redirect_uri}?code={authorization_code}&state={state}
   ```

4. **交换访问令牌**

   应用使用授权码换取访问令牌：

   ```http
   POST /oauth/token
   Content-Type: application/x-www-form-urlencoded
   
   grant_type=authorization_code&code={authorization_code}&redirect_uri={redirect_uri}&client_id={client_id}&client_secret={client_secret}
   ```

5. **获取访问令牌**

   认证服务器返回访问令牌：

   ```json
   {
     "access_token": "...",
     "token_type": "Bearer",
     "expires_in": 3600,
     "refresh_token": "...",
     "id_token": "..."
   }
   ```

6. **使用访问令牌访问资源**

   应用使用访问令牌访问受保护的资源。

---

## 附录

### Token 格式

本系统使用 JWT (JSON Web Token) 作为令牌格式。

**访问令牌 (Access Token)**:

- 类型: JWT
- 签名算法: RS256
- 有效期: 1 小时（可配置）

**刷新令牌 (Refresh Token)**:

- 类型: JWT
- 签名算法: RS256
- 有效期: 30 天（可配置）

**ID Token** (OIDC):

- 类型: JWT
- 签名算法: RS256
- 包含用户身份信息

### 权限范围 (Scopes)

| Scope | 描述 |
|-------|------|
| openid | 必须包含，表示 OIDC 请求 |
| profile | 访问用户基本资料 |
| email | 访问用户邮箱 |

### 安全建议

1. **使用 HTTPS**: 生产环境必须使用 HTTPS 协议
2. **保护客户端密钥**: 客户端密钥应安全存储，不要暴露在客户端代码中
3. **使用 state 参数**: 防止 CSRF 攻击
4. **验证重定向 URI**: 确保重定向 URI 在白名单内
5. **令牌过期策略**: 定期刷新令牌，及时撤销不再使用的令牌

---

## 更新日志

### v1.0.0 (2026-01-17)

- 初始版本
- 实现基础认证功能
- 实现 OAuth2.0 授权码模式
- 实现 OIDC 发现和 JWKS 端点
- 实现客户端管理功能
- 实现密钥轮换功能

---

## 联系方式

如有问题或建议，请联系开发团队。
