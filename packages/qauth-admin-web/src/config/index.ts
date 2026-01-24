import type { ThemeColor } from '@/types'

// 应用配置
export const APP_NAME = 'Quanta Auth'
export const DARK_MODE_SELECTOR = '.app-dark'

// 存储键
export const STORAGE_KEYS = {
   THEME: 'theme',
   RECENT_SEARCHES: 'qauth-recent-searches',
} as const

// 搜索配置
export const MAX_RECENT_SEARCHES = 5

// 动画配置
export const ANIMATION_DEFAULTS = {
   DURATION: 800,
   DECIMALS: 0,
   AUTO_START: true,
   DELAY: 0,
   PAD_START: false,
   USE_GROUPING: true,
} as const

// 缓动函数
export const easings = {
   linear: (t: number) => t,
   easeOutExpo: (t: number) => (t === 1 ? 1 : 1 - Math.pow(2, -10 * t)),
   easeOutCubic: (t: number) => 1 - Math.pow(1 - t, 3),
   easeOutQuart: (t: number) => 1 - Math.pow(1 - t, 4),
   easeOutElastic: (t: number) => {
      const c4 = (2 * Math.PI) / 3
      return t === 0 ? 0 : t === 1 ? 1 : Math.pow(2, -10 * t) * Math.sin((t * 10 - 0.75) * c4) + 1
   },
   easeOutBack: (t: number) => {
      const c1 = 1.70158
      const c3 = c1 + 1
      return 1 + c3 * Math.pow(t - 1, 3) + c1 * Math.pow(t - 1, 2)
   },
} as const

// 主题配色
interface ThemeColorConfig {
   bgLight: string
   textLight: string
   badgeLight: string
   lineLight: string
   bgDark: string
   glowDark: string
   textDark: string
   badgeDark: string
   lineDark: string
   accentDark: string
}

export const THEME_COLORS: Record<ThemeColor, ThemeColorConfig> = {
   orange: {
      bgLight: 'linear-gradient(145deg, #fff7ed 0%, #fed7aa 30%, #fdba74 80%, #fff7ed 120%)',
      textLight: '#9a3412',
      badgeLight: 'rgba(251, 146, 60, 0.25)',
      lineLight: '#ea580c',
      bgDark: 'linear-gradient(145deg, #1c1917 0%, #292524 50%, #44403c 100%)',
      glowDark: 'inset 0 1px 0 rgba(251, 146, 60, 0.15)',
      textDark: '#fdba74',
      badgeDark: 'rgba(251, 146, 60, 0.2)',
      lineDark: '#fb923c',
      accentDark: '#f97316',
   },
   blue: {
      bgLight: 'linear-gradient(145deg, #eff6ff 0%, #bfdbfe 30%, #93c5fd 80%, #eff6ff 120%)',
      textLight: '#1e40af',
      badgeLight: 'rgba(59, 130, 246, 0.2)',
      lineLight: '#2563eb',
      bgDark: 'linear-gradient(145deg, #0c1929 0%, #1e293b 50%, #334155 100%)',
      glowDark: 'inset 0 1px 0 rgba(59, 130, 246, 0.2)',
      textDark: '#93c5fd',
      badgeDark: 'rgba(59, 130, 246, 0.25)',
      lineDark: '#60a5fa',
      accentDark: '#3b82f6',
   },
   green: {
      bgLight: 'linear-gradient(145deg, #f0fdf4 0%, #bbf7d0 30%, #86efac 80%, #f0fdf4 120%)',
      textLight: '#166534',
      badgeLight: 'rgba(34, 197, 94, 0.2)',
      lineLight: '#16a34a',
      bgDark: 'linear-gradient(145deg, #052e16 0%, #14532d 50%, #166534 100%)',
      glowDark: 'inset 0 1px 0 rgba(34, 197, 94, 0.2)',
      textDark: '#86efac',
      badgeDark: 'rgba(34, 197, 94, 0.25)',
      lineDark: '#4ade80',
      accentDark: '#22c55e',
   },
   purple: {
      bgLight: 'linear-gradient(145deg, #faf5ff 0%, #e9d5ff 30%, #d8b4fe 80%, #faf5ff 120%)',
      textLight: '#7e22ce',
      badgeLight: 'rgba(168, 85, 247, 0.2)',
      lineLight: '#9333ea',
      bgDark: 'linear-gradient(145deg, #1a0a2e 0%, #2e1065 50%, #4c1d95 100%)',
      glowDark: 'inset 0 1px 0 rgba(168, 85, 247, 0.2)',
      textDark: '#d8b4fe',
      badgeDark: 'rgba(168, 85, 247, 0.25)',
      lineDark: '#c084fc',
      accentDark: '#a855f7',
   },
   cyan: {
      bgLight: 'linear-gradient(145deg, #ecfeff 0%, #a5f3fc 30%, #67e8f9 80%, #ecfeff 120%)',
      textLight: '#0e7490',
      badgeLight: 'rgba(6, 182, 212, 0.2)',
      lineLight: '#06b6d4',
      bgDark: 'linear-gradient(145deg, #042f2e 0%, #134e4a 50%, #115e59 100%)',
      glowDark: 'inset 0 1px 0 rgba(6, 182, 212, 0.25)',
      textDark: '#67e8f9',
      badgeDark: 'rgba(6, 182, 212, 0.25)',
      lineDark: '#22d3ee',
      accentDark: '#06b6d4',
   },
   red: {
      bgLight: 'linear-gradient(145deg, #fef2f2 0%, #fecaca 30%, #fca5a5 80%, #fef2f2 120%)',
      textLight: '#b91c1c',
      badgeLight: 'rgba(239, 68, 68, 0.2)',
      lineLight: '#dc2626',
      bgDark: 'linear-gradient(145deg, #1c0a0a 0%, #450a0a 50%, #7f1d1d 100%)',
      glowDark: 'inset 0 1px 0 rgba(239, 68, 68, 0.15)',
      textDark: '#fca5a5',
      badgeDark: 'rgba(239, 68, 68, 0.25)',
      lineDark: '#f87171',
      accentDark: '#ef4444',
   },
}

