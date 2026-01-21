# QAuth Server 测试架构分析文档

> **文档版本**: 1.0  
> **分析基于**: PRD v1.0  
> **分析日期**: 2026-01-21  
> **角色**: QA Architect（测试架构师）

---

## 1. 核心业务目标与价值链

### 1.1 业务目标金字塔

```mermaid
graph TB
    subgraph Strategic [战略目标 Strategic]
        S[统一企业/学校身份认证体系]
    end
    
    subgraph Business [业务目标 Business]
        B1[标准协议接入<br/>OAuth2/OIDC]
        B2[细粒度权限管控<br/>RBAC]
        B3[安全合规保障<br/>密钥轮换/哈希]
    end
    
    subgraph Operational [运营目标 Operational]
        O1[高可用服务<br/>健康检查/优雅停机]
        O2[可观测性<br/>日志/错误追踪]
        O3[可扩展性<br/>存储/Token后端]
    end
    
    S --> B1
    S --> B2
    S --> B3
    B1 --> O1
    B2 --> O2
    B3 --> O3
    
    style Strategic fill:#e1f5fe
    style Business fill:#fff3e0
    style Operational fill:#e8f5e9
```

### 1.2 价值链分析

| 价值阶段 | 核心活动 | 关键能力 | 测试关注点 |
|---------|---------|---------|-----------|
| **身份建立** | 用户注册 | 唯一性校验、密码安全存储 | 边界条件、冲突处理、安全性 |
| **身份验证** | 用户登录、OAuth授权 | 凭证验证、状态检查 | 正向流程、异常流程、状态机 |
| **令牌颁发** | JWT生成、Token管理 | 签名算法、过期机制 | 令牌完整性、时效性、安全性 |
| **权限控制** | RBAC校验 | 角色-权限映射 | 权限边界、越权测试 |
| **协议支持** | OIDC发现、JWKS | 标准协议合规 | 协议兼容性、互操作性 |
| **安全保障** | 密钥轮换、审计日志 | 生命周期管理 | 轮换正确性、日志完整性 |

### 1.3 关键业务指标 (KPIs)

| 指标类别 | 指标名称 | 测试验证方式 |
|---------|---------|-------------|
| **可用性** | 服务可用率 | 健康检查端点监控 |
| **性能** | 登录响应时间 | 负载测试 |
| **安全性** | 令牌有效性 | 过期/撤销测试 |
| **合规性** | OIDC标准符合度 | 协议一致性测试 |
| **可靠性** | 密钥轮换成功率 | 轮换流程测试 |

---

## 2. 关键参与角色 (Actors) 识别

### 2.1 角色全景图

```mermaid
graph TB
    subgraph 外部角色 [外部参与者]
        EU[终端用户<br/>End User]
        DEV[开发者<br/>Developer]
        ADMIN[管理员<br/>Administrator]
        CLIENT[第三方应用<br/>OAuth2 Client]
    end
    
    subgraph 系统角色 [系统内部角色]
        IDP[身份提供者<br/>Identity Provider]
        AS[授权服务器<br/>Authorization Server]
        RS[资源服务器<br/>Resource Server]
    end
    
    subgraph 基础设施 [基础设施]
        DB[(PostgreSQL)]
        REDIS[(Redis)]
        FS[文件存储]
    end
    
    EU -->|注册/登录| IDP
    EU -->|授权确认| AS
    DEV -->|接入OAuth2| CLIENT
    ADMIN -->|管理配置| RS
    CLIENT -->|请求授权| AS
    CLIENT -->|获取用户信息| RS
    
    IDP --> DB
    AS --> REDIS
    RS --> FS
```

### 2.2 角色详细定义

