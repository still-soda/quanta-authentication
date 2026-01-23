# QAuth Admin Web

基于 Vue 3 的认证系统管理后台。

## 技术栈

| 类别 | 技术 |
|:---|:---|
| 框架 | Vue 3.5 |
| 语言 | TypeScript |
| 构建工具 | Rolldown Vite |
| UI 组件 | PrimeVue 4.5 |
| 样式 | TailwindCSS 4 |
| 状态管理 | Pinia |
| 数据请求 | TanStack Vue Query + Axios |

## 功能模块

- **仪表盘** - 系统概览与数据统计
- **用户管理** - 用户列表、创建、编辑
- **角色权限** - 角色管理与权限分配
- **OAuth 应用** - OAuth2 客户端管理
- **组织架构** - 组织结构管理
- **审计日志** - 操作日志查询
- **系统设置** - 系统配置管理

## 项目结构

```
qauth-admin-web/
├── src/
│   ├── apis/                  # API 接口
│   ├── components/            # 组件
│   │   ├── dashboard/         # 仪表盘组件
│   │   ├── oauth/             # OAuth 组件
│   │   ├── roles/             # 角色组件
│   │   ├── users/             # 用户组件
│   │   └── shared/            # 公共组件
│   ├── composables/           # 组合式函数
│   ├── config/                # 应用配置
│   ├── layouts/               # 布局组件
│   ├── router/                # 路由配置
│   ├── stores/                # Pinia 状态
│   ├── types/                 # 类型定义
│   └── views/                 # 页面视图
├── public/                    # 静态资源
├── nginx/                     # Nginx 配置
├── vite.config.ts             # Vite 配置
└── tsconfig.json              # TypeScript 配置
```

## 快速开始

### 安装依赖

```bash
pnpm install
```

### 开发模式

```bash
pnpm dev
```

### 构建生产版本

```bash
pnpm build
```

### 预览构建结果

```bash
pnpm preview
```
