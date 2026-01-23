<p align="center">
  <img src="./images/logo.jpg" alt="Quanta Authentication" width="200" />
</p>

<h1 align="center">Quanta Authentication</h1>

<p align="center">
  <strong>实验室级身份认证与授权服务系统</strong>
</p>

<p align="center">
  基于 OAuth 2.0 / OpenID Connect 标准 · Go + Vue 3 · Monorepo 架构
</p>

---

> **内部项目** - 本项目为实验室内部使用，仅限授权人员访问。

## 项目简介

**Quanta Authentication (QAuth)** 是一套面向组织的统一身份认证服务系统，完整实现了 OAuth 2.0 和 OpenID Connect 协议，提供用户身份管理、单点登录（SSO）以及基于角色的访问控制（RBAC）能力。

## 技术栈

| 服务 | 技术 |
|:---|:---|
| **认证服务** | Go 1.25 · Gin · GORM · PostgreSQL 16 · Redis |
| **管理后台** | Vue 3.5 · TypeScript · PrimeVue · TailwindCSS 4 |
| **构建工具** | Turborepo · pnpm · Rolldown Vite |
| **部署** | Docker Compose · Nginx |

## 项目结构

```
quanta-authentication/
├── packages/
│   ├── qauth-server/          # Go 认证服务
│   └── qauth-admin-web/       # Vue 管理后台
├── docker/                    # Docker 配置
└── nginx/                     # Nginx 网关配置
```

## 快速开始

### 前置要求

- Node.js >= 18
- Go >= 1.25
- pnpm >= 9.0
- Docker & Docker Compose

### 使用 Docker Compose（推荐）

```bash
cd docker
docker-compose up -d
```

### 本地开发

```bash
# 安装依赖
pnpm install

# 启动所有服务
pnpm dev
```

### 访问地址

| 服务 | 地址 |
|:---|:---|
| 管理后台 | http://localhost |
| PostgreSQL | localhost:15432 |
| Redis | localhost:16379 |

---

<p align="center">
  <sub>Quanta Authentication · 实验室内部项目</sub>
</p>