| Actor ID | 角色名称 | 角色描述 | 主要交互 | 测试场景关注点 |
|----------|---------|---------|---------|---------------|
| **ACT-01** | 终端用户 (End User) | 通过学号/工号进行身份认证的最终用户 | 注册、登录、授权确认、Token刷新 | 用户旅程完整性、状态转换 |
| **ACT-02** | 系统管理员 (Admin) | 具有 `system_admin` 或 `system_super_admin` 角色的管理用户 | 用户管理、角色管理、权限分配、客户端管理 | 权限边界、越权防护 |
| **ACT-03** | 开发者 (Developer) | 接入OAuth2/OIDC协议的第三方应用开发者 | 客户端注册、OIDC发现、JWKS获取 | 协议兼容性、错误处理 |
| **ACT-04** | OAuth2客户端 (Client) | 第三方应用系统 | 授权码获取、Token交换、用户信息获取 | 授权流程、令牌生命周期 |
| **ACT-05** | 身份提供者 (IdP) | QAuth Server 作为身份提供者 | 身份验证、ID Token签发 | OIDC协议合规性 |
| **ACT-06** | 定时任务 (Scheduler) | 系统内部定时任务（密钥轮换） | 密钥生成、密钥状态切换 | 轮换正确性、Grace Period |

### 2.3 角色权限矩阵

| 功能模块 | 终端用户 | 系统管理员 | 超级管理员 | OAuth2客户端 |
|---------|:-------:|:---------:|:---------:|:-----------:|
| 用户注册 | ✅ | ✅ | ✅ | ❌ |
| 用户登录 | ✅ | ✅ | ✅ | ❌ |
| Token刷新 | ✅ | ✅ | ✅ | ✅ |
| OAuth授权 | ✅ | ✅ | ✅ | ❌ |
| 客户端管理 | ❌ | ✅ | ✅ | ❌ |
| 角色管理 | ❌ | ⚠️ 部分 | ✅ | ❌ |
| 权限分配 | ❌ | ⚠️ 部分 | ✅ | ❌ |
| 密钥管理 | ❌ | ❌ | ✅ | ❌ |
| 用户信息获取 | ❌ | ❌ | ❌ | ✅ (授权后) |

---

## 3. 系统核心模块及依赖关系

### 3.1 模块架构图

```mermaid
graph TB
    subgraph 接入层 [接入层 - Entry Layer]
        GW[API Gateway / Routes]
        MW[中间件层<br/>Auth | CORS | Logger | Recovery]
    end
    
    subgraph 业务层 [业务层 - Business Layer]
        AUTH[认证模块<br/>Auth Handler]
        OAUTH[OAuth2模块<br/>OAuth Handler]
        OIDC_H[OIDC模块<br/>OIDC Handler]
        ROLE[角色模块<br/>Role Handler]
        PERM[权限模块<br/>Permission Handler]
        FILE[文件模块<br/>File Handler]
    end
    
    subgraph 服务层 [服务层 - Service Layer]
        USER_SVC[用户服务<br/>UserService]
        OAUTH_SVC[OAuth服务<br/>OAuthService]
        OIDC_SVC[OIDC服务<br/>OIDCService]
        ROLE_SVC[角色服务<br/>RoleService]
        PERM_SVC[权限服务<br/>PermissionService]
        FILE_SVC[文件服务<br/>FileService]
        STORAGE_SVC[存储服务<br/>StorageService]
    end
    
    subgraph 基础设施层 [基础设施层 - Infrastructure Layer]
        DB_MOD[数据库模块<br/>Database]
        JWT_MOD[JWT模块<br/>JWT Utils]
        JWKS_MOD[JWKS模块<br/>Key Manager]
        ENCRYPT[加密模块<br/>Encrypt Utils]
    end
    
    subgraph 外部依赖 [外部依赖 - External Dependencies]
        PG[(PostgreSQL)]
        REDIS[(Redis)]
        FS[文件系统/S3]
    end
    
    GW --> MW
    MW --> AUTH & OAUTH & OIDC_H & ROLE & PERM & FILE
    
    AUTH --> USER_SVC
    OAUTH --> OAUTH_SVC
    OIDC_H --> OIDC_SVC
    ROLE --> ROLE_SVC
    PERM --> PERM_SVC
    FILE --> FILE_SVC
    
    USER_SVC --> DB_MOD & ENCRYPT
    OAUTH_SVC --> DB_MOD & JWT_MOD & JWKS_MOD
    OIDC_SVC --> JWKS_MOD
    ROLE_SVC --> DB_MOD
    PERM_SVC --> DB_MOD
    FILE_SVC --> STORAGE_SVC
    
    DB_MOD --> PG
    OAUTH_SVC --> REDIS
    STORAGE_SVC --> FS
```

