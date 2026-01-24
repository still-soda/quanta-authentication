# 类型定义目录

本目录包含项目的所有 TypeScript 类型定义，按功能模块进行分类组织。

## 文件结构

```
types/
├── index.ts           # 统一导出所有类型
├── common.ts          # 通用类型定义
├── auth.ts            # 认证相关类型
├── user.ts            # 用户相关类型
├── role.ts            # 角色和权限相关类型
├── oauth.ts           # OAuth 相关类型
├── audit.ts           # 审计日志相关类型
├── dashboard.ts       # 仪表盘相关类型
├── notification.ts    # 通知相关类型
├── organization.ts    # 组织架构相关类型
├── profile.ts         # 个人资料相关类型
└── settings.ts        # 系统设置相关类型
```

## 模块说明

### common.ts

通用类型定义，包括：

- 主题色类型 (`ThemeColor`, `ExtendedThemeColor`)
- 通用状态类型 (`Status`)
- 搜索相关类型 (`SearchCategory`, `SearchItem`, `SearchGroup`)
- 动画数字配置 (`AnimatedNumberOptions`)

### auth.ts

认证相关类型，包括：

- 认证用户信息 (`AuthUser`)
- 登录请求/响应 (`LoginRequest`, `LoginResponse`)
- Token 刷新 (`RefreshTokenRequest`, `RefreshTokenResponse`)
- 注册请求 (`RegisterRequest`)
- 忘记密码请求 (`ForgotPasswordRequest`)

### user.ts

用户管理相关类型，包括：

- 用户信息 (`User`)
- 用户表单数据 (`UserFormData`)

### role.ts

角色和权限管理相关类型，包括：

- 角色信息 (`Role`)
- 角色表单数据 (`RoleFormData`)
- 权限定义 (`Permission`)
- 权限组 (`PermissionGroup`)
- 用户角色 (`UserRole`)

### oauth.ts

OAuth 相关类型，包括：

- OAuth 应用 (`OAuthApp`)
- OAuth 应用表单数据 (`OAuthAppFormData`)
- OAuth 设置 (`OAuthSettings`)

### audit.ts

审计日志相关类型，包括：

- 审计模块枚举 (`AuditModule`)
- 审计操作枚举 (`AuditAction`)
- 审计状态枚举 (`AuditStatus`)
- 审计日志 (`AuditLog`)
- 审计日志查询过滤器 (`AuditLogFilter`)
- 审计活动 (`AuditActivity`)
- 审计统计 (`AuditStats`)
- 热门客户端 (`TopClient`)
- 后端响应类型 (`AuditStatsResponse`, `RecentActivityResponse`, `TopClientResponse`)

### dashboard.ts

仪表盘相关类型，包括：

- 统计卡片数据 (`StatCardData`, `SimpleStatData`, `MiniStatItem`)
- 活动记录 (`Activity`)
- 热门应用 (`TopApp`)
- Dashboard 响应类型 (`DashboardStats`, `AuthTrendData`, `UserDistData`)
- 后端响应类型 (`DashboardStatsResponse`)

### notification.ts

通知相关类型，包括：

- 通知类型枚举 (`NotificationType`)
- 通知信息 (`Notification`)
- 通知过滤器 (`NotificationFilter`)

### organization.ts

组织架构相关类型，包括：

- 组织节点数据 (`OrgNodeData`)
- 组织节点 (`OrgNode`)
- 组织成员表单数据 (`OrgMemberFormData`)

### profile.ts

个人资料相关类型，包括：

- 个人资料数据 (`ProfileData`)
- 个人资料更新数据 (`ProfileUpdateData`)
- 登录记录 (`LoginRecord`)
- 安全设置 (`SecuritySettings`)
- 密码修改数据 (`PasswordChangeData`)

### settings.ts

系统设置相关类型，包括：

- 设置组 (`SettingGroup`)
- 通用设置 (`GeneralSettings`)
- 安全设置配置 (`SecuritySettingsConfig`)
- 邮件设置 (`EmailSettings`)
- 存储设置 (`StorageSettings`)
- 所有设置 (`AllSettings`)

## 使用方式

### 导入类型

从 `@/types` 导入所需的类型：

```typescript
import type { User, Role, Permission } from '@/types'
import type { LoginRequest, LoginResponse } from '@/types'
import type { DashboardStats, Activity } from '@/types'
```

### 导入特定模块的类型

如果只需要特定模块的类型，也可以从具体文件导入：

```typescript
import type { User, UserFormData } from '@/types/user'
import type { AuthUser, LoginRequest } from '@/types/auth'
import type { OAuthApp } from '@/types/oauth'
```

## 命名规范

### 接口命名

- 数据模型：使用名词，如 `User`, `Role`, `OAuthApp`
- 表单数据：添加 `FormData` 后缀，如 `UserFormData`, `RoleFormData`
- 请求数据：添加 `Request` 后缀，如 `LoginRequest`, `RegisterRequest`
- 响应数据：添加 `Response` 后缀，如 `LoginResponse`, `AuditStatsResponse`
- 过滤器：添加 `Filter` 后缀，如 `AuditLogFilter`, `NotificationFilter`
- 设置配置：添加 `Settings` 或 `Config` 后缀，如 `SecuritySettings`, `OAuthSettings`

### 枚举命名

- 使用 `type` 定义字符串联合类型
- 使用大写形式，如 `'SUCCESS' | 'WARNING' | 'ERROR'`
- 类型名使用名词，如 `AuditModule`, `AuditAction`, `AuditStatus`

## 注意事项

1. **后端数据格式**：某些类型（如 `AuthUser`, `AuditStatsResponse` 等）使用下划线命名（snake_case），这是为了与后端 API 返回格式保持一致。

2. **类型导出**：所有类型都通过 `index.ts` 统一导出，确保导入路径的一致性。

3. **类型复用**：相关联的类型定义在同一文件中，便于维护和理解上下文。

4. **文档注释**：每个文件开头都有注释说明该文件包含的类型类别。