// SimpleStatCard 主题配色
export const SIMPLE_STAT_COLORS = {
   orange: {
      color: 'var(--p-orange-500)',
      bg: 'var(--p-orange-50)',
      bgDark: 'rgba(251, 146, 60, 0.15)',
   },
   blue: { color: 'var(--p-blue-500)', bg: 'var(--p-blue-50)', bgDark: 'rgba(59, 130, 246, 0.15)' },
   green: {
      color: 'var(--p-green-500)',
      bg: 'var(--p-green-50)',
      bgDark: 'rgba(34, 197, 94, 0.15)',
   },
   purple: {
      color: 'var(--p-purple-500)',
      bg: 'var(--p-purple-50)',
      bgDark: 'rgba(168, 85, 247, 0.15)',
   },
   cyan: { color: 'var(--p-cyan-500)', bg: 'var(--p-cyan-50)', bgDark: 'rgba(6, 182, 212, 0.15)' },
   red: { color: 'var(--p-red-500)', bg: 'var(--p-red-50)', bgDark: 'rgba(239, 68, 68, 0.15)' },
   gray: {
      color: 'var(--p-surface-500)',
      bg: 'var(--p-surface-100)',
      bgDark: 'rgba(107, 114, 128, 0.15)',
   },
} as const

// 通知类型配置
export const NOTIFICATION_TYPE_CONFIG = {
   system: {
      label: '系统',
      icon: 'pi pi-cog',
      color: 'text-blue-600 dark:text-blue-400',
      bgColor: 'bg-blue-100 dark:bg-blue-900/30',
   },
   security: {
      label: '安全',
      icon: 'pi pi-shield',
      color: 'text-green-600 dark:text-green-400',
      bgColor: 'bg-green-100 dark:bg-green-900/30',
   },
   user: {
      label: '用户',
      icon: 'pi pi-user',
      color: 'text-purple-600 dark:text-purple-400',
      bgColor: 'bg-purple-100 dark:bg-purple-900/30',
   },
   oauth: {
      label: 'OAuth',
      icon: 'pi pi-key',
      color: 'text-cyan-600 dark:text-cyan-400',
      bgColor: 'bg-cyan-100 dark:bg-cyan-900/30',
   },
   alert: {
      label: '警告',
      icon: 'pi pi-exclamation-triangle',
      color: 'text-red-600 dark:text-red-400',
      bgColor: 'bg-red-100 dark:bg-red-900/30',
   },
} as const

// 搜索类别配置
export const SEARCH_CATEGORY_CONFIG = {
   navigation: { icon: 'pi pi-compass', color: 'text-blue-500 bg-blue-50 dark:bg-blue-900/20' },
   action: { icon: 'pi pi-bolt', color: 'text-amber-500 bg-amber-50 dark:bg-amber-900/20' },
   user: { icon: 'pi pi-users', color: 'text-emerald-500 bg-emerald-50 dark:bg-emerald-900/20' },
   app: { icon: 'pi pi-box', color: 'text-violet-500 bg-violet-50 dark:bg-violet-900/20' },
   recent: { icon: 'pi pi-clock', color: 'text-surface-500 bg-surface-100 dark:bg-surface-700' },
} as const