### 3.2 模块详细定义

| 模块ID | 模块名称 | 职责描述 | 入口端点 | 上游依赖 | 下游依赖 |
|--------|---------|---------|---------|---------|---------|
| **MOD-01** | 认证模块 (Auth) | 用户注册、登录、令牌刷新 | `/_/v1/auth/*` | Routes, Middleware | UserService, JWT |
| **MOD-02** | OAuth2模块 | OAuth2授权码流程、令牌管理 | `/v1/oauth/*` | Routes, Middleware | OAuthService, Redis |
| **MOD-03** | OIDC模块 | OIDC发现文档、JWKS端点 | `/.well-known/*` | Routes | OIDCService, JWKS |
| **MOD-04** | 角色模块 | 角色CRUD | `/_/v1/roles/*` | Auth Middleware | RoleService, DB |
| **MOD-05** | 权限模块 | 权限授予/撤销 | `/_/v1/permissions/*` | Auth Middleware | PermissionService, DB |
| **MOD-06** | 客户端模块 | OAuth2客户端管理 | `/_/v1/clients/*` | Auth Middleware | OAuthService, DB |
| **MOD-07** | 文件模块 | 文件上传 | `/_/v1/resources/*` | Auth Middleware | FileService, Storage |
| **MOD-08** | 密钥模块 | JWKS密钥管理 | `/_/v1/jwks/*` | Auth Middleware | JWKS Manager |
| **MOD-09** | 健康检查 | 服务健康探测 | `/health` | Routes | - |

### 3.3 模块依赖矩阵

| 模块 | Auth | OAuth2 | OIDC | Role | Permission | Client | File | JWKS |
|------|:----:|:------:|:----:|:----:|:----------:|:------:|:----:|:----:|
| **Auth** | - | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| **OAuth2** | ✅ | - | ✅ | ✅ | ✅ | ❌ | ❌ | ✅ |
| **OIDC** | ❌ | ❌ | - | ❌ | ❌ | ❌ | ❌ | ✅ |
| **Role** | ✅ | ❌ | ❌ | - | ✅ | ❌ | ❌ | ❌ |
| **Permission** | ✅ | ❌ | ❌ | ✅ | - | ❌ | ❌ | ❌ |
| **Client** | ✅ | ❌ | ❌ | ❌ | ✅ | - | ❌ | ❌ |
| **File** | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | - | ❌ |
| **JWKS** | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | - |

> **图例**: ✅ = 依赖 | ❌ = 无依赖

### 3.4 关键依赖路径 (Critical Path)

```mermaid
flowchart LR
    subgraph 用户登录流程
        L1[Routes] --> L2[Auth Middleware] --> L3[Auth Handler] --> L4[UserService] --> L5[Database] --> L6[JWT Utils]
    end
```

```mermaid
flowchart LR
    subgraph OAuth2授权流程
        O1[Routes] --> O2[OAuth Handler] --> O3[OAuthService] --> O4[UserService + Database + Redis + JWKS]
    end
```

```mermaid
flowchart LR
    subgraph 权限校验流程
        P1[Auth Middleware] --> P2[JWT Utils] --> P3[PermissionService] --> P4[Database]
    end
```

---

## 4. 业务实体关系图

### 4.1 核心实体关系 (ER Diagram)

```mermaid
erDiagram
    Users ||--o{ UsersRoles : "拥有"
    Users ||--o{ LoginState : "产生"
    Users ||--o{ Organization : "隶属"
    Users ||--o| Images : "头像"
    Users ||--o{ Files : "上传"
    Users ||--o{ OAuth2Client : "创建"
    
    Roles ||--o{ UsersRoles : "分配给"
    Roles ||--o{ RolesPermissions : "包含"
    
    Permissions ||--o{ RolesPermissions : "授予"
    
    OAuth2Client ||--o{ LoginState : "触发"
    
    Organization ||--o| Organization : "上级"
    
    Images ||--|| Files : "关联"
    
    Users {
        uuid ID PK
        varchar StudentID UK "学号-业务主键"
        varchar Email UK
        varchar PasswordHash
        varchar Salt
        varchar Name
        varchar Phone UK
        varchar DisplayName
        uuid AvatarID FK
        enum Status "ACTIVE|LOCKED|BANNED"
        boolean EmailVerified
    }
    
    Roles {
        uuid ID PK
        varchar Code UK "角色代码"
        varchar Name
        varchar Description
        boolean IsSystem "系统内置"
    }
    
    Permissions {
        uuid ID PK
        varchar Resource "资源名称"
        int8 Action "0=C,1=R,2=U,3=D"
        varchar Code UK "权限代码"
        varchar Description
    }
    
    UsersRoles {
        uuid UsersID FK
        uuid RolesID FK
        timestamp AssignedAt
    }
    
    RolesPermissions {
        uuid RolesID FK
        uuid PermissionsID FK
    }
    
    OAuth2Client {
        uuid ID PK "ClientID"
        varchar Secret
        varchar Name
        varchar Domain "回调域名"
        jsonb Data "扩展数据"
    }
    
    LoginState {
        uuid ID PK
        uuid UserID FK
        uuid ClientID FK
        timestamp Time
        enum Type "PASSWORD|OAUTH2|REFRESH"
        varchar IP
        varchar UserAgent
        varchar Location
        boolean IsSuccess
        varchar FailReason
    }
    
    Organization {
        uuid ID PK
        uuid UserID FK
        uuid SuperiorID FK "自引用"
        varchar AncestorPath
        int Depth
        varchar OrgRole
        varchar Class
    }
    
    Files {
        uuid ID PK
        varchar StorageKey
        varchar Bucket
        varchar MimeType
        bigint SizeBytes
        uuid CreatorID FK
        boolean IsTemporary
    }
    
    Images {
        uuid ID PK
        int Width
        int Height
        uuid FileID FK
        uuid CreatorID FK
    }
```

### 4.2 实体分类与测试重点

| 实体类别 | 实体名称 | 核心约束 | 测试重点 |
|---------|---------|---------|---------|
| **核心实体** | Users | StudentID/Email/Phone唯一 | 唯一性冲突、状态机转换 |
| **核心实体** | Roles | Code唯一、IsSystem不可删除 | 系统角色保护 |
| **核心实体** | Permissions | Code唯一 | 权限码完整性 |
| **关联实体** | UsersRoles | 复合主键 | 多角色分配 |
| **关联实体** | RolesPermissions | 复合主键 | 权限继承 |
| **业务实体** | OAuth2Client | Secret安全性 | 客户端认证 |
| **审计实体** | LoginState | 历史记录不可篡改 | 审计完整性 |
| **层级实体** | Organization | 自引用树结构 | 层级遍历、路径正确性 |
| **资源实体** | Files / Images | 外键完整性 | 文件关联、孤儿清理 |

### 4.3 状态机模型

#### 4.3.1 用户状态机 (User Status)

```mermaid
stateDiagram-v2
    [*] --> ACTIVE : 注册成功
    
    ACTIVE --> LOCKED : 管理员锁定
    ACTIVE --> BANNED : 永久封禁
    
    LOCKED --> ACTIVE : 管理员解锁
    LOCKED --> BANNED : 升级为封禁
    
    BANNED --> [*] : 终态(不可恢复)
    
    note right of ACTIVE : 可正常登录/授权
    note right of LOCKED : 无法登录/授权
    note right of BANNED : 永久禁止访问
```