// 审计模块图标
export const AUDIT_MODULE_ICONS: Record<string, string> = {
   AUTH: 'pi pi-sign-in',
   OAUTH: 'pi pi-key',
   USER: 'pi pi-users',
   ROLE: 'pi pi-shield',
   PERMISSION: 'pi pi-lock',
   CLIENT: 'pi pi-box',
   SYSTEM: 'pi pi-cog',
}

// 用户角色选项
export const USER_ROLE_OPTIONS = [
   { label: '管理员', value: '管理员' },
   { label: '开发者', value: '开发者' },
   { label: '普通用户', value: '普通用户' },
] as const

// OAuth Scope 选项
export const OAUTH_SCOPE_OPTIONS = [
   { label: 'OpenID', value: 'openid' },
   { label: 'Profile', value: 'profile' },
   { label: 'Email', value: 'email' },
   { label: 'Admin', value: 'admin' },
   { label: 'Read Users', value: 'read:users' },
   { label: 'Write Users', value: 'write:users' },
   { label: 'Offline Access', value: 'offline_access' },
] as const

// 设置分组
export const SETTING_GROUPS = [
   { id: 'general', label: '基本设置', icon: 'pi pi-cog' },
   { id: 'security', label: '安全设置', icon: 'pi pi-shield' },
   { id: 'oauth', label: 'OAuth 配置', icon: 'pi pi-key' },
   { id: 'email', label: '邮件服务', icon: 'pi pi-envelope' },
   { id: 'storage', label: '存储配置', icon: 'pi pi-database' },
] as const

// 语言选项
export const LANGUAGE_OPTIONS = [
   { label: '简体中文', value: 'zh-CN' },
   { label: '繁体中文', value: 'zh-TW' },
   { label: 'English', value: 'en-US' },
   { label: '日本語', value: 'ja-JP' },
] as const

// 时区选项
export const TIMEZONE_OPTIONS = [
   { label: '(UTC+8) 中国标准时间', value: 'Asia/Shanghai' },
   { label: '(UTC+9) 日本标准时间', value: 'Asia/Tokyo' },
   { label: '(UTC+0) 格林威治时间', value: 'UTC' },
   { label: '(UTC-5) 美国东部时间', value: 'America/New_York' },
] as const

// 加密选项
export const ENCRYPTION_OPTIONS = [
   { label: 'TLS', value: 'tls' },
   { label: 'SSL', value: 'ssl' },
   { label: '无加密', value: 'none' },
] as const

// 存储类型选项
export const STORAGE_TYPE_OPTIONS = [
   { label: '本地存储', value: 'local' },
   { label: 'Amazon S3', value: 's3' },
   { label: '阿里云 OSS', value: 'oss' },
   { label: '腾讯云 COS', value: 'cos' },
] as const

// 职级选项
export const ORG_CLASS_OPTIONS = [
   { label: '高管层', value: '高管层' },
   { label: '管理层', value: '管理层' },
   { label: '员工', value: '员工' },
] as const

// 职级颜色
export const ORG_CLASS_COLORS: Record<string, string> = {
   高管层: 'bg-orange-100 text-orange-700 dark:bg-orange-900/30 dark:text-orange-400',
   管理层: 'bg-purple-100 text-purple-700 dark:bg-purple-900/30 dark:text-purple-400',
   员工: 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400',
}

// 默认设置值
export const DEFAULT_SETTINGS = {
   general: {
      siteName: 'Quanta 认证中心',
      siteDescription: '企业级统一身份认证平台',
      adminEmail: 'admin@example.com',
      defaultLanguage: 'zh-CN',
      timezone: 'Asia/Shanghai',
      maintenanceMode: false,
      allowRegistration: true,
      requireEmailVerification: true,
   },
   security: {
      passwordMinLength: 8,
      maxLoginAttempts: 5,
      lockoutDuration: 30,
      sessionTimeout: 60,
   },
   oauth: {
      accessTokenLifetime: 3600,
      refreshTokenLifetime: 604800,
      authCodeLifetime: 600,
      jwksRotationDays: 90,
   },
   email: {
      smtpPort: 587,
      smtpEncryption: 'tls',
   },
   storage: {
      maxFileSize: 10,
      allowedFileTypes: 'jpg,png,gif,pdf,doc,docx',
   },
} as const