#### 4.3.2 密钥状态机 (JWKS Key Status)

```mermaid
stateDiagram-v2
    [*] --> Active : 生成新密钥
    
    Active --> Rotating : 触发轮换
    
    Rotating --> Expired : Grace Period结束
    
    Expired --> [*] : 删除密钥
    
    note right of Active : 用于签名新Token
    note right of Rotating : 不签名新Token<br/>仍可验证旧Token
    note right of Expired : 准备删除
```

#### 4.3.3 OAuth2授权码状态机

```mermaid
stateDiagram-v2
    [*] --> Generated : 生成授权码
    
    Generated --> Exchanged : 兑换Token
    Generated --> Expired : 超时未使用
    
    Exchanged --> [*] : 一次性使用
    Expired --> [*] : 作废
    
    note right of Generated : code有效期短(通常10分钟)
    note right of Exchanged : 单次有效,防重放
```

---

## 5. 测试架构建议

### 5.1 测试金字塔

```mermaid
%%{init: {'theme': 'base', 'themeVariables': { 'fontSize': '14px'}}}%%
graph TB
    subgraph pyramid [测试金字塔]
        E2E["🔺 E2E Tests (5-10%)<br/>少量关键业务流程"]
        INT["🔶 Integration Tests (20-30%)<br/>API契约、模块集成"]
        UNIT["🟩 Unit Tests (60-70%)<br/>业务逻辑、边界条件"]
    end
    
    E2E --> INT --> UNIT
    
    style E2E fill:#ffcdd2,stroke:#c62828
    style INT fill:#fff9c4,stroke:#f9a825
    style UNIT fill:#c8e6c9,stroke:#2e7d32
```

### 5.2 关键测试场景识别

| 优先级 | 测试场景 | 测试类型 | 验证点 |
|:------:|---------|---------|-------|
| P0 | OAuth2授权码流程 | E2E | 完整授权流程、Token有效性 |
| P0 | 用户登录/Token刷新 | Integration | JWT签发、过期验证 |
| P0 | RBAC权限校验 | Integration | 角色-权限映射正确性 |
| P1 | 用户状态机转换 | Unit | 状态约束、非法转换拒绝 |
| P1 | 密钥轮换 | Integration | 新旧密钥共存期验证 |
| P1 | 唯一性约束 | Unit | 冲突检测、错误码 |
| P2 | OIDC协议合规性 | Contract | Discovery文档、Claims |
| P2 | 审计日志完整性 | Integration | LoginState记录 |

### 5.3 测试环境依赖

| 依赖项 | 测试环境方案 | 备注 |
|-------|-------------|------|
| PostgreSQL | Testcontainers / Docker | 数据隔离 |
| Redis | Testcontainers / Docker | Token存储 |
| 文件系统 | 内存文件系统 / 临时目录 | 隔离清理 |

---

## 6. 风险与待澄清项

### 6.1 架构风险点

| 风险ID | 风险描述 | 影响评估 | 测试缓解措施 |
|--------|---------|---------|-------------|
| R-01 | 密钥轮换期间新旧密钥验证逻辑 | 高 | Grace Period边界测试 |
| R-02 | 用户状态切换无暴露API | 中 | 数据库直接操作测试 |
| R-03 | Organization树结构性能 | 中 | 深层级遍历压测 |
| R-04 | Token存储Redis单点依赖 | 高 | 故障注入测试 |

### 6.2 待确认项 (从PRD继承)

| 项目 | 测试影响 | 建议 |
|------|---------|------|
| 头像上传功能 (TODO) | 无法测试picture claim | 待开发后补充 |
| 邮箱验证流程 (未实现) | EmailVerified字段测试受限 | Mock或跳过 |
| 密码重置功能 (未实现) | 安全性测试缺失 | 标记为已知Gap |
| 用户状态切换API (未暴露) | 状态机测试需绕行 | 直接DB操作 |

---

*文档版本: 1.0*  
*创建日期: 2026-01-21*  
*作者: QA Architect (AI-assisted)*
